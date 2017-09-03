package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// Config is in-code representation of config.yml.
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
		MaxOpenConns int    `yaml:"max_open_conns"`
	}
}

// Get reads config.yml and return filled Config struct. If any error occurs program is terminated.
func Get() Config {
	content, err := ioutil.ReadFile("config.yml")
	checkErr(err)

	conf := Config{}
	err = yaml.Unmarshal([]byte(content), &conf)
	checkErr(err)

	return conf
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
