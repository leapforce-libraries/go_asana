package asana

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

type Webhook struct {
	Guid               string          `json:"guid"`
	ResourceType       string          `json:"resource_type"`
	Active             bool            `json:"active"`
	Resource           WebhookResource `json:"resource"`
	Target             string          `json:"target"`
	CreatedAt          string          `json:"created_at"`
	LastFailureAt      string          `json:"last_failure_at"`
	LastFailureContent string          `json:"last_failure_content"`
	LastSuccessAt      string          `json:"last_success_at"`
	Filters            []WebhookFilter `json:"filters"`
}

type WebhookResource struct {
	Guid         string `json:"guid"`
	ResourceType string `json:"resource_type"`
	Name         string `json:"name"`
}

type WebhookFilter struct {
	ResourceType    string   `json:"resource_type"`
	ResourceSubtype string   `json:"resource_subtype"`
	Action          string   `json:"action"`
	Fields          []string `json:"fields"`
}

type EstablishWebhookConfig struct {
	Resource string `json:"resource"`
}

func (service *Service) EstablishWebhook(config *EstablishWebhookConfig) (*Webhook, *errortools.Error) {
	var response struct {
		Data Webhook `json:"data"`
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.url("webhooks"),
		BodyModel:     config,
		ResponseModel: &response,
	}
	_, _, _, e := service.getData(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &response.Data, nil
}
