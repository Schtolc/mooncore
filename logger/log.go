package logger

import (
	"github.com/Gurpartap/logrus-stack"
	"github.com/Schtolc/mooncore/config"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"os"
	"runtime/debug"
)

var (
	defaultFormatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	defaultStackLevels = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
)

func openLogFile(filename string) (logfile *os.File) {
	logfile, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal(string(debug.Stack()))
	}
	return
}

// Configure main logger params: [correct time, output and level]
func Init(conf config.Config) {
	logrus.SetFormatter(defaultFormatter)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(openLogFile(conf.Logs.Main))
	logrus.AddHook(logrus_stack.NewHook(defaultStackLevels, defaultStackLevels))
}

// Configure access logger params: [correct time, output, level, fields]
func Log(conf config.Config) echo.MiddlewareFunc {
	log := logrus.New()
	log.Formatter = defaultFormatter
	log.Out = openLogFile(conf.Logs.Access)
	logrus.AddHook(logrus_stack.NewHook(defaultStackLevels, defaultStackLevels))

	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			log.WithFields(logrus.Fields{
				"method": req.Method,
				"host":   req.Host,
				"url":    req.URL,
				"status": res.Status,
			}).Info()
			defer CatchError()
			return h(c)
		}
	}
}

// Catch panic and log stack trace
func CatchError() {
	if e := recover(); e != nil {
		logrus.Error(e)
	}
}
