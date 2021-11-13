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
package cmd

import (
	"context"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thmeitz/ksqldb-go"
	"github.com/thmeitz/ksqldb-go/net"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "print the dog stats",
}

func init() {
	pullCmd.Run = dogstats
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringP("dogsize", "d", "medium", "dogsizes are small|medium|large")
	if err := viper.BindPFlag("dogsize", pullCmd.Flags().Lookup("dogsize")); err != nil {
		log.Fatal(err)
	}
}

func dogstats(cmd *cobra.Command, args []string) {
	setLogger()
	host := viper.GetString("host")
	user := viper.GetString("username")
	password := viper.GetString("password")
	s := viper.GetString("dogsize")

	options := net.Options{
		Credentials: net.Credentials{Username: user, Password: password},
		BaseUrl:     host,
		AllowHTTP:   true,
	}

	kcl, err := ksqldb.NewClientWithOptions(options)
	if err != nil {
		log.Fatal(err)
	}
	defer kcl.Close()

	k := `select timestamptostring(windowstart,'yyyy-MM-dd HH:mm:ss','Europe/London') as window_start, 
	timestamptostring(windowend,'HH:mm:ss','Europe/London') as window_end, 
	dog_size, dogs_ct from dogs_by_size 
	where dog_size=?;`

	builder, err := ksqldb.DefaultQueryBuilder(k)
	if err != nil {
		log.Fatal(err)
	}

	stmnt, err := builder.Bind(s)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	_, r, err := kcl.Pull(ctx, *stmnt, true)
	if err != nil {
		log.Fatal(err)
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
			log.Infof("🐶 There are %v dogs size %v between %v and %v", dogsCt, dogSize, windowStart, windowEnd)
		}
	}
}
