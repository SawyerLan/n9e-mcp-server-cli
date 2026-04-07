package app

import (
	"context"

	"github.com/n9e/n9e-mcp-server/pkg/client"
	"github.com/n9e/n9e-mcp-server/pkg/types"
)

// ListBusiGroupsInput represents business groups list query parameters.
type ListBusiGroupsInput struct {
	Limit int `json:"limit,omitempty"`
	Page  int `json:"p,omitempty"`
}

// ListBusiGroups returns a paginated list of business groups the current user can access.
func ListBusiGroups(ctx context.Context, c *client.Client, input ListBusiGroupsInput) (types.PageResp[types.BusiGroup], error) {
	result, err := client.DoGet[[]types.BusiGroup](c, ctx, "/api/n9e/busi-groups", nil)
	if err != nil {
		return types.PageResp[types.BusiGroup]{}, err
	}

	items, total := SlicePage(result, input.Page, input.Limit)
	return types.PageResp[types.BusiGroup]{List: items, Total: total}, nil
}
