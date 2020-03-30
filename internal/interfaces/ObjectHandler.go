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

package interfaces

import (
	"reflect"

	"github.com/essenius/slim4go/internal/slimentity"
)

// ObjectHandler encapsulates all object handler functions including the collector, as well as the object type itself.
type ObjectHandler interface {
	Collector
	ObjectSerializer
	AddObjectByConstructor(instanceName string, constructor reflect.Value, args []string) error
	InvokeMemberOn(instance interface{}, method string, args *slimentity.SlimList) (slimentity.SlimEntity, error)
	InstancesWithPrefix(prefix string) []interface{}
}
