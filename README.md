# RustDesk Reviewed Source

Corresponding source for the reviewed Linux amd64 Docker images published on
2026-11-11. The original development repositories remain private; this repository
contains a source snapshot, not their Git history.

## Images

- Integrated server and API:
  `ghcr.io/5120xiaoge/rustdesk-server-s6-reviewed:reviewed-20261111-v1.1`
  (`sha256:720dbe25a0879475a6403d2a8d301f38b291897040bd9f58d969dbf4f51df093`)
- API:
  `ghcr.io/5120xiaoge/rustdesk-api-reviewed:reviewed-20261111-v1.1`
  (`sha256:009023db119d2b625574f19bde144495d790de3ee5d248393e8329bb56f1cf42`)

See [BUILD.md](BUILD.md) for build instructions and [deploy/README.md](deploy/README.md)
for the Compose configuration. The anonymous image verification workflow must pass
before treating the integrated image as anonymously pullable.

## Source Provenance

| Directory | Source | Original Commit |
| --- | --- | --- |
| `server/` | Modified `lejianwen/rustdesk-server`, forapi implementation | `f826b779acc2f09161d95807713f80fa643e0d70` |
| `server/libs/hbb_common/` | Vendored RustDesk hbb_common submodule | `d6b14975ffd35ed63528617a53795c479a6eaf13` |
| `api/` | Modified `lejianwen/rustdesk-api` | `b4c38b5c27c8b5951a6a7bc74efb3707fa2a8290` |
| `web/` | `lejianwen/rustdesk-api-web` | `3998c2a9213fcd047252776d0f0db33e6717026c` |

`Cargo.lock`, `go.sum`, the frontend lockfile, test cases, Dockerfiles, and original
publication workflows are included. The submodule is vendored so a private checkout
is not needed. Remaining dependencies use the included lockfiles and their recorded
public sources.

Upstream example private-key literals in five server README translations have been
replaced by placeholders in this snapshot. No executable source was changed during
snapshot preparation. Root documentation, deployment examples, and the anonymous
verification workflow are additional files, not files used in the original image.

## Review Scope

Focused fixes cover server/API JWT validation, signing-key logging, file-based key
validation, address-book input types, language initialization, and Docker packaging.
Original server and API publication workflows passed regression tests, locked
compilation, and container smoke/restart checks before image publication.

These checks are not an exhaustive security audit or a production deployment test.
Production keys, account databases, credentials, client signing material, and private
repository history are not included.

## Licenses

The modified server is AGPL-3.0; see `server/LICENSE`. The API is MIT; see
`api/LICENSE`. Other components and dependencies retain their included licenses and
notices. This repository does not replace those licenses with a single new license.
