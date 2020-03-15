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
	"testing"

	"github.com/essenius/slim4go/internal/assert"
	"github.com/essenius/slim4go/internal/slimentity"
)

func TestFunctionCallerCallFunction(t *testing.T) {
	caller := injectFunctionCaller(newSymbolTable())
	result1, err1 := caller.call(reflect.ValueOf(NewObjectWithPanic), "Object", []string{})
	assert.Equals(t, "Panic: Object creation failed", err1.Error(), "Panicking function")
	assert.Equals(t, nil, result1, "No result with panic")
	messengerInstance, err2 := caller.call(reflect.ValueOf(NewMessenger), "Messenger", []string{})
	assert.Equals(t, nil, err2, "No error calling function")
	assert.Equals(t, "*slimprocessor.Messenger", reflect.TypeOf(messengerInstance).String(), "Type of instance OK")
	result3, err3 := caller.call(reflect.ValueOf(NewMessenger), "Messenger", []string{"q"})
	assert.Equals(t, "Expected 0 parameter(s) but got 1", err3.Error(), "Create messenger with wrong parameter")
	assert.Equals(t, "", result3, "No result with parameter error")
}

func TestFunctionCallerTransformCallResult(t *testing.T) {
	output := []reflect.Value{}
	assert.Equals(t, "/__VOID__/", transformCallResult(output), "empty output")
	output = append(output, reflect.ValueOf(35))
	assert.Equals(t, "35", transformCallResult(output), "1 int output")
	output = append(output, reflect.ValueOf("test"))
	outputList := transformCallResult(output)
	assert.IsTrue(t, slimentity.IsSlimList(outputList), "is a SlimList")
	assert.Equals(t, 2, outputList.(*slimentity.SlimList).Length(), "list length")
	assert.Equals(t, "35", outputList.(*slimentity.SlimList).ElementAt(0), "list entry 0")
	assert.Equals(t, "test", outputList.(*slimentity.SlimList).ElementAt(1), "list entry 1")
}

func TestFunctionCallerMatchParamType(t *testing.T) {
	caller := newFunctionCaller(newParser(newSymbolTable()))
	factoryType := reflect.TypeOf(FixtureFactory{})
	factory := reflect.Zero(reflect.PtrTo(factoryType))
	assert.Equals(t, "*slimprocessor.FixtureFactory", factory.Type().String(), "Factory type O<")
	method := factory.MethodByName("NewOrder")
	assert.IsTrue(t, method.IsValid(), "method valid")
	params := []string{"test", "100"}
	_, err := caller.matchParamType(params, method)
	assert.Equals(t, "Expected 3 parameter(s) but got 2", err.Error(), "Wrong number of parameters")
	params = append(params, "25")
	result, err := caller.matchParamType(params, method)
	assert.Equals(t, nil, err, "No error")
	assert.Equals(t, "string", (*result)[0].Type().String(), "type of param 0 is string")
	assert.Equals(t, "test", (*result)[0].Interface(), "value of param 0 is test")
	assert.Equals(t, "float64", (*result)[1].Type().String(), "type of param 1 is float")
	assert.Equals(t, 100.0, (*result)[1].Interface(), "value of param 1 is 100.0")
	assert.Equals(t, "int", (*result)[2].Type().String(), "type of param 2 is int")
	assert.Equals(t, 25, (*result)[2].Interface(), "value of param 2 is 25")
	params = []string{"test", "100", "q"}
	_, err = caller.matchParamType(params, method)
	assert.Equals(t, "Could not convert 'q' to type 'int'", err.Error(), "invalid conversion")
}

func TestFunctionCallerMatchParamTypeVariadic(t *testing.T) {
	sumFunction := func(num ...int) int {
		sum := 0
		for _, value := range num {
			sum += value
		}
		return sum
	}
	printFunction := func(template string, args ...interface{}) string {
		return fmt.Sprintf(template, args...)
	}
	caller := newFunctionCaller(newParser(newSymbolTable()))
	args := []string{"2", "3", "5"}
	result1, err1 := caller.matchParamType(args, reflect.ValueOf(sumFunction))
	assert.Equals(t, nil, err1, "variadic1: No error")
	assert.Equals(t, 3, len(*result1), "variadic1: 3 params")
	assert.Equals(t, "int", (*result1)[0].Type().String(), "variadic: type of param is int")
	assert.Equals(t, 2, (*result1)[0].Interface().(int), "first param value is 2")
	assert.Equals(t, 3, (*result1)[1].Interface().(int), "second param value is 3")
	assert.Equals(t, 5, (*result1)[2].Interface().(int), "thirdd param value is 5")

	emptyArgs := []string{}
	result2, err2 := caller.matchParamType(emptyArgs, reflect.ValueOf(sumFunction))
	assert.Equals(t, nil, err2, "variadic2 empty: No error")
	assert.Equals(t, 0, len(*result2), "variadic2 empty: 0 params")

	args3 := []string{"param %v %v", "3", "5.5"}
	result3, err3 := caller.matchParamType(args3, reflect.ValueOf(printFunction))
	assert.Equals(t, nil, err3, "Variadic3 param interface: No error")
	assert.Equals(t, 3, len(*result3), "Variadic3 param interface: 3 params")
	assert.Equals(t, "string", (*result3)[0].Type().String(), "Variadic3 param interface: type of param[0] is string")
	assert.Equals(t, "int64", (*result3)[1].Type().String(), "Variadic3 param interface: type of param[0] is interface{}")
	assert.Equals(t, "float64", (*result3)[2].Type().String(), "Variadic3 param interface: type of param[0] is interface{}")
	assert.Equals(t, "param %v %v", (*result3)[0].Interface().(string), "first param value is ok")
	assert.Equals(t, int64(3), (*result3)[1].Interface(), "second param value is 3")
	assert.Equals(t, 5.5, (*result3)[2].Interface(), "third param value is 5.5")
}

func TestFunctionCallerIsObject(t *testing.T) {
	assert.IsTrue(t, !isObject(reflect.ValueOf("")), "empty string is no object")
	assert.IsTrue(t, isObject(reflect.ValueOf(&demoStruct1{})), "pointer to demoStruct1 is an object")
}

func TestFunctionCallerStringifyObjectsIn(t *testing.T) {
	test1 := "test1"
	assert.Equals(t, test1, stringifyObjectsIn(test1), "string")
	assert.IsTrue(t, !isObject(reflect.ValueOf(test1)), "test1 is no object")
	aDemoStruct1 := &demoStruct1{"demo"}
	aDemoStruct1.Parse("demo1")
	assert.IsTrue(t, isObject(reflect.ValueOf(aDemoStruct1)), "aDemoStruct1 is an object")
	assert.Equals(t, "demo1", stringifyObjectsIn(aDemoStruct1), "*struct with *ToString")
	aDemoStruct2 := demoStruct2{"demo2"}
	assert.IsTrue(t, isObject(reflect.ValueOf(aDemoStruct2)), "aDemoStruct2 is an object")
	assert.Equals(t, "demo2", stringifyObjectsIn(aDemoStruct2), "struct with ToString")
	ptrToADemoStruct2 := &aDemoStruct2
	assert.IsTrue(t, isObject(reflect.ValueOf(ptrToADemoStruct2)), "*aDemoStruct2 is an object")
	assert.Equals(t, "demo2", stringifyObjectsIn(ptrToADemoStruct2), "*struct with ToString")
	anEmptyStruct := emptyStruct{}
	assert.Equals(t, "slimprocessor.emptyStruct", stringifyObjectsIn(anEmptyStruct), "struct without ToString")
	list := slimentity.NewSlimListContaining([]slimentity.SlimEntity{"test2", aDemoStruct1, aDemoStruct2, anEmptyStruct})
	assert.IsTrue(t, !isObject(reflect.ValueOf(list)), "list is no object")
	assert.Equals(t, "[test2, demo1, demo2, slimprocessor.emptyStruct]",
		stringifyObjectsIn(list).(*slimentity.SlimList).ToString(), "list with objects")
}

func TestFunctionCallerValueToSlimEntity(t *testing.T) {

	testSliceList := func(list interface{}) {
		outList := valueToSlimEntity(reflect.ValueOf(list))
		assert.IsTrue(t, slimentity.IsSlimList(outList), "Is SlimList")
		assert.Equals(t, "[1, a, true]", outList.(*slimentity.SlimList).ToString(), "slice content OK")
	}

	assert.Equals(t, "null", valueToSlimEntity(reflect.ValueOf(slimentity.SlimEntity(nil))), "nil")
	assert.Equals(t, "1", valueToSlimEntity(reflect.ValueOf(1)), "int")
	assert.Equals(t, "Test", valueToSlimEntity(reflect.ValueOf("Test")), "string")
	assert.Equals(t, "1.2", valueToSlimEntity(reflect.ValueOf(1.2)), "float64")
	assert.Equals(t, "true", valueToSlimEntity(reflect.ValueOf(true)), "bool")
	assert.Equals(t, "func(*testing.T)", valueToSlimEntity(reflect.ValueOf(TestFunctionCallerValueToSlimEntity)), "func")

	aSlice := []interface{}{1, "a", true}
	aPointer := &aSlice
	testSliceList(aSlice)
	testSliceList(aPointer)
	anEmptyStruct := emptyStruct{}
	ptrToAnEmptyStruct := &anEmptyStruct
	assert.Equals(t, anEmptyStruct, valueToSlimEntity(reflect.ValueOf(anEmptyStruct)), "struct")
	assert.Equals(t, ptrToAnEmptyStruct, valueToSlimEntity(reflect.ValueOf(ptrToAnEmptyStruct)), "pointer to struct")
	aMap := make(map[string]int)
	aMap["size"] = 50
	result := "<table class=\"hash_table\">\n  <tr class=\"hash_row\">\n    <td class=\"hash_key\">size</td>\n    <td class=\"hash_value\">50</td>\n  </tr>\n</table>"
	assert.Equals(t, result, valueToSlimEntity(reflect.ValueOf(aMap)), "map")
}
