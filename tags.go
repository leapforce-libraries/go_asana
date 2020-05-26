package asana

import (
	"fmt"

	sentry "github.com/getsentry/sentry-go"
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
func (i *Asana) GetTags() ([]Tag, error) {
	urlStr := "%stags?opt_fields=%s"

	tags := []Tag{}

	url := fmt.Sprintf(urlStr, i.ApiURL, GetJSONTaggedFieldNames(Tag{}))
	//fmt.Println(url)

	_, response, err := i.Get(url, &tags)
	if err != nil {
		return nil, err
	}

	if response != nil {
		if response.Errors != nil {
			for _, e := range *response.Errors {
				message := fmt.Sprintf("Error in %v: %v", url, e.Message)
				if i.IsLive {
					sentry.CaptureMessage(message)
				}
				fmt.Println(message)
			}
		}
	}

	return tags, nil
}
