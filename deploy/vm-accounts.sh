#!/usr/bin/env bash
# Create the five public PA#3 test accounts on the deployment VM.
# Run it yourself on the VM after vm-bootstrap.sh:
#   sudo bash deploy/vm-accounts.sh
# It prints each email and password ONCE so they can go into PA3.md. The
# accounts are public by design (PA#3 rubric): they must hold no real data.
# Re-running fails for accounts that already exist; nothing is overwritten.
set -euo pipefail
APP=/opt/leetkob
[ "$(id -u)" -eq 0 ] || { echo "Run with sudo." >&2; exit 1; }
# shellcheck disable=SC1091
. /etc/leetkob/worker.env

word() { shuf -n1 -e amber birch cedar delta ember fjord glade harbor indigo juniper kestrel lagoon meadow nectar orchid prairie quartz raven sierra tundra; }
for spec in admin:admin ta1:ta ta2:ta student1:student student2:student; do
  name="${spec%%:*}"; role="${spec##*:}"
  email="pa3-${name}@leetkob.test"
  password="$(word)-$(word)-$(word)-$((RANDOM % 90 + 10))"
  DATABASE_URL="$DATABASE_URL" USER_EMAIL="$email" USER_PASSWORD="$password" USER_ROLE="$role" \
    "$APP/bin/create-user" >/dev/null
  echo "LK_${name^^}=${email}:${password}"
done
echo "Next, from your own machine: put these lines in an untracked file and run"
echo "  LK_URL=https://<domain> node scripts/pa3-walkthrough.mjs setup"
