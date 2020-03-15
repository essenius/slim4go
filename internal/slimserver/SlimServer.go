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
	"fmt"

	"github.com/essenius/slim4go/internal/slimentity"
	"github.com/essenius/slim4go/internal/slimlog"
	"github.com/essenius/slim4go/internal/slimprocessor"
	"github.com/essenius/slim4go/internal/slimprotocol"
)

// SlimServer is the main object.
type SlimServer struct {
	fixtures  *slimprocessor.FixtureMap
	messenger slimMessenger
	processor *slimprocessor.SlimProcessor
}

var slimServerInstance *SlimServer

// InjectSlimServer provides a Slim server instance.
func InjectSlimServer() *SlimServer {
	if slimServerInstance == nil {
		slimServerInstance = newSlimServer(slimprocessor.InjectFixtures(), injectMessenger(), slimprocessor.InjectSlimProcessor())
	}
	return slimServerInstance
}

func newSlimServer(fixtures *slimprocessor.FixtureMap, messenger slimMessenger, processor *slimprocessor.SlimProcessor) *SlimServer {
	server := new(SlimServer)
	server.fixtures = fixtures
	server.messenger = messenger
	server.processor = processor
	return server
}

// Serve The Slim Server fetching requests, processing them, and returning results.
func (server *SlimServer) Serve() error {

	if err1 := server.messenger.Listen(); err1 != nil {
		return err1
	}
	// not a mistake -- this is the only time that we don't use the size in the SLIM protocol
	if err2 := server.messenger.SendMessage(slimprotocol.Version()); err2 != nil {
		return err2
	}

	for {
		slimlog.Trace.Println("Reading")
		request, err3 := slimentity.ReadRequest(server.messenger)
		if err3 != nil {
			slimlog.Trace.Printf("Read error %v", err3)
			return err3
		}
		slimlog.Trace.Println("After read")
		slimlog.Trace.Println("Request: ", slimentity.Marshal(request)) //request.Marshal())
		if !slimentity.IsSlimList(request) {
			if request.(string) == slimprotocol.Bye() {
				return nil
			}
			return fmt.Errorf("Encountered unexpected command %v", request.(string))
		}
		responseMessage := server.processor.Process(request.(*slimentity.SlimList))
		marshalledResponse := slimentity.Marshal(responseMessage)
		slimlog.Trace.Println("Response: ", marshalledResponse)
		server.messenger.SendMessage(marshalledResponse)
	}
}

// RegisterFixture registers a type as fixture using a constructor.
func (server *SlimServer) RegisterFixture(fixtureName string, constructor interface{}) {
	server.fixtures.RegisterFixture(fixtureName, constructor)
}
