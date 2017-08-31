package cfg

import (
	"log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AppConfig struct {
	Hostbase struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	}
	Log struct {
		Access string `yaml:"access"`
	}
}

type BDConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Dialect  string `yaml:"dialect"`
}

func getPath(file_name string) string {
	return "/Users/v.suslova/src/mooncore/" + file_name
}

func GetAppConfig(config_name string) AppConfig {
	conf_path := getPath(config_name)
	conf := AppConfig{}

	content, err := ioutil.ReadFile(conf_path)
	check_err(err)

	err = yaml.Unmarshal([]byte(content), &conf)
	check_err(err)
	return conf
}
func GetBdConfig(config_name string) BDConfig {
	conf_path := getPath(config_name)
	conf := BDConfig{}

	content, err := ioutil.ReadFile(conf_path)
	check_err(err)

	err = yaml.Unmarshal([]byte(content), &conf)
	check_err(err)
	return conf
}

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}