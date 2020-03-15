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

	"github.com/essenius/slim4go/internal/slimentity"
)

// Definitions and constructors

type symbolTable map[string]interface{}

func injectSymbolTable() *symbolTable {
	return newSymbolTable()
}

func newSymbolTable() *symbolTable {
	symbols := make(symbolTable)
	return &symbols
}

// Methods

const symbolPattern = `[a-zA-Z][a-zA-Z0-9_]*`

func (symbols *symbolTable) IsValidSymbol(source string) bool {
	regex := regexp.MustCompile(`^\$` + symbolPattern + "$")
	return regex.MatchString(source)
}

func (symbols *symbolTable) IsValidSymbolName(source string) bool {
	regex := regexp.MustCompile("^" + symbolPattern + "$")
	return regex.MatchString(source)
}

func (symbols *symbolTable) NonTextSymbol(symbolName string) (interface{}, bool) {
	if symbols.IsValidSymbol(symbolName) {
		value, ok := symbols.ValueOf(symbolName)
		if ok {
			if reflect.TypeOf(value).Kind() != reflect.String {
				return value, true
			}
		}
	}
	return nil, false
}

func (symbols *symbolTable) ReplaceSymbolsIn(source string) string {
	regex := regexp.MustCompile(`\$` + symbolPattern)
	return regex.ReplaceAllStringFunc(source, symbols.ReplaceValue)
}

func (symbols *symbolTable) ReplaceSymbols(source interface{}) interface{} {
	if slimentity.IsSlimList(source) {
		sourceList := source.(*slimentity.SlimList)
		result := slimentity.NewSlimList()
		for _, value := range *sourceList {
			result.Append(symbols.ReplaceSymbols(value))
		}
		return result
	}
	return symbols.ReplaceSymbolsIn(source.(string))
}

func (symbols *symbolTable) ReplaceValue(symbolName string) string {
	if symbolValue, ok := symbols.ValueOf(symbolName); ok {
		symbolValueValue := reflect.ValueOf(symbolValue)
		if isObject(symbolValueValue) {
			return stringifyObject(symbolValueValue).(string)
		}
		return symbolValue.(string)
	}
	return symbolName
}

func (symbols *symbolTable) SetSymbol(symbol string, value interface{}) error {
	if symbols.IsValidSymbolName(symbol) {
		(*symbols)[symbol] = value
		return nil
	}
	return fmt.Errorf("Invalid symbol name: %v", symbol)
}

func (symbols *symbolTable) ValueOf(symbolName string) (interface{}, bool) {
	symbolValue, ok := (*symbols)[symbolName[1:]]
	return symbolValue, ok
}
