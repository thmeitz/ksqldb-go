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

package logging_test

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

// TestLevelFromString tests mapping string levels to slog.Level
func TestLevelFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{
			name:     "debug level",
			input:    "debug",
			expected: slog.LevelDebug,
		},
		{
			name:     "DEBUG uppercase",
			input:    "DEBUG",
			expected: slog.LevelDebug,
		},
		{
			name:     "info level",
			input:    "info",
			expected: slog.LevelInfo,
		},
		{
			name:     "INFO uppercase",
			input:    "INFO",
			expected: slog.LevelInfo,
		},
		{
			name:     "warn level",
			input:    "warn",
			expected: slog.LevelWarn,
		},
		{
			name:     "WARN uppercase",
			input:    "WARN",
			expected: slog.LevelWarn,
		},
		{
			name:     "warning level alias",
			input:    "warning",
			expected: slog.LevelWarn,
		},
		{
			name:     "error level",
			input:    "error",
			expected: slog.LevelError,
		},
		{
			name:     "ERROR uppercase",
			input:    "ERROR",
			expected: slog.LevelError,
		},
		{
			name:     "trace level maps to debug",
			input:    "trace",
			expected: slog.LevelDebug,
		},
		{
			name:     "invalid level defaults to info",
			input:    "invalid",
			expected: slog.LevelInfo,
		},
		{
			name:     "empty string defaults to info",
			input:    "",
			expected: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logging.LevelFromString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNewLogger_TextFormat tests creating logger with text format
func TestNewLogger_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.NewLogger(&buf, slog.LevelInfo, "text")

	require.NotNil(t, logger)

	// Log a test message
	logger.Info("test message", "key", "value")

	output := buf.String()
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "key=value")
	// Text format should not contain JSON
	assert.NotContains(t, output, "{\"")
}

// TestNewLogger_JSONFormat tests creating logger with JSON format
func TestNewLogger_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.NewLogger(&buf, slog.LevelInfo, "json")

	require.NotNil(t, logger)

	// Log a test message
	logger.Info("test message", "key", "value")

	output := buf.String()
	assert.NotEmpty(t, output)

	// Verify it's valid JSON
	var logEntry map[string]interface{}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	err := json.Unmarshal([]byte(lines[0]), &logEntry)
	require.NoError(t, err)

	// Verify JSON structure
	assert.Equal(t, "test message", logEntry["msg"])
	assert.Equal(t, "value", logEntry["key"])
	assert.Contains(t, logEntry, "time")
	assert.Contains(t, logEntry, "level")
}

// TestNewLogger_InvalidFormat_DefaultsToText tests fallback to text format
func TestNewLogger_InvalidFormat_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.NewLogger(&buf, slog.LevelInfo, "invalid")

	require.NotNil(t, logger)

	// Log a test message
	logger.Info("test message")

	output := buf.String()
	assert.Contains(t, output, "test message")
	// Should use text format by default
	assert.NotContains(t, output, "{\"")
}

// TestNewLogger_AllLevels tests logger with different log levels
func TestNewLogger_AllLevels(t *testing.T) {
	tests := []struct {
		name      string
		level     slog.Level
		logFunc   func(*slog.Logger)
		shouldLog bool
	}{
		{
			name:  "debug message at debug level",
			level: slog.LevelDebug,
			logFunc: func(l *slog.Logger) {
				l.Debug("debug message")
			},
			shouldLog: true,
		},
		{
			name:  "debug message at info level",
			level: slog.LevelInfo,
			logFunc: func(l *slog.Logger) {
				l.Debug("debug message")
			},
			shouldLog: false,
		},
		{
			name:  "info message at info level",
			level: slog.LevelInfo,
			logFunc: func(l *slog.Logger) {
				l.Info("info message")
			},
			shouldLog: true,
		},
		{
			name:  "warn message at warn level",
			level: slog.LevelWarn,
			logFunc: func(l *slog.Logger) {
				l.Warn("warn message")
			},
			shouldLog: true,
		},
		{
			name:  "info message at warn level",
			level: slog.LevelWarn,
			logFunc: func(l *slog.Logger) {
				l.Info("info message")
			},
			shouldLog: false,
		},
		{
			name:  "error message at error level",
			level: slog.LevelError,
			logFunc: func(l *slog.Logger) {
				l.Error("error message")
			},
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := logging.NewLogger(&buf, tt.level, "text")

			tt.logFunc(logger)

			output := buf.String()
			if tt.shouldLog {
				assert.NotEmpty(t, output)
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

// TestDefaultLogger tests the default logger configuration
func TestDefaultLogger(t *testing.T) {
	logger := logging.DefaultLogger()

	require.NotNil(t, logger)

	// Should be able to log without panicking
	logger.Info("test message from default logger")
	logger.Debug("debug message from default logger")
	logger.Error("error message from default logger")

	// Test passes if no panic occurs
}

// TestNewLogger_WithAttributes tests logger with additional attributes
func TestNewLogger_WithAttributes(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.NewLogger(&buf, slog.LevelInfo, "json")

	// Add attributes
	logger = logger.With("app", "ksqldb-go", "version", "v2.0.0")

	logger.Info("test message")

	output := buf.String()
	assert.Contains(t, output, "ksqldb-go")
	assert.Contains(t, output, "v2.0.0")
}

// TestNewLogger_StructuredLogging tests structured logging capabilities
func TestNewLogger_StructuredLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.NewLogger(&buf, slog.LevelInfo, "json")

	logger.Info("operation completed",
		"duration_ms", 150,
		"status", "success",
		"records", 42,
	)

	output := buf.String()
	assert.NotEmpty(t, output)

	var logEntry map[string]interface{}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	err := json.Unmarshal([]byte(lines[0]), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, float64(150), logEntry["duration_ms"])
	assert.Equal(t, "success", logEntry["status"])
	assert.Equal(t, float64(42), logEntry["records"])
}

// TestNewLogger_CaseInsensitiveFormat tests format string case handling
func TestNewLogger_CaseInsensitiveFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{"lowercase json", "json"},
		{"uppercase JSON", "JSON"},
		{"mixed case Json", "Json"},
		{"lowercase text", "text"},
		{"uppercase TEXT", "TEXT"},
		{"mixed case Text", "Text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := logging.NewLogger(&buf, slog.LevelInfo, tt.format)

			require.NotNil(t, logger)
			logger.Info("test")
			assert.NotEmpty(t, buf.String())
		})
	}
}
