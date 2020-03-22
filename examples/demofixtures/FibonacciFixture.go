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

import (
	"fmt"

	"github.com/essenius/slim4go/internal/slimlog"

	"github.com/essenius/slim4go/examples/demosut"
)

// FibonacciFixture is an example how to create a fixture for a decision table.
type FibonacciFixture struct {
	InputValue int64
	err        error
	result     int64
}

// NewFibonacciFixture is the constructor for FibonacciFixture.
func NewFibonacciFixture() *FibonacciFixture {
	return new(FibonacciFixture)
}

// Reset is called before processing a line in the decision table.
func (fixture *FibonacciFixture) Reset() {
	fixture.err = nil
	fixture.result = 0
}

// Execute is called after all values have been set, and before getting results.
func (fixture *FibonacciFixture) Execute() {
	fixture.result, fixture.err = demosut.Fibonacci(fixture.InputValue)
}

// Fibonacci runs the fibonacci function from the system under test.
func (fixture *FibonacciFixture) Fibonacci() string {
	if fixture.err != nil {
		return fixture.err.Error()
	}
	return fmt.Sprintf("%v", fixture.result)
}

// Table returns the full test table
func (fixture *FibonacciFixture) Table(table [][]string) {
	slimlog.Info.Printf("Table size: %v x %v", len(table), len(table[0]))
}
