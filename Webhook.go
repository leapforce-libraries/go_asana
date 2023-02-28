package asana

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
)

type Webhook struct {
	Gid                string          `json:"gid"`
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
	Gid          string `json:"gid"`
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
	Filters  []EstablishWebhookConfigFilter `json:"filters"`
	Resource string                         `json:"resource"`
	Target   string                         `json:"target"`
}

type EstablishWebhookConfigFilter struct {
	Action       string   `json:"action"`
	Fields       []string `json:"fields"`
	ResourceType string   `json:"resource_type"`
}

func (service *Service) EstablishWebhook(config *EstablishWebhookConfig) (*Webhook, *errortools.Error) {
	if config == nil {
		return nil, nil
	}

	var response struct {
		Data Webhook `json:"data"`
	}

	requestConfig := go_http.RequestConfig{
		Method: http.MethodPost,
		Url:    service.url("webhooks"),
		BodyModel: struct {
			Data EstablishWebhookConfig `json:"data"`
		}{*config},
		ResponseModel: &response,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &response.Data, nil
}

type GetWebhooksConfig struct {
	Workspace string  `json:"workspace"`
	Resource  *string `json:"resource"`
}

func (service *Service) GetWebhooks(config *GetWebhooksConfig) (*[]Webhook, *errortools.Error) {
	var values = url.Values{}

	values.Set("workspace", config.Workspace)
	if config.Resource != nil {
		values.Set("resource", *config.Resource)
	}

	var webhooks []Webhook

	for {
		var webhooks_ []Webhook

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("webhooks?%s", values.Encode())),
			ResponseModel: &webhooks_,
		}
		_, _, nextPage, e := service.getData(&requestConfig)
		if e != nil {
			return nil, e
		}

		webhooks = append(webhooks, webhooks_...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		values.Set("offset", nextPage.Offset)
	}

	return &webhooks, nil
}

func (service *Service) DeleteWebhook(webhookGid string) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		Url:    service.url(fmt.Sprintf("webhooks/%s", webhookGid)),
	}

	_, _, e := service.httpRequest(&requestConfig)

	return e
}
