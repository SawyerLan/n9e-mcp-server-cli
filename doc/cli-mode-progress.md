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

## 2. 当前总体状态

| 项目 | 状态 | 说明 |
| --- | --- | --- |
| CLI 模式方案梳理 | 已完成 | 总指南与模块化文档已建立 |
| CLI 模块化开发文档 | 已完成 | 已拆分为 01-07 模块文档 |
| 项目级开发约束 | 已完成 | 根目录 `AGENTS.md` 已创建 |
| 项目专用 skill | 已完成 | 已创建仓库内 `skills/n9e-cli-modular-development` |
| CLI 代码实现 | 未开始 | 尚未进入实际编码阶段 |

## 3. 模块状态

| 模块 | 文件 | 状态 | 说明 |
| --- | --- | --- | --- |
| 01 | `doc/cli-mode/01-overview-and-roadmap.md` | 已完成 | 已明确 v1 范围、阶段与边界 |
| 02 | `doc/cli-mode/02-config-auth-client.md` | 已完成 | 已明确配置、认证、client 复用原则 |
| 03 | `doc/cli-mode/03-shared-app-layer.md` | 已完成 | 已明确 `pkg/app` 共享层方向 |
| 04 | `doc/cli-mode/04-cli-framework.md` | 已完成 | 已明确 CLI 骨架与命令树建议 |
| 05 | `doc/cli-mode/05-readonly-commands-v1.md` | 已完成 | 已明确首批只读命令范围 |
| 06 | `doc/cli-mode/06-write-commands-and-safety.md` | 已完成 | 已明确写命令与安全策略 |
| 07 | `doc/cli-mode/07-output-errors-testing.md` | 已完成 | 已明确输出、错误和测试基线 |

说明：

- 以上“已完成”指文档规划完成
- 代码实现状态仍需单独推进

## 4. 里程碑状态

| 里程碑 | 状态 | 说明 |
| --- | --- | --- |
| A 基础可跑 | 未开始 | 代码尚未开始 |
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

## 6. 下一步建议

建议按下面顺序推进代码实现：

1. 模块 02：统一配置、认证和 client 构建入口
2. 模块 03：抽出 `pkg/app`，先迁移 `alerts` 与 `mutes`
3. 模块 04：接入 `cli` 根命令和基础骨架
4. 模块 05：实现首批只读命令

推荐第一个可执行编码目标：

- 打通 `n9e-mcp-server cli busi-groups list --output json`

原因：

- 它是最短查询链路
- 覆盖 CLI 入口、认证、共享层、输出四个关键点

## 7. 阻塞与待决策

当前无硬性阻塞。

后续编码前需要持续确认的事项：

- `pkg/app` 是否直接承载现有 input struct，还是同步迁移输入模型
- `--timeout` 是否在第一版真正落地到 `pkg/client`
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
