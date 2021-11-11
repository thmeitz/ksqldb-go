/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"github.com/Masterminds/log-go/impl/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thmeitz/ksqldb-go"
	"github.com/thmeitz/ksqldb-go/net"
)

// cstatsCmd represents the cstats command
var cstatsCmd = &cobra.Command{
	Use:   "cluster-status",
	Short: "get cluster status",
}

func init() {
	cstatsCmd.Run = cstats
	rootCmd.AddCommand(cstatsCmd)
}

func cstats(cmd *cobra.Command, args []string) {
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

	clusterStatus, err := kcl.GetClusterStatus()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", clusterStatus)

}
