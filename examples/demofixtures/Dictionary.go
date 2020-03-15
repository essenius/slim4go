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

// Dictionary is what the name suggests - a dictionary.
type Dictionary struct {
	dictionaryMap map[string]string
}

// NewDictionary is the constructor.
func NewDictionary() *Dictionary {
	return new(Dictionary)
}

// AddItem adds a key value pair to the dictionary.
func (dictionary *Dictionary) AddItem(key string, value string) {
	dictionary.dictionaryMap[key] = value
}

// Set initializes the dictionary with a set of key value pairs.
func (dictionary *Dictionary) Set(dict map[string]string) {
	dictionary.dictionaryMap = dict
}

// Get retrieves the dictionary.
func (dictionary *Dictionary) Get() map[string]string {
	return dictionary.dictionaryMap
}

// Contains returns whether the key exists in the dictionary.
func (dictionary *Dictionary) Contains(key string) bool {
	_, ok := dictionary.dictionaryMap[key]
	return ok
}

// GetValue returns the value belonging to a key
func (dictionary *Dictionary) GetValue(key string) string {
	value, ok := dictionary.dictionaryMap[key]
	if ok {
		return value
	}
	return ""
}
