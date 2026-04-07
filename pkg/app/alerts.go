package app

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/n9e/n9e-mcp-server/pkg/client"
	"github.com/n9e/n9e-mcp-server/pkg/types"
)

// ListActiveAlertsInput represents active alerts query parameters.
type ListActiveAlertsInput struct {
	Hours         int64  `json:"hours,omitempty"`
	Stime         int64  `json:"stime,omitempty"`
	Etime         int64  `json:"etime,omitempty"`
	Severity      string `json:"severity,omitempty"`
	Query         string `json:"query,omitempty"`
	Cate          string `json:"cate,omitempty"`
	RuleProds     string `json:"rule_prods,omitempty"`
	DatasourceIds string `json:"datasource_ids,omitempty"`
	RuleId        int64  `json:"rid,omitempty"`
	BusiGroupId   int64  `json:"bgid,omitempty"`
	Limit         int    `json:"limit,omitempty"`
	Page          int    `json:"p,omitempty"`
}

// ListActiveAlerts returns a paginated list of active alert events.
func ListActiveAlerts(ctx context.Context, c *client.Client, input ListActiveAlertsInput) (types.PageResp[types.AlertCurEvent], error) {
	params := url.Values{}
	if input.Hours > 0 {
		params.Set("hours", strconv.FormatInt(input.Hours, 10))
	}
	if input.Stime > 0 {
		params.Set("stime", strconv.FormatInt(input.Stime, 10))
	}
	if input.Etime > 0 {
		params.Set("etime", strconv.FormatInt(input.Etime, 10))
	}
	if input.Severity != "" {
		params.Set("severity", input.Severity)
	}
	if input.Query != "" {
		params.Set("query", input.Query)
	}
	if input.Cate != "" {
		params.Set("cate", input.Cate)
	}
	if input.RuleProds != "" {
		params.Set("rule_prods", input.RuleProds)
	}
	if input.DatasourceIds != "" {
		params.Set("datasource_ids", input.DatasourceIds)
	}
	if input.RuleId > 0 {
		params.Set("rid", strconv.FormatInt(input.RuleId, 10))
	}
	if input.BusiGroupId > 0 {
		params.Set("bgid", strconv.FormatInt(input.BusiGroupId, 10))
	}
	if input.Limit > 0 {
		params.Set("limit", strconv.Itoa(input.Limit))
	}
	if input.Page > 0 {
		params.Set("p", strconv.Itoa(input.Page))
	}

	return client.DoGet[types.PageResp[types.AlertCurEvent]](c, ctx, "/api/n9e/alert-cur-events/list", params)
}

// ListHistoryAlertsInput represents historical alerts query parameters.
type ListHistoryAlertsInput struct {
	Hours         int64  `json:"hours,omitempty"`
	Stime         int64  `json:"stime,omitempty"`
	Etime         int64  `json:"etime,omitempty"`
	Severity      int    `json:"severity,omitempty"`
	IsRecovered   int    `json:"is_recovered,omitempty"`
	Query         string `json:"query,omitempty"`
	Cate          string `json:"cate,omitempty"`
	RuleProds     string `json:"rule_prods,omitempty"`
	DatasourceIds string `json:"datasource_ids,omitempty"`
	BusiGroupId   int64  `json:"bgid,omitempty"`
	Limit         int    `json:"limit,omitempty"`
	Page          int    `json:"p,omitempty"`
}

// ListHistoryAlerts returns a paginated list of historical alert events.
func ListHistoryAlerts(ctx context.Context, c *client.Client, input ListHistoryAlertsInput) (types.PageResp[types.AlertHisEvent], error) {
	params := url.Values{}
	if input.Hours > 0 {
		params.Set("hours", strconv.FormatInt(input.Hours, 10))
	}
	if input.Stime > 0 {
		params.Set("stime", strconv.FormatInt(input.Stime, 10))
	}
	if input.Etime > 0 {
		params.Set("etime", strconv.FormatInt(input.Etime, 10))
	}
	if input.Severity != 0 {
		params.Set("severity", strconv.Itoa(input.Severity))
	}
	if input.IsRecovered != 0 {
		params.Set("is_recovered", strconv.Itoa(input.IsRecovered))
	}
	if input.Query != "" {
		params.Set("query", input.Query)
	}
	if input.Cate != "" {
		params.Set("cate", input.Cate)
	}
	if input.RuleProds != "" {
		params.Set("rule_prods", input.RuleProds)
	}
	if input.DatasourceIds != "" {
		params.Set("datasource_ids", input.DatasourceIds)
	}
	if input.BusiGroupId > 0 {
		params.Set("bgid", strconv.FormatInt(input.BusiGroupId, 10))
	}
	if input.Limit > 0 {
		params.Set("limit", strconv.Itoa(input.Limit))
	}
	if input.Page > 0 {
		params.Set("p", strconv.Itoa(input.Page))
	}

	return client.DoGet[types.PageResp[types.AlertHisEvent]](c, ctx, "/api/n9e/alert-his-events/list", params)
}

// GetActiveAlert returns a single active alert event by ID.
func GetActiveAlert(ctx context.Context, c *client.Client, eventId int64) (types.AlertCurEvent, error) {
	path := fmt.Sprintf("/api/n9e/alert-cur-event/%d", eventId)
	return client.DoGet[types.AlertCurEvent](c, ctx, path, nil)
}

// GetHistoryAlert returns a single historical alert event by ID.
func GetHistoryAlert(ctx context.Context, c *client.Client, eventId int64) (types.AlertHisEvent, error) {
	path := fmt.Sprintf("/api/n9e/alert-his-event/%d", eventId)
	return client.DoGet[types.AlertHisEvent](c, ctx, path, nil)
}

// ListAlertRulesInput represents alert rules list query parameters.
type ListAlertRulesInput struct {
	GroupId int64 `json:"group_id"`
	Limit   int   `json:"limit,omitempty"`
	Page    int   `json:"p,omitempty"`
}

// ListAlertRules returns a paginated list of alert rules for a business group.
func ListAlertRules(ctx context.Context, c *client.Client, input ListAlertRulesInput) (types.PageResp[types.AlertRule], error) {
	path := fmt.Sprintf("/api/n9e/busi-group/%d/alert-rules", input.GroupId)
	result, err := client.DoGet[[]types.AlertRule](c, ctx, path, nil)
	if err != nil {
		return types.PageResp[types.AlertRule]{}, err
	}

	items, total := SlicePage(result, input.Page, input.Limit)
	return types.PageResp[types.AlertRule]{List: items, Total: total}, nil
}

// GetAlertRule returns a single alert rule by ID.
func GetAlertRule(ctx context.Context, c *client.Client, ruleId int64) (types.AlertRule, error) {
	path := fmt.Sprintf("/api/n9e/alert-rule/%d", ruleId)
	return client.DoGet[types.AlertRule](c, ctx, path, nil)
}
