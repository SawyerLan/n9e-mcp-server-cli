package commands

import (
	"fmt"

	"github.com/n9e/n9e-mcp-server/internal/cli/output"
	"github.com/n9e/n9e-mcp-server/pkg/app"

	"github.com/spf13/cobra"
)

// NewUsersCommand creates the "users" command group.
func NewUsersCommand(buildCtx CtxBuilder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage users and user groups",
	}

	cmd.AddCommand(
		newUsersListCommand(buildCtx),
		newUsersGetCommand(buildCtx),
		newUserGroupsCommand(buildCtx),
	)
	return cmd
}

func newUsersListCommand(buildCtx CtxBuilder) *cobra.Command {
	var input app.ListUsersInput

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.ListUsers(cmd.Context(), ctx.Client, input)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().StringVar(&input.Query, "query", "", "Search keyword (matches username/nickname/email/phone)")
	cmd.Flags().IntVarP(&input.Limit, "limit", "l", 20, "Number of items per page")
	cmd.Flags().IntVarP(&input.Page, "page", "p", 1, "Page number")

	return cmd
}

func newUsersGetCommand(buildCtx CtxBuilder) *cobra.Command {
	var userId int64

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if userId <= 0 {
				return fmt.Errorf("--user-id is required and must be positive")
			}

			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.GetUser(cmd.Context(), ctx.Client, userId)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().Int64Var(&userId, "user-id", 0, "User ID")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func newUserGroupsCommand(buildCtx CtxBuilder) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "groups",
		Aliases: []string{"group"},
		Short:   "Manage user groups",
	}

	cmd.AddCommand(
		newUserGroupsListCommand(buildCtx),
		newUserGroupsGetCommand(buildCtx),
	)
	return cmd
}

func newUserGroupsListCommand(buildCtx CtxBuilder) *cobra.Command {
	var input app.ListUserGroupsInput

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List user groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.ListUserGroups(cmd.Context(), ctx.Client, input)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().StringVar(&input.Query, "query", "", "Search keyword for group name")
	cmd.Flags().IntVarP(&input.Limit, "limit", "l", 0, "Maximum number of groups to return (default 1500)")

	return cmd
}

func newUserGroupsGetCommand(buildCtx CtxBuilder) *cobra.Command {
	var groupId int64

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of a user group",
		RunE: func(cmd *cobra.Command, args []string) error {
			if groupId <= 0 {
				return fmt.Errorf("--group-id is required and must be positive")
			}

			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			result, err := app.GetUserGroup(cmd.Context(), ctx.Client, groupId)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().Int64Var(&groupId, "group-id", 0, "User group ID")
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}
