package osproxy

import (
	"os"

	"osproxy/api/v1alpha3"
	"osproxy/internal/utils"

	"gopkg.in/yaml.v3"
)

func parseConfig(filepath string) (config v1alpha3.OSProxyConfigT, err error) {
	configBytes, err := os.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	configBytes = []byte(os.ExpandEnv(string(configBytes)))

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return config, err
	}

	if _, ok := config.Proxy.Source.Buckets[utils.DefaultSourceKey]; !ok {
		config.Proxy.Source.Buckets[utils.DefaultSourceKey] = v1alpha3.BucketObjectConfigT{
			Bucket: "placeholder-bucket",
		}
	}

	return config, err
}
