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
	"io/ioutil"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thmeitz/ksqldb-go/parser"
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
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}

	if err := checkCmd.MarkFlagRequired("file"); err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}

}

func checksqlfile(cmd *cobra.Command, args []string) {
	setLogger()

	fname := viper.GetString("file")

	fbytes, err := ioutil.ReadFile(filepath.Clean(fname))
	if err != nil {
		slog.Error("failed to read file", "filename", fname, "error", err)
		os.Exit(1)
	}

	ksqlerr := parser.ParseSql(string(fbytes))
	if ksqlerr != nil {
		slog.Error("sql parser error", "error", ksqlerr)
		for _, e := range *ksqlerr {
			slog.Error("sql parser error at", "line", e.Line, "col", e.Column, "message", e.Msg)
		}
		return
	}
	slog.Info("no errors found")
}
