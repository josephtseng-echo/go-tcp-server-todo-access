package log

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/onrik/logrus/filename"
	"time"
)

type LoggerCfg struct {
	Dir  string
	Name string
}

func NewLogger(cfg *LoggerCfg) *logrus.Logger {
	nowTime := time.Now()
	nowTimeStr := nowTime.Format("2006-01-02_15")
	logInfoName := cfg.Dir + "/" + cfg.Name + "/" + nowTimeStr + ".info"
	logErrorName := cfg.Dir + "/" + cfg.Name + "/" + nowTimeStr + ".info"

	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "2006/02/01 15:04:05.0000"
	Formatter.FullTimestamp = true

	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  logInfoName,
		logrus.ErrorLevel: logErrorName,
	}
	Log := logrus.New()
	Log.Formatter = Formatter
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{
			DisableTimestamp : false,
		},
	))
	Log.Hooks.Add(filename.NewHook())

	return Log
}
