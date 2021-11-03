/*
Copyright © 2021 Robin Moffat & Contributors
Copyright © 2021 Thomas Meitz <thme219@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Parts of this apiclient are borrowed from Zalando Skipper
https://github.com/zalando/skipper/blob/master/net/httpclient.go

Zalando licence: MIT
https://github.com/zalando/skipper/blob/master/LICENSE
*/

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
