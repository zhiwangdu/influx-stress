// Command influx_stress runs the InfluxDB stress test tool.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"

	v2 "github.com/influxdata/influx-stress/stress/v2"
)

var (
	useV2      = flag.Bool("v2", true, "Use version 2 of stress tool. Deprecated: v2 is now the only supported engine")
	config     = flag.String("config", "", "The v2 IQL stress test file")
	cpuprofile = flag.String("cpuprofile", "", "Write the cpu profile to `filename`")
)

var unsupportedV1Flags = map[string]struct{}{
	"addr":             {},
	"database":         {},
	"db":               {},
	"retention-policy": {},
	"tags":             {},
}

func main() {
	registerUnsupportedV1Flags()
	flag.Parse()

	if err := validateFlags(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	file := *config
	if file == "" {
		file = "stress/v2/iql/file.iql"
	}
	v2.RunStress(file)
}

func registerUnsupportedV1Flags() {
	for name := range unsupportedV1Flags {
		flag.String(name, "", "unsupported legacy v1 flag")
	}
}

func validateFlags() error {
	if !*useV2 {
		return fmt.Errorf("v1 TOML mode has been removed; omit -v2 or pass -v2=true to run v2 IQL")
	}

	var legacy []string
	flag.Visit(func(f *flag.Flag) {
		if _, ok := unsupportedV1Flags[f.Name]; ok {
			legacy = append(legacy, "-"+f.Name)
		}
	})
	if len(legacy) != 0 {
		return fmt.Errorf("unsupported v1 flag(s): %s; use v2 IQL SET statements instead", strings.Join(legacy, ", "))
	}

	if *config != "" && strings.EqualFold(filepath.Ext(*config), ".toml") {
		return fmt.Errorf("v1 TOML config %q is no longer supported; pass a v2 .iql config", *config)
	}

	return nil
}
