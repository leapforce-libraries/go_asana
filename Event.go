package asana

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type EventsResponse struct {
	Data    *[]Event `json:"data"`
	Sync    *string  `json:"sync"`
	HasMore *bool    `json:"has_more"`
}

type EventsErrorResponse struct {
	Sync   *string `json:"sync"`
	Errors *[]*struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// Event stores Event from Service
//
type Event struct {
	User      Resource `json:"user"`
	CreatedAt string   `json:"created_at"`
	Type      string   `json:"type"`
	Action    string   `json:"action"`
	Resource  Resource `json:"resource"`
	Parent    Resource `json:"parent"`
	Change    struct {
		Field        string   `json:"field"`
		Action       string   `json:"action"`
		AddedValue   Resource `json:"added_value"`
		NewValue     Resource `json:"new_value"`
		RemovedValue Resource `json:"removed_value"`
	} `json:"change"`
}

type Resource struct {
	GID             string `json:"gid"`
	ResourceType    string `json:"resource_type"`
	Name            string `json:"name"`
	ResourceSubtype string `json:"resource_subtype"`
}

// GetEvents returns all events
//
func (service *Service) GetEventsByProject(projectID string, syncToken *string) (*[]Event, string, *errortools.Error) {
	eventsResponse := EventsResponse{}
	eventsErrorResponse := EventsErrorResponse{}

	params := url.Values{}
	params.Set("resource", projectID)
	if syncToken != nil {
		params.Set("sync", *syncToken)
	}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("events?%s", params.Encode())),
		ResponseModel: &eventsResponse,
		ErrorModel:    &eventsErrorResponse,
	}
	//fmt.Println(requestConfig.URL)
	_, response, e := service.get(&requestConfig)
	if syncToken == nil && response.StatusCode == 412 {
		if eventsErrorResponse.Sync == nil {
			return nil, "", errortools.ErrorMessage("eventsErrorResponse.Sync is nil")
		}
		return nil, *eventsErrorResponse.Sync, nil
	} else if e != nil {
		return nil, "", e
	}
	if eventsResponse.Sync == nil {
		return nil, "", errortools.ErrorMessage("eventsResponse.Sync is nil")
	}

	return eventsResponse.Data, *eventsResponse.Sync, nil
}
