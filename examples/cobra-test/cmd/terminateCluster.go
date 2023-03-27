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
	"github.com/Masterminds/log-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thmeitz/ksqldb-go"
	"github.com/thmeitz/ksqldb-go/net"
)

// terminateClusterCmd represents the terminateCluster command
var terminateClusterCmd = &cobra.Command{
	Use:   "terminate-cluster",
	Short: "terminates your cluster",
	Long: `If you don't need your ksqlDB cluster anymore, 
you can terminate the cluster and clean up the resources 
using this command. To terminate a ksqlDB cluster, 
first shut down all of the servers, except one.`,
	Run: terminateCluster,
}

func init() {
	rootCmd.AddCommand(terminateClusterCmd)
}

func terminateCluster(cmd *cobra.Command, args []string) {
	setLogger()
	host := viper.GetString("host")
	user := viper.GetString("username")
	password := viper.GetString("password")

	var result *ksqldb.KsqlResponseSlice

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

	if result, err = kcl.TerminateCluster("DOGS_BY_SIZE", "dogs"); err != nil {
		log.Fatal(err)
	}

	log.Infof("%+v", result)
}
