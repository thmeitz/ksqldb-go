/*
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
*/

package ksqldb_test

import (
	"testing"
)

// var (
// 	logger = logrus.NewStandard()
// )

func TestClientNotNil(t *testing.T) {
	// client := ksqldb.NewClient("http://example.com", "testuser", "testpassword", logger)
	// assert.NotNil(t, client)
}

func TestClientSanitizeQuery(t *testing.T) {
	//client := ksqldb.NewClient("http://example.com", "testuser", "testpassword", logger)
	//	sanitizedString := client.SanitizeQuery(`
	//
	//	This is the 	house of Nicolas
	//
	//`)
	//	assert.Equal(t, "This is the house of Nicolas", sanitizedString)
}
