package asana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
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
func (a *Asana) Get(url string, responseModel interface{}) (*http.Request, *http.Response, *NextPage, *errortools.Error) {
	return a.httpRequest(http.MethodGet, url, nil, responseModel)
}

func (a *Asana) httpRequest(httpMethod string, url string, bodyModel interface{}, responseModel interface{}) (*http.Request, *http.Response, *NextPage, *errortools.Error) {
	if utilities.IsNil(bodyModel) {
		return a.httpRequestWithBuffer(httpMethod, url, nil, responseModel)
	}

	b, err := json.Marshal(bodyModel)
	if err != nil {
		return nil, nil, nil, errortools.ErrorMessage(err)
	}

	return a.httpRequestWithBuffer(httpMethod, url, bytes.NewBuffer(b), responseModel)
}

func (a *Asana) httpRequestWithBuffer(httpMethod string, url string, body io.Reader, responseModel interface{}) (*http.Request, *http.Response, *NextPage, *errortools.Error) {
	client := &http.Client{}

	e := new(errortools.Error)

	request, err := http.NewRequest(httpMethod, url, body)
	e.SetRequest(request)
	if err != nil {
		e.SetMessage(err)
		return request, nil, nil, e
	}

	// Add authorization token to header
	bearer := fmt.Sprintf("Bearer %s", a.BearerToken)
	request.Header.Add("Authorization", bearer)
	request.Header.Set("Accept", "application/json")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	// Send out the HTTP request
	response, e := utilities.DoWithRetry(client, request, 10, 3)

	if response != nil {
		// Check HTTP StatusCode
		if response.StatusCode < 200 || response.StatusCode > 299 {
			fmt.Println(fmt.Sprintf("ERROR in %s", httpMethod))
			fmt.Println("url", url)
			fmt.Println("StatusCode", response.StatusCode)

			if e == nil {
				e = new(errortools.Error)
				e.SetRequest(request)
				e.SetResponse(response)
			}

			e.SetMessage(fmt.Sprintf("Server returned statuscode %v", response.StatusCode))
		}
	}

	if response.Body == nil {
		return request, response, nil, e
	}

	if e != nil {
		errorResponse := ErrorResponse{}
		err := a.unmarshalError(response, &errorResponse)
		errortools.CaptureInfo(err)

		b, _ := json.Marshal(errorResponse)
		e.SetExtra("error", string(b))

		return request, response, nil, e
	}

	res := Response{}

	if responseModel != nil {
		defer response.Body.Close()

		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			e.SetMessage(err)
			return request, response, nil, e
		}

		err = json.Unmarshal(b, &res)
		if err != nil {
			e.SetMessage(err)
			return request, response, nil, e
		}

		if *res.Data != nil {
			err = json.Unmarshal(*res.Data, &responseModel)
			if err != nil {
				e.SetMessage(err)
				return request, response, nil, e
			}
		}
	}

	return request, response, res.NextPage, nil
}

func (a *Asana) unmarshalError(response *http.Response, errorModel interface{}) *errortools.Error {
	if response == nil {
		return nil
	}
	if reflect.TypeOf(errorModel).Kind() != reflect.Ptr {
		return errortools.ErrorMessage("Type of errorModel must be a pointer.")
	}
	if reflect.ValueOf(errorModel).IsNil() {
		return nil
	}

	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	err = json.Unmarshal(b, &errorModel)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	return nil
}
