package asana

// ErrorResponse stores general Asana error response
//
type ErrorResponse struct {
	Errors []struct {
		Message string `json:"message"`
		Help    string `json:"help"`
	} `json:"errors"`
}
