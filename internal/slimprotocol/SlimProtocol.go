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

package slimprotocol

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

//AbortSuite returns the instruction to abort executing the suite.
func AbortSuite(reason string) string {
	return "__EXCEPTION__:ABORT_SLIM_SUITE:message:<<" + reason + ">>"
}

//AbortTest returns the instruction to abort executing the suite.
func AbortTest(reason string) string {
	return "__EXCEPTION__:ABORT_SLIM_TEST:message:<<" + reason + ">>"
}

// Bye is the incoming instruction to quit.
func Bye() string {
	return "bye"
}

// CouldNotInvokeConstructor returns the exception that the instructor could not be invoked.
func CouldNotInvokeConstructor(fixtureName string) string {
	return Exceptionf("COULD_NOT_INVOKE_CONSTRUCTOR %v", fixtureName)
}

// Exception returns the message in the FitNesse exception format.
// If the message starts with 'AbortTest:', the current test will be aborted, and with 'AbortSuite:' the suite will be aborted (both case insensitive).
func Exception(exception string) string {
	regex := regexp.MustCompile("(?i)^abort(test|suite):(.*)")
	match := regex.FindStringSubmatch(exception)
	if match == nil {
		return "__EXCEPTION__:message:<<" + exception + ">>"
	}
	if strings.EqualFold("test", match[1]) {
		return AbortTest(match[2])
	}
	return AbortSuite(match[2])
}

// Exceptionf returns a message using a formatting
func Exceptionf(template string, param ...interface{}) string {
	return Exception(fmt.Sprintf(template, param...))
}

// MalformedInstruction returns the exception that the insinstruction could not be parsed successfully.
func MalformedInstruction(instruction string) string {
	return Exceptionf("MALFORMED_INSTRUCTION %v", instruction)
}

// NoFixture returns an exception message that the fixture was not found.
func NoFixture(fixtureName string) string {
	return Exceptionf("NO_CLASS %v", fixtureName)
}

// NoConstructor returns an exception that no suitable constructor could be found.
func NoConstructor(fixtureName string) string {
	return Exceptionf("NO_CONSTRUCTOR %v", fixtureName)
}

// NoConverterForArgumentNumber returns an exception message that one of the parameters could not be converted. unused at this point.
func NoConverterForArgumentNumber(argumentType string) string {
	return Exceptionf("NO_CONVERTER_FOR_ARGUMENT_NUMBER %v", argumentType)
}

// NoInstance returns an exception that no instance could be found.
func NoInstance(instanceName string) string {
	return Exceptionf("NO_INSTANCE %v", instanceName)
}

// NoMethodInFixture returns an exception that no suitable method could be found in a fixture
func NoMethodInFixture(methodName string, fixtureName string, argsLength int) string {
	return Exceptionf("NO_METHOD_IN_CLASS %v[%v] %v", methodName, argsLength, fixtureName)
}

// Null returns the representation of an empty response
func Null() string {
	return "null"
}

// OK retuns a success response
func OK() string {
	return "OK"
}

// TimedOut returns that a timeout has occurred.
func TimedOut(timeout time.Duration) string {
	return Exceptionf("TIMED_OUT %v", int(timeout.Round(time.Second).Seconds()))
}

// Version returns the Slim protocol version.
func Version() string {
	return "Slim -- V0.5\n"
}

//Void returns that a call resulted in a void response.
func Void() string {
	return "/__VOID__/"
}
