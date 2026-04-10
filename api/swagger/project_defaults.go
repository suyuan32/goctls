package swagger

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type swaggerDefaults struct {
	Title string
	Host  string
}

var serviceStylePattern = regexp.MustCompile(`(?m)^\s*SERVICE_STYLE\s*(?::=|\?=|\+=|=)\s*([^#\r\n]+)`)

func resolveSwaggerDefaults(apiPath, makefilePath string) swaggerDefaults {
	makefileAbsPath, projectRoot, ok := resolveMakefilePath(apiPath, makefilePath)
	if !ok {
		return swaggerDefaults{}
	}

	serviceStyle, err := readServiceStyle(makefileAbsPath)
	if err != nil || serviceStyle == "" {
		return swaggerDefaults{}
	}

	defaults := swaggerDefaults{
		Title: serviceStyle,
	}
	host, port, err := readHostAndPortFromYAML(filepath.Join(projectRoot, "etc", serviceStyle+".yaml"))
	if err == nil && host != "" && port != "" {
		defaults.Host = host + ":" + port
	}
	return defaults
}

func resolveMakefilePath(apiPath, makefilePath string) (string, string, bool) {
	apiAbsPath, err := filepath.Abs(apiPath)
	if err != nil {
		return "", "", false
	}

	projectRoot, ok := projectRootFromAPIPath(apiAbsPath)
	if !ok {
		return "", "", false
	}

	if strings.TrimSpace(makefilePath) != "" {
		makefileAbsPath, err := filepath.Abs(makefilePath)
		if err != nil {
			return "", "", false
		}
		return makefileAbsPath, projectRoot, true
	}

	return filepath.Join(projectRoot, "Makefile"), projectRoot, true
}

func projectRootFromAPIPath(apiPath string) (string, bool) {
	dir := filepath.Dir(apiPath)
	for {
		if strings.EqualFold(filepath.Base(dir), "desc") {
			parent := filepath.Dir(dir)
			if parent == dir {
				return "", false
			}
			return parent, true
		}
		next := filepath.Dir(dir)
		if next == dir {
			return "", false
		}
		dir = next
	}
}

func readServiceStyle(makefilePath string) (string, error) {
	data, err := os.ReadFile(makefilePath)
	if err != nil {
		return "", err
	}

	matches := serviceStylePattern.FindSubmatch(data)
	if len(matches) < 2 {
		return "", nil
	}

	value := strings.TrimSpace(string(matches[1]))
	value = strings.Trim(value, `"'`)
	return value, nil
}

func readHostAndPortFromYAML(configPath string) (string, string, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", "", err
	}

	var raw interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return "", "", err
	}

	hostValue, hostOK := findValueByKey(raw, "host")
	portValue, portOK := findValueByKey(raw, "port")
	if !hostOK || !portOK {
		return "", "", nil
	}

	host := normalizeHost(hostValue)
	port := normalizePort(portValue)
	if host == "" || port == "" {
		return "", "", fmt.Errorf("invalid host/port in %s", configPath)
	}

	return host, port, nil
}

func findValueByKey(v interface{}, keyName string) (interface{}, bool) {
	switch val := v.(type) {
	case map[interface{}]interface{}:
		for key, item := range val {
			if strings.EqualFold(fmt.Sprint(key), keyName) {
				return item, true
			}
		}
		for _, item := range val {
			if found, ok := findValueByKey(item, keyName); ok {
				return found, true
			}
		}
	case map[string]interface{}:
		for key, item := range val {
			if strings.EqualFold(key, keyName) {
				return item, true
			}
		}
		for _, item := range val {
			if found, ok := findValueByKey(item, keyName); ok {
				return found, true
			}
		}
	case []interface{}:
		for _, item := range val {
			if found, ok := findValueByKey(item, keyName); ok {
				return found, true
			}
		}
	}
	return nil, false
}

func normalizeHost(v interface{}) string {
	return strings.Trim(strings.TrimSpace(fmt.Sprint(v)), `"'`)
}

func normalizePort(v interface{}) string {
	switch val := v.(type) {
	case int:
		return strconv.Itoa(val)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float32:
		return strconv.Itoa(int(val))
	case float64:
		return strconv.Itoa(int(val))
	case string:
		return strings.Trim(strings.TrimSpace(val), `"'`)
	default:
		return strings.Trim(strings.TrimSpace(fmt.Sprint(v)), `"'`)
	}
}
