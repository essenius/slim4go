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

import "time"

// Waiter introduces a delay to simulate long running activities.
type Waiter struct {
}

// NewWaiter is the constructor for Waiter.
func NewWaiter() *Waiter {
	return new(Waiter)
}

// Wait does the actual waiting.
func (waiter *Waiter) Wait(delay int64) {
	time.Sleep(time.Duration(delay) * time.Second)
}
