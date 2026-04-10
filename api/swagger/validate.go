package swagger

import (
	"strconv"
	"strings"

	apiSpec "github.com/suyuan32/goctls/api/spec"
)

type validateRules struct {
	required bool
	min      *float64
	max      *float64
	length   *float64
}

func validateRulesFromTags(tags *apiSpec.Tags) validateRules {
	validateTag, err := tags.Get(tagValidate)
	if err != nil || validateTag == nil {
		return validateRules{}
	}

	items := make([]string, 0, len(validateTag.Options)+1)
	if len(validateTag.Name) > 0 {
		items = append(items, validateTag.Name)
	}
	items = append(items, validateTag.Options...)

	var rules validateRules
	for _, raw := range items {
		item := strings.TrimSpace(raw)
		if len(item) == 0 {
			continue
		}

		switch {
		case item == "required":
			rules.required = true
		case strings.HasPrefix(item, "min="):
			if value, ok := parseFloatValidateValue(item, "min="); ok {
				rules.min = &value
			}
		case strings.HasPrefix(item, "max="):
			if value, ok := parseFloatValidateValue(item, "max="); ok {
				rules.max = &value
			}
		case strings.HasPrefix(item, "len="):
			if value, ok := parseFloatValidateValue(item, "len="); ok {
				rules.length = &value
			}
		}
	}

	return rules
}

func applyValidateRulesToBounds(tp apiSpec.Type, rules validateRules, minimum, maximum **float64,
	minLength, maxLength, minItems, maxItems **int64) {
	if rules.min == nil && rules.max == nil && rules.length == nil {
		return
	}

	switch sampleTypeFromGoType(Context{}, tp) {
	case swaggerTypeInteger, swaggerTypeNumber:
		if *minimum == nil && rules.min != nil {
			*minimum = rules.min
		}
		if *maximum == nil && rules.max != nil {
			*maximum = rules.max
		}
		if rules.length != nil {
			if *minimum == nil {
				*minimum = rules.length
			}
			if *maximum == nil {
				*maximum = rules.length
			}
		}
	case swaggerTypeString:
		if *minLength == nil {
			if rules.length != nil {
				if value, ok := toInt64IfWhole(*rules.length); ok {
					*minLength = &value
				}
			} else if rules.min != nil {
				if value, ok := toInt64IfWhole(*rules.min); ok {
					*minLength = &value
				}
			}
		}
		if *maxLength == nil {
			if rules.length != nil {
				if value, ok := toInt64IfWhole(*rules.length); ok {
					*maxLength = &value
				}
			} else if rules.max != nil {
				if value, ok := toInt64IfWhole(*rules.max); ok {
					*maxLength = &value
				}
			}
		}
	case swaggerTypeArray:
		if *minItems == nil {
			if rules.length != nil {
				if value, ok := toInt64IfWhole(*rules.length); ok {
					*minItems = &value
				}
			} else if rules.min != nil {
				if value, ok := toInt64IfWhole(*rules.min); ok {
					*minItems = &value
				}
			}
		}
		if *maxItems == nil {
			if rules.length != nil {
				if value, ok := toInt64IfWhole(*rules.length); ok {
					*maxItems = &value
				}
			} else if rules.max != nil {
				if value, ok := toInt64IfWhole(*rules.max); ok {
					*maxItems = &value
				}
			}
		}
	}
}

func hasValidateRequired(tags *apiSpec.Tags) bool {
	return validateRulesFromTags(tags).required
}

func parseFloatValidateValue(item, prefix string) (float64, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(item, prefix))
	if len(value) == 0 {
		return 0, false
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func toInt64IfWhole(value float64) (int64, bool) {
	if value != float64(int64(value)) {
		return 0, false
	}
	return int64(value), true
}
