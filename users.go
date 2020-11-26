package asana

import (
	"fmt"
	"strconv"

	sentry "github.com/getsentry/sentry-go"
	errortools "github.com/leapforce-libraries/go_errortools"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// User stores User from Asana
//
type User struct {
	ID           string          `json:"gid"`
	Name         string          `json:"name"`
	ResourceType string          `json:"resource_type"`
	Email        string          `json:"email"`
	Photo        Photo           `json:"photo"`
	Workspaces   []CompactObject `json:"workspaces"`
}

// Photo stores Photo from Asana
//
type Photo struct {
	Image128x128 string `json:"image_128x128"`
	Image21x21   string `json:"image_21x21"`
	Image27x27   string `json:"image_27x27"`
	Image36x36   string `json:"image_36x36"`
	Image60x60   string `json:"image_60x60"`
}

// GetUsersByWorkspaceID returns all users for a specific workspace
//
func (i *Asana) GetUsersByWorkspaceID(workspaceID string) ([]User, *errortools.Error) {
	return i.GetUsersInternal(workspaceID)
}

// GetUsersInternal is the generic function retrieving users from Asana
//
func (i *Asana) GetUsersInternal(workspaceID string) ([]User, *errortools.Error) {
	urlStr := "%sworkspaces/%s/users?limit=%s%s&opt_fields=%s"
	limit := 100
	offset := ""
	//rowCount := limit
	batch := 0

	users := []User{}

	for batch == 0 || offset != "" {
		batch++
		//fmt.Printf("Batch %v for WorkspaceID %v\n", batch, workspaceID)

		url := fmt.Sprintf(urlStr, i.ApiURL, workspaceID, strconv.Itoa(limit), offset, utilities.GetTaggedFieldNames("json", User{}))
		//fmt.Println(url)

		ts := []User{}

		nextPage, response, e := i.Get(url, &ts)
		if e != nil {
			return nil, e
		}

		if response != nil {
			if response.Errors != nil {
				for _, e := range *response.Errors {
					message := fmt.Sprintf("Error for WorkspaceID %v: %v", workspaceID, e.Message)
					if i.IsLive {
						sentry.CaptureMessage(message)
					}
					fmt.Println(message)
				}
			}
		}

		for _, t := range ts {
			users = append(users, t)
		}

		//rowCount = len(ts)
		offset = ""
		if nextPage != nil {
			offset = fmt.Sprintf("&offset=%s", nextPage.Offset)
		}
	}

	if len(users) == 0 {
		users = nil
	}

	return users, nil
}

// GetUsers returns all users
//
/*
func (i *Asana) GetUsers() ([]User, error) {
	urlStr := "%susers?opt_fields=%s"

	users := []User{}

	url := fmt.Sprintf(urlStr, i.ApiURL, GetJSONTaggedFieldNames(User{}))
	//fmt.Println(url)

	_, response, err := i.Get(url, &users)
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

	return users, nil
}*/
