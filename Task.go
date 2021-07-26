package asana

import (
	"fmt"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Task stores Task from Service
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

type GetTasksConfig struct {
	ProjectID     *string
	ModifiedSince *time.Time
}

// GetTasks returns all tasks
//
func (service *Service) GetTasks(config *GetTasksConfig) ([]Task, *errortools.Error) {
	tasks := []Task{}

	params := url.Values{}
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", Task{}))
	params.Set("limit", fmt.Sprintf("%v", limitDefault))
	if config.ProjectID != nil {
		params.Set("project", *config.ProjectID)
	}
	if config.ModifiedSince != nil {
		params.Set("modified_since", config.ModifiedSince.Format("2006-01-02T15:04:05"))
	}

	for {
		_tasks := []Task{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("tasks?%s", params.Encode())),
			ResponseModel: &tasks,
		}
		_, _, nextPage, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		tasks = append(tasks, _tasks...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		params.Set("offset", nextPage.Offset)
	}

	return tasks, nil
}

// GetSubTasks returns all subtasks of a parent task
//
func (service *Service) GetSubTasks(taskID string) ([]Task, *errortools.Error) {
	tasks := []Task{}

	params := url.Values{}
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", Task{}))

	for {
		_tasks := []Task{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("tasks/%s/subtasks?%s", taskID, params.Encode())),
			ResponseModel: &tasks,
		}
		_, _, nextPage, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		tasks = append(tasks, _tasks...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		params.Set("offset", nextPage.Offset)
	}

	return tasks, nil
}

/*
// GetTasksByProjectID returns all tasks for a specific project
//
func (service *Service) GetTasksByProjectID(projectID string, projectIDsDone *[]string, modifiedSince *time.Time) ([]Task, *errortools.Error) {
	urlStr := "tasks?project=%s&limit=%s%s&opt_fields=%s%s"
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

		urlPath := fmt.Sprintf(urlStr, projectID, strconv.Itoa(limit), offset, utilities.GetTaskgedTaskNames("json", Task{}), _modifiedSince)
		//fmt.Println(url)

		nextPage, e := service.getTasksInternal(urlPath, &tasks, projectIDsDone, false)
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

func (service *Service) getTasksInternal(url string, tasks *[]Task, projectIDsDone *[]string, subtasks bool) (*NextPage, *errortools.Error) {
	ts := []Task{}

	requestConfig := go_http.RequestConfig{
		URL:           url,
		ResponseModel: &ts,
	}
	_, _, nextPage, e := service.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	if tasks != nil {
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
				urlSub := fmt.Sprintf(urlSubStr, t.ID, utilities.GetTaskgedTaskNames("json", Task{}))
				_, e := service.getTasksInternal(urlSub, tasks, projectIDsDone, true)
				if e != nil {
					return nil, e
				}
			}
		}
	}

	return nextPage, nil
}
*/
