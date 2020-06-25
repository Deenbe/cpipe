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

package core

// Writer is the interface between a data stream and a destination service.
type Writer interface {
	Write([]interface{}) error
}

// Deserializer converts input data to a map with typed information.
type Deserializer interface {
	Deserialize(string) (interface{}, error)
}

// Transformer transforms an input data into a specific structure required by
// writer.
type Transformer interface {
	Transform(interface{}) (interface{}, error)
}

// Config used to pass core configuration parameters around.
type Config struct {
	File                    string
	MaxWriters              int
	BatchSize               int
	AutoFlushTimeoutSeconds int
}
