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
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DebugWriter struct {
}

func (w *DebugWriter) Write(in []interface{}) error {
	for _, i := range in {
		m := i.(map[string]*dynamodb.AttributeValue)
		for k, v := range m {
			_, err := os.Stdout.Write([]byte(fmt.Sprintf("%s %v\n", k, v.GoString())))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewDebugWriter() *DebugWriter {
	return &DebugWriter{}
}
