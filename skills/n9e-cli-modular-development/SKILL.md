---
name: n9e-cli-modular-development
description: Guide modular development for the `n9e-mcp-server` project, especially when adding CLI mode, extracting shared `pkg/app` logic, organizing `internal/cli` commands, preserving MCP stdio compatibility, or updating the project's CLI development progress record. Use when working in `D:\new-code\n9e-mcp-server-cli` on any CLI-related module.
---

# N9E CLI Modular Development

## Overview

Use this skill to develop CLI mode for `n9e-mcp-server` one module at a time without breaking existing MCP behavior.
Treat the repo docs as the implementation map and always update the progress record after meaningful work.

## Required Reading

Before changing code, read these files in order:

1. `D:\new-code\n9e-mcp-server-cli\AGENTS.md`
2. `D:\new-code\n9e-mcp-server-cli\doc\cli-mode-index.md`
3. `D:\new-code\n9e-mcp-server-cli\doc\cli-mode-progress.md`

Then read only the module doc needed for the current task:

- `D:\new-code\n9e-mcp-server-cli\doc\cli-mode\01-overview-and-roadmap.md`
- `D:\new-code\n9e-mcp-server-cli\doc\cli-mode\02-config-auth-client.md`
- `D:\new-code\n9e-mcp-server-cli\doc\cli-mode\03-shared-app-layer.md`
- `D:\new-code\n9e-mcp-server-cli\doc\cli-mode\04-cli-framework.md`
- `D:\new-code\n9e-mcp-server-cli\doc\cli-mode\05-readonly-commands-v1.md`
- `D:\new-code\n9e-mcp-server-cli\doc\cli-mode\06-write-commands-and-safety.md`
- `D:\new-code\n9e-mcp-server-cli\doc\cli-mode\07-output-errors-testing.md`

## Workflow

### 1. Lock the module scope

Work on one module at a time.
Do not mix large changes across config, shared business logic, CLI framework, and write commands unless the current module explicitly requires it.

### 2. Preserve compatibility first

Keep these rules intact:

- Root command still defaults to `stdio`
- `N9E_TOKEN`, `N9E_BASE_URL`, `N9E_TOOLSETS`, `N9E_READ_ONLY` remain supported
- Existing MCP tool names and behavior stay unchanged unless explicitly requested

### 3. Prefer shared logic over duplication

When CLI needs existing MCP functionality:

- move reusable logic into a protocol-neutral layer such as `pkg/app`
- keep `pkg/api` focused on MCP schema and handler adaptation
- keep CLI rendering and command wiring under `internal/cli`

Do not copy Nightingale request-building logic into both MCP handlers and CLI commands.

### 4. Follow the recommended implementation order

Implement in this order unless the user explicitly redirects:

1. config/auth unification
2. shared `pkg/app` extraction
3. CLI framework
4. read-only commands
5. write commands
6. output polish and tests

### 5. Keep outputs and errors consistent

For CLI work:

- send normal output to `stdout`
- send errors to `stderr`
- prefer `json` correctness before `table`
- preserve API error context when possible

### 6. Update progress after meaningful work

After completing a task, update:

- `D:\new-code\n9e-mcp-server-cli\doc\cli-mode-progress.md`

Record:

- date
- module
- status change
- what changed
- next recommended step
- blockers or open decisions if any

### 7. Live verification against real N9E

After completing each module's code, **must** build the binary and run representative commands against the user's real N9E environment to verify correctness before committing.

Steps:

1. Build: `/home/corebug/.local/go/bin/go build -o /tmp/n9e-mcp-server ./cmd/n9e-mcp-server/`
2. Run at least 2-3 representative CLI commands from the module. `N9E_TOKEN` and `N9E_BASE_URL` are already set in the shell environment (via `~/.bashrc`), so just run:
   ```
   /tmp/n9e-mcp-server cli <command> ...
   ```
3. Confirm the output is valid JSON and contains expected data
4. If any command fails or returns unexpected results, fix the issue before committing
5. Record the verification results (which commands were tested, data summary) in the progress doc change log

Do not skip this step. If the N9E environment is unreachable, document the gap explicitly in the commit message.

### 8. Commit before handoff

Unless the user explicitly says not to commit, finishing a meaningful implementation task includes creating a Git commit.

Rules:

- commit after code/docs changes, progress updates, and live verification are in place
- use a focused commit that contains only the intended task files
- do not include unrelated dirty-worktree changes from the user
- if verification is blocked, still commit the work and mention the testing gap clearly in the final message and, when useful, in the commit message
- if you truly cannot make a safe commit, say so explicitly instead of silently leaving changes uncommitted

## Progress Record Rules

When updating `doc/cli-mode-progress.md`:

- keep statuses short and factual
- update module status and milestone status together when relevant
- append a concise entry to the change log
- do not erase prior history unless it is clearly wrong

## Done Criteria

A module is only "completed" when all of the following are true:

- code or docs for that module are in place
- behavior matches the module doc
- `go build ./...` and `go test ./...` pass
- live verification against real N9E (using env vars from `~/.bashrc`) was performed and results recorded
- `doc/cli-mode-progress.md` was updated
- a Git commit was created unless the user explicitly asked not to commit

## Avoid

- changing root default behavior away from `stdio`
- inventing a second config system
- mixing CLI rendering into shared business logic
- implementing write commands before read paths and guards are stable
- skipping progress updates after major development steps
