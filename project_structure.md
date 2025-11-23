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
├── distributed_app_design.md # Architecture documentation
├── test_e2e.sh             # End-to-End test script
│
├── apps/                   # The source code for your services
│   ├── auth/               # Auth Service (Go + HTTP REST)
│   │   ├── cmd/            # Main entrypoint
│   │   ├── internal/       # Business logic (db, models, utils)
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── IMPLEMENTATION.md
│   │
│   ├── api/                # API Service (Go + HTTP REST)
│   │   ├── cmd/
│   │   ├── internal/       # Business logic (db, models, kafka, middleware)
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── IMPLEMENTATION.md
│   │
│   ├── worker/             # Worker Service (Go + Kafka Consumer)
│   │   ├── cmd/
│   │   ├── internal/       # Business logic (db, models, processor)
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── IMPLEMENTATION.md
│   │
│   └── frontend/           # Frontend UI (React + TailwindCSS)
│       ├── src/            # React source code
│       ├── public/
│       ├── package.json
│       ├── Dockerfile
│       ├── nginx.conf
│       └── IMPLEMENTATION.md
│
└── deploy/                 # Infrastructure as Code (Future)
    ├── k8s/                # Kubernetes Manifests
    └── helm/               # Helm Charts
```

## 3. How it works
1.  **`apps/`**: Each folder here is a self-contained service. They have their own `go.mod` and internal packages (`internal/`) to ensure loose coupling.
2.  **`frontend/`**: The React application that interacts with the API service.
3.  **`deploy/`**: Keeps your infrastructure separate from your code.

## 4. Summary
-   **Repositories to create**: **1** (Just `pixelflow`).
-   **Architecture**: Monorepo with independent, decoupled services.
