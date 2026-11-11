# Deploy with Docker Compose

This example is for Debian x86_64 and the delivered custom RustDesk 1.5.1 clients:
ID `rust.520kx.dpdns.org:22116`, relay `rust.520kx.dpdns.org:22117`, and API
`https://rust.520kx.dpdns.org`.

## Prerequisites

- The integrated GHCR image must be publicly pullable.
- `/data/rustdesk/server` must be the existing server directory containing the
  matching `id_ed25519` and `id_ed25519.pub` pair, not a newly created empty directory.
- `/data/rustdesk/api` must be the existing API database directory. If either path
  differs, edit the bind source first. Bind mounts do not auto-create missing paths.
- Check the public key locally against the delivered clients. Never send the private
  key. Starting with a different key pair will break those clients.
- The domain must resolve to the server, and firewall rules must allow TCP 22115-22119
  and UDP 22116.
- A host-side HTTPS reverse proxy must forward the domain's API/admin traffic to
  `http://127.0.0.1:21114`. If the proxy runs in Docker, join its network and use the
  service name/port instead; container localhost is not the host.
- Resolve existing container/port conflicts and back up existing data before an
  upgrade. This example does not stop, remove, migrate, or replace existing services.
- A new high-entropy JWT key is required. Do not reuse keys posted in a chat or
  upstream examples. A changed JWT key invalidates existing sessions.

To inspect mounts without dumping environment variables:

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Ports}}'
docker ps -q | xargs -r docker inspect --format \
  '{{.Name}} {{range .Mounts}}{{.Source}} -> {{.Destination}}; {{end}}'
```

## Start

Run from this `deploy/` directory only after checking the prerequisites.
Enter the new JWT key privately at the hidden prompt; do not put its literal value
in a command, the Compose file, shell history, a repository, or chat.

```bash
set -eu
read -rsp 'New private JWT key: ' RUSTDESK_API_JWT_KEY
printf '\n'
export RUSTDESK_API_JWT_KEY
test -s /data/rustdesk/server/id_ed25519
test -s /data/rustdesk/server/id_ed25519.pub
test "$(tr -d '\r\n' < /data/rustdesk/server/id_ed25519.pub)" = \
  '5rfSBmytgCwTBn7VkE5muQHfnRJkybSf1++Muvz4fEM='
test -d /data/rustdesk/api
docker compose config --quiet
docker compose pull
docker compose up -d
docker compose ps
unset RUSTDESK_API_JWT_KEY
```

Only `config --quiet` is used: ordinary `docker compose config` prints interpolated
secrets. The key is stored in container environment metadata during operation;
restrict Docker/root access and never share full inspection output. After leaving
the shell, re-enter the same privately retained key for future Compose updates.

The API retains accounts only when its existing database is reused. A first-time
database initializes an administrator differently and is not the upgrade path here.
This repository does not publish passwords or provide unsafe password-in-argument
reset commands.
