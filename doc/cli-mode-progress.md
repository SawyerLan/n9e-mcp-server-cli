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

- 当前建议模块：模块 03，共享业务层 `pkg/app`
- 当前最小目标：先打通 `busi-groups list` 的共享逻辑抽取，不要同时展开完整 CLI 树
- 本轮优先阅读：
  - `doc/cli-mode/03-shared-app-layer.md`
  - `pkg/api/busi_groups.go`
  - `pkg/types/types.go`
  - `internal/config/config.go`
- 本轮通常不需要先读：
  - 模块 05/06/07 文档
  - 整仓 `rg --files`
  - 与 `busi-groups` 无关的 API 文件
- 验证前置限制：
  - 当前会话环境缺少 `go` / `gofmt`
  - 正常开发流里不要反复搜索 Go 工具链位置，直接记录验证缺口即可

## 2. 当前总体状态

| 项目 | 状态 | 说明 |
| --- | --- | --- |
| CLI 模式方案梳理 | 已完成 | 总指南与模块化文档已建立 |
| CLI 模块化开发文档 | 已完成 | 已拆分为 01-07 模块文档 |
| 项目级开发约束 | 已完成 | 根目录 `AGENTS.md` 已创建 |
| 项目专用 skill | 已完成 | 已创建仓库内 `skills/n9e-cli-modular-development` |
| CLI 代码实现 | 进行中 | 已启动模块 02，完成共享配置与 client 入口抽取 |

## 3. 模块状态

| 模块 | 文件 | 状态 | 说明 |
| --- | --- | --- | --- |
| 01 | `doc/cli-mode/01-overview-and-roadmap.md` | 已完成 | 已明确 v1 范围、阶段与边界 |
| 02 | `doc/cli-mode/02-config-auth-client.md` | 已实现 | 已新增 `internal/config` 并让 `stdio` 复用统一配置与 client 构造，待补跑 Go 测试 |
| 03 | `doc/cli-mode/03-shared-app-layer.md` | 未开始 | 文档已完成，代码尚未抽取 `pkg/app` |
| 04 | `doc/cli-mode/04-cli-framework.md` | 未开始 | 文档已完成，`cli` 子命令骨架尚未接入 |
| 05 | `doc/cli-mode/05-readonly-commands-v1.md` | 未开始 | 文档已完成，只读命令尚未实现 |
| 06 | `doc/cli-mode/06-write-commands-and-safety.md` | 未开始 | 文档已完成，写命令尚未实现 |
| 07 | `doc/cli-mode/07-output-errors-testing.md` | 未开始 | 文档已完成，CLI 输出与退出码体系尚未接入 |

说明：

- 模块状态优先反映代码实现进度
- 模块 01 保持为规划完成，其余模块按代码推进情况更新

## 4. 里程碑状态

| 里程碑 | 状态 | 说明 |
| --- | --- | --- |
| A 基础可跑 | 进行中 | 配置、认证和 client 共享入口已就绪，CLI 骨架尚未接入 |
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
- 根据本次会话复盘，新增“快速启动卡片”并收紧下一步建议，减少重复读文档和无效检索

## 6. 下一步建议

建议按下面顺序推进代码实现：

1. 模块 03：抽出 `pkg/app`，先迁移 `busi-groups` 查询链路
2. 模块 04：接入 `cli` 根命令和基础骨架
3. 模块 05：先实现 `busi-groups list --output json`
4. 再迁移 `alerts` 与 `mutes` 相关共享逻辑

推荐第一个可执行编码目标：

- 打通 `n9e-mcp-server cli busi-groups list --output json`

原因：

- 它是最短查询链路
- 覆盖 CLI 入口、认证、共享层、输出四个关键点

## 7. 阻塞与待决策

当前无硬性产品阻塞，但存在本地验证缺口：

- 当前会话环境未提供 `go` / `gofmt`，无法执行 `go test ./...` 或格式化检查
- 模块 02 代码需在具备 Go 工具链的环境中补跑测试与 `stdio` 启动回归

后续编码前需要持续确认的事项：

- `pkg/app` 是否直接承载现有 input struct，还是同步迁移输入模型
- `--timeout` 已在 client 工厂层落地，后续只需决定何时接入 CLI flag
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
