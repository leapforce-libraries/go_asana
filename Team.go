package asana

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Team stores Team from Service
//
type Team struct {
	Id              string `json:"gid"`
	ResourceType    string `json:"resource_type"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	HtmlDescription string `json:"html_description"`
	Organization    Object `json:"organization"`
	PermalinkUrl    string `json:"permalink_url"`
}

// GetTeams returns all teams
//
func (service *Service) GetTeamsByWorkspace(workspaceID string) ([]Team, *errortools.Error) {
	teams := []Team{}

	params := url.Values{}
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", Team{}))
	params.Set("limit", fmt.Sprintf("%v", limitDefault))

	for {
		_teams := []Team{}

		requestConfig := go_http.RequestConfig{
			Url:           service.url(fmt.Sprintf("organizations/%s/teams?%s", workspaceID, params.Encode())),
			ResponseModel: &_teams,
		}
		_, _, nextPage, e := service.getData(&requestConfig)
		if e != nil {
			return nil, e
		}

		teams = append(teams, _teams...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		params.Set("offset", nextPage.Offset)
	}

	return teams, nil
}

/*
// GetTeamsByWorkspaceID returns all teams for a specific team
//
func (service *Service) GetTeamsByWorkspaceID(workspaceID string) ([]Team, *errortools.Error) {
	return service.getTeamsInternal(workspaceID)
}

// getTeamsInternal is the generic function retrieving teams from Service
//
func (service *Service) getTeamsInternal(workspaceID string) ([]Team, *errortools.Error) {
	urlStr := "organizations/%s/teams?limit=%s%s&opt_fields=%s"
	limit := 100
	offset := ""
	//rowCount := limit
	batch := 0

	teams := []Team{}

	for batch == 0 || offset != "" {
		batch++
		//fmt.Printf("Batch %v for WorkspaceID %v\n", batch, workspaceID)

		urlPath := fmt.Sprintf(urlStr, workspaceID, strconv.Itoa(limit), offset, utilities.GetTeamgedTeamNames("json", Team{}))
		//fmt.Println(url)

		ts := []Team{}

		_, _, nextPage, e := service.Get(urlPath, &ts)
		if e != nil {
			return nil, e
		}

		for _, t := range ts {
			teams = append(teams, t)
		}

		//rowCount = len(ts)
		offset = ""
		if nextPage != nil {
			offset = fmt.Sprintf("&offset=%s", nextPage.Offset)
		}
	}

	if len(teams) == 0 {
		teams = nil
	}

	return teams, nil
}
*/
