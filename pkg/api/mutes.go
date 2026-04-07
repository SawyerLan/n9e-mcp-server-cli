package api

import (
	"context"

	"github.com/n9e/n9e-mcp-server/pkg/app"
	"github.com/n9e/n9e-mcp-server/pkg/client"
	"github.com/n9e/n9e-mcp-server/pkg/toolset"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetMuteInput represents get single mute rule parameters
type GetMuteInput struct {
	GroupId int64 `json:"group_id"`
	MuteId  int64 `json:"mute_id"`
}


// RegisterMutesToolset registers alert mutes toolset
func RegisterMutesToolset(group *toolset.ToolsetGroup, getClient client.GetClientFunc) {
	ts := toolset.NewToolset("mutes", "Alert mute/silence management tools")

	ts.AddReadTools(
		listMutesTool(getClient),
		getMuteTool(getClient),
	)

	ts.AddWriteTools(
		createMuteTool(getClient),
		updateMuteTool(getClient),
	)

	group.AddToolset(ts)
}

func listMutesTool(getClient client.GetClientFunc) toolset.ServerTool {
	return toolset.NewServerTool(
		mcp.Tool{
			Name:        "list_mutes",
			Description: "List alert mutes/silences for a business group",
			Annotations: &mcp.ToolAnnotations{
				Title:        "List Alert Mutes",
				ReadOnlyHint: true,
			},
			InputSchema: &jsonschema.Schema{
				Type:     "object",
				Required: []string{"group_id"},
				Properties: map[string]*jsonschema.Schema{
					"group_id": {
						Type:        "integer",
						Description: "Business group ID",
					},
					"limit": {
						Type:        "integer",
						Description: "Page size (default 20)",
					},
					"p": {
						Type:        "integer",
						Description: "Page number (starts from 1)",
					},
				},
			},
		},
		toolset.MakeToolHandler(func(ctx context.Context, req *mcp.CallToolRequest, input app.ListMutesInput) (*mcp.CallToolResult, error) {
			if input.GroupId <= 0 {
				return toolset.NewToolResultError("group_id is required and must be positive"), nil
			}

			c := getClient(ctx)
			if c == nil {
				return toolset.NewToolResultError("failed to get n9e client from context"), nil
			}

			result, err := app.ListMutes(ctx, c, input)
			if err != nil {
				return toolset.NewToolResultError(err.Error()), nil
			}

			return toolset.MarshalResult(result), nil
		}),
	)
}

func getMuteTool(getClient client.GetClientFunc) toolset.ServerTool {
	return toolset.NewServerTool(
		mcp.Tool{
			Name:        "get_mute",
			Description: "Get details of a specific alert mute by ID",
			Annotations: &mcp.ToolAnnotations{
				Title:        "Get Alert Mute",
				ReadOnlyHint: true,
			},
			InputSchema: &jsonschema.Schema{
				Type:     "object",
				Required: []string{"group_id", "mute_id"},
				Properties: map[string]*jsonschema.Schema{
					"group_id": {
						Type:        "integer",
						Description: "Business group ID",
					},
					"mute_id": {
						Type:        "integer",
						Description: "Alert mute ID",
					},
				},
			},
		},
		toolset.MakeToolHandler(func(ctx context.Context, req *mcp.CallToolRequest, input GetMuteInput) (*mcp.CallToolResult, error) {
			if input.GroupId <= 0 {
				return toolset.NewToolResultError("group_id is required and must be positive"), nil
			}
			if input.MuteId <= 0 {
				return toolset.NewToolResultError("mute_id is required and must be positive"), nil
			}

			c := getClient(ctx)
			if c == nil {
				return toolset.NewToolResultError("failed to get n9e client from context"), nil
			}

			result, err := app.GetMute(ctx, c, input.GroupId, input.MuteId)
			if err != nil {
				return toolset.NewToolResultError(err.Error()), nil
			}

			return toolset.MarshalResult(result), nil
		}),
	)
}

func createMuteTool(getClient client.GetClientFunc) toolset.ServerTool {
	return toolset.NewServerTool(
		mcp.Tool{
			Name:        "create_mute",
			Description: "Create a new alert mute/silence rule. Use mute_time_type=0 for time range mode (btime/etime), or mute_time_type=1 for periodic mode (periodic_mutes).",
			Annotations: &mcp.ToolAnnotations{
				Title:           "Create Alert Mute",
				ReadOnlyHint:    false,
				DestructiveHint: toolset.BoolPtr(false),
			},
			InputSchema: &jsonschema.Schema{
				Type:     "object",
				Required: []string{"group_id", "cause", "btime", "etime"},
				Properties: map[string]*jsonschema.Schema{
					"group_id": {
						Type:        "integer",
						Description: "Business group ID",
					},
					"note": {
						Type:        "string",
						Description: "Note/title for the mute rule",
					},
					"cate": {
						Type:        "string",
						Description: "Category (e.g., prometheus, host, elasticsearch)",
					},
					"prod": {
						Type:        "string",
						Description: "Product type (e.g., metric, host, loki)",
					},
					"datasource_ids": {
						Type:        "array",
						Description: "Datasource IDs to match (empty means all)",
						Items:       &jsonschema.Schema{Type: "integer"},
					},
					"cluster": {
						Type:        "string",
						Description: "Cluster name filter",
					},
					"tags": {
						Type:        "array",
						Description: "Tag filters. Each filter has key, func (==, !=, in, not in, =~, !~), and value",
						Items: &jsonschema.Schema{
							Type: "object",
							Properties: map[string]*jsonschema.Schema{
								"key":   {Type: "string", Description: "Tag key"},
								"func":  {Type: "string", Description: "Operator: ==, !=, in, not in, =~, !~"},
								"value": {Type: "string", Description: "Tag value (for 'in'/'not in', space-separated values)"},
							},
						},
					},
					"cause": {
						Type:        "string",
						Description: "Reason/description for the mute",
					},
					"btime": {
						Type:        "integer",
						Description: "Start time Unix timestamp",
					},
					"etime": {
						Type:        "integer",
						Description: "End time Unix timestamp",
					},
					"severities": {
						Type:        "array",
						Description: "Severity levels to match (1=critical, 2=warning, 3=info). Empty means all.",
						Items:       &jsonschema.Schema{Type: "integer"},
					},
					"disabled": {
						Type:        "integer",
						Description: "Disabled status (0=enabled, 1=disabled)",
					},
					"mute_time_type": {
						Type:        "integer",
						Description: "Mute time type (0=time range, 1=periodic)",
					},
					"periodic_mutes": {
						Type:        "array",
						Description: "Periodic mute rules (when mute_time_type=1)",
						Items: &jsonschema.Schema{
							Type: "object",
							Properties: map[string]*jsonschema.Schema{
								"enable_stime":        {Type: "string", Description: "Start time in HH:MM format"},
								"enable_etime":        {Type: "string", Description: "End time in HH:MM format"},
								"enable_days_of_week": {Type: "string", Description: "Days of week (0-6, space-separated, 0=Sunday)"},
							},
						},
					},
				},
			},
		},
		toolset.MakeToolHandler(func(ctx context.Context, req *mcp.CallToolRequest, input app.CreateMuteInput) (*mcp.CallToolResult, error) {
			c := getClient(ctx)
			if c == nil {
				return toolset.NewToolResultError("failed to get n9e client from context"), nil
			}

			result, err := app.CreateMute(ctx, c, input)
			if err != nil {
				return toolset.NewToolResultError(err.Error()), nil
			}

			return toolset.MarshalResult(result), nil
		}),
	)
}

func updateMuteTool(getClient client.GetClientFunc) toolset.ServerTool {
	return toolset.NewServerTool(
		mcp.Tool{
			Name:        "update_mute",
			Description: "Update an existing alert mute/silence rule",
			Annotations: &mcp.ToolAnnotations{
				Title:           "Update Alert Mute",
				ReadOnlyHint:    false,
				DestructiveHint: toolset.BoolPtr(false),
			},
			InputSchema: &jsonschema.Schema{
				Type:     "object",
				Required: []string{"group_id", "mute_id", "cause", "btime", "etime"},
				Properties: map[string]*jsonschema.Schema{
					"group_id": {
						Type:        "integer",
						Description: "Business group ID",
					},
					"mute_id": {
						Type:        "integer",
						Description: "Alert mute ID to update",
					},
					"note": {
						Type:        "string",
						Description: "Note/title for the mute rule",
					},
					"cate": {
						Type:        "string",
						Description: "Category (e.g., prometheus, host, elasticsearch)",
					},
					"prod": {
						Type:        "string",
						Description: "Product type (e.g., metric, host, loki)",
					},
					"datasource_ids": {
						Type:        "array",
						Description: "Datasource IDs to match (empty means all)",
						Items:       &jsonschema.Schema{Type: "integer"},
					},
					"cluster": {
						Type:        "string",
						Description: "Cluster name filter",
					},
					"tags": {
						Type:        "array",
						Description: "Tag filters. Each filter has key, func (==, !=, in, not in, =~, !~), and value",
						Items: &jsonschema.Schema{
							Type: "object",
							Properties: map[string]*jsonschema.Schema{
								"key":   {Type: "string", Description: "Tag key"},
								"func":  {Type: "string", Description: "Operator: ==, !=, in, not in, =~, !~"},
								"value": {Type: "string", Description: "Tag value (for 'in'/'not in', space-separated values)"},
							},
						},
					},
					"cause": {
						Type:        "string",
						Description: "Reason/description for the mute",
					},
					"btime": {
						Type:        "integer",
						Description: "Start time Unix timestamp",
					},
					"etime": {
						Type:        "integer",
						Description: "End time Unix timestamp",
					},
					"severities": {
						Type:        "array",
						Description: "Severity levels to match (1=critical, 2=warning, 3=info). Empty means all.",
						Items:       &jsonschema.Schema{Type: "integer"},
					},
					"disabled": {
						Type:        "integer",
						Description: "Disabled status (0=enabled, 1=disabled)",
					},
					"mute_time_type": {
						Type:        "integer",
						Description: "Mute time type (0=time range, 1=periodic)",
					},
					"periodic_mutes": {
						Type:        "array",
						Description: "Periodic mute rules (when mute_time_type=1)",
						Items: &jsonschema.Schema{
							Type: "object",
							Properties: map[string]*jsonschema.Schema{
								"enable_stime":        {Type: "string", Description: "Start time in HH:MM format"},
								"enable_etime":        {Type: "string", Description: "End time in HH:MM format"},
								"enable_days_of_week": {Type: "string", Description: "Days of week (0-6, space-separated, 0=Sunday)"},
							},
						},
					},
				},
			},
		},
		toolset.MakeToolHandler(func(ctx context.Context, req *mcp.CallToolRequest, input app.UpdateMuteInput) (*mcp.CallToolResult, error) {
			c := getClient(ctx)
			if c == nil {
				return toolset.NewToolResultError("failed to get n9e client from context"), nil
			}

			result, err := app.UpdateMute(ctx, c, input)
			if err != nil {
				return toolset.NewToolResultError(err.Error()), nil
			}

			return toolset.MarshalResult(result), nil
		}),
	)
}
