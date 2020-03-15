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
	"io"
	"os"
	"time"

	"github.com/essenius/slim4go/internal/slimcontext"
	"github.com/essenius/slim4go/internal/slimlog"
)

const defaultTimeout = 30 * time.Second

type slimMessenger interface {
	io.Reader
	Listen() error
	SendMessage(message string) error
}

func injectMessenger() slimMessenger {
	context := slimcontext.InjectContext()
	slimlog.Trace.Printf("Timeout is %v", context.ConnectionTimeout)
	return newSlimMessenger(context.Port, context.ConnectionTimeout)
}

func newSlimMessenger(port int, timeout time.Duration) slimMessenger {
	var messenger slimMessenger

	slimUsesPipe := port == 1
	if slimUsesPipe {
		slimlog.Trace.Println("Using pipes")
		messenger = newSlimPipe(os.Stdin, os.Stdout, slimlog.Info, timeout)
	} else {
		slimlog.Trace.Println("Using socket on port ", port)
		messenger = newSlimSocket(port, slimlog.Info, timeout)
	}
	return messenger
}
