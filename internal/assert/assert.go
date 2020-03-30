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

package assert

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func prefix(pc uintptr, file string, line int, ok bool) string {
	return fmt.Sprintf("%v:%v", filepath.Base(file), line)
}

// Equals returns whether expected and actual have the same types and values.
func Equals(t *testing.T, expected interface{}, actual interface{}, description string) {
	expectedValue := realNil(expected)
	actualValue := realNil(actual)
	if expectedValue != actualValue {
		var expectedType, actualType reflect.Type
		if actualValue != nil {
			actualType = reflect.ValueOf(actual).Type()
		}
		if expected != nil {
			expectedType = reflect.ValueOf(expected).Type()
		}
		prefix := prefix(runtime.Caller(1)) + fmt.Sprintf(" (%v): expected ", description)
		if expected == nil {
			t.Fatalf("%v nil but got '%v' (type %v)", prefix, actual, actualType)
		}
		if actual == nil {
			t.Fatalf("%v '%v' (type %v) but got nil", prefix, expected, expectedType)
		}
		t.Fatalf("%v '%v' (type %v) but got '%v' (type %v)", prefix, expected, expectedType, actual, actualType)
	}
}

// IsTrue returhs whether actual is true.
func IsTrue(t *testing.T, actual bool, description string) {
	if !actual {
		t.Fatalf("%v (%v): not true", prefix(runtime.Caller(1)), description)
	}
}

// Panics asserts whether a function panics.
func Panics(t *testing.T, testFunction func(), expectedPanicMessage string, description string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%v: No panic", description)
		} else {
			Equals(t, expectedPanicMessage, r, fmt.Sprintf("%v (%v) Panic message", prefix(runtime.Caller(1)), description))
		}
	}()

	testFunction()
}

// Resolve the trickyness that values can be typed nil, which is not equal to nil
func realNil(input interface{}) interface{} {
	if input == nil {
		return nil
	}
	value := reflect.ValueOf(input)
	if value.Type().Kind() == reflect.Ptr && value.IsNil() {
		return nil
	}
	return input
}
