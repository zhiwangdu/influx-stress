# influx_stress

Standalone extraction of the InfluxDB v1.x `influx_stress` command.

## Build and Test

```bash
go test ./...
go build ./...
go build -o /tmp/influx_stress ./cmd/influx_stress
```

The module uses `github.com/influxdata/influxdb1-client` and `github.com/influxdata/influxql` as external dependencies instead of importing packages from the InfluxDB source repository.
