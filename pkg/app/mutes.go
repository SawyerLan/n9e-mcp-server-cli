package app

import (
	"context"
	"fmt"

	"github.com/n9e/n9e-mcp-server/pkg/client"
	"github.com/n9e/n9e-mcp-server/pkg/types"
)

// CreateMuteInput represents parameters for creating an alert mute.
type CreateMuteInput struct {
	GroupId       int64                `json:"group_id"`
	Note          string               `json:"note"`
	Cate          string               `json:"cate,omitempty"`
	Prod          string               `json:"prod,omitempty"`
	DatasourceIds []int64              `json:"datasource_ids,omitempty"`
	Cluster       string               `json:"cluster,omitempty"`
	Tags          []types.TagFilter    `json:"tags,omitempty"`
	Cause         string               `json:"cause"`
	Btime         int64                `json:"btime"`
	Etime         int64                `json:"etime"`
	Severities    []int                `json:"severities,omitempty"`
	Disabled      int                  `json:"disabled,omitempty"`
	MuteTimeType  int                  `json:"mute_time_type,omitempty"`
	PeriodicMutes []types.PeriodicMute `json:"periodic_mutes,omitempty"`
}

// UpdateMuteInput represents parameters for updating an alert mute.
type UpdateMuteInput struct {
	GroupId       int64                `json:"group_id"`
	MuteId        int64                `json:"mute_id"`
	Note          string               `json:"note"`
	Cate          string               `json:"cate,omitempty"`
	Prod          string               `json:"prod,omitempty"`
	DatasourceIds []int64              `json:"datasource_ids,omitempty"`
	Cluster       string               `json:"cluster,omitempty"`
	Tags          []types.TagFilter    `json:"tags,omitempty"`
	Cause         string               `json:"cause"`
	Btime         int64                `json:"btime"`
	Etime         int64                `json:"etime"`
	Severities    []int                `json:"severities,omitempty"`
	Disabled      int                  `json:"disabled,omitempty"`
	MuteTimeType  int                  `json:"mute_time_type,omitempty"`
	PeriodicMutes []types.PeriodicMute `json:"periodic_mutes,omitempty"`
}

// CreateMuteResult holds the result of a successful mute creation.
type CreateMuteResult struct {
	Id      int64  `json:"id"`
	Message string `json:"message"`
}

// UpdateMuteResult holds the result of a successful mute update.
type UpdateMuteResult struct {
	Id      int64  `json:"id"`
	Message string `json:"message"`
}

func muteBody(note, cate, prod, cluster, cause string, datasourceIds []int64, tags []types.TagFilter,
	btime, etime int64, severities []int, disabled, muteTimeType int, periodicMutes []types.PeriodicMute) map[string]any {
	return map[string]any{
		"note":           note,
		"cate":           cate,
		"prod":           prod,
		"datasource_ids": datasourceIds,
		"cluster":        cluster,
		"tags":           tags,
		"cause":          cause,
		"btime":          btime,
		"etime":          etime,
		"severities":     severities,
		"disabled":       disabled,
		"mute_time_type": muteTimeType,
		"periodic_mutes": periodicMutes,
	}
}

// CreateMute creates a new alert mute rule.
func CreateMute(ctx context.Context, c *client.Client, input CreateMuteInput) (CreateMuteResult, error) {
	if input.GroupId <= 0 {
		return CreateMuteResult{}, fmt.Errorf("group_id is required and must be positive")
	}
	if input.Cause == "" {
		return CreateMuteResult{}, fmt.Errorf("cause is required")
	}
	if input.MuteTimeType == 0 {
		if input.Btime <= 0 || input.Etime <= 0 {
			return CreateMuteResult{}, fmt.Errorf("btime and etime are required for time range mode")
		}
		if input.Btime >= input.Etime {
			return CreateMuteResult{}, fmt.Errorf("btime must be less than etime")
		}
	}

	body := muteBody(input.Note, input.Cate, input.Prod, input.Cluster, input.Cause,
		input.DatasourceIds, input.Tags, input.Btime, input.Etime,
		input.Severities, input.Disabled, input.MuteTimeType, input.PeriodicMutes)

	path := fmt.Sprintf("/api/n9e/busi-group/%d/alert-mutes", input.GroupId)
	id, err := client.DoPost[int64](c, ctx, path, body)
	if err != nil {
		return CreateMuteResult{}, err
	}

	return CreateMuteResult{Id: id, Message: "Alert mute created successfully"}, nil
}

// UpdateMute updates an existing alert mute rule.
func UpdateMute(ctx context.Context, c *client.Client, input UpdateMuteInput) (UpdateMuteResult, error) {
	if input.GroupId <= 0 {
		return UpdateMuteResult{}, fmt.Errorf("group_id is required and must be positive")
	}
	if input.MuteId <= 0 {
		return UpdateMuteResult{}, fmt.Errorf("mute_id is required and must be positive")
	}
	if input.Cause == "" {
		return UpdateMuteResult{}, fmt.Errorf("cause is required")
	}
	if input.MuteTimeType == 0 {
		if input.Btime <= 0 || input.Etime <= 0 {
			return UpdateMuteResult{}, fmt.Errorf("btime and etime are required for time range mode")
		}
		if input.Btime >= input.Etime {
			return UpdateMuteResult{}, fmt.Errorf("btime must be less than etime")
		}
	}

	body := muteBody(input.Note, input.Cate, input.Prod, input.Cluster, input.Cause,
		input.DatasourceIds, input.Tags, input.Btime, input.Etime,
		input.Severities, input.Disabled, input.MuteTimeType, input.PeriodicMutes)

	path := fmt.Sprintf("/api/n9e/busi-group/%d/alert-mute/%d", input.GroupId, input.MuteId)
	_, err := client.DoPut[any](c, ctx, path, body)
	if err != nil {
		return UpdateMuteResult{}, err
	}

	return UpdateMuteResult{Id: input.MuteId, Message: "Alert mute updated successfully"}, nil
}

// ListMutesInput represents alert mutes list query parameters.
type ListMutesInput struct {
	GroupId int64 `json:"group_id"`
	Limit   int   `json:"limit,omitempty"`
	Page    int   `json:"p,omitempty"`
}

// ListMutes returns a paginated list of alert mutes for a business group.
func ListMutes(ctx context.Context, c *client.Client, input ListMutesInput) (types.PageResp[types.AlertMute], error) {
	path := fmt.Sprintf("/api/n9e/busi-group/%d/alert-mutes", input.GroupId)
	result, err := client.DoGet[[]types.AlertMute](c, ctx, path, nil)
	if err != nil {
		return types.PageResp[types.AlertMute]{}, err
	}

	items, total := SlicePage(result, input.Page, input.Limit)
	return types.PageResp[types.AlertMute]{List: items, Total: total}, nil
}

// GetMute returns a single alert mute by ID.
func GetMute(ctx context.Context, c *client.Client, groupId, muteId int64) (types.AlertMute, error) {
	path := fmt.Sprintf("/api/n9e/busi-group/%d/alert-mute/%d", groupId, muteId)
	return client.DoGet[types.AlertMute](c, ctx, path, nil)
}
