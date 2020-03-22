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
	"testing"

	"github.com/essenius/slim4go/internal/assert"
	"github.com/essenius/slim4go/internal/slimentity"
)

type Messenger struct {
	MessageField string
}

func NewMessenger() *Messenger {
	return new(Messenger)
}

func (messenger *Messenger) SetMessage(message string) {
	messenger.MessageField = message
}

func (messenger *Messenger) Message() string {
	return messenger.MessageField
}

const instanceName = "scriptTableActor"

func initProcessorAndLibrary(t *testing.T) (*slimStatementProcessor, *StandardLibrary) {
	// this is normally a single instance but we want to start fresh during testing
	objectCollectionInstance = nil
	processor := injectStatementProcessor().(*slimStatementProcessor)
	processor.fixtureRegistry().AddFixture(NewMessenger)
	processor.fixtureRegistry().AddNamespace("slimprocessor")
	assert.Equals(t, 1, processor.objects().Length(), "Initial Length of stack = 1 (libraryStandard")
	library := processor.objects().objectNamed("libraryStandard").instance().(*StandardLibrary)
	assert.Equals(t, "OK", processor.doMake(instanceName, "slimprocessor.Messenger", slimentity.NewSlimList()), "Make Messenger in initProcessorAndLibrary")
	assert.Equals(t, "/__VOID__/", processor.doCall(instanceName, "SetMessage", slimentity.NewSlimListContaining([]slimentity.SlimEntity{"Hello world"})), "Call Setin initProcessorAndLibrary")
	return processor, library
}

func TestStandardLibraryStack(t *testing.T) {

	processor, library := initProcessorAndLibrary(t)
	assert.IsTrue(t, library != nil, "library found")
	fixture1 := library.GetFixture()
	processor.setSymbol("fixture1", fixture1)
	assert.Equals(t, "*slimprocessor.Messenger", reflect.TypeOf(fixture1).String(), "Fixture type Messenger OK")
	assert.Equals(t, 0, library.actors.Length(), "Length of stack = 0")
	assert.Equals(t, "/__VOID__/", processor.doCall(instanceName, "SetMessage",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"Hello world"})), "Call Set before push")
	library.PushFixture()
	assert.Equals(t, 1, library.actors.Length(), "Length of stack = 1 after push")
	assert.Equals(t, "/__VOID__/", processor.doCall(instanceName, "SetMessage",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"Bye Bye"})), "Call Set after push")
	assert.Equals(t, "Bye Bye", processor.doCall(instanceName, "Message", slimentity.NewSlimList()), "Call Get before pop")
	library.PopFixture()
	assert.Equals(t, 0, library.actors.Length(), "Length of stack = 0 after pop")
	assert.Equals(t, "Hello world", processor.doCall(instanceName, "Message", slimentity.NewSlimList()), "Call Get after pop")
	assert.Equals(t, "echo", library.Echo("echo"), "Echo")
	assert.Equals(t, "OK", processor.doMake(instanceName, "Messenger", slimentity.NewSlimList()), "Make Messenger before making $fixture1")
	assert.Equals(t, "", processor.doCall(instanceName, "Message", slimentity.NewSlimList()), "Check value before making $fixture1")
	assert.Equals(t, "OK", processor.doMake(instanceName, "$fixture1", slimentity.NewSlimList()), "Make $fixture1")
	assert.Equals(t, "Hello world", processor.doCall(instanceName, "Message", slimentity.NewSlimList()), "Call Get after making $fixture1")
	assert.Equals(t, "__EXCEPTION__:message:<<Actor stack empty>>", library.PopFixture(), "Pop fixture on empty stack")
}

func TestSlimLibaryCloneSymbol(t *testing.T) {
	processor, library := initProcessorAndLibrary(t)
	assert.Equals(t, "clone", library.CloneSymbol("clone"), "clone")
	// A clone of a symbol should really be a clone, not a pointer to the same instance.
	s := "string in variable"
	p1 := &s
	p2 := library.CloneSymbol(p1).(*string)
	assert.Equals(t, "string in variable", *p2, "Clone string in pointer to variable")
	*p2 = "Changed"
	assert.Equals(t, "string in variable", *p1, "Original pointer element didn't change")
	assert.Equals(t, "Changed", *p2, "Clone pointer element did change")
	// Same thing if the clone is an object
	processor.setSymbol("testSymbol", library.GetFixture())
	clonedItem := processor.doCall(instanceName, "CloneSymbol", slimentity.NewSlimListContaining([]slimentity.SlimEntity{"$testSymbol"}))
	clonedMessage := clonedItem.(*Messenger)
	assert.Equals(t, "Hello world", clonedMessage.Message(), "cloned instnace has original instance's content")
	clonedMessage.SetMessage("Goodbye")
	assert.Equals(t, "Hello world", processor.doCall(instanceName, "Message", slimentity.NewSlimList()), "Original instance should not be changed")
}
