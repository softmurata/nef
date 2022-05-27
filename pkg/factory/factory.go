package factory

import (
	"fmt"
	"io/ioutil"

	"github.com/softmurata/nef/internal/logger"
	"gopkg.in/yaml.v2"
)

var (
	NefConfig  Config
	Configured bool
)

func init() {
	Configured = false
}

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	if content, err := ioutil.ReadFile(f); err != nil {
		return err
	} else {
		NefConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &NefConfig); yamlErr != nil {
			return yamlErr
		}
		Configured = true
	}

	return nil
}

func CheckConfigVersion() error {
	currentVersion := NefConfig.GetVersion()

	if currentVersion != NEF_EXPECTED_CONFIG_VERSION {
		return fmt.Errorf("config version is [%s], but expected is [%s].",
			currentVersion, NEF_EXPECTED_CONFIG_VERSION)
	}

	logger.CfgLog.Infof("config version [%s]", currentVersion)

	return nil
}
