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
	"testing"

	"github.com/essenius/slim4go/internal/assert"
)

func TestSlimObjectWithPrefix(t *testing.T) {
	objects := newObjectMap()
	objects.addObject("test1", 1, "fixture1")
	objects.addObject("test2", 2, "fixture1")
	objects.addObject("library1", 3, "fixture2")
	objects.addObject("all1", 4, "all2")
	objects.addObject("library2", 5, "fixture4")
	libraries := objects.objectsWithPrefix("library")
	assert.Equals(t, 2, libraries.Length(), "Length OK")
	assert.Equals(t, 3, (*libraries)["library1"].instance, "entry 1 exists")
	assert.Equals(t, 5, (*libraries)["library2"].instance, "entry 2 exists")
}

func TestSlimObject(t *testing.T) {
	objects := newObjectMap()
	assert.Equals(t, nil, objects.AnyObject(), "AnyObject on empty map returns nil")
	objects.addObject("test1", 1, "fixture1")
	assert.Equals(t, "test1", objects.AnyObject().instanceName, "AnyObject on map with one entry returns that entry")
	err := objects.setObjectInstance("nonexisting", 2)
	assert.IsTrue(t, nil != err, "Error occurred")
	assert.Equals(t, "instance not found", err.Error(), "Error message OK")
}
