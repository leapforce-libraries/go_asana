package asana

import (
	"fmt"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Task stores Task from Asana
//
type Task struct {
	ID                    string            `json:"gid"`
	Name                  string            `json:"name"`
	ResourceType          string            `json:"resource_type"`
	ApprovalStatus        string            `json:"approval_status"`
	AssigneeStatus        string            `json:"assignee_status"`
	Completed             bool              `json:"completed"`
	CompletedAt           string            `json:"completed_at"`
	CompletedBy           CompactObject     `json:"completed_by"`
	CreatedAt             string            `json:"created_at"`
	Dependencies          []CompactObject   `json:"dependencies"`
	Dependents            []CompactObject   `json:"dependents"`
	DueAt                 string            `json:"due_at"`
	DueOn                 string            `json:"due_on"`
	External              External          `json:"external"`
	Hearted               bool              `json:"hearted"`
	Hearts                []UserList        `json:"hearts"`
	HTMLNotes             string            `json:"html_notes"`
	IsRenderedAsSeparator bool              `json:"is_rendered_as_separator"`
	Liked                 bool              `json:"liked"`
	Likes                 []UserList        `json:"likes"`
	Memberships           []Membership      `json:"memberships"`
	ModifiedAt            string            `json:"modified_at"`
	Notes                 string            `json:"notes"`
	NumHearts             int               `json:"num_hearts"`
	NumLikes              int               `json:"num_likes"`
	NumSubtasks           int               `json:"num_subtasks"`
	ResourceSubtype       string            `json:"resource_subtype"`
	StartOn               string            `json:"start_on"`
	Assignee              CompactObject     `json:"assignee"`
	CustomFields          []CustomFieldTask `json:"custom_fields"`
	Followers             []CompactObject   `json:"followers"`
	Parent                CompactObject     `json:"parent"`
	Projects              []CompactObject   `json:"projects"`
	Tags                  []CompactObject   `json:"tags"`
	Workspace             CompactObject     `json:"workspace"`
}

// GetTasksByProjectID returns all tasks for a specific project
//
func (i *Asana) GetTasksByProjectID(projectID string, projectIDsDone *[]string, modifiedSince *time.Time) ([]Task, *errortools.Error) {
	urlStr := "%stasks?project=%s&limit=%s%s&opt_fields=%s%s"
	limit := 100
	offset := ""
	//rowCount := limit
	batch := 0

	tasks := []Task{}

	for batch == 0 || offset != "" {
		batch++
		//fmt.Printf("Batch %v for ProjectID %v\n", batch, projectID)

		_modifiedSince := ""
		if modifiedSince != nil {
			_modifiedSince = fmt.Sprintf("&modified_since=%s", modifiedSince.Format("2006-01-02T15:04:05"))
		}

		url := fmt.Sprintf(urlStr, i.ApiURL, projectID, strconv.Itoa(limit), offset, utilities.GetTaggedTagNames("json", Task{}), _modifiedSince)
		//fmt.Println(url)

		nextPage, e := i.GetTasksInternal(url, &tasks, projectIDsDone, false)
		if e != nil {
			return nil, e
		}

		offset = ""
		if nextPage != nil {
			offset = fmt.Sprintf("&offset=%s", nextPage.Offset)
		}
	}

	if len(tasks) == 0 {
		tasks = nil
	}

	return tasks, nil
}

func (i *Asana) GetTasksInternal(url string, tasks *[]Task, projectIDsDone *[]string, subtasks bool) (*NextPage, *errortools.Error) {
	urlSubStr := "%stasks/%s/subtasks?opt_fields=%s"

	ts := []Task{}

	_, _, nextPage, e := i.Get(url, &ts)
	if e != nil {
		return nil, e
	}

	if tasks != nil {
		//tasks2 := *tasks
		//fmt.Println("len(tasks2)", len(tasks2))

		for _, t := range ts {
			taskFound := false

			if projectIDsDone != nil {
				if len(*projectIDsDone) > 0 {
				out:
					for _, proj := range t.Projects {
						for _, pid := range *projectIDsDone {
							if proj.ID == pid {
								taskFound = true
								break out
							}
						}
					}

					if taskFound {
						fmt.Println("duplicate TaskID: ", t.ID)
						continue
					}
				}
			}

			if !subtasks {
				if t.Parent.ResourceType != "project" && t.Parent.ResourceType != "" {
					fmt.Println("invalid Parent.ResourceType: ", t.Parent.ResourceType)
					continue
				}
			}

			*tasks = append(*tasks, t)

			if t.NumSubtasks > 0 {
				urlSub := fmt.Sprintf(urlSubStr, i.ApiURL, t.ID, utilities.GetTaggedTagNames("json", Task{}))
				_, e := i.GetTasksInternal(urlSub, tasks, projectIDsDone, true)
				if e != nil {
					return nil, e
				}
			}
		}
	}

	return nextPage, nil
}
