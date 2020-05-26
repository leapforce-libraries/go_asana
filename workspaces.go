package asana

import (
	"fmt"
	"strconv"

	sentry "github.com/getsentry/sentry-go"
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
func (i *Asana) GetWorkspaces() ([]Workspace, error) {
	return i.GetWorkspacesInternal()
}

// GetWorkspacesInternal is the generic function retrieving workspaces from Asana
//
func (i *Asana) GetWorkspacesInternal() ([]Workspace, error) {
	urlStr := "%sworkspaces?limit=%s%s&opt_fields=%s"
	limit := 100
	offset := ""
	//rowCount := limit
	batch := 0

	workspaces := []Workspace{}

	for batch == 0 || offset != "" {
		batch++
		//fmt.Printf("Batch %v for ProjectID %v\n", batch, projectID)

		url := fmt.Sprintf(urlStr, i.ApiURL, strconv.Itoa(limit), offset, GetJSONTaggedFieldNames(Workspace{}))
		//fmt.Println(url)

		ts := []Workspace{}

		nextPage, response, err := i.Get(url, &ts)
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
