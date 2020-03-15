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
	"strings"
)

// Definitions and constructors for Fixture

// fixure is what FitNesse documentation often calls a class. The constructor is a function that creates the object.
// Using interface{} is the only type I know of that can contain a function with an unspecified number of parameters and return values.
// We need this because the design is based on allowing maximum freedom for the fixture developers.
type fixture struct {
	name        string
	constructor interface{}
}

func newFixture(fixtureName string, constructor interface{}) *fixture {
	aFixture := new(fixture)
	aFixture.name = fixtureName
	aFixture.constructor = constructor
	return aFixture
}

// FixtureMap holds all registered fixtures.
type FixtureMap map[string]*fixture

var fixtureInstance *FixtureMap

// InjectFixtures injects a FixtureMap.
func InjectFixtures() *FixtureMap {
	if fixtureInstance == nil {
		fixtureMap := make(FixtureMap)
		fixtureInstance = &fixtureMap
	}
	return fixtureInstance
}

// Methods for Fixture

func (fixtures *FixtureMap) fixtureNamed(fixtureName string) *fixture {
	aFixture, _ := (*fixtures)[fixtureName]
	return aFixture
}

// RegisterFixture registers a fixture with associated constructor function.
func (fixtures *FixtureMap) RegisterFixture(fixtureName string, constructor interface{}) {
	aFixture := newFixture(fixtureName, constructor)
	(*fixtures)[fixtureName] = aFixture
}

// Definitions and constructors for object. An object is an instantiated fixture.
// It points to an instance which is typically a pointer to a struct.

type object struct {
	instanceName string
	instance     interface{}
	fixtureName  string
}

type objectMap map[string]*object

func newObject(instanceName string, instance interface{}, fixtureName string) *object {
	anObject := new(object)
	anObject.instanceName = instanceName
	anObject.instance = instance
	anObject.fixtureName = fixtureName
	return anObject
}

func newObjectMap() *objectMap {
	anObjectMap := make(objectMap)
	return &anObjectMap
}

func injectObjects() *objectMap {
	return newObjectMap()
}

// Methods for object

func (objects *objectMap) AnyObject() *object {
	for _, value := range *objects {
		return value
	}
	return nil
}

func (objects *objectMap) addObject(instanceName string, instance interface{}, fixtureName string) {
	anObject := newObject(instanceName, instance, fixtureName)
	(*objects)[instanceName] = anObject
}

func (objects *objectMap) Length() int {
	return len(*objects)
}

func (objects *objectMap) objectNamed(instanceName string) *object {
	anObject, _ := (*objects)[instanceName]
	return anObject
}

func (objects *objectMap) objectsWithPrefix(prefix string) *objectMap {
	result := newObjectMap()
	for _, anObject := range *objects {
		if strings.HasPrefix(anObject.instanceName, prefix) {
			(*result)[anObject.instanceName] = anObject
		}
	}
	return result
}

func (objects *objectMap) setObjectInstance(instanceName string, instance interface{}) error {
	anObject := objects.objectNamed(instanceName)
	if anObject == nil {
		return fmt.Errorf("instance not found")
	}
	anObject.instance = instance
	return nil
}

// TODO: implement graceful naming (Set/Get, capitals, underscores?)
