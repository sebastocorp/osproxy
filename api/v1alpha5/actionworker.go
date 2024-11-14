package v1alpha5

import "time"

//--------------------------------------------------------------
// ACTION WORKER CONFIG
//--------------------------------------------------------------

type ActionWorkerConfigT struct {
	Loglevel       string         `yaml:"loglevel"`
	PoolCapacity   int            `yaml:"poolCapacity"`
	Type           string         `yaml:"type"`
	ScrapeInterval time.Duration  `yaml:"scrapeInterval"`
	Request        RequestConfigT `yaml:"request"`
}

type RequestConfigT struct {
	URL string `yaml:"url"`
}
