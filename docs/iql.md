# Influx Stress Tool V2

```
$ influx_stress -config examples/iql/file.iql
```

This stress tool works from list of InfluxQL-esque statements. The language has been extended to allow for some basic templating of fields, tags and measurements in both line protocol and query statements.

By default the test outputs a human readable report to `STDOUT` and records test statistics in an active installation of InfluxDB at `localhost:8086`.

To set state variables for the test such as the address of the Influx node use the following syntax:

```
# The values listed below are the default values for each of the parameters
 
# Pipe delineated list of addresses. For cluster: [192.168.0.10:8086|192.168.0.2:8086|192.168.0.3:8086]
# Queries and writes are round-robin to the configured addresses.
SET Addresses [localhost:8086]

# False (default) uses http, true uses https
SET SSL [false]

# Username for targeted influx server or cluster
SET Username []

# Password for targeted influx server or cluster
SET Password []

# Database to target for queries and writes. Works like the InfluxCLI USE
SET Database [stress]

# Precision for the data being written
# Only s and ns supported
SET Precision [s]

# Date the first written point will be timestamped
SET StartDate [2016-01-01]

# Size of batches to send to InfluxDB
SET BatchSize [5000]

# Time to wait between sending batches
SET WriteInterval [0s]

# Time to wait between sending queries
SET QueryInterval [0s]

# Number of concurrent writers
SET WriteConcurrency [15]

# Number of concurrent readers
SET QueryConcurrency [5]
```

The values in the example are also the defaults.

Valid line protocol will be forwarded right to the server making setting up your testing environment easy:

```
CREATE DATABASE thing

ALTER RETENTION POLICY default ON thing DURATION 1h REPLICATION 1

SET database [thing]
```

You can write points like this:
```
INSERT mockCpu
cpu,
host=server-[int inc(0) 10000],location=[str rand(8) 1000]
value=[float rand(1000) 0]
100000 10s

Explained:

# INSERT keyword kicks off the statement, next to it is the name of the statement for reporting and templated query generation
INSERT mockCpu
# Measurement
cpu,
# Tags - separated by commas. Tag values can be templates, mixed template and fixed values
host=server-[float rand(100) 10000],location=[int inc(0) 1000],fixed=[fix|fid|dor|pom|another_tag_value]
# Fields - separated by commas either templates, mixed template and fixed values
value=[float inc(0) 0]
# 'Timestamp' - Number of points to insert into this measurement and the amount of time between points
100000 10s
```

Each template contains 3 parts: a datatype (`str`, `float`, or `int`), a function with one non-negative integer argument, and the number of unique generated values to cache and cycle. A final count of `0` means values advance after each full series cycle.

Supported functions:

- `int`: `inc(n)` starts at `n` and increments, `dec(n)` starts at `n` and decrements, `const(n)` repeats `n`, `rand(n)` is uniform in `[0,n)`, and `zipf(n)` is a skewed distribution over `[0,n)`.
- `float`: `inc(n)`, `dec(n)`, `const(n)`, `rand(n)`, and `zipf(n)` emit integer-looking float field values. `randf(n)` emits decimal values in `[0,n)`, `normal(n)` emits normally distributed decimal noise with scale `n`, `sin(n)` emits a decimal sine wave with amplitude `n`, and `walk(n)` emits a decimal random walk with step magnitude `n`.
- `str`: `rand(n)` emits line-safe random hex with length `2*floor(n/2)`, `inc(n)` emits increasing decimal strings, `id(n)` emits `id-` plus a zero-padded counter of width `n`, and `hash(n)` emits deterministic hex with length `2*floor(n/2)`.

Functions with an upper bound (`rand`, `randf`, `zipf`) require `n > 0`, except `str rand(0)` which emits an empty string.

To run multiple insert statements at once:
```
GO INSERT devices
devices,
city=[str rand(8) 10],country=[str rand(8) 25],device_id=[str rand(10) 1000]
lat=[float rand(90) 0],lng=[float rand(120) 0],temp=[float rand(40) 0]
10000000 10s

GO INSERT devices2
devices2,
city=[str rand(8) 10],country=[str rand(8) 25],device_id=[str rand(10) 1000]
lat=[float rand(90) 0],lng=[float rand(120) 0],temp=[float rand(40) 0]
10000000 10s

WAIT
```

Fastest point generation and write load requires 3-4 running `GO INSERT` statements at a time.

You can run queries like this:

```
QUERY cpu
SELECT mean(value) FROM cpu WHERE host='server-1'
DO 1000
```

### Output:
Output for config file in this repo:
```
[√] "CREATE DATABASE thing" -> 1.806785ms
[√] "CREATE DATABASE thing2" -> 1.492504ms
SET Database = 'thing'
SET Precision = 's'
Go Write Statement:                    mockCpu
  Points/Sec:                          245997
  Resp Time Average:                   173.354445ms
  Resp Time Standard Deviation:        123.80344ms
  95th Percentile Write Response:      381.363503ms
  Average Request Bytes:               276110
  Successful Write Reqs:               20
  Retries:                             0
Go Query Statement:                    mockCpu
  Resp Time Average:                   3.140803ms
  Resp Time Standard Deviation:        2.292328ms
  95th Percentile Read Response:       5.915437ms
  Query Resp Bytes Average:            16 bytes
  Successful Queries:                  10
WAIT -> 406.400059ms
SET DATABASE = 'thing2'
Go Write Statement:                    devices
  Points/Sec:                          163348
  Resp Time Average:                   132.553789ms
  Resp Time Standard Deviation:        149.397972ms
  95th Percentile Write Response:      567.987467ms
  Average Request Bytes:               459999
  Successful Write Reqs:               20
  Retries:                             0
Go Write Statement:                    devices2
  Points/Sec:                          160078
  Resp Time Average:                   133.303097ms
  Resp Time Standard Deviation:        144.352404ms
  95th Percentile Write Response:      560.565066ms
  Average Request Bytes:               464999
  Successful Write Reqs:               20
  Retries:                             0
Go Query Statement:                    fooName
  Resp Time Average:                   1.3307ms
  Resp Time Standard Deviation:        640.249µs
  95th Percentile Read Response:       2.668ms
  Query Resp Bytes Average:            16 bytes
  Successful Queries:                  10
WAIT -> 624.585319ms
[√] "DROP DATABASE thing" -> 991.088464ms
[√] "DROP DATABASE thing2" -> 421.362831ms
```

### Next Steps:

##### Documentation
- Parser behavior and proper `.iql` syntax
- How the templated query generation works
- Collection of tested `.iql` files to simulate different loads
  
##### Performance
- `Commune`, a stuct to enable templated Query generation, is blocking writes when used, look into performance.
- Templated query generation is currently in a quazi-working state. See the above point.
