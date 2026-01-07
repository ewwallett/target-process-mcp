# Target Process MCP Server

A Model Context Protocol (MCP) server for interacting with Target Process, built in Go.

## Tools

| Tool | Description |
|------|-------------|
| `list_my_tickets` | List all tickets assigned to you that are currently In Progress |
| `get_ticket_details` | Get full details of a specific ticket by ID, including comments |
| `summarize_my_work` | Get a summary of your current workload grouped by type |

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `TP_ACCESS_TOKEN` | Yes | - | Your Target Process access token |
| `TP_USER_ID` | No | `188` | Your Target Process user ID |
| `TP_BASE_URL` | No | `https://thettcgroup.tpondemand.com` | Target Process instance URL |

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
        "TP_USER_ID": "your-user-id"
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
