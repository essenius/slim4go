package slim4go

// TODO:
// * More flexibility in Setter/Getters
// * Use of public fields
// * Use of fixture factory to enable import of related fixtures
// * Break bidirectional relationship in Parser and FunctionCaller

import (
	"os"

	"github.com/essenius/slim4go/internal/slimcontext"
	"github.com/essenius/slim4go/internal/slimlog"
	"github.com/essenius/slim4go/internal/slimserver"
)

//Server provides the Slim server
func Server() *slimserver.SlimServer {
	// We need to do this as early as possible.
	// It gets the command line parameters and initializes the log
	slimcontext.InjectContext().Initialize(os.Args)
	return slimserver.InjectSlimServer()
}

// Serve runs the Slim Server process.
func Serve() {
	if err := Server().Serve(); err != nil {
		slimlog.Error.Print(err)
	}
}

// RegisterFixture registers a type as fixture using a constructor.
func RegisterFixture(fixtureName string, constructor interface{}) {
	Server().RegisterFixture(fixtureName, constructor)
}
