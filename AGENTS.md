# extend-vivox-authorization-service-go

An Extend Service Extension app for **Vivox authorization** written in Go. Generates Vivox access tokens for voice chat integration, exposed as REST endpoints through AGS's API gateway.

This is a template project — clone it, replace the sample logic in the service implementation, and deploy.

## Build & Test

```bash
make build                           # Build the project
go test ./...                        # Run unit tests
docker compose up --build            # Run locally with Docker
make proto                           # Regenerate proto code
```

Linting: `golangci-lint run` (config in `.golangci.yml`).

## Architecture

Game clients reach this app through AGS via auto-generated REST endpoints:

```
Game Client → AGS Gateway → [REST] → gRPC-Gateway → [gRPC] → This App
```

The proto file defines both the gRPC service and the REST mapping (via `google.api.http` annotations). The gRPC-Gateway automatically generates an OpenAPI spec and REST proxy from the proto.

The sample implementation generates HS256-signed JWT tokens for Vivox operations (login, join, join-muted, kick) with configurable issuer, domain, and signing key.

### Key Files

| Path | Purpose |
|---|---|
| `main.go` | Entry point — starts gRPC server, wires interceptors and observability |
| `pkg/service/myService.go` | **Service implementation** — your custom logic goes here |
| `pkg/service/vivoxToken.go` | **Service implementation** — your custom logic goes here |
| `pkg/proto/permission.proto` | gRPC service definition (user-defined, add your endpoints here) |
| `pkg/proto/service.proto` | gRPC service definition (user-defined, add your endpoints here) |
| `pkg/pb/` | Generated code from proto (do not hand-edit) |
| `pkg/common/` | Auth interceptor, tracing, logging utilities |
| `docker-compose.yaml` | Local development setup |
| `.env.template` | Environment variable template |

### Vivox-Specific Notes

Vivox tokens are HS256-signed JWTs with claims specific to each operation type (login, join, join-muted, kick). The signing key, issuer, and domain are configured via `VIVOX_ISSUER`, `VIVOX_DOMAIN`, and `VIVOX_SIGNING_KEY` environment variables. Each token type encodes a SIP-format URI that identifies the user and/or channel. Token expiry defaults to 90 seconds and is configurable via `VIVOX_DEFAULT_EXPIRY`.

## Rules

See `.agents/rules/` for coding conventions, commit standards, and proto file policies.

## Environment

Copy `.env.template` to `.env` and fill in your credentials.

| Variable | Description |
|---|---|
| `AB_BASE_URL` | AccelByte base URL (e.g. `https://test.accelbyte.io`) |
| `AB_NAMESPACE` | Target namespace |
| `AB_CLIENT_ID` | OAuth client ID |
| `AB_CLIENT_SECRET` | OAuth client secret |
| `PLUGIN_GRPC_SERVER_AUTH_ENABLED` | Enable gRPC auth (`true` by default) |
| `BASE_PATH` | Custom base path for REST endpoints |
| `VIVOX_ISSUER` | Vivox issuer ID |
| `VIVOX_DOMAIN` | Vivox domain |
| `VIVOX_SIGNING_KEY` | Vivox HS256 signing key |

## Dependencies

- [AccelByte Go SDK](https://github.com/AccelByte/accelbyte-go-sdk) — AGS platform SDK and gRPC plugin utilities
