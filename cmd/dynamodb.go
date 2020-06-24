/*
Copyright © 2020 cpipe contributors

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
	"cpipe/core"
	"cpipe/dynamodb"

	"github.com/spf13/cobra"
)

// dynamodbCmd represents the dynamodb command
var dynamodbCmd = &cobra.Command{
	Use:   "dynamodb",
	Short: "Stream data into dynamodb",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		writer := dynamodb.NewDebugWriter()
		transformer := &dynamodb.Transformer{}
		deserializer := &core.JsonDeserializer{}
		core.Start("", deserializer, transformer, writer, 1, 256, 10)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dynamodbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dynamodbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dynamodbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
