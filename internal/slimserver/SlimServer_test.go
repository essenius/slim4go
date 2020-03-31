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
	"strings"
	"testing"
	"time"

	"github.com/essenius/slim4go/examples/demofixtures"

	"github.com/essenius/slim4go/internal/assert"
	"github.com/essenius/slim4go/internal/fixture"
	"github.com/essenius/slim4go/internal/slimprocessor"
	"github.com/essenius/slim4go/internal/standardlibrary"
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

// Read receives a number of bytes from a string slice
func (messenger *testMessenger) Read(buffer []byte) (int, error) {
	if len(messenger.input) < messenger.readIndex+1 {
		return 0, fmt.Errorf("EOF")
	}
	messenger.reader = strings.NewReader(messenger.input[messenger.readIndex])
	messenger.readIndex++
	readBytes, err := io.ReadAtLeast(messenger.reader, buffer, 1)
	return readBytes, err
}

// Listen emulates listening. It returns an error if description prefix is ListenError.
func (messenger *testMessenger) Listen() error {
	messenger.readIndex = 0
	messenger.writeIndex = 0
	if strings.HasPrefix(messenger.description, "ListenError") {
		return fmt.Errorf(messenger.description)
	}
	return nil
}

// SendMessage emulates sending a message. It returns an error if desciption prefix is SendError.
func (messenger *testMessenger) SendMessage(message string) error {
	assert.Equals(messenger.t, messenger.output[messenger.writeIndex], message,
		fmt.Sprintf("%v line %v", messenger.description, messenger.writeIndex))
	messenger.writeIndex++
	if strings.HasPrefix(messenger.description, "SendError") {
		return fmt.Errorf(messenger.description)
	}
	return nil
}

func TestServeErrorResponses(t *testing.T) {
	messenger1 := newTestMessenger(t, []string{}, []string{}, "ListenError")
	slimServer1 := NewSlimServer(nil, messenger1, nil)
	err1 := slimServer1.Serve()
	assert.Equals(t, "ListenError", err1.Error(), "Error listening")

	messenger2 := newTestMessenger(t, []string{"a"}, []string{"Slim -- V0.5\n"}, "SendError")
	slimServer2 := NewSlimServer(nil, messenger2, nil)
	err2 := slimServer2.Serve()
	assert.Equals(t, "SendError", err2.Error(), "Error sending")

	messenger3 := newTestMessenger(t, []string{"000005:bogus"}, []string{"Slim -- V0.5\n"}, "Test 2 - bogus message")
	slimServer3 := NewSlimServer(nil, messenger3, nil)
	err3 := slimServer3.Serve()
	assert.Equals(t, "Encountered unexpected command 'bogus'", err3.Error(), "Error sending")

	messenger4 := newTestMessenger(t, []string{"000005:bye"}, []string{"Slim -- V0.5\n"}, "Test 3 - size wrong")
	slimServer4 := NewSlimServer(nil, messenger4, nil)
	err4 := slimServer4.Serve()
	assert.Equals(t, "readExactBytes: Expected 5 bytes from Slim client, but got 3", err4.Error(), "Error sending")
}

func TestServerServe(t *testing.T) {
	testInput := []string{
		"000472:[000004:" +
			"000109:[000004:000015:scriptTable_0_0:000004:make:000016:scriptTableActor:000033:demofixtures.TemperatureConverter:]:" +
			"000127:[000007:000015:scriptTable_0_1:000013:callAndAssign:000004:temp:000016:scriptTableActor:000009:ConvertTo:000004:68 F:000001:C:]:" +
			"000093:[000005:000015:scriptTable_0_2:000004:call:000016:scriptTableActor:000004:echo:000005:$temp:]:" +
			"000102:[000006:000015:scriptTable_0_3:000004:call:000016:scriptTableActor:000009:convertTo:000000::000001:K:]:]",
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

	// Spell out what the injector is expected to do.
	registry := fixture.NewRegistry()
	symbols := slimprocessor.NewSymbolTable()
	parser := slimprocessor.NewParser(symbols)
	objectHandler := slimprocessor.NewObjectHandler(parser)
	standardLibrary := standardlibrary.New(standardlibrary.NewActorStack(), objectHandler)
	objectHandler.Add("libraryStandard", standardLibrary)
	parser.SetObjectSerializer(objectHandler)
	processor := slimprocessor.NewStatementProcessor(registry, objectHandler, parser, symbols)
	interpreter := slimprocessor.NewSlimInterpreter(processor, time.Second)

	messenger1 := newTestMessenger(t, testInput, expectedOutput, "Test 1")
	slimServer1 := NewSlimServer(registry, messenger1, interpreter)
	slimServer1.RegisterFixturesFrom(demofixtures.NewTemperatureFactory())
	slimServer1.Serve()
	assert.Equals(t, "Could not add fixture '1'", slimServer1.RegisterFixture(1).Error(), "Wrong argument for RegisterFixture returns an error")
}
