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

package slimprocessor

import (
	"reflect"

	"github.com/essenius/slim4go/internal/slimentity"
	"github.com/essenius/slim4go/internal/slimprotocol"
)

// Definitions and constructors

const scriptTableActorName = "scriptTableActor"

// StandardLibrary is the library that gets added to the library list by default.
type StandardLibrary struct {
	actors  *actorStack
	objects *objectCollection
}

// NewStandardLibrary instantiates a new StandardLibrary.
func newStandardLibrary(actors *actorStack, objects *objectCollection) *StandardLibrary {
	standardLibrary := new(StandardLibrary)
	standardLibrary.actors = actors
	standardLibrary.objects = objects
	return standardLibrary
}

// Methods

// CloneSymbol creates a clone of a symbol. If the symbol points to a struct, a copy of that struct is made.
func (standardLibrary *StandardLibrary) CloneSymbol(input interface{}) interface{} {
	inputValue := reflect.ValueOf(input)
	if inputValue.Kind() == reflect.Ptr {
		inpuElem := inputValue.Elem()
		targetValue := reflect.New(inpuElem.Type())
		targetElem := targetValue.Elem()
		targetElem.Set(inpuElem)
		return targetValue.Interface()
	}
	return inputValue.Interface()
}

// Echo returns the input.
func (standardLibrary *StandardLibrary) Echo(input interface{}) interface{} {
	return input
}

// GetFixture gets the currently executed fixture.
func (standardLibrary *StandardLibrary) GetFixture() interface{} {
	scriptTableActor := standardLibrary.objects.objectNamed(scriptTableActorName)
	return (*scriptTableActor).instance()
}

// PopFixture pops a fixture from the stack.
func (standardLibrary *StandardLibrary) PopFixture() slimentity.SlimEntity {
	fixture := standardLibrary.actors.Pop()
	if fixture != nil {
		return standardLibrary.objects.setObjectInstance(scriptTableActorName, fixture)
	}
	return slimprotocol.Exception("Actor stack empty")
}

// PushFixture pushes a fixture on the stack.
func (standardLibrary *StandardLibrary) PushFixture() slimentity.SlimEntity {
	currentFixture := standardLibrary.GetFixture()
	newFixture := standardLibrary.CloneSymbol(currentFixture)
	standardLibrary.actors.Push(currentFixture)
	return standardLibrary.objects.setObjectInstance(scriptTableActorName, newFixture)
}
