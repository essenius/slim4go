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

package slimentity

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/essenius/slim4go/internal/apperrors"
)

const (
	terminator     = ':'
	listStarter    = '['
	listTerminator = ']'
)

type slimReader struct {
	*bufio.Reader
}

func newSlimReader(reader io.Reader) *slimReader {
	aSlimReader := new(slimReader)
	aSlimReader.Reader = bufio.NewReader(reader)
	return aSlimReader
}

func (reader *slimReader) skipByte(expected byte) {
	character, err := reader.ReadByte()

	if err != nil {
		panic("SkipByte: No input available")
	}
	if character != expected {
		panic(fmt.Sprintf("skipByte: Expected '%v' but found '%v'", string(expected), string(character)))
	}
}

func (reader *slimReader) readLength() int {
	// lengths are normally 6 characters, but can be more. So we generalize and search for the terminator
	lengthString, err1 := reader.ReadString(terminator)
	if err1 != nil {
		panic(fmt.Sprintf("readLength: Could not find next delimiter '%c' (%v)", terminator, err1.Error()))
	}
	lengthString = lengthString[0 : len(lengthString)-1]
	length, err2 := strconv.Atoi(lengthString)
	if err2 != nil {
		panic(fmt.Sprintf("readLength: Could not interpret length '%v'", lengthString))
	}
	return length
}

func (reader *slimReader) readExactBytes(numberOfBytes int) []byte {
	// we're using bytes on purpose here. The line sizes are in bytes.
	// Note that this means it isn't necessarily the number of characters - SLIM uses UTF-8 in text.

	buffer := make([]byte, numberOfBytes)
	if n, err := io.ReadAtLeast(reader, buffer, numberOfBytes); err != nil {
		panic(fmt.Errorf("readExactBytes: Expected %v bytes from Slim client, but got %v", numberOfBytes, n))
	}
	return buffer
}

func (reader *slimReader) isStartOfList() (bool, int) {
	const minCharsForList = 9 // smallest possible list is [000000:]
	// we need this check since Peek hangs if it tries to go beyond the buffer
	if reader.Buffered() < minCharsForList {
		return false, 0
	}
	for peekLength := minCharsForList - 1; ; peekLength++ {
		potentialListStart, _ := reader.Peek(peekLength) // can't go wrong due to above check
		if potentialListStart[0] != listStarter {
			return false, 0
		}
		listLength, err2 := strconv.Atoi(string(potentialListStart[1 : peekLength-1]))
		if err2 != nil {
			return false, 0
		}
		if potentialListStart[peekLength-1] == terminator {
			_ = reader.readExactBytes(peekLength)
			return true, listLength
		}
	}

}

// ReadRequest is the entry point to read a request message from FitNesse.
func ReadRequest(reader io.Reader) (out interface{}, err error) {
	var aSlimReader = newSlimReader(reader)
	return aSlimReader.readRequest()
}

func (reader *slimReader) readRequest() (out SlimEntity, err error) {
	defer func() {
		if r := recover(); r != nil {
			out = nil
			err = errors.New(apperrors.ErrorToString(r))
		}
	}()

	lengthInBytes := reader.readLength()
	if isStart, numberOfItems := reader.isStartOfList(); isStart {
		list := NewSlimList()
		for line := 0; line < numberOfItems; line++ {
			listEntry, err := reader.readRequest()
			if err != nil {
				return nil, err
			}
			reader.skipByte(terminator)
			list.Append(listEntry)
		}
		reader.skipByte(listTerminator)
		return list, nil
	}
	listEntry := reader.readExactBytes(lengthInBytes)
	return string(listEntry), nil
}

// Marshal converts a SlimEntity to its Slim serialized representation
func Marshal(entity SlimEntity) string {
	if !IsSlimList(entity) {
		entry := convertNull(fmt.Sprintf("%v", entity))
		return fmt.Sprintf("%06d%c%s", len(entry), terminator, entry)
	}
	var stringBuilder strings.Builder
	list := entity.(*SlimList)
	stringBuilder.WriteString(fmt.Sprintf("%c%06d%c", listStarter, list.Length(), terminator))
	for _, listEntry := range *list {
		stringBuilder.WriteString(Marshal(listEntry))
		stringBuilder.WriteRune(terminator)
	}
	stringBuilder.WriteRune(listTerminator)
	result := stringBuilder.String()
	return Marshal(result)
}

func convertNull(message interface{}) string {
	if message == nil {
		return "null"
	}
	return message.(string)
}
