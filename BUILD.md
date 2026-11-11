# Build the Reviewed Images

Use a Linux amd64 host with Git, Docker, Go 1.23.10, Node.js 20, Rust 1.99.1,
the `x86_64-unknown-linux-musl` Rust target, `musl-gcc`, and `protoc` available.
The original builds used Ubuntu 24.04 GitHub runners and S6 Overlay 3.2.0.0.
Acquire missing tools from their official sources; review installations separately.
The original build workflows are included under `api/.github/workflows/` and
`server/.github/workflows/`.

Dependency locks are retained. The base Alpine 3.22 image and original reviewed API
base are pinned by digest. S6 archives are checked against their release SHA256 files.
The build is functionally repeatable, not claimed to produce byte-identical binaries
or image digests, since package indexes and other runner components can change.

Run these commands from this repository's root.

## API and Admin Frontend

```bash
set -eu
(
  cd web
  npm ci --ignore-scripts
  npm run build
)
mkdir -p api/resources/admin
cp -a web/dist/. api/resources/admin/
(
  cd api
  go mod verify
  go test -mod=readonly -race ./lib/jwt ./lib/bootstrap ./http/controller/api ./global
  go vet -mod=readonly ./lib/jwt ./lib/bootstrap
  mkdir -p amd64/release/data amd64/release/runtime
  printf '%s\n' reviewed-20261111-v1.2 > resources/version
  GOOS=linux GOARCH=amd64 CC=musl-gcc CGO_ENABLED=1 CGO_LDFLAGS="-static" \
    go build -mod=readonly -ldflags="-s -w" -o amd64/release/apimain ./cmd/apimain.go
  cp -a conf resources docs amd64/release/
  docker build --platform linux/amd64 --build-arg BUILDARCH=amd64 \
    -t rustdesk-api-reviewed:local .
)
```

## Server and S6 Integration

```bash
set -eu
(
  cd server
  export CARGO_TARGET_X86_64_UNKNOWN_LINUX_MUSL_LINKER=musl-gcc
  export CC_x86_64_unknown_linux_musl=musl-gcc
  cargo test --locked --test jwt_tests --test key_file_tests
  cargo build --locked --release --target x86_64-unknown-linux-musl --bins
  install -m 0755 target/x86_64-unknown-linux-musl/release/hbbs docker/rootfs/usr/bin/
  install -m 0755 target/x86_64-unknown-linux-musl/release/hbbr docker/rootfs/usr/bin/
  install -m 0755 target/x86_64-unknown-linux-musl/release/rustdesk-utils docker/rootfs/usr/bin/
  install -D -m 0644 LICENSE docker/rootfs/usr/share/licenses/rustdesk-server/LICENSE
  docker build --platform linux/amd64 -f docker/Dockerfile \
    --build-arg API_IMAGE=rustdesk-api-reviewed:local \
    --build-arg S6_ARCH=x86_64 --build-arg S6_OVERLAY_VERSION=3.2.0.0 \
    -t rustdesk-server-s6-reviewed:local docker
)
```

To build with the exact published API base instead of rebuilding it, replace the
`API_IMAGE` build argument with:

```text
ghcr.io/5120xiaoge/rustdesk-api-reviewed@sha256:3499576cdc034c3e9f280193ad8c9a58d2d260f0a9d53e26a6dacc851b6e082a
```

All commands above are local builds and tests; they do not push images or change a
production server. The anonymous verification workflow uses an isolated temporary
Docker credential directory and an ephemeral test container, never production data.
