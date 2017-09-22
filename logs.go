package main

import (
	"github.com/Gurpartap/logrus-stack"
	"github.com/Schtolc/mooncore/utils"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

var (
	defaultFormatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	defaultStackLevels = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
)

func openLogFile(filename string) *os.File {
	logfile, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	return logfile
}

// InitLogs sets logger format and hooks, redirects stdout and strerr to main logfile
func InitLogs(config utils.Config) {
	logrus.SetFormatter(defaultFormatter)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.AddHook(logrus_stack.NewHook(defaultStackLevels, defaultStackLevels))

	logfile := openLogFile(config.Logs.Main)

	if err := syscall.Dup2(int(logfile.Fd()), int(os.Stderr.Fd())); err != nil {
		logrus.Fatal(err)
	}
	if err := syscall.Dup2(int(logfile.Fd()), int(os.Stdout.Fd())); err != nil {
		logrus.Fatal(err)
	}
}

// GetAccessConfig returns config for access logs used in echo middleware
func GetAccessConfig(filename string) middleware.LoggerConfig {
	var config = middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format:  "${time_rfc3339} ${host} ${method} ${uri} ${status}\n",
		Output:  openLogFile(filename),
	}
	return config
}
