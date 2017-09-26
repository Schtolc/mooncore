package main

import (
	"github.com/Gurpartap/logrus-stack"
	"github.com/Schtolc/mooncore/config"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

func openLogFile(filename string) *os.File {
	logfile, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	return logfile
}

// InitLogs sets logger format and hooks, redirects stdout and strerr to main logfile
func InitLogs(config config.Config) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.AddHook(logrus_stack.StandardHook())

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
	var parameters = middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format:  "${time_rfc3339} ${host} ${method} ${uri} ${status}\n",
		Output:  openLogFile(filename),
	}
	return parameters
}
