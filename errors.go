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

package ksqldb

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = errors.New("no result found")
)

type ResponseError struct {
	ErrType string `json:"@type"`
	ErrCode int    `json:"error_code"`
	Message string `json:"message"`
}

// Error gets the error string without new lines from ResponseError
func (e ResponseError) Error() string {
	// I don't like error messages with new lines
	e.Message = strings.ReplaceAll(e.Message, "\n", " ")
	return fmt.Sprintf("%v", e.Message)
}
