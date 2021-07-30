package asana

type EnumOption struct {
	ID           string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Color        string `json:"color"`
	Enabled      bool   `json:"enabled"`
	Name         string `json:"name"`
}
