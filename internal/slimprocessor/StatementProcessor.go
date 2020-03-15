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

// Comments
// TODO * Method matcher needs to get smarter - work with Get, Set, etc.

import (
	"fmt"
	"reflect"
	"strings"

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
	actors() *actorStack
	fixtures() *FixtureMap
	doCall(instanceName, methodName string, args *slimentity.SlimList) slimentity.SlimEntity
	doImport(value string) slimentity.SlimEntity
	doMake(instanceName, fixtureName string, args *slimentity.SlimList) slimentity.SlimEntity
	objects() *objectMap
	funcCaller() *functionCaller
	setSymbol(symbol string, value interface{})
	symbols() *symbolTable
}

type slimStatementProcessor struct {
	theActors   *actorStack
	theFixtures *FixtureMap
	theObjects  *objectMap
	//theParser   *parser
	theCaller  *functionCaller
	theSymbols *symbolTable
}

func injectStatementProcessor() statementProcessor {
	symbolTable := injectSymbolTable()
	processor := newStatementProcessor(injectActors(), InjectFixtures(), injectObjects(), injectFunctionCaller(symbolTable), symbolTable)
	processor.objects().addObject("libraryStandard", newStandardLibrary(processor.actors(), processor.objects(), processor.symbols()), "StandardLibrary")
	return processor
}

func newStatementProcessor(actors *actorStack, fixtures *FixtureMap, objects *objectMap, caller *functionCaller, symbols *symbolTable) statementProcessor {
	processor := new(slimStatementProcessor)
	processor.theFixtures = fixtures
	processor.theObjects = objects
	processor.theSymbols = symbols
	processor.theCaller = caller
	processor.theActors = actors
	return processor
}

// Methods. First the different tables

func (processor *slimStatementProcessor) actors() *actorStack {
	return processor.theActors
}

func (processor *slimStatementProcessor) fixtures() *FixtureMap {
	return processor.theFixtures
}

func (processor *slimStatementProcessor) objects() *objectMap {
	return processor.theObjects
}

func (processor *slimStatementProcessor) symbols() *symbolTable {
	return processor.theSymbols
}

func (processor *slimStatementProcessor) funcCaller() *functionCaller {
	return processor.theCaller
}

// Helper methods

/*func (processor *slimStatementProcessor) callFunction(function reflect.Value, name string, args *slimentity.SlimList) (returnEntity slimentity.SlimEntity, err error) {
	arguments, err := matchParamType(args, function, processor.parser())
	if err != nil {
		return "", err
	}
	// The function we call might panic (after all, FitNesse is a testing framework). Be ready for that.
	defer func() {
		if panicData := recover(); panicData != nil {
			returnEntity = nil
			err = fmt.Errorf("Panic: %v", errorToString(panicData))
		}
	}()
	returnValue := function.Call(*arguments)
	return transformCallResult(returnValue), nil
} */

func (processor *slimStatementProcessor) callMethod(instance interface{}, methodName string, args *slimentity.SlimList) (slimentity.SlimEntity, error) {
	// We only use exported methods. Since in the Java convention, methods are in camelCase, we need to capitalize the first letter.
	MethodName := strings.Title(methodName)
	instanceValue := reflect.ValueOf(instance)

	method := instanceValue.MethodByName(MethodName)
	if !method.IsValid() {
		return "", &notFoundError{"method", MethodName}
	}

	return processor.funcCaller().call(method, methodName, slimentity.ToSlice(args))
}

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
	var err error
	object, err := processor.findObject(instanceName)
	if object != nil {
		var result slimentity.SlimEntity
		result, err = processor.callMethod(object.instance, methodName, args)
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
		result, err := processor.callMethod(library.instance, methodName, args)
		if err == nil {
			return result
		}
	}
	if notFoundErr, ok := err.(*notFoundError); ok {
		if notFoundErr.entity == "instance" {
			return slimprotocol.NoInstance(notFoundErr.description)
		}
	}
	return slimprotocol.NoMethodInFixture(methodName, object.fixtureName, args.Length())
}

func (processor *slimStatementProcessor) doImport(value string) slimentity.SlimEntity {
	return slimprotocol.OK()
}

func (processor *slimStatementProcessor) doMake(instanceName, fixtureName string, args *slimentity.SlimList) slimentity.SlimEntity {

	if instance, ok := processor.symbols().NonTextSymbol(fixtureName); ok {
		processor.objects().addObject(instanceName, instance, fixtureName)
	} else {
		resolvedFixtureName := processor.theSymbols.ReplaceSymbolsIn(fixtureName)
		aFixture := processor.fixtures().fixtureNamed(resolvedFixtureName)
		if aFixture == nil {
			return slimprotocol.NoFixture(resolvedFixtureName)
		}
		constructor := reflect.ValueOf(aFixture.constructor)
		instance, err := processor.funcCaller().call(constructor, resolvedFixtureName+" Constructor", slimentity.ToSlice(args))

		if err == nil {
			processor.objects().addObject(instanceName, instance, fixtureName)
			return slimprotocol.OK()
		}
		return slimprotocol.CouldNotInvokeConstructor(strings.ReplaceAll(fixtureName+":"+err.Error(), " ", "_"))
	}
	return slimprotocol.OK()
}

func (processor *slimStatementProcessor) setSymbol(symbol string, value interface{}) {
	processor.symbols().SetSymbol(symbol, value)
}
