# Deploy with Docker Compose

This original-style Compose example is for a new Debian x86_64 server.
Replace `117.24.15.98` with the server's actual public address and adjust published
ports and advertised addresses together. YAML is configuration, not a Bash command.

## Prerequisites

- The integrated GHCR image must be publicly pullable.
- Short bind mounts auto-create missing server/API directories on first deployment.
  When both server key files are absent, hbbs creates a new pair. A partial or invalid
  pair remains an error; never delete existing keys to bypass it.
- Existing database accounts and keys are retained when the same directories are
  reused. Back up existing data privately before an upgrade.
- Firewall/security-group rules must allow the selected TCP ports and UDP ID port.
- Plain HTTP is provided to match the upstream example, but exposes login traffic
  on untrusted networks. Use an HTTPS reverse proxy for production API/admin access
  and set the advertised API URL accordingly.
- Resolve existing container/port conflicts and back up existing data before an
  upgrade. This example does not stop, remove, migrate, or replace existing services.
- With `MUST_LOGIN=Y` and no explicit JWT, the image generates a private persistent
  JWT under the API data mount. An explicit JWT environment value still takes
  precedence. Never reuse example or previously disclosed values.

To inspect mounts without dumping environment variables:

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Ports}}'
docker ps -q | xargs -r docker inspect --format \
  '{{.Name}} {{range .Mounts}}{{.Source}} -> {{.Destination}}; {{end}}'
```

## Start

Use the YAML in a Compose management panel, or save it as `docker-compose.yml` and
run from its directory. No manual JWT entry or pre-existing key pair is needed for
a new server.

```bash
docker compose config --quiet
docker compose pull
docker compose up -d --wait --wait-timeout 180
docker compose ps
```

## Private Bootstrap Credentials

First installation creates administrator `admin` and a random initial password.
Inside the container the private files are:

- `/app/data/.bootstrap/jwt.key`
- `/app/data/.bootstrap/admin.password`

On the example host these are under `/data/rustdesk/api/.bootstrap/`.
The directory is mode 0700 and the files are mode 0600. Passwords and JWT keys are
not printed to application logs. Keep Docker/root access private and never upload
these files, include them in a public backup, or paste their contents into a chat.
Inspect the initial password only in your private server terminal when needed.

To set your own administrator password without including it in command arguments,
use a hidden prompt and standard input:

```bash
(
  set -eu
  set +x
  read -rsp 'New administrator password (12-32 bytes): ' ADMIN_PASSWORD </dev/tty
  printf '\n'
  printf '%s\n' "$ADMIN_PASSWORD" |
    docker compose exec -T rustdesk ./apimain reset-admin-pwd --stdin
  unset ADMIN_PASSWORD
)
```

Restart/recreation reuses the JWT and does not reset existing administrators.
The initial password file is an installation record, not an updated copy of a
password subsequently changed through the UI or CLI.

## Clients

Read only the server public key:

```bash
docker compose exec -T rustdesk cat /data/id_ed25519.pub
```

Clients must use this public key and the advertised addresses. Previously delivered
custom clients target a different domain/221xx setup and public key. This deployment
does not rebuild them or prove their embedded settings are GUI-overridable.
