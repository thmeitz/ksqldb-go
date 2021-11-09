/*
Copyright © 2021 Robin Moffat & Contributors
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

package ksqldb

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = errors.New("no result found")
)

type Error struct {
	ErrType string `json:"@type"`
	ErrCode int    `json:"error_code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	// I don't like error messages with new lines
	e.Message = strings.ReplaceAll(e.Message, "\n", " ")
	return fmt.Sprintf("%v", e.Message)
}

func (e *Error) Is(target error) bool {
	_, ok := target.(*Error)
	return ok
}
