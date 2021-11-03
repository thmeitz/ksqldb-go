package ksqldb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thmeitz/ksqldb-go"
)

func TestClientDebugModeAfterClientInit(t *testing.T) {
	t.Parallel()
	client := ksqldb.NewClient("http://example.com", "testuser", "testpassword", logger)
	assert.False(t, client.GetDebugMode(), false)
}

// func TestClientDebugMethod(t *testing.T) {
// 	t.Parallel()
// 	client := ksqldb.NewClient("http://example.com", "testuser", "testpassword")
// 	client.Debug()
// 	assert.True(t, client.GetDebugMode())
// }

func TestClientSetDebug(t *testing.T) {
	t.Parallel()
	client := ksqldb.NewClient("http://example.com", "testuser", "testpassword", logger)
	client.SetDebug(true)
	assert.True(t, client.GetDebugMode())
	client.SetDebug(false)
	assert.False(t, client.GetDebugMode())
}
