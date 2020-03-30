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
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/essenius/slim4go/internal/apperrors"
	"github.com/essenius/slim4go/internal/interfaces"
	"github.com/essenius/slim4go/internal/slimentity"
	"github.com/essenius/slim4go/internal/slimprotocol"
)

// Definitions and constructors for object. An object is an instantiated fixture.
// It points to an instance which is typically a pointer to a struct.

type object struct {
	instanceValue reflect.Value
	parser        interfaces.Parser
}

// NewObject creates a new object instance.
func newObject(instanceValue reflect.Value, parser interfaces.Parser) *object {
	anObject := new(object)
	anObject.instanceValue = instanceValue
	anObject.parser = parser
	return anObject
}

// Helpers

func alternativeName(name string, prefix string) string {
	if hasFieldPrefix(name, prefix) {
		return name[3:]
	}
	return prefix + name
}

func getField(field reflect.Value, name string) (slimentity.SlimEntity, error) {
	if field.CanInterface() {
		return slimentity.TransformCallResult([]reflect.Value{field}), nil
	}
	return nil, fmt.Errorf("Can't get value for '%v'", name)
}

func hasFieldPrefix(name string, prefix string) bool {
	return len(name) > 3 && strings.HasPrefix(name, prefix) && unicode.IsUpper(rune(name[3]))
}

func memberNamesFor(name string, argcount int) []string {
	returnValue := []string{name}
	switch argcount {
	case 0:
		returnValue = append(returnValue, alternativeName(name, "Get"))
	case 1:
		returnValue = append(returnValue, alternativeName(name, "Set"))
	}
	return returnValue
}

func paramTypeFor(methodType reflect.Type, paramIndex int) reflect.Type {
	if methodType.IsVariadic() && paramIndex >= methodType.NumIn()-1 {
		return methodType.In(methodType.NumIn() - 1).Elem()
	}
	return methodType.In(paramIndex)
}

// Methods

func (anObject *object) instance() interface{} {
	return anObject.instanceValue.Interface()
}

// InvokeMember invokes a function or sets/gets a field.
func (anObject *object) InvokeMember(memberName string, args *slimentity.SlimList) (slimentity.SlimEntity, error) {
	// We can only use exported methods or fields, which start with a capital.
	// Since in the Java convention that FitNesse uses, methods are in camelCase, we need to capitalize the first letter.
	names := memberNamesFor(strings.Title(memberName), args.Length())
	for _, name := range names {
		method := anObject.instanceValue.MethodByName(name)
		if method.IsValid() {
			return anObject.parser.CallFunction(method, slimentity.ToSlice(args))
		}
	}
	if result, err := anObject.tryField(names, args); err == nil {
		return result, nil
	}
	return "", &apperrors.NotFoundError{Entity: "member", Description: memberName}
}

// Serialize returns the state of an object into a string format.
func (anObject *object) Serialize() string {
	entity, err := anObject.InvokeMember("ToString", slimentity.NewSlimList())
	if err == nil {
		return entity.(string)
	}
	return anObject.instanceValue.Type().String()
}

func (anObject *object) setField(field reflect.Value, value slimentity.SlimEntity, name string) (slimentity.SlimEntity, error) {
	if field.CanSet() {
		fieldType := field.Type()
		result, err := anObject.parser.Parse(slimentity.ToString(value), fieldType)
		if err == nil {
			field.Set(reflect.ValueOf(result))
			return slimprotocol.Void(), nil
		}
	}
	return nil, fmt.Errorf("Can't set value for '%v'", name)
}

func (anObject *object) tryField(fieldNames []string, args *slimentity.SlimList) (slimentity.SlimEntity, error) {
	if anObject.instanceValue.Kind() == reflect.Ptr {
		elemObject := newObject(anObject.instanceValue.Elem(), anObject.parser)
		return elemObject.tryField(fieldNames, args)
	}
	if anObject.instanceValue.Kind() != reflect.Struct {
		return nil, &apperrors.NotFoundError{Entity: "Field", Description: fieldNames[0]}
	}
	for _, name := range fieldNames {
		field := anObject.instanceValue.FieldByName(name)
		if field.IsValid() {
			switch args.Length() {
			case 0:
				return getField(field, name)
			case 1:
				return anObject.setField(field, args.ElementAt(0), name)
			}
		}
	}
	return nil, &apperrors.NotFoundError{Entity: "Field", Description: fieldNames[0]}
}
