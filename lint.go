package pbc

import (
	"fmt"
)

/*
  A basic warning datatype with simple messages.

  The empty string is used as the NULL-space of this type, meaning there is no warning.
*/
type Warning string

const NilWarning = Warning("")

func NewWarning(format string, e ...interface{}) Warning {
	return Warning(fmt.Sprintf(format, e...))
}
