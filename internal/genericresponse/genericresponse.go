package genericresponse

import "fmt"

type GenericResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message,omitempty"`
}

func (e *GenericResponse) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}
