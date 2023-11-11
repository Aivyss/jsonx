package tag

import (
	"github.com/aivyss/jsonx/definitions"
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
