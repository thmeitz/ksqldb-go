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

package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thmeitz/ksqldb-go/internal/logging"
)

// TestCobraCommands_TextFormat tests commands with text log format
func TestCobraCommands_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.NewLogger(&buf, slog.LevelInfo, "text")

	require.NotNil(t, logger)

	// Test that logger produces text output
	logger.Info("test command executed", "command", "health", "status", "success")

	output := buf.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "test command executed")
	assert.Contains(t, output, "command=health")
	assert.Contains(t, output, "status=success")

	// Verify it's NOT JSON
	assert.NotContains(t, output, "{\"")
}

// TestCobraCommands_JSONFormat tests commands with JSON log format
func TestCobraCommands_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.NewLogger(&buf, slog.LevelInfo, "json")

	require.NotNil(t, logger)

	// Test that logger produces JSON output
	logger.Info("test command executed",
		"command", "health",
		"status", "success",
		"duration_ms", 150,
	)

	output := buf.String()
	assert.NotEmpty(t, output)

	// Verify it's valid JSON
	var logEntry map[string]interface{}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	require.NotEmpty(t, lines)

	err := json.Unmarshal([]byte(lines[0]), &logEntry)
	require.NoError(t, err)

	// Verify JSON structure
	assert.Equal(t, "test command executed", logEntry["msg"])
	assert.Equal(t, "health", logEntry["command"])
	assert.Equal(t, "success", logEntry["status"])
	assert.Equal(t, float64(150), logEntry["duration_ms"])
}

// TestCobraCommands_LogLevels tests all log levels work correctly
func TestCobraCommands_LogLevels(t *testing.T) {
	tests := []struct {
		name        string
		levelString string
		level       slog.Level
	}{
		{
			name:        "debug level",
			levelString: "debug",
			level:       slog.LevelDebug,
		},
		{
			name:        "info level",
			levelString: "info",
			level:       slog.LevelInfo,
		},
		{
			name:        "warn level",
			levelString: "warn",
			level:       slog.LevelWarn,
		},
		{
			name:        "error level",
			levelString: "error",
			level:       slog.LevelError,
		},
		{
			name:        "trace level (maps to debug)",
			levelString: "trace",
			level:       slog.LevelDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test LevelFromString conversion
			level := logging.LevelFromString(tt.levelString)
			assert.Equal(t, tt.level, level)

			// Test logger with this level
			var buf bytes.Buffer
			logger := logging.NewLogger(&buf, level, "text")
			require.NotNil(t, logger)

			// Log at different levels
			logger.Debug("debug message")
			logger.Info("info message")
			logger.Warn("warn message")
			logger.Error("error message")

			// Verify appropriate messages appear based on level
			output := buf.String()
			switch tt.level {
			case slog.LevelDebug:
				assert.Contains(t, output, "debug message")
				assert.Contains(t, output, "info message")
			case slog.LevelInfo:
				assert.NotContains(t, output, "debug message")
				assert.Contains(t, output, "info message")
			case slog.LevelWarn:
				assert.NotContains(t, output, "info message")
				assert.Contains(t, output, "warn message")
			case slog.LevelError:
				assert.NotContains(t, output, "warn message")
				assert.Contains(t, output, "error message")
			}
		})
	}
}

// TestCobraCommands_LoggerInitialization tests logger initialization scenarios
func TestCobraCommands_LoggerInitialization(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		level    string
		validate func(*testing.T, *slog.Logger, *bytes.Buffer)
	}{
		{
			name:   "text format with debug level",
			format: "text",
			level:  "debug",
			validate: func(t *testing.T, logger *slog.Logger, buf *bytes.Buffer) {
				logger.Debug("debug test")
				assert.Contains(t, buf.String(), "debug test")
			},
		},
		{
			name:   "json format with info level",
			format: "json",
			level:  "info",
			validate: func(t *testing.T, logger *slog.Logger, buf *bytes.Buffer) {
				logger.Info("info test")
				output := buf.String()
				assert.NotEmpty(t, output)
				var logEntry map[string]interface{}
				lines := strings.Split(strings.TrimSpace(output), "\n")
				err := json.Unmarshal([]byte(lines[0]), &logEntry)
				require.NoError(t, err)
				assert.Equal(t, "info test", logEntry["msg"])
			},
		},
		{
			name:   "text format with error level",
			format: "text",
			level:  "error",
			validate: func(t *testing.T, logger *slog.Logger, buf *bytes.Buffer) {
				logger.Info("info test")
				logger.Error("error test")
				output := buf.String()
				assert.NotContains(t, output, "info test")
				assert.Contains(t, output, "error test")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := logging.LevelFromString(tt.level)
			var buf bytes.Buffer
			logger := logging.NewLogger(&buf, level, tt.format)
			require.NotNil(t, logger)

			tt.validate(t, logger, &buf)
		})
	}
}

// TestCobraCommands_StructuredLogging tests structured logging for commands
func TestCobraCommands_StructuredLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.NewLogger(&buf, slog.LevelInfo, "json")

	// Simulate logging from a cobra command
	logger.Info("command started", "cmd", "pull")
	logger.Info("connecting to ksqldb",
		"host", "localhost:8088",
		"timeout_ms", 5000,
	)
	logger.Info("query executed",
		"sql", "SELECT * FROM stream",
		"rows", 42,
		"duration_ms", 250,
	)
	logger.Info("command completed", "cmd", "pull", "status", "success")

	output := buf.String()
	assert.NotEmpty(t, output)

	// Verify all log entries are valid JSON
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Len(t, lines, 4)

	for _, line := range lines {
		var logEntry map[string]interface{}
		err := json.Unmarshal([]byte(line), &logEntry)
		require.NoError(t, err)
		assert.Contains(t, logEntry, "msg")
		assert.Contains(t, logEntry, "level")
		assert.Contains(t, logEntry, "time")
	}
}

// TestCobraCommands_DefaultLogger tests the default logger
func TestCobraCommands_DefaultLogger(t *testing.T) {
	logger := logging.DefaultLogger()
	require.NotNil(t, logger)

	// Test that default logger doesn't panic
	logger.Info("test from default logger")
	logger.Debug("debug from default logger")
	logger.Error("error from default logger")

	// Test passes if no panic
}

// TestCobraCommands_LoggerWithContext tests logger with contextual attributes
func TestCobraCommands_LoggerWithContext(t *testing.T) {
	var buf bytes.Buffer
	baseLogger := logging.NewLogger(&buf, slog.LevelInfo, "json")

	// Add command-specific context
	cmdLogger := baseLogger.With(
		"cmd", "health",
		"version", "v2.0.0",
		"host", "localhost:8088",
	)

	cmdLogger.Info("checking server health")
	cmdLogger.Info("server responding", "status_code", 200)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	require.Len(t, lines, 2)

	// Both log entries should have the context attributes
	for _, line := range lines {
		var logEntry map[string]interface{}
		err := json.Unmarshal([]byte(line), &logEntry)
		require.NoError(t, err)
		assert.Equal(t, "health", logEntry["cmd"])
		assert.Equal(t, "v2.0.0", logEntry["version"])
		assert.Equal(t, "localhost:8088", logEntry["host"])
	}
}

// TestCobraCommands_FormatCaseInsensitive tests case-insensitive format handling
func TestCobraCommands_FormatCaseInsensitive(t *testing.T) {
	formats := []string{"text", "TEXT", "Text", "json", "JSON", "Json"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			logger := logging.NewLogger(&buf, slog.LevelInfo, format)
			require.NotNil(t, logger)

			logger.Info("test message")
			assert.NotEmpty(t, buf.String())
		})
	}
}

// TestCobraCommands_InvalidFormatDefaultsToText tests invalid format fallback
func TestCobraCommands_InvalidFormatDefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.NewLogger(&buf, slog.LevelInfo, "invalid-format")

	require.NotNil(t, logger)

	logger.Info("test message")

	output := buf.String()
	assert.Contains(t, output, "test message")
	// Should default to text format (not JSON)
	assert.NotContains(t, output, "{\"")
}
