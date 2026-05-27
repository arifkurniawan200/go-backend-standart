# Go Backend Standards

Scalable Go backend project structure with **Clean Architecture**, **Git workflow standards**, and **best practices**.

## 📁 Project Structure

```
myapp/
├── cmd/                    # Entry point per service
│   └── api/
│       └── main.go
├── internal/               # Private application code
│   ├── domain/             # Entity, value object (no dependencies)
│   ├── repository/         # Data access interface + implementation
│   ├── usecase/            # Business logic
│   └── handler/            # HTTP/gRPC handlers
├── pkg/                    # Public packages (sharable across projects)
│   ├── validator/
│   ├── response/
│   └── logger/
├── config/                 # Configuration
├── scripts/                # Build, migration, code gen scripts
├── .github/
│   └── workflows/
│       └── ci.yml
├── Makefile
├── golangci.yml
└── go.mod
```

## 🚀 Quick Start

```bash
# Install dependencies
go mod tidy

# Run
go run cmd/api/main.go

# Build
make build

# Test
make test

# Lint
make lint
```

## 📋 Git Workflow

### Branch Naming
- `feature/` — New features
- `bugfix/` — Bug fixes
- `hotfix/` — Production fixes
- `chore/` — Maintenance tasks
- `refactor/` — Code refactoring

### Commit Convention (Conventional Commits)
```
feat: add user authentication
fix: handle nil pointer in repository
chore: upgrade Go to 1.22
docs: update API documentation
refactor: extract validation to separate package
```

### PR Requirements
- Minimum 1 approval
- All CI checks passed
- No merge conflicts
- Squash merge to main

## 🛠️ Standards

### 1. Dependency Injection via Interface
```go
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

type UserUsecase struct {
    repo UserRepository // loosely coupled
}
```

### 2. Context Propagation
Always pass `context.Context` for timeout, cancellation, and tracing.

### 3. Graceful Shutdown
```go
srv := &http.Server{Addr: ":8080"}
go func() { srv.ListenAndServe() }()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

### 4. Worker Pool Pattern
```go
func workerPool(ctx context.Context, jobs <-chan Job, workers int) {
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                process(job)
            }
        }()
    }
    wg.Wait()
}
```

## ⚡ Optimization Tips

1. **JSON encoding** — Use `sonic` or `json-iterator` for speed-critical paths
2. **Connection pooling** — Set `SetMaxOpenConns`, `SetMaxIdleConns` on DB
3. **Batch insert** for bulk DB writes
4. **Use `pprof`** for CPU/memory/goroutine profiling

## 📦 Tools

| Tool | Purpose |
|------|---------|
| `gofmt` / `goimports` | Code formatting |
| `golangci-lint` | Linting + static analysis |
| `go vet` | Static analysis |
| `make` | Task automation |

## 🐳 Docker & Traefik

### Architecture

```
                    ┌─────────────────────────────┐
                    │          Traefik            │
                    │   (Reverse Proxy / Gateway) │
                    │                             │
  Internet ────────►│  :80 ─► Router ─► Service  │
                    │                     │       │
                    │                     ▼       │
                    │              [Middleware]    │
                    │           Rate Limit        │
                    │           Strip Prefix      │
                    │           CORS              │
                    │           Compress          │
                    └─────────┬───────────────────┘
                              │ HTTP
           ┌──────────────────┼──────────────────┐
           │                  ▼                  │
           │  ┌────────────────────────────┐   │
           │  │   Go API (Backend)         │   │
           │  │   Port: 8080              │   │
           │  │   /health → healthy       │   │
           │  │   /api/v1/* → handlers    │   │
           │  └────────────────────────────┘   │
           └───────────────────────────────────┘
```

### Quick Start with Docker

```bash
# Start all services (Traefik + API)
make docker-up

# Test routing
make docker-test

# View Traefik logs
make docker-logs-traefik

# Stop services
make docker-down
```

### Endpoints

| URL | Description |
|-----|-------------|
| `http://localhost/api/v1/users` | API endpoint (via Traefik) |
| `http://localhost/health` | Health check |
| `http://localhost:8080/dashboard/` | Traefik Dashboard |

### Middleware Pipeline

```
Request → Rate Limit → Strip Prefix → CORS → Compress → Backend
          (100 req/s)  (/api/v1 → /)  (OPTIONS)  (gzip)
```

### Files

```
docker/
├── docker-compose.yml          # Dev compose
├── docker-compose.prod.yml     # Prod compose
├── Dockerfile                  # Multi-stage build
├── traefik/
│   ├── traefik.yml             # Static config (entrypoints, providers)
│   └── dynamic/
│       └── services.yml        # Dynamic config (routers, services, middleware)
└── .env.example
```

## 📄 License

MIT
