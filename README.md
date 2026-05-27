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

## 📄 License

MIT
