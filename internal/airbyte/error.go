package airbyte

import (
	"fmt"
	"os"
)

type airbyteError struct {
	Msg   string
	Field string
	Error error
}

func NewAirbyteError(msg string, field string, err error) *airbyteError {
	return &airbyteError{
		Msg:   msg,
		Field: field,
		Error: err,
	}
}

func (e *airbyteError) Print() {
	fmt.Printf("\n ERROR: %s", e.Msg)
	fmt.Printf("\n FIELD: %s", e.Field)
	fmt.Printf("\n DEBUG: %v", e)
	os.Exit(1)
}
