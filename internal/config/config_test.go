package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/n9e/n9e-mcp-server/pkg/client"
	"github.com/n9e/n9e-mcp-server/pkg/toolset"

	"github.com/spf13/cobra"
)

func TestLoadDefaults(t *testing.T) {
	v := NewViper()
	cmd := &cobra.Command{Use: "test"}

	if err := BindRuntimeFlags(cmd, v); err != nil {
		t.Fatalf("BindRuntimeFlags() error = %v", err)
	}
	if err := BindCLIFlags(cmd, v); err != nil {
		t.Fatalf("BindCLIFlags() error = %v", err)
	}

	cfg := Load(v)

	if cfg.Token != "" {
		t.Fatalf("Token = %q, want empty", cfg.Token)
	}
	if cfg.BaseURL != DefaultBaseURL {
		t.Fatalf("BaseURL = %q, want %q", cfg.BaseURL, DefaultBaseURL)
	}
	if !reflect.DeepEqual(cfg.EnabledToolsets, toolset.DefaultToolsets) {
		t.Fatalf("EnabledToolsets = %v, want %v", cfg.EnabledToolsets, toolset.DefaultToolsets)
	}
	if cfg.Output != DefaultOutput {
		t.Fatalf("Output = %q, want %q", cfg.Output, DefaultOutput)
	}
	if cfg.Timeout != client.DefaultTimeout {
		t.Fatalf("Timeout = %s, want %s", cfg.Timeout, client.DefaultTimeout)
	}
}

func TestLoadEnvOverridesDefaults(t *testing.T) {
	t.Setenv("N9E_TOKEN", "env-token")
	t.Setenv("N9E_BASE_URL", "https://env.example")
	t.Setenv("N9E_TOOLSETS", "users,mutes")
	t.Setenv("N9E_READ_ONLY", "true")

	v := NewViper()
	cmd := &cobra.Command{Use: "test"}

	if err := BindRuntimeFlags(cmd, v); err != nil {
		t.Fatalf("BindRuntimeFlags() error = %v", err)
	}
	if err := BindCLIFlags(cmd, v); err != nil {
		t.Fatalf("BindCLIFlags() error = %v", err)
	}

	cfg := Load(v)

	if cfg.Token != "env-token" {
		t.Fatalf("Token = %q, want env-token", cfg.Token)
	}
	if cfg.BaseURL != "https://env.example" {
		t.Fatalf("BaseURL = %q, want https://env.example", cfg.BaseURL)
	}
	if want := []string{"users", "mutes"}; !reflect.DeepEqual(cfg.EnabledToolsets, want) {
		t.Fatalf("EnabledToolsets = %v, want %v", cfg.EnabledToolsets, want)
	}
	if !cfg.ReadOnly {
		t.Fatalf("ReadOnly = false, want true")
	}
}

func TestLoadFlagsOverrideEnv(t *testing.T) {
	t.Setenv("N9E_TOKEN", "env-token")
	t.Setenv("N9E_BASE_URL", "https://env.example")
	t.Setenv("N9E_TOOLSETS", "users,mutes")
	t.Setenv("N9E_READ_ONLY", "false")

	v := NewViper()
	cmd := &cobra.Command{Use: "test"}

	if err := BindRuntimeFlags(cmd, v); err != nil {
		t.Fatalf("BindRuntimeFlags() error = %v", err)
	}
	if err := BindCLIFlags(cmd, v); err != nil {
		t.Fatalf("BindCLIFlags() error = %v", err)
	}

	if err := cmd.ParseFlags([]string{
		"--token=flag-token",
		"--base-url=https://flag.example",
		"--toolsets=alerts,busi-groups",
		"--read-only",
		"--log-file=/tmp/n9e.log",
		"--output=table",
		"--timeout=45s",
		"--yes",
		"--quiet",
	}); err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}

	cfg := Load(v)

	if cfg.Token != "flag-token" {
		t.Fatalf("Token = %q, want flag-token", cfg.Token)
	}
	if cfg.BaseURL != "https://flag.example" {
		t.Fatalf("BaseURL = %q, want https://flag.example", cfg.BaseURL)
	}
	if want := []string{"alerts", "busi-groups"}; !reflect.DeepEqual(cfg.EnabledToolsets, want) {
		t.Fatalf("EnabledToolsets = %v, want %v", cfg.EnabledToolsets, want)
	}
	if !cfg.ReadOnly {
		t.Fatalf("ReadOnly = false, want true")
	}
	if cfg.LogFilePath != "/tmp/n9e.log" {
		t.Fatalf("LogFilePath = %q, want /tmp/n9e.log", cfg.LogFilePath)
	}
	if cfg.Output != "table" {
		t.Fatalf("Output = %q, want table", cfg.Output)
	}
	if cfg.Timeout != 45*time.Second {
		t.Fatalf("Timeout = %s, want 45s", cfg.Timeout)
	}
	if !cfg.Yes {
		t.Fatalf("Yes = false, want true")
	}
	if !cfg.Quiet {
		t.Fatalf("Quiet = false, want true")
	}
}
