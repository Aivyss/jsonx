package jsonx

import (
	"encoding/json"
	"errors"
	"github.com/aivyss/jsonx/definitions"
	"github.com/aivyss/jsonx/tag"
	"github.com/aivyss/jsonx/validate"
	"github.com/aivyss/typex"
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

	if validationErr := Validate(*v); validationErr != nil {
		return nil, validationErr
	}

	return v, nil
}

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func RegisterCustomAnnotation(annotationName string, validateFunc definitions.AnnotationValidate) error {
	return definitions.RegisterCustomAnnotation(annotationName, validateFunc)
}

// Validate
// don't input pointer type
func Validate[T any](v T) error {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Struct {
		return errors.New("use only struct and its pointer")
	}

	// tag validation
	typeOf := reflect.TypeOf(v)
	for i := 0; i < typeOf.NumField(); i++ {
		// annotation validation
		fieldTag := typeOf.Field(i).Tag

		if annotations := fieldTag.Get("annotation"); annotations != "" {
			if err := tag.ValidateAnnotationTag(
				annotations,
				reflect.ValueOf(v).Field(i).Interface(),
			); err != nil {
				return err
			}
		}

		// regex validation
		if pattern := fieldTag.Get("pattern"); pattern != "" {
			if err := tag.RegexTag(
				pattern,
				reflect.ValueOf(v).Field(i).Interface(),
			); err != nil {
				return err
			}
		}
	}

	typeOfElem := reflect.TypeOf(&v).Elem()
	// default validation
	if validator, ok := validatorMap[typeOfElem]; ok {
		if vd, ok := validator.(validate.Validator[T]); ok {
			if err := vd.Validate(v); err != nil {
				return err
			}
		}
	}

	// ordered validations
	validators := orderedValidatorMap.Get(typeOfElem)
	vSlice := make([]validate.OrderedValidator[T], 0, len(validators))
	for _, v := range validators {
		if v2, ok := v.(validate.OrderedValidator[T]); ok {
			vSlice = append(vSlice, v2)
		}
	}
	sort.Slice(vSlice, func(i, j int) bool {
		return vSlice[i].Order() < vSlice[j].Order()
	})
	for _, validator := range vSlice {
		if err := validator.Validate(v); err != nil {
			return err
		}
	}

	return nil
}

func Close() {
	validatorMap = map[reflect.Type]any{}
	orderedValidatorMap = typex.NewMultiMap[reflect.Type, any]()
}
