package jsonx

import (
	"github.com/aivyss/jsonx/common"
	"reflect"
)

var orderedValidatorMap = common.NewMultiMap[reflect.Type, any]()
var validatorMap = map[reflect.Type]any{}
