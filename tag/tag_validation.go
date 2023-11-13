package tag

import (
	"errors"
	"github.com/aivyss/jsonx/definitions"
	"regexp"
	"strings"
)

func ValidateAnnotationTag(tagValue string, value any) error {
	annotations := strings.Split(strings.TrimSpace(tagValue), "@")[1:]

	for _, annoStr := range annotations {
		annotation, err := definitions.ConvertToAnnotation(annoStr)
		if err != nil {
			return err
		}

		if err := annotation.Validate(value); err != nil {
			return err
		}
	}

	return nil
}

func RegexTag(pattern string, value any) error {
	s := ""
	switch value.(type) {
	case string:
		s = value.(string)
	case *string:
		s = *value.(*string)
	default:
		return errors.New("wrong field type")
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return errors.New("wrong regular expression")
	}

	if matched := regex.Match([]byte(s)); !matched {
		return errors.New("not matched (pattern)")
	}

	return nil
}
