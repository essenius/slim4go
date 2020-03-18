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
)

type objectFactory struct {
	theParser *parser
}

func injectObjectFactory(aParser *parser) *objectFactory {
	return newObjectFactory(aParser)
}

func newObjectFactory(aParser *parser) *objectFactory {
	factory := new(objectFactory)
	factory.theParser = aParser
	return factory
}

func (factory *objectFactory) NewObject(instanceValue reflect.Value) *object {
	return newObject(instanceValue, factory.theParser)
}

func (factory *objectFactory) ConstructObject(constructor reflect.Value, args []string) (*object, error) {
	instance, err := factory.theParser.callFunction(constructor, args)
	if err == nil {
		return factory.NewObject(reflect.ValueOf(instance)), nil
	}
	return nil, err
}
