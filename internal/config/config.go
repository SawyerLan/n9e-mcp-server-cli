package config

import (
	"strings"
	"time"

	"github.com/n9e/n9e-mcp-server/pkg/client"
	"github.com/n9e/n9e-mcp-server/pkg/toolset"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	EnvPrefix      = "N9E"
	DefaultBaseURL = "http://localhost:17000"
	DefaultOutput  = "json"
)

// Config contains shared runtime configuration for stdio and CLI mode.
type Config struct {
	Token           string
	BaseURL         string
	EnabledToolsets []string
	ReadOnly        bool
	LogFilePath     string
	Output          string
	Timeout         time.Duration
	Yes             bool
	Quiet           bool
}

// NewViper returns a configured viper instance with the shared defaults.
func NewViper() *viper.Viper {
	v := viper.New()
	ConfigureViper(v)
	SetDefaults(v)
	return v
}

// ConfigureViper enables shared env var semantics.
func ConfigureViper(v *viper.Viper) {
	v.SetEnvPrefix(EnvPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

// SetDefaults applies shared defaults so flag/env precedence stays consistent.
func SetDefaults(v *viper.Viper) {
	v.SetDefault("base_url", DefaultBaseURL)
	v.SetDefault("toolsets", append([]string(nil), toolset.DefaultToolsets...))
	v.SetDefault("read_only", false)
	v.SetDefault("log_file", "")
	v.SetDefault("output", DefaultOutput)
	v.SetDefault("timeout", client.DefaultTimeout)
	v.SetDefault("yes", false)
	v.SetDefault("quiet", false)
}

// BindRuntimeFlags binds the shared runtime flags used by stdio and CLI mode.
func BindRuntimeFlags(cmd *cobra.Command, v *viper.Viper) error {
	flags := cmd.PersistentFlags()
	flags.String("token", "", "Nightingale API token (env: N9E_TOKEN)")
	flags.String("base-url", DefaultBaseURL, "Nightingale API base URL (env: N9E_BASE_URL)")
	flags.StringSlice("toolsets", toolset.DefaultToolsets, "Enabled toolsets (env: N9E_TOOLSETS)")
	flags.Bool("read-only", false, "Read-only mode, disable write operations (env: N9E_READ_ONLY)")
	flags.String("log-file", "", "Log file path (default: stderr)")

	return bindFlagSet(v, flags,
		"token",
		"base_url=base-url",
		"toolsets",
		"read_only=read-only",
		"log_file=log-file",
	)
}

// BindCLIFlags binds CLI-only flags on top of the shared runtime flags.
func BindCLIFlags(cmd *cobra.Command, v *viper.Viper) error {
	flags := cmd.PersistentFlags()
	flags.String("output", DefaultOutput, "Output format")
	flags.Duration("timeout", client.DefaultTimeout, "Request timeout")
	flags.Bool("yes", false, "Skip interactive confirmation")
	flags.Bool("quiet", false, "Suppress non-essential output")

	return bindFlagSet(v, flags,
		"output",
		"timeout",
		"yes",
		"quiet",
	)
}

// Load reads the shared runtime configuration from viper.
func Load(v *viper.Viper) Config {
	return Config{
		Token:           v.GetString("token"),
		BaseURL:         v.GetString("base_url"),
		EnabledToolsets: loadStringSlice(v, "toolsets"),
		ReadOnly:        v.GetBool("read_only"),
		LogFilePath:     v.GetString("log_file"),
		Output:          v.GetString("output"),
		Timeout:         v.GetDuration("timeout"),
		Yes:             v.GetBool("yes"),
		Quiet:           v.GetBool("quiet"),
	}
}

// NewAPIClient creates a Nightingale client from shared config.
func NewAPIClient(cfg Config, userAgent string) (*client.Client, error) {
	return client.NewClientWithOptions(cfg.Token, cfg.BaseURL, userAgent, client.Options{
		Timeout: cfg.Timeout,
	})
}

func bindFlagSet(v *viper.Viper, flags *pflag.FlagSet, bindings ...string) error {
	for _, binding := range bindings {
		key := binding
		flagName := binding

		if before, after, ok := strings.Cut(binding, "="); ok {
			key = before
			flagName = after
		}

		if err := v.BindPFlag(key, flags.Lookup(flagName)); err != nil {
			return err
		}
	}

	return nil
}

func loadStringSlice(v *viper.Viper, key string) []string {
	raw := v.Get(key)
	switch value := raw.(type) {
	case nil:
		return nil
	case []string:
		return append([]string(nil), value...)
	case string:
		if value == "" {
			return nil
		}

		parts := strings.Split(value, ",")
		items := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				items = append(items, part)
			}
		}
		return items
	default:
		return append([]string(nil), v.GetStringSlice(key)...)
	}
}
