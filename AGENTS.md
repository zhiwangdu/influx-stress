# Repository Guidelines

## Project Structure & Module Organization

This repository is a standalone Go module for the InfluxDB `influx_stress` tool. The CLI entry point is in `cmd/influx_stress/`. The supported runner lives under `internal/`, including the app entrypoint in `internal/app/`, runtime statements in `internal/engine/`, the InfluxDB backend in `internal/influx/`, IQL parsing in `internal/iql/`, and sample `.iql` files in `examples/iql/`. Tests are colocated with source files as `*_test.go`.

## Build, Test, and Development Commands

- `go test ./...`: runs all package tests.
- `go build ./...`: verifies all packages compile.
- `go build -o /tmp/influx_stress ./cmd/influx_stress`: builds the CLI binary.
- `/tmp/influx_stress`: runs the default v2 IQL file at `examples/iql/file.iql`.
- `/tmp/influx_stress -config examples/iql/default.iql`: runs a specific v2 IQL configuration.
- `/tmp/influx_stress -v2 -config examples/iql/default.iql`: supported compatibility form; `-v2` is optional.

Use a local InfluxDB instance at `http://localhost:8086` when exercising the default runtime paths.

## Coding Style & Naming Conventions

Use standard Go formatting: run `gofmt` on changed `.go` files before committing. Keep package names short and lowercase, matching the current internal packages (`app`, `engine`, `influx`, `iql`, `report`, `workload`). Exported identifiers should have clear Go doc comments when they form part of a package API. Prefer explicit error handling and small helper functions over hidden global behavior.

## Testing Guidelines

Add or update colocated `*_test.go` files for parser, statement, client, and configuration changes. Follow existing table-driven test patterns where present, especially in `internal/engine` and `internal/iql`. Run `go test ./...` before opening a pull request. For changes that require a live InfluxDB server, document the manual command and target configuration used.

## Commit & Pull Request Guidelines

This checkout does not include Git history, so use concise imperative commit messages such as `Add query parser test` or `Fix stress client response handling`. Pull requests should include a short problem statement, a summary of code changes, test results, and any manual InfluxDB validation steps. Link related issues when available and include screenshots only for changes to visual assets such as `docs/influx_stress_v2.png`.

## Configuration & Safety Notes

Stress-test configs can create, reset, and write large volumes of data. Review database names, addresses, concurrency, and reset flags before running against shared or production InfluxDB instances.
