package configures

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConsoleConfig struct {
	Port        int    `yaml:"port"`
	AdminSecret string `yaml:"adminSecret"`

	Log struct {
		LogPath string `yaml:"logPath"`
		LogName string `yaml:"logName"`
	} `yaml:"log"`

	Mysql struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Address  string `yaml:"address"`
		DbName   string `yaml:"name"`
		Debug    bool   `yaml:"debug"`
	} `yaml:"mysql"`

	ImApiDomain string `yaml:"imApiDomain"`
}

var Config ConsoleConfig
var Env string

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

func InitConfigures() error {
	cfBytes, err := os.ReadFile("conf/config.yml")
	if err == nil {
		var conf ConsoleConfig
		yaml.Unmarshal(cfBytes, &conf)
		Config = conf
		if Config.Port <= 0 {
			Config.Port = 8090
		}
		return nil
	} else {
		return err
	}
}
