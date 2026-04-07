package commands

import (
	"fmt"

	"github.com/n9e/n9e-mcp-server/internal/cli/output"
	"github.com/n9e/n9e-mcp-server/pkg/app"

	"github.com/spf13/cobra"
)

// NewMutesCommand creates the "mutes" command group.
func NewMutesCommand(buildCtx CtxBuilder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mutes",
		Short: "Manage alert mutes/silences",
	}

	cmd.AddCommand(
		newMutesListCommand(buildCtx),
		newMutesGetCommand(buildCtx),
	)
	return cmd
}

func newMutesListCommand(buildCtx CtxBuilder) *cobra.Command {
	var input app.ListMutesInput

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List alert mutes for a business group",
		RunE: func(cmd *cobra.Command, args []string) error {
			if input.GroupId <= 0 {
				return fmt.Errorf("--group-id is required and must be positive")
			}

			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.ListMutes(cmd.Context(), ctx.Client, input)
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

func newMutesGetCommand(buildCtx CtxBuilder) *cobra.Command {
	var (
		groupId int64
		muteId  int64
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of an alert mute",
		RunE: func(cmd *cobra.Command, args []string) error {
			if groupId <= 0 {
				return fmt.Errorf("--group-id is required and must be positive")
			}
			if muteId <= 0 {
				return fmt.Errorf("--mute-id is required and must be positive")
			}

			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.GetMute(cmd.Context(), ctx.Client, groupId, muteId)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().Int64Var(&groupId, "group-id", 0, "Business group ID")
	cmd.Flags().Int64Var(&muteId, "mute-id", 0, "Alert mute ID")
	_ = cmd.MarkFlagRequired("group-id")
	_ = cmd.MarkFlagRequired("mute-id")

	return cmd
}
