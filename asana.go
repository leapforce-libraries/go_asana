package asana

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	types "github.com/Leapforce-nl/go_types"
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

func New(apiURL string, bearerToken string, isLive bool) (*Asana, error) {
	i := new(Asana)

	if apiURL == "" {
		return nil, &types.ErrorString{"Asana ApiUrl not provided"}
	}
	if bearerToken == "" {
		return nil, &types.ErrorString{"Asana Token not provided"}
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
func (i *Asana) Get(url string, model interface{}) (*NextPage, *Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", "Bearer "+i.BearerToken)

	// Send out the HTTP request
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	response := Response{}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, &response, err
	}

	if response.Data == nil {
		return nil, &response, nil
	}

	err = json.Unmarshal(*response.Data, &model)
	if err != nil {
		return nil, &response, err
	}

	return response.NextPage, &response, nil
}
