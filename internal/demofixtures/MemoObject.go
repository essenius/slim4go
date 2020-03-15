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

package demofixtures

import (
	"math/rand"
	"time"
)

var random *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// MemoObject is used to show fixture life cycle concepts.
type MemoObject struct {
	id   int
	data interface{}
}

// NewMemoObject is the constructor for MemoObject.
func NewMemoObject(data ...interface{}) *MemoObject {
	memoObject := new(MemoObject)
	memoObject.id = random.Int()
	if len(data) > 0 {
		memoObject.data = data[0]
	}
	return memoObject
}

// Data returns the value of the data field.
func (memoObject MemoObject) Data() interface{} {
	return memoObject.data
}

// ID returns the object id.
func (memoObject MemoObject) ID() interface{} {
	return memoObject.id
}

// SetData sets the data field.
func (memoObject MemoObject) SetData(data interface{}) {
	memoObject.data = data
}
