package jsonx

import (
	"github.com/aivyss/jsonx/errors"
	"github.com/aivyss/typex"
	"reflect"
)

var orderedValidatorMap = typex.NewMultiMap[reflect.Type, any]()
var validatorMap = map[reflect.Type]any{}
var fieldErrMap = map[string]errors.FieldError{}
