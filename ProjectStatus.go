package asana

import (
	a_types "github.com/leapforce-libraries/go_asana/types"
)

// ProjectStatus stores ProjectStatus from Asana
//
type ProjectStatus struct {
	ID           string                 `json:"gid"`
	ResourceType string                 `json:"resource_type"`
	Title        string                 `json:"title"`
	Color        string                 `json:"color"`
	HTMLText     string                 `json:"html_text"`
	Text         string                 `json:"text"`
	Author       Object                 `json:"author"`
	CreatedAt    a_types.DateTimeString `json:"created_at"`
	CreatedBy    Object                 `json:"created_by"`
	ModifiedAt   a_types.DateTimeString `json:"modified_at"`
}
