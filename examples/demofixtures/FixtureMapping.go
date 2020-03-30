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

// FixtureMapping shows how FitNesse Call functions can be mapped to members (fields/methods).
// Only public fields or methods can be used

// FixtureMapping is used to show how to map fields and methods to FitNesse script statements.
type FixtureMapping struct {
	Field []string
}

// NewFixtureMapping is the constructor for Waiter.
func NewFixtureMapping() *FixtureMapping {
	return new(FixtureMapping)
}

// Method1 sets Field.
func (mapping *FixtureMapping) Method1(input []string) {
	mapping.Field = input
}

// Method2 gets Field.
func (mapping *FixtureMapping) Method2() []string {
	return mapping.Field
}
