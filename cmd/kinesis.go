/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"cpipe/kinesis"
	"cpipe/utils"

	"github.com/spf13/cobra"
)

var stream string

// kinesisCmd represents the kinesis command
var kinesisCmd = &cobra.Command{
	Use:   "kinesis",
	Short: "Stream data into a Kinesis data stream",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var writer core.Writer
		writer = kinesis.NewWriter(utils.NewSDKSession(), stream)

		transformer := &kinesis.Transformer{}
		deserializer := &core.JsonDeserializer{}
		config := &core.Config{File: "", BatchSize: batchSize, MaxWriters: maxWriters, AutoFlushTimeoutSeconds: flushTimeoutSeconds}
		if config.BatchSize > 500 {
			config.BatchSize = 500
		}
		core.Start(deserializer, transformer, writer, config)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(kinesisCmd)

	kinesisCmd.Flags().StringVar(&stream, "stream", "", "name of stream")
	kinesisCmd.MarkFlagRequired("stream")
}
