package commands

import (
	"fmt"

	"github.com/n9e/n9e-mcp-server/internal/cli/output"
	"github.com/n9e/n9e-mcp-server/pkg/app"

	"github.com/spf13/cobra"
)

// NewAlertsCommand creates the "alerts" command group.
func NewAlertsCommand(buildCtx CtxBuilder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alerts",
		Short: "Manage alerts and alert rules",
	}

	cmd.AddCommand(
		newAlertsActiveCommand(buildCtx),
		newAlertsHistoryCommand(buildCtx),
		newAlertRulesCommand(buildCtx),
	)
	return cmd
}

func newAlertsActiveCommand(buildCtx CtxBuilder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active",
		Short: "Manage active alerts",
	}

	cmd.AddCommand(
		newAlertsActiveListCommand(buildCtx),
		newAlertsActiveGetCommand(buildCtx),
	)
	return cmd
}

func newAlertsActiveListCommand(buildCtx CtxBuilder) *cobra.Command {
	var input app.ListActiveAlertsInput

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List active alert events",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.ListActiveAlerts(cmd.Context(), ctx.Client, input)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().Int64Var(&input.Hours, "hours", 0, "Lookback hours")
	cmd.Flags().Int64Var(&input.Stime, "stime", 0, "Start time Unix timestamp")
	cmd.Flags().Int64Var(&input.Etime, "etime", 0, "End time Unix timestamp")
	cmd.Flags().StringVar(&input.Severity, "severity", "", "Severity levels comma-separated (1=critical, 2=warning, 3=info)")
	cmd.Flags().StringVar(&input.Query, "query", "", "Search keyword")
	cmd.Flags().StringVar(&input.Cate, "cate", "", "Alert category")
	cmd.Flags().StringVar(&input.RuleProds, "rule-prods", "", "Product types comma-separated")
	cmd.Flags().StringVar(&input.DatasourceIds, "datasource-ids", "", "Datasource IDs comma-separated")
	cmd.Flags().Int64Var(&input.RuleId, "rule-id", 0, "Alert rule ID")
	cmd.Flags().Int64Var(&input.BusiGroupId, "group-id", 0, "Business group ID")
	cmd.Flags().IntVarP(&input.Limit, "limit", "l", 20, "Number of items per page")
	cmd.Flags().IntVarP(&input.Page, "page", "p", 1, "Page number")

	return cmd
}

func newAlertsActiveGetCommand(buildCtx CtxBuilder) *cobra.Command {
	var eventId int64

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of an active alert event",
		RunE: func(cmd *cobra.Command, args []string) error {
			if eventId <= 0 {
				return fmt.Errorf("--event-id is required and must be positive")
			}

			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.GetActiveAlert(cmd.Context(), ctx.Client, eventId)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().Int64Var(&eventId, "event-id", 0, "Alert event ID")
	_ = cmd.MarkFlagRequired("event-id")

	return cmd
}

func newAlertsHistoryCommand(buildCtx CtxBuilder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Manage historical alerts",
	}

	cmd.AddCommand(
		newAlertsHistoryListCommand(buildCtx),
		newAlertsHistoryGetCommand(buildCtx),
	)
	return cmd
}

func newAlertsHistoryListCommand(buildCtx CtxBuilder) *cobra.Command {
	var input app.ListHistoryAlertsInput

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List historical alert events",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.ListHistoryAlerts(cmd.Context(), ctx.Client, input)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().Int64Var(&input.Hours, "hours", 0, "Lookback hours")
	cmd.Flags().Int64Var(&input.Stime, "stime", 0, "Start time Unix timestamp")
	cmd.Flags().Int64Var(&input.Etime, "etime", 0, "End time Unix timestamp")
	cmd.Flags().IntVar(&input.Severity, "severity", 0, "Severity level (-1=all, 1=critical, 2=warning, 3=info)")
	cmd.Flags().IntVar(&input.IsRecovered, "is-recovered", 0, "Recovery status (-1=all, 0=not recovered, 1=recovered)")
	cmd.Flags().StringVar(&input.Query, "query", "", "Search keyword")
	cmd.Flags().StringVar(&input.Cate, "cate", "", "Alert category")
	cmd.Flags().StringVar(&input.RuleProds, "rule-prods", "", "Product types comma-separated")
	cmd.Flags().StringVar(&input.DatasourceIds, "datasource-ids", "", "Datasource IDs comma-separated")
	cmd.Flags().Int64Var(&input.BusiGroupId, "group-id", 0, "Business group ID")
	cmd.Flags().IntVarP(&input.Limit, "limit", "l", 20, "Number of items per page")
	cmd.Flags().IntVarP(&input.Page, "page", "p", 1, "Page number")

	return cmd
}

func newAlertsHistoryGetCommand(buildCtx CtxBuilder) *cobra.Command {
	var eventId int64

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of a historical alert event",
		RunE: func(cmd *cobra.Command, args []string) error {
			if eventId <= 0 {
				return fmt.Errorf("--event-id is required and must be positive")
			}

			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.GetHistoryAlert(cmd.Context(), ctx.Client, eventId)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().Int64Var(&eventId, "event-id", 0, "Alert event ID")
	_ = cmd.MarkFlagRequired("event-id")

	return cmd
}

func newAlertRulesCommand(buildCtx CtxBuilder) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rules",
		Aliases: []string{"rule"},
		Short:   "Manage alert rules",
	}

	cmd.AddCommand(
		newAlertRulesListCommand(buildCtx),
		newAlertRulesGetCommand(buildCtx),
	)
	return cmd
}

func newAlertRulesListCommand(buildCtx CtxBuilder) *cobra.Command {
	var input app.ListAlertRulesInput

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List alert rules for a business group",
		RunE: func(cmd *cobra.Command, args []string) error {
			if input.GroupId <= 0 {
				return fmt.Errorf("--group-id is required and must be positive")
			}

			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.ListAlertRules(cmd.Context(), ctx.Client, input)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().Int64Var(&input.GroupId, "group-id", 0, "Business group ID")
	cmd.Flags().IntVarP(&input.Limit, "limit", "l", 20, "Number of items per page")
	cmd.Flags().IntVarP(&input.Page, "page", "p", 1, "Page number")
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

func newAlertRulesGetCommand(buildCtx CtxBuilder) *cobra.Command {
	var ruleId int64

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of an alert rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if ruleId <= 0 {
				return fmt.Errorf("--rule-id is required and must be positive")
			}

			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.GetAlertRule(cmd.Context(), ctx.Client, ruleId)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().Int64Var(&ruleId, "rule-id", 0, "Alert rule ID")
	_ = cmd.MarkFlagRequired("rule-id")

	return cmd
}
