package logger

import (
	"github.com/labstack/echo/middleware"
	"log"
	"os"
)

func Configure(filename string) middleware.LoggerConfig {
	logfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	var config = middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format:  "${time} ${host} ${method} ${uri} ${status}\n",
		Output:  logfile,
	}
	return config
}
