# Jump Backend Monorepo

This monorepo powers Jump's financial services platform for French freelancers and contractors. Business logic is heavily focused on French employment law compliance, salary calculations (portage salarial / coopérative), and financial operations.

## Development Commands

### Build and Quality Tools

- `task lint` — Run golangci-lint (config: `.golangci.custom.yaml`)
  - Can be run directly in modified folders to save time
- `task mocks` — Generate mocks with mockery (config: `.mockery.yml`)
  - Run from project root — it generates all mocks regardless
- `task openapi` — Generate OpenAPI v3 specs with swag
  - Requires `main.go` in current directory
  - **For grouped services**, run from the sub-service directory (e.g. `services/assignments/assignments/`), not the parent
  - **MANDATORY**: Must be run whenever an endpoint is added, removed, or its swaggo comments are updated — run it in the service directory before finishing any task that touches HTTP handlers
  - The generated `specs/v3.yml` must be included in the commit

### Testing

- `go test ./...` — run all tests
- `go test ./services/<service-name>/...` — run tests for a specific service

## Architecture Overview

### Service Organization

Microservices monorepo with three types of services in `services/`:

- **API Services**: HTTP REST APIs (e.g., `users/`, `invoices/`) — deployed with the `backend-service` Helm chart
- **Daemon Services** (legacy): Long-running background workers with internal cron scheduling (e.g., `payroll-daemon/`, `banking-daemon/`) — deployed with the `backend-service` Helm chart. Migration to cronjob services is in progress.
- **Cronjob Services**: Scheduled jobs managed by Kubernetes CronJobs (e.g., `invoices-cronjobs/`, `payroll-cronjobs/`) — deployed with the `backend-cronjob` Helm chart. Preferred over daemons for scheduled work.

Some services are standalone (flat structure), while others are grouped under a parent directory when they share a domain. For example:

- `services/missions/` contains `missions/` (API), `missions-cronjobs/` (cronjob)
- `services/payroll/` contains `payroll-api/`, `payroll-daemon/`, and `payroll-cronjobs/`
- `services/onboarding/` contains `onboarding/` (API) and `onboarding-cronjobs/`
- `services/assignments/` contains `assignments/` (API) and `assignments-cronjobs/`

**Naming rule for grouped services**: the API sub-directory is named exactly like the parent (e.g. `services/missions/missions/`). Do **not** use names like `missions-api` because `go:embed` does not allow `../` paths, so the binary must live in a sibling directory of `internal/` and `pkg/`.

### CI Workflows for grouped services

Each sub-service gets its own deployment workflow in `.github/workflows/`:

```
.github/workflows/
├── <parent>.deployment.yml           # e.g. missions.deployment.yml
├── <parent>-cronjobs.deployment.yml  # e.g. missions-cronjobs.deployment.yml
```

**Important**:
- `service_path` is `./services/<parent>/<sub-service>` (e.g. `./services/missions/missions`)
- `service_name` is the deployed binary name (e.g. `missions`)
- Tags follow the pattern `<service_name>_*.*.*`
- The `paths` trigger **must** include both the sub-service directory and the shared `internal/` sibling so that changes to shared domain code redeploy every dependent binary

### Cronjob Services

Cronjob services are command-line binaries invoked by Kubernetes CronJobs. They follow the same grouped/flat structure as API services but with a few key differences:

- **Helm chart**: `backend-cronjob` (not `backend-service`)
- **CI**: `migration: false` — no migration image is built
- **No HTTP layer**: No `handlers/`, no `jgin.Engine`, no `specs/v3.yml`
- **M2M auth is still required**: If the cronjob calls other internal services, it needs `AuthURL` and `AuthClientSecret` in its `configuration.go` to build an `httpauth.NewAuth` token
- **`kube/values.yaml`**: Declares `cronjobs:` with `name`, `schedule`, and `args`
- **`.gitignore`**: Must exclude the binary at its exact path (e.g. `assignments/assignments` and `assignments-cronjobs/assignments-cronjobs`)

### Standard Service Structure

```
services/<service-name>/
├── main.go                 # Entry point
├── configuration.go        # Service config
├── internal/
│   ├── domain/           # Business entities & interfaces
│   ├── handlers/         # HTTP handlers
│   ├── usecases/         # Business logic
│   ├── infrastructure/   # External integrations
│   └── mocks/           # Generated test mocks
├── pkg/                  # Public API types & clients
├── db/                   # Database migrations (Sqitch)
├── kube/                 # Kubernetes manifests
└── specs/v3.yml         # OpenAPI specification
```

### Shared Libraries (`common/`)

- `jgin/` — HTTP server with auth, tracing, monitoring middleware
- `jdb/` — Database utilities and query builders
- `jmb/` — Message broker (SQS) wrapper
- `jsimulator/` — French salary calculation engine (coop/portage)
- `clients/` — External service clients (Auth0, HubSpot, Swan, etc.)
- `jutils/` — Business validations (IBAN, SIRET, phone numbers)

### Common Patterns and Pitfalls

- **`go:embed` paths**: The `//go:embed` directive does not allow `../` relative paths. This is why grouped services put the API binary in a sibling directory of `internal/` (e.g. `services/missions/missions/`), not `missions-api/`.
- **Multi-value filters**: To accept multiple filter values in a query parameter (e.g. `?status=pending&status=ended`), use a slice with `form:"status"` in the filter struct and pass `pq.Array(...)` to the SQL query.
- **Cronjobs iterating entities**: When a cronjob loops over stored entities that may have optional foreign keys (e.g. `ContractID *uuid.UUID`), always check for `nil` before dereferencing — otherwise the job will panic on newly created records.
- **Time-based business logic**: Inject `jclock.Clock` into usecases so tests can use `jclock.NewFakeClock(...)` to deterministically test date-based transitions (e.g. status changes based on `start_date`/`end_date`).
- **Detecting a 404 from an HTTP client**: use `jhttp.IsNotFoundAPIError(err)` (and `jhttp.IsHTTPAPIError(err, status)` for other codes) instead of manually comparing `errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound`.
- **`pkg/` vs `internal/`**: `pkg/` is reserved for the service's public cross-service API. A helper used only within the service (e.g. a category-range predicate) belongs in `internal/domain`.

### Architecture Patterns

- **Clean Architecture**: Domain → Usecases → Infrastructure separation
- **Dependency Injection**: Interface-based with extensive mocking
- **Event-Driven**: Async communication via SQS message queues
- **Shared Kernel**: Common business logic centralized in `common/`

## Workflow

### Linear

- Issues related to this repository should be labeled **Back-end task**.
- Issue titles are prefixed with the target service name: `[service-name] Description`.
  - Multiple services can be indicated with `+` or `/`: `[onboarding+dashboard]`, `[missions/pandadoc]`.
- Sub-issues are used regularly to break down work.
- **Projects**:
  - **Tech Backlog** — maintenance, dependency updates, bug fixes, and tech debt.
  - **Side projects** — small, well-defined user stories coming from product.
  - Larger features get their own dedicated project that runs over the course of the year.
- **Creating Linear issues**:
  - If work is already in progress (code changes exist), assign the issue to the current user and set status to **Implem in progress**.
  - If work hasn't started yet, don't assign the issue but set status to **Ready for implem** (unless otherwise specified).
  - Default to these statuses unless the user explicitly requests a different status.
- **Issue descriptions**:
  - Write descriptions **prescriptively** (what needs to be done), not retrospectively (what was done).
  - Use imperative mood: "Add X", "Create Y", "Remove Z" — not "Added X", "Created Y", "Removed Z".
  - Structure with clear sections (objective, tasks, constraints/context).
  - Use proper Markdown formatting with actual line breaks.

### GitHub

- **Branch naming**: Free-form. Branches typically include the Linear ticket ID (e.g., `feature/pt-5078-...`, `feat/pt-5074-fees-plan-spec`, `fix/...`).
- **PR titles**: Use the service name as prefix. For service-related changes: `<service-name>: <description>` (e.g., `assignments: add status filter to ListAssignments`).
  - For changes that don't touch a specific service (CI, docs, repo-wide tooling, AI skills, etc.), keep conventional prefixes: `docs:`, `ci:`, `ai:`, ...
  - Examples:
    - `assignments: add assignment statuses and UpdateStatuses cronjob`
    - `contracts: add GetContract client method`
    - `docs: update AGENTS.md with commit conventions`
    - `ci: add assignments-cronjobs deployment workflow`
- **PR descriptions**: Contains an API impact recap — list new APIs and the changes made to existing ones (mark each as breaking or non-breaking). Do not duplicate the full Linear issue content.

## Testing Guidelines

**MANDATORY**: You **MUST** use the `/test` skill when writing any tests in this codebase. This skill contains all project-specific testing conventions and patterns.

```bash
# To use:
Use the skill tool to load the "test" skill, then follow its instructions.
```

Full documentation: `.claude/skills/test/SKILL.md`

**HTTP client tests**: Always use `common/jhttp/httpmock` — **never** `httptest.NewServer`. See `common/clients/api_gouv/geo/client_test.go` for the canonical pattern.

**MANDATORY**: Whenever you add or modify a method in any layer (handler, usecase, store/infrastructure), you **MUST** write or update the corresponding test in the same task — without waiting to be asked. This applies to:
- New store methods (e.g., a new SQL query function)
- New or modified usecase methods
- New or modified HTTP handlers

## Code Comments

**MANDATORY**: Do not add comments that just restate what the code already says (e.g. explaining what a regex matches, what a struct groups together, what an if-condition checks). Only comment on the non-obvious *why* (a workaround, a business rule, a surprising constraint). Prefer no comment over a restating one.

## Error Wrapping

**MANDATORY**: In all non-pure methods (HTTP clients, handlers, usecases, stores, and any method that calls external dependencies), always use `defer jerror.Wrap(&err)` for error context — never use a local `errMsg` constant or manual `fmt.Errorf("funcName has failed: %w", err)` wrapping at each return site.

```go
// ✅ DO
func (c *client) CreateClient(ctx context.Context, req pkg.RequestCreateClient) (id uuid.UUID, err error) {
    defer jerror.Wrap(&err)
    // inner errors are wrapped concisely, without repeating the function name:
    return uuid.Nil, fmt.Errorf("json.Marshal err: %w", err)
}

// ❌ DON'T
func (c *client) CreateClient(ctx context.Context, req pkg.RequestCreateClient) (uuid.UUID, error) {
    const errMsg = "client.CreateClient has failed"
    return uuid.Nil, fmt.Errorf("%s: json.Marshal err: %w", errMsg, err)
}
```

`jerror.Wrap` automatically prepends the calling function's name (e.g. `"client.(*client).CreateClient has failed: ..."`) and uses `%w`, so `errors.Is`/`errors.As` chains are preserved.

## Third-party Documentations

- [Calendly API](https://developer.calendly.com/api-docs/d7755e2f9e5fe-calendly-api)
- [HubSpot API](https://developers.hubspot.com/docs/reference/api/overview)
- [Intercom API](https://developers.intercom.com/docs/references/rest-api/api.intercom.io)
- [Kenko API](https://kenko.readme.io/)
- [PandaDoc API](https://developers.pandadoc.com/reference/about)
- [Stripe API](https://docs.stripe.com/api)

## Cursor Cloud specific instructions

This repo is part of the **Jump 4-repo end-to-end Cloud Agent environment**. Cursor clones these four repos as siblings (on this env under `/agent/repos/<repo>`):

- `frontend-monorepo` (default branch `main`)
- `backend-interview` (default branch `master`)
- `application-infrastructure` (default branch `main`)
- `terraform-modules` (default branch `main`)

One feature often spans several of them (e.g. a backend API/service here, service/API permissions in `application-infrastructure`, emails in `terraform-modules`, matching UI in `frontend-monorepo`).

**Workflow**
- Open **one PR per repo you change**, and cross-link the PRs in their descriptions.
- Base each PR on that repo's own default branch (do not force a shared branch name).
- Never merge, force-push, or use a personal GitHub PAT — GitHub access comes from the Cursor GitHub App / environment token.
- Do not commit secrets, `.env`, or credentials.

**Install / test here**
- Toolchain is installed by the environment; if needed: `command -v task || go install github.com/go-task/task/v3/cmd/task@v3.44.1`.
- Build: `go build ./...`. Tests: `task test` (or `go test ./...`). Lint: `task lint`. Mocks: `task mocks`.
- `task openapi` works because `swag` is pre-installed over HTTPS (the Taskfile otherwise clones the fork over SSH, which agents can't use).

**Staging hours (frontend)**
After **20:00 Europe/Paris** and on **weekends**, the dev API is down — run frontend apps with their `staging` scripts (e.g. `pnpm --filter <app> staging`) instead of `dev`.
