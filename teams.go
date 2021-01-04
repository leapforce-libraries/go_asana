package asana

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Team stores Team from Asana
//
type Team struct {
	ID              string        `json:"gid"`
	Name            string        `json:"name"`
	ResourceType    string        `json:"resource_type"`
	Description     string        `json:"description"`
	HTMLDescription string        `json:"html_description"`
	Organization    CompactObject `json:"organization"`
}

// GetTeamsByWorkspaceID returns all teams for a specific team
//
func (i *Asana) GetTeamsByWorkspaceID(workspaceID string) ([]Team, *errortools.Error) {
	return i.GetTeamsInternal(workspaceID)
}

// GetTeamsInternal is the generic function retrieving teams from Asana
//
func (i *Asana) GetTeamsInternal(workspaceID string) ([]Team, *errortools.Error) {
	urlStr := "%sorganizations/%s/teams?limit=%s%s&opt_fields=%s"
	limit := 100
	offset := ""
	//rowCount := limit
	batch := 0

	teams := []Team{}

	for batch == 0 || offset != "" {
		batch++
		//fmt.Printf("Batch %v for WorkspaceID %v\n", batch, workspaceID)

		url := fmt.Sprintf(urlStr, i.ApiURL, workspaceID, strconv.Itoa(limit), offset, utilities.GetTaggedTagNames("json", Team{}))
		//fmt.Println(url)

		ts := []Team{}

		_, _, nextPage, e := i.Get(url, &ts)
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
