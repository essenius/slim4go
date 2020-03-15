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
	"reflect"
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
func NewSlimListContaining(listToAdd []SlimEntity) *SlimList {
	list := NewSlimList()
	for _, entry := range listToAdd {
		list.Append(entry)
	}
	return list
}

// Functions

// IsSlimList checks whether entity is a SlimList.
func IsSlimList(entity SlimEntity) bool {
	slimListType := reflect.PtrTo(reflect.TypeOf((*SlimList)(nil)).Elem())
	return reflect.TypeOf(entity) == slimListType
}

// ToSlice converts a list to a slice of strings. Only works for simple lists.
// Nested lists will be serialized with brackets and commas.
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

// Pop removes the first item from the list.
func (list *SlimList) Pop() interface{} {
	if len(*list) == 0 {
		panic("Can't pop from empty list")
	}
	returnValue := (*list)[0]
	*list = (*list)[1:]
	return returnValue
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

// TypeOfElementAt returns the type of the element at position index.
func (list *SlimList) TypeOfElementAt(index int) string {
	return reflect.TypeOf(list.ElementAt(index)).Name()
}
