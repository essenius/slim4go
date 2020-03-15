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

package slimlog

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	logFileName = flag.String("log", "slim4go.log", "Slim4Go log file name")
	// Info logs informational messages.
	Info *log.Logger
	// Error logs errors.
	Error *log.Logger
	// Trace logs debugging info.
	Trace *log.Logger
)

// making sure Trace works in tests.
var _ = func() error {
	if Trace == nil {
		Trace = log.New(os.Stderr, "Trace:", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return nil
}()

// Initialize sets up the loggers.
func Initialize(slimUsesPipe bool) {
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Could not open", *logFile, ". Slim4Go start Failed")
		os.Exit(1)
	}
	Trace = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)

	if slimUsesPipe {
		Info = log.New(os.Stderr, "SOUT :", 0)
		Error = log.New(os.Stderr, "SERR :", 0)
	} else {
		Info = log.New(os.Stdout, "", 0)
		Error = log.New(os.Stderr, "", 0)
	}
	Trace.Println("Setup trace for production")
}
