package customerrors

import "fmt"

type MultipartError struct {
	OpenErr bool
	Err     error
	Source  string
}

func (m *MultipartError) Error() string {
	return fmt.Sprintf("File Open: %v | Source: %v | %v", m.OpenErr, m.Source, m.Err.Error())
}

type HTTPError struct {
	StatusCode uint16
	Err        error
	Source     string
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("Status: %v | Source: %v | %v", h.StatusCode, h.Source, h.Err.Error())
}
