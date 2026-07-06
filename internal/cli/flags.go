// Package cli parses command-line options for influx_stress.
package cli

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

const DefaultConfig = "examples/iql/file.iql"

// Options contains validated command-line options.
type Options struct {
	Config     string
	CPUProfile string
}

var unsupportedV1Flags = map[string]struct{}{
	"addr":             {},
	"database":         {},
	"db":               {},
	"retention-policy": {},
	"tags":             {},
}

// Parse validates CLI arguments. The -v2 flag is accepted for compatibility,
// but v2 IQL is now the only supported runner.
func Parse(args []string, output io.Writer) (Options, error) {
	if output == nil {
		output = io.Discard
	}

	fs := flag.NewFlagSet("influx_stress", flag.ContinueOnError)
	fs.SetOutput(output)

	useV2 := fs.Bool("v2", true, "use v2 IQL runner; retained for compatibility")
	config := fs.String("config", "", "v2 IQL stress test file")
	cpuProfile := fs.String("cpuprofile", "", "write the CPU profile to filename")

	for name := range unsupportedV1Flags {
		fs.String(name, "", "unsupported legacy v1 flag")
	}

	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}
	if !*useV2 {
		return Options{}, fmt.Errorf("v1 TOML mode has been removed; omit -v2 or pass -v2=true to run v2 IQL")
	}

	var legacy []string
	fs.Visit(func(f *flag.Flag) {
		if _, ok := unsupportedV1Flags[f.Name]; ok {
			legacy = append(legacy, "-"+f.Name)
		}
	})
	if len(legacy) != 0 {
		return Options{}, fmt.Errorf("unsupported v1 flag(s): %s; use v2 IQL SET statements instead", strings.Join(legacy, ", "))
	}

	if *config != "" && strings.EqualFold(filepath.Ext(*config), ".toml") {
		return Options{}, fmt.Errorf("v1 TOML config %q is no longer supported; pass a v2 .iql config", *config)
	}

	opts := Options{
		Config:     *config,
		CPUProfile: *cpuProfile,
	}
	if opts.Config == "" {
		opts.Config = DefaultConfig
	}
	return opts, nil
}
