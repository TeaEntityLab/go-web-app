package errtrace

import (
	"fmt"
)

type Error struct {
	ErrRef    error  `json:"raw"`
	ErrCode   int    `json:"code"`
	ErrTitle  string `json:"title"`
	ErrDetail string `json:"detail"`

	StatusCode int `json:"status"`
}

func (e Error) Error() string {
	if e.ErrRef != nil {
		return e.ErrRef.Error()
	}

	return fmt.Sprintf("error: %d status: %d  (%s):%s\nraw: %v", e.ErrCode, e.StatusCode, e.ErrTitle, e.ErrDetail, e.ErrRef)
}
