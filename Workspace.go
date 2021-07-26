package asana

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Workspace stores Workspace from Service
//
type Workspace struct {
	ID             string   `json:"gid"`
	Name           string   `json:"name"`
	ResourceType   string   `json:"resource_type"`
	EmailDomains   []string `json:"email_domains"`
	IsOrganization bool     `json:"is_organization"`
}

// GetWorkspaces returns all workspaces
//
func (service *Service) GetWorkspaces() ([]Workspace, *errortools.Error) {
	workspaces := []Workspace{}

	params := url.Values{}
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", Workspace{}))
	params.Set("limit", fmt.Sprintf("%v", limitDefault))

	for {
		_workspaces := []Workspace{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("workspaces?%s", params.Encode())),
			ResponseModel: &workspaces,
		}
		_, _, nextPage, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		workspaces = append(workspaces, _workspaces...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		params.Set("offset", nextPage.Offset)
	}

	return workspaces, nil
}