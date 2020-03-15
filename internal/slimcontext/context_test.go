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

package slimcontext

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/essenius/slim4go/internal/assert"
)

func TestContextNew(t *testing.T) {
	noCallback := func(err error) {
		t.Fatalf("Unexpected callback with error '%v'", err.Error())
	}

	callbackCount := 0
	callback1 := func(err error) {
		assert.Equals(t, "invalid value \"q\" for flag -s: parse error", err.Error(), "Error message in callback1")
		callbackCount++
	}
	callback2 := func(err error) {
		assert.Equals(t, "port 'a' should be numerical", err.Error(), "Error message in callback2")
		callbackCount++
	}
	args := []string{"slim4go", "-s", "7", "-t", "20", "1"}

	contextOk := newContext()
	contextOk.ErrorAction = noCallback
	contextOk.Initialize(args)
	assert.Equals(t, 1, contextOk.Port, "port == 1")
	assert.Equals(t, time.Duration(7)*time.Second, contextOk.InstructionTimeout, "instruction timeout is 7 seconds")
	assert.Equals(t, time.Duration(20)*time.Second, contextOk.ConnectionTimeout, "connection timeout is 20 seconds")
	args[3] = "8475"
	contextOk.Initialize(args)
	assert.Equals(t, 1, contextOk.Port, "Initialize not executed a second time")

	args[2] = "q"
	contextErr1 := newContext()
	contextErr1.ErrorAction = callback1
	contextErr1.Initialize(args)
	assert.Equals(t, 8475, contextErr1.Port, "port ok")
	assert.Equals(t, 1, callbackCount, "callback called once")

	args[2] = "2"
	callbackCount = 0
	args[3] = "a"
	contextErr2 := newContext()
	contextErr2.ErrorAction = callback2
	contextErr2.Initialize(args)
	assert.Equals(t, 1, contextErr2.Port, "port defaulted to 1 with invalid port")
	assert.Equals(t, 1, callbackCount, "callback called once")

}

func TestContextParsePort(t *testing.T) {
	args := []string{}
	port1, err1 := parsePort(args)
	assert.Equals(t, 1, port1, "port = 1 with empty args")
	assert.Equals(t, "Missing port specification. Expected params [-s timeout] port", err1.Error(), "no arguments")
	args = []string{"a"}
	port2, err2 := parsePort(args)
	assert.Equals(t, 1, port2, "port = 1 with non-numerical port")
	assert.Equals(t, "port 'a' should be numerical", err2.Error(), "err != nil with port = a")
	args[0] = "-5"
	port3, err3 := parsePort(args)
	assert.Equals(t, 1, port3, "port = -5")
	assert.Equals(t, "port '-5' should be non-negative", err3.Error(), "err = nil with port = -5")
	args[0] = "8475"
	port4, err4 := parsePort(args)
	assert.Equals(t, 8475, port4, "port = 8475")
	assert.Equals(t, nil, err4, "err = nil with port = 8475")
	args[0] = "1"
	port5, err5 := parsePort(args)
	assert.Equals(t, 1, port5, "port = 1 ")
	assert.Equals(t, nil, err5, "err == nil with port = 1")
}

func TestContextInject(t *testing.T) {
	fmt.Printf("%v", os.Args)
	context := InjectContext()
	assert.Equals(t, time.Duration(0), context.InstructionTimeout, "Timeout not initialized yet")
	assert.Equals(t, time.Duration(0), context.ConnectionTimeout, "Timeout not initialized yet")
}
