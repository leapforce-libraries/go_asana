package asana

import (
	"fmt"
	"strconv"

	sentry "github.com/getsentry/sentry-go"
)

// Task stores Task from Asana
//
type Task struct {
	ID           string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}

// GetTasksByProjectID returns all tasks for a specific project
//
func (i *Asana) GetTasksByProjectID(projectID int) ([]Task, error) {
	return i.GetTasksInternal(projectID)
}

// GetTasksInternal is the generic function retrieving tasks from Asana
//
func (i *Asana) GetTasksInternal(projectID int) ([]Task, error) {
	urlStr := "%stasks?project=%s&limit=%s%s"
	limit := 100
	offset := ""
	//rowCount := limit
	batch := 0

	tasks := []Task{}

	for batch == 0 || offset != "" {
		batch++
		//fmt.Printf("Batch %v for ProjectID %v\n", batch, projectID)

		url := fmt.Sprintf(urlStr, i.ApiURL, strconv.Itoa(projectID), strconv.Itoa(limit), offset)
		//fmt.Println(url)

		ts := []Task{}

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
		}

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
