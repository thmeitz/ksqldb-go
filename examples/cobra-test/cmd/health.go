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
	"fmt"

	"github.com/Masterminds/log-go"
	"github.com/Masterminds/log-go/impl/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thmeitz/ksqldb-go"
	"github.com/thmeitz/ksqldb-go/net"
)

// healthCmd represents the serverhealth command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "display the server state of your servers",
}

func init() {
	healthCmd.Run = health
	rootCmd.AddCommand(healthCmd)
}

func health(cmd *cobra.Command, args []string) {
	setLogger()

	host := viper.GetString("host")
	user := viper.GetString("username")
	password := viper.GetString("password")

	log.Current = logrus.NewStandard()

	options := net.Options{
		Credentials: net.Credentials{Username: user, Password: password},
		BaseUrl:     host,
	}

	kcl, err := ksqldb.NewClientWithOptions(options)
	if err != nil {
		log.Fatal(err)
	}
	defer kcl.Close()

	health, err := kcl.GetServerStatus(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Overall healthiness   : %v", GoodOrBad(*health.IsHealthy)))
	fmt.Println(fmt.Sprintf("Kafka healthiness     : %v", GoodOrBad(*health.Details.Kafka.IsHealthy)))
	fmt.Println(fmt.Sprintf("Metastore healthiness : %v", GoodOrBad(*health.Details.Metastore.IsHealthy)))
}

func GoodOrBad(healthiness bool) string {
	if healthiness {
		return "healthy"
	}
	return "unhealthy"
}
