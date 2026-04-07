package cli

import (
	"fmt"

	"github.com/n9e/n9e-mcp-server/internal/cli/commands"
	"github.com/n9e/n9e-mcp-server/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cliUserAgent = "n9e-cli"

// NewCLICommand creates the "cli" subcommand with all child commands registered.
func NewCLICommand(v *viper.Viper, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cli",
		Short: "Interactive CLI mode for Nightingale",
		Long:  "Query and manage Nightingale resources from the command line",
	}

	// Register CLI-only flags so they appear in help and get bound via viper.
	flags := cmd.PersistentFlags()
	flags.String("output", config.DefaultOutput, "Output format (json)")
	flags.Duration("timeout", 0, "Request timeout")
	flags.Bool("yes", false, "Skip interactive confirmation")
	flags.Bool("quiet", false, "Suppress non-essential output")

	// Bind CLI-only flags to viper.
	_ = v.BindPFlag("output", flags.Lookup("output"))
	_ = v.BindPFlag("timeout", flags.Lookup("timeout"))
	_ = v.BindPFlag("yes", flags.Lookup("yes"))
	_ = v.BindPFlag("quiet", flags.Lookup("quiet"))

	// Build CLIContext lazily inside each command's RunE via a helper.
	buildCtx := func() (*commands.CLIContext, error) {
		cfg := config.Load(v)
		if cfg.Token == "" {
			return nil, fmt.Errorf("N9E_TOKEN is required. Set it via --token flag or N9E_TOKEN environment variable")
		}
		c, err := config.NewAPIClient(cfg, cliUserAgent+"/"+version)
		if err != nil {
			return nil, err
		}
		return &commands.CLIContext{Config: cfg, Client: c}, nil
	}

	// Register domain commands
	cmd.AddCommand(commands.NewBusiGroupsCommand(buildCtx))

	return cmd
}
