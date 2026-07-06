# influx_stress

Standalone `influx_stress` command using the v2 IQL stress test runner.

## Build and Test

```bash
go test ./...
go build ./...
go build -o /tmp/influx_stress ./cmd/influx_stress
```

## Run

```bash
/tmp/influx_stress
/tmp/influx_stress -config examples/iql/default.iql
/tmp/influx_stress -v2 -config examples/iql/default.iql
```

The `-v2` flag is accepted for compatibility but no longer required. Legacy v1 TOML configs and v1-only flags are rejected.

The module uses `github.com/influxdata/influxdb1-client` and `github.com/influxdata/influxql` as external dependencies instead of importing packages from the InfluxDB source repository.
