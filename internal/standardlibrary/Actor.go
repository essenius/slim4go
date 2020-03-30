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

// Definitions and constructors

// ActorStack is a stack of actors.
type ActorStack []interface{}

// NewActorStack returns a new ActorStack.
func NewActorStack() *ActorStack {
	stack := make(ActorStack, 0)
	return &stack
}

// Methods

// Head returns the top of the stack without removing it.
func (actors *ActorStack) Head() interface{} {
	length := len(*actors)
	if length == 0 {
		return nil
	}
	return (*actors)[0]
}

// Length returns the number of items in the stack.
func (actors *ActorStack) Length() int {
	return len(*actors)
}

// Pop returns the top item and removes it from the stack.
func (actors *ActorStack) Pop() interface{} {
	anActor := actors.Head()
	if anActor == nil {
		return nil
	}
	*actors = (*actors)[1:]
	return anActor
}

// Push pushes a new item onto the stack
func (actors *ActorStack) Push(instance interface{}) {
	*actors = append(ActorStack{instance}, *actors...)
}
