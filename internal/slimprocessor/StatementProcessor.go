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
	"reflect"
	"strings"

	"github.com/essenius/slim4go/internal/apperrors"
	"github.com/essenius/slim4go/internal/interfaces"
	"github.com/essenius/slim4go/internal/slimentity"
	"github.com/essenius/slim4go/internal/slimprotocol"
)

// Definitions and constructors

// SlimStatementProcessor is the implementation of the StatementProcessor interface.
type SlimStatementProcessor struct {
	registry interfaces.Registry
	objects  interfaces.ObjectHandler
	parser   interfaces.Parser
	symbols  interfaces.SymbolCollector
}

// NewStatementProcessor returns a new SlimStatementProcesspr.
func NewStatementProcessor(registry interfaces.Registry, objects interfaces.ObjectHandler, parser interfaces.Parser, symbols interfaces.SymbolCollector) *SlimStatementProcessor {
	processor := new(SlimStatementProcessor)
	processor.registry = registry
	processor.objects = objects
	processor.parser = parser
	processor.symbols = symbols
	return processor
}

// Interface methods

// DoCall calls a method (or property) on an instance.
func (processor *SlimStatementProcessor) DoCall(instanceName, methodName string, args *slimentity.SlimList) slimentity.SlimEntity {
	instance := processor.objects.Get(instanceName)
	// The instance can be nil if the test solely relies on the libraries. So that's not a fatal error.
	var result slimentity.SlimEntity
	var err1 error
	if instance != nil {
		result, err1 = processor.objects.InvokeMemberOn(instance, methodName, args)
	} else {
		err1 = &apperrors.NotFoundError{Entity: "instance", Description: instanceName}
	}
	if err1 == nil {
		return result
	}
	if _, ok := err1.(*apperrors.NotFoundError); !ok {
		return slimprotocol.Exception(err1.Error())
	}
	// no object found or no method found on the object instance. Try via the libraries
	libraries := processor.objects.InstancesWithPrefix("library")
	var err2 error
	for _, library := range libraries {
		result, err2 = processor.objects.InvokeMemberOn(library, methodName, args)
		if err2 == nil {
			return result
		}
	}
	if notFoundErr, ok := err1.(*apperrors.NotFoundError); ok {
		// If the instance was not found, best to return that message.
		if notFoundErr.Entity == "instance" {
			return slimprotocol.NoInstance(notFoundErr.Description)
		}
	}
	return slimprotocol.NoMethodInFixture(methodName, reflect.TypeOf(instance).String(), args.Length())
}

// DoImport executes an Slim Import command
func (processor *SlimStatementProcessor) DoImport(value string) slimentity.SlimEntity {
	processor.registry.AddNamespace(value)
	return slimprotocol.OK()
}

// DoMake executes a Make command (creating a new instance).
func (processor *SlimStatementProcessor) DoMake(instanceName, fixtureName string, args *slimentity.SlimList) slimentity.SlimEntity {
	if instance, ok := processor.symbols.NonTextSymbol(fixtureName); ok {
		processor.objects.Add(instanceName, instance)
		return slimprotocol.OK()
	}
	resolvedFixtureName := processor.parser.ReplaceSymbolsIn(fixtureName)
	constructor := processor.registry.FixtureNamed(resolvedFixtureName)
	if constructor == nil {
		return slimprotocol.NoFixture(resolvedFixtureName)
	}
	constructorValue := reflect.ValueOf(constructor)
	if err := processor.objects.AddObjectByConstructor(instanceName, constructorValue, slimentity.ToSlice(args)); err != nil {
		return slimprotocol.CouldNotInvokeConstructor(strings.ReplaceAll(fixtureName+":"+err.Error(), " ", "_"))
	}
	return slimprotocol.OK()
}

// SetSymbol sets a value in the symbol table.
func (processor *SlimStatementProcessor) SetSymbol(symbol string, value interface{}) {
	processor.symbols.Set(symbol, value)
}

// SerializeObjectsIn converts all occurrences of an object in a Slim list to their string representations.
func (processor *SlimStatementProcessor) SerializeObjectsIn(input slimentity.SlimEntity) slimentity.SlimEntity {
	if slimentity.IsSlimList(input) {
		result := slimentity.NewSlimList()
		list := input.(*slimentity.SlimList)
		for _, entry := range *list {
			result.Append(processor.SerializeObjectsIn(entry))
		}
		return result
	}

	if slimentity.IsObject(input) {
		return processor.objects.Serialize(input)
	}
	return input
}
