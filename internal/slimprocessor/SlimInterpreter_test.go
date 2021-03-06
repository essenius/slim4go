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

package slimprocessor

import (
	"fmt"
	"testing"
	"time"

	"github.com/essenius/slim4go/internal/assert"
	"github.com/essenius/slim4go/internal/interfaces"
	"github.com/essenius/slim4go/internal/slimentity"
)

type MockStatementProcessor struct {
	SetSymbolCalls       int
	RegisterFixtureCalls int
}

func (mock *MockStatementProcessor) FixtureRegistry() interfaces.Registry {
	return nil
}

func (mock *MockStatementProcessor) DoCall(instanceName, methodName string, args *slimentity.SlimList) slimentity.SlimEntity {
	return fmt.Sprintf("Call %v %v(%v)", instanceName, methodName, args.ToString())
}

func (mock *MockStatementProcessor) DoImport(path string) slimentity.SlimEntity {
	if path == "wait" {
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
	return fmt.Sprintf("Import %v", path)
}

func (mock *MockStatementProcessor) DoMake(instanceName, fixtureName string, args *slimentity.SlimList) slimentity.SlimEntity {
	return fmt.Sprintf("Make %v %v(%v)", instanceName, fixtureName, args.ToString())
}

func (mock *MockStatementProcessor) Objects() interfaces.Collector {
	return nil
}

func (mock *MockStatementProcessor) Parser() interfaces.Parser {
	return nil
}

func (mock *MockStatementProcessor) SerializeObjectsIn(input slimentity.SlimEntity) slimentity.SlimEntity {
	return input
}

func (mock *MockStatementProcessor) SetSymbol(symbol string, value interface{}) {
	mock.SetSymbolCalls++
}

func MakeInstructionList(instruction ...slimentity.SlimEntity) *slimentity.SlimList {
	instructionList := slimentity.NewSlimListContaining(instruction)
	command := slimentity.NewSlimList()
	command.Append(instructionList)
	return command
}

func TestSlimInterpreterExecute1(t *testing.T) {
	MockStatementProcessor := new(MockStatementProcessor)
	slimInterpreter := NewSlimInterpreter(MockStatementProcessor, time.Duration(10)*time.Second)
	importList := MakeInstructionList("import1", "import", "test")
	assert.Equals(t, `[[import1, Import test]]`, slimInterpreter.Process(importList).ToString(), "Import")
	makeList := MakeInstructionList("make1", "make", "instance1", "fixture", "arg1", "arg2")
	assert.Equals(t, `[[make1, Make instance1 fixture([arg1, arg2])]]`, slimInterpreter.Process(makeList).ToString(), "Make")
	callList := MakeInstructionList("call1", "call", "instance1", "method1", "arg1")
	assert.Equals(t, `[[call1, Call instance1 method1([arg1])]]`, slimInterpreter.Process(callList).ToString(), "Call")
	assert.Equals(t, 0, MockStatementProcessor.SetSymbolCalls, "SetSymbol not called")
	callAndAssignList := MakeInstructionList("callAndAssign1", "callAndAssign", "symbol1", "instance1", "method2")
	assert.Equals(t, `[[callAndAssign1, Call instance1 method2([])]]`, slimInterpreter.Process(callAndAssignList).ToString(), "CallAndAssign")
	assert.Equals(t, 1, MockStatementProcessor.SetSymbolCalls, "SetSymbol called once")
	assignList := MakeInstructionList("assign1", "assign", "symbol2", "value2")
	assert.Equals(t, `[[assign1, OK]]`, slimInterpreter.Process(assignList).ToString(), "Assign")
	assert.Equals(t, 2, MockStatementProcessor.SetSymbolCalls, "SetSymbol called twice")
}

func TestSlimInterpreterTimeout(t *testing.T) {
	MockStatementProcessor := new(MockStatementProcessor)
	slimInterpreter := NewSlimInterpreter(MockStatementProcessor, time.Duration(1)*time.Nanosecond)
	importList := MakeInstructionList("import1", "import", "wait")
	assert.Equals(t, `[[import1, __EXCEPTION__:message:<<TIMED_OUT 0>>]]`, slimInterpreter.Process(importList).ToString(), "Import with timeout")
}

func TestSlimInterpreterMalformedInstructions(t *testing.T) {
	MockStatementProcessor := new(MockStatementProcessor)
	slimInterpreter := NewSlimInterpreter(MockStatementProcessor, time.Duration(7)*time.Second)
	importList := MakeInstructionList("import1", "import")
	assert.Equals(t, `[[import1, __EXCEPTION__:message:<<MALFORMED_INSTRUCTION [import1, import]>>]]`, slimInterpreter.Process(importList).ToString(), "Import invalid")
	makeList := MakeInstructionList("make1", "make", "instance1")
	assert.Equals(t, `[[make1, __EXCEPTION__:message:<<MALFORMED_INSTRUCTION [make1, make, instance1]>>]]`, slimInterpreter.Process(makeList).ToString(), "Make invalid")
	callList := MakeInstructionList("call1", "call", "instance1")
	assert.Equals(t, `[[call1, __EXCEPTION__:message:<<MALFORMED_INSTRUCTION [call1, call, instance1]>>]]`, slimInterpreter.Process(callList).ToString(), "Call invalid")
	callAndAssignList := MakeInstructionList("callAndAssign1", "callAndAssign", "symbol1", "instance1")
	assert.Equals(t, `[[callAndAssign1, __EXCEPTION__:message:<<MALFORMED_INSTRUCTION [callAndAssign1, callAndAssign, symbol1, instance1]>>]]`, slimInterpreter.Process(callAndAssignList).ToString(), "CallAndAssign invalid")
	assignList := MakeInstructionList("assign1", "assign")
	assert.Equals(t, `[[assign1, __EXCEPTION__:message:<<MALFORMED_INSTRUCTION [assign1, assign]>>]]`, slimInterpreter.Process(assignList).ToString(), "Assign invalid")
	nullList := MakeInstructionList()
	assert.Equals(t, `[[__EXCEPTION__:message:<<MALFORMED_INSTRUCTION []>>]]`, slimInterpreter.Process(nullList).ToString(), "Null")
	unknownCommandList := MakeInstructionList("unknown1", "unknown")
	assert.Equals(t, `[[unknown1, __EXCEPTION__:message:<<MALFORMED_INSTRUCTION [unknown1, unknown]>>]]`, slimInterpreter.Process(unknownCommandList).ToString(), "unknown Command")
	noCommandList := MakeInstructionList("bogus")
	assert.Equals(t, `[[bogus, __EXCEPTION__:message:<<MALFORMED_INSTRUCTION [bogus]>>]]`, slimInterpreter.Process(noCommandList).ToString(), "no command")
}

