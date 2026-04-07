# 模块 03：共享业务层 pkg/app 设计

## 1. 模块目标

本模块用于把当前散落在 `pkg/api` 里的业务调用逻辑抽出来，形成 CLI 和 MCP 共用的执行层。

核心目标：

- 业务逻辑只写一遍
- CLI 不直接依赖 MCP 结果包装
- MCP handler 变成轻量适配层

## 2. 当前问题

当前 `pkg/api` 同时承担：

- 参数定义
- MCP schema
- 参数校验
- Nightingale API 调用
- MCP 结果包装

问题在于：

- CLI 如果直接复用，会被迫依赖 `mcp.CallToolResult`
- context 注入式拿 client 对 CLI 不友好
- 业务逻辑与协议适配逻辑耦合过深

## 3. 目标结构

建议新增：

- `pkg/app/alerts.go`
- `pkg/app/mutes.go`
- `pkg/app/targets.go`
- `pkg/app/users.go`
- `pkg/app/common.go`

后续再逐步扩到其他领域。

## 4. 接口设计原则

### 原则 1：不依赖 MCP

`pkg/app` 不应返回：

- `mcp.CallToolResult`
- `toolset.ServerTool`

只返回业务对象和 error。

### 原则 2：显式依赖 client

建议签名：

```go
func ListActiveAlerts(ctx context.Context, c *client.Client, input ListActiveAlertsInput) (types.PageResp[types.AlertCurEvent], error)
```

而不是从 `context` 里隐式取 client。

### 原则 3：入参尽量复用

短期可以复用现有输入结构体；中期建议把输入结构体也迁到 `pkg/app` 或更中立的包。

### 原则 4：保留校验逻辑

原来在 `pkg/api` 里的校验应迁到 `pkg/app`，避免 CLI 再抄一遍。

## 5. 推荐迁移顺序

如果目标是最快打通第一条 CLI 查询链路，可以先做一个更薄的例外切片：

- `busi-groups list`

原因：

- 当前逻辑最短
- 几乎不涉及复杂筛选条件
- 适合作为 `pkg/app` 提取和 CLI 接线的第一块样板

在这个最小切片跑通后，再按下面顺序扩大迁移范围：

第一批建议迁移：

- `alerts`
- `mutes`

原因：

- `alerts` 覆盖读接口典型模式
- `mutes` 同时覆盖读写模式

第二批再迁移：

- `targets`
- `users`
- `busi_groups`

## 6. MCP 层改造方式

改造后，`pkg/api` 每个 tool handler 只做三件事：

1. 输入解析
2. 调用 `pkg/app`
3. 将返回结果转成 MCP 文本

这样可以显著减轻 `pkg/api` 的复杂度。

## 7. 风险点

### 风险 1：迁移时改坏现有 MCP 行为

规避：

- 先做无行为变化重构
- 让 `pkg/api` 仍然输出相同 JSON 结果

### 风险 2：输入模型迁移过大

规避：

- 第一阶段先只迁“逻辑”，不强行迁所有 input struct

### 风险 3：公共层边界过宽

规避：

- `pkg/app` 只做用例编排
- 不把 CLI 输出逻辑塞进去

## 8. 本模块建议产出

- 首批 `pkg/app` 领域文件
- 明确的 service 函数签名
- 从 `pkg/api` 到 `pkg/app` 的调用替换

## 9. 验收标准

- `pkg/api` 中读写核心逻辑已迁到 `pkg/app`
- CLI 可以直接调用 `pkg/app`
- MCP 输出行为不变
