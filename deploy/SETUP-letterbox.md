# Letterbox setup on liu

One-time wiring to expose `letterbox.pointegrity.com` and have it
write to a persistent SQLite file under `/home/liu/infra/`.

## 1. DNS

Add an A record at the DNS provider for `pointegrity.com`:

    letterbox.pointegrity.com  →  <liu's public IP>

(Same target as `pouch.pointegrity.com` etc.)

## 2. Workspace on liu

    sudo -u liu mkdir -p /home/liu/infra/pointegrityws/data

`make deploy` (the target added in pointegrity/Makefile) clones the
repo into `/home/liu/infra/pointegrityws/pointegrity` and builds the
binary into `./build/letterbox`. Mirrors the other apps' layout.

## 3. /etc/default/letterbox

Copy `deploy/letterbox.env.example` to `/etc/default/letterbox` and
generate a real admin key:

    sudo cp deploy/letterbox.env.example /etc/default/letterbox
    sudo openssl rand -base64 32 | sed -i "s|replace-with-32.*|$(cat)|" /etc/default/letterbox
    sudo chmod 600 /etc/default/letterbox
    sudo chown root:root /etc/default/letterbox

## 4. systemd unit

    sudo cp deploy/letterbox.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl enable --now letterbox
    sudo systemctl status letterbox

## 5. Nginx Proxy Manager (UI at port 81)

In the NPM UI, add a new Proxy Host:

- Domain Names: `letterbox.pointegrity.com`
- Scheme: `http`
- Forward Hostname / IP: `172.18.0.1`
- Forward Port: `3737`
- Block common exploits: ON
- Websockets Support: not needed
- SSL tab: request a new Let's Encrypt cert; force SSL; HTTP/2

## 6. Smoke test

    curl -i https://letterbox.pointegrity.com/healthz
    # expect: HTTP/2 200 + body "ok"

    curl -i -X POST -d 'email=you@example.com&source=manual' \
      https://letterbox.pointegrity.com/subscribe
    # expect: 303 See Other → https://www.pointegrity.com/journal/subscribed/

## 7. Pulling the subscriber list

    KEY=$(sudo grep LETTERBOX_ADMIN_KEY /etc/default/letterbox | cut -d= -f2)
    curl -s "https://letterbox.pointegrity.com/list?key=$KEY" > subscribers.csv

Pipe the CSV into your mail client of choice when sending the
monthly letter. (No automated sending in v1 — keep the writing
practice unhurried.)
