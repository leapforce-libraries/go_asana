package asana

type Object struct {
	ID           string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}
