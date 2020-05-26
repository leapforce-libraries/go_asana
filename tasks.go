package asana

import (
	"fmt"
	"strconv"

	sentry "github.com/getsentry/sentry-go"
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
func (i *Asana) GetTasksByProjectID(projectID string) ([]Task, error) {
	return i.GetTasksBatched(projectID)
}

func (i *Asana) GetTasksInternal(url string, tasks *[]Task) (*NextPage, error) {
	urlSubStr := "%stasks/%s/subtasks?opt_fields=%s"

	ts := []Task{}

	nextPage, response, err := i.Get(url, &ts)
	if err != nil {
		return nil, err
	}

	if response != nil {
		if response.Errors != nil {
			for _, e := range *response.Errors {
				message := fmt.Sprintf("Error for %v: %v", url, e.Message)
				if i.IsLive {
					sentry.CaptureMessage(message)
				}
				fmt.Println(message)
			}
		}
	}

	if tasks != nil {
		tasks2 := *tasks
		//fmt.Println("len(tasks2)", len(tasks2))

		for _, t := range ts {
			tasks2 = append(tasks2, t)

			if t.NumSubtasks > 0 {
				urlSub := fmt.Sprintf(urlSubStr, i.ApiURL, t.ID, GetJSONTaggedFieldNames(Task{}))
				i.GetTasksInternal(urlSub, &tasks2)
				//fmt.Println("t.NumSubtasks", t.NumSubtasks)
			}
		}
		//fmt.Println("len(tasks2)", len(tasks2))

		*tasks = tasks2
	}

	return nextPage, nil
}

// GetTasksInternal is the generic function retrieving tasks from Asana
//
func (i *Asana) GetTasksBatched(projectID string) ([]Task, error) {
	urlStr := "%stasks?project=%s&limit=%s%s&opt_fields=%s"
	limit := 100
	offset := ""
	//rowCount := limit
	batch := 0

	tasks := []Task{}

	for batch == 0 || offset != "" {
		batch++
		//fmt.Printf("Batch %v for ProjectID %v\n", batch, projectID)

		url := fmt.Sprintf(urlStr, i.ApiURL, projectID, strconv.Itoa(limit), offset, GetJSONTaggedFieldNames(Task{}))
		//fmt.Println(url)

		nextPage, err := i.GetTasksInternal(url, &tasks)
		if err != nil {
			return nil, err
		}

		//fmt.Println(len(tasks))

		/*ts := []Task{}

		nextPage, response, err := i.Get(url, &ts)
		if err != nil {
			return nil, err
		}

		if response != nil {
			if response.Errors != nil {
				for _, e := range *response.Errors {
					message := fmt.Sprintf("Error for ProjectID %v: %v", projectID, e.Message)
					if i.IsLive {
						sentry.CaptureMessage(message)
					}
					fmt.Println(message)
				}
			}
		}

		for _, t := range ts {
			tasks = append(tasks, t)
		}*/

		//rowCount = len(ts)
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
