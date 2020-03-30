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

func TestSlimActorStack(t *testing.T) {
	actors := new(ActorStack)
	assert.Equals(t, 0, actors.Length(), "Length is 0 at start")
	assert.Equals(t, nil, actors.Head(), "Head is nil at start")
	assert.Equals(t, nil, actors.Pop(), "Pop is nil at start")
	actors.Push("one")
	assert.Equals(t, 1, actors.Length(), "Length is 1 after first push")
	assert.Equals(t, "one", actors.Head(), "content of first push OK")
	actors.Push("two")
	assert.Equals(t, 2, actors.Length(), "Length is 2 after second push")
	assert.Equals(t, "two", actors.Head(), "name of second push OK")
	assert.Equals(t, "one", (*actors)[1], "name of first push OK")
	actor3 := actors.Pop()
	assert.Equals(t, 1, actors.Length(), "Length is 1 after first pop")
	assert.Equals(t, "one", actors.Head(), "first push now on top")
	assert.Equals(t, "two", actor3, "Popped entry is correct")
	actor4 := actors.Pop()
	assert.Equals(t, 0, actors.Length(), "Length is 0 after second pop")
	assert.Equals(t, "one", actor4, "Second popped entry is correct")
}
