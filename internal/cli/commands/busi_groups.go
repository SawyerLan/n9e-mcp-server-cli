package commands

import (
	"github.com/n9e/n9e-mcp-server/internal/cli/output"
	"github.com/n9e/n9e-mcp-server/internal/config"
	"github.com/n9e/n9e-mcp-server/pkg/app"
	"github.com/n9e/n9e-mcp-server/pkg/client"

	"github.com/spf13/cobra"
)

// CLIContext holds shared dependencies for all CLI commands.
type CLIContext struct {
	Config config.Config
	Client *client.Client
}

// CtxBuilder is a function that lazily builds a CLIContext.
type CtxBuilder = func() (*CLIContext, error)

// NewBusiGroupsCommand creates the "busi-groups" command group.
func NewBusiGroupsCommand(buildCtx CtxBuilder) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "busi-groups",
		Aliases: []string{"bg"},
		Short:   "Manage business groups",
	}

	cmd.AddCommand(newBusiGroupsListCommand(buildCtx))
	return cmd
}

func newBusiGroupsListCommand(buildCtx CtxBuilder) *cobra.Command {
	var (
		limit int
		page  int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List business groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.ListBusiGroups(cmd.Context(), ctx.Client, app.ListBusiGroupsInput{
				Limit: limit,
				Page:  page,
			})
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 20, "Number of items per page")
	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")

	return cmd
}
