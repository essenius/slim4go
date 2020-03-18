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
)

type objectMap map[string]*object

type objectCollection struct {
	objectMap *objectMap
	factory   *objectFactory
}

func injectObjectCollection(aParser *parser) *objectCollection {
	collection := newObjectCollection(injectObjectFactory(aParser))
	standardLibrary := newStandardLibrary(injectActors(), collection)
	collection.addObject("libraryStandard", standardLibrary)
	return collection
}

func newObjectCollection(factory *objectFactory) *objectCollection {
	collection := new(objectCollection)
	collection.objectMap = newObjectMap()
	collection.factory = factory
	return collection
}

func newObjectMap() *objectMap {
	aMap := make(objectMap)
	return &aMap
}

var objectCollectionInstance *objectCollection

// Methods for objectMap

func (collection *objectCollection) AnyObject() *object {
	for _, value := range *collection.objectMap {
		return value
	}
	return nil
}

func (collection *objectCollection) addObject(instanceName string, instance interface{}) {
	anObject := collection.factory.NewObject(reflect.ValueOf(instance))
	(*collection.objectMap)[instanceName] = anObject
}

func (collection *objectCollection) addObjectByConstructor(instanceName string, constructor reflect.Value, args []string) error {
	anObject, err := collection.factory.ConstructObject(constructor, args)
	if err == nil {
		(*collection.objectMap)[instanceName] = anObject
		return nil
	}
	return err
}

func (collection *objectCollection) Length() int {
	return len(*collection.objectMap)
}

func (collection *objectCollection) objectNamed(instanceName string) *object {
	anObject, _ := (*collection.objectMap)[instanceName]
	return anObject
}

func (collection *objectCollection) objectsWithPrefix(prefix string) *objectMap {
	result := newObjectMap()
	for instanceName, instance := range *collection.objectMap {
		if strings.HasPrefix(instanceName, prefix) {
			(*result)[instanceName] = instance
		}
	}
	return result
}

func (collection *objectCollection) setObjectInstance(instanceName string, instance interface{}) error {
	anObject := collection.objectNamed(instanceName)
	if anObject == nil {
		return fmt.Errorf("instance not found")
	}
	anObject.instanceValue = reflect.ValueOf(instance)
	return nil
}
