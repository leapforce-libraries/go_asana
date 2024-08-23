package asana

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	a_types "github.com/leapforce-libraries/go_asana/types"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Task stores Task from Service
type Task struct {
	Id                    string                  `json:"gid"`
	ResourceType          string                  `json:"resource_type"`
	Name                  string                  `json:"name"`
	ApprovalStatus        string                  `json:"approval_status"`
	AssigneeStatus        string                  `json:"assignee_status"`
	Completed             bool                    `json:"completed"`
	CompletedAt           string                  `json:"completed_at"`
	CompletedBy           Object                  `json:"completed_by"`
	CreatedAt             a_types.DateTimeString  `json:"created_at"`
	CreatedBy             Object                  `json:"created_by"`
	Dependencies          []ObjectCompact         `json:"dependencies"`
	Dependents            []ObjectCompact         `json:"dependents"`
	DueAt                 *a_types.DateTimeString `json:"due_at"`
	DueOn                 *a_types.DateString     `json:"due_on"`
	External              External                `json:"external"`
	HtmlNotes             string                  `json:"html_notes"`
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
	PermalinkUrl          string                  `json:"permalink_url"`
	Projects              []Object                `json:"projects"`
	Tags                  []Object                `json:"tags"`
	Workspace             Object                  `json:"workspace"`
}

type GetTaskConfig struct {
	TaskId string
	Fields *string
}

// GetTask returns a specific tasks
func (service *Service) GetTask(config *GetTaskConfig) (*Task, *errortools.Error) {
	var values = url.Values{}

	var task Task

	if config.Fields != nil {
		values.Set("opt_fields", *config.Fields)
	}

	requestConfig := go_http.RequestConfig{
		Url:           service.url(fmt.Sprintf("tasks/%s?%s", config.TaskId, values.Encode())),
		ResponseModel: &task,
	}
	_, _, _, e := service.getData(&requestConfig)
	if e != nil {
		return nil, e
	}
	return &task, nil
}

type GetTasksConfig struct {
	ProjectID     *string
	ModifiedSince *time.Time
}

// GetTasks returns all tasks
func (service *Service) GetTasks(config *GetTasksConfig) ([]Task, *errortools.Error) {
	var tasks []Task

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
		var _tasks []Task

		requestConfig := go_http.RequestConfig{
			Url:           service.url(fmt.Sprintf("tasks?%s", params.Encode())),
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

// GetSubTasks returns all subtasks of a parent task
func (service *Service) GetSubTasks(taskID string) ([]Task, *errortools.Error) {
	var tasks []Task

	params := url.Values{}
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", Task{}))

	for {
		var _tasks []Task

		requestConfig := go_http.RequestConfig{
			Url:           service.url(fmt.Sprintf("tasks/%s/subtasks?%s", taskID, params.Encode())),
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
type SortByField string

const (
	SortByFieldDueDate     SortByField = "due_date"
	SortByFieldCreatedAt   SortByField = "created_at"
	SortByFieldCompletedAt SortByField = "completed_at"
	SortByFieldLikes       SortByField = "likes"
	SortByFieldModifiedAt  SortByField = "modified_at"
)*/

type SearchTasksConfig struct {
	WorkspaceID       string
	CreatedAtBefore   *time.Time
	CreatedAtAfter    *time.Time
	ModifiedAtBefore  *time.Time
	ModifiedAtAfter   *time.Time
	CompletedAtBefore *time.Time
	CompletedAtAfter  *time.Time
	Fields            *[]string
	Values            *url.Values
}

func (service *Service) SearchTasks(config *SearchTasksConfig) ([]Task, *errortools.Error) {
	var tasks []Task

	params := url.Values{}
	if config.Values != nil {
		params = *config.Values
	}
	if config.Fields != nil {
		params.Set("opt_fields", strings.Join(*config.Fields, ",")+",created_at")
	} else {
		params.Set("opt_fields", utilities.GetTaggedTagNames("json", Task{}))
	}
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
	if config.CompletedAtBefore != nil {
		params.Set("completed_at.before", config.CompletedAtBefore.Format(DateTimeLayout))
	}
	if config.CompletedAtAfter != nil {
		params.Set("completed_at.after", config.CompletedAtAfter.Format(DateTimeLayout))
	}

	for {
		var _tasks []Task

		requestConfig := go_http.RequestConfig{
			Url:           service.url(fmt.Sprintf("workspaces/%s/tasks/search?%s", config.WorkspaceID, params.Encode())),
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

type DeleteTaskConfig struct {
	TaskID string
}

func (service *Service) DeleteTask(config *DeleteTaskConfig) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		Url:    service.url(fmt.Sprintf("tasks/%s", config.TaskID)),
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}

type NewTask struct {
	ApprovalStatus  string                  `json:"approval_status,omitempty"`
	Assignee        string                  `json:"assignee,omitempty"`
	AssigneeSection string                  `json:"assignee_section,omitempty"`
	AssigneeStatus  string                  `json:"assignee_status,omitempty"`
	Completed       bool                    `json:"completed,omitempty"`
	CompletedBy     *Object                 `json:"completed_by,omitempty"`
	CustomFields    map[string]string       `json:"custom_fields,omitempty"`
	DueAt           *a_types.DateTimeString `json:"due_at,omitempty"`
	DueOn           *a_types.DateString     `json:"due_on,omitempty"`
	External        *External               `json:"external,omitempty"`
	Followers       []string                `json:"followers,omitempty"`
	HtmlNotes       string                  `json:"html_notes,omitempty"`
	Liked           bool                    `json:"liked,omitempty"`
	Name            string                  `json:"name,omitempty"`
	Notes           string                  `json:"notes,omitempty"`
	Parent          string                  `json:"parent,omitempty"`
	Projects        []string                `json:"projects,omitempty"`
	ResourceSubtype string                  `json:"resource_subtype,omitempty"`
	StartAt         *a_types.DateTimeString `json:"start_at,omitempty"`
	StartOn         *a_types.DateString     `json:"start_on,omitempty"`
	Tags            []string                `json:"tags,omitempty"`
	Workspace       string                  `json:"workspace,omitempty"`
}

func (service *Service) CreateTask(task *NewTask) (*Task, *errortools.Error) {
	if task == nil {
		return nil, nil
	}

	var createdTask Task

	requestConfig := go_http.RequestConfig{
		Method: http.MethodPost,
		Url:    service.url("tasks"),
		BodyModel: struct {
			Data *NewTask `json:"data"`
		}{task},
		ResponseModel: &createdTask,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &createdTask, nil
}

func (service *Service) UpdateTask(taskId int64, task *NewTask) (*Task, *errortools.Error) {
	if task == nil {
		return nil, nil
	}

	var updatedTask Task

	requestConfig := go_http.RequestConfig{
		Method: http.MethodPut,
		Url:    service.url(fmt.Sprintf("tasks/%v", taskId)),
		BodyModel: struct {
			Data *NewTask `json:"data"`
		}{task},
		ResponseModel: &updatedTask,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &updatedTask, nil
}
