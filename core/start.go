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
	"context"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

func Start(file string, deserializer Deserializer, transformer Transformer, writer Writer, batchSize, maxWriters, timeoutSeconds int) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	chanSignal := make(chan os.Signal, 1)
	c := make(chan interface{})
	bw := NewBatchWriter(c, batchSize, maxWriters, time.After(time.Second*time.Duration(timeoutSeconds)), writer)
	p := NewPump(deserializer, transformer, c)

	signal.Notify(chanSignal, os.Interrupt)

	reader := os.Stdin
	if file != "" {
		var err error
		reader, err = os.Open(file)
		if err != nil {
			panic(err)
		}
	}

	go bw.Start(ctx)
	go p.Start(ctx, reader)

	select {
	case e := <-bw.Done():
		if e != nil {
			log.Error(e)
		}
	case e := <-p.Done():
		if e != nil {
			log.Error(e)
		}
	case <-chanSignal:
	}

	cancelFunc()

	for e := range bw.Done() {
		log.Error(e)
	}

	for e := range p.Done() {
		log.Error(e)
	}
}
