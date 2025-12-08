# Repository Guidelines

This guide helps contributors work efficiently with the **Muse backend** Go codebase. Follow the conventions below to keep the project consistent, maintainable, and DDD‑compatible.

## Project Structure \& Module Organization
```text
backend/
├─ cmd/                # Application entry point (`main.go`)
├─ internal/           # Domain‑driven design layers
│   ├─ domain/         # Entities, value objects, domain services
│   ├─ application/    # Use‑case / service layer
│   ├─ infrastructure/ # DB, external APIs, adapters
│   └─ transport/      # HTTP, bot, and other transport adapters
├─ pkg/                # Shared utilities and third‑party wrappers
├─ sql/                # SQL schema and migrations
├─ tests/              # Unit and integration tests (mirrors `internal/`)
├─ e2e/                # End‑to‑end tests (tavern/pytest, parallelizable)
└─ Taskfile.yml        # Task runner definitions
```
* Add new code under the appropriate `internal/*` layer to respect DDD boundaries.
* Test files should mirror the package layout, e.g. `internal/domain/user.go` → `tests/domain/user_test.go`.
* Place future end‑to‑end specifications under `e2e/` using **tavern**. They can be run in parallel with `pytest -n auto`.

## Build, Test, and Development Commands
The project uses **Taskfile** as the primary task runner:

| Task | Command |
|---|---|
| `task test` | Run all Go tests (`go test -v ./...`). |
| `task lint` | Lint code using `golangci-lint`. |
| `task sqlc` | Generate type‑safe Go from SQL (`sqlc generate`). |
| `task swagger` | Generate Swagger docs (`swag init -g cmd/main.go`). |
| `task e2e` | Run tavern end‑to‑end tests (`pytest -n auto e2e`). |

All tasks are defined in **Taskfile.yml** at the repository root.

## Coding Style \& Naming Conventions
* **Indentation** – 4 spaces, no tabs.
* **Line length** – Max 100 characters.
* **Naming** – `snake_case` for variables/functions, `PascalCase` for types, `UPPER_SNAKE_CASE` for constants.
* **Formatting** – `go fmt ./...` and `go vet` are run via the `lint` task.
* **Imports** – Grouped as standard library, third‑party, then local packages, each block separated by a blank line.

## Testing Guidelines
* **Unit/Integration** – Use Go's built‑in `testing` package; run via `task test`.
* **End‑to‑End** – Write tavern YAML specs in `e2e/`. Run with `task e2e` to execute with pytest in parallel.
* **Coverage** – Aim for ≥ 85 % total coverage; `go test -cover ./...` reports the percentage.

## Commit \& Pull Request Guidelines
* **Commit messages** – Follow the conventional‑commits format:
  `type(scope): short description`
  e.g., `feat(auth): add JWT refresh endpoint`.
* **PR requirements** – Provide a clear description, link related issue, and ensure the CI pipeline passes before merging.
* **Review** – At least one approval from a maintainer is required; address all review comments before final merge.

---
Adhering to these guidelines speeds up onboarding and keeps the codebase clean. Thanks for contributing!

