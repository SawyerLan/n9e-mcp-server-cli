# 模块 07：输出、错误处理与测试

## 1. 模块目标

本模块用于统一 CLI 的用户体验和工程质量，避免后续每个命令各写各的输出与报错。

## 2. 输出格式建议

最低支持：

- `json`
- `table`
- `text`

建议第一阶段默认：

- 默认输出 `json`

原因：

- 最利于自动化
- 最容易保持一致
- 最不容易因为字段变化导致表格返工

## 3. Renderer 设计建议

建议抽象：

```go
type Renderer interface {
    Render(v any) error
}
```

推荐实现：

- `JSONRenderer`
- `TableRenderer`
- `TextRenderer`

## 4. 错误处理建议

需要统一三类问题：

- 参数错误
- API / 网络错误
- 业务错误

建议退出码：

- 参数错误：2
- API / 网络错误：1
- 成功：0

建议错误输出到 `stderr`。

## 5. API 错误展示建议

当前 `pkg/client` 已有 `APIError`，建议 CLI 直接利用。

错误展示至少包含：

- method
- path
- status
- request_id

这样便于排查服务端问题。

## 6. 表格输出建议

不要一开始给所有命令做表格。

推荐优先支持：

- `busi-groups list`
- `alerts active list`
- `targets list`

详情型命令建议继续用 `json`。

## 7. 测试建议

### 单元测试

建议覆盖：

- `pkg/app`
- `internal/config`
- `internal/cli/output`

### 集成测试

建议覆盖：

- CLI 参数解析
- 认证参数透传
- 典型命令执行路径

### 回归测试

建议确认：

- MCP `stdio` 现有行为未受影响

## 8. 风险点

### 风险 1：命令输出风格不统一

规避：

- 所有命令统一走 renderer

### 风险 2：错误信息太少

规避：

- 复用 `APIError`

### 风险 3：只改 CLI 没测 MCP

规避：

- 至少保留一轮 MCP 回归验证

## 9. 本模块建议产出

- Renderer 体系
- 统一错误处理入口
- 退出码约定
- 测试清单和测试布局

## 10. 验收标准

- CLI 正常输出走 `stdout`
- 错误输出走 `stderr`
- 参数错误与 API 错误能区分退出码
- 至少有基础单测和关键路径验证
