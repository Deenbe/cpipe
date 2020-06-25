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

package dynamodb

import (
	"net"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Writer struct {
	table      string
	ddbClient  *dynamodb.DynamoDB
	retryDelay time.Duration
}

func (w *Writer) Write(in []interface{}) error {
	requests := make([]*dynamodb.WriteRequest, 0, len(in))
	for _, i := range in {
		requests = append(requests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: i.(map[string]*dynamodb.AttributeValue),
			},
		})
	}
	return w.writeAll(map[string][]*dynamodb.WriteRequest{
		w.table: requests,
	})
}

func (w *Writer) writeAll(requests map[string][]*dynamodb.WriteRequest) error {
	for len(requests) > 0 {
		output, err := w.ddbClient.BatchWriteItem(&dynamodb.BatchWriteItemInput{
			RequestItems: requests,
		})

		if err != nil {
			return err
		}
		requests = output.UnprocessedItems
		if len(requests) > 0 {
			time.Sleep(w.retryDelay)
		}
	}

	return nil
}

func optimisedHTTPClient() *http.Client {
	// TODO: Make this configurable
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport: t,
	}
}

func NewWriter(table string, retryDelay time.Duration) *Writer {
	s := session.Must(session.NewSession(&aws.Config{
		HTTPClient: optimisedHTTPClient(),
	}))
	c := dynamodb.New(s)
	return &Writer{
		table:      table,
		ddbClient:  c,
		retryDelay: retryDelay,
	}
}
