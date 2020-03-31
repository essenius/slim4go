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
	"fmt"
	"strings"
	"testing"

	"github.com/essenius/slim4go/internal/assert"
)

func TestSlimMarshallerReadRequest(t *testing.T) {
	var testcases = []struct {
		request     string
		listLength  int
		firstValue  string
		roundtrip   string
		description string
	}{
		{"000000:", 0, "", "000000:", "Empty message"},
		{"000003:bye", 0, "bye", "000003:bye", "Bye message"},
		{"000026:[000001:000009:Hi there.:]", 1, "Hi there.", "000026:[000001:000009:Hi there.:]", "message with ASCII only"},
		{"000029:[000001:000012:Hi JRÜ€©:]", 1, "Hi JRÜ€©", "000029:[000001:000012:Hi JRÜ€©:]", "message with multi-byte characters"},
		{"000022:[0000001:00000002:Hi:]", 1, "Hi", "000019:[000001:000002:Hi:]", "message with lengths longer than 6 digits"},
		{"000017:[000001:000000::]", 1, "", "000017:[000001:000000::]", "List with empty entry"},
		{"000027:[000001:000010:[[a, b, c]:]", 1, "[[a, b, c]", "000027:[000001:000010:[[a, b, c]:]", "message with table spec in string"},
	}
	for _, testcase := range testcases {
		stringReader := strings.NewReader(testcase.request)
		output, err := ReadRequest(stringReader)
		assert.Equals(t, nil, err, fmt.Sprintf("no error for %v", testcase.description))
		if testcase.listLength == 0 {
			assert.Equals(t, testcase.firstValue, output.(string), testcase.description)
		} else {
			assert.Equals(t, testcase.listLength, output.(*SlimList).Length(), fmt.Sprintf("instruction count for %v", testcase.description))
			assert.Equals(t, testcase.firstValue, (output.(*SlimList)).ElementAt(0).(string), fmt.Sprintf("First line for %v", testcase.description))
		}
		assert.Equals(t, testcase.roundtrip, Marshal(output), fmt.Sprintf("Round trip works for %v", testcase.description))
	}
}

func TestSlimMarshallerReadRequestFail(t *testing.T) {
	var testcases = []struct {
		request      string
		errorMessage string
		description  string
	}{
		{"", "readLength: Could not find next delimiter ':' (EOF)", "Empty message"},
		{"00a:", "readLength: Could not interpret length '00a'", "Wrong length spec"},
		{"0000017:[000001:", "readExactBytes: Expected 17 bytes from Slim client, but got 8", "Incomplete message"},
		{"000026:[000001:000009:Hi there.:q", "skipByte: Expected ']' but found 'q'", "Wrong list delimiter"},
		{"000026:[000001:000009:Hi there.:", "SkipByte: No input available", "missing final delimiter"},
		{"", "readLength: Could not find next delimiter ':' (EOF)", "Empty message"},
	}
	for _, testcase := range testcases {
		stringReader := strings.NewReader(testcase.request)
		_, err := ReadRequest(stringReader)
		assert.Equals(t, testcase.errorMessage, err.Error(), testcase.description)
	}
}

func TestSlimMarshallerReadRequestMultiple(t *testing.T) {
	const slimLineIn = "000037:[000002:000009:Hi there.:000003:Bye:]"
	stringReader := strings.NewReader(slimLineIn)
	entity, err := ReadRequest(stringReader)
	list := entity.(*SlimList)
	assert.Equals(t, nil, err, "no error")
	assert.Equals(t, 2, list.Length(), "length == 2")
	assert.Equals(t, "Hi there.", list.ElementAt(0).(string), "list[0] matches")
	assert.Equals(t, "Bye", list.ElementAt(1).(string), "list[1] matches")
	slimLineOut := Marshal(list)
	assert.Equals(t, slimLineIn, slimLineOut, "Round trip works")
}

func TestSlimMarshallerReadRequestRecursive(t *testing.T) {
	const slimLineIn = "000106:[000001:000089:[000004:000017:decisionTable_0_0:000004:make:000015:decisionTable_0:000012:Hi_JRÜ€©:]:]"
	stringReader := strings.NewReader(slimLineIn)
	entity, err := ReadRequest(stringReader)
	list := entity.(*SlimList)
	assert.Equals(t, nil, err, "no error")
	assert.Equals(t, 1, list.Length(), "list length 1")
	subList := list.ElementAt(0).(*SlimList)
	assert.Equals(t, 4, subList.Length(), "Sublist length 4")
	assert.Equals(t, "decisionTable_0_0", subList.ElementAt(0), "First element")
	assert.Equals(t, "make", subList.ElementAt(1), "Second element")
	assert.Equals(t, "decisionTable_0", subList.ElementAt(2), "Third element")
	assert.Equals(t, "Hi_JRÜ€©", subList.ElementAt(3), "Fourth element")
}

func TestSlimMarshallerConvertNull(t *testing.T) {
	assert.Equals(t, "null", convertNull(nil), "nil")
	assert.Equals(t, "a", convertNull("a"), "non-null")
}
