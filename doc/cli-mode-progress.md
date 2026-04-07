# CLI 模式开发进度记录

## 1. 说明

本文档用于记录 `n9e-mcp-server` 项目 CLI 模式的开发进度，目标是：

- 让开发状态对所有协作者可见
- 避免模块推进过程中丢失上下文
- 记录阶段性决策、已完成项、阻塞项和下一步建议

更新规则：

- 每次完成一个明确任务后更新
- 同时更新模块状态和变更记录
- 如果存在阻塞或待确认事项，写进“阻塞与待决策”

相关文档：

- [CLI 模式文档索引](./cli-mode-index.md)
- [CLI 模式开发总指南](./cli-mode-development-guide.md)

## 1.1 快速启动卡片

给下一位 agent 的最小启动信息：

- 当前建议模块：模块 05，只读命令扩展（alerts、targets 等）
- 当前最小目标：参照 `busi-groups list` 模式，实现 `alerts list` 或 `targets list` 命令
- 本轮优先阅读：
  - `doc/cli-mode/05-readonly-commands-v1.md`
  - `pkg/api/` 下对应领域的 MCP handler（提取共享逻辑到 `pkg/app`）
  - `internal/cli/commands/busi_groups.go`（参考模式）
- 本轮通常不需要先读：
  - 模块 06/07 文档
  - 整仓 `rg --files`
- Go 工具链：已确认可用（`/home/corebug/.local/go/bin/go`）

## 2. 当前总体状态

| 项目 | 状态 | 说明 |
| --- | --- | --- |
| CLI 模式方案梳理 | 已完成 | 总指南与模块化文档已建立 |
| CLI 模块化开发文档 | 已完成 | 已拆分为 01-07 模块文档 |
| 项目级开发约束 | 已完成 | 根目录 `AGENTS.md` 已创建 |
| 项目专用 skill | 已完成 | 已创建仓库内 `skills/n9e-cli-modular-development` |
| CLI 代码实现 | 进行中 | 模块 02-04 已完成，CLI 骨架与首个命令已就绪 |

## 3. 模块状态

| 模块 | 文件 | 状态 | 说明 |
| --- | --- | --- | --- |
| 01 | `doc/cli-mode/01-overview-and-roadmap.md` | 已完成 | 已明确 v1 范围、阶段与边界 |
| 02 | `doc/cli-mode/02-config-auth-client.md` | 已实现 | 已新增 `internal/config` 并让 `stdio` 复用统一配置与 client 构造，待补跑 Go 测试 |
| 03 | `doc/cli-mode/03-shared-app-layer.md` | 进行中 | busi-groups 共享逻辑已抽取到 `pkg/app`，其余领域待后续迁移 |
| 04 | `doc/cli-mode/04-cli-framework.md` | 已实现 | CLI 骨架、CLIContext、output 层、busi-groups list 命令均已就绪，已通过实际环境验证 |
| 05 | `doc/cli-mode/05-readonly-commands-v1.md` | 未开始 | 文档已完成，只读命令尚未实现 |
| 06 | `doc/cli-mode/06-write-commands-and-safety.md` | 未开始 | 文档已完成，写命令尚未实现 |
| 07 | `doc/cli-mode/07-output-errors-testing.md` | 未开始 | 文档已完成，CLI 输出与退出码体系尚未接入 |

说明：

- 模块状态优先反映代码实现进度
- 模块 01 保持为规划完成，其余模块按代码推进情况更新

## 4. 里程碑状态

| 里程碑 | 状态 | 说明 |
| --- | --- | --- |
| A 基础可跑 | 已完成 | 配置、认证、client、CLI 骨架、首个命令均已就绪并通过验证 |
| B 首批可用 | 未开始 | 只读命令尚未实现 |
| C 具备写能力 | 未开始 | 写命令尚未实现 |
| D 可交付 | 未开始 | 测试与 README 更新尚未开始 |

## 5. 已完成事项

### 2026-04-07

- 阅读项目结构、入口、MCP 运行方式、toolset 分层和 npm 包装层
- 新增 [CLI 模式开发总指南](./cli-mode-development-guide.md)
- 在总指南中补充 CLI 对 `N9E_TOKEN` / `N9E_BASE_URL` 的兼容要求
- 新增 [CLI 模式文档索引](./cli-mode-index.md)
- 新增 `doc/cli-mode/` 下 01-07 模块文档
- 新增项目级开发约束文件 `AGENTS.md`
- 新建项目专用 skill：`skills/n9e-cli-modular-development`
- 完成项目内 skill 结构校验
- 尝试清理用户目录下的临时 skill 副本，但未成功自动删除
- 新增 `internal/config` 统一配置结构、env/flag 绑定和默认值收敛
- 为 `pkg/client` 增加可选超时配置入口，并保持原有 `NewClient` 兼容
- 将 `cmd/n9e-mcp-server` 与 `internal/server` 切换到共享配置和 client 构建路径
- 为 `internal/config` 和 `pkg/client` 增加模块 02 相关测试
- 当前环境缺少 `go` / `gofmt`，未能执行 `go test ./...`，需在具备 Go 工具链的环境补跑验证
- 根据本次会话复盘，新增”快速启动卡片”并收紧下一步建议，减少重复读文档和无效检索
- 新增 `pkg/app` 包：`common.go`（SlicePage 通用分页）、`busi_groups.go`（ListBusiGroups 共享逻辑）
- 重构 `pkg/api/busi_groups.go`，MCP handler 改为调用 `pkg/app.ListBusiGroups`
- 新增 `pkg/app/common_test.go`，SlicePage 单元测试通过
- 确认 Go 工具链可用，`go build ./...` 和 `go test ./...` 全部通过

## 6. 下一步建议

建议按下面顺序推进代码实现：

1. 模块 03 继续 + 模块 05：逐个领域迁移共享逻辑到 `pkg/app`，同时实现对应 CLI 只读命令
2. 推荐首个目标：`alerts list` 或 `targets list`
3. 模块 06：写命令与安全守卫

推荐下一个可执行编码目标：

- 参照 busi-groups 模式，将 alerts 共享逻辑抽取到 `pkg/app`，并实现 `cli alerts list`

原因：

- CLI 骨架和命令注册模式已就绪
- 只需在 `pkg/app` 新增共享逻辑 + 在 `commands/` 新增命令文件

## 7. 阻塞与待决策

当前无硬性阻塞。

已确认事项：

- Go 工具链可用，`go build ./...` 和 `go test ./...` 均通过
- `pkg/app` 承载独立的 input struct（如 `ListBusiGroupsInput`），不复用 MCP 层定义
- `--timeout` 已在 client 工厂层落地，后续只需决定何时接入 CLI flag

待决策事项：

- `table` 输出第一版覆盖哪些命令
- 写命令第一版是否只支持必要 flag，还是同步支持结构化输入

## 8. 变更记录

### 2026-04-07

- 初始化 CLI 模式文档体系
- 初始化项目级开发约束
- 初始化项目专用 skill
- 完成项目内 skill 校验
- 记录用户目录临时副本待人工清理
- 初始化 CLI 模式进度记录文档
- 启动模块 02 代码实现，落地共享配置与 client 构建入口
- 优化 handoff 文档结构，补充快速启动信息
- 模块 03 首个切片完成：`pkg/app` busi-groups 共享逻辑抽取与 MCP 层适配
- 模块 04 完成：CLI 框架骨架接入
  - 新增 `internal/cli/root.go`：CLI 根命令，绑定 CLI-only flags（output/timeout/yes/quiet）
  - 新增 `internal/cli/commands/busi_groups.go`：CLIContext 定义 + busi-groups list 命令
  - 新增 `internal/cli/output/renderer.go`：JSON 输出渲染层
  - 在 `cmd/n9e-mcp-server/main.go` 注册 `cli` 子命令
  - `go build ./...` 和 `go test ./...` 全部通过
  - 已通过实际 N9E 环境验证：`cli busi-groups list` 正确返回数据并分页
