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
type FixtureMapping struct {
	PublicProperty  []string
	privateProperty string
}

// NewFixtureMapping is the constructor for Waiter.
func NewFixtureMapping() *FixtureMapping {
	return new(FixtureMapping)
}

// PublicMethod1 sets PublicProperty.
func (mapping *FixtureMapping) PublicMethod1(input []string) {
	mapping.PublicProperty = input
}

// PublicMethod2 gets PublicProperty.
func (mapping *FixtureMapping) PublicMethod2() []string {
	return mapping.PublicProperty
}

// PublicMethod3 sets PrivateProperty.
func (mapping *FixtureMapping) PublicMethod3(input string) {
	mapping.privateProperty = input
}
