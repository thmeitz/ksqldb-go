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

package internal

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ValidateUrl checks the url; url must not contain a trailing slash
func GetUrl(path string) (*url.URL, error) {
	trimmedPath := strings.TrimSuffix(path, "/")
	u, err := url.Parse(trimmedPath)
	if err != nil {
		return nil, fmt.Errorf("can't parse url: %w", err)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("invalid host name given")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("invalid url scheme given")
	}
	return u, nil
}

// SanitizeQuery sanitizes the given content
//
// eventually we can use the KSqlParser to rewrite the query
func SanitizeQuery(content string) string {
	r := regexp.MustCompile(`\s+`)
	content = r.ReplaceAllString(content, " ")
	return strings.TrimSpace(content)
}
