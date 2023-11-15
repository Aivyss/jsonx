package jsonx

import (
	"github.com/aivyss/jsonx/tag"
	"reflect"
	"time"
)

func tagValidation[V any](v V) error {
	typeOf := reflect.TypeOf(v)

	for i := 0; i < typeOf.NumField(); i++ {
		kind := typeOf.Field(i).Type.Kind()

		switch kind {
		case reflect.Pointer:
			if typeOf.Field(i).Type.Elem().Kind() == reflect.Struct {
				if typeOf.Field(i).Type.Elem() != reflect.TypeOf(time.Time{}) {
					if err := tagValidation(reflect.ValueOf(v).Field(i).Interface()); err != nil {
						return err
					}
				}

				return nil
			}
		case reflect.Struct:
			if typeOf.Field(i).Type != reflect.TypeOf(time.Time{}) {
				if err := tagValidation(reflect.ValueOf(v).Field(i).Interface()); err != nil {
					return err
				}

				return nil
			}
		}

		// annotation validation
		fieldTag := typeOf.Field(i).Tag

		if annotations := fieldTag.Get("annotation"); annotations != "" {
			if err := tag.ValidateAnnotationTag(
				annotations,
				reflect.ValueOf(v).Field(i).Interface(),
			); err != nil {
				return exchangeIfFieldError(fieldTag, err)
			}
		}

		// regex validation
		if pattern := fieldTag.Get("pattern"); pattern != "" {
			if err := tag.RegexTag(
				pattern,
				reflect.ValueOf(v).Field(i).Interface(),
			); err != nil {
				return exchangeIfFieldError(fieldTag, err)
			}
		}
	}

	return nil
}
