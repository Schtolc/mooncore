package dependencies

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	// This is a way to get path of the current file in runtime. It's nessary,
	// because binary and tests run from different dirs and we need to find config file
	// by relative path in both cases. Taken from https://stackoverflow.com/a/38644571
	_, b, _, _  = runtime.Caller(0)
	projectRoot = filepath.Join(filepath.Dir(b), "..")
)

// Config is in-code representation of config.yml.
type Config struct {
	Server struct {
		Hostbase struct {
			Port string `yaml:"port"`
			Host string `yaml:"host"`
		}
		Auth struct {
			Lifetime int `yaml:"lifetime"`
		}
	}
	Database struct {
		User         string `yaml:"user"`
		Dbname       string `yaml:"dbname"`
		Dialect      string `yaml:"dialect"`
		MaxOpenConns int    `yaml:"max_open_conns"`
	}
	Logs struct {
		Access string `yaml:"access"`
		Main   string `yaml:"main"`
	}
}

var configInstance *Config
var configMutex = &sync.Mutex{}

func getConfig() *Config {
	content, err := ioutil.ReadFile(filepath.Join(projectRoot, "config.yml"))
	if err != nil {
		logrus.Fatal(err)
	}
	conf := &Config{}
	err = yaml.Unmarshal([]byte(content), conf)
	if err != nil {
		logrus.Fatal(err)
	}
	return conf
}

// ConfigInstance returns config instance
func ConfigInstance() *Config {
	if configInstance != nil {
		return configInstance
	}
	configMutex.Lock()
	defer configMutex.Unlock()
	if configInstance == nil {
		configInstance = getConfig()
	}
	return configInstance
}
