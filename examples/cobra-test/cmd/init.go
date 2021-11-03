package cmd

import (
	"fmt"
	"os"

	"github.com/Masterminds/log-go"
	lgrs "github.com/Masterminds/log-go/impl/logrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	lg = logrus.New()
)

func setLogger() {
	setLogFormat()
	setLogLevel()
	log.Current = lgrs.New(lg)
}

func setLogLevel() {
	loglevel := viper.GetString("loglevel")
	switch loglevel {
	case "info":
		lg.Level = logrus.InfoLevel
		return
	case "debug":
		lg.Level = logrus.DebugLevel
		return
	case "trace":
		lg.Level = logrus.TraceLevel
		return
	case "error":
		lg.Level = logrus.ErrorLevel
		return
	}
	fmt.Println("unknown log level")
	os.Exit(1)
}

func setLogFormat() {
	formatter := viper.GetString("logformat")
	switch formatter {
	case "text":
		lg.Formatter = &logrus.TextFormatter{}
		return
	case "json":
		fmt.Println("json formatter")
		lg.Formatter = &logrus.JSONFormatter{}
		return
	}
	fmt.Println("unknown formatter")
	os.Exit(1)
}
