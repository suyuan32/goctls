package swagger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiSpec "github.com/suyuan32/goctls/api/spec"
)

func TestResolveSwaggerDefaults(t *testing.T) {
	root := t.TempDir()
	descDir := filepath.Join(root, "desc")
	etcDir := filepath.Join(root, "etc")
	require.NoError(t, os.MkdirAll(descDir, 0755))
	require.NoError(t, os.MkdirAll(etcDir, 0755))

	apiPath := filepath.Join(descDir, "user.api")
	require.NoError(t, os.WriteFile(apiPath, []byte("syntax = \"v1\"\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "Makefile"), []byte("SERVICE_STYLE = user-api\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(etcDir, "user-api.yaml"), []byte("Host: 0.0.0.0\nPort: 8888\n"), 0644))

	defaults := resolveSwaggerDefaults(apiPath, "")
	assert.Equal(t, "user-api", defaults.Title)
	assert.Equal(t, "0.0.0.0:8888", defaults.Host)
}

func TestResolveSwaggerDefaults_WithCustomMakefileRelativePath(t *testing.T) {
	workspace := t.TempDir()
	apiPath := filepath.Join(workspace, "any", "path", "api.api")
	require.NoError(t, os.MkdirAll(filepath.Dir(apiPath), 0755))
	require.NoError(t, os.WriteFile(apiPath, []byte("syntax = \"v1\"\n"), 0644))

	coreRoot := filepath.Join(workspace, "core")
	require.NoError(t, os.MkdirAll(filepath.Join(coreRoot, "etc"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(coreRoot, "Makefile"), []byte("SERVICE_STYLE = core-service\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(coreRoot, "etc", "core-service.yaml"), []byte("Host: 127.0.0.1\nPort: 9001\n"), 0644))

	oldWD, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Chdir(oldWD)
	})
	require.NoError(t, os.Chdir(workspace))

	defaults := resolveSwaggerDefaults(apiPath, filepath.Join("core", "Makefile"))
	assert.Equal(t, "core-service", defaults.Title)
	assert.Equal(t, "127.0.0.1:9001", defaults.Host)
}

func TestSpec2Swagger_GoZeroInfoTakesPriority(t *testing.T) {
	swaggerDoc, err := spec2Swagger(&apiSpec.ApiSpec{
		Info: apiSpec.Info{
			Properties: map[string]string{
				propertyKeyTitle: "Demo API",
				propertyKeyHost:  "api.example.com",
			},
		},
	}, swaggerDefaults{
		Title: "service-style-title",
		Host:  "localhost:8888",
	})
	require.NoError(t, err)
	assert.Equal(t, "Demo API", swaggerDoc.Info.Title)
	assert.Equal(t, "api.example.com", swaggerDoc.Host)
}
