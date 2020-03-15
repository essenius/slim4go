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
	"testing"

	"github.com/essenius/slim4go/internal/assert"
	"github.com/essenius/slim4go/internal/slimentity"
	. "github.com/essenius/slim4go/internal/slimentity"
)

var symbols *symbolTable

func TestSymbolTableIsValidSymbolName(t *testing.T) {
	symbols := newSymbolTable()
	assert.IsTrue(t, symbols.IsValidSymbolName("test1"), "test1 is valid")
	assert.IsTrue(t, symbols.IsValidSymbolName("Test1Q"), "Test1Q is valid")
	assert.IsTrue(t, !symbols.IsValidSymbolName("$test"), "$test is invalid")
	assert.IsTrue(t, !symbols.IsValidSymbolName("1test"), "1test is invalid")
}

func TestSymbolTable(t *testing.T) {
	symbols := newSymbolTable()

	assert.Equals(t, nil, symbols.SetSymbol("test1", "value1"), "SetSymbol test1 to strinf")
	assert.Equals(t, nil, symbols.SetSymbol("test2", "value2"), "SetSymbol test2 to string")
	aDemoStruct1 := new(demoStruct1)
	aDemoStruct1.Parse("hi from aDemoStruct1")
	assert.Equals(t, nil, symbols.SetSymbol("test3", aDemoStruct1), "SetSymbol test3 to object")
	assert.Equals(t, "Invalid symbol name: $_test3", symbols.SetSymbol("$_test3", "_value3").Error(), "invalid name $_test3")
	assert.Equals(t, "value2", symbols.ReplaceValue("$test2"), "ReplaceValue($test2) returns value=2")
	assert.Equals(t, "$_test3", symbols.ReplaceValue("$_test3"), "ReplaceValue($_test1) returns the input")
	assert.Equals(t, "value1", symbols.ReplaceValue("$test1"), "ReplaceValue($test1) returns value1")
	assert.Equals(t, "hi from aDemoStruct1", symbols.ReplaceValue("$test3"), "ReplaceValue($test3) returns ToString() value")
	assert.Equals(t, "$test4", symbols.ReplaceValue("$test4"), "Nonexisting symbol returns the input")
	assert.Equals(t, "$", symbols.ReplaceSymbolsIn("$"), "just a $ returns input")
	assert.Equals(t, "1€$", symbols.ReplaceSymbolsIn("1€$"), "ending with $ returns input")

	assert.Equals(t, "we see value1, value2, 'hi from aDemoStruct1', $test4 and $_test3 returned as-is",
		symbols.ReplaceSymbolsIn("we see $test1, $test2, '$test3', $test4 and $_test3 returned as-is"), "3 replacements")
	listIn := NewSlimListContaining([]slimentity.SlimEntity{"$$test1", "this is $test2", "$test3", "$test1-$"})
	listOut := symbols.ReplaceSymbols(listIn).(*SlimList)
	assert.Equals(t, `[$value1, this is value2, hi from aDemoStruct1, value1-$]`, listOut.ToString(), "list values replaced OK")
}

func TestSymbolTableNonString(t *testing.T) {
	symbols := newSymbolTable()
	(*symbols)["test1"] = NewMessenger()
	assert.Equals(t, nil, symbols.SetSymbol("test2", "text2"), "SetSymbol with text")
	result, ok := symbols.NonTextSymbol("$test1")
	assert.IsTrue(t, ok, "Identified non-text symbol")
	assert.Equals(t, "*slimprocessor.Messenger", reflect.TypeOf(result).String(), "type correctly identified")

	result, ok = symbols.NonTextSymbol("$test2")
	assert.IsTrue(t, !ok, "Identified text symbol")
}
