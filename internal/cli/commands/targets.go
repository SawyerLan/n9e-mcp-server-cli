package commands

import (
	"github.com/n9e/n9e-mcp-server/internal/cli/output"
	"github.com/n9e/n9e-mcp-server/pkg/app"

	"github.com/spf13/cobra"
)

// NewTargetsCommand creates the "targets" command group.
func NewTargetsCommand(buildCtx CtxBuilder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "targets",
		Short: "Manage monitored targets/hosts",
	}

	cmd.AddCommand(newTargetsListCommand(buildCtx))
	return cmd
}

func newTargetsListCommand(buildCtx CtxBuilder) *cobra.Command {
	var input app.ListTargetsInput

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List monitored targets",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.ListTargets(cmd.Context(), ctx.Client, input)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().StringVar(&input.GroupIds, "group-ids", "", "Business group IDs comma-separated")
	cmd.Flags().StringVar(&input.Query, "query", "", "Search keyword (matches ident/tags)")
	cmd.Flags().IntVarP(&input.Limit, "limit", "l", 20, "Number of items per page")
	cmd.Flags().IntVarP(&input.Page, "page", "p", 1, "Page number")
	cmd.Flags().Int64Var(&input.Downtime, "downtime", 0, "Filter by downtime in seconds")
	cmd.Flags().StringVar(&input.DatasourceIds, "datasource-ids", "", "Datasource IDs comma-separated")

	return cmd
}
