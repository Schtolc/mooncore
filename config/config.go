package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Server struct {
		Hostbase struct {
			Port string `yaml:"port"`
			Host string `yaml:"host"`
		}
	}
	Database struct {
		User         string `yaml:"user"`
		Host         string `yaml:"host"`
		Dbname       string `yaml:"dbname"`
		Dialect      string `yaml:"dialect"`
		MaxOpenConns int    `yaml:"max_open_conns"`
	}
	Logs struct {
		Access      string `yaml:"access"`
		Main        string `yaml:"main"`
		TelegramBot struct {
			ChatName  string `yaml:"chatname"`
			AuthToken string `yaml:"authtoken"`
			ChatId    string `yaml:"chatid"`
		}
	}
}

// Get yaml config
func Get() (conf Config) {
	content, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	conf = Config{}
	err = yaml.Unmarshal([]byte(content), &conf)
	if err != nil {
		log.Fatal(err)
	}
	return
}
