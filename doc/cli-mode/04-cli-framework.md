# 模块 04：CLI 命令骨架与命令树

## 1. 模块目标

本模块用于搭起 CLI 的外层框架，重点包括：

- 新增 `cli` 子命令
- 组织命令树
- 注入共享配置和 client
- 承接后续具体命令实现

## 2. 顶层策略

根命令默认行为不能变：

- `n9e-mcp-server` 仍默认进入 `stdio`

CLI 通过子命令接入：

- `n9e-mcp-server cli`

这是兼容现有 MCP 客户端的前提。

## 3. 推荐目录

建议新增：

```text
internal/cli/
  root.go
  context.go
  output/
  commands/
    alerts.go
    targets.go
    mutes.go
    users.go
```

## 4. 推荐命令树

建议按业务领域组织：

```text
n9e-mcp-server cli alerts ...
n9e-mcp-server cli targets ...
n9e-mcp-server cli mutes ...
n9e-mcp-server cli users ...
n9e-mcp-server cli busi-groups ...
n9e-mcp-server cli datasource ...
n9e-mcp-server cli notify-rules ...
n9e-mcp-server cli alert-subscribes ...
n9e-mcp-server cli event-pipelines ...
```

不要按 HTTP 动词拆，否则会很快变乱。

## 5. 全局参数建议

CLI 根命令建议持有：

- `--token`
- `--base-url`
- `--toolsets`
- `--read-only`
- `--output`
- `--timeout`
- `--yes`
- `--quiet`

## 6. 推荐上下文对象

建议为 CLI 建一个上下文结构，承载共享依赖，例如：

```go
type CLIContext struct {
    Config   config.Config
    Client   *client.Client
    Renderer output.Renderer
}
```

作用：

- 避免每个命令重复造对象
- 后续更容易测

## 7. 命令注册建议

每个领域一个注册函数，例如：

```go
func NewAlertsCommand(ctx *CLIContext) *cobra.Command
func NewMutesCommand(ctx *CLIContext) *cobra.Command
```

这样不会把所有命令都堆进 `main.go`。

## 8. 风险点

### 风险 1：命令树太早铺太大

规避：

- 先搭骨架和首批命令
- 其余领域先预留位置

### 风险 2：输出逻辑散到命令里

规避：

- 命令只负责拿数据
- 渲染交给 `output` 层

### 风险 3：命令初始化依赖混乱

规避：

- 统一使用 `CLIContext`

## 9. 本模块建议产出

- `cli` 子命令
- 基础命令树
- CLIContext
- 命令注册模式

## 10. 验收标准

- `n9e-mcp-server cli --help` 正常
- CLI 全局参数可见
- 首批命令可以挂到统一骨架上
