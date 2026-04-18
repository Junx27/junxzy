# Junxzy CLI

[![CI](https://github.com/Junx27/junxzy/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Junx27/junxzy/actions/workflows/ci.yml)
[![CD](https://github.com/Junx27/junxzy/actions/workflows/cd.yml/badge.svg)](https://github.com/Junx27/junxzy/actions/workflows/cd.yml)
[![Coverage](https://codecov.io/gh/Junx27/junxzy/branch/main/graph/badge.svg)](https://codecov.io/gh/Junx27/junxzy)

Junxzy is a Go-based CLI for generating modules, initializing microservice projects, and simulating the project scaffolding flow from the terminal.

## Requirements

- Go 1.25.4 or newer
- Git

## Installation

### Option 1: Build from source

```bash
git clone https://github.com/Junx27/junxzy.git
cd junxzy
go mod tidy
go build -o junxzy .
```

### Option 2: Install with Go

If you publish a tagged release, you can install it with:

```bash
go install github.com/Junx27/junxzy@latest
```

## Usage

Run the CLI after building:

```bash
./junxzy
```

The application starts an interactive REPL with commands such as:

- `help`
- `clear`
- `exit`
- `init <project-name>`
- `make:module <module-name>`
- `simulate`

## Development

Run tests locally:

```bash
go test ./...
```

Run coverage locally:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Latest Test and Coverage

Last checked: 2026-04-18

### Test result

```text
go test ./...
PASS
```

### Coverage summary

```text
github.com/Junx27/junxzy                 100.0%
github.com/Junx27/junxzy/cmd/commands      43.0%
github.com/Junx27/junxzy/internal/cli      87.9%
github.com/Junx27/junxzy/internal/file     80.0%
github.com/Junx27/junxzy/internal/generator 74.6%
github.com/Junx27/junxzy/internal/ui      100.0%
TOTAL                                      67.6%
```

## CI/CD

This repository includes GitHub Actions workflows for:

- lint
- test
- coverage
- release builds on tag push
