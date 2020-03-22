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
	"reflect"
	"strings"

	"github.com/essenius/slim4go/internal/fixture"
	"github.com/essenius/slim4go/internal/slimentity"
	"github.com/essenius/slim4go/internal/slimprotocol"
)

// Definitions and constructors

type notFoundError struct {
	entity      string
	description string
}

func (aNotFoundError *notFoundError) Error() string {
	return fmt.Sprintf("%v: %v not found", aNotFoundError.description, aNotFoundError.entity)
}

type statementProcessor interface {
	fixtureRegistry() *fixture.Registry
	doCall(instanceName, methodName string, args *slimentity.SlimList) slimentity.SlimEntity
	doImport(value string) slimentity.SlimEntity
	doMake(instanceName, fixtureName string, args *slimentity.SlimList) slimentity.SlimEntity
	objects() *objectCollection
	parser() *parser
	setSymbol(symbol string, value interface{}) //remove
	serializeObjectsIn(slimentity.SlimEntity) slimentity.SlimEntity
}

type slimStatementProcessor struct {
	theRegistry *fixture.Registry
	theObjects  *objectCollection
	theParser   *parser
}

func injectStatementProcessor() statementProcessor {
	parser := injectParser()
	processor := newStatementProcessor(fixture.InjectRegistry(), injectObjectCollection(parser), parser)
	return processor
}

func newStatementProcessor(registry *fixture.Registry, objects *objectCollection, aParser *parser) statementProcessor {
	processor := new(slimStatementProcessor)
	processor.theRegistry = registry
	processor.theObjects = objects
	processor.theParser = aParser
	return processor
}

// Methods. First the different tables

func (processor *slimStatementProcessor) fixtureRegistry() *fixture.Registry {
	return processor.theRegistry
}

func (processor *slimStatementProcessor) objects() *objectCollection {
	return processor.theObjects
}

func (processor *slimStatementProcessor) parser() *parser {
	return processor.theParser
}

// Helper methods

func (processor *slimStatementProcessor) findObject(instanceName string) (*object, error) {
	// instance can't use symbols, so this becomes quite easy
	anObject := processor.objects().objectNamed(instanceName)
	if anObject != nil {
		return anObject, nil
	}
	return nil, &notFoundError{"instance", instanceName}
}

// Interface methods

func (processor *slimStatementProcessor) doCall(instanceName, methodName string, args *slimentity.SlimList) slimentity.SlimEntity {
	anObject, err := processor.findObject(instanceName)
	if anObject != nil {
		var result slimentity.SlimEntity
		result, err = anObject.InvokeMember(methodName, args)
		if err == nil {
			return result
		}
	}
	if _, ok := err.(*notFoundError); !ok {
		return slimprotocol.Exception(err.Error())
	}
	// no object found or no method found on the object instance. Try via the libraries
	libraries := processor.objects().objectsWithPrefix("library")
	for _, library := range *libraries {
		result, err := library.InvokeMember(methodName, args)
		if err == nil {
			return result
		}
	}
	if notFoundErr, ok := err.(*notFoundError); ok {
		if notFoundErr.entity == "instance" {
			return slimprotocol.NoInstance(notFoundErr.description)
		}
	}
	return slimprotocol.NoMethodInFixture(methodName, reflect.TypeOf(anObject.instance()).String(), args.Length())
}

func (processor *slimStatementProcessor) doImport(value string) slimentity.SlimEntity {
	processor.fixtureRegistry().AddNamespace(value)
	return slimprotocol.OK()
}

func (processor *slimStatementProcessor) doMake(instanceName, fixtureName string, args *slimentity.SlimList) slimentity.SlimEntity {
	if instance, ok := processor.parser().symbols.NonTextSymbol(fixtureName); ok {
		processor.objects().addObject(instanceName, instance)
		return slimprotocol.OK()
	}
	resolvedFixtureName := processor.parser().ReplaceSymbolsIn(fixtureName)
	constructor := processor.fixtureRegistry().FixtureNamed(resolvedFixtureName)
	if constructor == nil {
		return slimprotocol.NoFixture(resolvedFixtureName)
	}
	constructorValue := reflect.ValueOf(constructor)
	if err := processor.objects().addObjectByConstructor(instanceName, constructorValue, slimentity.ToSlice(args)); err != nil {
		return slimprotocol.CouldNotInvokeConstructor(strings.ReplaceAll(fixtureName+":"+err.Error(), " ", "_"))
	}
	return slimprotocol.OK()
}

func (processor *slimStatementProcessor) setSymbol(symbol string, value interface{}) {
	processor.parser().symbols.SetSymbol(symbol, value)
}

func (processor *slimStatementProcessor) serializeObjectsIn(input slimentity.SlimEntity) slimentity.SlimEntity {
	if slimentity.IsSlimList(input) {
		result := slimentity.NewSlimList()
		list := input.(*slimentity.SlimList)
		for _, entry := range *list {
			result.Append(processor.serializeObjectsIn(entry))
		}
		return result
	}
	inputValue := reflect.ValueOf(input)
	if slimentity.IsObject(inputValue) {
		anObject := processor.objects().factory.NewObject(inputValue)
		return anObject.serialize()
	}
	return input
}
