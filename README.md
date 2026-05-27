# Go Backend Standards

Scalable Go backend project structure with **Clean Architecture**, **Git workflow standards**, **Traefik integration**, and **observability**.

## 📁 Project Structure

```
myapp/
├── cmd/
│   ├── api/                    # Main API service
│   │   └── main.go
│   └── auth/                  # Auth service
│       └── main.go
├── internal/
│   ├── domain/               # Entity, value object
│   ├── repository/           # Data access (interface-based)
│   ├── usecase/              # Business logic
│   ├── handler/              # HTTP handlers
│   └── middleware/           # JWT, auth middleware
├── pkg/
│   ├── validator/
│   ├── response/
│   └── logger/
├── config/
├── docker/
│   ├── traefik/
│   │   ├── traefik.yml
│   │   └── dynamic/
│   │       └── services.yml
│   └── monitoring/
│       ├── prometheus/
│       └── grafana/
├── docker-compose.yml          # Dev (Traefik + API)
├── docker-compose.prod.yml     # Production
├── docker-compose.multi.yml    # Multi-service (API + Auth)
├── docker-compose.monitoring.yml # Monitoring (Prometheus + Grafana)
├── Dockerfile
├── Dockerfile.auth
├── Makefile
└── go.mod
```

## 🚀 Quick Start

```bash
# Install dependencies
go mod tidy

# Run API
go run cmd/api/main.go

# Run Auth service
go run cmd/auth/main.go

# Build
make build
make build-auth

# Test
make test
```

## 🐳 Docker Compose Options

### 1. Basic (Traefik + API)
```bash
make docker-up        # Start
make docker-down      # Stop
make docker-test      # Test routing
make docker-logs      # View logs
```

### 2. Production
```bash
make docker-up-prod   # Start with resource limits
make docker-down-prod  # Stop
```

### 3. Multi-Service (API + Auth)
```bash
make docker-up-multi   # Start API + Auth
make docker-down-multi  # Stop
make docker-logs-api   # API logs
make docker-logs-auth  # Auth logs
```

### 4. Monitoring (Prometheus + Grafana)
```bash
make docker-up-monitoring   # Start monitoring
make docker-down-monitoring  # Stop
# Prometheus: http://localhost:9090
# Grafana:    http://localhost:3000 (admin/admin)
```

## 🌐 Endpoints

| URL | Service | Description |
|-----|---------|-------------|
| `http://localhost/api/v1/users` | API | REST endpoints |
| `http://localhost/auth/login` | Auth | Login endpoint |
| `http://localhost/auth/register` | Auth | Register endpoint |
| `http://localhost/health` | API | Health check |
| `http://localhost/auth/health` | Auth | Auth health check |
| `http://localhost:8080/dashboard/` | Traefik | Dashboard |
| `http://localhost:9090` | Prometheus | Metrics |
| `http://localhost:3000` | Grafana | Dashboards |

## 🔒 HTTPS & Security

Traefik dengan Let's Encrypt untuk automatic HTTPS:
- HTTP → HTTPS redirect
- Certificate auto-renewal
- Security headers (HSTS, X-Frame-Options, etc.)

## ⚡ Middleware Pipeline

```
Request → Rate Limit → Strip Prefix → CORS → Security Headers → Compress → Backend
          (100/s)      /api/v1→/     (OPTIONS)   (HSTS, etc)      (gzip)
```

## 🔄 Load Balancing & Resilience

- **Circuit Breaker**: `NetworkErrorRatio() > 0.30`
- **Retries**: 3 attempts with exponential backoff
- **Health Checks**: Every 10s per service
- **Connection pooling**: Configured per service

## 📊 Monitoring Stack

### Prometheus Metrics
- Traefik request rate, latency, errors
- Service-level metrics
- Custom application metrics (future)

### Grafana Dashboards
- Traefik Overview (request rate, latency, errors)
- Pre-provisioned datasource & dashboards
- Default credentials: `admin/admin`

## 🔐 Auth Service

JWT-based authentication:
- `/auth/login` — Returns access + refresh tokens
- `/auth/register` — Create new user
- `/auth/refresh` — Refresh access token
- `/auth/validate` — Validate token (for forwardAuth)

## 📋 Git Workflow

### Branch Naming
- `feature/` — New features
- `bugfix/` — Bug fixes
- `hotfix/` — Production fixes
- `chore/` — Maintenance tasks

### Commit Convention
```
feat: add user authentication
fix: handle nil pointer in repository
chore: upgrade Go to 1.22
docs: update API documentation
refactor: extract validation
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
    repo UserRepository
}
```

### 2. Context Propagation
```go
func (uc *UserUsecase) FindByID(ctx context.Context, id string) (*User, error)
```

### 3. Graceful Shutdown
```go
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

1. **JSON encoding** — Use `sonic` or `json-iterator`
2. **Connection pooling** — Set `SetMaxOpenConns`, `SetMaxIdleConns`
3. **Batch insert** for bulk DB writes
4. **Use `pprof`** for profiling

## 📦 Tools

| Tool | Purpose |
|------|---------|
| `gofmt` / `goimports` | Code formatting |
| `golangci-lint` | Linting + static analysis |
| `go vet` | Static analysis |
| `make` | Task automation |
| `docker compose` | Container orchestration |

## 📄 License

MIT
