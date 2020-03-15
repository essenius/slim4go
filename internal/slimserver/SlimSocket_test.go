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
	"net"
	"testing"

	"github.com/essenius/slim4go/internal/assert"
)

func TestSlimSocketUsedPort(t *testing.T) {
	var (
		logBuffer bytes.Buffer
		logger    = log.New(&logBuffer, "logger: ", log.Lshortfile)
		socket    = newSlimSocket(-1, logger, 3e10)
	)
	err := socket.Listen()
	assert.Equals(t, "listen tcp: address -1: invalid port", err.Error(), "Error logged")
}

func TestSlimSocketTimeout(t *testing.T) {
	var (
		logBuffer bytes.Buffer
		logger    = log.New(&logBuffer, "logger: ", log.Lshortfile)
		socket    = newSlimSocket(8485, logger, 1e6)
	)
	err := socket.Listen()
	assert.Equals(t, "Timeout (1ms) waiting for a connection", err.Error(), "Error logged")
}

func TestSlimSocketOk(t *testing.T) {
	var (
		logBuffer   bytes.Buffer
		logger      = log.New(&logBuffer, "logger: ", log.Lshortfile)
		socket      = newSlimSocket(8485, logger, 3e10)
		readBuffer  = make([]byte, 25)
		writeBuffer = make([]byte, 25)
	)

	go func() {
		socket.Listen()
		count, err1 := socket.Read(readBuffer)
		assert.Equals(t, nil, err1, "no error in goroutine")
		assert.Equals(t, 20, count, "bytres read")
		assert.Equals(t, "some text to be sent", string(readBuffer)[:count], "content read")
		socket.SendMessage("message")
	}()

	conn, err2 := net.Dial("tcp", ":8485")
	if err2 != nil {
		t.Error("could not connect: ", err2)
	}
	defer conn.Close()

	payload := []byte("some text to be sent")
	conn.Write(payload)

	count, err3 := conn.Read(writeBuffer)
	assert.Equals(t, nil, err3, "no error3")
	assert.Equals(t, 7, count, "bytes sent")
	assert.Equals(t, "message", string(writeBuffer)[:count], "content sent")
}
