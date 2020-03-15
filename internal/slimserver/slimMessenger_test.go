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

package slimserver

import (
	"reflect"
	"testing"

	"github.com/essenius/slim4go/internal/slimcontext"

	"github.com/essenius/slim4go/internal/assert"
)

func TestSlimMessengerNew(t *testing.T) {
	context := slimcontext.InjectContext()
	context.InstructionTimeout = 0 // force  re-init
	context.Initialize([]string{"app", "1"})
	pipeMessenger := injectMessenger()
	assert.Equals(t, 1, context.Port, "Port is indeed 1")
	assert.Equals(t, reflect.TypeOf(new(slimPipe)), reflect.TypeOf(pipeMessenger), "Port 1 results in slimPipe")

	context.InstructionTimeout = 0
	context.Initialize([]string{"app", "8485"})
	socketMessenger := injectMessenger()
	assert.Equals(t, reflect.TypeOf(new(slimSocket)), reflect.TypeOf(socketMessenger), "Port 8485 results in slimSocket")
	assert.Equals(t, 8485, socketMessenger.(*slimSocket).port, "Port OK in slimSocket")
}
