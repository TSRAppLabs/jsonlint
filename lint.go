package jsonlint

import (
	"fmt"
)

/*
  A basic warning datatype with simple messages.

  The empty string is used as the NULL-space of this type, meaning there is no warning.
*/
type Warning []string

/*
  Creates a new Warning with printf formatting
*/
func NewWarning(format string, e ...interface{}) Warning {
	return Warning([]string{fmt.Sprintf(format, e...)})
}
