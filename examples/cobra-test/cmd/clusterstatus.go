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
	"log/slog"
	"os"

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

	options := net.Options{
		Credentials: net.Credentials{Username: user, Password: password},
		BaseUrl:     host,
	}

	kcl, err := ksqldb.NewClientWithOptions(options)
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
	defer kcl.Close()

	clusterStatus, err := kcl.GetClusterStatus(context.Background())
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}

	fmt.Printf("%+v", clusterStatus)
}
