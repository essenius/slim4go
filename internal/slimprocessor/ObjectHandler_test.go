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
)

func TestObjectHandlerWithPrefix(t *testing.T) {
	objects := NewObjectHandler(nil)
	objects.objectMap = newObjectMap()
	objects.Add("test1", 1)
	objects.Add("test2", 2)
	objects.Add("library1", 3)
	objects.Add("all1", 4)
	objects.Add("library2", 5)
	libraries := objects.InstancesWithPrefix("library")
	assert.Equals(t, 2, len(libraries), "Length OK")
	assert.Equals(t, 3, libraries[0], "entry 1 exists")
	assert.Equals(t, 5, libraries[1], "entry 2 exists")
}

func TestObjectHandler(t *testing.T) {
	objects := NewObjectHandler(nil)
	err := objects.Set("nonexisting", 2)
	assert.IsTrue(t, nil != err, "Error occurred")
	assert.Equals(t, "instance not found", err.Error(), "Error message OK")
}

func TestObjectHandlerSerialize(t *testing.T) {
	parser := NewParser(NewSymbolTable())
	objectHandler := NewObjectHandler(parser)
	parser.SetObjectSerializer(objectHandler)
	_, err := objectHandler.Deserialize(reflect.TypeOf(NewOrder()), "bogus")
	assert.Equals(t, "Panic: Parse failed", err.Error(), "Failing Parse")
}
