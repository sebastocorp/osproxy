package osproxy

import (
	"os"
	"regexp"

	"osproxy/api/v1alpha5"

	"gopkg.in/yaml.v3"
)

func expandEnv(input string) string {
	re := regexp.MustCompile(`\${ENV:([A-Za-z_][A-Za-z0-9_]*)}\$`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		key := re.FindStringSubmatch(match)[1]
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return match
	})
}

func parseConfig(filepath string) (config v1alpha5.OSProxyConfigT, err error) {
	configBytes, err := os.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	configBytes = []byte(expandEnv(string(configBytes)))

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return config, err
	}

	return config, err
}
