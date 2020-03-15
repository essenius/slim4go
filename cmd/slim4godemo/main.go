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

package main

import (
	"github.com/essenius/slim4go"
	"github.com/essenius/slim4go/examples/demofixtures"
)

func main() {
	slim4go.RegisterFixture("Array", demofixtures.NewArray)
	slim4go.RegisterFixture("Counter", demofixtures.NewCounter)
	slim4go.RegisterFixture("Dictionary", demofixtures.NewDictionary)
	slim4go.RegisterFixture("FibonacciFixture", demofixtures.NewFibonacciFixture)
	slim4go.RegisterFixture("MemoObject", demofixtures.NewMemoObject)
	slim4go.RegisterFixture("TableFixture", demofixtures.NewTableFixture)
	slim4go.RegisterFixture("Temperature", demofixtures.NewTemperature)
	slim4go.RegisterFixture("TemperatureConverter", demofixtures.NewTemperatureConverter)
	slim4go.RegisterFixture("TestQuery", demofixtures.NewTestQuery)
	slim4go.RegisterFixture("Waiter", demofixtures.NewWaiter)
	slim4go.Serve()
}
