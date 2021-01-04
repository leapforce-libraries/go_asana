package asana

import (
	"fmt"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Workspace stores Workspace from Asana
//
type Workspace struct {
	ID             string   `json:"gid"`
	Name           string   `json:"name"`
	ResourceType   string   `json:"resource_type"`
	EmailDomains   []string `json:"email_domains"`
	IsOrganization bool     `json:"is_organization"`
}

// GetWorkspacesByProjectID returns all workspaces for a specific project
//
func (i *Asana) GetWorkspaces() ([]Workspace, *errortools.Error) {
	return i.getWorkspacesInternal()
}

// getWorkspacesInternal is the generic function retrieving workspaces from Asana
//
func (i *Asana) getWorkspacesInternal() ([]Workspace, *errortools.Error) {
	urlStr := "workspaces?limit=%s%s&opt_fields=%s"
	limit := 100
	offset := ""
	//rowCount := limit
	batch := 0

	workspaces := []Workspace{}

	for batch == 0 || offset != "" {
		batch++
		//fmt.Printf("Batch %v for ProjectID %v\n", batch, projectID)

		urlPath := fmt.Sprintf(urlStr, strconv.Itoa(limit), offset, utilities.GetTaggedTagNames("json", Workspace{}))
		//fmt.Println(url)

		ts := []Workspace{}

		_, _, nextPage, e := i.Get(urlPath, &ts)
		if e != nil {
			return nil, e
		}

		for _, t := range ts {
			workspaces = append(workspaces, t)
		}

		//rowCount = len(ts)
		offset = ""
		if nextPage != nil {
			offset = fmt.Sprintf("&offset=%s", nextPage.Offset)
		}
	}

	if len(workspaces) == 0 {
		workspaces = nil
	}

	return workspaces, nil
}
