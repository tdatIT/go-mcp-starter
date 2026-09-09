# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

A minimal MCP (Model Context Protocol) server in Go that exposes two US National Weather Service tools (`get_forecast`, `get_alerts`) over the stdio transport. Built on `github.com/modelcontextprotocol/go-sdk` v1.7.0. Single `package main`, flat layout, no tests yet.

## Commands

```bash
go run .                  # run the server (reads MCP JSON-RPC on stdin — no interactive output)
go build -o mcp-go-starter-kit .
go vet ./...
gofmt -l .                # list unformatted files
go test ./...             # no tests exist yet
go test -run TestName ./  # single test, once tests are added
```

Note `go run .` (not `go run main.go`) — the tools, types, and HTTP helper live in sibling files and `go run main.go` fails to compile.

## Architecture

Four files, one per layer:

- **main.go** — server construction and registration. `mcp.AddTool(server, &mcp.Tool{...}, Handler)` derives each tool's JSON input schema by reflecting on the handler's third parameter, so a tool's public contract is defined by its input struct in `domain.go`, not by anything written in `main.go`. Also holds the `NWSAPIBase` / `UserAgent` constants.
- **domain.go** — all wire types. NWS API response shapes plus the tool input structs (`ForecastInput`, `AlertsInput`). `jsonschema:"..."` tags on input fields become the parameter descriptions the LLM sees — treat them as prompt text, not comments.
- **caller.go** — `makeNWSRequest[T]` (generic HTTP GET + JSON decode against the NWS API, sets the required `User-Agent` and `Accept: application/geo+json` headers) and the `formatAlert` / `formatPeriod` string formatters.
- **tool_exec.go** — the handler functions. Signature is `func(ctx, *mcp.CallToolRequest, Input) (*mcp.CallToolResult, any, error)`.

### Error convention

Handlers **never** return a non-nil `error`. Every failure path returns a `CallToolResult` with a human-readable `TextContent` message and `nil, nil` for the remaining returns, so the model receives an explanation instead of a protocol-level fault. Keep new handlers consistent with this.

`get_forecast` is a two-hop call: `/points/{lat},{lon}` yields a forecast URL, which is then fetched; it returns at most the next 5 periods.

## Conventions

- Go 1.26 toolchain; the code uses modern idioms (`range` over int, `cmp.Or`, generics). Prefer these over older equivalents.
- Anything written to stdout corrupts the stdio JSON-RPC stream. Log to stderr (`log` defaults there) — never `fmt.Println`.
- The `mcp-go-starter-kit` binary at the repo root is a build artifact; don't edit or commit it.
- Not a git repository yet.

## Running this server inside Claude Code

`.mcp.json` at the repo root registers this project's own server as `mcp-go-starter-kit` (`go run .`); `.claude/settings.json` pre-approves it via `enabledMcpjsonServers`. Restart Claude Code after changing either. Its tools then appear as `mcp__mcp-go-starter-kit__get_forecast` / `mcp__mcp-go-starter-kit__get_alerts`.
