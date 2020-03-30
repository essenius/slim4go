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
	"github.com/essenius/slim4go/internal/context"
	"github.com/essenius/slim4go/internal/fixture"
	"github.com/essenius/slim4go/internal/interfaces"
	"github.com/essenius/slim4go/internal/slimlog"
	"github.com/essenius/slim4go/internal/slimprocessor"
	"github.com/essenius/slim4go/internal/slimserver"
	"github.com/essenius/slim4go/internal/standardlibrary"
)

// ActorStack injects an ActorStack object
func ActorStack() *standardlibrary.ActorStack {
	return standardlibrary.NewActorStack()
}

var contextInstance *context.Context

// Context provides the application context (e.g. command line parameters)
func Context() *context.Context {
	if contextInstance == nil {
		contextInstance = context.New()
	}
	return contextInstance
}

var messengerInstance interfaces.SlimMessenger

// Messenger provides a Messenger instance that can send messages to the Slim client.
func Messenger() interfaces.SlimMessenger {
	if messengerInstance == nil {
		context := Context()
		slimlog.Trace.Printf("Timeout is %v", context.ConnectionTimeout)
		messengerInstance = slimserver.NewSlimMessenger(context.Port, context.ConnectionTimeout)
	}
	return messengerInstance
}

var objectHandlerInstance *slimprocessor.ObjectHandler

// ObjectHandler injects an ObjectHandler (single instance)
func ObjectHandler() *slimprocessor.ObjectHandler {
	if objectHandlerInstance == nil {
		// This is a bit tricky as both Parser and StandardLibrary need this ObjectHandler.
		// StandardLibrary is no issue as it has to be injected when objectHandlerInstance exists,
		// but Parser is, as we'd like to inject it via the constructor (TODO).
		// For now we use dependency injection of ObjectHandler into Parser via a method.
		parser := Parser()
		objectHandlerInstance = slimprocessor.NewObjectHandler(parser)
		parser.SetObjectSerializer(objectHandlerInstance)
		objectHandlerInstance.Add("libraryStandard", StandardLibrary())
	}
	return objectHandlerInstance
}

var parserInstance *slimprocessor.Parser

// Parser injects an Parser (single instance)
func Parser() *slimprocessor.Parser {
	if parserInstance == nil {
		parserInstance = slimprocessor.NewParser(SymbolTable())
	}
	return parserInstance
}

var registryInstance *fixture.Registry

// Registry returns a Registry (single instance).
func Registry() *fixture.Registry {
	if registryInstance == nil {
		registryInstance = fixture.NewRegistry()
	}
	return registryInstance
}

// SlimInterpreter injects a Slim Interpreter
func SlimInterpreter() *slimprocessor.SlimInterpreter {
	return slimprocessor.NewSlimInterpreter(StatementProcessor(), Context().InstructionTimeout)
}

// SlimServer provides the Slim server instance.
func SlimServer() *slimserver.SlimServer {
	if slimServerInstance == nil {
		slimServerInstance = slimserver.NewSlimServer(Registry(), Messenger(), SlimInterpreter())
	}
	return slimServerInstance
}

// StandardLibrary injects a StandardLibrary.
func StandardLibrary() *standardlibrary.StandardLibrary {
	return standardlibrary.New(ActorStack(), ObjectHandler())
}

// StatementProcessor injects a StatementProcessor.
func StatementProcessor() *slimprocessor.SlimStatementProcessor {
	return slimprocessor.NewStatementProcessor(Registry(), ObjectHandler(), Parser(), SymbolTable())
}

var symbolTableInstance *slimprocessor.SymbolTable

// SymbolTable injects a symbol table.
func SymbolTable() *slimprocessor.SymbolTable {
	if symbolTableInstance == nil {
		symbolTableInstance = slimprocessor.NewSymbolTable()
	}
	return symbolTableInstance
}

var slimServerInstance *slimserver.SlimServer
