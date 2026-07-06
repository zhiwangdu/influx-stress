// Command influx_stress runs the InfluxDB stress test tool.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/influxdata/influx-stress/internal/app"
	"github.com/influxdata/influx-stress/internal/cli"
)

func main() {
	opts, err := cli.Parse(os.Args[1:], os.Stderr)
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if opts.CPUProfile != "" {
		f, err := os.Create(opts.CPUProfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	app.RunStress(opts.Config)
}
