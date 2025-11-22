# Project Structure & Repository Strategy

## 1. Recommendation: The Monorepo
For a single developer (or small team) building a distributed system, I **strongly recommend a Monorepo** (One single Git repository).

**Why?**
-   **Shared Code**: You can easily share Protobuf definitions (`.proto`) and Go utility packages between services without managing complex `go.mod` replace directives or private module registries.
-   **Unified Build**: One `Makefile` to run everything. `make up` starts the whole world.
-   **Atomic Changes**: You can change the API contract (Proto) and the Server implementation in a single Commit/PR.
-   **Easier Learning**: You don't need to context-switch between 4 different VS Code windows.

---

## 2. Directory Structure
We will organize the code to look like "separate apps" but living together.

```text
pixelflow/                  # Root of the repository
├── Makefile                # Master makefile (build, run, test all)
├── docker-compose.yml      # Run everything locally (Kafka, Mongo, Postgres, Apps)
├── README.md               # Project documentation
├── .github/                # GitHub Actions (CI/CD)
│
├── apps/                   # The source code for your services
│   ├── auth/               # Auth Service (Go + gRPC)
│   │   ├── cmd/            # Main entrypoint
│   │   ├── internal/       # Business logic
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   ├── api/                # API Service (Go + HTTP)
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   ├── worker/             # Worker Service (Go + Kafka Consumer)
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   └── frontend/           # React App (Vite)
│       ├── src/
│       ├── package.json
│       └── Dockerfile
│
├── pkg/                    # Shared Go Code (Libraries)
│   ├── pb/                 # Generated Protobuf Go code (shared between Auth & API)
│   └── logger/             # Common logging setup
│
├── proto/                  # Protocol Buffer Definitions
│   └── auth.proto          # The contract between API and Auth
│
└── deploy/                 # Infrastructure as Code
    ├── k8s/                # Kubernetes Manifests
    └── helm/               # Helm Charts
```

## 3. How it works
1.  **`apps/`**: Each folder here is effectively a "microservice". They have their own `go.mod` (or we can use a workspace `go.work` file in the root).
2.  **`pkg/`**: This is the magic of the monorepo. The `auth` service and `api` service can both import `github.com/yourname/pixelflow/pkg/pb`.
3.  **`deploy/`**: Keeps your infrastructure separate from your code.

## 4. Summary
-   **Repositories to create**: **1** (Just `pixelflow`).
-   **Complexity**: Low overhead, high velocity.
