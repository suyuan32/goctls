package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suyuan32/goctls/api/spec"
)

func TestSpec2PathsWithRootRoute(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		routePath    string
		expectedPath string
	}{
		{
			name:         "prefix with root route",
			prefix:       "/api/v1/shoppings",
			routePath:    "/",
			expectedPath: "/api/v1/shoppings",
		},
		{
			name:         "prefix with sub route",
			prefix:       "/api/v1/shoppings",
			routePath:    "/list",
			expectedPath: "/api/v1/shoppings/list",
		},
		{
			name:         "empty prefix with root route",
			prefix:       "",
			routePath:    "/",
			expectedPath: "/",
		},
		{
			name:         "empty prefix with sub route",
			prefix:       "",
			routePath:    "/list",
			expectedPath: "/list",
		},
		{
			name:         "prefix with trailing slash and root route",
			prefix:       "/api/v1/shoppings/",
			routePath:    "/",
			expectedPath: "/api/v1/shoppings",
		},
		{
			name:         "prefix without leading slash and root route",
			prefix:       "api/v1/shoppings",
			routePath:    "/",
			expectedPath: "/api/v1/shoppings",
		},
		{
			name:         "single level prefix with root route",
			prefix:       "/api",
			routePath:    "/",
			expectedPath: "/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := spec.Service{
				Groups: []spec.Group{
					{
						Annotation: spec.Annotation{
							Properties: map[string]string{
								propertyKeyPrefix: tt.prefix,
							},
						},
						Routes: []spec.Route{
							{
								Method:  "get",
								Path:    tt.routePath,
								Handler: "TestHandler",
							},
						},
					},
				},
			}

			ctx := testingContext(t)
			paths := spec2Paths(ctx, srv)

			assert.Contains(t, paths.Paths, tt.expectedPath,
				"Expected path %s not found in generated paths. Got: %v",
				tt.expectedPath, paths.Paths)
		})
	}
}

func TestSpec2PathSummaryUsesRouteCommentByDefault(t *testing.T) {
	route := spec.Route{
		Method:  "get",
		Path:    "/users",
		Handler: "GetUsers",
		Comment: []string{"// list users"},
	}
	group := spec.Group{}

	pathItem := spec2Path(testingContext(t), group, route)

	if assert.NotNil(t, pathItem.Get) {
		assert.Equal(t, "list users", pathItem.Get.Summary)
		assert.Equal(t, "list users", pathItem.Get.Description)
	}
}

func TestSpec2PathSummaryAndDescriptionUseRouteDocByDefault(t *testing.T) {
	route := spec.Route{
		Method:  "get",
		Path:    "/users",
		Handler: "GetUsers",
		Doc:     []string{"// fetch user list"},
	}
	group := spec.Group{}

	pathItem := spec2Path(testingContext(t), group, route)

	if assert.NotNil(t, pathItem.Get) {
		assert.Equal(t, "fetch user list", pathItem.Get.Summary)
		assert.Equal(t, "fetch user list", pathItem.Get.Description)
	}
}

func TestSpec2PathSchemesFallbackToGlobalByDefault(t *testing.T) {
	route := spec.Route{
		Method:  "get",
		Path:    "/students",
		Handler: "GetStudents",
	}
	group := spec.Group{}

	pathItem := spec2Path(testingContext(t), group, route)

	if assert.NotNil(t, pathItem.Get) {
		assert.Nil(t, pathItem.Get.Schemes)
	}
}

func TestSpec2PathSchemesCanBeSetPerRoute(t *testing.T) {
	route := spec.Route{
		Method:  "get",
		Path:    "/students",
		Handler: "GetStudents",
		AtDoc: spec.AtDoc{
			Properties: map[string]string{
				propertyKeySchemes: "https",
			},
		},
	}
	group := spec.Group{}

	pathItem := spec2Path(testingContext(t), group, route)

	if assert.NotNil(t, pathItem.Get) {
		assert.Equal(t, []string{"https"}, pathItem.Get.Schemes)
	}
}
