// Copyright 2020 Rik Essenius
//
//   Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//   except in compliance with the License. You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software distributed under the License
//   is distributed on an "AS IS" BASIS WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and limitations under the License.

package slimserver

import (
	"fmt"
	"io"
	"log"
	"time"
)

type slimPipe struct {
	reader     io.Reader
	writer     io.Writer
	infoLogger *log.Logger
	timeout    time.Duration
}

func newSlimPipe(reader io.Reader, writer io.Writer, infoLogger *log.Logger, timeout time.Duration) *slimPipe {
	pipe := new(slimPipe)
	pipe.infoLogger = infoLogger
	pipe.reader = reader
	pipe.writer = writer
	pipe.timeout = timeout
	return pipe
}

// Listen starts listening on Stdin
func (pipe *slimPipe) Listen() error {
	pipe.infoLogger.Println("Listening on Stdin")
	return nil
}

// SendMessage sends a message on Stdout
func (pipe *slimPipe) SendMessage(message string) error {
	_, err := io.WriteString(pipe.writer, message)
	return err
}

// Read receives a number of bytes from Stdin
func (pipe *slimPipe) Read(buffer []byte) (int, error) {
	errChannel := make(chan error)
	countChannel := make(chan int)
	go func() {
		readBytes, err := io.ReadAtLeast(pipe.reader, buffer, 1)
		if err != nil {
			errChannel <- err
		} else {
			countChannel <- readBytes
		}
		close(errChannel)
		close(countChannel)
	}()
	select {
	case count := <-countChannel:
		return count, nil
	case err := <-errChannel:
		return 0, err
	case <-time.After(pipe.timeout):
		return 0, fmt.Errorf("Timeout (%v)", pipe.timeout)
	}
}
