package swagger

import (
	"github.com/go-openapi/spec"
	apiSpec "github.com/suyuan32/goctls/api/spec"
)

func propertiesFromType(ctx Context, tp apiSpec.Type) (spec.SchemaProperties, []string) {
	var (
		properties     = map[string]spec.Schema{}
		requiredFields []string
	)
	switch val := tp.(type) {
	case apiSpec.PointerType:
		return propertiesFromType(ctx, val.Type)
	case apiSpec.ArrayType:
		return propertiesFromType(ctx, val.Value)
	case apiSpec.DefineStruct, apiSpec.NestedStruct:
		rangeMemberAndDo(ctx, val, func(tag *apiSpec.Tags, required bool, member apiSpec.Member) {
			validateRules := validateRulesFromTags(tag)
			if validateRules.required {
				required = true
			}

			var (
				jsonTagString                      = member.Name
				minimum, maximum                   *float64
				exclusiveMinimum, exclusiveMaximum bool
				minLength, maxLength               *int64
				minItems, maxItems                 *int64
				example, defaultValue              any
				enum                               []any
			)
			pathTag, _ := tag.Get(tagPath)
			if pathTag != nil {
				return
			}
			formTag, _ := tag.Get(tagForm)
			if formTag != nil {
				return
			}
			headerTag, _ := tag.Get(tagHeader)
			if headerTag != nil {
				return
			}

			jsonTag, _ := tag.Get(tagJson)
			if jsonTag != nil {
				jsonTagString = jsonTag.Name
				minimum, maximum, exclusiveMinimum, exclusiveMaximum = rangeValueFromOptions(jsonTag.Options)
				example = exampleValueFromOptions(ctx, jsonTag.Options, member.Type)
				defaultValue = defValueFromOptions(ctx, jsonTag.Options, member.Type)
				enum = enumsValueFromOptions(jsonTag.Options)
			}
			applyValidateRulesToBounds(member.Type, validateRules, &minimum, &maximum, &minLength, &maxLength, &minItems, &maxItems)

			if required {
				requiredFields = append(requiredFields, jsonTagString)
			}

			schema := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					Example: example,
				},
				SchemaProps: spec.SchemaProps{
					Description:          formatComment(member.Comment),
					Type:                 typeFromGoType(ctx, member.Type),
					Default:              defaultValue,
					Maximum:              maximum,
					ExclusiveMaximum:     exclusiveMaximum,
					Minimum:              minimum,
					ExclusiveMinimum:     exclusiveMinimum,
					MaxLength:            maxLength,
					MinLength:            minLength,
					MaxItems:             maxItems,
					MinItems:             minItems,
					Enum:                 enum,
					AdditionalProperties: mapFromGoType(ctx, member.Type),
				},
			}

			switch sampleTypeFromGoType(ctx, member.Type) {
			case swaggerTypeArray:
				schema.Items = itemsFromGoType(ctx, member.Type)
				// Special handling for arrays with useDefinitions
				if ctx.UseDefinitions {
					// For arrays, check if the array element (not the array itself) contains a struct
					if arrayType, ok := member.Type.(apiSpec.ArrayType); ok {
						if structName, containsStruct := containsStruct(arrayType.Value); containsStruct {
							// Set the $ref inside the items, not at the schema level
							schema.Items = &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: spec.MustCreateRef(getRefName(structName)),
									},
								},
							}
						}
					}
				}
			case swaggerTypeObject:
				p, r := propertiesFromType(ctx, member.Type)
				schema.Properties = p
				schema.Required = r
				// For objects with useDefinitions, set $ref at schema level
				if ctx.UseDefinitions {
					structName, containsStruct := containsStruct(member.Type)
					if containsStruct {
						schema.SchemaProps.Ref = spec.MustCreateRef(getRefName(structName))
					}
				}
			default:
				// For non-array, non-object types, apply useDefinitions logic
				if ctx.UseDefinitions {
					structName, containsStruct := containsStruct(member.Type)
					if containsStruct {
						schema.SchemaProps.Ref = spec.MustCreateRef(getRefName(structName))
					}
				}
			}

			properties[jsonTagString] = schema
		})
	}

	return properties, requiredFields
}

func containsStruct(tp apiSpec.Type) (string, bool) {
	switch val := tp.(type) {
	case apiSpec.PointerType:
		return containsStruct(val.Type)
	case apiSpec.ArrayType:
		return containsStruct(val.Value)
	case apiSpec.DefineStruct:
		return val.RawName, true
	case apiSpec.MapType:
		return containsStruct(val.Value)
	default:
		return "", false
	}
}

func getRefName(typeName string) string {
	return "#/definitions/" + typeName
}
