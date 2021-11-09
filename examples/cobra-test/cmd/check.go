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
	"io/ioutil"
	"path/filepath"

	"github.com/Masterminds/log-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thmeitz/ksqldb-go"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check a <example>.ksql file with the integrated parser",
}

func init() {
	checkCmd.Run = checksqlfile
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringP("file", "f", "", "filename to check sql syntax")
	if err := viper.BindPFlag("file", checkCmd.Flags().Lookup("file")); err != nil {
		log.Fatal(err)
	}

	if err := checkCmd.MarkFlagRequired("file"); err != nil {
		log.Fatal(err)
	}

}

func checksqlfile(cmd *cobra.Command, args []string) {
	setLogger()

	fname := viper.GetString("file")

	fbytes, err := ioutil.ReadFile(filepath.Clean(fname))
	if err != nil {
		log.Fatalf("%v %w", fname, err)
	}

	ksqlerr := ksqldb.ParseKSQL(string(fbytes))
	if ksqlerr != nil {
		log.Errorw("sql parser error", log.Fields{"error": ksqlerr})
		for _, e := range *ksqlerr {
			log.Errorw("sql parser error at", log.Fields{"line": e.Line, "col": e.Column, "message": e.Msg})
		}
		return
	}
	log.Info("no errors found")
}
