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

package standardlibrary

import (
	"testing"

	"github.com/essenius/slim4go/internal/assert"
)

type MockCollection struct{}

func (collection *MockCollection) Add(name string, value interface{}) {
}

func (collection *MockCollection) Get(name string) interface{} {
	return name
}

func (collection *MockCollection) Set(name string, value interface{}) error {
	return nil
}

func (collection *MockCollection) Length() int {
	return 1
}

const instanceName = "scriptTableActor"

func TestStandardLibraryStack(t *testing.T) {
	library := New(NewActorStack(), new(MockCollection))
	assert.Equals(t, scriptTableActorName, library.GetFixture(), "GetFixture returns right result")
	assert.Equals(t, 0, library.actors.Length(), "Initial actors length == 0")
	assert.Equals(t, nil, library.PushFixture(), "Push fixture succeeds")
	assert.Equals(t, 1, library.actors.Length(), "Initial actors length == 1 after push")
	assert.Equals(t, nil, library.PopFixture(), "Pop fixture succeeds")
	assert.Equals(t, 0, library.actors.Length(), "Initial actors length == 0 after pop")
	assert.Equals(t, "__EXCEPTION__:message:<<Actor stack empty>>", library.PopFixture(), "Pop fixture fails on empty stack")
}

func TestSlimLibaryCloneSymbol(t *testing.T) {
	library := New(NewActorStack(), nil)
	assert.Equals(t, "clone", library.CloneSymbol("clone"), "clone")
	// A clone of a symbol should really be a clone, not a pointer to the same instance.
	s := "string in variable"
	p1 := &s
	p2 := library.CloneSymbol(p1).(*string)
	assert.Equals(t, "string in variable", *p2, "Clone string in pointer to variable")
	*p2 = "Changed"
	assert.Equals(t, "string in variable", *p1, "Original pointer element didn't change")
	assert.Equals(t, "Changed", *p2, "Clone pointer element did change")
	assert.Equals(t, "echo", library.Echo("echo"), "Echo")
}
