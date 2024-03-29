package asana

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Tag stores Tag from Service
//
type Tag struct {
	Id           string   `json:"gid"`
	ResourceType string   `json:"resource_type"`
	Color        string   `json:"color"`
	Followers    []Object `json:"followers"`
	Name         string   `json:"name"`
	PermalinkUrl string   `json:"permalink_url"`
	Workspace    Object   `json:"workspace"`
}

// GetTags returns all tags
//
func (service *Service) GetTagsByWorkspace(workspaceID string) ([]Tag, *errortools.Error) {
	tags := []Tag{}

	params := url.Values{}
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", Tag{}))
	params.Set("limit", fmt.Sprintf("%v", limitDefault))

	for {
		_tags := []Tag{}

		requestConfig := go_http.RequestConfig{
			Url:           service.url(fmt.Sprintf("workspaces/%s/tags?%s", workspaceID, params.Encode())),
			ResponseModel: &_tags,
		}
		_, _, nextPage, e := service.getData(&requestConfig)
		if e != nil {
			return nil, e
		}

		tags = append(tags, _tags...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		params.Set("offset", nextPage.Offset)
	}

	return tags, nil
}
