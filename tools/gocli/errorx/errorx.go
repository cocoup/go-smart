package errorx

import (
	"fmt"
	"strings"
)

var errorFormat = `go-smart error: %+v
%s`

// GocliError represents a go-smart error.
type GocliError struct {
	message []string
	err     error
}

func (e *GocliError) Error() string {
	detail := wrapMessage(e.message...)
	return fmt.Sprintf(errorFormat, e.err, detail)
}

// Wrap wraps an error with go-smart version and message.
func Wrap(err error, message ...string) error {
	e, ok := err.(*GocliError)
	if ok {
		return e
	}

	return &GocliError{
		message: message,
		err:     err,
	}
}

func wrapMessage(message ...string) string {
	if len(message) == 0 {
		return ""
	}
	return fmt.Sprintf(`message: %s`, strings.Join(message, "\n"))
}
