# 模块 05：只读命令 V1 实施

## 1. 模块目标

本模块用于落首批 CLI 只读命令，目标不是一次把所有查询命令写完，而是优先验证架构。

## 2. 首批建议命令

建议优先实现：

- `busi-groups list`
- `alerts active list`
- `alerts active get`
- `targets list`
- `users list`
- `mutes list`

这批命令足以覆盖：

- 列表查询
- 详情查询
- 分页参数
- 过滤参数
- 本地分页和服务端分页两类场景

## 3. 命令落地建议

### `busi-groups list`

价值：

- 最简单的全量读命令
- 适合打通 CLI 最短路径

### `alerts active list`

价值：

- 覆盖时间范围、分页、severity 等典型参数

### `alerts active get`

价值：

- 覆盖详情型命令

### `targets list`

价值：

- 覆盖 query、gids、downtime 等筛选场景

### `users list`

价值：

- 覆盖用户列表型输出

### `mutes list`

价值：

- 为后续写命令铺路

## 4. 输出建议

第一阶段建议：

- 默认输出 `json`
- 可先不实现所有表格展示

如果要优先做表格，推荐只给下面命令做：

- `busi-groups list`
- `alerts active list`
- `targets list`

## 5. 分页注意事项

当前并不是所有接口都原生支持服务端分页。

例如这些命令存在本地分页行为：

- `list_mutes`
- `list_busi_groups`

CLI 应保留现有行为，不要擅自改成“全部服务端分页”。

## 6. 参数设计建议

保持与现有 API / MCP 参数语义尽量一致：

- 避免重新发明命名
- 允许适当做 CLI 友好的 flag 别名

例如：

- `--group-id`
- `--mute-id`
- `--severity`
- `--hours`

## 7. 风险点

### 风险 1：命令名与 MCP 语义偏差过大

规避：

- 保持领域与资源层级一致

### 风险 2：命令里直接写 HTTP 调用

规避：

- CLI 只调用 `pkg/app`

### 风险 3：过早优化 table 输出

规避：

- 先保证 `json` 稳定

## 8. 本模块建议产出

- 首批只读命令代码
- `json` 输出链路跑通
- 至少 1 到 2 个列表命令支持基础表格输出

## 9. 验收标准

- 首批命令都可执行
- 环境变量认证可直接使用
- 分页和筛选行为符合现有接口语义
- CLI 可作为 shell 脚本输入源使用
