/*
Copyright © 2021 Thomas Meitz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

func setLogger() {
	level := parseLoglevel()
	handler := createHandler(level)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func parseLoglevel() slog.Level {
	loglevel := viper.GetString("loglevel")
	switch loglevel {
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	case "trace":
		return slog.LevelDebug // slog doesn't have trace, use debug
	case "error":
		return slog.LevelError
	default:
		fmt.Println("unknown log level")
		os.Exit(1)
		return slog.LevelInfo // unreachable
	}
}

func createHandler(level slog.Level) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	formatter := viper.GetString("logformat")
	switch formatter {
	case "text":
		return slog.NewTextHandler(os.Stdout, opts)
	case "json":
		return slog.NewJSONHandler(os.Stdout, opts)
	default:
		fmt.Println("unknown formatter")
		os.Exit(1)
		return nil // unreachable
	}
}
