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
	"fmt"
	"reflect"
	"strings"

	"github.com/essenius/slim4go/internal/slimlog"
)

// Definitions and constructors for Fixture

// A fixure is what FitNesse documentation often calls a class.
// A constructor is a function that creates the object, and a factory is an object that provides constructor methods
// Using interface{} is the only type I know of that can contain a function with an unspecified number of parameters and return values.
// We need this because the design is based on allowing maximum freedom for the fixture developers.

type anyMap map[string]interface{}

// Registry defines the fixture registry.
type Registry struct {
	constructor anyMap
	factory     anyMap
	namespace   []string
}

// NewRegistry creates a new fixture registry.
func NewRegistry() *Registry {
	registry := new(Registry)
	registry.constructor = make(anyMap)
	registry.factory = make(anyMap)
	registry.namespace = []string{}
	return registry
}

var registryInstance *Registry

// InjectRegistry returns a registry.
func InjectRegistry() *Registry {
	if registryInstance == nil {
		registryInstance = NewRegistry()
	}
	return registryInstance
}

// internal utility functions

func typeWithoutPointer(name string) string {
	if strings.HasPrefix(name, "*") {
		return name[1:]
	}
	return name
}

// This assumes there is only one return parameter, which makes sense for constructors
func fixtureNameFromConstructor(constructor interface{}) string {
	constructorType := reflect.TypeOf(constructor)
	if constructorType.Kind() != reflect.Func {
		return ""
	}
	firstOutFieldType := constructorType.Out(0).String()
	return typeWithoutPointer(firstOutFieldType)
}

// FixtureNamed returns the fixture with the specified name.
func (registry *Registry) FixtureNamed(fixtureName string) interface{} {
	prefixes := []string{""}
	for _, namespace := range registry.namespace {
		prefixes = append(prefixes, namespace+".")
	}
	for _, prefix := range prefixes {
		nameWithNamespace := prefix + fixtureName
		fixture, ok := registry.constructor[nameWithNamespace]
		if ok {
			return fixture
		}
	}
	return nil
}

// AddFixture registers a fixture definition via its constructor function.
func (registry *Registry) AddFixture(fixtureConstructor interface{}) error {
	fixtureName := fixtureNameFromConstructor(fixtureConstructor)
	if fixtureName != "" {
		registry.constructor[fixtureName] = fixtureConstructor
		return nil
	}
	return fmt.Errorf("Could not add fixture '%v'", fixtureConstructor)
}

// AddFixturesFrom registers fixtures via a fixture factory (pointer to instantiated object with NewXxxx pointer receivers).
func (registry *Registry) AddFixturesFrom(fixtureFactory interface{}) error {
	factoryType := reflect.TypeOf(fixtureFactory)
	for i := 0; i < factoryType.NumMethod(); i++ {
		method := factoryType.Method(i)
		if strings.HasPrefix(method.Name, "New") {
			fixtureName := method.Name[3:]
			slimlog.Trace.Printf("Found %v", fixtureName)
			registry.AddFixture(reflect.ValueOf(fixtureFactory).Method(i).Interface())
		}
	}
	return nil
}

// AddNamespace adds a namespace to the registry (fixture prefix to take into account with searching).
func (registry *Registry) AddNamespace(newNamespace string) {
	for _, value := range registry.namespace {
		if value == newNamespace {
			return
		}
	}
	registry.namespace = append(registry.namespace, newNamespace)
}
