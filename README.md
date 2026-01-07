# Target Process MCP Server

A Model Context Protocol (MCP) server for interacting with Target Process, built in Go.

## Tools

| Tool | Description |
|------|-------------|
| `list_my_tickets` | List all tickets assigned to you that are currently In Progress |
| `get_ticket_details` | Get full details of a specific ticket by ID, including comments |

## Configuration

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `TP_ACCESS_TOKEN` | Yes | Your Target Process access token |
| `TP_USER_ID` | Yes | Your Target Process user ID |
| `TP_BASE_URL` | Yes | Target Process instance URL |

### Claude Code MCP Config

Add to your `.claude/mcp.json`:

```json
{
  "mcpServers": {
    "target-process": {
      "type": "stdio",
      "command": "C:\\path\\to\\target-process-mcp.exe",
      "env": {
        "TP_ACCESS_TOKEN": "your-access-token",
        "TP_USER_ID": "your-user-id",
        "TP_BASE_URL": "https://your-instance.tpondemand.com"
      }
    }
  }
}
```

## Building

```bash
go build -o target-process-mcp.exe
```

## Dependencies

- [mcp-go](https://github.com/mark3labs/mcp-go) - Go SDK for Model Context Protocol
