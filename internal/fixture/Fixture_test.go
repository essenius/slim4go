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

package fixture

import (
	"reflect"
	"testing"

	"github.com/essenius/slim4go/internal/assert"
)

// Test fixtures and factories

type FixtureFactory struct{}

func NewFixtureFactory() *FixtureFactory {
	return new(FixtureFactory)
}

func (factory *FixtureFactory) NewOrder() *Order {
	return new(Order)
}

func (factory *FixtureFactory) NewMessenger() *Messenger {
	return new(Messenger)
}

type Order struct {
}

func NewOrder() *Order {
	return new(Order)
}

type Messenger struct {
}

func TestFixtureNamespace(t *testing.T) {
	registryInstance = nil
	registry := InjectRegistry()
	registry.AddNamespace("test1")
	assert.Equals(t, 1, len(registry.namespace), "one namespace")
	registry.AddNamespace("test2")
	assert.Equals(t, 2, len(registry.namespace), "two namespaces")
	registry.AddNamespace("test1")
	assert.Equals(t, 2, len(registry.namespace), "still two namespaces after adding already exiting one")
	assert.Equals(t, "test1", registry.namespace[0], "First namespace OK")
	assert.Equals(t, "test2", registry.namespace[1], "Second namespace OK")
	registry.AddNamespace("fixture")
	registry.AddFixture(NewOrder)
	order1 := registry.FixtureNamed("Order")
	assert.IsTrue(t, order1 != nil, "Order found without namespace spec")
	order2 := registry.FixtureNamed("fixture.Order")
	assert.IsTrue(t, order2 != nil, "Order found with namespace spec")
	assert.Equals(t, reflect.TypeOf(order1).String(), reflect.TypeOf(order2).String(), "Order constructor signatures are equal")
	assert.Equals(t, nil, registry.FixtureNamed("bogus"), "Unknown fixture name returns nil")
}

func TestFixtureRegisterFixtures(t *testing.T) {
	registryInstance = nil
	registry := InjectRegistry()
	assert.Equals(t, nil, registry.AddFixture(NewOrder), "Add Fixure NewOrder succeeds")
	assert.Equals(t, 1, len(registry.constructor), "one fixture added")
	assert.Equals(t, nil, registry.AddFixturesFrom(NewFixtureFactory()), "AddFixturesFromFactory succeeded")
	assert.Equals(t, 2, len(registry.constructor), "two more fixtures added, but one already existed")
	order := registry.FixtureNamed("fixture.Order")
	assert.IsTrue(t, order != nil, "Order constructor exists")
	messenger := registry.FixtureNamed("fixture.Messenger")
	assert.IsTrue(t, messenger != nil, "Messenger constructor exists")
	messengerValue := reflect.ValueOf(messenger)
	assert.Equals(t, "func", messengerValue.Kind().String(), "Messenger constructor is a func")
	object := messengerValue.Call([]reflect.Value{})
	assert.Equals(t, "*fixture.Messenger", object[0].Type().String(), "Messenger object created via constructor")
	assert.Equals(t, "Could not add fixture '1'", registry.AddFixture(1).Error(), "Add invalid fixture")
}

func TestFixtureTypeWithoutPointer(t *testing.T) {
	assert.Equals(t, "test1.test2", typeWithoutPointer("*test1.test2"), "with pointer")
	assert.Equals(t, "test1.test2", typeWithoutPointer("test1.test2"), "without pointer")
	assert.Equals(t, "", typeWithoutPointer(""), "empty")
}
