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
package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thmeitz/ksqldb-go"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push dogs example",
}

func init() {
	pushCmd.Run = push
	rootCmd.AddCommand(pushCmd)
}

func push(cmd *cobra.Command, args []string) {
	setLogger()
	host := viper.GetString("host")
	user := viper.GetString("username")
	password := viper.GetString("password")

	options := ksqldb.Options{
		Credentials: ksqldb.Credentials{Username: user, Password: password},
		BaseUrl:     host,
		AllowHTTP:   true,
	}

	client, err := ksqldb.NewClient(options, log.Current)
	if err != nil {
		log.Fatal(errors.Unwrap(err))
	}

	// You don't need to parse your ksql statement; Client.Pull parses it for you
	k := "SELECT ROWTIME, ID, NAME, DOGSIZE, AGE FROM DOGS EMIT CHANGES;"

	rc := make(chan ksqldb.Row)
	hc := make(chan ksqldb.Header, 1)

	// This Go routine will handle rows as and when they
	// are sent to the channel
	go func() {
		var dataTs float64
		var id string
		var name string
		var dogSize string
		var age string
		for row := range rc {
			if row != nil {
				// Should do some type assertions here
				dataTs = row[0].(float64)
				id = row[1].(string)
				name = row[2].(string)
				dogSize = row[3].(string)
				age = row[4].(string)

				// Handle the timestamp
				t := int64(dataTs)
				ts := time.Unix(t/1000, 0).Format(time.RFC822)

				log.Infof("🐾 New dog at %v: '%v' is %v and %v (id %v)\n", ts, name, dogSize, age, id)
			}
		}

	}()

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	e := ksqldb.Push(client, ctx, k, rc, hc)

	client.Close()

	if e != nil {
		log.Fatal(e)
	}
}
