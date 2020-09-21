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

import (
	"bufio"
	"context"
	"io"
)

type Pump struct {
	deserializer Deserializer
	transformer  Transformer
	output       chan<- interface{}
	done         chan error
}

func (p *Pump) Start(ctx context.Context, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for {
		select {
		case <-ctx.Done():
			p.exit(nil)
			return
		default:
			if !scanner.Scan() {
				p.exit(nil)
				return
			}
			input := scanner.Text()
			data, err := p.deserializer.Deserialize(input)
			if err != nil {
				data = input
			}

			transformed, err := p.transformer.Transform(data)
			if err != nil {
				p.exit(err)
				return
			}

			select {
			case <-ctx.Done():
				p.exit(nil)
				return
			case p.output <- transformed:
			}
		}
	}
}

func (p *Pump) Done() <-chan error {
	return p.done
}

func (p *Pump) exit(err error) {
	if err != nil {
		p.done <- err
	}
	close(p.done)
}

func NewPump(deserializer Deserializer, transformer Transformer, output chan<- interface{}) *Pump {
	return &Pump{
		deserializer: deserializer,
		transformer:  transformer,
		output:       output,
		done:         make(chan error),
	}
}
