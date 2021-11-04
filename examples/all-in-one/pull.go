/*
Copyright © 2021 Robin Moffat & Contributors

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
	"context"
	"fmt"
	"time"

	"github.com/thmeitz/ksqldb-go"
)

func getDogStats(client *ksqldb.Client, s string) (e error) {

	k := "SELECT TIMESTAMPTOSTRING(WINDOWSTART,'yyyy-MM-dd HH:mm:ss','Europe/London') AS WINDOW_START, TIMESTAMPTOSTRING(WINDOWEND,'HH:mm:ss','Europe/London') AS WINDOW_END, DOG_SIZE, DOGS_CT FROM DOGS_BY_SIZE WHERE DOG_SIZE='" + s + "';"

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	_, r, e := ksqldb.Pull(client, ctx, k, true)

	if e != nil {
		// handle the error better here, e.g. check for no rows returned
		return fmt.Errorf("error running pull request against ksqlDB:\n%w", e)
	}

	var windowStart string
	var windowEnd string
	var dogSize string
	var dogsCt float64
	for _, row := range r {
		if row != nil {
			// Should do some type assertions here
			windowStart = row[0].(string)
			windowEnd = row[1].(string)
			dogSize = row[2].(string)
			dogsCt = row[3].(float64)
			fmt.Printf("🐶 There are %v dogs size %v between %v and %v\n", dogsCt, dogSize, windowStart, windowEnd)
		}
	}
	return nil
}
