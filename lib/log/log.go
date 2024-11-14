package log

import (
	"io"

	"github.com/labstack/echo/v4"
	"gopkg.in/natefinch/lumberjack.v2"
)

func To(c echo.Context) echo.Logger {
	return c.Echo().Logger
}

func ToFile() io.Writer {
	return &lumberjack.Logger{
		Filename:   "ggts_log.txt",
		MaxSize:    10, // megabytes
		MaxBackups: 3,  // files to keep
		MaxAge:     28, // days
	}
}
