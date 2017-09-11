package logger

import (
	"github.com/Gurpartap/logrus-stack"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	defaultFormatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	defaultStackLevels = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
)

// OpenLogFile for add; create if not exists
func OpenLogFile(filename string) (logfile *os.File) {
	logfile, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	return
}

// Init main logger with params: [correct time, output and level]
func init() {
	logrus.SetFormatter(defaultFormatter)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.AddHook(logrus_stack.NewHook(defaultStackLevels, defaultStackLevels))
}

// Configure Access Log
func Configure(filename string) middleware.LoggerConfig {
	var config = middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format:  "${time_rfc3339} ${host} ${method} ${uri} ${status}\n",
		Output:  OpenLogFile(filename),
	}
	return config
}
