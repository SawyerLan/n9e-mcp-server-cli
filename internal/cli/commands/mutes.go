package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/n9e/n9e-mcp-server/internal/cli/output"
	"github.com/n9e/n9e-mcp-server/pkg/app"
	"github.com/n9e/n9e-mcp-server/pkg/types"

	"github.com/spf13/cobra"
)

// checkWriteAllowed returns an error if the CLI is in read-only mode.
func checkWriteAllowed(ctx *CLIContext) error {
	if ctx.Config.ReadOnly {
		return fmt.Errorf("write operation denied: --read-only mode is enabled")
	}
	return nil
}

// confirmAction prompts the user for confirmation unless --yes is set.
// Returns nil if confirmed, error otherwise.
func confirmAction(ctx *CLIContext, message string) error {
	if ctx.Config.Yes {
		return nil
	}
	fmt.Fprintf(os.Stderr, "%s [y/N]: ", message)
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read confirmation: %w", err)
	}
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "y" && answer != "yes" {
		return fmt.Errorf("operation cancelled by user")
	}
	return nil
}

// NewMutesCommand creates the "mutes" command group.
func NewMutesCommand(buildCtx CtxBuilder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mutes",
		Short: "Manage alert mutes/silences",
	}

	cmd.AddCommand(
		newMutesListCommand(buildCtx),
		newMutesGetCommand(buildCtx),
		newMutesCreateCommand(buildCtx),
		newMutesUpdateCommand(buildCtx),
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

// parseTagFlags parses tag flags in "key=value" format into TagFilter slice (using == operator).
func parseTagFlags(raw []string) []types.TagFilter {
	var tags []types.TagFilter
	for _, s := range raw {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) == 2 {
			tags = append(tags, types.TagFilter{Key: parts[0], Func: "==", Value: parts[1]})
		}
	}
	return tags
}

func newMutesCreateCommand(buildCtx CtxBuilder) *cobra.Command {
	var (
		groupId       int64
		note          string
		cate          string
		prod          string
		datasourceIds []int64
		tagFlags      []string
		cause         string
		duration      string
		btime         int64
		etime         int64
		severities    []int
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new alert mute/silence rule",
		Long: `Create a new alert mute rule. Time can be specified either as:
  --btime/--etime (Unix timestamps), or
  --duration (e.g., "30m", "2h") which sets btime=now, etime=now+duration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			if err := checkWriteAllowed(ctx); err != nil {
				return err
			}

			// Handle duration shortcut
			if duration != "" {
				d, err := time.ParseDuration(duration)
				if err != nil {
					return fmt.Errorf("invalid --duration: %w", err)
				}
				now := time.Now()
				btime = now.Unix()
				etime = now.Add(d).Unix()
			}

			input := app.CreateMuteInput{
				GroupId:       groupId,
				Note:          note,
				Cate:          cate,
				Prod:          prod,
				DatasourceIds: datasourceIds,
				Tags:          parseTagFlags(tagFlags),
				Cause:         cause,
				Btime:         btime,
				Etime:         etime,
				Severities:    severities,
			}

			// Show summary and confirm
			btimeStr := time.Unix(input.Btime, 0).Format("2006-01-02 15:04:05")
			etimeStr := time.Unix(input.Etime, 0).Format("2006-01-02 15:04:05")
			summary := fmt.Sprintf("Create mute: group=%d, cause=%q, time=%s ~ %s, tags=%d",
				input.GroupId, input.Cause, btimeStr, etimeStr, len(input.Tags))
			if err := confirmAction(ctx, summary); err != nil {
				return err
			}

			result, err := app.CreateMute(cmd.Context(), ctx.Client, input)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().Int64Var(&groupId, "group-id", 0, "Business group ID (required)")
	cmd.Flags().StringVar(&note, "note", "", "Note/title for the mute rule")
	cmd.Flags().StringVar(&cate, "cate", "", "Category (e.g., prometheus, host)")
	cmd.Flags().StringVar(&prod, "prod", "", "Product type (e.g., metric, host)")
	cmd.Flags().Int64SliceVar(&datasourceIds, "datasource-ids", nil, "Datasource IDs (comma-separated)")
	cmd.Flags().StringArrayVar(&tagFlags, "tag", nil, "Tag filter in key=value format (repeatable)")
	cmd.Flags().StringVar(&cause, "cause", "", "Reason for the mute (required)")
	cmd.Flags().StringVar(&duration, "duration", "", "Mute duration from now (e.g., 30m, 2h)")
	cmd.Flags().Int64Var(&btime, "btime", 0, "Start time Unix timestamp")
	cmd.Flags().Int64Var(&etime, "etime", 0, "End time Unix timestamp")
	cmd.Flags().IntSliceVar(&severities, "severities", nil, "Severity levels (1=critical, 2=warning, 3=info)")
	_ = cmd.MarkFlagRequired("group-id")
	_ = cmd.MarkFlagRequired("cause")

	return cmd
}

func newMutesUpdateCommand(buildCtx CtxBuilder) *cobra.Command {
	var (
		groupId       int64
		muteId        int64
		note          string
		cate          string
		prod          string
		datasourceIds []int64
		tagFlags      []string
		cause         string
		duration      string
		btime         int64
		etime         int64
		severities    []int
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing alert mute/silence rule",
		Long: `Update an existing alert mute rule. Time can be specified either as:
  --btime/--etime (Unix timestamps), or
  --duration (e.g., "30m", "2h") which sets btime=now, etime=now+duration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildCtx()
			if err != nil {
				return err
			}

			if err := checkWriteAllowed(ctx); err != nil {
				return err
			}

			// Handle duration shortcut
			if duration != "" {
				d, err := time.ParseDuration(duration)
				if err != nil {
					return fmt.Errorf("invalid --duration: %w", err)
				}
				now := time.Now()
				btime = now.Unix()
				etime = now.Add(d).Unix()
			}

			input := app.UpdateMuteInput{
				GroupId:       groupId,
				MuteId:        muteId,
				Note:          note,
				Cate:          cate,
				Prod:          prod,
				DatasourceIds: datasourceIds,
				Tags:          parseTagFlags(tagFlags),
				Cause:         cause,
				Btime:         btime,
				Etime:         etime,
				Severities:    severities,
			}

			btimeStr := time.Unix(input.Btime, 0).Format("2006-01-02 15:04:05")
			etimeStr := time.Unix(input.Etime, 0).Format("2006-01-02 15:04:05")
			summary := fmt.Sprintf("Update mute: group=%d, mute_id=%d, cause=%q, time=%s ~ %s",
				input.GroupId, input.MuteId, input.Cause, btimeStr, etimeStr)
			if err := confirmAction(ctx, summary); err != nil {
				return err
			}

			result, err := app.UpdateMute(cmd.Context(), ctx.Client, input)
			if err != nil {
				return err
			}

			return output.Render(ctx.Config.Output, result)
		},
	}

	cmd.Flags().Int64Var(&groupId, "group-id", 0, "Business group ID (required)")
	cmd.Flags().Int64Var(&muteId, "mute-id", 0, "Alert mute ID (required)")
	cmd.Flags().StringVar(&note, "note", "", "Note/title for the mute rule")
	cmd.Flags().StringVar(&cate, "cate", "", "Category (e.g., prometheus, host)")
	cmd.Flags().StringVar(&prod, "prod", "", "Product type (e.g., metric, host)")
	cmd.Flags().Int64SliceVar(&datasourceIds, "datasource-ids", nil, "Datasource IDs (comma-separated)")
	cmd.Flags().StringArrayVar(&tagFlags, "tag", nil, "Tag filter in key=value format (repeatable)")
	cmd.Flags().StringVar(&cause, "cause", "", "Reason for the mute (required)")
	cmd.Flags().StringVar(&duration, "duration", "", "Mute duration from now (e.g., 30m, 2h)")
	cmd.Flags().Int64Var(&btime, "btime", 0, "Start time Unix timestamp")
	cmd.Flags().Int64Var(&etime, "etime", 0, "End time Unix timestamp")
	cmd.Flags().IntSliceVar(&severities, "severities", nil, "Severity levels (1=critical, 2=warning, 3=info)")
	_ = cmd.MarkFlagRequired("group-id")
	_ = cmd.MarkFlagRequired("mute-id")
	_ = cmd.MarkFlagRequired("cause")

	return cmd
}
