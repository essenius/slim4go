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

	"github.com/essenius/slim4go/internal/slimcontext"
	"github.com/essenius/slim4go/internal/slimentity"
	"github.com/essenius/slim4go/internal/slimprotocol"
)

// Definitions and constructors

// SlimProcessor processes incomming SlimLists and dispatches commands to its statement processor.
type SlimProcessor struct {
	processor statementProcessor
	timeout   time.Duration
}

// InjectSlimProcessor is the IoC entry point to provide a SlimProcessor.
func InjectSlimProcessor() *SlimProcessor {
	return newSlimProcessor(injectStatementProcessor(), slimcontext.InjectContext().InstructionTimeout)
}

func newSlimProcessor(processor statementProcessor, timeout time.Duration) *SlimProcessor {
	slimProcessor := new(SlimProcessor)
	slimProcessor.processor = processor
	slimProcessor.timeout = timeout
	return slimProcessor
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

func (slimProcessor *SlimProcessor) dispatch(instruction *slimentity.SlimList) slimentity.SlimEntity {
	command := instruction.StringAt(1)
	switch command {
	case "assign":
		return slimProcessor.doAssign(instruction)
	case "call":
		return slimProcessor.doCall(instruction, noAssign)
	case "callAndAssign":
		return slimProcessor.doCall(instruction, assign)
	case "import":
		return slimProcessor.doImport(instruction)
	case "make":
		return slimProcessor.doMake(instruction)
	default:
		return malformedInstruction(instruction)
	}
}

func (slimProcessor *SlimProcessor) dispatchWithTimeout(instruction *slimentity.SlimList) slimentity.SlimEntity {
	resultChannel := make(chan slimentity.SlimEntity, 1)
	instructionTimer := time.NewTimer(slimProcessor.timeout)
	go func() {
		returnValue := slimProcessor.dispatch(instruction)
		resultChannel <- returnValue
	}()
	select {
	case result := <-resultChannel:
		instructionTimer.Stop()
		return result
	case <-instructionTimer.C:
		return slimprotocol.TimedOut(slimProcessor.timeout)
	}
}

func (slimProcessor *SlimProcessor) doAssign(instruction *slimentity.SlimList) string {
	if instruction.Length() < 3 {
		return malformedInstruction(instruction)
	}
	symbolName := instruction.StringAt(2)
	value := instruction.StringAt(3)
	slimProcessor.processor.setSymbol(symbolName, value)
	return slimprotocol.OK()
}

const (
	noAssign = 4
	assign   = 5
)

func (slimProcessor *SlimProcessor) doCall(instruction *slimentity.SlimList, minLength int) slimentity.SlimEntity {
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
	result := slimProcessor.processor.doCall(instanceName, methodName, args)
	if minLength == assign {
		slimProcessor.processor.setSymbol(symbolName, result)
	}
	return slimProcessor.processor.serializeObjectsIn(result)
}

func (slimProcessor *SlimProcessor) doImport(instruction *slimentity.SlimList) slimentity.SlimEntity {
	if instruction.Length() < 3 {
		return malformedInstruction(instruction)
	}
	pathName := instruction.StringAt(2)
	return slimProcessor.processor.doImport(pathName)
}

func (slimProcessor *SlimProcessor) doMake(instruction *slimentity.SlimList) slimentity.SlimEntity {
	if instruction.Length() < 4 {
		return malformedInstruction(instruction)
	}
	instanceName := instruction.StringAt(2)
	fixtureName := instruction.StringAt(3)
	args := instruction.TailAt(4)
	return slimProcessor.processor.doMake(instanceName, fixtureName, args)
}

// Process takes an incoming set of instructions, dispatches to statement processor, and retrieves the result.
func (slimProcessor *SlimProcessor) Process(instructions *slimentity.SlimList) *slimentity.SlimList {
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
					result := slimProcessor.dispatchWithTimeout(instructionList)
					addResult(results, id, result)
				}
			}
		} else {
			results.Append(malformedInstruction(instruction))
		}
	}
	return results
}
