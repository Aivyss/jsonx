package jsonx

import (
	"encoding/json"
	"errors"
	"github.com/aivyss/jsonx/tag"
	"github.com/aivyss/jsonx/validate"
	"reflect"
	"sort"
)

func RegisterValidator[T any](v validate.Validator[T]) {
	typeOf := reflect.TypeOf(new(T)).Elem()
	validatorMap[typeOf] = v
}

func RegisterOrderedValidator[T any](v validate.OrderedValidator[T]) {
	typeOf := reflect.TypeOf(new(T)).Elem()
	orderedValidatorMap.Put(typeOf, v)
}

func Unmarshal[V any](data []byte) (*V, error) {
	v := new(V)
	err := json.Unmarshal(data, v)
	if err != nil {
		return nil, errors.Join(errors.New("fail to unmarshal"), err)
	}

	// tag validation
	typeOf := reflect.TypeOf(*v)
	for i := 0; i < typeOf.NumField(); i++ {
		err := tag.ValidateAnnotationTag(typeOf.Field(i).Tag.Get("annotation"), reflect.ValueOf(*v).Field(i).Interface())
		if err != nil {
			return nil, err
		}
	}

	typeOfElem := reflect.TypeOf(v).Elem()
	// default validation
	if validator, ok := validatorMap[typeOfElem]; ok {
		if vd, ok := validator.(validate.Validator[V]); ok {
			if err := vd.Validate(*v); err != nil {
				return nil, err
			}
		}
	}

	// ordered validations
	validators := orderedValidatorMap.Get(typeOfElem)
	vSlice := make([]validate.OrderedValidator[V], 0, len(validators))
	for _, v := range validators {
		if v2, ok := v.(validate.OrderedValidator[V]); ok {
			vSlice = append(vSlice, v2)
		}
	}
	sort.Slice(vSlice, func(i, j int) bool {
		return vSlice[i].Order() < vSlice[j].Order()
	})
	for _, validator := range vSlice {
		if err := validator.Validate(*v); err != nil {
			return nil, err
		}
	}

	return v, nil
}

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func Close() {
	validatorMap = map[reflect.Type]any{}
	orderedValidatorMap.Clean()
}
