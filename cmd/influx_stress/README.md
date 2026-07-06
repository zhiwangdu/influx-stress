# `influx_stress`

If you run into any issues with this tool please mention @jackzampolin when you create an issue.

## Ways to run

### `influx_stress`
This runs the v2 stress test with the default IQL file at `examples/iql/file.iql`.

### `influx_stress -config someConfig.iql`
This runs the v2 stress test with a valid IQL file located at `someConfig.iql`.

### `influx_stress -v2 -config someConfig.iql`
This is still supported for compatibility. The `-v2` flag is now optional because v2 is the only supported runner.

## Flags

### `-config` string
The relative path to the v2 IQL stress test file.

`default` = `examples/iql/file.iql`

### `-cpuprofile` filename
Writes the result of Go's cpu profile to filename

`default` = no profiling

### `-v2`
Compatibility no-op. v2 is always used.

## Unsupported Legacy Flags

The v1 TOML runner has been removed. Flags such as `-db`, `-addr`, `-database`, `-retention-policy`, and `-tags` now return an error. Use v2 `SET` statements in IQL files instead.
