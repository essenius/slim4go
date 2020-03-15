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

package demofixtures

// Array is a simple demo fixture for Slim4Go.
type Array struct {
	intArray     []int
	stringMatrix [][]string
}

// NewArray is the constructor for Array
func NewArray() *Array {
	return new(Array)
}

// SetIntArray sets a one dimensional int array.
func (array *Array) SetIntArray(input []int) {
	array.intArray = input
}

// IntArray gets a one dimensional int array.
func (array *Array) IntArray() []int {
	return array.intArray
}

// SetStringMatrix sets a two dimensional string array.
func (array *Array) SetStringMatrix(input [][]string) {
	array.stringMatrix = input
}

// StringMatrix gets a two dimensional string array.
func (array *Array) StringMatrix() [][]string {
	return array.stringMatrix
}
