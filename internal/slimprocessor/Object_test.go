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

func (order *Order) SetProduct(productID string, unitPrice float64) {
	order.ProductID = productID
	order.UnitPrice = unitPrice
}

func TestObjectWithPrefix(t *testing.T) {
	objects := injectObjectCollection(injectParser())
	// clean out object map -- objectcollection is a single instance
	objects.objectMap = newObjectMap()
	objects.addObject("test1", 1)
	objects.addObject("test2", 2)
	objects.addObject("library1", 3)
	objects.addObject("all1", 4)
	objects.addObject("library2", 5)
	libraries := objects.objectsWithPrefix("library")
	assert.Equals(t, 2, len(*libraries), "Length OK")
	assert.Equals(t, 3, (*libraries)["library1"].instance(), "entry 1 exists")
	assert.Equals(t, 5, (*libraries)["library2"].instance(), "entry 2 exists")
}

func TestObject(t *testing.T) {
	objects := injectObjectCollection(injectParser())
	// clean out object map -- objectcollection is a single instance
	objects.objectMap = newObjectMap()
	assert.Equals(t, nil, objects.AnyObject(), "AnyObject on empty map returns nil")
	objects.addObject("test1", 1)
	assert.Equals(t, 1, objects.AnyObject().instance(), "AnyObject on collection with one entry returns that entry")
	err := objects.setObjectInstance("nonexisting", 2)
	assert.IsTrue(t, nil != err, "Error occurred")
	assert.Equals(t, "instance not found", err.Error(), "Error message OK")
}

func TestObjectTryField(t *testing.T) {
	anObject := newObject(reflect.ValueOf(NewOrder()), injectParser())
	fields1 := []string{"SetUnits", "Units"}
	entity1, err1 := anObject.tryField(fields1, slimentity.NewSlimListContaining([]slimentity.SlimEntity{"35"}))
	assert.Equals(t, nil, err1, "SetUnits returns no error")
	assert.Equals(t, "/__VOID__/", entity1, "SetUnits returns void")

	fields2 := []string{"UnitPrice", "SetUnitPrice"}
	entity2, err2 := anObject.tryField(fields2, slimentity.NewSlimListContaining([]slimentity.SlimEntity{"2.71"}))
	assert.Equals(t, nil, err2, "SetUnitPrice returns no error")
	assert.Equals(t, "/__VOID__/", entity2, "SetUnitPrice returns void")

	fields3 := []string{"GetUnits", "Units"}
	entity3, err3 := anObject.tryField(fields3, slimentity.NewSlimList())
	assert.Equals(t, nil, err3, "GetUnits returns no error")
	assert.Equals(t, "35", slimentity.ToString(entity3), "GetUnits returns right value")

	fields4 := []string{"UnitPrice", "GetUnitPrice"}
	entity4, err4 := anObject.tryField(fields4, slimentity.NewSlimList())
	assert.Equals(t, nil, err4, "no error get")
	assert.Equals(t, "2.71", slimentity.ToString(entity4), "UnitPrice returns right value")
}
