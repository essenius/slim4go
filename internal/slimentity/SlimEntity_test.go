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

package slimentity

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/essenius/slim4go/internal/assert"
)

var slimList SlimList

type emptyStruct struct{}

func TestSlimListBaseTests(t *testing.T) {
	var list0 *SlimList = nil
	assert.Equals(t, 0, list0.Length(), "list1 has length 0")
	list1 := NewSlimList()
	list2 := NewSlimList()
	assert.IsTrue(t, list1.Equals(list2), "two empty lists are equal")
	assert.Equals(t, "[]", list1.ToString(), "Empty list ToString")
	assert.Equals(t, "[]", fmt.Sprintf("%v", ToSlice(list1)), "empty list to slice")

	assert.Equals(t, 0, list1.Length(), "list1 has length 0")
	var entity SlimEntity = "a"
	assert.Equals(t, "a", ToString(entity), "Entity with string value ToString")
	list1.Append(entity)

	assert.Equals(t, "[a]", list1.ToString(), "List with 1 item ToString")
	assert.Equals(t, 1, list1.Length(), "list1 has length 1")
	assert.Equals(t, "a", list1.ElementAt(0), "List1.ElementAt(0)==a")
	assert.IsTrue(t, !list1.Equals(list2), "empty list not equal to list with one element")
	list2.Append("a")
	assert.IsTrue(t, list1.Equals(list2), "two lists with the same single value")
	list1.Append("2")
	assert.Equals(t, "[a, 2]", ToString(list1), "List with 2 items ToString")

	list2.Append(2)
	assert.IsTrue(t, !list1.Equals(list2), "list with entries of different types")
	assert.Equals(t, "2", list1.StringAt(1), "list1.ElementAt(1)=='2'")
	assert.Equals(t, 2, list2.ElementAt(1), "list2.ElementAt(1)==2")
	assert.Equals(t, 2, list2.Length(), "list2 has length 2")

	list2 = NewSlimList()
	list2.Append("a")
	list2.Append("2")
	assert.IsTrue(t, list1.Equals(list2), "lists with two equal elements")
	list1.Append("a")
	list1.Append("b")
	list3 := list1.TailAt(2)
	assert.Equals(t, 4, list1.Length(), "list1.Length() == 4 after adding a slice with length 4")
	assert.Equals(t, 2, list3.Length(), "Tail list.Length == 2")
	assert.Equals(t, "a", list3.ElementAt(0), "Tail head is 'a'")
	assert.Equals(t, "b", list3.ElementAt(1), "Tail tail is 'b'")
	list1.Append(list3)
	assert.Equals(t, "[a, 2, a, b, [a, b]]", list1.ToString(), "nested list")
	assert.Equals(t, "[a 2 a b [a, b]]", fmt.Sprintf("%v", ToSlice(list1)), "nested list to slice")

	list4 := list1.ElementAt(4).(*SlimList)
	assert.Equals(t, 2, list4.Length(), "List in list length == 2")
	assert.Equals(t, "a", list4.ElementAt(0), "List in list element 0 == 'a'")
	list5 := list3.TailAt(2)
	assert.Equals(t, 0, list5.Length(), "empty tail")
	list6 := NewSlimListContaining([]SlimEntity{"1", "2"})
	assert.Equals(t, "[1, 2]", list6.ToString(), "NewSlimListContaining created")
	assert.Equals(t, "[1 2]", fmt.Sprintf("%v", ToSlice(list6)), "to slice")
}

func TestSlimEntityIsObject(t *testing.T) {
	assert.IsTrue(t, !IsObject(""), "empty string is no object")
	assert.IsTrue(t, IsObject(t), "t (pointer to struct T) is an object")
}

func TestSlimEntityIsObjectValue(t *testing.T) {
	assert.IsTrue(t, !isObjectValue(reflect.ValueOf("")), "empty string is no object")
	assert.IsTrue(t, isObjectValue(reflect.ValueOf(t)), "t (pointer to struct T) is an object")
}

func TestSlimEntityTransformCallResult(t *testing.T) {
	output := []reflect.Value{}
	assert.Equals(t, "/__VOID__/", TransformCallResult(output), "empty output")
	output = append(output, reflect.ValueOf(35))
	assert.Equals(t, "35", TransformCallResult(output), "1 int output")
	output = append(output, reflect.ValueOf("test"))
	outputList := TransformCallResult(output)
	assert.IsTrue(t, IsSlimList(outputList), "is a SlimList")
	assert.Equals(t, 2, outputList.(*SlimList).Length(), "list length")
	assert.Equals(t, "35", outputList.(*SlimList).ElementAt(0), "list entry 0")
	assert.Equals(t, "test", outputList.(*SlimList).ElementAt(1), "list entry 1")
}

func TestSlimEntityValueToSlimEntity(t *testing.T) {
	testSliceList := func(list interface{}) {
		outList := valueToSlimEntity(reflect.ValueOf(list))
		assert.IsTrue(t, IsSlimList(outList), "Is SlimList")
		assert.Equals(t, "[1, a, true]", outList.(*SlimList).ToString(), "slice content OK")
	}

	assert.Equals(t, "null", valueToSlimEntity(reflect.ValueOf(SlimEntity(nil))), "nil")
	assert.Equals(t, "1", valueToSlimEntity(reflect.ValueOf(1)), "int")
	assert.Equals(t, "Test", valueToSlimEntity(reflect.ValueOf("Test")), "string")
	assert.Equals(t, "1.2", valueToSlimEntity(reflect.ValueOf(1.2)), "float64")
	assert.Equals(t, "true", valueToSlimEntity(reflect.ValueOf(true)), "bool")
	assert.Equals(t, "func(*testing.T)", valueToSlimEntity(reflect.ValueOf(TestSlimEntityValueToSlimEntity)), "func")

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
