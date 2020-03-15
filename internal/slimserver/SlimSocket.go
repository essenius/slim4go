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
	"log"
	"net"
	"strconv"
	"time"

	"github.com/essenius/slim4go/internal/slimlog"
)

type slimSocket struct {
	port       int
	connection net.Conn
	logger     *log.Logger
	timeout    time.Duration
}

var socketInstance *slimSocket

func newSlimSocket(port int, infoLogger *log.Logger, timeout time.Duration) *slimSocket {
	if socketInstance == nil {
		socketInstance = new(slimSocket)
	}
	socketInstance.port = port
	socketInstance.logger = infoLogger
	socketInstance.timeout = timeout
	return socketInstance
}

func (socket *slimSocket) accept(listener net.Listener, channel chan error) {
	connection, err := listener.Accept()
	if err == nil {
		socket.connection = connection
		origin := socket.connection.RemoteAddr().String()
		socket.logger.Println("Connection origin: ", origin)
		channel <- nil
	}
	channel <- err
}

// Listen sets up a socket connection and starts listening.
func (socket *slimSocket) Listen() error {
	slimlog.Trace.Printf("Port %v", socket.port)
	listenSpec := ":" + strconv.Itoa(socket.port)
	socket.logger.Println("Listening at tcp ", listenSpec)
	listener, err1 := net.Listen("tcp", listenSpec)
	if err1 == nil {
		errChannel := make(chan error, 1)
		acceptTimer := time.NewTimer(socket.timeout)
		go socket.accept(listener, errChannel)
		select {
		case err2 := <-errChannel:
			acceptTimer.Stop()
			return err2
		case <-acceptTimer.C:
			listener.Close()
			return fmt.Errorf("Timeout (%v) waiting for a connection", socket.timeout)
		}
	}
	return err1
}

// Read reads a message from the socket connection.
func (socket *slimSocket) Read(buffer []byte) (int, error) {
	socket.connection.SetReadDeadline(time.Now().Add(socket.timeout))
	retrievedBytes, err := socket.connection.Read(buffer)
	return retrievedBytes, err
}

// SendMessage writes a message to the socket connection.
func (socket *slimSocket) SendMessage(message string) error {
	//slimlog.Trace.Printf("socket port: %v", socket.port)
	//slimlog.Trace.Printf("socket connection local address: %v", socket.connection.LocalAddr())
	socket.connection.SetWriteDeadline(time.Now().Add(socket.timeout))
	_, err := socket.connection.Write([]byte(message))
	return err
}
