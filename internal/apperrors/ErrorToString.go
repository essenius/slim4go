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

package apperrors

import (
	"fmt"
)

// ErrorToString converts an error response (e.g. from a panic) to a string.
func ErrorToString(message interface{}) string {
	if err, ok := message.(error); ok {
		return err.Error()
	}
	return fmt.Sprintf("%v", message)
}
