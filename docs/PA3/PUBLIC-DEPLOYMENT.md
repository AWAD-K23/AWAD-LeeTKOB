# Public deployment runbook (PA#3 → PA#5)

One Ubuntu VM runs everything: Caddy (HTTPS, static frontend, `/api` proxy), the Go API,
PostgreSQL, and the **worker as a systemd service on the host**, because only the worker may talk
to Docker to start sandbox containers. The model key lives only in `/etc/leetkob/worker.env`.

```
Internet ──443──▶ Caddy (web) ──▶ API ──▶ PostgreSQL (127.0.0.1:5433, not public)
                                             ▲
                         leetkob-worker (systemd, host) ──docker──▶ sandbox containers
                                 │                                   (no network, no mounts,
                                 └──HTTPS──▶ Gemini (key in /etc/leetkob/worker.env)   no env)
```

The deployment must stay reachable from **Sun 29 Nov 2026** until the PA#4 marks are out (the
PA#4 red team attacks it), and PA#5 needs it again through the oral on 24/31 Dec.

## 1. Create the VM (you, in the cloud console)

- Ubuntu Server **24.04 LTS**, x86-64, **2 vCPU / 4 GB RAM** or more, 30 GB disk.
- Inbound firewall: TCP **22** (your IP only if the console allows it), **80** and **443** from anywhere.
- A DNS name pointing at the VM's public IP. Free options:
  - Azure: the VM's public IP → *Configuration* → **DNS name label** gives `<label>.<region>.cloudapp.azure.com`.
  - Any provider: a free DuckDNS subdomain (`<name>.duckdns.org`) set to the VM's IP.
- Use a **static** public IP, so the address survives a VM restart.

## 2. Install (you, on the VM)

```bash
sudo apt-get update && sudo apt-get install -y git
git clone https://github.com/AWAD-K23/AWAD-LeeTKOB.git /tmp/leetkob-setup
sudo DOMAIN=<your.dns.name> bash /tmp/leetkob-setup/deploy/vm-bootstrap.sh
```

The script installs Docker Engine, creates the `leetkob` service user, clones the repo to
`/opt/leetkob`, writes `/opt/leetkob/.env` (random database password; no model key) and
`/etc/leetkob/worker.env` (**`MODEL_OFF=1`, empty key**), builds the sandbox image, starts
`compose.pa3.yaml`, builds `bin/worker` and `bin/create-user` in a throwaway Go container, and
starts `leetkob-worker`. It never deletes volumes and never overwrites an existing `.env`.

Check: `curl -sI https://<your.dns.name>/` → `200`; `curl -s https://<your.dns.name>/api/me` → `401`.

## 3. Test accounts (you, on the VM)

```bash
sudo bash /opt/leetkob/deploy/vm-accounts.sh
```

It prints five `LK_…=email:password` lines once. They are public by design (the rubric puts them in
`PA3.md`); they hold no real data. Keep a copy in an **untracked** file on your machine, e.g.
`../pa3-accounts.env`, never inside the repository folder.

## 4. Turn the model on (you, after checking the cap)

1. In Google AI Studio, open the **AWAD-LeeTKOB** project, check the monthly spending cap
   (₫100,000 in `img/evidence/gemini-monthly-cap.png`) is still set, and screenshot it with today's date.
2. Create or copy that project's API key. Do not paste it into a chat or a commit.
3. `sudoedit /etc/leetkob/worker.env` → set `MODEL_API_KEY=…` and `MODEL_OFF=0`.
4. `sudo systemctl restart leetkob-worker && sudo journalctl -u leetkob-worker -n 20 --no-pager`.

The key is now readable only by root and the `leetkob` group. The API container, the browser
bundle and the sandbox never receive it (`compose.pa3.yaml` passes no model variable to `api`;
`TestAttacks/secrets` and `TestContainerPolicy` cover the sandbox).

## 5. Walk through it (from your own machine)

```bash
set -a; . ../pa3-accounts.env; set +a
LK_URL=https://<your.dns.name> node scripts/pa3-walkthrough.mjs setup
LK_URL=https://<your.dns.name> node scripts/pa3-walkthrough.mjs walk --out docs/PA3/evidence/public-walk.json
```

`setup` makes class A (TA1 + both students) and class B (TA2). `walk` creates the `sort-merge` lab
as TA1, submits `s01` (compliant), `s04` (violation) and `s09` (injection), waits for judging, prints
TA1's review queue, and checks the refusals (TA2 → 404, student → 403). Then record the same flow
in the browser for the ≤5-minute video.

Queue rule Q1 on the deployment: `sudo systemctl stop leetkob-worker`, submit in the browser
(status stays *Queued*), `sudo systemctl start leetkob-worker`, then count the rows read-only:

```bash
cd /opt/leetkob && sudo docker compose --project-name leetkob -f compose.pa3.yaml exec -T postgres \
  sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT s.id, count(DISTINCT j.submission_id) judged, count(v.*) verdicts FROM submissions s LEFT JOIN judge_results j ON j.submission_id=s.id LEFT JOIN constraint_verdicts v ON v.submission_id=s.id GROUP BY s.id ORDER BY 1"'
```

## 6. Update to the submitted tag

```bash
sudo DOMAIN=<your.dns.name> REF=pa3 bash /opt/leetkob/deploy/vm-bootstrap.sh
```

## Operating notes

- Restart safely: `sudo systemctl restart leetkob-worker`; `docker compose … restart`. **Never**
  `docker compose down -v` (deletes the database volume).
- Logs: `sudo journalctl -u leetkob-worker -f`; `sudo docker compose --project-name leetkob -f /opt/leetkob/compose.pa3.yaml logs -f api web`.
- Cost exposure: the accounts are public, so anyone can create labs and submit. Each lab is capped
  by `MODEL_LAB_BUDGET_MICRO_USD` (default US$1, reserved before each call); the provider's monthly
  project cap is the backstop for many labs. Turn the model off (`MODEL_OFF=1`, restart the worker)
  if the AI Studio usage page shows unexpected spend.
