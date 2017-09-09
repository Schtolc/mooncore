package logger

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/labstack/echo"
	"github.com/rossmcdonald/telegram_hook"
	"github.com/sirupsen/logrus"
	"os"
	"runtime/debug"
)

var (
	defaultFormatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
)

func openLogFile(filename string) (logfile *os.File) {
	logfile, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	return
}
func configureTelegramHook(conf config.Config) (hook logrus.Hook) {
	hook, err := telegram_hook.NewTelegramHook(
		conf.Logs.TelegramBot.ChatName,
		conf.Logs.TelegramBot.AuthToken,
		conf.Logs.TelegramBot.ChatId,
	)
	if err != nil {
		logrus.Fatal(err)
	}
	return
}

// Configure main logger params: [correct time, output and level]
func Init(conf config.Config) {
	logrus.SetFormatter(defaultFormatter)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(openLogFile(conf.Logs.Main))
	logrus.AddHook(configureTelegramHook(conf))
}

// Configure access logger params: [correct time, output, level, fields]
func Log(conf config.Config) echo.MiddlewareFunc {
	log := logrus.New()
	log.Formatter = defaultFormatter
	log.Out = openLogFile(conf.Logs.Access)
	logrus.AddHook(configureTelegramHook(conf))

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

func CatchError() {
	if e := recover(); e != nil {
		logrus.WithFields(logrus.Fields{
			"error": e,
		}).Error(string(debug.Stack()))
	}
}
