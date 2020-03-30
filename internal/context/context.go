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

package context

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/essenius/slim4go/internal/slimlog"
)

// Definitions, constructors and wiring

// Context provides the application context (command line params, log init)
type Context struct {
	Port               int
	InstructionTimeout time.Duration
	ConnectionTimeout  time.Duration

	// ErrorAction enables overriding exit in tests
	ErrorAction func(err error)
}

func exit(err error) {
	slimlog.Trace.Println(err.Error())
	slimlog.Error.Fatalln(err.Error())
}

var theContext *Context

// Initialize injects the command line arguments. We can't do that in the constructor
// because we want to replace os.Args by a plain string slice during testing
// (os.Args returns somthing different during testing)
func (context *Context) Initialize(args []string) {
	// prevent multiple initializations
	if context.InstructionTimeout != time.Duration(0) {
		return
	}
	var commandLine = flag.NewFlagSet("slim", flag.ContinueOnError)
	var instructionTimeoutPtr = commandLine.Float64("s", 10, "Instruction timeout")
	var connectionTimeoutPtr = commandLine.Float64("t", 30, "Connection timeout")
	// we handle errors after initializing the logger
	err1 := commandLine.Parse(args[1:])
	var err2 error
	context.Port, err2 = parsePort(commandLine.Args())
	slimlog.Initialize(context.Port == 1)
	if err1 != nil {
		context.ErrorAction(err1)
	}
	if err2 != nil {
		context.ErrorAction(err2)
	}
	context.InstructionTimeout = time.Duration(*instructionTimeoutPtr * float64(time.Second))
	context.ConnectionTimeout = time.Duration(*connectionTimeoutPtr * float64(time.Second))
}

// New creates a new Context
func New() *Context {
	context := new(Context)
	context.ErrorAction = exit
	return context
}

func parsePort(args []string) (int, error) {
	// default the port to 1, as then a fatal error comes through no matter if pipes or sockets are used
	port := 1
	if len(args) == 0 {
		return 1, fmt.Errorf("Missing port specification. Expected params [-s timeout] port")
	}
	portString := args[0]
	var err error
	port, err = strconv.Atoi(portString)
	if err != nil {
		return 1, fmt.Errorf("port '%v' should be numerical", portString)
	}
	if port < 0 {
		return 1, fmt.Errorf("port '%v' should be non-negative", portString)
	}
	return port, nil
}
