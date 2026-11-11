# RustDesk Reviewed Source

Corresponding source for the reviewed Linux amd64 Docker images published on
2026-11-11. The original development repositories remain private; this repository
contains a source snapshot, not their Git history.

## Images

- Integrated server and API:
  `ghcr.io/5120xiaoge/rustdesk-server-s6-reviewed:reviewed-20261111-v1.2`
  (`sha256:bfa791c8cb6662ce157ca4c583df101bdcc44a5d97106d6bd1b3d259b8cf1aae`)
- API:
  `ghcr.io/5120xiaoge/rustdesk-api-reviewed:reviewed-20261111-v1.2`
  (`sha256:3499576cdc034c3e9f280193ad8c9a58d2d260f0a9d53e26a6dacc851b6e082a`)

See [BUILD.md](BUILD.md) for build instructions and [deploy/README.md](deploy/README.md)
for the Compose configuration. The anonymous image verification workflow must pass
before treating the integrated image as anonymously pullable.

## Source Provenance

| Directory | Source | Original Commit |
| --- | --- | --- |
| `server/` | Modified `lejianwen/rustdesk-server`, forapi implementation | `ccfde77610fda5a9c93c6ba0af3567366d515c9a` |
| `server/libs/hbb_common/` | Vendored RustDesk hbb_common submodule | `d6b14975ffd35ed63528617a53795c479a6eaf13` |
| `api/` | Modified `lejianwen/rustdesk-api` | `b6d5a1bd5a9dbee1182d23d87bfc046d891d899f` |
| `web/` | `lejianwen/rustdesk-api-web` | `3998c2a9213fcd047252776d0f0db33e6717026c` |

`Cargo.lock`, `go.sum`, the frontend lockfile, test cases, Dockerfiles, and original
publication workflows are included. The submodule is vendored so a private checkout
is not needed. Remaining dependencies use the included lockfiles and their recorded
public sources.

Upstream example private-key literals in five server README translations have been
replaced by placeholders in this snapshot. No executable source was changed during
snapshot preparation. Root documentation, deployment examples, and the anonymous
verification workflow are additional files, not files used in the original image.

The upstream sample device records have since been removed from
`server/db_v2.sqlite3`. It retains only the schema needed by SQLx compile-time
queries and is not a production database. Deployment examples now require a
`RUSTDESK_PUBLIC_HOST` value instead of embedding a server address. These privacy
changes do not alter the previously published images.

## Review Scope

Focused fixes cover server/API JWT validation, signing-key logging, file-based key
validation, address-book input types, language initialization, and Docker packaging.
Version v1.2 adds persistent automatic JWT initialization for mandatory-login
integrated deployments, private initial administrator credentials, and stdin-based
administrator password reset. Existing account databases are not reset.
Original server and API publication workflows passed regression tests, locked
compilation, and container smoke/restart checks before image publication.
The v1.2 integration workflow additionally verified fresh bind-directory creation,
administrator authentication, shared JWT signing, container recreation, password
reset persistence, credential permissions, and absence of plaintext secret logs.

These checks are not an exhaustive security audit or a production deployment test.
Production keys, account databases, credentials, client signing material, and private
repository history are not included.

## Licenses

The modified server is AGPL-3.0; see `server/LICENSE`. The API is MIT; see
`api/LICENSE`. Other components and dependencies retain their included licenses and
notices. This repository does not replace those licenses with a single new license.
