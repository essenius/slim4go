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

// Counter is a simple demo fixture for Slim4Go
type Counter struct {
	count int64
}

// NewCounter is the constructor
func NewCounter(args ...int64) *Counter {
	counter := new(Counter)
	if len(args) > 0 {
		counter.count = args[0]
	}
	return counter
}

// SetCount sets the counter to a specific value
func (counter *Counter) SetCount(count int64) {
	counter.count = count
}

// CountUp increments the counter
func (counter *Counter) CountUp() {
	counter.count++
}

// Value returns the current counter value
func (counter *Counter) Value() int64 {
	return counter.count
}
