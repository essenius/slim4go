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
	"testing"

	"github.com/essenius/slim4go/internal/assert"
)

type demoStruct1 struct {
	message string
}

func (demo1 *demoStruct1) ToString() string {
	return demo1.message
}

func (demo1 *demoStruct1) Parse(input string) {
	demo1.message = input
}

type demoStruct2 struct {
	message string
}

func (demo2 demoStruct2) ToString() string {
	return demo2.message
}

func (demo2 *demoStruct2) Parse(input string) {
	demo2.message = input
}

type emptyStruct struct{}

func TestParserIsPredefined(t *testing.T) {
	assertPredefined := func(isPredefined bool, value interface{}, description string) {
		assert.IsTrue(t, !isObject(reflect.ValueOf("")), "empty string is no object")
		assert.IsTrue(t, isObject(reflect.ValueOf(&demoStruct1{})), "pointer to demoStruct1 is an object")
	}
	aSlice := []string{}
	aPointer := &aSlice
	aStruct := demoStruct1{}
	assertPredefined(true, 1, "int")
	assertPredefined(true, "a", "string")
	assertPredefined(true, true, "bool")
	assertPredefined(true, 1.2, "float")
	assertPredefined(false, aSlice, "slice")
	assertPredefined(false, aPointer, "pointer")
	assertPredefined(false, aStruct, "struct")
}

// TODO: change other test helper functions to use this pattern
func TestParserToMatchingCloseBracket(t *testing.T) {
	assertMatchingCloseBracket := func(input string, expectedResult string, expectedRest string, errorMessage string, description string) {
		actualResult, actualRest, err := toMatchingClosingBracket(input)
		if errorMessage == "" {
			assert.Equals(t, nil, err, fmt.Sprintf("No error for %v", description))
		} else {
			assert.Equals(t, errorMessage, err.Error(), fmt.Sprintf("Error for %v", description))
		}

		assert.Equals(t, expectedResult, actualResult, fmt.Sprintf("result for %v", description))
		assert.Equals(t, expectedRest, actualRest, fmt.Sprintf("rest for %v", description))
	}
	assertMatchingCloseBracket("]", "", "", "", "empty list")
	assertMatchingCloseBracket("1, 2]", "1, 2", "", "", "1 dimensional list")
	assertMatchingCloseBracket("[1, 2], [3, 4]], rest", "[1, 2], [3, 4]", ", rest", "", "2 dimensional list")
	assertMatchingCloseBracket("", "", "", "Could not find matching ']' in '['", "error")
}

func TestParserParse(t *testing.T) {
	var i int
	var f float64
	var ui uint32
	var s string
	var r rune
	var b byte
	var ok bool
	assertParse := func(aParser *parser, expectation interface{}, source string, target interface{}, description string) {
		result, err := aParser.parse(source, reflect.ValueOf(target).Type())
		assert.Equals(t, nil, err, fmt.Sprintf("%v NoError", description))
		assert.Equals(t, expectation, result, fmt.Sprintf("%v Value", description))
	}
	symbols := newSymbolTable()
	aParser := newParser(symbols)
	symbols.SetSymbol("test", "text")
	aDemoStruct1 := new(demoStruct1)
	aDemoStruct1.Parse("qwe123")
	symbols.SetSymbol("obj", aDemoStruct1)
	assertParse(aParser, false, "False", ok, "bool")
	assertParse(aParser, 23, "23", i, "int")
	assertParse(aParser, 3.14, "3.14", f, "float64")
	assertParse(aParser, uint32(4095), "0o7777", ui, "uint32")
	assertParse(aParser, rune(51966), "0xCAFE", r, "rune")
	assertParse(aParser, byte(10), "0b1010", b, "byte")
	assertParse(aParser, "text", "text", s, "string")

	var c complex128
	_, err := aParser.parse("1+j", reflect.ValueOf(c).Type())
	assert.Equals(t, "Could not convert '1+j' to type 'complex128'", err.Error(), "Unable to convert")
	result, err2 := aParser.parse("$obj", reflect.TypeOf(aDemoStruct1))
	assert.Equals(t, nil, err2, "err2 is nil with retrieving obj")
	assert.Equals(t, "qwe123", result.(*demoStruct1).ToString(), "Get data from retieved non-text symbol")
	anEmptyStruct := new(emptyStruct)
	_, err3 := aParser.parse("$obj", reflect.TypeOf(anEmptyStruct))
	assert.Equals(t, "Symbol '$obj' of type '*slimprocessor.demoStruct1' not assignable to type '*slimprocessor.emptyStruct'", err3.Error(), "Symbol content not assignable to type")
	aMap := make(map[rune]uint)
	outMap, err4 := aParser.parse("<table></table>", reflect.TypeOf(aMap))
	assert.IsTrue(t, err4 == nil, "no error expected with empty map")
	assert.Equals(t, 0, len(outMap.(map[rune]uint)), "map len=0")
	// testing both parsePtr and parseSlice
	aSlice := &[]string{}
	outSlice, err5 := aParser.parse("[]", reflect.TypeOf(aSlice))
	assert.IsTrue(t, err5 == nil, "no error expected with slice")
	assert.Equals(t, 0, len(*outSlice.(*[]string)), "slice len=0")
	func1 := TestParserParse
	_, err6 := aParser.parse("[]", reflect.TypeOf(func1))
	assert.Equals(t, "Don't know how to resolve '[]' into 'func(*testing.T)'", err6.Error(), "error expected with func")
}

func TestParserParseFixture(t *testing.T) {
	// THis is a bit convoluted because parser needs to know its caller and vice versa
	// TODO: break this bidirectional relationship
	caller := injectFunctionCaller(newSymbolTable())
	aDemoStruct1 := new(demoStruct1)
	aDemoStruct1.Parse("demo1")
	aParser := caller.theParser
	result1, err1 := aParser.parseFixture("new value 1", reflect.TypeOf(aDemoStruct1))
	assert.Equals(t, nil, err1, "no err1")
	assert.Equals(t, "new value 1", result1.(*demoStruct1).message, "value set in parse demo1")
	assert.Equals(t, "demo1", aDemoStruct1.message, "original value not changed demo1")
	aDemoStruct2 := demoStruct2{"demo"}
	aDemoStruct2.Parse("demo2")
	result2, err2 := aParser.parseFixture("new value 2", reflect.TypeOf(aDemoStruct2))
	assert.Equals(t, nil, err2, "no err2")
	assert.Equals(t, "new value 2", result2.(demoStruct2).message, "value set in parse demo2")
	assert.Equals(t, "demo2", aDemoStruct2.message, "original value not changed demo2")
}

func TestParserParseHTMLTable(t *testing.T) {
	assertParseHTMLTable := func(input string, hasError bool, expected string, description string) {
		result, err := parseHTMLTable(input)
		if hasError {
			assert.IsTrue(t, err != nil, fmt.Sprintf("error not nil for %v", description))
			assert.Equals(t, expected, err.Error(), fmt.Sprintf("%v: error value", description))
		} else {
			assert.IsTrue(t, err == nil, fmt.Sprintf("%v: error nil", description))
			assert.Equals(t, expected, fmt.Sprintf("%v", result), fmt.Sprintf("%v: result", description))
		}
	}
	assertParseHTMLTable("<table><tr><td>id</td><td>123</td></tr><tr><td>name</td><td>Charlie</td></tr></table>",
		false, "[[id 123] [name Charlie]]", "table without fixtures and without spaces")
	tableString :=
		"<table class=\"hash_table\">\n" +
			"  <tr class=\"hash_row\">\n    <td class=\"hash_key\">id</td>\n    <td class=\"hash_value\">321</td>\n  </tr>\n" +
			"  <tr class=\"hash_row\">\n    <td class=\"hash_key\">name</td>\n    <td class=\"hash_value\">Parker</td>\n  </tr>\n</table>"
	assertParseHTMLTable(tableString, false, "[[id 321] [name Parker]]", "table with class attributes and spaces")
	assertParseHTMLTable("bogus", true, "Could not parse 'bogus' as an HTML table", "no table")
	assertParseHTMLTable("<table>bogus</table>", false, "[]", "table without valid tr/td tag returns empty")
	assertParseHTMLTable("<table><tr>bogus</tr></table>", false, "[]", "no td tags")
}

func TestParserParseMap(t *testing.T) {
	aParser := newParser(newSymbolTable())
	aMap := make(map[rune]uint)
	input1 := "<table><tr><td>1</td><td>5</td></tr><tr><td>4</td><td>25</td><tr></table>"
	result1, err1 := aParser.parseMap(input1, reflect.TypeOf(aMap))
	assert.Equals(t, nil, err1, "ParseMap err")
	assert.Equals(t, uint(5), result1.(map[rune]uint)[1], "ParseMap First entry match")
	assert.Equals(t, uint(25), result1.(map[rune]uint)[4], "ParseMap Second entry match")
	_, err2 := aParser.parseMap("1", reflect.TypeOf(aMap))
	assert.IsTrue(t, nil != err2, "err2 != nil")
	assert.Equals(t, "'1' is not a valid specification for 'map[int32]uint'", err2.Error(), "ParseMap invalid input")
	_, err3 := aParser.parseMap("<table><tr><td>1</td></tr></table>", reflect.TypeOf(aMap))
	assert.Equals(t, "row '[1]' in hash 'map[int32]uint' does not have two cells", err3.Error(), "ParseMap row not a list")
	_, err4 := aParser.parseMap("<table><tr><td>a</td><td>b</td></tr></table>", reflect.TypeOf(aMap))
	assert.Equals(t, "Could not parse key 'a' in hash 'map[int32]uint'", err4.Error(), "ParseMap wrong key type")
	_, err5 := aParser.parseMap("<table><tr><td>1</td><td>b</td></tr></table>", reflect.TypeOf(aMap))
	assert.Equals(t, "Could not parse value 'b' in hash 'map[int32]uint'", err5.Error(), "ParseMap wrong value type")
	input6 := "<table class=\"hash_table\">\n" +
		"  <tr class=\"hash_row\">\n    <td class=\"hash_key\">5</td>\n    <td class=\"hash_value\">25</td>\n  </tr>\n" +
		"  <tr class=\"hash_row\">\n    <td class=\"hash_key\">10</td>\n    <td class=\"hash_value\">100</td>\n  </tr>\n</table>"
	result6, err6 := aParser.parseMap(input6, reflect.TypeOf(aMap))
	assert.Equals(t, nil, err6, "No error in mapping from HTML string")
	outMap := result6.(map[rune]uint)
	assert.Equals(t, 2, len(outMap), "ParseMam from HTLM len OK")
	assert.Equals(t, uint(25), outMap[5], "ParseMap from HTML First entry match")
	assert.Equals(t, uint(100), outMap[10], "ParseMap from HTML Second entry match")
}

func TestParserParsePredefined(t *testing.T) {
	assertParsePredefined := func(expected interface{}, valueToParse string, targetType reflect.Type, description string) {
		aParser := newParser(newSymbolTable())
		response, err := aParser.parsePredefined(valueToParse, targetType)
		assert.Equals(t, nil, err, fmt.Sprintf("%v error", description))
		assert.Equals(t, expected, response, fmt.Sprintf("%v value", description))
	}
	assertParsePredefined(2, "2", reflect.TypeOf(1), "int")
	assertParsePredefined(uint32(4095), "0o7777", reflect.TypeOf(uint32(3)), "uint32")
	assertParsePredefined(rune(51966), "0xCAFE", reflect.TypeOf(rune(0)), "rune")
	assertParsePredefined(byte(10), "0b1010", reflect.TypeOf(byte(0)), "float")
	assertParsePredefined(3.14, "3.14", reflect.TypeOf(2.0), "float")
	assertParsePredefined(false, "false", reflect.TypeOf(true), "bool")
	assertParsePredefined("text", "text", reflect.TypeOf(""), "string")
	aParser := newParser(newSymbolTable())
	_, err := aParser.parsePredefined("1+j", reflect.TypeOf(complex128(1)))
	assert.Equals(t, "Could not convert '1+j' to type 'complex128'", err.Error(), "complex error")
}

func TestParserParsePtr(t *testing.T) {
	aParser := newParser(newSymbolTable())
	aParser.caller = newFunctionCaller(aParser)
	aDemoStruct1 := new(demoStruct1)
	result, err := aParser.parsePtr("text2", reflect.TypeOf(aDemoStruct1))
	assert.Equals(t, nil, err, fmt.Sprintf("%v error", "parsePtr Err"))
	assert.Equals(t, "text2", result.(*demoStruct1).ToString(), "ParsePtr value")
	anEmptyStruct := new(emptyStruct)
	_, err = aParser.parsePtr("ok", reflect.TypeOf(anEmptyStruct))
	assert.IsTrue(t, err != nil, "Err = nil")
	assert.Equals(t, "No method Parse found for type '*slimprocessor.emptyStruct'", err.Error(), "ParsePtr Err Empty")
}

func TestParserParseEmptyAndWrongSlices(t *testing.T) {
	aParser := newParser(newSymbolTable())
	slice := []int{}
	result, err1 := aParser.parseSlice("[]", reflect.TypeOf(slice))
	assert.Equals(t, nil, err1, "no err")
	sliceOut := result.([]int)
	assert.Equals(t, 0, len(sliceOut), "len")
	_, err2 := aParser.parseSlice("", reflect.TypeOf(slice))
	assert.IsTrue(t, err2 != nil, "err2 != nil")
	assert.Equals(t, "'' is not an array", err2.Error(), "wrong array")
	_, err3 := aParser.parseSlice("[1", reflect.TypeOf(slice))
	assert.IsTrue(t, err2 != nil, "err3 != nil")
	assert.Equals(t, "Could not find matching ']' in '[1'", err3.Error(), "wrong array")
	_, err4 := aParser.parseSlice("[[]", reflect.TypeOf(slice))
	assert.IsTrue(t, err2 != nil, "err4 != nil")
	assert.Equals(t, "Could not find matching ']' in '[[]'", err4.Error(), "wrong array")
}

func TestParserParse1DimensionalSlice(t *testing.T) {
	aParser := newParser(newSymbolTable())
	slice := []int{}
	result, err := aParser.parseSlice("[1, 2, 3, 5, 8]", reflect.TypeOf(slice))
	assert.Equals(t, nil, err, "No err")
	assert.Equals(t, "[1 2 3 5 8]", fmt.Sprintf("%v", result.([]int)), "result OK")
}

func TestParserParse1DimensionalSliceWithWrongValue(t *testing.T) {
	aParser := newParser(newSymbolTable())
	slice := []int{}
	_, err := aParser.parseSlice("[1, 2, a]", reflect.TypeOf(slice))
	assert.Equals(t, "Can't parse 'a' as element for slice '[]int'", err.Error(), "error message")
}

func TestParserParse2DimensionalSlice(t *testing.T) {
	aParser := newParser(newSymbolTable())
	slice := [][]int{}
	result, err := aParser.parseSlice("[[1, 2], [3, 5], [8, 13]]", reflect.TypeOf(slice))
	assert.Equals(t, nil, err, "No err")
	assert.Equals(t, "[[1 2] [3 5] [8 13]]", fmt.Sprintf("%v", result.([][]int)), "result OK")
}

func TestParserParseToInferredType(t *testing.T) {
	aParser := newParser(newSymbolTable())
	assert.Equals(t, int64(5), aParser.parseToInferredType("5"), "Int")
	maxUint := ^uint(0)
	assert.Equals(t, uint64(maxUint), aParser.parseToInferredType(fmt.Sprintf("%v", maxUint)), "Uint")
	assert.Equals(t, float64(3.14), aParser.parseToInferredType("3.14"), "Float64")
	assert.Equals(t, false, aParser.parseToInferredType("False"), "Bool")
	assert.Equals(t, "q", aParser.parseToInferredType("q"), "String")
}
