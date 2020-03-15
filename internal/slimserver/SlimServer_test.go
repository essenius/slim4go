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
	"reflect"
	"strings"
	"testing"

	"github.com/essenius/slim4go/internal/demofixtures"

	"github.com/essenius/slim4go/internal/assert"
	"github.com/essenius/slim4go/internal/slimcontext"
	"github.com/essenius/slim4go/internal/slimlog"
	"github.com/essenius/slim4go/internal/slimprocessor"
)

type testMessenger struct {
	reader      *strings.Reader
	readIndex   int
	writeIndex  int
	input       []string
	output      []string
	t           *testing.T
	description string
}

func newTestMessenger(t *testing.T, input []string, output []string, description string) *testMessenger {
	messenger := new(testMessenger)
	messenger.input = input
	messenger.output = output
	messenger.t = t
	messenger.description = description
	return messenger
}

// Read receives a number of bytes from Stdin
func (messenger *testMessenger) Read(buffer []byte) (int, error) {
	if len(messenger.input) < messenger.readIndex {
		assert.IsTrue(messenger.t, false, "EOF")
	}
	messenger.reader = strings.NewReader(messenger.input[messenger.readIndex])
	messenger.readIndex++
	readBytes, err := io.ReadAtLeast(messenger.reader, buffer, 1)
	fmt.Printf("%v: Read is returning %v bytes", messenger.description, readBytes)
	return readBytes, err
}

// Listen starts listening on Stdin
func (messenger *testMessenger) Listen() error {
	messenger.readIndex = 0
	messenger.writeIndex = 0
	if strings.HasPrefix(messenger.description, "ListenError") {
		return fmt.Errorf(messenger.description)
	}
	return nil
}

// SendMessage sends a message on Stdout
func (messenger *testMessenger) SendMessage(message string) error {
	slimlog.Trace.Println(message)
	assert.Equals(messenger.t, messenger.output[messenger.writeIndex], message,
		fmt.Sprintf("%v line %v", messenger.description, messenger.writeIndex))
	messenger.writeIndex++
	if strings.HasPrefix(messenger.description, "SendError") {
		return fmt.Errorf(messenger.description)
	}
	return nil
}

func TestSlimServerInject(t *testing.T) {
	context := slimcontext.InjectContext()
	context.Initialize([]string{"app", "1"})
	server := InjectSlimServer()
	assert.IsTrue(t, server.processor != nil, "processor created")
	assert.Equals(t, reflect.TypeOf(new(slimPipe)), reflect.TypeOf(server.messenger), "messenger is slimPipe")

	context.Initialize([]string{"app", "8495"})
	server2 := InjectSlimServer()
	assert.Equals(t, *server, *server2, "second call returns the same object")
	assert.Equals(t, reflect.TypeOf(new(slimPipe)), reflect.TypeOf(server2.messenger), "second initialization ignored")
}

func TestServerListenError(t *testing.T) {
	args := []string{"slim4go", "-s", "1", "2"}
	context := slimcontext.InjectContext()
	context.ErrorAction = func(err error) {
		t.Fatal("Unexpected call to ErrorAction")
	}
	context.Initialize(args)
	context.Port = -1
	messenger1 := newTestMessenger(t, []string{}, []string{}, "ListenError")
	slimServer1 := newSlimServer(slimprocessor.InjectFixtures(), messenger1, slimprocessor.InjectSlimProcessor())
	err1 := slimServer1.Serve()
	assert.Equals(t, "ListenError", err1.Error(), "Error listening")

	messenger2 := newTestMessenger(t, []string{"a"}, []string{"Slim -- V0.5\n"}, "SendError")
	slimServer2 := newSlimServer(slimprocessor.InjectFixtures(), messenger2, slimprocessor.InjectSlimProcessor())
	err2 := slimServer2.Serve()
	assert.Equals(t, "SendError", err2.Error(), "Error sending")
}

func TestServerServe(t *testing.T) {
	testInput := []string{
		"000459:[000004:" +
			"000096:[000004:000015:scriptTable_0_0:000004:make:000016:scriptTableActor:000020:TemperatureConverter:]:" +
			"000127:[000007:000015:scriptTable_0_1:000013:callAndAssign:000004:temp:000016:scriptTableActor:000009:ConvertTo:000004:68 F:000001:C:]:" +
			"000093:[000005:000015:scriptTable_0_2:000004:call:000016:scriptTableActor:000004:echo:000005:$temp:]:" +
			"000102:[000006:000015:scriptTable_0_3:000004:call:000016:scriptTableActor:000009:ConvertTo:000000::000001:K:]:]",
		"000003:bye",
	}

	expectedOutput := []string{
		"Slim -- V0.5\n",
		"000287:[000004:" +
			"000042:[000002:000015:scriptTable_0_0:000002:OK:]:" +
			"000042:[000002:000015:scriptTable_0_1:000002:20:]:" +
			"000042:[000002:000015:scriptTable_0_2:000002:20:]:" +
			"000120:[000002:000015:scriptTable_0_3:000080:__EXCEPTION__:message:<<Panic: Expected float with suffix F, C or K but got ''>>:]:]",
	}

	args := []string{"slim4go", "-s", "1", "1"}

	context := slimcontext.InjectContext()
	context.ErrorAction = func(err error) {
		t.Fatal("Unexpected call to ErrorAction")
	}
	context.Initialize(args)
	messenger1 := newTestMessenger(t, testInput, expectedOutput, "Test 1")
	slimServer1 := newSlimServer(slimprocessor.InjectFixtures(), messenger1, slimprocessor.InjectSlimProcessor())
	slimServer1.RegisterFixture("TemperatureConverter", demofixtures.NewTemperatureConverter)
	slimServer1.Serve()

	messenger2 := newTestMessenger(t, []string{"000005:bogus"}, []string{"Slim -- V0.5\n"}, "Test 2 - bogus message")
	slimServer2 := newSlimServer(slimprocessor.InjectFixtures(), messenger2, slimprocessor.InjectSlimProcessor())
	slimServer2.Serve()

	messenger3 := newTestMessenger(t, []string{"000005:bye"}, []string{"Slim -- V0.5\n"}, "Test 3 - size wrong")
	slimServer3 := newSlimServer(slimprocessor.InjectFixtures(), messenger3, slimprocessor.InjectSlimProcessor())
	slimServer3.Serve()

}
