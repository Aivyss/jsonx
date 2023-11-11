package definitions

import (
	"errors"
	"github.com/aivyss/jsonx/constant"
	"github.com/aivyss/typex/util"
	"reflect"
	"regexp"
	"strings"
)

var (
	defaultAnnotations = map[string]Annotation{
		"NotBlank":         {name: "NotBlank", Validate: notBlank},
		"NotEmpty":         {name: "NotEmpty", Validate: notEmpty},
		"Required":         {name: "Required", Validate: required},
		"Email":            {name: "Email", Validate: email},
		"NotContainsNil":   {name: "NotContainsNil", Validate: notContainsNil},
		"NotContainsEmpty": {name: "NotContainsEmpty", Validate: notContainsEmpty},
		"NotContainsBlank": {name: "NotContainsBlank", Validate: notContainsBlank},
	}
	customAnnotations = map[string]Annotation{}
)

type Annotation struct {
	name     string
	Validate AnnotationValidate
}

func ConvertToAnnotation(v string) (*Annotation, error) {
	anno, defaultOk := defaultAnnotations[v]
	if defaultOk {
		return &anno, nil
	}

	anno, customOk := customAnnotations[v]
	if customOk {
		return &anno, nil
	}

	return nil, errors.New("invalid annotation")
}

func RegisterCustomAnnotation(annotationName string, validateFunc AnnotationValidate) error {
	annotation := Annotation{
		name:     annotationName,
		Validate: validateFunc,
	}

	_, ok := defaultAnnotations[annotation.name]
	if ok {
		return errors.New("duplicate annotation name with one of default annotation")
	}

	customAnnotations[annotation.name] = annotation

	return nil
}

// notEmpty
// @NotEmpty
func notEmpty(v any) error {
	var err error = nil

	switch v.(type) {
	case *string:
		if util.IsNil(v) {
			err = errors.New("@NotEmpty nil value")
		} else {
			if *(v.(*string)) == "" {
				err = errors.New("@NotEmpty empty value")
			}
		}
	case string:
		if v.(string) == "" {
			err = errors.New("@NotEmpty empty value")
		}
	default:
		err = errors.New("@NotEmpty wrong type")
	}

	return err
}

// notBlank
// @NotBlank
func notBlank(v any) error {
	var err error = nil

	switch v.(type) {
	case *string:
		if util.IsNil(v) {
			err = errors.New("@NotBlank nil value")
		} else {
			if strings.TrimSpace(*(v.(*string))) == "" {
				err = errors.New("@NotBlank blank value")
			}
		}
	case string:
		if strings.TrimSpace(v.(string)) == "" {
			err = errors.New("@NotBlank empty value")
		}
	default:
		err = errors.New("@NotBlank wrong type")
	}

	return err
}

// requried
// @Required
func required(v any) error {
	if util.IsNil(v) {
		return errors.New("@Required")
	}

	return nil
}

// email
// @Email
func email(v any) error {
	var err error = nil

	switch v.(type) {
	case *string:
		if util.IsNil(v) {
			err = errors.New("@Email nil value")
		} else {
			matched, err := regexp.MatchString(constant.EmailRegex, *(v.(*string)))
			if err != nil {
				err = errors.Join(err, errors.New("@Email not email format"))
			}

			if !matched {
				err = errors.New("@Email not email format")
			}
		}
	case string:
		matched, err := regexp.MatchString(constant.EmailRegex, v.(string))
		if err != nil {
			err = errors.Join(err, errors.New("@Email not email format"))
		}

		if !matched {
			err = errors.New("@Email not email format")
		}
	default:
		err = errors.New("@Email wrong type")
	}

	return err
}

// notContainsNil
// @NotContainsNil
func notContainsNil(v any) error {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Slice {
		return errors.New("@NotContainsNil not slice type")
	}

	for i := 0; i < valueOf.Len(); i++ {
		elem := valueOf.Index(i)
		if elem.IsNil() {
			return errors.New("@NotContainsNil nil value is not allowed")
		}
	}

	return nil
}

// notContainsEmpty
// @NotContainsEmpty
func notContainsEmpty(v any) error {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Slice {
		return errors.New("@NotContainsEmpty not slice type")
	}

	for i := 0; i < valueOf.Len(); i++ {
		elem := valueOf.Index(i)
		err := notEmpty(elem.Interface())
		if err != nil {
			return errors.Join(err, errors.New("@NotContainsEmpty"))
		}
	}

	return nil
}

// notContainsBlank
// @NotContainsBlank
func notContainsBlank(v any) error {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Slice {
		return errors.New("@NotContainsBlank not slice type")
	}

	for i := 0; i < valueOf.Len(); i++ {
		elem := valueOf.Index(i)
		err := notBlank(elem.Interface())
		if err != nil {
			return errors.Join(err, errors.New("@NotContainsBlank"))
		}
	}

	return nil
}
