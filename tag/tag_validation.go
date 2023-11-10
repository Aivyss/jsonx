package tag

import (
	"errors"
	"github.com/aivyss/typex/util"
	"strings"
)

func ValidateAnnotationTag(tagValue string, value any) error {
	annotations := strings.Split(tagValue, "@")

	var err error = nil
	for _, annotation := range annotations {
		switch annotation {
		case "NotBlank":
			err = NotBlank(value)
		case "NotEmpty":
			err = NotEmptyString(value)
		case "Required":
			err = Required(value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// NotEmptyString
// @NotEmpty
func NotEmptyString(v any) error {
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

// NotBlank
// @NotBlank
func NotBlank(v any) error {
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

// Required
// @Required
func Required(v any) error {
	if util.IsNil(v) {
		return errors.New("@Required")
	}

	return nil
}
