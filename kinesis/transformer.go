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
	"github.com/google/uuid"
)

type Transformer struct {
	partitionKey string
}

func (t *Transformer) Transform(in interface{}) (interface{}, error) {
	r := make(map[string]interface{})
	r["PartitionKey"] = nil
	r["Data"] = in

	// If incoming data is in json format, try to extract the
	// partition key from the content.
	if m, ok := in.(map[string]interface{}); ok {
		if k, ok := m[t.partitionKey]; ok {
			r["PartitionKey"] = k
		}
	}

	// If we could not resolve the partition key in input,
	// assign an uuid.
	if r["PartitionKey"] == nil {
		r["PartitionKey"] = uuid.New().String()
	}

	return r, nil
}
