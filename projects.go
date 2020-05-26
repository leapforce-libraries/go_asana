package asana

import (
	"fmt"

	sentry "github.com/getsentry/sentry-go"
)

// Project stores Project from Asana
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

// GetProjects returns all projects
//
func (i *Asana) GetProjects() ([]Project, error) {
	urlStr := "%sprojects?opt_fields=%s"

	projects := []Project{}

	url := fmt.Sprintf(urlStr, i.ApiURL, GetJSONTaggedFieldNames(Project{}))
	//fmt.Println(url)

	_, response, err := i.Get(url, &projects)
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

	return projects, nil
}
