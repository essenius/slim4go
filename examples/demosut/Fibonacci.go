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

package demosut

import "fmt"

// Fibonacci calculates a Fibonacci number.
func Fibonacci(factor int64) (int64, error) {
	if factor < 0 {
		return 0, fmt.Errorf("Input can't be negative")
	}
	if factor < 2 {
		return factor, nil
	}

	var previousNumber, currentNumber, i int64 = 1, 1, 0
	for i = 2; i < factor; i++ {
		previousNumber, currentNumber = currentNumber, previousNumber+currentNumber
		if currentNumber < 0 {
			return 0, fmt.Errorf("Overflow")
		}
	}
	return currentNumber, nil
}
