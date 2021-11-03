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

package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thmeitz/ksqldb-go/internal"
)

var tests = []struct {
	name    string
	url     string
	message string
}{
	{"empty url", "", "invalid host name given"},
	{"check localhost with port", "http://localhost:8123", ""},
	{"check localhost without port", "http://localhost", ""},
	{"invalid url scheme", "httpx://google.com", "invalid url scheme given"},
	{"invalid protocol", "stomp://localhost", "invalid url scheme given"},
	{"invalid host", "httpx://", "invalid host name given"},
	{"empty url scheme", "://hostname", "can't parse url: parse \"://hostname\": missing protocol scheme"},
	{"check url parser", "http://ahsd^öf023as", "can't parse url: parse \"http://ahsd^öf023as\": invalid character \"^\" in host name"},
}

func TestGetUrlValid(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := internal.GetUrl(tt.url)
			// fmt.Println(fmt.Sprintf("%v. %v: %v: %v => isValid:%v error: %v", idx, tt.name, tt.url, tt.want, valid, err))
			if err != nil {
				assert.Equal(t, tt.message, err.Error())
			}
		})
	}
}
