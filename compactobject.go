package asana

// CompactObject stores 3 fielded compacted object from Asana
//
type CompactObject struct {
	ID           string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}

// CompactObject2 stores 2 fielded compacted object from Asana
//
type CompactObject2 struct {
	ID           string `json:"gid"`
	ResourceType string `json:"resource_type"`
}
