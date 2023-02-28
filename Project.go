package asana

import (
	"fmt"
	"net/url"

	a_types "github.com/leapforce-libraries/go_asana/types"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Project stores Project from Service
type Project struct {
	Id                  string                 `json:"gid"`
	Name                string                 `json:"name"`
	ResourceType        string                 `json:"resource_type"`
	Archived            bool                   `json:"archived"`
	Color               *string                `json:"color"`
	CreatedAt           a_types.DateTimeString `json:"created_at"`
	CurrentStatus       ProjectStatus          `json:"current_status"`
	CustomFieldSettings []ObjectCompact        `json:"custom_field_settings"`
	DefaultView         string                 `json:"default_view"`
	DueOn               *a_types.DateString    `json:"due_on"`
	HtmlNotes           string                 `json:"html_notes"`
	IsTemplate          bool                   `json:"is_template"`
	Members             []Object               `json:"members"`
	ModifiedAt          a_types.DateTimeString `json:"modified_at"`
	Notes               string                 `json:"notes"`
	Public              bool                   `json:"public"`
	StartOn             *a_types.DateString    `json:"start_on"`
	Workspace           Object                 `json:"workspace"`
	CustomFields        []CustomFieldProject   `json:"custom_fields"`
	Followers           []Object               `json:"followers"`
	Icon                *string                `json:"icon"`
	Owner               Object                 `json:"owner"`
	PermalinkUrl        string                 `json:"permalink_url"`
	Team                Object                 `json:"team"`
}

type GetProjectsConfig struct {
	WorkspaceID *string
	Archived    *bool
}

// GetProjects returns all projects
func (service *Service) GetProjects(config *GetProjectsConfig) ([]Project, *errortools.Error) {
	projects := []Project{}

	params := url.Values{}
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", Project{}))

	if config != nil {
		if config.WorkspaceID != nil {
			params.Set("workspace", *config.WorkspaceID)
			params.Set("limit", fmt.Sprintf("%v", limitDefault)) // pagination only if workspace is specified
		}
		if config.Archived != nil {
			params.Set("archived", fmt.Sprintf("%v", *config.Archived))
		}
	}

	for {
		_projects := []Project{}

		requestConfig := go_http.RequestConfig{
			Url:           service.url(fmt.Sprintf("projects?%s", params.Encode())),
			ResponseModel: &_projects,
		}
		_, _, nextPage, e := service.getData(&requestConfig)
		if e != nil {
			return nil, e
		}

		projects = append(projects, _projects...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		params.Set("offset", nextPage.Offset)
	}

	return projects, nil
}
