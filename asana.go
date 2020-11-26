package asana

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// type
//
type Asana struct {
	ApiURL      string
	BearerToken string
	IsLive      bool
}

// Response represents highest level of exactonline api response
//
type Response struct {
	Data     *json.RawMessage `json:"data,omitempty"`
	NextPage *NextPage        `json:"next_page,omitempty"`
	Errors   *[]AsanaError    `json:"errors,omitempty"`
}

// NextPage contains info for batched data retrieval
//
type NextPage struct {
	Offset string `json:"offset"`
	Path   string `json:"path"`
	URI    string `json:"uri"`
}

// AsanaError contains error info
//
type AsanaError struct {
	Message string `json:"message"`
	Help    string `json:"help"`
}

func New(apiURL string, bearerToken string, isLive bool) (*Asana, *errortools.Error) {
	i := new(Asana)

	if apiURL == "" {
		return nil, errortools.ErrorMessage("Asana ApiUrl not provided")
	}
	if bearerToken == "" {
		return nil, errortools.ErrorMessage("Asana Token not provided")
	}

	i.ApiURL = apiURL
	i.BearerToken = bearerToken
	i.IsLive = isLive

	if !strings.HasSuffix(i.ApiURL, "/") {
		i.ApiURL = i.ApiURL + "/"
	}

	return i, nil
}

// generic Get method
//
func (i *Asana) Get(url string, model interface{}) (*NextPage, *Response, *errortools.Error) {
	client := &http.Client{}

	e := new(errortools.Error)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	e.SetRequest(req)
	if err != nil {
		e.SetMessage(err)
		return nil, nil, e
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", "Bearer "+i.BearerToken)

	// Send out the HTTP request
	res, err := utilities.DoWithRetry(client, req, 10, 5)
	e.SetResponse(res)
	if err != nil {
		e.SetMessage(err)
		return nil, nil, e
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	response := Response{}

	err = json.Unmarshal(b, &response)
	if err != nil {
		e.SetMessage(err)
		return nil, nil, e
	}

	if response.Data == nil {
		return nil, &response, nil
	}

	err = json.Unmarshal(*response.Data, &model)
	if err != nil {
		e.SetMessage(err)
		return nil, nil, e
	}

	return response.NextPage, &response, nil
}

func (a *Asana) captureErrors(e *errortools.Error, response *Response) {
	if response != nil {
		if response.Errors != nil {
			ee := []string{}
			for _, err := range *response.Errors {
				ee = append(ee, fmt.Sprintf("%s\n%s", err.Message, err.Help))
			}

			e.SetMessage(strings.Join(ee, "\n\n"))
			errortools.CaptureMessage(e, a.IsLive)
		}
	}
}
