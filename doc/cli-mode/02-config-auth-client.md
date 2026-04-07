# 模块 02：配置、认证与客户端构建

## 1. 模块目标

本模块要解决 CLI 和 MCP 共用配置入口的问题，重点是：

- CLI 是否继续兼容环境变量认证
- flag 和环境变量的优先级怎么定义
- `*client.Client` 在什么地方构建

## 2. 当前现状

当前入口在 `cmd/n9e-mcp-server/main.go`，已有这些基础能力：

- `viper.SetEnvPrefix("N9E")`
- `viper.AutomaticEnv()`
- `--token`
- `--base-url`
- `--toolsets`
- `--read-only`

这说明 CLI 模式完全可以继承现有认证方式，不需要额外引入新机制。

## 3. 目标状态

CLI 和 MCP 共用一套配置语义：

- `--token` / `N9E_TOKEN`
- `--base-url` / `N9E_BASE_URL`
- `--toolsets` / `N9E_TOOLSETS`
- `--read-only` / `N9E_READ_ONLY`

建议新增 CLI 自有配置：

- `--output`
- `--timeout`
- `--yes`
- `--quiet`

## 4. 优先级规则

必须明确以下优先级：

1. 命令行 flag
2. 环境变量
3. 默认值

示例：

- 如果同时有 `--token` 和 `N9E_TOKEN`，以 `--token` 为准
- 如果没有传 `--base-url`，但设置了 `N9E_BASE_URL`，CLI 应自动使用环境变量

## 5. 推荐实现方式

建议新增统一配置结构，例如：

```go
type Config struct {
    Token           string
    BaseURL         string
    EnabledToolsets []string
    ReadOnly        bool
    LogFilePath     string
    Output          string
    Timeout         time.Duration
    Yes             bool
    Quiet           bool
}
```

建议新增一个集中配置构造入口，例如：

- `internal/config/config.go`

职责：

- 读取 `viper`
- 统一做默认值收敛
- 返回 `Config`

## 6. 客户端构建建议

建议不要在每个 CLI 命令里单独创建 `client.Client`。

推荐做法：

- 在 CLI 根初始化时创建
- 或通过工厂函数统一构造

例如：

```go
func NewClientFromConfig(cfg Config, userAgent string) (*client.Client, error)
```

这样后续如果要支持 `--timeout`，只需要调整这一个入口。

## 7. 风险点

### 风险 1：CLI 和 MCP 配置行为不一致

后果：

- 用户在 MCP 可用的环境变量，CLI 里却失效

规避：

- CLI 必须复用相同的 `viper` key 和 env 映射规则

### 风险 2：超时参数无处落地

后果：

- 文档有 `--timeout`，实现却是假的

规避：

- 在本模块决定 `pkg/client` 是否需要支持可配置超时

### 风险 3：命令中重复创建 client

后果：

- 结构散乱，后续很难统一 mock 和测试

规避：

- 把 client 构建收敛到根层或工厂层

## 8. 本模块建议产出

- 统一 `Config` 结构
- 统一配置读取逻辑
- 统一 client 构造逻辑
- 文档化优先级规则

## 9. 验收标准

- CLI 不传 `--token` 时，可直接读取 `N9E_TOKEN`
- CLI 不传 `--base-url` 时，可直接读取 `N9E_BASE_URL`
- CLI 与 MCP 的 key 命名保持一致
- `Config` 可被 CLI 和 `stdio` 模式共享
