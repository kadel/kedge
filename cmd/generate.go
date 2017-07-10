/*
Copyright 2017 The Kedge Authors All rights reserved.

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
	"os"

	pkgcmd "github.com/kedgeproject/kedge/pkg/cmd"
	"github.com/spf13/cobra"
)

// Variables
var (
	AppFiles []string
)

// Represents the "generate" command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Kubernetes resources from App definition",
	Run: func(cmd *cobra.Command, args []string) {
		if err := pkgcmd.Generate(AppFiles); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func init() {
	generateCmd.Flags().StringArrayVarP(&AppFiles, "files", "f", []string{}, "input files")
	RootCmd.AddCommand(generateCmd)
}
