package swagger

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	apiSpec "github.com/suyuan32/goctls/api/spec"
)

func TestValidateRulesFromTags(t *testing.T) {
	tags, err := apiSpec.Parse(`json:"name,optional" validate:"required,min=2,max=10,len=6"`)
	assert.NoError(t, err)

	rules := validateRulesFromTags(tags)
	assert.True(t, rules.required)
	if assert.NotNil(t, rules.min) {
		assert.Equal(t, 2.0, *rules.min)
	}
	if assert.NotNil(t, rules.max) {
		assert.Equal(t, 10.0, *rules.max)
	}
	if assert.NotNil(t, rules.length) {
		assert.Equal(t, 6.0, *rules.length)
	}
}

func TestParametersFromType_ValidateMinMaxLen(t *testing.T) {
	ctx := testingContext(t)
	testStruct := apiSpec.DefineStruct{
		RawName: "Request",
		Members: []apiSpec.Member{
			{
				Name: "Name",
				Type: apiSpec.PrimitiveType{RawName: "string"},
				Tag:  `form:"name" validate:"min=2,max=10"`,
			},
			{
				Name: "Codes",
				Type: apiSpec.ArrayType{
					RawName: "[]string",
					Value:   apiSpec.PrimitiveType{RawName: "string"},
				},
				Tag: `form:"codes" validate:"len=3"`,
			},
		},
	}

	params := parametersFromType(ctx, http.MethodGet, testStruct)
	assert.Len(t, params, 2)

	assert.EqualValues(t, 2, *params[0].MinLength)
	assert.EqualValues(t, 10, *params[0].MaxLength)
	assert.EqualValues(t, 3, *params[1].MinItems)
	assert.EqualValues(t, 3, *params[1].MaxItems)
}

func TestPropertiesFromType_ValidateAndRequired(t *testing.T) {
	ctx := testingContext(t)
	testStruct := apiSpec.DefineStruct{
		RawName: "Request",
		Members: []apiSpec.Member{
			{
				Name: "Age",
				Type: apiSpec.PrimitiveType{RawName: "int"},
				Tag:  `json:"age,optional" validate:"required,min=1,max=120"`,
			},
		},
	}

	properties, required := propertiesFromType(ctx, testStruct)
	age := properties["age"]

	if assert.NotNil(t, age.Minimum) {
		assert.Equal(t, 1.0, *age.Minimum)
	}
	if assert.NotNil(t, age.Maximum) {
		assert.Equal(t, 120.0, *age.Maximum)
	}
	assert.Equal(t, []string{"age"}, required)
}

func TestPropertiesFromType_DescriptionFromDocs(t *testing.T) {
	ctx := testingContext(t)
	testStruct := apiSpec.DefineStruct{
		RawName: "Request",
		Members: []apiSpec.Member{
			{
				Name: "Name",
				Type: apiSpec.PrimitiveType{RawName: "string"},
				Tag:  `json:"name"`,
				Docs: []string{"// user name"},
			},
		},
	}

	properties, _ := propertiesFromType(ctx, testStruct)
	assert.Equal(t, "user name", properties["name"].Description)
}

func TestParametersFromType_DescriptionFromDocs(t *testing.T) {
	ctx := testingContext(t)
	testStruct := apiSpec.DefineStruct{
		RawName: "Request",
		Members: []apiSpec.Member{
			{
				Name: "Name",
				Type: apiSpec.PrimitiveType{RawName: "string"},
				Tag:  `form:"name"`,
				Docs: []string{"// user name for query"},
			},
		},
	}

	params := parametersFromType(ctx, http.MethodGet, testStruct)
	if assert.Len(t, params, 1) {
		assert.Equal(t, "user name for query", params[0].Description)
	}
}
