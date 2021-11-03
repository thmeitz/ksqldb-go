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

	"github.com/Masterminds/log-go/impl/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/thmeitz/ksqldb-go"
)

var (
	logger = logrus.NewStandard()
)

func TestClientNotNil(t *testing.T) {
	client, _ := ksqldb.NewClient(ksqldb.Options{}, logger)
	assert.NotNil(t, client)
}

func TestClientPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The Client did not panic")
		}
	}()
	ksqldb.NewClient(ksqldb.Options{BaseUrl: "sf"}, logger)
}

func TestClientSanitizeQuery(t *testing.T) {
	client, _ := ksqldb.NewClient(ksqldb.Options{}, logger)
	sanitizedString := client.SanitizeQuery(`
	
		This is the 	house of Nicolas
	
	`)
	assert.Equal(t, "This is the house of Nicolas", sanitizedString)
}
