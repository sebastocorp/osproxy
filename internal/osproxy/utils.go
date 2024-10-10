package osproxy

import (
	"os"

	"osproxy/api/v1alpha4"

	"gopkg.in/yaml.v3"
)

func parseConfig(filepath string) (config v1alpha4.OSProxyConfigT, err error) {
	configBytes, err := os.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	configBytes = []byte(os.ExpandEnv(string(configBytes)))

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return config, err
	}

	return config, err
}
