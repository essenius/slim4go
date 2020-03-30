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

	"github.com/essenius/slim4go/internal/fixture"
	"github.com/essenius/slim4go/internal/standardlibrary"

	"github.com/essenius/slim4go/internal/assert"
	"github.com/essenius/slim4go/internal/slimentity"
)

const instanceName = "scriptTableActor"

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

func (messenger *Messenger) Panic() {
	panic(messenger.Message())
}

func initProcessorAndLibrary(t *testing.T) (*SlimStatementProcessor, *standardlibrary.StandardLibrary) {
	symbols := NewSymbolTable()
	parser := NewParser(symbols)
	objectHandler := NewObjectHandler(parser)
	parser.SetObjectSerializer(objectHandler)
	processor := NewStatementProcessor(fixture.NewRegistry(), objectHandler, parser, symbols)

	objectHandler.Add("libraryStandard", standardlibrary.New(standardlibrary.NewActorStack(), objectHandler))
	assert.Equals(t, 1, objectHandler.Length(), "Length of object collection = 1 (libraryStandard)")
	library := objectHandler.Get("libraryStandard").(*standardlibrary.StandardLibrary)
	assert.IsTrue(t, library != nil, "library found")

	processor.registry.AddFixture(NewMessenger)
	assert.Equals(t, 1, processor.registry.Length(), "Length of registry = 1 (Messenger)")
	processor.registry.AddNamespace("slimprocessor")

	assert.Equals(t, "OK", processor.DoMake(instanceName, "slimprocessor.Messenger", slimentity.NewSlimList()), "Make Messenger in initProcessorAndLibrary")
	assert.Equals(t, "/__VOID__/", processor.DoCall(instanceName, "SetMessage", slimentity.NewSlimListContaining([]slimentity.SlimEntity{"Hello world"})), "Call Set in initProcessorAndLibrary")
	return processor, library
}

func TestStatememtProcessorStandardLibrary(t *testing.T) {
	processor, library := initProcessorAndLibrary(t)

	fixture1 := library.GetFixture()
	processor.SetSymbol("fixture1", fixture1)
	assert.Equals(t, "*slimprocessor.Messenger", reflect.TypeOf(fixture1).String(), "Fixture type Messenger OK")
	assert.Equals(t, "/__VOID__/", processor.DoCall(instanceName, "SetMessage", slimentity.NewSlimListContaining([]slimentity.SlimEntity{"Hello world"})), "Call Set before push")
	library.PushFixture()
	assert.Equals(t, "/__VOID__/", processor.DoCall(instanceName, "SetMessage",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"Bye Bye"})), "Call Set after push")
	assert.Equals(t, "Bye Bye", processor.DoCall(instanceName, "Message", slimentity.NewSlimList()), "Call Get before pop")
	library.PopFixture()
	assert.Equals(t, "Hello world", processor.DoCall(instanceName, "Message", slimentity.NewSlimList()), "Call Get after pop")
	assert.Equals(t, "echo", library.Echo("echo"), "Echo")
	assert.Equals(t, "OK", processor.DoMake(instanceName, "Messenger", slimentity.NewSlimList()), "Make Messenger before making $fixture1")
	assert.Equals(t, "", processor.DoCall(instanceName, "Message", slimentity.NewSlimList()), "Check value before making $fixture1")
	assert.Equals(t, "OK", processor.DoMake(instanceName, "$fixture1", slimentity.NewSlimList()), "Make $fixture1")
	assert.Equals(t, "Hello world", processor.DoCall(instanceName, "Message", slimentity.NewSlimList()), "Call Get after making $fixture1")
	assert.Equals(t, "__EXCEPTION__:message:<<Actor stack empty>>", library.PopFixture(), "Pop fixture on empty stack")
}
func TestStatementProcessorMakeMessenger(t *testing.T) {
	processor, _ := initProcessorAndLibrary(t)
	processor.SetSymbol("test1", "TestResponse")
	assert.Equals(t, "TestResponse", processor.DoCall("instance1", "CloneSymbol",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"$test1"})), "Call cloneSymbol without creating an instance first")

	processor.registry.AddFixture(NewMessenger)
	processor.DoImport("slimprocessor")
	assert.Equals(t, "OK", processor.DoMake("instance1", "Messenger", slimentity.NewSlimList()), "Make")
	assert.Equals(t, "/__VOID__/", processor.DoCall("instance1", "SetMessage",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"Hello world"})), "Call Set Message (method)")
	assert.Equals(t, "Hello world", processor.DoCall("instance1", "Message", slimentity.NewSlimList()), "Call Message (method)")
	assert.Equals(t, "Hello world", processor.DoCall("instance1", "GetMessageField", slimentity.NewSlimList()), "Call Get Message Field (field)")
	assert.Equals(t, "Hello world", processor.DoCall("instance1", "MessageField", slimentity.NewSlimList()), "Call Message Field (field)")
	assert.Equals(t, "/__VOID__/", processor.DoCall("instance1", "SetMessageField",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"Goodbye"})), "Call Set Message Field (field)")
	assert.Equals(t, "Goodbye", processor.DoCall("instance1", "Message", slimentity.NewSlimList()), "Call Message (method)")

	processor.SetSymbol("fixture", "Messenger")
	assert.Equals(t, "OK", processor.DoMake("instance1", "$fixture", slimentity.NewSlimList()),
		"Remake an existing instance overwrites it without error. It uses a string symbol as fixture name")
	assert.Equals(t, "", processor.DoCall("instance1", "Message", slimentity.NewSlimList()), "Call Get after creating new instance1")
	processor.SetSymbol("message", "Bye bye")
	processor.SetSymbol("method", "SetMessage")
	assert.Equals(t, "/__VOID__/", processor.DoCall("instance1", "SetMessage",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"$message"})), "Call Set with symbol in args")
	assert.Equals(t, "Bye bye", processor.DoCall("instance1", "Message", slimentity.NewSlimList()), "Call Get after setting with symbols")
	assert.Equals(t, "__EXCEPTION__:message:<<Panic: Bye bye>>",
		processor.DoCall("instance1", "Panic", slimentity.NewSlimList()), "Panic is caught and reported")
	assert.Equals(t, "__EXCEPTION__:message:<<Expected 1 parameter(s) but got 0>>",
		processor.DoCall("instance1", "SetMessage", slimentity.NewSlimList()), "Call Set with empty parameter set")
	assert.Equals(t, "SetMessage", processor.DoCall("instance1", "CloneSymbol",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"$method"})), "Call cloneSymbol on instance1")
	assert.Equals(t, "__EXCEPTION__:message:<<COULD_NOT_INVOKE_CONSTRUCTOR Messenger:Expected_0_parameter(s)_but_got_1>>",
		processor.DoMake("wronginstance", "Messenger", slimentity.NewSlimListContaining([]slimentity.SlimEntity{"5"})),
		"wrong number of parameters for constructor")
}

func TestStatementProcessorMakeOrder(t *testing.T) {
	processor, _ := initProcessorAndLibrary(t)
	processor.registry.AddFixture(NewOrder)
	processor.DoImport("fixture")
	assert.Equals(t, "OK", processor.DoMake("instance1", "Order", slimentity.NewSlimList()), "Make Order")
	assert.Equals(t, "/__VOID__/", processor.DoCall("instance1", "SetProduct",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"cup", "0.50"})), "Call SetProduct")
	assert.Equals(t, "/__VOID__/", processor.DoCall("instance1", "SetUnits",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"200"})), "Call SetUnits")
	assert.Equals(t, "100", processor.DoCall("instance1", "Price", slimentity.NewSlimList()), "Call Price")
	assert.Equals(t, "__EXCEPTION__:message:<<NO_CLASS nonexisting>>",
		processor.DoMake("instance2", "nonexisting", slimentity.NewSlimList()), "Make a nonexisting fixture")
	assert.Equals(t, "__EXCEPTION__:message:<<NO_INSTANCE nonexisting>>",
		processor.DoCall("nonexisting", "Price", slimentity.NewSlimList()), "Price on nonexisting instance")
	assert.Equals(t, "__EXCEPTION__:message:<<NO_METHOD_IN_CLASS Nonexisting[0] *slimprocessor.Order>>",
		processor.DoCall("instance1", "Nonexisting", slimentity.NewSlimList()), "Nonexisting method on existing instance")
	assert.Equals(t, "__EXCEPTION__:message:<<COULD_NOT_INVOKE_CONSTRUCTOR Order:Expected_0_parameter(s)_but_got_1>>",
		processor.DoMake("instance3", "Order", slimentity.NewSlimListContaining([]slimentity.SlimEntity{"entry"})),
		"Use a constructor with wrong number of parameters")
}

func TestStatementProcessorMakeObjectWithPanic(t *testing.T) {
	processor, _ := initProcessorAndLibrary(t)
	processor.registry.AddFixture(NewObjectWithPanic)
	assert.Equals(t, "__EXCEPTION__:message:<<COULD_NOT_INVOKE_CONSTRUCTOR int:Panic:_Object_creation_failed>>",
		processor.DoMake("instance1", "int", slimentity.NewSlimList()), "Make Object With Panic")
}

func TestStatementProcessorSerializeObjectsIn(t *testing.T) {
	processor, _ := initProcessorAndLibrary(t)
	test1 := "test1"
	assert.Equals(t, test1, processor.SerializeObjectsIn(test1), "string")
	assert.IsTrue(t, !slimentity.IsObject(test1), "test1 is no object")
	aDemoStruct1 := &demoStruct1{"demo"}
	aDemoStruct1.Parse("demo1")
	assert.IsTrue(t, slimentity.IsObject(aDemoStruct1), "aDemoStruct1 is an object")
	assert.Equals(t, "demo1", processor.SerializeObjectsIn(aDemoStruct1), "*struct with *ToString")
	aDemoStruct2 := demoStruct2{"demo2"}
	assert.IsTrue(t, slimentity.IsObject(aDemoStruct2), "aDemoStruct2 is an object")
	assert.Equals(t, "demo2", processor.SerializeObjectsIn(aDemoStruct2), "struct with ToString")
	ptrToADemoStruct2 := &aDemoStruct2
	assert.IsTrue(t, slimentity.IsObject(ptrToADemoStruct2), "*aDemoStruct2 is an object")
	assert.Equals(t, "demo2", processor.SerializeObjectsIn(ptrToADemoStruct2), "*struct with ToString")
	anEmptyStruct := emptyStruct{}
	assert.Equals(t, "slimprocessor.emptyStruct", processor.SerializeObjectsIn(anEmptyStruct), "struct without ToString")
	list := slimentity.NewSlimListContaining([]slimentity.SlimEntity{"test2", aDemoStruct1, aDemoStruct2, anEmptyStruct})
	assert.IsTrue(t, !slimentity.IsObject(list), "list is no object")
	assert.Equals(t, "[test2, demo1, demo2, slimprocessor.emptyStruct]",
		processor.SerializeObjectsIn(list).(*slimentity.SlimList).ToString(), "list with objects")
}
