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

	"github.com/essenius/slim4go/internal/apperrors"
	"github.com/essenius/slim4go/internal/interfaces"
	"github.com/essenius/slim4go/internal/slimentity"
)

type objectMap map[string]*object

// ObjectHandler contains the instantiated fixtures (i.e. objects) and provides functions to handle them.
type ObjectHandler struct {
	objectMap *objectMap
	parser    interfaces.Parser
}

// NewObjectHandler creates a new ObjectCollection.
func NewObjectHandler(parser interfaces.Parser) *ObjectHandler {
	handler := new(ObjectHandler)
	handler.objectMap = newObjectMap()
	handler.parser = parser
	return handler
}

func newObjectMap() *objectMap {
	aMap := make(objectMap)
	return &aMap
}

// Interface implementation for ObjectHandler

// Add adds a key value pair to the collection.
func (handler *ObjectHandler) Add(instanceName string, value interface{}) {
	handler.addObject(instanceName, value)
}

// AddObjectByConstructor adds a new object to the collection by calling a constructor
func (handler *ObjectHandler) AddObjectByConstructor(instanceName string, constructor reflect.Value, args []string) error {
	anObject, err := handler.constructObject(constructor, args)
	if err == nil {
		(*handler.objectMap)[instanceName] = anObject
		return nil
	}
	return err
}

// Deserialize takes a string representation and converts it into an instantiated object of the required type by calling its Parse method.
// If it's called with a Struct type, we create a pointer first as Parse is always a pointer receiver. Then we dereference the result.
func (handler *ObjectHandler) Deserialize(objectType reflect.Type, input string) (interface{}, error) {
	var parseType reflect.Type
	if objectType.Kind() == reflect.Ptr {
		parseType = objectType.Elem()
	} else {
		// The input type is not a pointer, and Parse is a pointer receiver. New creates a pointer to the type
		parseType = objectType
	}
	instance := reflect.New(parseType).Interface()
	_, err := handler.InvokeMemberOn(instance, "Parse", slimentity.NewSlimListContaining([]slimentity.SlimEntity{input}))
	if _, ok := err.(*apperrors.NotFoundError); ok {
		return nil, toErrorf("No method Parse found for type '%v'", reflect.TypeOf(instance))
	}
	if err != nil {
		return nil, err
	}

	if objectType.Kind() == reflect.Ptr {
		return instance, nil
	}
	// now we have a pointer, but we need the element
	return reflect.ValueOf(instance).Elem().Interface(), nil
}

// Get returns the entry via its key.
func (handler *ObjectHandler) Get(instanceName string) interface{} {
	object := handler.objectNamed(instanceName)
	if object == nil {
		return nil
	}
	return object.instance()
}

// Length returns the number of items in the collection.
func (handler *ObjectHandler) Length() int {
	return len(*handler.objectMap)
}

// InvokeMemberOn finds an instance, and invokes a member on it with the given parameters.
func (handler *ObjectHandler) InvokeMemberOn(instance interface{}, memberName string, args *slimentity.SlimList) (slimentity.SlimEntity, error) {
	anObject := handler.newObject(reflect.ValueOf(instance))
	//var result slimentity.SlimEntity
	result, err := anObject.InvokeMember(memberName, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// InstancesWithPrefix returns all instances of which the name starts with the prefix.
func (handler *ObjectHandler) InstancesWithPrefix(prefix string) []interface{} {
	result := make([]interface{}, 0)
	for instanceName, instance := range *handler.objectMap {
		if strings.HasPrefix(instanceName, prefix) {
			result = append(result, instance.instance())
		}
	}
	return result
}

// Serialize returns a serialized representation of an instantiated object (i.e. call its ToString method).
func (handler *ObjectHandler) Serialize(instance interface{}) string {
	anObject := handler.newObject(reflect.ValueOf(instance))
	return anObject.Serialize()
}

// Set sets an existing enty in the collection to a new value.
func (handler *ObjectHandler) Set(instanceName string, instance interface{}) error {
	anObject := handler.objectNamed(instanceName)
	if anObject == nil {
		return fmt.Errorf("instance not found")
	}
	anObject.instanceValue = reflect.ValueOf(instance)
	return nil
}

// Other methods

func (handler *ObjectHandler) addObject(instanceName string, instance interface{}) {
	anObject := handler.newObject(reflect.ValueOf(instance))
	(*handler.objectMap)[instanceName] = anObject
}

func (handler *ObjectHandler) constructObject(constructor reflect.Value, args []string) (*object, error) {
	instance, err := handler.parser.CallFunction(constructor, args)
	if err == nil {
		return handler.newObject(reflect.ValueOf(instance)), nil
	}
	return nil, err
}

func (handler *ObjectHandler) newObject(value reflect.Value) *object {
	return newObject(value, handler.parser)
}

func (handler *ObjectHandler) objectNamed(instanceName string) *object {
	anObject, _ := (*handler.objectMap)[instanceName]
	return anObject
}
