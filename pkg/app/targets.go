package app

import (
	"context"
	"net/url"
	"strconv"

	"github.com/n9e/n9e-mcp-server/pkg/client"
	"github.com/n9e/n9e-mcp-server/pkg/types"
)

// ListTargetsInput represents monitored objects list query parameters.
type ListTargetsInput struct {
	GroupIds      string `json:"gids,omitempty"`
	Query         string `json:"query,omitempty"`
	Limit         int    `json:"limit,omitempty"`
	Page          int    `json:"p,omitempty"`
	Downtime      int64  `json:"downtime,omitempty"`
	DatasourceIds string `json:"datasource_ids,omitempty"`
}

// ListTargets returns a paginated list of monitored targets.
func ListTargets(ctx context.Context, c *client.Client, input ListTargetsInput) (types.PageResp[types.Target], error) {
	params := url.Values{}
	if input.GroupIds != "" {
		params.Set("gids", input.GroupIds)
	}
	if input.Query != "" {
		params.Set("query", input.Query)
	}
	if input.Limit > 0 {
		params.Set("limit", strconv.Itoa(input.Limit))
	}
	if input.Page > 0 {
		params.Set("p", strconv.Itoa(input.Page))
	}
	if input.Downtime > 0 {
		params.Set("downtime", strconv.FormatInt(input.Downtime, 10))
	}
	if input.DatasourceIds != "" {
		params.Set("datasource_ids", input.DatasourceIds)
	}

	return client.DoGet[types.PageResp[types.Target]](c, ctx, "/api/n9e/targets", params)
}
