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
	"testing"
	"time"

	"github.com/essenius/slim4go/internal/assert"
)

func TestSlimProtocolException(t *testing.T) {
	assert.Equals(t, "__EXCEPTION__:message:<<error>>", Exception("error"), "Exception")
	assert.Equals(t, "__EXCEPTION__:ABORT_SLIM_SUITE:message:<<quit>>", Exception("AbortSuite:quit"), "Abort Suite")
	assert.Equals(t, "__EXCEPTION__:ABORT_SLIM_TEST:message:<<Quit>>", Exception("aborttest:Quit"), "Abort Test")
}

func TestSlimProtocolExceptionf(t *testing.T) {
	assert.Equals(t, "bye", Bye(), "Bye")
	assert.Equals(t, "__EXCEPTION__:message:<<COULD_NOT_INVOKE_CONSTRUCTOR myFixture>>", CouldNotInvokeConstructor("myFixture"), "Could not invoke constructor")
	assert.Equals(t, "__EXCEPTION__:message:<<MALFORMED_INSTRUCTION qwe>>", MalformedInstruction("qwe"), "Malformed Instruction")
	assert.Equals(t, "__EXCEPTION__:message:<<NO_CLASS testFixture>>", NoFixture("testFixture"), "No Fixture")
	assert.Equals(t, "__EXCEPTION__:message:<<NO_CONVERTER_FOR_ARGUMENT_NUMBER 1>>", NoConverterForArgumentNumber("1"), "No Covnerter For Argument Number")
	assert.Equals(t, "__EXCEPTION__:message:<<NO_CONSTRUCTOR NewObject>>", NoConstructor("NewObject"), "No Constructor")
	assert.Equals(t, "__EXCEPTION__:message:<<NO_INSTANCE myInstance>>", NoInstance("myInstance"), "No Instance")
	assert.Equals(t, "__EXCEPTION__:message:<<NO_METHOD_IN_CLASS myMethod[2] myFixture>>", NoMethodInFixture("myMethod", "myFixture", 2), "No Method In Fixture")
	assert.Equals(t, "null", Null(), "Null")
	assert.Equals(t, "OK", OK(), "OK")
	duration, _ := time.ParseDuration("500s")
	assert.Equals(t, "__EXCEPTION__:message:<<TIMED_OUT 500>>", TimedOut(duration), "Timed out")
	assert.Equals(t, "Slim -- V0.5\n", Version(), "Version")
	assert.Equals(t, "/__VOID__/", Void(), "Void")
}
