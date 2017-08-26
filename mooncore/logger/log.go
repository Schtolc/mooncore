package logger

import(
	"github.com/labstack/echo/middleware"
	"os"
	"log"
)


func Configure(filename string) middleware.LoggerConfig {
	logfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644 )
	if err != nil {
		log.Fatal(err)
	}
	var config = middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `"method=${method}, "` +
			`"id=${id}, "` +
			`"host=${host}, "` +
			`"uri=${uri}, "` +
			`"status=${status}, "` +
			`"time=${time_rfc3339_nano}, "` +
			`"bytes_in=${bytes_in}, "` +
			`"bytes_out=${bytes_out}, "` +
			"\n",
		Output:  logfile,
	}
	return config
}


