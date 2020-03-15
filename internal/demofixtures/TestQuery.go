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

// TestQuery shows how to implement a FitNesse Query fixture.
type TestQuery struct {
	max int
}

// NewTestQuery is the constructor for TestQuery.
func NewTestQuery(max int) *TestQuery {
	testQuery := new(TestQuery)
	testQuery.max = max
	return testQuery
}

// Query fulfils the FitNesse query interface.
func (testQuery TestQuery) Query() [][][]interface{} {
	rowList := [][][]interface{}{}
	for i := 1; i <= testQuery.max; i++ {
		entry1 := []interface{}{"n", i}
		entry2 := []interface{}{"2n", 2 * i}
		row := [][]interface{}{entry1, entry2}
		rowList = append(rowList, row)
	}
	return rowList
}
