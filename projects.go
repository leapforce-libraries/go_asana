package asana

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	utilities "github.com/leapforce-libraries/go_utilities"
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
func (i *Asana) GetProjects() ([]Project, *errortools.Error) {
	urlStr := "%sprojects?opt_fields=%s"

	projects := []Project{}

	url := fmt.Sprintf(urlStr, i.ApiURL, utilities.GetTaggedFieldNames("json", Project{}))
	//fmt.Println(url)

	_, response, e := i.Get(url, &projects)
	if e != nil {
		return nil, e
	}

	i.captureErrors(response)

	return projects, nil
}
