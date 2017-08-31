package cfg

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type AppConfig struct {
	Hostbase struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	}
	Log struct {
		Access string `yaml:"access"`
	}
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		Dialect  string `yaml:"dialect"`
	}
}

func getPath(file_name string) string {
	dir, err := os.Getwd()
	check_err(err)
	return dir + "/" + file_name
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

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
