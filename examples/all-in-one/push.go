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

func getDogUpdates(client *ksqldb.Client) (err error) {

	rc := make(chan ksqldb.Row)
	hc := make(chan ksqldb.Header, 1)

	k := "SELECT ROWTIME, ID, NAME, DOGSIZE, AGE FROM DOGS EMIT CHANGES;"

	// This Go routine will handle rows as and when they
	// are sent to the channel
	go func() {
		var DATA_TS float64
		var ID string
		var NAME string
		var DOG_SIZE string
		var AGE string
		for row := range rc {
			if row != nil {
				// Should do some type assertions here
				DATA_TS = row[0].(float64)
				ID = row[1].(string)
				NAME = row[2].(string)
				DOG_SIZE = row[3].(string)
				AGE = row[4].(string)

				// Handle the timestamp
				t := int64(DATA_TS)
				ts := time.Unix(t/1000, 0).Format(time.RFC822)

				fmt.Printf("🐾New dog at %v: '%v' is %v and %v (id %v)\n", ts, NAME, DOG_SIZE, AGE, ID)
			}
		}

	}()

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	e := client.Push(ctx, k, rc, hc)

	if e != nil {
		// handle the error better here, e.g. check for no rows returned
		return fmt.Errorf("error running push request against ksqlDB:\n%v", e)
	}

	return nil
}
