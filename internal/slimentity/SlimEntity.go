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

	"github.com/essenius/slim4go/internal/slimprotocol"
)

// Definitions and constructors

// SlimEntity is either a string, an object or a SlimList.
type SlimEntity interface{}

// SlimList is a list of SlimEntities
type SlimList []SlimEntity

// NewSlimList creates a new SlimList.
func NewSlimList() *SlimList {
	var newList = make(SlimList, 0, 0)
	return &newList
}

// NewSlimListContaining creates a new SlimList with defined items.
func NewSlimListContaining(listToAdd SlimList) *SlimList {
	list := NewSlimList()
	for _, entry := range listToAdd {
		list.Append(entry)
	}
	return list
}

// Functions

// IsObject returns whether an instance could be a valid object type
func IsObject(instance interface{}) bool {
	return IsObjectType(reflect.TypeOf(instance))
}

// IsObjectType returns whether the type could be that of an object.
func IsObjectType(inputType reflect.Type) bool {
	inputKind := inputType.Kind()
	return inputKind == reflect.Struct || (inputKind == reflect.Ptr && IsObjectType(inputType.Elem()))
}

func isObjectValue(inputValue reflect.Value) bool {
	return IsObjectType(inputValue.Type())
}

func isPredefinedType(inputType reflect.Type) bool {
	return inputType.Name() != "" && inputType.PkgPath() == ""
}

// IsSlimList checks whether entity is a SlimList.
func IsSlimList(entity SlimEntity) bool {
	slimListType := reflect.PtrTo(reflect.TypeOf((*SlimList)(nil)).Elem())
	return reflect.TypeOf(entity) == slimListType
}

// ToSlice converts a list to a slice of strings. Only works for simple lists.
// Nested lists will be serialized with brackets and commas.
// This allows for slices as parameters in functions (the basis for parse is a string)
func ToSlice(list *SlimList) []string {
	result := []string{}
	for _, entry := range *list {
		result = append(result, ToString(entry))
	}
	return result
}

// ToString converts entity to string representation.
func ToString(entity SlimEntity) string {
	if IsSlimList(entity) {
		return entity.(*SlimList).ToString()
	}
	return entity.(string)
}

// Methods

//Append adds an item to the list.
func (list *SlimList) Append(value interface{}) {
	*list = append(*list, value)
}

// ElementAt returns the element at the index.
func (list *SlimList) ElementAt(index int) interface{} {
	return (*list)[index]
}

// Equals returns whether the contents of the lists are equal.
func (list *SlimList) Equals(other *SlimList) bool {
	if list.Length() != other.Length() {
		return false
	}
	for i, value := range *list {
		if value != (*other)[i] {
			return false
		}
	}
	return true
}

// Length returns the number of items in the list.
func (list *SlimList) Length() int {
	if list == nil {
		return 0
	}
	return len(*list)
}

// StringAt returns the element at the index converted to string.
func (list *SlimList) StringAt(index int) string {
	return list.ElementAt(index).(string)
}

// TailAt returns a list containing everything from the indes to the end.
func (list *SlimList) TailAt(index int) *SlimList {
	tail := (*list)[index:]
	return &tail
}

// ToString returns a string representation of the list.
func (list *SlimList) ToString() string {
	result := "["

	for i, entry := range *list {
		if IsSlimList(entry) {
			result += entry.(*SlimList).ToString()
		} else {
			result += entry.(string)
		}

		if i != len(*list)-1 {
			result += ", "
		}
	}
	result += "]"
	return result
}

// TransformCallResult converts the result of a call to a string representation, or an object pointer.
func TransformCallResult(callOutput []reflect.Value) SlimEntity {
	count := len(callOutput)
	if count == 0 {
		return slimprotocol.Void()
	}
	if count == 1 {
		return valueToSlimEntity(callOutput[0])
	}
	resultList := new(SlimList)
	for result := 0; result < count; result++ {
		resultList.Append(valueToSlimEntity(callOutput[result]))
	}
	return resultList
}

func valueToSlimEntity(inputValue reflect.Value) SlimEntity {
	if !inputValue.IsValid() {
		return slimprotocol.Null()
	}
	// For predefined types, use fmt.Sprintf
	if isPredefinedType(inputValue.Type()) {
		return fmt.Sprintf("%v", inputValue.Interface())
	}
	if isObjectValue(inputValue) {
		return inputValue.Interface()
	}

	switch inputValue.Kind() {
	case reflect.Ptr, reflect.Interface:
		// This is a non-object pointer. Resolve the element
		return valueToSlimEntity(inputValue.Elem())
	// Unravel arrays and slices
	case reflect.Array, reflect.Slice:
		result := NewSlimList()
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
