# AGENTS.md

## Purpose

This file guides contributors and coding agents working in `n9e-mcp-server`.

Primary goals:

- Preserve existing MCP `stdio` behavior
- Keep the codebase easy to extend for CLI mode
- Reuse business logic instead of duplicating it across protocols
- Prefer small, verifiable, low-risk changes

## Project Summary

This project is a Go-based MCP server for Nightingale.

Current architecture:

- `cmd/n9e-mcp-server`
  Program entrypoint and Cobra/Viper wiring
- `internal`
  Runtime assembly for server mode
- `pkg/api`
  MCP tool registration, schemas, and handlers
- `pkg/client`
  Nightingale HTTP client and API error handling
- `pkg/toolset`
  Toolset grouping, validation helpers, and MCP result helpers
- `pkg/types`
  Shared DTOs for Nightingale API responses
- `npm`
  Platform package wrappers for distributing the binary

## Development Priorities

When making changes, prefer this order of importance:

1. Do not break existing MCP behavior
2. Keep configuration and auth backward compatible
3. Avoid business logic duplication
4. Keep new code easy to test
5. Keep CLI additions isolated from MCP transport concerns

## Non-Negotiable Compatibility Rules

- The default command behavior must remain equivalent to `n9e-mcp-server stdio`
- Existing env vars must keep working:
  - `N9E_TOKEN`
  - `N9E_BASE_URL`
  - `N9E_TOOLSETS`
  - `N9E_READ_ONLY`
- Existing MCP tool names and request semantics should not change unless explicitly requested
- Existing npm wrapper behavior must not be broken

## Configuration Rules

The project already uses `cobra` + `viper`.

Required behavior:

1. Command-line flags override environment variables
2. Environment variables override code defaults
3. CLI mode must support the same auth inputs as MCP mode

If you add CLI functionality, preserve support for:

- `--token` / `N9E_TOKEN`
- `--base-url` / `N9E_BASE_URL`
- `--toolsets` / `N9E_TOOLSETS`
- `--read-only` / `N9E_READ_ONLY`

## Architecture Rules

### Keep protocol adapters thin

`pkg/api` is the MCP adapter layer.

Do:

- Keep MCP schema definitions there
- Keep MCP tool annotations there
- Keep MCP result wrapping there

Do not:

- Keep large chunks of Nightingale API business logic there long term
- Put CLI-specific output logic there

### Prefer a shared application layer for reusable logic

For new CLI work, extract shared business logic into a protocol-neutral layer, recommended as `pkg/app`.

That layer should:

- Accept `context.Context`
- Accept an explicit `*client.Client`
- Return typed business objects from `pkg/types`
- Perform reusable validation and request building

That layer should not:

- Return `mcp.CallToolResult`
- Depend on Cobra
- Depend on MCP middleware context injection

### Keep CLI concerns isolated

CLI-specific responsibilities should stay in a CLI-focused package, recommended under `internal/cli`.

Examples:

- command tree
- flag definitions
- output rendering
- interactive confirmations
- exit code handling

## CLI Mode Guidance

If implementing CLI mode, use `n9e-mcp-server cli ...` as the entry path.

Recommended development order:

1. Unify config and auth loading
2. Extract shared business logic into `pkg/app`
3. Add CLI framework and root command
4. Implement read-only commands first
5. Add write commands after read path is stable
6. Add output polish and tests

Recommended first commands:

- `busi-groups list`
- `alerts active list`
- `alerts active get`
- `targets list`
- `users list`
- `mutes list`

Write commands should come later:

- `mutes create`
- `mutes update`

## Output and Error Handling

For CLI work:

- Normal command output goes to `stdout`
- Errors and logs go to `stderr`

Recommended output progression:

1. Make `json` output correct first
2. Add `table` only for high-value list commands
3. Keep detail commands on structured output unless there is a strong reason otherwise

Recommended exit codes:

- `0` success
- `1` API/network/runtime failure
- `2` invalid arguments or invalid local usage

When surfacing API errors, preserve useful context such as:

- method
- path
- status
- request id

## Read-Only and Write Safety

The current MCP server enforces read-only mode by not registering write tools.
CLI mode does not get that protection automatically.

If adding CLI write commands:

- Enforce `--read-only` before execution
- Support `--yes` for non-interactive confirmation bypass
- Prefer explicit success output with resource identifiers
- Be conservative with destructive or state-changing operations

## Testing Guidance

Before finishing a meaningful change, test the smallest relevant scope.

Preferred checks:

- `go test ./...`
- focused package tests when iterating

When changing MCP behavior, verify:

- stdio mode still starts
- tool registration still works
- existing env var auth still works

When changing CLI behavior, verify:

- config precedence
- output mode behavior
- exit codes
- read-only/write protections

## Editing Guidance

- Keep changes minimal and localized
- Reuse existing helpers in `pkg/client` and `pkg/toolset` where appropriate
- Do not duplicate API calling logic across MCP handlers and CLI commands
- Prefer adding shared helpers over copying request-building code
- Keep comments short and only where the code would otherwise be unclear

## Documentation Guidance

If you change CLI design or structure, update the related docs in `doc/`.

Relevant docs:

- `doc/cli-mode-development-guide.md`
- `doc/cli-mode-index.md`
- `doc/cli-mode/`

Use the modular docs as the source of truth for CLI planning and sequencing.

## Things To Avoid

- Changing default root command behavior away from `stdio`
- Hardcoding credentials
- Introducing a second, inconsistent config system
- Mixing MCP transport concerns into CLI code
- Mixing CLI rendering concerns into reusable business logic
- Rewriting stable packages without a clear reason
- Large refactors without preserving behavior and verifying incrementally

## Good Change Pattern

A good change in this repo usually looks like:

1. Adjust or add shared business logic
2. Keep protocol layer thin
3. Add or update focused tests
4. Update docs if behavior or structure changed

## If You Are Unsure

Default to the least risky path:

- preserve MCP behavior
- reuse existing config/auth behavior
- extract reusable logic instead of duplicating it
- implement read paths before write paths
