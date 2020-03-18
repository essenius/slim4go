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

// TODO: complete experiment: create a fixture factory instead of specifying all fixture constructors
// That would also enable using Import to import all fixture specs in a factory

func NewFixtureFactory() *FixtureFactory {
	return new(FixtureFactory)
}

type FixtureFactory struct{}

func (factory *FixtureFactory) NewOrder(productID string, unitPrice float64, units int) *Order {
	order := new(Order)
	order.ProductID = productID
	order.UnitPrice = unitPrice
	order.Units = units
	return order
}

func (factory *FixtureFactory) NewMessenger() interface{} {
	return new(Messenger)
}

// end experiment

type Order struct {
	ProductID string
	Units     int
	UnitPrice float64
}

func NewOrder() *Order {
	return new(Order)
}

func (order *Order) Price() float64 {
	return order.UnitPrice * float64(order.Units)
}

func (order *Order) SetUnits(units int) {
	order.Units = units
}

func (order *Order) SetProduct(productID string, unitPrice float64) {
	order.ProductID = productID
	order.UnitPrice = unitPrice
}

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

func NewObjectWithPanic() interface{} {
	panic("Object creation failed")
}

func TestStatementProcessorMakeMessenger(t *testing.T) {
	processor := injectStatementProcessor()
	processor.setSymbol("test1", "TestResponse")
	assert.Equals(t, "TestResponse", processor.doCall("instance1", "CloneSymbol",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"$test1"})), "Call cloneSymbol without creating an instance first")

	processor.fixtures().RegisterFixture("Messenger", NewMessenger)
	assert.Equals(t, "OK", processor.doMake("instance1", "Messenger", slimentity.NewSlimList()), "Make")
	assert.Equals(t, "/__VOID__/", processor.doCall("instance1", "SetMessage",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"Hello world"})), "Call Set Message (method)")
	assert.Equals(t, "Hello world", processor.doCall("instance1", "Message", slimentity.NewSlimList()), "Call Message (method)")
	assert.Equals(t, "Hello world", processor.doCall("instance1", "GetMessageField", slimentity.NewSlimList()), "Call Get Message Field (field)")
	assert.Equals(t, "Hello world", processor.doCall("instance1", "MessageField", slimentity.NewSlimList()), "Call Message Field (field)")
	assert.Equals(t, "/__VOID__/", processor.doCall("instance1", "SetMessageField",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"Goodbye"})), "Call Set Message Field (field)")
	assert.Equals(t, "Goodbye", processor.doCall("instance1", "Message", slimentity.NewSlimList()), "Call Message (method)")

	processor.setSymbol("fixture", "Messenger")
	assert.Equals(t, "OK", processor.doMake("instance1", "$fixture", slimentity.NewSlimList()), "Remake an existing instance overwrites it without error. It uses a string symbol as fixture name")
	assert.Equals(t, "", processor.doCall("instance1", "Message", slimentity.NewSlimList()), "Call Get after creating new instance1")
	processor.setSymbol("message", "Bye bye")
	processor.setSymbol("method", "SetMessage")
	assert.Equals(t, "/__VOID__/", processor.doCall("instance1", "SetMessage",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"$message"})), "Call Set with symbol in args")
	assert.Equals(t, "Bye bye", processor.doCall("instance1", "Message", slimentity.NewSlimList()), "Call Get after setting with symbols")
	assert.Equals(t, "__EXCEPTION__:message:<<Panic: Bye bye>>", processor.doCall("instance1", "Panic", slimentity.NewSlimList()), "Panic is caught and reported")
	assert.Equals(t, "__EXCEPTION__:message:<<Expected 1 parameter(s) but got 0>>", processor.doCall("instance1", "SetMessage", slimentity.NewSlimList()), "Call Set with empty parameter set")
	assert.Equals(t, "SetMessage", processor.doCall("instance1", "CloneSymbol",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"$method"})), "Call cloneSymbol on instance1")

	processor.fixtures().RegisterFixture("Bogus", 1)

	assert.Equals(t, "__EXCEPTION__:message:<<COULD_NOT_INVOKE_CONSTRUCTOR Bogus:int_is_not_a_function>>",
		processor.doMake("wronginstance", "Bogus", slimentity.NewSlimList()), "Bogus constructor")
	assert.Equals(t, "__EXCEPTION__:message:<<COULD_NOT_INVOKE_CONSTRUCTOR Messenger:Expected_0_parameter(s)_but_got_1>>",
		processor.doMake("wronginstance", "Messenger", slimentity.NewSlimListContaining([]slimentity.SlimEntity{"5"})), "wrong number of parameters for constructor")
}

func TestStatementProcessorMakeOrder(t *testing.T) {
	processor := injectStatementProcessor()
	processor.fixtures().RegisterFixture("Order", NewOrder)
	assert.Equals(t, "OK", processor.doMake("instance1", "Order", slimentity.NewSlimList()), "Make Order")
	assert.Equals(t, "/__VOID__/", processor.doCall("instance1", "SetProduct",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"cup", "0.50"})), "Call SetProduct")
	assert.Equals(t, "/__VOID__/", processor.doCall("instance1", "SetUnits",
		slimentity.NewSlimListContaining([]slimentity.SlimEntity{"200"})), "Call SetUnits")
	assert.Equals(t, "100", processor.doCall("instance1", "Price", slimentity.NewSlimList()), "Call Price")
	assert.Equals(t, "__EXCEPTION__:message:<<NO_CLASS nonexisting>>", processor.doMake("instance2", "nonexisting", slimentity.NewSlimList()), "Make a nonexisting fixture")
	assert.Equals(t, "__EXCEPTION__:message:<<NO_INSTANCE nonexisting>>", processor.doCall("nonexisting", "Price", slimentity.NewSlimList()), "Price on nonexisting instance")
	assert.Equals(t, "__EXCEPTION__:message:<<NO_METHOD_IN_CLASS Nonexisting[0] *slimprocessor.Order>>", processor.doCall("instance1", "Nonexisting", slimentity.NewSlimList()), "Nonexisting method on existing instance")
	assert.Equals(t, "__EXCEPTION__:message:<<COULD_NOT_INVOKE_CONSTRUCTOR Order:Expected_0_parameter(s)_but_got_1>>",
		processor.doMake("instance3", "Order", slimentity.NewSlimListContaining([]slimentity.SlimEntity{"entry"})), "Use a constructor with wrong number of parameters")
}

func TestStatementProcessorMakeObjectWithPanic(t *testing.T) {
	processor := injectStatementProcessor()
	processor.fixtures().RegisterFixture("Panic", NewObjectWithPanic)
	assert.Equals(t, "__EXCEPTION__:message:<<COULD_NOT_INVOKE_CONSTRUCTOR Panic:Panic:_Object_creation_failed>>",
		processor.doMake("instance1", "Panic", slimentity.NewSlimList()), "Make Object With Panic")
}

func TestStatementProcessorSerializeObjectsIn(t *testing.T) {
	processor := injectStatementProcessor()
	test1 := "test1"
	assert.Equals(t, test1, processor.serializeObjectsIn(test1), "string")
	assert.IsTrue(t, !slimentity.IsObject(reflect.ValueOf(test1)), "test1 is no object")
	aDemoStruct1 := &demoStruct1{"demo"}
	aDemoStruct1.Parse("demo1")
	assert.IsTrue(t, slimentity.IsObject(reflect.ValueOf(aDemoStruct1)), "aDemoStruct1 is an object")
	assert.Equals(t, "demo1", processor.serializeObjectsIn(aDemoStruct1), "*struct with *ToString")
	aDemoStruct2 := demoStruct2{"demo2"}
	assert.IsTrue(t, slimentity.IsObject(reflect.ValueOf(aDemoStruct2)), "aDemoStruct2 is an object")
	assert.Equals(t, "demo2", processor.serializeObjectsIn(aDemoStruct2), "struct with ToString")
	ptrToADemoStruct2 := &aDemoStruct2
	assert.IsTrue(t, slimentity.IsObject(reflect.ValueOf(ptrToADemoStruct2)), "*aDemoStruct2 is an object")
	assert.Equals(t, "demo2", processor.serializeObjectsIn(ptrToADemoStruct2), "*struct with ToString")
	anEmptyStruct := emptyStruct{}
	assert.Equals(t, "slimprocessor.emptyStruct", processor.serializeObjectsIn(anEmptyStruct), "struct without ToString")
	list := slimentity.NewSlimListContaining([]slimentity.SlimEntity{"test2", aDemoStruct1, aDemoStruct2, anEmptyStruct})
	assert.IsTrue(t, !slimentity.IsObject(reflect.ValueOf(list)), "list is no object")
	assert.Equals(t, "[test2, demo1, demo2, slimprocessor.emptyStruct]",
		processor.serializeObjectsIn(list).(*slimentity.SlimList).ToString(), "list with objects")
}
