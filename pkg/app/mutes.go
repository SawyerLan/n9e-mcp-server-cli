package app

import (
	"context"
	"fmt"

	"github.com/n9e/n9e-mcp-server/pkg/client"
	"github.com/n9e/n9e-mcp-server/pkg/types"
)

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
