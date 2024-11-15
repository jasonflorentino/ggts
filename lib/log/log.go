package log

import (
	"ggts/lib/env"
	"io"

	"github.com/labstack/echo/v4"
	gommonlog "github.com/labstack/gommon/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func To(c echo.Context) echo.Logger {
	return c.Echo().Logger
}

func ToFile(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10, // megabytes
		MaxBackups: 3,  // files to keep
		MaxAge:     28, // days
	}
}

func Lvl() gommonlog.Lvl {
	var logLevel gommonlog.Lvl = env.LogLevel()
	var logLevelIsSet bool = logLevel != gommonlog.Lvl(0)
	if logLevelIsSet {
		return logLevel
	} else if env.IsProd() {
		return gommonlog.INFO
	} else {
		return gommonlog.DEBUG
	}
}
