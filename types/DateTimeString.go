package types

import (
	"encoding/json"
	"fmt"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

const (
	layoutDateTime string = "2006-01-02T15:04:05.000Z"
)

type DateTimeString time.Time

func (d *DateTimeString) UnmarshalJSON(b []byte) error {
	var returnError = func() error {
		errortools.CaptureError(fmt.Sprintf("Cannot parse '%s' to DateTimeString", string(b)))
		return nil
	}

	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return returnError()
	}

	if s == "" || s == "0000-00-00 00:00:00.000" {
		d = nil
		return nil
	}

	_t, err := time.Parse(layoutDateTime, s)
	if err != nil {
		return returnError()
	}

	*d = DateTimeString(_t)
	return nil
}

func (d *DateTimeString) MarshalJSON() ([]byte, error) {
	if d == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(time.Time(*d).Format(layoutDateTime))
}

func (d *DateTimeString) ValuePtr() *time.Time {
	if d == nil {
		return nil
	}

	_d := time.Time(*d)
	return &_d
}

func (d DateTimeString) Value() time.Time {
	return time.Time(d)
}
