package definitions

import (
	"errors"
	"github.com/aivyss/jsonx/constant"
	"github.com/aivyss/typex"
	"github.com/aivyss/typex/types"
	"reflect"
	"regexp"
	"strings"
	"time"
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
		"Positive":         {name: "Positive", Validate: positive},
		"Negative":         {name: "Negative", Validate: negative},
		"PositiveOrZero":   {name: "PositiveOrZero", Validate: positiveOrZero},
		"NegativeOrZero":   {name: "NegativeOrZero", Validate: negativeOrZero},
		"Future":           {name: "Future", Validate: future},
		"Present":          {name: "Present", Validate: present},
		"Past":             {name: "Past", Validate: past},
		"FutureOrPresent":  {name: "FutureOrPresent", Validate: futureOrPresent},
		"PastOrPresent":    {name: "PastOrPresent", Validate: pastOrPresent},
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
		if types.IsNil(v) {
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
		if types.IsNil(v) {
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
	if types.IsNil(v) {
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
		if types.IsNil(v) {
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

// positive
// @Positive
func positive(v any) error {
	nilErr := errors.New("@Positive nil value")
	notPositiveErr := errors.New("@Positive not positive value")
	var err error = nil

	typex.Opt(v).IfPresent(func(v any) {
		switch v.(type) {
		case int8:
			if v.(int8) <= 0 {
				err = notPositiveErr
			}
		case int16:
			if v.(int16) <= 0 {
				err = notPositiveErr
			}
		case int32:
			if v.(int32) <= 0 {
				err = notPositiveErr
			}
		case int64:
			if v.(int64) <= 0 {
				err = notPositiveErr
			}
		case int:
			if v.(int) <= 0 {
				err = notPositiveErr
			}
		case float32:
			if v.(float32) <= 0 {
				err = notPositiveErr
			}
		case float64:
			if v.(float64) <= 0 {
				err = notPositiveErr
			}

		case *int8:
			i := v.(*int8)
			if *i <= 0 {
				err = notPositiveErr
			}
		case *int16:
			i := v.(*int16)
			if *i <= 0 {
				err = notPositiveErr
			}
		case *int32:
			i := v.(*int32)
			if *i <= 0 {
				err = notPositiveErr
			}
		case *int64:
			i := v.(*int64)
			if *i <= 0 {
				err = notPositiveErr
			}
		case *int:
			i := v.(*int)
			if *i <= 0 {
				err = notPositiveErr
			}
		case *float32:
			i := v.(*float32)
			if *i <= 0 {
				err = notPositiveErr
			}
		case *float64:
			i := v.(*float64)
			if *i <= 0 {
				err = notPositiveErr
			}
		default:
			err = errors.New("@Positive not number type")
		}
	}).ElseDo(func() {
		err = nilErr
	})

	return err
}

// positiveOrZero
// @PositiveOrZero
func positiveOrZero(v any) error {
	nilErr := errors.New("@PositiveOrZero nil value")
	notPositiveErr := errors.New("@PositiveOrZero negative value")
	var err error = nil

	typex.Opt(v).IfPresent(func(v any) {
		switch v.(type) {
		case int8:
			if v.(int8) < 0 {
				err = notPositiveErr
			}
		case int16:
			if v.(int16) < 0 {
				err = notPositiveErr
			}
		case int32:
			if v.(int32) < 0 {
				err = notPositiveErr
			}
		case int64:
			if v.(int64) < 0 {
				err = notPositiveErr
			}
		case int:
			if v.(int) < 0 {
				err = notPositiveErr
			}
		case float32:
			if v.(float32) < 0 {
				err = notPositiveErr
			}
		case float64:
			if v.(float64) < 0 {
				err = notPositiveErr
			}

		case *int8:
			i := v.(*int8)
			if *i < 0 {
				err = notPositiveErr
			}
		case *int16:
			i := v.(*int16)
			if *i < 0 {
				err = notPositiveErr
			}
		case *int32:
			i := v.(*int32)
			if *i < 0 {
				err = notPositiveErr
			}
		case *int64:
			i := v.(*int64)
			if *i < 0 {
				err = notPositiveErr
			}
		case *int:
			i := v.(*int)
			if *i < 0 {
				err = notPositiveErr
			}
		case *float32:
			i := v.(*float32)
			if *i < 0 {
				err = notPositiveErr
			}
		case *float64:
			i := v.(*float64)
			if *i < 0 {
				err = notPositiveErr
			}
		default:
			err = errors.New("@PositiveOrZero not number type")
		}
	}).ElseDo(func() {
		err = nilErr
	})

	return err
}

// negative
// @Negative
func negative(v any) error {
	nilErr := errors.New("@Negative nil value")
	notNegativeErr := errors.New("@Negative not negative value")
	var err error = nil

	typex.Opt(v).IfPresent(func(v any) {
		switch v.(type) {
		case int8:
			if v.(int8) >= 0 {
				err = notNegativeErr
			}
		case int16:
			if v.(int16) >= 0 {
				err = notNegativeErr
			}
		case int32:
			if v.(int32) >= 0 {
				err = notNegativeErr
			}
		case int64:
			if v.(int64) >= 0 {
				err = notNegativeErr
			}
		case int:
			if v.(int) >= 0 {
				err = notNegativeErr
			}
		case float32:
			if v.(float32) >= 0 {
				err = notNegativeErr
			}
		case float64:
			if v.(float64) >= 0 {
				err = notNegativeErr
			}

		case *int8:
			i := v.(*int8)
			if *i >= 0 {
				err = notNegativeErr
			}
		case *int16:
			i := v.(*int16)
			if *i >= 0 {
				err = notNegativeErr
			}
		case *int32:
			i := v.(*int32)
			if *i >= 0 {
				err = notNegativeErr
			}
		case *int64:
			i := v.(*int64)
			if *i >= 0 {
				err = notNegativeErr
			}
		case *int:
			i := v.(*int)
			if *i >= 0 {
				err = notNegativeErr
			}
		case *float32:
			i := v.(*float32)
			if *i >= 0 {
				err = notNegativeErr
			}
		case *float64:
			i := v.(*float64)
			if *i >= 0 {
				err = notNegativeErr
			}
		default:
			err = errors.New("@Negative not number type")
		}
	}).ElseDo(func() {
		err = nilErr
	})

	return err
}

// negativeOrZero
// @NegativeOrZero
func negativeOrZero(v any) error {
	nilErr := errors.New("@NegativeOrZero nil value")
	notNegativeErr := errors.New("@NegativeOrZero positive value")
	var err error = nil

	typex.Opt(v).IfPresent(func(v any) {
		switch v.(type) {
		case int8:
			if v.(int8) > 0 {
				err = notNegativeErr
			}
		case int16:
			if v.(int16) > 0 {
				err = notNegativeErr
			}
		case int32:
			if v.(int32) > 0 {
				err = notNegativeErr
			}
		case int64:
			if v.(int64) > 0 {
				err = notNegativeErr
			}
		case int:
			if v.(int) > 0 {
				err = notNegativeErr
			}
		case float32:
			if v.(float32) > 0 {
				err = notNegativeErr
			}
		case float64:
			if v.(float64) > 0 {
				err = notNegativeErr
			}

		case *int8:
			i := v.(*int8)
			if *i > 0 {
				err = notNegativeErr
			}
		case *int16:
			i := v.(*int16)
			if *i > 0 {
				err = notNegativeErr
			}
		case *int32:
			i := v.(*int32)
			if *i > 0 {
				err = notNegativeErr
			}
		case *int64:
			i := v.(*int64)
			if *i > 0 {
				err = notNegativeErr
			}
		case *int:
			i := v.(*int)
			if *i > 0 {
				err = notNegativeErr
			}
		case *float32:
			i := v.(*float32)
			if *i > 0 {
				err = notNegativeErr
			}
		case *float64:
			i := v.(*float64)
			if *i > 0 {
				err = notNegativeErr
			}
		default:
			err = errors.New("@NegativeOrZero not number type")
		}
	}).ElseDo(func() {
		err = nilErr
	})

	return err
}

// future
// @Future
func future(v any) error {
	notFutureErr := errors.New("@Future not future time")
	var err error = nil
	typex.Opt(v).IfPresent(func(v any) {
		now := time.Now()

		switch v.(type) {
		case time.Time:
			t := v.(time.Time)
			if t.Before(now) || equal(t, now) {
				err = notFutureErr
			}
		case *time.Time:
			t := v.(*time.Time)
			if t.Before(now) || equal(*t, now) {
				err = notFutureErr
			}
		default:
			err = errors.New("@Future wrong type")
		}
	}).ElseDo(func() {
		err = errors.New("@Future nil value")
	})

	return err
}

// futureOrPresent
// @FutureOrPresent
func futureOrPresent(v any) error {
	notFutureOrPresentErr := errors.New("@FutureOrPresent past time")
	var err error = nil

	typex.Opt(v).IfPresent(func(v any) {
		now := time.Now()

		switch v.(type) {
		case time.Time:
			t := v.(time.Time)
			if !t.After(now) && !equal(t, now) {
				err = notFutureOrPresentErr
			}
		case *time.Time:
			t := v.(*time.Time)
			if !t.After(now) && !equal(*t, now) {
				err = notFutureOrPresentErr
			}
		default:
			err = errors.New("@FutureOrPresent wrong type")
		}
	}).ElseDo(func() {
		err = errors.New("@FutureOrPresent nil value")
	})

	return err
}

// pastOrPresent
// @PastOrPresent
func pastOrPresent(v any) error {
	notPastOrPresentErr := errors.New("@PastOrPresent future time")
	var err error = nil

	typex.Opt(v).IfPresent(func(v any) {
		now := time.Now()

		switch v.(type) {
		case time.Time:
			t := v.(time.Time)
			if !t.Before(now) && !equal(t, now) {
				err = notPastOrPresentErr
			}
		case *time.Time:
			t := v.(*time.Time)
			if !t.Before(now) && !equal(*t, now) {
				err = notPastOrPresentErr
			}
		default:
			err = errors.New("@PastOrPresent wrong type")
		}
	}).ElseDo(func() {
		err = errors.New("@PastOrPresent nil value")
	})

	return err
}

// present
// @Present
func present(v any) error {
	notPresentErr := errors.New("@Present not present time")
	var err error = nil
	typex.Opt(v).IfPresent(func(v any) {
		now := time.Now()

		switch v.(type) {
		case time.Time:
			t := v.(time.Time)
			if !equal(t, now) {
				err = notPresentErr
			}
		case *time.Time:
			t := v.(*time.Time)
			if !equal(*t, now) {
				err = notPresentErr
			}
		default:
			err = errors.New("@Present wrong type")
		}
	}).ElseDo(func() {
		err = errors.New("@Present nil value")
	})

	return err
}

// past
// @Past
func past(v any) error {
	notPastErr := errors.New("@Past not past time")
	var err error = nil
	typex.Opt(v).IfPresent(func(v any) {
		now := time.Now()

		switch v.(type) {
		case time.Time:
			t := v.(time.Time)
			if equal(t, now) || t.After(now) {
				err = notPastErr
			}
		case *time.Time:
			t := v.(*time.Time)
			if equal(*t, now) || t.After(now) {
				err = notPastErr
			}
		default:
			err = errors.New("@Past wrong type")
		}
	}).ElseDo(func() {
		err = errors.New("@Past nil value")
	})

	return err
}

func equal(t1 time.Time, t2 time.Time) bool {
	return t1.Unix()-t2.Unix() == 0
}
