package asana

import (
	"fmt"
	"net/url"
	"time"

	a_types "github.com/leapforce-libraries/go_asana/types"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Task stores Task from Service
//
type Task struct {
	ID                    string                  `json:"gid"`
	ResourceType          string                  `json:"resource_type"`
	Name                  string                  `json:"name"`
	ApprovalStatus        string                  `json:"approval_status"`
	AssigneeStatus        string                  `json:"assignee_status"`
	Completed             bool                    `json:"completed"`
	CompletedAt           string                  `json:"completed_at"`
	CompletedBy           Object                  `json:"completed_by"`
	CreatedAt             a_types.DateTimeString  `json:"created_at"`
	Dependencies          []ObjectCompact         `json:"dependencies"`
	Dependents            []ObjectCompact         `json:"dependents"`
	DueAt                 *a_types.DateTimeString `json:"due_at"`
	DueOn                 *a_types.DateString     `json:"due_on"`
	External              External                `json:"external"`
	HTMLNotes             string                  `json:"html_notes"`
	IsRenderedAsSeparator bool                    `json:"is_rendered_as_separator"`
	Liked                 bool                    `json:"liked"`
	Likes                 []UserList              `json:"likes"`
	Memberships           []Membership            `json:"memberships"`
	ModifiedAt            a_types.DateTimeString  `json:"modified_at"`
	Notes                 string                  `json:"notes"`
	NumLikes              int64                   `json:"num_likes"`
	NumSubtasks           int64                   `json:"num_subtasks"`
	ResourceSubtype       string                  `json:"resource_subtype"`
	StartOn               *a_types.DateString     `json:"start_on"`
	Assignee              Object                  `json:"assignee"`
	CustomFields          []CustomFieldTask       `json:"custom_fields"`
	Followers             []Object                `json:"followers"`
	Parent                *Object                 `json:"parent"`
	PermalinkURL          string                  `json:"permalink_url"`
	Projects              []Object                `json:"projects"`
	Tags                  []Object                `json:"tags"`
	Workspace             Object                  `json:"workspace"`
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
			ResponseModel: &_tasks,
		}
		_, _, nextPage, e := service.getData(&requestConfig)
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
			ResponseModel: &_tasks,
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
}*/

/*
type SortByField string

const (
	SortByFieldDueDate     SortByField = "due_date"
	SortByFieldCreatedAt   SortByField = "created_at"
	SortByFieldCompletedAt SortByField = "completed_at"
	SortByFieldLikes       SortByField = "likes"
	SortByFieldModifiedAt  SortByField = "modified_at"
)*/

type SearchTasksConfig struct {
	WorkspaceID      string
	CreatedAtBefore  *time.Time
	CreatedAtAfter   *time.Time
	ModifiedAtBefore *time.Time
	ModifiedAtAfter  *time.Time
	//SortByField      *SortByField
	//SortAscending    *bool
}

// GetTasks returns all tasks
//
func (service *Service) SearchTasks(config *SearchTasksConfig) ([]Task, *errortools.Error) {
	tasks := []Task{}

	params := url.Values{}
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", Task{}))
	params.Set("limit", fmt.Sprintf("%v", limitDefault))
	params.Set("sort_by", "created_at")
	params.Set("sort_ascending", "true")
	if config.CreatedAtBefore != nil {
		params.Set("created_at.before", config.CreatedAtBefore.Format(DateTimeLayout))
	}
	if config.CreatedAtAfter != nil {
		params.Set("created_at.after", config.CreatedAtAfter.Format(DateTimeLayout))
	}
	if config.ModifiedAtBefore != nil {
		params.Set("modified_at.before", config.ModifiedAtBefore.Format(DateTimeLayout))
	}
	if config.ModifiedAtAfter != nil {
		params.Set("modified_at.after", config.ModifiedAtAfter.Format(DateTimeLayout))
	}

	for {
		_tasks := []Task{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("workspaces/%s/tasks/search?%s", config.WorkspaceID, params.Encode())),
			ResponseModel: &_tasks,
		}
		_, _, _, e := service.getData(&requestConfig)
		if e != nil {
			return nil, e
		}

		tasks = append(tasks, _tasks...)

		if len(_tasks) == 0 {
			break
		}

		lastTask := _tasks[len(_tasks)-1]

		// assumption: no tasks with same CreatedAt at millisecond level
		params.Set("created_at.after", lastTask.CreatedAt.Value().Add(time.Millisecond).Format(DateTimeLayout))
	}

	return tasks, nil
}
