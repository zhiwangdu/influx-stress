package cli

import (
	"io"
	"strings"
	"testing"
)

func TestParseDefaultsToV2Config(t *testing.T) {
	opts, err := Parse(nil, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Config != DefaultConfig {
		t.Fatalf("wrong default config: %q", opts.Config)
	}
}

func TestParseAcceptsCompatibleV2Flag(t *testing.T) {
	opts, err := Parse([]string{"-v2", "-config", "examples/iql/default.iql"}, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Config != "examples/iql/default.iql" {
		t.Fatalf("wrong config: %q", opts.Config)
	}
}

func TestParseRejectsDisabledV2(t *testing.T) {
	_, err := Parse([]string{"-v2=false"}, io.Discard)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "v1 TOML mode has been removed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRejectsV1Flags(t *testing.T) {
	_, err := Parse([]string{"-db", "stress"}, io.Discard)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unsupported v1 flag(s): -db") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRejectsTomlConfig(t *testing.T) {
	_, err := Parse([]string{"-config", "stress.toml"}, io.Discard)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "v1 TOML config") {
		t.Fatalf("unexpected error: %v", err)
	}
}
