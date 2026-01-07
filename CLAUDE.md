# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
go build -o target-process-mcp.exe
```

## Architecture

This is an MCP (Model Context Protocol) server that exposes Target Process API functionality as MCP tools. It runs over stdio and is designed to be used with Claude Code or other MCP-compatible clients.

**File structure:**
- `main.go` - MCP server setup, tool registration, and request handlers
- `tp_client.go` - Target Process API client (`TPClient`) and response formatting functions
- `types.go` - Go structs mapping to Target Process API entities

**Key dependency:** `github.com/mark3labs/mcp-go` - Go SDK for Model Context Protocol

**MCP tools exposed:**
- `list_my_tickets` - Lists in-progress tickets assigned to the configured user
- `get_ticket_details` - Fetches a specific ticket by ID with comments

## Configuration

All environment variables are required (set in MCP client config, not shell):
- `TP_ACCESS_TOKEN` - Target Process access token
- `TP_USER_ID` - User ID for filtering tickets (only this user's tickets are accessible)
- `TP_BASE_URL` - Target Process instance URL
