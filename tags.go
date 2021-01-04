package asana

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Tag stores Tag from Asana
//
type Tag struct {
	ID           string          `json:"gid"`
	Name         string          `json:"name"`
	ResourceType string          `json:"resource_type"`
	Color        string          `json:"color"`
	Followers    []CompactObject `json:"followers"`
	Workspace    CompactObject   `json:"workspace"`
}

// GetTags returns all tags
//
func (i *Asana) GetTags() ([]Tag, *errortools.Error) {
	urlStr := "tags?opt_fields=%s"

	tags := []Tag{}

	urlPath := fmt.Sprintf(urlStr, utilities.GetTaggedTagNames("json", Tag{}))
	//fmt.Println(url)

	_, _, _, e := i.Get(urlPath, &tags)
	if e != nil {
		return nil, e
	}

	return tags, nil
}
