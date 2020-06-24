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
	"time"

	"github.com/spf13/cobra"
)

var table string
var retryDelaySeconds int

// dynamodbCmd represents the dynamodb command
var dynamodbCmd = &cobra.Command{
	Use:   "dynamodb",
	Short: "Stream data into dynamodb",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var writer core.Writer
		if enableDebugWriter {
			writer = dynamodb.NewDebugWriter()
		} else {
			writer = dynamodb.NewWriter(table, time.Duration(retryDelaySeconds)*time.Second)
		}

		transformer := &dynamodb.Transformer{}
		deserializer := &core.JsonDeserializer{}
		if batchSize > 25 {
			batchSize = 25
		}
		core.Start("", deserializer, transformer, writer, batchSize, maxWriters, flushTimeoutSeconds)
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
	dynamodbCmd.Flags().StringVar(&table, "table", "", "dynamodb table name")
	dynamodbCmd.Flags().IntVar(&retryDelaySeconds, "retry-delay-seconds", 5, "number of seconds to wait before resending unprocessed items in a batch")
	dynamodbCmd.MarkFlagRequired("table")
}
