package main

import (
	"fmt"
	"os"

	"github.com/n9e/n9e-mcp-server/internal"
	"github.com/n9e/n9e-mcp-server/internal/config"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	cfgViper = config.NewViper()
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "n9e-mcp-server",
	Short: "Nightingale MCP Server",
	Long:  "MCP (Model Context Protocol) server for Nightingale",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default to run in stdio mode
		return runStdio(cmd, args)
	},
}

var stdioCmd = &cobra.Command{
	Use:   "stdio",
	Short: "Run in stdio mode (default)",
	Long:  "Run the MCP server using stdin/stdout for communication",
	RunE:  runStdio,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("n9e-mcp-server %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

func init() {
	if err := config.BindRuntimeFlags(rootCmd, cfgViper); err != nil {
		panic(err)
	}

	// Add subcommands
	rootCmd.AddCommand(stdioCmd)
	rootCmd.AddCommand(versionCmd)
}

func runStdio(cmd *cobra.Command, args []string) error {
	cfg := config.Load(cfgViper)
	if cfg.Token == "" {
		return fmt.Errorf("N9E_TOKEN is required. Set it via --token flag or N9E_TOKEN environment variable")
	}

	return internal.RunStdioServer(internal.StdioServerConfig{
		Version: version,
		Config:  cfg,
	})
}
