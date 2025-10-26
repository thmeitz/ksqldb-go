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

package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// LevelFromString converts a string log level to slog.Level.
// Supported levels: debug, info, warn, warning, error, trace (maps to debug).
// Invalid or empty strings default to info level.
// Case-insensitive.
func LevelFromString(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "trace":
		// slog doesn't have trace, map to debug
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

// NewLogger creates a new slog.Logger with the specified output, level, and format.
// Format can be "text" or "json" (case-insensitive). Invalid formats default to text.
func NewLogger(w io.Writer, level slog.Level, format string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	default:
		// Default to text format for invalid or "text" values
		handler = slog.NewTextHandler(w, opts)
	}

	return slog.New(handler)
}

// DefaultLogger returns a default logger that writes to stderr with info level and text format.
func DefaultLogger() *slog.Logger {
	return NewLogger(os.Stderr, slog.LevelInfo, "text")
}
