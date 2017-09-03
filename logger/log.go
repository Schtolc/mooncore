package logger

import (
	"github.com/labstack/echo/middleware"
	"log"
	"os"
)

// Configure is a build-in middleware function responsible for access logs in echo framework.
func Configure(filename string) middleware.LoggerConfig {
	logfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	var config = middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format:  "${time_rfc3339} ${host} ${method} ${uri} ${status}\n",
		Output:  logfile,
	}
	return config
}
