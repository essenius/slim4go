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

type demoFactory struct{}

func main() {
	slim4go.RegisterFixturesFrom(demofixtures.NewTemperatureFactory())
	slim4go.RegisterFixture(demofixtures.NewArray)
	slim4go.RegisterFixture(demofixtures.NewCounter)
	slim4go.RegisterFixture(demofixtures.NewDictionary)
	slim4go.RegisterFixture(demofixtures.NewFibonacciFixture)
	slim4go.RegisterFixture(demofixtures.NewFixtureMapping)
	slim4go.RegisterFixture(demofixtures.NewMemoObject)
	slim4go.RegisterFixture(demofixtures.NewTableFixture)
	slim4go.RegisterFixture(demofixtures.NewTestQuery)
	slim4go.RegisterFixture(demofixtures.NewWaiter)
	slim4go.Serve()
}
