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

package inject

import (
	"testing"

	"github.com/essenius/slim4go/examples/demofixtures"
	"github.com/essenius/slim4go/internal/assert"
)

func TestInjectSlimServer(t *testing.T) {
	context := Context()
	context.Initialize([]string{"app", "1"})
	server := SlimServer()
	// TODO: bit of a lazy test. Can use e.g. a mock messenger like in SlimServer.
	assert.Equals(t, nil, server.RegisterFixturesFrom(demofixtures.NewTemperatureFactory()), "Registering fixtures succeeded")
}
