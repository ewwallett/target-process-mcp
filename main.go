package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var tpClient *TPClient

func main() {
	// Initialize the TP client
	var err error
	tpClient, err = NewTPClient()
	if err != nil {
		// Log to stderr since stdout is used for MCP protocol
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize TP client: %v\n", err)
		fmt.Fprintf(os.Stderr, "The server will start but API calls will fail until TP_ACCESS_TOKEN is set.\n")
	}

	// Create the MCP server
	s := server.NewMCPServer(
		"target-process",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register tools
	registerTools(s)

	// Run the server over stdio
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func registerTools(s *server.MCPServer) {
	// Tool: list_my_tickets
	listTicketsTool := mcp.NewTool("list_my_tickets",
		mcp.WithDescription("List all Target Process tickets assigned to you that are currently In Progress. Returns full details including description and recent comments."),
	)
	s.AddTool(listTicketsTool, listMyTicketsHandler)

	// Tool: get_ticket_details
	getTicketTool := mcp.NewTool("get_ticket_details",
		mcp.WithDescription("Get full details of a specific Target Process ticket by its ID, including description and all comments."),
		mcp.WithNumber("ticket_id",
			mcp.Description("The numeric ID of the ticket (e.g., 12345)"),
			mcp.Required(),
		),
	)
	s.AddTool(getTicketTool, getTicketDetailsHandler)

}

func listMyTicketsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if tpClient == nil {
		return mcp.NewToolResultError("TP client not initialized. Please set TP_ACCESS_TOKEN environment variable."), nil
	}

	tickets, err := tpClient.GetMyInProgressTickets()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error fetching tickets: %v", err)), nil
	}

	return mcp.NewToolResultText(FormatTicketsList(tickets)), nil
}

func getTicketDetailsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if tpClient == nil {
		return mcp.NewToolResultError("TP client not initialized. Please set TP_ACCESS_TOKEN environment variable."), nil
	}

	// Extract ticket_id from arguments
	args := req.GetArguments()
	ticketIDRaw, ok := args["ticket_id"]
	if !ok {
		return mcp.NewToolResultError("ticket_id is required"), nil
	}

	var ticketID int
	switch v := ticketIDRaw.(type) {
	case float64:
		ticketID = int(v)
	case int:
		ticketID = v
	default:
		return mcp.NewToolResultError(fmt.Sprintf("ticket_id must be a number, got %T", ticketIDRaw)), nil
	}

	if ticketID == 0 {
		return mcp.NewToolResultError("ticket_id must be a positive number"), nil
	}

	ticket, err := tpClient.GetTicketDetails(ticketID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error fetching ticket %d: %v", ticketID, err)), nil
	}

	return mcp.NewToolResultText(FormatTicket(*ticket)), nil
}

