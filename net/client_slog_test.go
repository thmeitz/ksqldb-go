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

package net_test

import (
	"bytes"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thmeitz/ksqldb-go/net"
)

// TestNewHTTPClient_WithSlogLogger tests that NewHTTPClient accepts *slog.Logger
func TestNewHTTPClient_WithSlogLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	options := net.Options{
		BaseUrl: "http://localhost:8088",
	}

	client, err := net.NewHTTPClient(options, logger)

	require.NoError(t, err)
	require.NotNil(t, client)

	// Ensure HTTPClient interface is implemented
	var _ net.HTTPClient = &client

	// Verify the client is properly configured
	assert.Equal(t, "http://localhost:8088", client.GetUrl(""))
}

// TestNewHTTPClient_WithNilLogger_UsesDefault tests that nil logger defaults to slog.Default()
func TestNewHTTPClient_WithNilLogger_UsesDefault(t *testing.T) {
	options := net.Options{
		BaseUrl: "http://localhost:8088",
	}

	client, err := net.NewHTTPClient(options, nil)

	require.NoError(t, err)
	require.NotNil(t, client)

	// Verify the client is properly configured even with nil logger
	assert.Equal(t, "http://localhost:8088", client.GetUrl(""))
}

// TestNewHTTPClient_WithSlogLogger_TextFormat tests text format logging
func TestNewHTTPClient_WithSlogLogger_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	options := net.Options{
		BaseUrl:   "http://localhost:8088",
		AllowHTTP: true,
	}

	client, err := net.NewHTTPClient(options, logger)

	require.NoError(t, err)
	require.NotNil(t, client)

	// If logger is used in the future, verify text format
	output := buf.String()
	// For now, just verify it doesn't crash with text handler
	_ = output
}

// TestNewHTTPClient_WithSlogLogger_JSONFormat tests JSON format logging
func TestNewHTTPClient_WithSlogLogger_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	options := net.Options{
		BaseUrl:   "http://localhost:8088",
		AllowHTTP: true,
	}

	client, err := net.NewHTTPClient(options, logger)

	require.NoError(t, err)
	require.NotNil(t, client)

	// If logger is used in the future, verify JSON format
	output := buf.String()
	// For now, just verify it doesn't crash with JSON handler
	_ = output
}

// TestClient_LogsConcurrently tests that logger is safe for concurrent use
func TestClient_LogsConcurrently(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	options := net.Options{
		BaseUrl:   "http://localhost:8088",
		AllowHTTP: true,
	}

	client, err := net.NewHTTPClient(options, logger)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Test concurrent operations that might use logger
	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(n int) {
			defer wg.Done()
			// Create multiple clients concurrently
			_, err := net.NewHTTPClient(options, logger)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()
	// Test passes if no race conditions occur
}

// TestNewHTTPClient_WithSlogLogger_AllLevels tests all log levels
func TestNewHTTPClient_WithSlogLogger_AllLevels(t *testing.T) {
	levels := []slog.Level{
		slog.LevelDebug,
		slog.LevelInfo,
		slog.LevelWarn,
		slog.LevelError,
	}

	options := net.Options{
		BaseUrl:   "http://localhost:8088",
		AllowHTTP: true,
	}

	for _, level := range levels {
		t.Run(level.String(), func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
				Level: level,
			}))

			client, err := net.NewHTTPClient(options, logger)

			require.NoError(t, err)
			require.NotNil(t, client)
		})
	}
}

// TestNewHTTPClient_InvalidURL_WithSlog tests error handling with slog logger
func TestNewHTTPClient_InvalidURL_WithSlog(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	options := net.Options{
		BaseUrl: "invalid-url",
	}

	client, err := net.NewHTTPClient(options, logger)

	require.Error(t, err)
	require.NotNil(t, client)
	assert.Contains(t, strings.ToLower(err.Error()), "invalid")
}

// TestNewHTTPClient_WithCredentials_Slog tests BasicAuth with slog logger
func TestNewHTTPClient_WithCredentials_Slog(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	options := net.Options{
		BaseUrl:   "http://localhost:8088",
		AllowHTTP: true,
		Credentials: net.Credentials{
			Username: "testuser",
			Password: "testpass",
		},
	}

	client, err := net.NewHTTPClient(options, logger)

	require.NoError(t, err)
	require.NotNil(t, client)

	// Verify BasicAuth is set
	auth := client.BasicAuth()
	assert.NotEmpty(t, auth)
}
