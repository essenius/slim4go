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

package slim4go

// TODO:
// * Split up Parser
// * Use more packages, see fixture

import (
	"os"

	"github.com/essenius/slim4go/internal/inject"
	"github.com/essenius/slim4go/internal/slimlog"
	"github.com/essenius/slim4go/internal/slimserver"
)

//Server provides the Slim server
func Server() *slimserver.SlimServer {
	// We need to do this as early as possible.
	// It gets the command line parameters and initializes the log
	inject.Context().Initialize(os.Args)
	return inject.SlimServer()
}

// Serve runs the Slim Server process.
func Serve() {
	if err := Server().Serve(); err != nil {
		slimlog.Error.Print(err)
	}
}

// RegisterFixture registers a type as fixture using a constructor func.
func RegisterFixture(constructor interface{}) error {
	return Server().RegisterFixture(constructor)
}

// RegisterFixturesFrom registers a number of fixtures using a factory (having pointer receivers named NewXxx).
func RegisterFixturesFrom(factory interface{}) error {
	return Server().RegisterFixturesFrom(factory)
}
