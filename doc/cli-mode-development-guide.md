# n9e-mcp-server CLI 模式开发指南

## 1. 文档目标

本文用于指导在当前 `n9e-mcp-server` 项目上新增 CLI 模式，目标是做到：

- 不破坏现有 MCP `stdio` 模式的行为和用法
- 最大化复用现有 Nightingale API 调用能力
- 让后续 CLI 开发具备清晰的分层、可测试性和可扩展性
- 为后续新增命令、输出格式和写操作留出稳定演进空间

## 2. 当前项目现状

### 2.1 入口与运行模式

当前入口位于 `cmd/n9e-mcp-server/main.go`，只有两类可见命令：

- 根命令：默认执行 `stdio` 模式
- `version`

现状说明：

- 根命令没有显式区分“服务模式”和“命令模式”，默认直接进入 MCP `stdio`
- 全局配置通过 `cobra + viper` 读取，核心参数有 `token`、`base-url`、`toolsets`、`read-only`、`log-file`
- 因此后续增加 CLI 模式时，最稳妥的方式是新增一个 `cli` 子命令，而不是改变根命令默认行为

### 2.2 核心分层

当前代码大致可以拆成四层：

1. 命令层
   `cmd/n9e-mcp-server/main.go`

2. 运行模式装配层
   `internal/server.go`

3. API 工具定义层
   `pkg/api/*.go`

4. HTTP 客户端与数据模型层
   `pkg/client/*.go`
   `pkg/types/types.go`

实际调用链路如下：

```text
cobra/viper
  -> runStdio
  -> internal.RunStdioServer
  -> internal.NewMCPServer
  -> api.DefaultToolsetGroup
  -> pkg/api 各 Toolset handler
  -> pkg/client.DoGet/DoPost/DoPut
  -> Nightingale HTTP API
```

### 2.3 MCP 模式当前职责

`internal/server.go` 主要负责：

- 创建 `pkg/client.Client`
- 创建 `mcp.Server`
- 通过 middleware 将 `Client` 注入 `context`
- 启用指定 toolset
- 注册所有 MCP tools
- 以 `mcp.StdioTransport` 运行服务

这个层本身没有业务逻辑，主要是装配层，未来可以继续保留。

### 2.4 业务能力当前分布

业务能力目前集中在 `pkg/api`，每个文件对应一个 toolset。当前共有 9 个 toolset、27 个 tools：

| Toolset | 说明 | 工具数量 | 写操作 |
| --- | --- | ---: | --- |
| `alerts` | 告警查询、规则查询 | 6 | 否 |
| `targets` | 监控对象查询 | 1 | 否 |
| `datasource` | 数据源查询 | 1 | 否 |
| `mutes` | 屏蔽规则查询与变更 | 4 | 是 |
| `busi_groups` | 业务组查询 | 1 | 否 |
| `notify_rules` | 通知规则查询 | 2 | 否 |
| `alert_subscribes` | 告警订阅查询 | 3 | 否 |
| `event_pipelines` | 事件流水线查询 | 5 | 否 |
| `users` | 用户与用户组查询 | 4 | 否 |

当前只有 `mutes` 暴露写操作：

- `create_mute`
- `update_mute`

### 2.5 现有实现模式

`pkg/api` 中每个工具基本都遵循同一个模式：

1. 定义输入结构体
2. 定义 MCP `Tool` 元数据和 JSON Schema
3. 在 handler 内做参数校验
4. 从 `context` 中拿 `*client.Client`
5. 组装 URL / Query / Body
6. 调用 `client.DoGet/DoPost/DoPut`
7. 使用 `toolset.MarshalResult` 转成 MCP 文本结果

这套模式对 MCP 很合适，但对 CLI 复用有一个明显问题：

- 业务逻辑和 MCP 适配逻辑混在一起了

例如：

- 参数定义既承担了“业务输入”作用，也承担了“JSON Schema”描述作用
- 返回值最后被固定包装成 `mcp.CallToolResult`
- 取 `Client` 的方式依赖 MCP middleware 注入的 `context`

这会导致 CLI 如果直接复用 `pkg/api`，代码会很别扭。

## 3. 当前代码对 CLI 模式的可复用点

CLI 模式不需要从零开始，以下部分都值得直接复用：

### 3.1 `pkg/client`

这是 CLI 最应该直接复用的部分，原因：

- 已经统一处理了 Token、BaseURL、User-Agent
- 已经封装了重试、超时、429、5xx 等行为
- 已经抽象出 `DoGet/DoPost/DoPut/DoDelete`
- `APIError` 已经带有路径、状态码、参数、请求体、请求 ID 等诊断信息

### 3.2 `pkg/types`

这是 CLI 输出层的基础模型，建议继续作为“远端 API DTO”使用。

### 3.3 `pkg/toolset` 的部分能力

以下能力可继续复用或借鉴：

- `DefaultToolsets`
- `ValidateTimeRange`
- `ValidatePagination`
- `ValidateSeverity`
- `SlicePage`

但以下部分不建议直接被 CLI 依赖：

- `ServerTool`
- `MakeToolHandler`
- `MarshalResult`
- `NewToolResultText`
- `NewToolResultError`

这些更偏向 MCP 适配层。

## 4. CLI 模式开发的主要问题

如果直接在当前结构上“边加命令边调 `pkg/api`”，后面大概率会出现以下问题：

### 4.1 业务逻辑重复

CLI 和 MCP 都会分别写一遍：

- 参数校验
- Query / Body 构造
- 调用 `client.DoXxx`
- 结果分页或格式转换

### 4.2 输出层耦合

当前 MCP 返回结果固定是 JSON 文本。CLI 需要的往往是：

- `json`
- `table`
- `yaml` 或 `text`

如果没有单独的结果层，CLI 命令会很快变成“又查数据又排版”的大函数。

### 4.3 配置入口不统一

当前配置解析只服务于 MCP 模式。CLI 模式接入后，如果继续把逻辑堆在 `main.go`，会让配置分支越来越乱。

### 4.4 写操作安全性不足

CLI 场景比 MCP 更容易直接被人手工执行，因此写命令至少需要考虑：

- `--read-only` 下禁止写操作
- 是否支持 `--yes` 跳过确认
- 默认输出是否要显示资源 ID 和摘要
- 错误时是否给出清晰的可操作提示

## 5. 推荐的目标架构

建议把后续结构调整成“协议适配层”和“业务执行层”分离。

### 5.1 建议分层

推荐新增一层应用服务层，例如：

```text
cmd/n9e-mcp-server/
  main.go

internal/
  server.go
  config/
  cli/
    root.go
    output/
    commands/

pkg/
  api/          # MCP 适配层
  app/          # 建议新增：CLI/MCP 共享业务逻辑
  client/
  toolset/
  types/
```

其中职责建议如下：

- `pkg/app`
  只做业务调用编排，不依赖 MCP，也不依赖 Cobra

- `pkg/api`
  只保留 MCP tool 的 schema、handler、annotation、结果包装

- `internal/cli`
  只处理命令树、flag、输出渲染、交互确认、exit code

### 5.2 推荐的数据流

目标调用链建议改成：

```text
CLI 命令 / MCP handler
  -> pkg/app 某个用例函数
  -> pkg/client.DoGet/DoPost/DoPut
  -> Nightingale API
```

这样 MCP 和 CLI 只是在“输入来源”和“结果展示方式”上不同，业务调用只写一遍。

### 5.3 建议新增的共享服务层

建议按领域拆分，而不是按协议拆分，例如：

- `pkg/app/alerts.go`
- `pkg/app/mutes.go`
- `pkg/app/users.go`
- `pkg/app/targets.go`

每个文件提供明确的方法，例如：

```go
func ListActiveAlerts(ctx context.Context, c *client.Client, input ListActiveAlertsInput) (types.PageResp[types.AlertCurEvent], error)
func GetMute(ctx context.Context, c *client.Client, groupID, muteID int64) (types.AlertMute, error)
func CreateMute(ctx context.Context, c *client.Client, input CreateMuteInput) (int64, error)
```

这里的关键点是：

- 返回真实业务对象，而不是 `mcp.CallToolResult`
- 接收显式 `*client.Client`，不要依赖 context 注入
- 保持输入结构体可被 CLI 和 MCP 共用

## 6. CLI 模式建议的命令设计

### 6.1 顶层策略

为保证兼容性，建议保留现有默认行为不变：

- `n9e-mcp-server` 仍然默认等价于 `n9e-mcp-server stdio`

新增：

- `n9e-mcp-server cli ...`

这样不会影响现有 Cursor / MCP 客户端配置。

### 6.2 推荐命令树

建议 CLI 命令按领域组织，而不是按 HTTP 动词组织：

```text
n9e-mcp-server cli alerts active list
n9e-mcp-server cli alerts active get --eid 123
n9e-mcp-server cli alerts history list
n9e-mcp-server cli alerts rules list --group-id 1
n9e-mcp-server cli alerts rules get --arid 10

n9e-mcp-server cli targets list
n9e-mcp-server cli datasource list
n9e-mcp-server cli busi-groups list

n9e-mcp-server cli mutes list --group-id 1
n9e-mcp-server cli mutes get --group-id 1 --mute-id 10
n9e-mcp-server cli mutes create ...
n9e-mcp-server cli mutes update ...

n9e-mcp-server cli notify-rules list
n9e-mcp-server cli notify-rules get --id 1

n9e-mcp-server cli alert-subscribes list --group-id 1
n9e-mcp-server cli alert-subscribes list-by-gids --gids 1,2
n9e-mcp-server cli alert-subscribes get --sid 3

n9e-mcp-server cli event-pipelines list
n9e-mcp-server cli event-pipelines get --id 1
n9e-mcp-server cli event-pipelines executions list --pipeline-id 1
n9e-mcp-server cli event-pipelines executions get --exec-id xxx

n9e-mcp-server cli users list
n9e-mcp-server cli users get --id 1
n9e-mcp-server cli user-groups list
n9e-mcp-server cli user-groups get --id 1
```

这种设计和现有 toolset 结构天然对应，迁移成本最低。

### 6.3 全局 Flag 建议

建议 CLI 与 MCP 共享以下全局参数：

- `--token`
- `--base-url`
- `--toolsets`
- `--read-only`

建议 CLI 专有增加：

- `--output json|table|text`
- `--timeout`
- `--yes`
- `--quiet`

其中：

- `--output json` 适合脚本调用
- `--output table` 适合人工查看
- `--yes` 只对写操作生效

### 6.4 认证与配置优先级

CLI 模式需要明确继承当前基于环境变量的认证方式，至少继续支持：

- `N9E_TOKEN`
- `N9E_BASE_URL`
- `N9E_TOOLSETS`
- `N9E_READ_ONLY`

建议优先级保持为：

1. 命令行 flag
2. 环境变量
3. 代码默认值

也就是说：

- `--token` 优先于 `N9E_TOKEN`
- `--base-url` 优先于 `N9E_BASE_URL`
- 未显式传 flag 时，CLI 应可直接读取现有环境变量完成认证

这点非常重要，因为当前 MCP 接入方式本身就是通过环境变量传递凭证，CLI 模式如果不兼容这一点，会破坏现有用户心智和自动化脚本习惯。

### 6.5 输出格式建议

CLI 最少支持两种输出：

1. `json`
   直接输出结构化数据，方便脚本消费

2. `table`
   给常见列表命令做简洁表格

建议第一阶段不要在所有命令上都追求漂亮表格，而是：

- 默认 `json`
- 高频 list 命令再补 `table`

这样能先把命令体系和共享逻辑搭起来。

## 7. 推荐的代码落地步骤

建议按“先抽共享逻辑，再接 CLI”的顺序推进。

### 阶段 1：抽离共享业务层

目标：不改变现有 MCP 行为，只做内部重构。

建议动作：

- 新增 `pkg/app`
- 将 `pkg/api` 中“参数校验 + 请求构造 + client.DoXxx”提取到 `pkg/app`
- `pkg/api` handler 改为调用 `pkg/app`
- 保持现有 tool name、schema、返回内容不变

完成这一阶段后，CLI 才有稳定复用基础。

### 阶段 2：增加 CLI 根命令和基础设施

建议新增：

- `internal/cli/root.go`
- `internal/cli/output/`
- `internal/config/`

目标：

- 统一创建 `*client.Client`
- 统一处理输出格式
- 统一处理错误和退出码

### 阶段 3：优先实现只读命令

建议首批 CLI 命令优先实现：

- `busi-groups list`
- `alerts active list`
- `alerts active get`
- `targets list`
- `users list`
- `mutes list`

原因：

- 能覆盖列表、详情、分页、过滤三类常见模式
- 不涉及交互确认，验证架构更轻

### 阶段 4：补写操作命令

建议后续再实现：

- `mutes create`
- `mutes update`

附加要求：

- `--read-only` 时直接拒绝执行
- 未传 `--yes` 时可增加确认提示
- 失败时原样打印 `APIError` 关键信息

### 阶段 5：补测试和文档

建议补充：

- `pkg/app` 单元测试
- `internal/cli/output` 单元测试
- 关键 CLI 命令的集成测试
- README 中新增 CLI 用法章节

## 8. 具体重构建议

### 8.1 输入结构体复用策略

当前很多输入结构体已经定义在 `pkg/api` 中。短期可以先复用，但中期更建议把这些输入结构体迁移到共享层，例如：

- `pkg/app/alerts.go` 中定义 `ListActiveAlertsInput`
- `pkg/api/alerts.go` 和 CLI 命令都引用它

否则后续 `pkg/api` 仍然会承担共享模型职责，边界不够清晰。

### 8.2 配置解析收敛

建议把配置收敛成一个统一结构，例如：

```go
type Config struct {
    Token           string
    BaseURL         string
    EnabledToolsets []string
    ReadOnly        bool
    LogFilePath     string
    Output          string
    Timeout         time.Duration
}
```

这样：

- `stdio` 模式使用其中一部分字段
- `cli` 模式使用另一部分字段
- 客户端创建逻辑也可以统一

### 8.3 错误处理建议

CLI 模式建议统一错误出口：

- 参数错误：返回 exit code 2
- API/网络错误：返回 exit code 1
- 成功：返回 0

输出建议：

- 默认打印简洁错误
- `APIError` 时带上 method、path、status、request_id
- 避免直接把大块请求体完整刷屏

### 8.4 输出层建议

建议抽象一个简单 renderer：

```go
type Renderer interface {
    Render(v any) error
}
```

实现：

- `JSONRenderer`
- `TableRenderer`
- `TextRenderer`

命令层只负责拿数据，不负责组织打印细节。

## 9. 开发时的注意事项

### 9.1 兼容性原则

必须确保以下行为不变：

- 根命令默认仍进入 `stdio`
- 现有环境变量名不变
- 现有 `stdio` 参数不变
- 现有 npm 包调用方式不变

当前 npm 包只是根据平台选择二进制并透传参数，因此只要默认行为不变，CLI 子命令的加入不会影响已有用户。

### 9.2 读写控制

当前 `read-only` 在 MCP 模式下通过“是否注册写工具”来限制。CLI 模式没有这个机制，因此必须单独在命令执行前检查。

建议实现统一守卫函数，例如：

```go
func EnsureWritable(cfg Config) error
```

### 9.3 分页行为

当前部分接口是服务端分页，部分是先取全量后在本地分页，例如：

- `list_alert_rules`
- `list_mutes`
- `list_notify_rules`
- `list_busi_groups`

CLI 设计时要明确保留这一行为，不要误以为所有接口都支持服务端分页。

### 9.4 日志与输出分离

`stdio` 模式当前会写日志，并把运行提示打印到 `stderr`。CLI 模式也应遵循：

- 正常结果输出到 `stdout`
- 日志和错误输出到 `stderr`

否则会影响 shell 管道和脚本处理。

## 10. 建议的首批改动清单

如果准备正式开工，建议按下面顺序提交：

1. 新增共享配置结构，收敛 `main.go` 中的配置读取
2. 新增 `pkg/app`，先迁移 `alerts` 和 `mutes`
3. 修改 `pkg/api` 让 MCP handler 走 `pkg/app`
4. 新增 `cli` 根命令和输出层
5. 落地首批只读命令
6. 再补 `mutes create/update`
7. 更新 README 和 npm 使用说明

## 11. 推荐的首批 CLI 验收标准

满足以下条件即可认为 CLI 模式第一版可用：

- `n9e-mcp-server` 默认行为不变
- `n9e-mcp-server cli busi-groups list --output json` 可正常工作
- `n9e-mcp-server cli alerts active list` 支持常见过滤参数
- `n9e-mcp-server cli mutes list --group-id 1` 可正常输出
- `--token` / `N9E_TOKEN`、`--base-url` / `N9E_BASE_URL` 在 CLI 下同样生效
- `--read-only` 下写命令被拒绝
- CLI 的错误输出可定位到具体 API 请求

## 12. 一句话结论

这个项目已经具备很好的 CLI 开发基础，真正需要补的不是 HTTP 能力，而是“把 MCP 适配层和共享业务层拆开”。推荐先抽 `pkg/app` 作为 CLI 与 MCP 的共同执行层，再增加 `cli` 子命令和输出层，这样后续命令数量再多，结构也不会失控。
