package factory

import (
	logger_util "github.com/free5gc/util/logger"
	"github.com/softmurata/freeopenapi/models"
)

const (
	NEF_EXPECTED_CONFIG_VERSION = "1.0.1"
)

type Config struct {
	Info          *Info               `yaml:"info" valid:"required"`
	Configuration *Configuration      `yaml:"configuration" valid:"required"`
	Logger        *logger_util.Logger `yaml:"logger" valid:"optional"`
}

type Info struct {
	Version     string `yaml:"version,omitempty" valid:"type(string),required"`
	Description string `yaml:"description,omitempty" valid:"type(string),optional"`
}

const (
	NEF_DEFAULT_IPV4     = "127.0.0.15"
	NEF_DEFAULT_PORT     = "8000"
	NEF_DEFAULT_PORT_INT = 8000
)

type Configuration struct {
	NefName         string               `yaml:"nefName,omitempty"`
	Sbi             *Sbi                 `yaml:"sbi" valid:"required"`
	ServiceNameList []models.ServiceName `yaml:"serviceNameList"`
	NrfUri          string               `yaml:"nrfUri" valid:"url,required"`
}

type Sbi struct {
	Scheme       models.UriScheme `yaml:"scheme" valid:"scheme,required"`
	RegisterIPv4 string           `yaml:"registerIPv4,omitempty" valid:"host,optional"` // IP that is registered at NRF.
	// IPv6Addr string `yaml:"ipv6Addr,omitempty"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty" valid:"host,optional"` // IP used to run the server in the node.
	Port        int    `yaml:"port" valid:"port,required"`
	Tls         *Tls   `yaml:"tls,omitempty" valid:"optional"`
}

type Tls struct {
	Pem string `yaml:"pem,omitempty" valid:"type(string),minstringlength(1),required"`
	Key string `yaml:"key,omitempty" valid:"type(string),minstringlength(1),required"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}
