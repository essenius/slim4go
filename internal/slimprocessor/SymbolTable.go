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

package slimprocessor

import (
	"fmt"
	"reflect"
	"regexp"
)

// Definitions and constructors

// SymbolTable contains the FitNesse symbols.
type SymbolTable map[string]interface{}

// NewSymbolTable creates a new Symbol table.
func NewSymbolTable() *SymbolTable {
	symbols := make(SymbolTable)
	return &symbols
}

// Methods

const symbolPattern = `[a-zA-Z][a-zA-Z0-9_]*`

func (symbols *SymbolTable) isValidSymbol(source string) bool {
	regex := regexp.MustCompile(`^\$` + symbolPattern + "$")
	return regex.MatchString(source)
}

// IsValidSymbolName returns whether the input (without the $) is a valid symbol name
func (symbols *SymbolTable) IsValidSymbolName(source string) bool {
	regex := regexp.MustCompile("^" + symbolPattern + "$")
	return regex.MatchString(source)
}

// NonTextSymbol returhs whether the symbol is valid and contains something else than a string.
func (symbols *SymbolTable) NonTextSymbol(symbolName string) (interface{}, bool) {
	if symbols.isValidSymbol(symbolName) {
		value, ok := symbols.ValueOf(symbolName)
		if ok {
			if reflect.TypeOf(value).Kind() != reflect.String {
				return value, true
			}
		}
	}
	return nil, false
}

// Add adds an entry to the symbol table. TODO: not used. Optimize interfaces.
func (symbols *SymbolTable) Add(symbolName string, value interface{}) {
	symbols.Set(symbolName, value)
}

// Get gets an entry from the symbol table.
func (symbols *SymbolTable) Get(symbolName string) interface{} {
	symbolValue, _ := symbols.ValueOf(symbolName)
	return symbolValue
}

// Set sets an entry in the symbol table.
func (symbols *SymbolTable) Set(symbolName string, value interface{}) error {
	return symbols.SetSymbol(symbolName, value)
}

// Length gets the number of items in the symbol table. TODO: not used. Optimize interfaces.
func (symbols *SymbolTable) Length() int {
	return len(*symbols)
}

// SetSymbol should be eliminated.
func (symbols *SymbolTable) SetSymbol(symbol string, value interface{}) error {
	if symbols.IsValidSymbolName(symbol) {
		(*symbols)[symbol] = value
		return nil
	}
	return fmt.Errorf("Invalid symbol name: %v", symbol)
}

// ValueOf should be eliminated.
func (symbols *SymbolTable) ValueOf(symbolName string) (interface{}, bool) {
	symbolValue, ok := (*symbols)[symbolName[1:]]
	return symbolValue, ok
}
