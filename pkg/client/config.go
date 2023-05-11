package client

import (
	"encoding/json"
	"os"
)

type Config map[string]string

func ParseConfig(configFile string) (config *Config, err error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return
	}
	config = &Config{}
	err = json.Unmarshal(data, &config)
	return
}
