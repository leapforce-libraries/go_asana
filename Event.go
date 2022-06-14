package asana

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	a_types "github.com/leapforce-libraries/go_asana/types"
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
	Action string `json:"action"`
	Change struct {
		Action       string          `json:"action"`
		AddedValue   json.RawMessage `json:"added_value"`
		Field        string          `json:"field"`
		NewValue     json.RawMessage `json:"new_value"`
		RemovedValue json.RawMessage `json:"removed_value"`
	} `json:"change"`
	CreatedAt a_types.DateTimeString `json:"created_at"`
	Parent    Object                 `json:"parent"`
	Resource  Object                 `json:"resource"`
	User      Object                 `json:"user"`
}

// GetEvents returns all events
//
func (service *Service) GetEventsByProject(projectID string, syncToken *string) (*[]Event, *string, *http.Response, *errortools.Error) {
	eventsResponse := EventsResponse{}
	eventsErrorResponse := EventsErrorResponse{}

	params := url.Values{}
	params.Set("resource", projectID)
	if syncToken != nil {
		params.Set("sync", *syncToken)
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("events?%s", params.Encode())),
		ResponseModel: &eventsResponse,
		ErrorModel:    &eventsErrorResponse,
	}
	//fmt.Println(requestConfig.Url)
	_, response, e := service.httpRequest(&requestConfig)
	sync := eventsResponse.Sync
	if sync == nil {
		sync = eventsErrorResponse.Sync
	}
	return eventsResponse.Data, sync, response, e
}
