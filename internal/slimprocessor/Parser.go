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
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/essenius/slim4go/internal/apperrors"
	"github.com/essenius/slim4go/internal/interfaces"
	"github.com/essenius/slim4go/internal/slimentity"

	"golang.org/x/net/html"
)

// Type definitions and constructors

// Parser provides parsing fuctions
type Parser struct {
	objectSerializer interfaces.ObjectSerializer
	symbols          interfaces.SymbolCollector
}

// NewParser creates a new Parser.
func NewParser(symbols interfaces.SymbolCollector) *Parser {
	aParser := new(Parser)
	aParser.symbols = symbols
	return aParser
}

// SetObjectSerializer injects the object serializer. We need to do it this way because of a bidirectional relationship.
func (aParser *Parser) SetObjectSerializer(objectSerializer interfaces.ObjectSerializer) {
	aParser.objectSerializer = objectSerializer
}

// Helper functions

func isPredefinedType(inputType reflect.Type) bool {
	return inputType.Name() != "" && inputType.PkgPath() == ""
}

func parseHTMLTable(input string) ([][]string, error) {
	tokenizer := html.NewTokenizer(strings.NewReader(input))
	table := [][]string{}
	row := []string{}
	_ = tokenizer.Next()
	if tokenizer.Token().Data != "table" {
		return table, toErrorf("Could not parse '%v' as an HTML table", input)
	}
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			if len(row) > 0 {
				table = append(table, row)
			}
			if tokenizer.Err() == io.EOF {
				return table, nil
			}
			return table, tokenizer.Err()
		case html.StartTagToken:
			if tagName, _ := tokenizer.TagName(); string(tagName) == "td" {
				content := tokenizer.Next()
				if content == html.TextToken {
					text := (string)(tokenizer.Text())
					cell := strings.TrimSpace(text)
					row = append(row, cell)
				}
			}
		case html.EndTagToken:
			if tagName, _ := tokenizer.TagName(); string(tagName) == "tr" {
				if len(row) > 0 {
					table = append(table, row)
					row = []string{}
				}
			}
		}
	}
}

func splitOnNextComma(input string) (string, string) {
	var nextComma int
	if nextComma = strings.Index(input, ", "); nextComma == -1 {
		return input, ""
	}
	return input[:nextComma], input[nextComma+2:]
}

func toErrorf(template string, param ...interface{}) error {
	return fmt.Errorf(template, param...)
}

func toMatchingClosingBracket(input string) (string, string, error) {
	nesting := 0
	for i, char := range input {
		if char == ']' && nesting == 0 {
			return input[:i], input[i+1:], nil
		}
		if char == ']' {
			nesting--
		} else if char == '[' {
			nesting++
		}
	}
	return "", "", toErrorf("Could not find matching ']' in '[%v'", input)
}

// Methods

// CallFunction calls a function including parsing/marshalling the input parameters and transforming the output.
// TODO: This is part of the bidirectional dependency issue.
func (aParser *Parser) CallFunction(function reflect.Value, args []string) (returnEntity slimentity.SlimEntity, err error) {
	arguments, err := aParser.matchParamType(args, function)
	if err != nil {
		return "", err
	}
	// The function we call might panic (after all, FitNesse is a testing framework). Be ready for that.
	defer func() {
		if panicData := recover(); panicData != nil {
			returnEntity = nil
			err = fmt.Errorf("Panic: %v", apperrors.ErrorToString(panicData))
		}
	}()
	returnValue := function.Call(*arguments)
	return slimentity.TransformCallResult(returnValue), nil
}

func (aParser *Parser) matchParamType(paramIn []string, method reflect.Value) (*[]reflect.Value, error) {
	result := []reflect.Value{}
	methodType := method.Type()
	if methodType.Kind() != reflect.Func {
		return nil, toErrorf("%v is not a function", methodType.String())
	}
	numParams := methodType.NumIn()
	var paramCountMatch bool
	if methodType.IsVariadic() {
		paramCountMatch = len(paramIn) >= numParams-1
	} else {
		paramCountMatch = len(paramIn) == numParams
	}
	if !paramCountMatch {
		return nil, toErrorf("Expected %v parameter(s) but got %v", numParams, len(paramIn))
	}
	for paramIndex, param := range paramIn {
		paramType := paramTypeFor(methodType, paramIndex)
		resultValue, err := aParser.Parse(param, paramType)
		if err != nil {
			return nil, err
		}
		result = append(result, reflect.ValueOf(resultValue))
	}
	return &result, nil
}

// Parse takes an input string and parses it to the desired type.
func (aParser *Parser) Parse(input string, targetType reflect.Type) (interface{}, error) {
	if isPredefinedType(targetType) {
		resolvedInput := aParser.ReplaceSymbolsIn(input)
		return aParser.parsePredefined(resolvedInput, targetType)
	}
	if symbolValue, ok := aParser.symbols.NonTextSymbol(input); ok {
		if reflect.TypeOf(symbolValue).AssignableTo(targetType) {
			return symbolValue, nil
		}
		return nil, toErrorf("Symbol '%v' of type '%v' not assignable to type '%v'", input, reflect.TypeOf(symbolValue), targetType)
	}
	// target is no predefined type, input is no list, and no Symbol as Object.
	// Check if it needs to be put in an interface - then we need to infer the type.
	resolvedInput := aParser.ReplaceSymbolsIn(input)
	resolvedInputType := reflect.TypeOf(resolvedInput)
	if resolvedInputType.Kind() == reflect.String && targetType.Kind() == reflect.Interface {
		return aParser.parseToInferredType(resolvedInput), nil
	}

	// See if the target is a fixture we can parse into
	if slimentity.IsObjectType(targetType) {
		return aParser.parseFixture(resolvedInput, targetType)
	}
	switch targetType.Kind() {
	case reflect.Map:
		return aParser.parseMap(resolvedInput, targetType)
	case reflect.Slice:
		result, err := aParser.parseSlice(resolvedInput, targetType)
		return result, err
	case reflect.Ptr:
		return aParser.parsePtr(input, targetType)
	default:
		return nil, toErrorf("Don't know how to resolve '%v' into '%v'", resolvedInput, targetType)
	}
}

// parseFixture tries to parse a fixture (struct) by calling its Parse(string) pointer receiver.
// this is a bit tricky, as parseFixture can be called when inputType is the struct type itself.
// If so, we we need to create the necessary pointer, make the call, and dereference afterwards
func (aParser *Parser) parseFixture(input string, inputType reflect.Type) (interface{}, error) {
	return aParser.objectSerializer.Deserialize(inputType, input)
}

// parseMap converts a hash table (rows of two columns) into a Map of the specified type.
// It uses the HTML table format for this, as specified by Slim
func (aParser *Parser) parseMap(input string, targetType reflect.Type) (interface{}, error) {
	matrix, err := parseHTMLTable(input)
	if err != nil {
		return nil, toErrorf("'%v' is not a valid specification for '%v'", input, targetType)
	}
	length := len(matrix)
	returnValue := reflect.MakeMapWithSize(targetType, length)
	for _, row := range matrix {
		if len(row) != 2 {
			return nil, toErrorf("row '%v' in hash '%v' does not have two cells", row, targetType)
		}
		var key, value interface{}
		var err error
		if key, err = aParser.Parse(row[0], targetType.Key()); err != nil {
			return nil, toErrorf("Could not parse key '%v' in hash '%v'", row[0], targetType)
		}
		if value, err = aParser.Parse(row[1], targetType.Elem()); err != nil {
			return nil, toErrorf("Could not parse value '%v' in hash '%v'", row[1], targetType)
		}
		returnValue.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
	}
	return returnValue.Interface(), nil
}

// parsePredefined takes the resolved input string and tries to convert it to the specified predefined type.
// Cannot handle complex numbers since there is no ParseComplex at this time
func (aParser *Parser) parsePredefined(resolvedInput string, targetType reflect.Type) (interface{}, error) {
	var result interface{}
	var err error
	switch targetType.Kind() {
	case reflect.Bool:
		result, err = strconv.ParseBool(resolvedInput)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result, err = strconv.ParseInt(resolvedInput, 0, targetType.Bits())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result, err = strconv.ParseUint(resolvedInput, 0, targetType.Bits())
	case reflect.Float32, reflect.Float64:
		result, err = strconv.ParseFloat(resolvedInput, targetType.Bits())
	case reflect.String:
		result, err = resolvedInput, nil
	default:
		result, err = nil, errors.New("")
	}
	if err == nil {
		return reflect.ValueOf(result).Convert(targetType).Interface(), nil
	}
	return nil, toErrorf("Could not convert '%v' to type '%v'", resolvedInput, targetType.String())
}

func (aParser *Parser) parsePtr(input string, targetType reflect.Type) (interface{}, error) {
	result, err := aParser.Parse(input, targetType.Elem())
	if err != nil {
		return nil, err
	}
	pointer := reflect.New(reflect.TypeOf(result))
	pointer.Elem().Set(reflect.ValueOf(result))
	return pointer.Interface(), nil
}

func (aParser *Parser) parseSlice(input string, targetType reflect.Type) (interface{}, error) {
	input = strings.TrimSpace(input)
	if len(input) > 0 && input[0] == '[' {
		var entry interface{}
		var err error
		if entry, _, err = aParser.parseSubslice(input[1:], targetType); err != nil {
			return nil, err
		}
		return entry, nil
	}
	return nil, toErrorf("'%v' is not an array", input)
}

func (aParser *Parser) parseSliceInternal(input string, targetType reflect.Type) (interface{}, error) {
	slice := reflect.MakeSlice(reflect.SliceOf(targetType.Elem()), 0, 0)
	for i := 0; ; i++ {
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			return slice.Interface(), nil
		}
		var (
			sliceValue interface{}
			next       string
			err        error
		)
		if input[0] == '[' {
			if sliceValue, next, err = aParser.parseSubslice(input[1:], targetType.Elem()); err != nil {
				return nil, err
			}
			_, input = splitOnNextComma(next)
		} else {
			var entry string
			entry, input = splitOnNextComma(input)
			sliceValue, err = aParser.Parse(entry, targetType.Elem())
			if err != nil {
				return nil, toErrorf("Can't parse '%v' as element for slice '%v'", entry, targetType)
			}
		}
		slice = reflect.Append(slice, reflect.ValueOf(sliceValue))
	}
}

func (aParser *Parser) parseSubslice(input string, targetType reflect.Type) (interface{}, string, error) {
	entry, next, err1 := toMatchingClosingBracket(input)
	if err1 != nil {
		return nil, "", err1
	}
	result, err2 := aParser.parseSliceInternal(entry, targetType)
	if err2 != nil {
		return nil, "", err2
	}
	return result, next, nil
}

func (aParser *Parser) parseToInferredType(input string) interface{} {
	if result, err := strconv.ParseInt(input, 0, 0); err == nil {
		return result
	}
	if result, err := strconv.ParseUint(input, 0, 0); err == nil {
		return result
	}
	if result, err := strconv.ParseFloat(input, 0); err == nil {
		return result
	}
	if result, err := strconv.ParseBool(input); err == nil {
		return result
	}
	return input
}

// ReplaceSymbolsIn replaces all occurrences of symbols in a string to their value.
func (aParser *Parser) ReplaceSymbolsIn(source string) string {
	regex := regexp.MustCompile(`\$` + symbolPattern)
	return regex.ReplaceAllStringFunc(source, aParser.replaceSymbolValue)
}

// ReplaceSymbols does ReplaceSymbolsIn recursively for a SlimList.
func (aParser *Parser) ReplaceSymbols(source interface{}) interface{} {
	if slimentity.IsSlimList(source) {
		sourceList := source.(*slimentity.SlimList)
		result := slimentity.NewSlimList()
		for _, value := range *sourceList {
			result.Append(aParser.ReplaceSymbols(value))
		}
		return result
	}
	return aParser.ReplaceSymbolsIn(source.(string))
}

func (aParser *Parser) replaceSymbolValue(symbolName string) string {
	if symbolValue := aParser.symbols.Get(symbolName); symbolValue != nil {
		if slimentity.IsObject(symbolValue) {
			return aParser.objectSerializer.Serialize(symbolValue)
		}
		return symbolValue.(string)
	}
	return symbolName
}
