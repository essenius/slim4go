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
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/essenius/slim4go/internal/assert"
)

func TestSlimPipeOk(t *testing.T) {
	var (
		logBuffer   bytes.Buffer
		logger      = log.New(&logBuffer, "logger: ", log.Lshortfile)
		readBuffer  = strings.NewReader("some text to be read")
		writeBuffer bytes.Buffer
		pipe        = newSlimPipe(readBuffer, &writeBuffer, logger, 3e10)
		pipeBuffer  = make([]byte, 25)
	)

	pipe.Listen()
	count, err := pipe.Read(pipeBuffer)
	assert.Equals(t, nil, err, "no error")
	assert.Equals(t, 20, count, "count")
	assert.Equals(t, "some text to be read", string(pipeBuffer)[:count], "content")
	pipe.SendMessage("message")
	assert.Equals(t, "message", writeBuffer.String(), "Write")

	readBuffer = strings.NewReader("")
	_, err1 := pipe.Read(pipeBuffer)
	assert.Equals(t, "EOF", err1.Error(), "Error")

}

func TestSlimPipeTimeout(t *testing.T) {
	var (
		logBuffer   bytes.Buffer
		logger      = log.New(&logBuffer, "logger: ", log.Lshortfile)
		writeBuffer bytes.Buffer
		pipe        = newSlimPipe(new(NeverEndingReader), &writeBuffer, logger, 1e6)
		pipeBuffer  = make([]byte, 25)
	)

	pipe.Listen()
	_, err := pipe.Read(pipeBuffer)
	assert.IsTrue(t, err != nil, "Error occurred")
	assert.Equals(t, "Timeout (1ms)", err.Error(), "Timeout")
}

type NeverEndingReader struct{}

func (reader *NeverEndingReader) Read(p []byte) (n int, err error) {
	for {
	}
}
