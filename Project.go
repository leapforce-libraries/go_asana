package asana

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Project stores Project from Service
//
type Project struct {
	ID            string        `json:"gid"`
	Name          string        `json:"name"`
	ResourceType  string        `json:"resource_type"`
	Archived      bool          `json:"archived"`
	Color         string        `json:"color"`
	CreatedAt     string        `json:"created_at"`
	CurrentStatus ProjectStatus `json:"current_status"`
	//CustomFieldSettings string `json:"custom_field_settings"`
	DefaultView  string               `json:"default_view"`
	DueDate      string               `json:"due_date"`
	DueOn        string               `json:"due_on"`
	HTMLNotes    string               `json:"html_notes"`
	IsTemplate   bool                 `json:"is_template"`
	Members      []CompactObject      `json:"members"`
	ModifiedAt   string               `json:"modified_at"`
	Notes        string               `json:"notes"`
	Public       bool                 `json:"public"`
	StartOn      string               `json:"start_on"`
	Workspace    CompactObject        `json:"workspace"`
	CustomFields []CustomFieldProject `json:"custom_fields"`
	Followers    []CompactObject      `json:"followers"`
	Icon         string               `json:"icon"`
	Owner        CompactObject        `json:"owner"`
	Team         CompactObject        `json:"team"`
}

type GetProjectsConfig struct {
	WorkspaceID *string
}

// GetProjects returns all projects
//
func (service *Service) GetProjects(config *GetProjectsConfig) ([]Project, *errortools.Error) {
	projects := []Project{}

	params := url.Values{}
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", Project{}))

	if config != nil {
		if config.WorkspaceID != nil {
			params.Set("workspace", *config.WorkspaceID)
			params.Set("limit", fmt.Sprintf("%v", limitDefault)) // pagination only if workspace is specified
		}
	}

	for {
		_projects := []Project{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("projects?%s", params.Encode())),
			ResponseModel: &projects,
		}
		_, _, nextPage, e := service.get(&requestConfig)
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
