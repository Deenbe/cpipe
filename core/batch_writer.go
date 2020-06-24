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
	"time"
)

type BatchWriter struct {
	input      <-chan interface{}
	batchSize  int
	maxWriters int
	timer      <-chan time.Time
	done       chan error
	writer     Writer
}

func (w *BatchWriter) Start(ctx context.Context) {
	nextBatch := make([]interface{}, 0, w.batchSize)
	writersCount := 0
	writerCompletion := make(chan error)

	for {
		if w.maxWriters == writersCount {
			// We have exhausted the max number of writers.
			// Wait for one of them to successfully complete
			// or the context is cancelled.
			select {
			case <-ctx.Done():
				w.waitForWriters(writersCount, writerCompletion)
				return
			case e := <-writerCompletion:
				writersCount--
				if e != nil {
					w.done <- e
					w.waitForWriters(writersCount, writerCompletion)
					return
				}
			}
		}

		select {
		case <-ctx.Done():
			w.waitForWriters(writersCount, writerCompletion)
			return
		case e := <-writerCompletion:
			writersCount--
			if e != nil {
				w.done <- e
				w.waitForWriters(writersCount, writerCompletion)
				return
			}
		case i := <-w.input:
			nextBatch = append(nextBatch, i)
			if len(nextBatch) == w.batchSize {
				writersCount++
				w.writeBatch(nextBatch, writerCompletion)
				nextBatch = make([]interface{}, 0, w.batchSize)
			}
		case <-w.timer:
			if len(nextBatch) > 0 {
				w.writeBatch(nextBatch, writerCompletion)
				nextBatch = make([]interface{}, w.batchSize)
			}
		}
	}
}

func (w *BatchWriter) Done() <-chan error {
	return w.done
}

func (w *BatchWriter) waitForWriters(count int, wc <-chan error) {
	for count != 0 {
		e := <-wc
		if e != nil {
			w.done <- e
		}
		count--
	}
	close(w.done)
}

func (w *BatchWriter) writeBatch(b []interface{}, writerCompletion chan<- error) {
	go func() {
		e := w.writer.Write(b)
		writerCompletion <- e
	}()
}

func NewBatchWriter(input <-chan interface{}, batchSize, maxWriters int, timer <-chan time.Time, writer Writer) *BatchWriter {
	return &BatchWriter{
		input:      input,
		batchSize:  batchSize,
		maxWriters: maxWriters,
		timer:      timer,
		done:       make(chan error, batchSize),
		writer:     writer,
	}
}
