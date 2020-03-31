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
)

func TestSymbolTableIsValidSymbolName(t *testing.T) {
	symbols := NewSymbolTable()
	assert.IsTrue(t, symbols.IsValidSymbolName("test1"), "test1 is valid")
	assert.IsTrue(t, symbols.IsValidSymbolName("Test1Q"), "Test1Q is valid")
	assert.IsTrue(t, !symbols.IsValidSymbolName("$test"), "$test is invalid")
	assert.IsTrue(t, !symbols.IsValidSymbolName("1test"), "1test is invalid")
}

func TestSymbolTable(t *testing.T) {
	symbols := NewSymbolTable()

	assert.Equals(t, nil, symbols.Set("test1", "value1"), "Set test1 to strinf")
	assert.Equals(t, nil, symbols.Set("test2", "value2"), "Set test2 to string")
	aDemoStruct1 := new(demoStruct1)
	aDemoStruct1.Parse("hi from aDemoStruct1")
	assert.Equals(t, nil, symbols.Set("test3", aDemoStruct1), "Set test3 to object")
	assert.Equals(t, "Invalid symbol name: $_test3", symbols.Set("$_test3", "_value3").Error(), "invalid name $_test3")
}

func TestSymbolTableNonString(t *testing.T) {
	symbols := NewSymbolTable()
	(*symbols)["test1"] = NewMessenger()
	assert.Equals(t, nil, symbols.Set("test2", "text2"), "Set with text")
	result, ok := symbols.NonTextSymbol("$test1")
	assert.IsTrue(t, ok, "Identified non-text symbol")
	assert.Equals(t, "*slimprocessor.Messenger", reflect.TypeOf(result).String(), "type correctly identified")

	result, ok = symbols.NonTextSymbol("$test2")
	assert.IsTrue(t, !ok, "Identified text symbol")
}
