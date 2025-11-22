# Examples

This page provides practical examples of Pace configurations for various use cases.

## Go Project

A complete build, test, and deployment pipeline for a Go application:

```pace
# Variables
var app_name = "myapp"
var version = "1.0.0"
var build_dir = "bin"

# Default task
default build

# Aliases
alias b build
alias t test
alias d dev

# Development task
task dev {
    description "Run development server with auto-reload"
    watch true
    inputs ["**/*.go"]
    command "go run cmd/server/main.go"
}

# Test task
task test {
    description "Run all tests with coverage"
    command "go test -v -cover ./..."
    inputs ["**/*.go", "**/*_test.go"]
    cache true
}

# Lint task
hook lint {
    description "Run golangci-lint"
    command "golangci-lint run ./..."
}

# Build task
task build {
    description "Build the application"
    command "go build -ldflags '-X main.Version=${version}' -o ${build_dir}/${app_name} cmd/server/main.go"
    before ["test", "lint"]
    inputs ["**/*.go", "go.mod", "go.sum"]
    outputs ["${build_dir}/${app_name}"]
    cache true
    env {
        "CGO_ENABLED" "0"
    }
}

# Docker build
task docker {
    description "Build Docker image"
    command "docker build -t ${app_name}:${version} ."
    before ["build"]
    inputs ["Dockerfile", "${build_dir}/${app_name}"]
}

# Deploy task
task deploy {
    description "Deploy to production"
    command "./scripts/deploy.sh ${version}"
    before ["docker"]
    on_success ["notify"]
    timeout "15m"
    retry 2
    retry_delay "10s"
}

hook notify {
    description "Send deployment notification"
    command "echo 'Deployed ${app_name} v${version}'"
}
```

## Multi-Language Monorepo

Build both backend and frontend projects:

```pace
default all

# Build everything
task all {
    description "Build backend and frontend"
    before ["backend", "frontend"]
    parallel true
}

# Backend (Go)
task backend {
    description "Build Go backend"
    command "go build -o bin/server cmd/server/main.go"
    inputs ["cmd/**/*.go", "internal/**/*.go", "go.mod"]
    outputs ["bin/server"]
    cache true
}

task backend-test {
    description "Run backend tests"
    command "go test ./..."
    inputs ["**/*.go"]
}

# Frontend (Node.js)
task frontend {
    description "Build React frontend"
    working_dir "frontend"
    command "npm run build"
    before ["frontend-install"]
    inputs ["frontend/src/**/*", "frontend/package.json"]
    outputs ["frontend/dist/**/*"]
    cache true
}

hook frontend-install {
    description "Install frontend dependencies"
    working_dir "frontend"
    command "npm install"
}

task frontend-dev {
    description "Start frontend dev server"
    working_dir "frontend"
    watch true
    inputs ["frontend/src/**/*.{ts,tsx,css}"]
    command "npm run dev"
}

# Development
task dev {
    description "Run both backend and frontend in dev mode"
    before ["backend", "frontend-dev"]
    parallel true
}
```

## Rust Project

Build, test, and package a Rust application:

```pace
default build

# Aliases
alias b build
alias t test
alias r run

# Format code
hook fmt {
    description "Format code with rustfmt"
    command "cargo fmt"
}

# Lint code
hook clippy {
    description "Run clippy linter"
    command "cargo clippy -- -D warnings"
}

# Run tests
task test {
    description "Run all tests"
    command "cargo test --all-features"
    before ["fmt", "clippy"]
    inputs ["src/**/*.rs", "tests/**/*.rs", "Cargo.toml"]
    cache true
}

# Build debug
task build {
    description "Build debug binary"
    command "cargo build"
    before ["test"]
    inputs ["src/**/*.rs", "Cargo.toml"]
    outputs ["target/debug/myapp"]
    cache true
}

# Build release
task release {
    description "Build optimized release binary"
    command "cargo build --release"
    before ["test"]
    inputs ["src/**/*.rs", "Cargo.toml"]
    outputs ["target/release/myapp"]
    env {
        "RUSTFLAGS" "-C target-cpu=native"
    }
}

# Run application
task run {
    description "Run the application"
    command "cargo run"
}

# Benchmark
task bench {
    description "Run benchmarks"
    command "cargo bench"
}
```

## Python Project

Development workflow for a Python application:

```pace
default test

# Aliases
alias t test
alias l lint
alias f format

# Setup virtual environment
hook venv {
    description "Create virtual environment"
    command "python -m venv .venv"
}

# Install dependencies
hook install {
    description "Install dependencies"
    command ".venv/bin/pip install -r requirements.txt"
    before ["venv"]
}

# Format code
hook format {
    description "Format code with black"
    command ".venv/bin/black src/ tests/"
}

# Lint code
hook lint {
    description "Lint with flake8 and mypy"
    command """
        .venv/bin/flake8 src/ tests/
        .venv/bin/mypy src/
    """
}

# Run tests
task test {
    description "Run pytest with coverage"
    command ".venv/bin/pytest tests/ --cov=src --cov-report=html"
    before ["install", "lint"]
    inputs ["src/**/*.py", "tests/**/*.py"]
    cache true
}

# Run application
task run {
    description "Run the application"
    command ".venv/bin/python -m src.main"
    before ["install"]
}

# Development with auto-reload
task dev {
    description "Run with auto-reload"
    watch true
    inputs ["src/**/*.py"]
    command ".venv/bin/python -m src.main --reload"
}
```

## CI/CD Pipeline

Comprehensive CI/CD configuration:

```pace
default ci

# CI pipeline
task ci {
    description "Run full CI pipeline"
    before ["lint", "test", "build", "security-scan"]
    parallel false
}

# Code quality checks
hook lint {
    description "Run all linters"
    command """
        echo 'Running linters...'
        golangci-lint run
        go fmt ./...
    """
}

# Security scanning
task security-scan {
    description "Run security checks"
    command """
        gosec ./...
        go list -json -m all | nancy sleuth
    """
    continue_on_error true
}

# Run tests
task test {
    description "Run tests with coverage"
    command "go test -v -race -coverprofile=coverage.out ./..."
    outputs ["coverage.out"]
    timeout "10m"
}

# Build application
task build {
    description "Build for multiple platforms"
    command """
        GOOS=linux GOARCH=amd64 go build -o dist/app-linux-amd64 main.go
        GOOS=darwin GOARCH=amd64 go build -o dist/app-darwin-amd64 main.go
        GOOS=windows GOARCH=amd64 go build -o dist/app-windows-amd64.exe main.go
    """
    before ["test"]
    inputs ["**/*.go"]
    outputs ["dist/*"]
    cache true
}

# Create release
task release {
    description "Create GitHub release"
    command "./scripts/create-release.sh"
    before ["ci"]
}
```

## Documentation Site

Build and serve documentation:

```pace
default build

# Aliases
alias s serve
alias d dev

# Install dependencies
hook install {
    description "Install documentation tools"
    command "pip install mkdocs mkdocs-material"
}

# Build docs
task build {
    description "Build documentation site"
    command "mkdocs build"
    before ["install"]
    inputs ["docs/**/*.md", "mkdocs.yml"]
    outputs ["site/**/*"]
    cache true
}

# Serve locally
task serve {
    description "Serve documentation locally"
    command "mkdocs serve"
    before ["install"]
}

# Development with auto-reload
task dev {
    description "Develop with live reload"
    watch true
    inputs ["docs/**/*.md", "mkdocs.yml"]
    command "mkdocs serve"
    before ["install"]
}

# Deploy to GitHub Pages
task deploy {
    description "Deploy to GitHub Pages"
    command "mkdocs gh-deploy --force"
    before ["build"]
}
```

## Database Migrations

Manage database migrations:

```pace
# Variables
var db_url = "${DATABASE_URL}"

# Run migrations
task migrate {
    description "Run database migrations"
    command "migrate -path migrations -database ${db_url} up"
    env {
        "DATABASE_URL" "postgres://localhost/mydb?sslmode=disable"
    }
}

# Rollback migration
task migrate-down {
    description "Rollback last migration"
    command "migrate -path migrations -database ${db_url} down 1"
}

# Create new migration
task migrate-create {
    description "Create new migration"
    args {
        required ["name"]
    }
    command "migrate create -ext sql -dir migrations -seq $name"
}

# Reset database
task db-reset {
    description "Reset database (drop and recreate)"
    command """
        dropdb mydb
        createdb mydb
    """
    after ["migrate"]
}

# Seed database
hook seed {
    description "Seed database with test data"
    command "psql ${db_url} < seeds/test-data.sql"
}
```

## Environment-Specific Tasks

Different tasks for different environments:

```pace
# Variables
var env = "${ENVIRONMENT}"

# Development build
task dev {
    description "Build for development"
    command "go build -o bin/app-dev main.go"
    env {
        "ENVIRONMENT" "development"
        "DEBUG" "true"
    }
}

# Staging build
task staging {
    description "Build for staging"
    command "go build -ldflags '-X main.Env=staging' -o bin/app-staging main.go"
    before ["test"]
    env {
        "ENVIRONMENT" "staging"
    }
}

# Production build
task production {
    description "Build for production"
    command "go build -ldflags '-X main.Env=production -s -w' -o bin/app-prod main.go"
    before ["test", "security-scan"]
    env {
        "ENVIRONMENT" "production"
        "CGO_ENABLED" "0"
    }
}

# Deploy with environment
task deploy {
    description "Deploy to specified environment"
    args {
        required ["env"]
    }
    command "./scripts/deploy.sh $env"
    timeout "20m"
    retry 2
}
```

## Next Steps

- [Configuration Reference](configuration.md) - Learn about all available options
- [Commands Reference](commands/list.md) - Explore all CLI commands
- [Quick Start Guide](quick-start.md) - Get started quickly
