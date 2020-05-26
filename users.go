package asana

import (
	"fmt"

	sentry "github.com/getsentry/sentry-go"
)

// User stores User from Asana
//
type User struct {
	ID           string          `json:"gid"`
	Name         string          `json:"name"`
	ResourceType string          `json:"resource_type"`
	Email        string          `json:"email"`
	Photo        Photo           `json:"photo"`
	Workspaces   []CompactObject `json:"workspaces"`
}

// Photo stores Photo from Asana
//
type Photo struct {
	Image128x128 string `json:"image_128x128"`
	Image21x21   string `json:"image_21x21"`
	Image27x27   string `json:"image_27x27"`
	Image36x36   string `json:"image_36x36"`
	Image60x60   string `json:"image_60x60"`
}

// GetUsers returns all users
//
func (i *Asana) GetUsers() ([]User, error) {
	urlStr := "%susers?opt_fields=%s"

	users := []User{}

	url := fmt.Sprintf(urlStr, i.ApiURL, GetJSONTaggedFieldNames(User{}))
	//fmt.Println(url)

	_, response, err := i.Get(url, &users)
	if err != nil {
		return nil, err
	}

	if response != nil {
		if response.Errors != nil {
			for _, e := range *response.Errors {
				message := fmt.Sprintf("Error in %v: %v", url, e.Message)
				if i.IsLive {
					sentry.CaptureMessage(message)
				}
				fmt.Println(message)
			}
		}
	}

	return users, nil
}
