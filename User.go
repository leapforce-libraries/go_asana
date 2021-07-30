package asana

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// User stores User from Service
//
type User struct {
	ID           string   `json:"gid"`
	ResourceType string   `json:"resource_type"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Photo        *Photo   `json:"photo"`
	Workspaces   []Object `json:"workspaces"`
}

// Photo stores Photo from Service
//
type Photo struct {
	Image128x128 string `json:"image_128x128"`
	Image21x21   string `json:"image_21x21"`
	Image27x27   string `json:"image_27x27"`
	Image36x36   string `json:"image_36x36"`
	Image60x60   string `json:"image_60x60"`
}

type GetUsersConfig struct {
	WorkspaceID *string
	TeamID      *string
}

// GetUsers returns all users
//
func (service *Service) GetUsers(config *GetUsersConfig) ([]User, *errortools.Error) {
	users := []User{}

	params := url.Values{}
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", User{}))
	params.Set("limit", fmt.Sprintf("%v", limitDefault))
	if config.WorkspaceID != nil {
		params.Set("workspace", *config.WorkspaceID)
	}
	if config.TeamID != nil {
		params.Set("team", *config.TeamID)
	}

	for {
		_users := []User{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("users?%s", params.Encode())),
			ResponseModel: &_users,
		}
		_, _, nextPage, e := service.getData(&requestConfig)
		if e != nil {
			return nil, e
		}

		users = append(users, _users...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		params.Set("offset", nextPage.Offset)
	}

	return users, nil
}

/*
// GetUsersByWorkspaceID returns all users for a specific workspace
//
func (service *Service) GetUsersByWorkspaceID(workspaceID string) ([]User, *errortools.Error) {
	return service.getUsersInternal(workspaceID)
}

// getUsersInternal is the generic function retrieving users from Service
//
func (service *Service) getUsersInternal(workspaceID string) ([]User, *errortools.Error) {
	urlStr := "users?workspace=%s&limit=%s%s&opt_fields=%s"
	limit := 100
	offset := ""
	//rowCount := limit
	batch := 0

	users := []User{}

	for batch == 0 || offset != "" {
		batch++
		//fmt.Printf("Batch %v for WorkspaceID %v\n", batch, workspaceID)

		urlPath := fmt.Sprintf(urlStr, workspaceID, strconv.Itoa(limit), offset, utilities.GetTaggedTagNames("json", User{}))
		//fmt.Println(url)

		ts := []User{}

		_, _, nextPage, e := service.Get(urlPath, &ts)
		if e != nil {
			return nil, e
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
*/
