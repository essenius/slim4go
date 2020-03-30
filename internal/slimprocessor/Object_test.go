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
	ProductID  string
	Units      int
	UnitPrice  float64
	unexported bool
}

func NewOrder() *Order {
	return new(Order)
}

func (order *Order) Price() float64 {
	return order.UnitPrice * float64(order.Units)
}

func (order *Order) Parse(input string) {
	panic("Parse failed")
}

func (order *Order) SetProduct(productID string, unitPrice float64) {
	order.ProductID = productID
	order.UnitPrice = unitPrice
}

type MockParser struct{}

func (parser MockParser) CallFunction(function reflect.Value, args []string) (slimentity.SlimEntity, error) {
	return nil, nil
}

func (parser MockParser) Parse(input string, targetType reflect.Type) (interface{}, error) {
	return reflect.Zero(targetType).Interface(), nil
}

func (parser MockParser) ReplaceSymbolsIn(fixtureName string) string {
	return ""
}

func TestObjectTryField(t *testing.T) {
	parser := new(MockParser)
	anObject := newObject(reflect.ValueOf(NewOrder()), parser)
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
	assert.Equals(t, "0", slimentity.ToString(entity3), "GetUnits returns right value")

	fields4 := []string{"UnitPrice", "GetUnitPrice"}
	entity4, err4 := anObject.tryField(fields4, slimentity.NewSlimList())
	assert.Equals(t, nil, err4, "no error get")
	assert.Equals(t, "0", slimentity.ToString(entity4), "UnitPrice returns right value")

	fields5 := []string{"unexported"}
	_, err5 := anObject.tryField(fields5, slimentity.NewSlimList())
	assert.Equals(t, "Can't get value for 'unexported'", err5.Error(), "error getting unexported")

	fields6 := []string{"unexported"}
	_, err6 := anObject.tryField(fields6, slimentity.NewSlimListContaining([]slimentity.SlimEntity{"true"}))
	assert.Equals(t, "Can't set value for 'unexported'", err6.Error(), "error setting unexported")

	fields7 := []string{"bogus"}
	_, err7 := anObject.tryField(fields7, slimentity.NewSlimList())
	assert.Equals(t, "bogus: Field not found", err7.Error(), "nonexistig field")

	bogusObject := newObject(reflect.ValueOf(1), parser)
	_, err8 := bogusObject.tryField(fields7, slimentity.NewSlimList())
	assert.Equals(t, "bogus: Field not found", err8.Error(), "object is not a (pointer to a) struct")

}
