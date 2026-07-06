# Influx Stress Architecture

The tool is organized around an IQL script runner. The command parses a file into runtime statements, executes them against an InfluxDB target, records response metadata, and renders a report.

## Packages

- `cmd/influx_stress`: CLI flags, compatibility checks, and handoff to the app runner.
- `internal/app`: top-level orchestration for parsing, execution, result flush, and report printing.
- `internal/iql`: IQL file splitting and statement parsing. The `parse` subpackage contains the scanner/parser for individual statements.
- `internal/engine`: runtime statements such as `SET`, `GO`, `WAIT`, `INSERT`, `QUERY`, `EXEC`, and raw InfluxQL.
- `internal/workload`: INSERT workload generation: templates, built-in functions, timestamp generation, and point rendering.
- `internal/influx`: runtime state, InfluxDB HTTP query/write client, directives, packages, tracing, and result collection.
- `internal/report`: report aggregation, response-time statistics, and output formatting.
- `examples/iql`: runnable sample IQL scripts.

## Flow

1. `cmd/influx_stress` validates CLI flags and calls `app.RunStress`.
2. `internal/app` calls `iql.ParseStatements` to turn the IQL file into `engine.Statement` values.
3. `engine.Statement.Run` sends write/query packages or directives through `internal/influx`.
4. The Influx client performs HTTP requests, records response metadata, and decrements tracers.
5. After execution, each statement asks `internal/report` to format its results.

## Statement Interface

```go
type Statement interface {
	Run(s *influxclient.StressTest)
	Report(s *influxclient.StressTest) string
	SetID(s string)
}
```

`GO` wraps another statement in a goroutine. `WAIT` blocks on the shared runtime wait group. `INSERT` delegates point rendering to `internal/workload`; `QUERY` can either run a static query or use commune data from a companion insert.

## Runtime State

`internal/influx.StressTest` is the statement-facing runtime. It owns the configured precision, start date, batch size, package/directive channels, result channel, commune map, and reporting database client.

The unexported HTTP client in `internal/influx` handles query/write concurrency, retry behavior, write intervals, query intervals, target address selection, and response-point creation.

## Workload Generation

`internal/workload` compiles INSERT templates into append-based renderers. Built-in functions emit line-protocol-safe values without per-point `fmt.Sprintf` allocation. The trailing IQL count controls whether generated values are cached and cycled or advanced after each full series cycle.

## Reporting

`internal/report` receives response columns and values from the reporting database and computes response-time percentiles, averages, standard deviation, success counts, retry counts, and byte averages. Runtime statements keep only the responsibility of locating their result series and passing it to the report package.
