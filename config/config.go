package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Server struct {
		Hostbase struct {
			Port string `yaml:"port"`
			Host string `yaml:"host"`
		}
		Logs struct {
			Access string `yaml:"access"`
		}
	}
	Database struct {
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		Dbname       string `yaml:"dbname"`
		Dialect      string `yaml:"dialect"`
		MaxOpenConns int    `yaml:"max_open_connections"`
	}
}

func Get() Config {
	content, err := ioutil.ReadFile("config.yaml")
	check_err(err)

	conf := Config{}
	err = yaml.Unmarshal([]byte(content), &conf)
	check_err(err)

	return conf
}

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
