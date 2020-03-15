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

	"github.com/essenius/slim4go/internal/slimentity"
	"github.com/essenius/slim4go/internal/slimprotocol"
	"github.com/essenius/slim4go/internal/utilities"
)

type functionCaller struct {
	theParser *parser
}

func newFunctionCaller(aParser *parser) *functionCaller {
	caller := new(functionCaller)
	caller.theParser = aParser
	return caller
}

func injectFunctionCaller(symbols *symbolTable) *functionCaller {
	aParser := newParser(symbols)
	caller := newFunctionCaller(aParser)
	aParser.caller = caller
	caller.theParser = aParser
	return caller
}

// Helper functions

func isObjectType(inputType reflect.Type) bool {
	inputKind := inputType.Kind()
	return inputKind == reflect.Struct || (inputKind == reflect.Ptr && isObjectType(inputType.Elem()))
}

func isObject(inputValue reflect.Value) bool {
	return isObjectType(inputValue.Type())
}

func paramTypeFor(methodType reflect.Type, paramIndex int) reflect.Type {
	if methodType.IsVariadic() && paramIndex >= methodType.NumIn()-1 {
		return methodType.In(methodType.NumIn() - 1).Elem()
	}
	return methodType.In(paramIndex)
}

func stringifyObject(inputValue reflect.Value) slimentity.SlimEntity {
	const toStringFunction = "ToString"
	toStringMethod := inputValue.MethodByName(toStringFunction)

	if toStringMethod.IsValid() {
		callResult := toStringMethod.Call([]reflect.Value{})
		return transformCallResult(callResult)
	}
	return inputValue.Type().String()
}

func stringifyObjectsIn(input slimentity.SlimEntity) slimentity.SlimEntity {
	if slimentity.IsSlimList(input) {
		result := slimentity.NewSlimList()
		list := input.(*slimentity.SlimList)
		for _, entry := range *list {
			result.Append(stringifyObjectsIn(entry))
		}
		return result
	}
	inputValue := reflect.ValueOf(input)
	if isObject(inputValue) {
		return stringifyObject(inputValue)
	}
	return input
}

// Convert the result of a call to a string representation, or an object pointer.
func transformCallResult(callOutput []reflect.Value) slimentity.SlimEntity {
	count := len(callOutput)
	if count == 0 {
		return slimprotocol.Void()
	}
	if count == 1 {
		return valueToSlimEntity(callOutput[0])
	}
	resultList := new(slimentity.SlimList)
	for result := 0; result < count; result++ {
		resultList.Append(valueToSlimEntity(callOutput[result]))
	}
	return resultList
}

func valueToSlimEntity(inputValue reflect.Value) slimentity.SlimEntity {
	if !inputValue.IsValid() {
		return slimprotocol.Null()
	}
	// For predefined types, use fmt.Sprintf
	if isPredefinedType(inputValue.Type()) {
		return fmt.Sprintf("%v", inputValue.Interface())
	}
	if isObject(inputValue) {
		return inputValue.Interface()
	}

	switch inputValue.Kind() {
	case reflect.Ptr, reflect.Interface:
		// This is a non-object pointer. Resolve the element
		return valueToSlimEntity(inputValue.Elem())
	// Unravel arrays and slices
	case reflect.Array, reflect.Slice:
		result := slimentity.NewSlimList()
		for i := 0; i < inputValue.Len(); i++ {
			entry := valueToSlimEntity(inputValue.Index(i))
			result.Append(entry)
		}
		return result
	// If we can't do anything else, return the type
	case reflect.Map:
		tableTemplate := "<table class=\"hash_table\">\n%v</table>"
		rowTemplate := "  <tr class=\"hash_row\">\n    <td class=\"hash_key\">%v</td>\n    <td class=\"hash_value\">%v</td>\n  </tr>\n"
		result := ""
		iterator := inputValue.MapRange()
		for iterator.Next() {
			result += fmt.Sprintf(rowTemplate, iterator.Key(), iterator.Value())
		}
		return fmt.Sprintf(tableTemplate, result)
	default:
		return inputValue.Type().String()
	}
}

// Methods

func (caller *functionCaller) call(function reflect.Value, name string, args []string /* *slimentity.SlimList */) (returnEntity slimentity.SlimEntity, err error) {
	arguments, err := caller.matchParamType(args, function)
	if err != nil {
		return "", err
	}
	// The function we call might panic (after all, FitNesse is a testing framework). Be ready for that.
	defer func() {
		if panicData := recover(); panicData != nil {
			returnEntity = nil
			err = fmt.Errorf("Panic: %v", utilities.ErrorToString(panicData))
		}
	}()
	returnValue := function.Call(*arguments)
	return transformCallResult(returnValue), nil
}

func (caller *functionCaller) matchParamType(paramIn []string, method reflect.Value) (*[]reflect.Value, error) {
	result := []reflect.Value{}
	methodType := method.Type()
	if methodType.Kind() != reflect.Func {
		return nil, toErrorf("%v is not a function", methodType.String())
	}
	numParams := methodType.NumIn()
	var paramCountMatch bool
	if methodType.IsVariadic() {
		paramCountMatch = len(paramIn) >= numParams-1
	} else {
		paramCountMatch = len(paramIn) == numParams
	}
	if !paramCountMatch {
		return nil, toErrorf("Expected %v parameter(s) but got %v", numParams, len(paramIn))
	}
	for paramIndex, param := range paramIn {
		paramType := paramTypeFor(methodType, paramIndex)
		resultValue, err := caller.theParser.parse(param, paramType)
		if err != nil {
			return nil, err
		}
		result = append(result, reflect.ValueOf(resultValue))
	}
	return &result, nil
}
