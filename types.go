package main

// EntityType represents the type of a TP entity (UserStory, Bug, Task, etc.)
type EntityType struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// EntityState represents the state of a TP entity (In Progress, Done, etc.)
type EntityState struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// Project represents a TP project
type Project struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// Priority represents the priority of a TP entity
type Priority struct {
	ID         int    `json:"Id"`
	Name       string `json:"Name"`
	Importance int    `json:"Importance"`
}

// Release represents a TP release
type Release struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// User represents a TP user
type User struct {
	ID        int    `json:"Id"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`
}

// Comment represents a comment on a TP entity
type Comment struct {
	ID          int    `json:"Id"`
	Description string `json:"Description"`
	Owner       *User  `json:"Owner,omitempty"`
}

// Attachment represents a file attached to a TP entity
type Attachment struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
	Uri  string `json:"Uri"`
}

// AttachmentsResponse represents the API response for attachments
type AttachmentsResponse struct {
	Items []Attachment `json:"Items"`
}

// CommentsCollection represents the comments collection in API response
type CommentsCollection struct {
	Items []Comment `json:"Items"`
}

// Assignable represents a work item (UserStory, Bug, Task, etc.)
type Assignable struct {
	ID          int                 `json:"Id"`
	Name        string              `json:"Name"`
	Description string              `json:"Description"`
	EntityType  *EntityType         `json:"EntityType,omitempty"`
	EntityState *EntityState        `json:"EntityState,omitempty"`
	Project     *Project            `json:"Project,omitempty"`
	Priority    *Priority           `json:"Priority,omitempty"`
	Release     *Release            `json:"Release,omitempty"`
	Comments    *CommentsCollection `json:"Comments,omitempty"`
}

// AssignablesResponse represents the API response for assignables
type AssignablesResponse struct {
	Items []Assignable `json:"Items"`
	Next  string       `json:"Next,omitempty"`
}

// LoggedUserResponse represents the response from /Users/loggeduser
type LoggedUserResponse struct {
	ID        int    `json:"Id"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`
}

// Config holds the configuration for the TP client
type Config struct {
	BaseURL     string
	AccessToken string
	UserID      int
}
