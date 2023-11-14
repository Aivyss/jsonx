package errors

import (
	"encoding/json"
)

type FieldError struct {
	defaultMsg string
	name       string
}

func NewFieldErr(errorName, defaultMsg string) *FieldError {
	return &FieldError{
		name:       errorName,
		defaultMsg: defaultMsg,
	}
}

func (e *FieldError) Name() string {
	return e.name
}

func (e *FieldError) DefaultMsg() string {
	return e.defaultMsg
}

func (e *FieldError) Error() string {
	j, _ := json.Marshal(errorJsonStruct{
		FrameworkName: "jsonx",
		Name:          e.name,
		Msg:           e.defaultMsg,
	})

	return string(j)
}
