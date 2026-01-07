package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// TPClient is the Target Process API client
type TPClient struct {
	config     Config
	httpClient *http.Client
}

// NewTPClient creates a new Target Process client
func NewTPClient() (*TPClient, error) {
	token := os.Getenv("TP_ACCESS_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TP_ACCESS_TOKEN environment variable is required")
	}

	baseURL := os.Getenv("TP_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("TP_BASE_URL environment variable is required")
	}

	userIDStr := os.Getenv("TP_USER_ID")
	if userIDStr == "" {
		return nil, fmt.Errorf("TP_USER_ID environment variable is required")
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("TP_USER_ID must be a valid integer: %w", err)
	}

	client := &TPClient{
		config: Config{
			BaseURL:     baseURL,
			AccessToken: token,
			UserID:      userID,
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return client, nil
}

// makeRequest makes an authenticated request to the TP API
func (c *TPClient) makeRequest(endpoint string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(c.config.BaseURL + "/api/v1" + endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := u.Query()
	q.Set("access_token", c.config.AccessToken)
	q.Set("format", "json")
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetLoggedUser returns the currently logged in user
func (c *TPClient) GetLoggedUser() (*LoggedUserResponse, error) {
	data, err := c.makeRequest("/Users/loggeduser", nil)
	if err != nil {
		return nil, err
	}

	var user LoggedUserResponse
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user: %w", err)
	}

	return &user, nil
}

// GetMyInProgressTickets returns all tickets assigned to the user that are In Progress
func (c *TPClient) GetMyInProgressTickets() ([]Assignable, error) {
	// Build the where clause for filtering
	// Use AssignedUser.Id for primary assignee only (not team/secondary assignments)
	whereClause := fmt.Sprintf(
		"(AssignedUser.Id eq %d) and (EntityState.Name eq 'In Progress')",
		c.config.UserID,
	)

	// Include fields we want
	includeFields := "[Id,Name,Description,EntityType,EntityState,Project,Priority]"

	params := map[string]string{
		"where":   whereClause,
		"include": includeFields,
		"take":    "100",
	}

	data, err := c.makeRequest("/Assignables", params)
	if err != nil {
		return nil, err
	}

	var response AssignablesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse assignables: %w", err)
	}

	return response.Items, nil
}

// GetTicketDetails returns detailed information about a specific ticket
func (c *TPClient) GetTicketDetails(ticketID int) (*Assignable, error) {
	includeFields := "[Id,Name,Description,EntityType,EntityState,Project,Priority,Comments[Id,Description,Owner[Id,FirstName,LastName,Email]]]"

	params := map[string]string{
		"include": includeFields,
	}

	endpoint := fmt.Sprintf("/Assignables/%d", ticketID)
	data, err := c.makeRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	var assignable Assignable
	if err := json.Unmarshal(data, &assignable); err != nil {
		return nil, fmt.Errorf("failed to parse assignable: %w", err)
	}

	return &assignable, nil
}

// FormatTicket formats a single ticket for display
func FormatTicket(a Assignable) string {
	var sb strings.Builder

	entityType := "Unknown"
	if a.EntityType != nil {
		entityType = a.EntityType.Name
	}

	state := "Unknown"
	if a.EntityState != nil {
		state = a.EntityState.Name
	}

	project := "Unknown"
	if a.Project != nil {
		project = a.Project.Name
	}

	priority := "None"
	if a.Priority != nil {
		priority = a.Priority.Name
	}

	sb.WriteString(fmt.Sprintf("## #%d: %s\n", a.ID, a.Name))
	sb.WriteString(fmt.Sprintf("**Type:** %s | **State:** %s | **Priority:** %s\n", entityType, state, priority))
	sb.WriteString(fmt.Sprintf("**Project:** %s\n", project))

	if a.Description != "" {
		// Truncate long descriptions
		desc := a.Description
		if len(desc) > 500 {
			desc = desc[:500] + "..."
		}
		sb.WriteString(fmt.Sprintf("\n**Description:**\n%s\n", desc))
	}

	if a.Comments != nil && len(a.Comments.Items) > 0 {
		sb.WriteString(fmt.Sprintf("\n**Comments (%d):**\n", len(a.Comments.Items)))
		for i, comment := range a.Comments.Items {
			if i >= 5 { // Limit to 5 most recent comments
				sb.WriteString(fmt.Sprintf("... and %d more comments\n", len(a.Comments.Items)-5))
				break
			}
			owner := "Unknown"
			if comment.Owner != nil {
				owner = fmt.Sprintf("%s %s", comment.Owner.FirstName, comment.Owner.LastName)
			}
			desc := comment.Description
			if len(desc) > 200 {
				desc = desc[:200] + "..."
			}
			sb.WriteString(fmt.Sprintf("- %s: %s\n", owner, desc))
		}
	}

	return sb.String()
}

// FormatTicketsList formats a list of tickets for display
func FormatTicketsList(tickets []Assignable) string {
	if len(tickets) == 0 {
		return "No in-progress tickets found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# My In-Progress Tickets (%d)\n\n", len(tickets)))

	for _, ticket := range tickets {
		sb.WriteString(FormatTicket(ticket))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

