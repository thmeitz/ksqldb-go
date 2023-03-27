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
	"fmt"

	"github.com/Masterminds/log-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thmeitz/ksqldb-go"
	"github.com/thmeitz/ksqldb-go/net"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays your server infos",
}

func init() {
	infoCmd.Run = info
	rootCmd.AddCommand(infoCmd)
}

func info(cmd *cobra.Command, args []string) {
	setLogger()
	host := viper.GetString("host")
	user := viper.GetString("username")
	password := viper.GetString("password")

	options := net.Options{
		Credentials: net.Credentials{Username: user, Password: password},
		BaseUrl:     host,
	}

	client, err := net.NewHTTPClient(options, nil)
	if err != nil {
		log.Fatal(err)
	}

	kcl, err := ksqldb.NewClient(&client)
	if err != nil {
		log.Fatal(err)
	}
	defer kcl.Close()

	info, err := kcl.GetServerInfo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("===== as console output")
	fmt.Println(fmt.Sprintf("Version        : %v", info.Version))
	fmt.Println(fmt.Sprintf("KSQLServiceID  : %v", info.KsqlServiceID))
	fmt.Println(fmt.Sprintf("KafkaClusterID : %v", info.KafkaClusterID))
	fmt.Println(fmt.Sprintf("ServerStatus   : %v", info.ServerStatus))
	fmt.Println("===== as info log")
	log.Current.Infow("server info", log.Fields{"version": info.Version, "ksqlServiceId": info.KsqlServiceID, "kafkaClusterId": info.KafkaClusterID, "serverStatus": info.ServerStatus})
}
