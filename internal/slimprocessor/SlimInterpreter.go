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
	"time"

	"github.com/essenius/slim4go/internal/interfaces"

	"github.com/essenius/slim4go/internal/slimentity"
	"github.com/essenius/slim4go/internal/slimprotocol"
)

// Definitions and constructors

// SlimInterpreter processes incomming SlimLists and dispatches commands to its statement processor.
type SlimInterpreter struct {
	processor interfaces.StatementProcessor
	timeout   time.Duration
}

// NewSlimInterpreter creates a new Slim interpreter.
func NewSlimInterpreter(processor interfaces.StatementProcessor, timeout time.Duration) *SlimInterpreter {
	slimInterpreter := new(SlimInterpreter)
	slimInterpreter.processor = processor
	slimInterpreter.timeout = timeout
	return slimInterpreter
}

// Helper methods

func addResult(list *slimentity.SlimList, entry ...slimentity.SlimEntity) {
	sublist := slimentity.NewSlimListContaining(entry)
	list.Append(sublist)
}

func malformedInstruction(instruction interface{}) string {
	if slimentity.IsSlimList(instruction) {
		return slimprotocol.MalformedInstruction(instruction.(*slimentity.SlimList).ToString())
	}
	return slimprotocol.MalformedInstruction(instruction.(string))
}

// Methods

func (slimInterpreter *SlimInterpreter) dispatch(instruction *slimentity.SlimList) slimentity.SlimEntity {
	command := instruction.StringAt(1)
	switch command {
	case "assign":
		return slimInterpreter.doAssign(instruction)
	case "call":
		return slimInterpreter.doCall(instruction, noAssign)
	case "callAndAssign":
		return slimInterpreter.doCall(instruction, assign)
	case "import":
		return slimInterpreter.DoImport(instruction)
	case "make":
		return slimInterpreter.DoMake(instruction)
	default:
		return malformedInstruction(instruction)
	}
}

func (slimInterpreter *SlimInterpreter) dispatchWithTimeout(instruction *slimentity.SlimList) slimentity.SlimEntity {
	resultChannel := make(chan slimentity.SlimEntity, 1)
	instructionTimer := time.NewTimer(slimInterpreter.timeout)
	go func() {
		returnValue := slimInterpreter.dispatch(instruction)
		resultChannel <- returnValue
	}()
	select {
	case result := <-resultChannel:
		instructionTimer.Stop()
		return result
	case <-instructionTimer.C:
		return slimprotocol.TimedOut(slimInterpreter.timeout)
	}
}

func (slimInterpreter *SlimInterpreter) doAssign(instruction *slimentity.SlimList) string {
	if instruction.Length() < 3 {
		return malformedInstruction(instruction)
	}
	symbolName := instruction.StringAt(2)
	value := instruction.StringAt(3)
	slimInterpreter.processor.SetSymbol(symbolName, value)
	return slimprotocol.OK()
}

const (
	noAssign = 4
	assign   = 5
)

func (slimInterpreter *SlimInterpreter) doCall(instruction *slimentity.SlimList, minLength int) slimentity.SlimEntity {
	if instruction.Length() < minLength {
		return malformedInstruction(instruction)
	}
	startIndex := 2
	var symbolName string
	if minLength == assign {
		symbolName = instruction.StringAt(startIndex)
		startIndex++
	}
	instanceName := instruction.StringAt(startIndex)
	methodName := instruction.StringAt(startIndex + 1)
	args := instruction.TailAt(startIndex + 2)
	result := slimInterpreter.processor.DoCall(instanceName, methodName, args)
	if minLength == assign {
		slimInterpreter.processor.SetSymbol(symbolName, result)
	}
	return slimInterpreter.processor.SerializeObjectsIn(result)
}

// DoImport executes an Import instruction.
func (slimInterpreter *SlimInterpreter) DoImport(instruction *slimentity.SlimList) slimentity.SlimEntity {
	if instruction.Length() < 3 {
		return malformedInstruction(instruction)
	}
	pathName := instruction.StringAt(2)
	return slimInterpreter.processor.DoImport(pathName)
}

// DoMake executes a Make instruction.
func (slimInterpreter *SlimInterpreter) DoMake(instruction *slimentity.SlimList) slimentity.SlimEntity {
	if instruction.Length() < 4 {
		return malformedInstruction(instruction)
	}
	instanceName := instruction.StringAt(2)
	fixtureName := instruction.StringAt(3)
	args := instruction.TailAt(4)
	return slimInterpreter.processor.DoMake(instanceName, fixtureName, args)
}

// Process takes an incoming set of instructions, dispatches to statement processor, and retrieves the result.
func (slimInterpreter *SlimInterpreter) Process(instructions *slimentity.SlimList) *slimentity.SlimList {
	results := slimentity.NewSlimList()
	for _, instruction := range *instructions {
		if slimentity.IsSlimList(instruction) {
			instructionList := instruction.(*slimentity.SlimList)
			if instructionList.Length() == 0 {
				addResult(results, malformedInstruction(instructionList))
			} else {
				id := instructionList.StringAt(0)

				if instructionList.Length() == 1 {
					addResult(results, id, malformedInstruction(instructionList))
				} else {
					result := slimInterpreter.dispatchWithTimeout(instructionList)
					addResult(results, id, result)
				}
			}
		} else {
			results.Append(malformedInstruction(instruction))
		}
	}
	return results
}
