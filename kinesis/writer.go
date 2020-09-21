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

package kinesis

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

type Writer struct {
	stream        string
	kinesisClient *kinesis.Kinesis
}

func (w *Writer) Write(items []interface{}) error {
	records := make([]*kinesis.PutRecordsRequestEntry, 0)

	for _, i := range items {
		m := i.(map[string]interface{})
		d, err := json.Marshal(m["Data"])
		if err != nil {
			return err
		}
		e := &kinesis.PutRecordsRequestEntry{
			Data:         d,
			PartitionKey: aws.String(m["PartitionKey"].(string)),
		}
		records = append(records, e)
	}

	input := &kinesis.PutRecordsInput{
		StreamName: &w.stream,
		Records:    records,
	}

	out, err := w.kinesisClient.PutRecords(input)
	if err != nil {
		return err
	}

	for *out.FailedRecordCount > 0 {
		retry := make([]*kinesis.PutRecordsRequestEntry, 0)
		for i, f := range out.Records {
			if f.ErrorCode != nil {
				retry = append(retry, records[i])
			}
		}
		input.Records = retry
		records = retry
		out, err = w.kinesisClient.PutRecords(input)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewWriter(s *session.Session, stream string) *Writer {
	kc := kinesis.New(s)
	return &Writer{
		stream:        stream,
		kinesisClient: kc,
	}
}
