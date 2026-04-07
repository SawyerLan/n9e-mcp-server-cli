package app

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/n9e/n9e-mcp-server/pkg/client"
	"github.com/n9e/n9e-mcp-server/pkg/types"
)

// ListUsersInput represents users list query parameters.
type ListUsersInput struct {
	Query string `json:"query,omitempty"`
	Limit int    `json:"limit,omitempty"`
	Page  int    `json:"p,omitempty"`
}

// ListUsers returns a paginated list of users.
func ListUsers(ctx context.Context, c *client.Client, input ListUsersInput) (types.PageResp[types.User], error) {
	params := url.Values{}
	if input.Query != "" {
		params.Set("query", input.Query)
	}
	if input.Limit > 0 {
		params.Set("limit", strconv.Itoa(input.Limit))
	}
	if input.Page > 0 {
		params.Set("p", strconv.Itoa(input.Page))
	}

	return client.DoGet[types.PageResp[types.User]](c, ctx, "/api/n9e/users", params)
}

// GetUser returns a single user by ID.
func GetUser(ctx context.Context, c *client.Client, userId int64) (types.User, error) {
	path := fmt.Sprintf("/api/n9e/user/%d/profile", userId)
	return client.DoGet[types.User](c, ctx, path, nil)
}

// ListUserGroupsInput represents user groups list query parameters.
type ListUserGroupsInput struct {
	Query string `json:"query,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

// ListUserGroups returns a list of user groups.
func ListUserGroups(ctx context.Context, c *client.Client, input ListUserGroupsInput) ([]types.UserGroup, error) {
	params := url.Values{}
	if input.Query != "" {
		params.Set("query", input.Query)
	}
	if input.Limit > 0 {
		params.Set("limit", strconv.Itoa(input.Limit))
	}

	return client.DoGet[[]types.UserGroup](c, ctx, "/api/n9e/user-groups", params)
}

// GetUserGroup returns a single user group with its members.
func GetUserGroup(ctx context.Context, c *client.Client, groupId int64) (types.UserGroupDetail, error) {
	path := fmt.Sprintf("/api/n9e/user-group/%d", groupId)
	return client.DoGet[types.UserGroupDetail](c, ctx, path, nil)
}
