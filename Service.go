package asana

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	apiName      string = "Asana"
	apiURL       string = "https://app.asana.com/api/1.0"
	limitDefault uint64 = 100
)

// type
//
type Service struct {
	bearerToken string
	httpService *go_http.Service
}

// Response represents highest level of exactonline api response
//
type Response struct {
	Data     *json.RawMessage `json:"data"`
	NextPage *NextPage        `json:"next_page"`
}

// NextPage contains info for batched data retrieval
//
type NextPage struct {
	Offset string `json:"offset"`
	Path   string `json:"path"`
	URI    string `json:"uri"`
}

type ServiceConfig struct {
	BearerToken string
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.BearerToken == "" {
		return nil, errortools.ErrorMessage("Service BearerToken not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		bearerToken: serviceConfig.BearerToken,
		httpService: httpService,
	}, nil
}

func (service *Service) httpRequest(httpMethod string, requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", service.bearerToken))
	(*requestConfig).NonDefaultHeaders = &header

	// add error model
	errorResponse := ErrorResponse{}
	(*requestConfig).ErrorModel = &errorResponse

	request, response, e := service.httpService.HTTPRequest(httpMethod, requestConfig)
	if len(errorResponse.Errors) > 0 {
		messages := []string{}
		for _, message := range errorResponse.Errors {
			messages = append(messages, message.Message)
		}
		e.SetMessage(strings.Join(messages, "\n"))
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiURL, path)
}

func (service *Service) get(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *NextPage, *errortools.Error) {
	_response := Response{}

	_requestConfig := go_http.RequestConfig{
		URL:           requestConfig.URL,
		ResponseModel: &_response,
	}

	request, response, e := service.httpRequest(http.MethodGet, &_requestConfig)
	if e != nil {
		return request, response, nil, e
	}

	if _response.Data != nil {
		err := json.Unmarshal(*_response.Data, requestConfig.ResponseModel)
		if err != nil {
			e := errortools.ErrorMessage(err)
			e.SetRequest(request)
			e.SetResponse(response)
			return request, response, nil, e
		}
	}

	return request, response, _response.NextPage, e
}

func (service *Service) APIName() string {
	return apiName
}

func (service *Service) APIKey() string {
	return service.bearerToken
}

func (service *Service) APICallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) APIReset() {
	service.httpService.ResetRequestCount()
}
