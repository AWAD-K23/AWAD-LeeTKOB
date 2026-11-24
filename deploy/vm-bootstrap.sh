#!/usr/bin/env bash
# Set up or update the public LeeTKOB deployment on one Ubuntu 24.04 VM.
# Run it yourself on the VM (it changes a production host):
#   sudo DOMAIN=your.host.name bash deploy/vm-bootstrap.sh          # first install
#   sudo DOMAIN=your.host.name REF=pa3 bash deploy/vm-bootstrap.sh  # move to the pa3 tag
# Safe to re-run: it never deletes volumes, never overwrites existing .env files,
# and never prints passwords or the model key.
set -euo pipefail

: "${DOMAIN:?Set DOMAIN to the public DNS name of the VM}"
REPO="${REPO:-https://github.com/AWAD-K23/AWAD-LeeTKOB.git}"
REF="${REF:-main}"
APP=/opt/leetkob
ETC=/etc/leetkob
[ "$(id -u)" -eq 0 ] || { echo "Run with sudo." >&2; exit 1; }

echo "== 1. Packages (Docker Engine from Docker's apt repository)"
if ! command -v docker >/dev/null; then
  apt-get update
  apt-get install -y ca-certificates curl git
  install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
  chmod a+r /etc/apt/keyrings/docker.asc
  # shellcheck disable=SC1091
  . /etc/os-release
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu ${VERSION_CODENAME} stable" \
    > /etc/apt/sources.list.d/docker.list
  apt-get update
  apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
fi
command -v git >/dev/null || apt-get install -y git

echo "== 2. Service account for the worker (the only process that talks to Docker)"
id leetkob >/dev/null 2>&1 || useradd --system --home-dir /var/lib/leetkob --create-home --shell /usr/sbin/nologin leetkob
usermod -aG docker leetkob

echo "== 3. Code at ${REF}"
[ -d "$APP/.git" ] || git clone "$REPO" "$APP"
git -C "$APP" fetch --tags --force origin
if git -C "$APP" rev-parse -q --verify "refs/tags/$REF" >/dev/null; then
  git -C "$APP" checkout --detach "refs/tags/$REF"
else
  git -C "$APP" checkout --detach "origin/$REF"
fi
COMMIT="$(git -C "$APP" rev-parse --short HEAD)"
echo "   deploying $COMMIT"

echo "== 4. Compose settings ($APP/.env: database and web only, never the model key)"
if [ ! -f "$APP/.env" ]; then
  DB_PASSWORD="$(openssl rand -hex 24)"
  cat > "$APP/.env" <<EOF
POSTGRES_DB=leetkob
POSTGRES_USER=leetkob
POSTGRES_PASSWORD=${DB_PASSWORD}
APP_ORIGIN=https://${DOMAIN}
COOKIE_SECURE=true
PUBLIC_ADDRESS=${DOMAIN}
WEB_PORT=80
WEB_TLS_PORT=443
EOF
fi
# The worker starts in $APP and reads this file, so its group may read it.
chown root:leetkob "$APP/.env"
chmod 640 "$APP/.env"

echo "== 5. Worker settings ($ETC/worker.env: the only place the model key lives)"
install -d -m 750 -o root -g leetkob "$ETC"
if [ ! -f "$ETC/worker.env" ]; then
  # shellcheck disable=SC1091
  . "$APP/.env"
  cat > "$ETC/worker.env" <<EOF
DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5433/${POSTGRES_DB}?sslmode=disable
SANDBOX_IMAGE=leetkob-sandbox:pa3
# Keep 1 until the provider spending cap is confirmed; then set 0 and fill the key with: sudoedit $ETC/worker.env
MODEL_OFF=1
MODEL_NAME=gemini-3.5-flash-lite
MODEL_API_KEY=
MODEL_LAB_BUDGET_MICRO_USD=1000000
EOF
fi
chown root:leetkob "$ETC/worker.env"
chmod 640 "$ETC/worker.env"

echo "== 6. Sandbox image and the web/API/database stack"
docker build -t leetkob-sandbox:pa3 "$APP/sandbox"
docker compose --project-name leetkob -f "$APP/compose.pa3.yaml" up -d --build

echo "== 7. Worker and account tool binaries (built in a throwaway Go container)"
install -d -m 755 "$APP/bin"
docker run --rm -e CGO_ENABLED=0 -v "$APP:/src" -w /src/backend golang:1.26 \
  sh -c 'go build -buildvcs=false -o /src/bin/worker ./cmd/worker && go build -buildvcs=false -o /src/bin/create-user ./cmd/create-user'
chmod 755 "$APP/bin/worker" "$APP/bin/create-user"

echo "== 8. Worker service"
install -m 644 "$APP/deploy/leetkob-worker.service" /etc/systemd/system/leetkob-worker.service
systemctl daemon-reload
systemctl enable leetkob-worker >/dev/null
systemctl restart leetkob-worker

echo "== 9. Status"
docker compose --project-name leetkob -f "$APP/compose.pa3.yaml" ps --format '{{.Service}}: {{.State}}'
systemctl is-active leetkob-worker
echo "Deployed $COMMIT at https://${DOMAIN} (Caddy fetches the certificate on the first HTTPS request)."
