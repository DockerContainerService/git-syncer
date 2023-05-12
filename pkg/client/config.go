package client

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type Config map[string]string

func ParseConfig(configFile string) (config *Config, err error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		logrus.Errorf("Read config file %s failed: %v", configFile, err)
		err = fmt.Errorf("parse config file failed")
		return
	}
	config = &Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		logrus.Errorf("Parse config file %s failed: %v", configFile, err)
		err = fmt.Errorf("parse config file failed")
		return
	}
	return
}
