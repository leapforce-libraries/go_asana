package asana

// ProjectStatus stores ProjectStatus from Asana
//
type ProjectStatus struct {
	ID           string        `json:"gid"`
	ResourceType string        `json:"resource_type"`
	Title        string        `json:"title"`
	Author       CompactObject `json:"author"`
	Color        string        `json:"color"`
	HTMLText     string        `json:"html_text"`
	ModifiedAt   string        `json:"modified_at"`
	Text         string        `json:"text"`
	CreatedAt    string        `json:"created_at"`
	CreatedBy    CompactObject `json:"created_by"`
}
