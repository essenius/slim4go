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

// TODO: this file is pretty large. See if we can split it up logically

package slimprocessor

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/essenius/slim4go/internal/assert"
	"github.com/essenius/slim4go/internal/slimentity"
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

func NewObjectWithPanic() int {
	panic("Object creation failed")
}

func NewOrderWithParams(productID string, unitPrice float64, units int) {}

func initParser() *Parser {
	parser := NewParser(NewSymbolTable())
	serializer := NewObjectHandler(parser)
	parser.SetObjectSerializer(serializer)
	return parser
}

func TestParserCallFunction(t *testing.T) {
	parser := initParser()
	result1, err1 := parser.CallFunction(reflect.ValueOf(NewObjectWithPanic), []string{})
	assert.Equals(t, "Panic: Object creation failed", err1.Error(), "Panicking function")
	assert.Equals(t, nil, result1, "No result with panic")
	messengerInstance, err2 := parser.CallFunction(reflect.ValueOf(NewMessenger), []string{})
	assert.Equals(t, nil, err2, "No error calling function")
	assert.Equals(t, "*slimprocessor.Messenger", reflect.TypeOf(messengerInstance).String(), "Type of instance OK")
	result3, err3 := parser.CallFunction(reflect.ValueOf(NewMessenger), []string{"q"})
	assert.Equals(t, "Expected 0 parameter(s) but got 1", err3.Error(), "Create messenger with wrong parameter")
	assert.Equals(t, "", result3, "No result with parameter error")
}

func TestParserIsPredefined(t *testing.T) {
	assertPredefined := func(isPredefined bool, value interface{}, description string) {
		assert.Equals(t, isPredefined, isPredefinedType(reflect.TypeOf(value)), description)
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

func TestParserMatchParamType(t *testing.T) {
	method := reflect.ValueOf(NewOrderWithParams)
	assert.IsTrue(t, method.IsValid(), "method valid")
	params := []string{"test", "100"}
	symbols := NewSymbolTable()
	parser := NewParser(symbols)
	_, err := parser.matchParamType(params, method)
	assert.Equals(t, "Expected 3 parameter(s) but got 2", err.Error(), "Wrong number of parameters")
	params = append(params, "25")
	result, err := parser.matchParamType(params, method)
	assert.Equals(t, nil, err, "No error")
	assert.Equals(t, "string", (*result)[0].Type().String(), "type of param 0 is string")
	assert.Equals(t, "test", (*result)[0].Interface(), "value of param 0 is test")
	assert.Equals(t, "float64", (*result)[1].Type().String(), "type of param 1 is float")
	assert.Equals(t, 100.0, (*result)[1].Interface(), "value of param 1 is 100.0")
	assert.Equals(t, "int", (*result)[2].Type().String(), "type of param 2 is int")
	assert.Equals(t, 25, (*result)[2].Interface(), "value of param 2 is 25")
	params = []string{"test", "100", "q"}
	_, err = parser.matchParamType(params, method)
	assert.Equals(t, "Could not convert 'q' to type 'int'", err.Error(), "invalid conversion")
}

func TestParserMatchParamTypeVariadic(t *testing.T) {
	sumFunction := func(num ...int) int {
		sum := 0
		for _, value := range num {
			sum += value
		}
		return sum
	}
	printFunction := func(template string, args ...interface{}) string {
		return fmt.Sprintf(template, args...)
	}

	symbols := NewSymbolTable()
	parser := NewParser(symbols)

	args := []string{"2", "3", "5"}
	result1, err1 := parser.matchParamType(args, reflect.ValueOf(sumFunction))
	assert.Equals(t, nil, err1, "variadic1: No error")
	assert.Equals(t, 3, len(*result1), "variadic1: 3 params")
	assert.Equals(t, "int", (*result1)[0].Type().String(), "variadic: type of param is int")
	assert.Equals(t, 2, (*result1)[0].Interface().(int), "first param value is 2")
	assert.Equals(t, 3, (*result1)[1].Interface().(int), "second param value is 3")
	assert.Equals(t, 5, (*result1)[2].Interface().(int), "thirdd param value is 5")

	emptyArgs := []string{}
	result2, err2 := parser.matchParamType(emptyArgs, reflect.ValueOf(sumFunction))
	assert.Equals(t, nil, err2, "variadic2 empty: No error")
	assert.Equals(t, 0, len(*result2), "variadic2 empty: 0 params")

	args3 := []string{"param %v %v", "3", "5.5"}
	result3, err3 := parser.matchParamType(args3, reflect.ValueOf(printFunction))
	assert.Equals(t, nil, err3, "Variadic3 param interface: No error")
	assert.Equals(t, 3, len(*result3), "Variadic3 param interface: 3 params")
	assert.Equals(t, "string", (*result3)[0].Type().String(), "Variadic3 param interface: type of param[0] is string")
	assert.Equals(t, "int64", (*result3)[1].Type().String(), "Variadic3 param interface: type of param[0] is interface{}")
	assert.Equals(t, "float64", (*result3)[2].Type().String(), "Variadic3 param interface: type of param[0] is interface{}")
	assert.Equals(t, "param %v %v", (*result3)[0].Interface().(string), "first param value is ok")
	assert.Equals(t, int64(3), (*result3)[1].Interface(), "second param value is 3")
	assert.Equals(t, 5.5, (*result3)[2].Interface(), "third param value is 5.5")
}

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
	assertParse := func(parser *Parser, expectation interface{}, source string, target interface{}, description string) {
		result, err := parser.Parse(source, reflect.ValueOf(target).Type())
		assert.Equals(t, nil, err, fmt.Sprintf("%v NoError", description))
		assert.Equals(t, expectation, result, fmt.Sprintf("%v Value", description))
	}
	parser := initParser()
	parser.symbols.Set("test", "text")
	aDemoStruct1 := new(demoStruct1)
	aDemoStruct1.Parse("qwe123")
	parser.symbols.Set("obj", aDemoStruct1)
	assertParse(parser, false, "False", ok, "bool")
	assertParse(parser, 23, "23", i, "int")
	assertParse(parser, 3.14, "3.14", f, "float64")
	assertParse(parser, uint32(4095), "0o7777", ui, "uint32")
	assertParse(parser, rune(51966), "0xCAFE", r, "rune")
	assertParse(parser, byte(10), "0b1010", b, "byte")
	assertParse(parser, "text", "text", s, "string")

	var c complex128
	_, err := parser.Parse("1+j", reflect.ValueOf(c).Type())
	assert.Equals(t, "Could not convert '1+j' to type 'complex128'", err.Error(), "Unable to convert")
	result, err2 := parser.Parse("$obj", reflect.TypeOf(aDemoStruct1))
	assert.Equals(t, nil, err2, "err2 is nil with retrieving obj")
	assert.Equals(t, "qwe123", result.(*demoStruct1).ToString(), "Get data from retieved non-text symbol")
	anEmptyStruct := new(emptyStruct)
	_, err3 := parser.Parse("$obj", reflect.TypeOf(anEmptyStruct))
	assert.Equals(t, "Symbol '$obj' of type '*slimprocessor.demoStruct1' not assignable to type '*slimprocessor.emptyStruct'", err3.Error(), "Symbol content not assignable to type")
	aMap := make(map[rune]uint)
	outMap, err4 := parser.Parse("<table></table>", reflect.TypeOf(aMap))
	assert.IsTrue(t, err4 == nil, "no error expected with empty map")
	assert.Equals(t, 0, len(outMap.(map[rune]uint)), "map len=0")
	// testing both parsePtr and parseSlice
	aSlice := &[]string{}
	outSlice, err5 := parser.Parse("[]", reflect.TypeOf(aSlice))
	assert.IsTrue(t, err5 == nil, "no error expected with slice")
	assert.Equals(t, 0, len(*outSlice.(*[]string)), "slice len=0")
	func1 := TestParserParse
	_, err6 := parser.Parse("[]", reflect.TypeOf(func1))
	assert.Equals(t, "Don't know how to resolve '[]' into 'func(*testing.T)'", err6.Error(), "error expected with func")
}

func TestParserParseFixture(t *testing.T) {
	aDemoStruct1 := new(demoStruct1)
	aDemoStruct1.Parse("demo1")
	parser := initParser()
	result1, err1 := parser.parseFixture("new value 1", reflect.TypeOf(aDemoStruct1))
	assert.Equals(t, nil, err1, "no err1")
	assert.Equals(t, "new value 1", result1.(*demoStruct1).message, "value set in parse demo1")
	assert.Equals(t, "demo1", aDemoStruct1.message, "original value not changed demo1")
	aDemoStruct2 := demoStruct2{"demo"}
	aDemoStruct2.Parse("demo2")
	result2, err2 := parser.parseFixture("new value 2", reflect.TypeOf(aDemoStruct2))
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
	parser := NewParser(NewSymbolTable())
	aMap := make(map[rune]uint)
	input1 := "<table><tr><td>1</td><td>5</td></tr><tr><td>4</td><td>25</td><tr></table>"
	result1, err1 := parser.parseMap(input1, reflect.TypeOf(aMap))
	assert.Equals(t, nil, err1, "ParseMap err")
	assert.Equals(t, uint(5), result1.(map[rune]uint)[1], "ParseMap First entry match")
	assert.Equals(t, uint(25), result1.(map[rune]uint)[4], "ParseMap Second entry match")
	_, err2 := parser.parseMap("1", reflect.TypeOf(aMap))
	assert.IsTrue(t, nil != err2, "err2 != nil")
	assert.Equals(t, "'1' is not a valid specification for 'map[int32]uint'", err2.Error(), "ParseMap invalid input")
	_, err3 := parser.parseMap("<table><tr><td>1</td></tr></table>", reflect.TypeOf(aMap))
	assert.Equals(t, "row '[1]' in hash 'map[int32]uint' does not have two cells", err3.Error(), "ParseMap row not a list")
	_, err4 := parser.parseMap("<table><tr><td>a</td><td>b</td></tr></table>", reflect.TypeOf(aMap))
	assert.Equals(t, "Could not parse key 'a' in hash 'map[int32]uint'", err4.Error(), "ParseMap wrong key type")
	_, err5 := parser.parseMap("<table><tr><td>1</td><td>b</td></tr></table>", reflect.TypeOf(aMap))
	assert.Equals(t, "Could not parse value 'b' in hash 'map[int32]uint'", err5.Error(), "ParseMap wrong value type")
	input6 := "<table class=\"hash_table\">\n" +
		"  <tr class=\"hash_row\">\n    <td class=\"hash_key\">5</td>\n    <td class=\"hash_value\">25</td>\n  </tr>\n" +
		"  <tr class=\"hash_row\">\n    <td class=\"hash_key\">10</td>\n    <td class=\"hash_value\">100</td>\n  </tr>\n</table>"
	result6, err6 := parser.parseMap(input6, reflect.TypeOf(aMap))
	assert.Equals(t, nil, err6, "No error in mapping from HTML string")
	outMap := result6.(map[rune]uint)
	assert.Equals(t, 2, len(outMap), "ParseMam from HTLM len OK")
	assert.Equals(t, uint(25), outMap[5], "ParseMap from HTML First entry match")
	assert.Equals(t, uint(100), outMap[10], "ParseMap from HTML Second entry match")
}

func TestParserParsePredefined(t *testing.T) {
	assertParsePredefined := func(expected interface{}, valueToParse string, targetType reflect.Type, description string) {
		parser := NewParser(NewSymbolTable())
		response, err := parser.parsePredefined(valueToParse, targetType)
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
	parser := NewParser(NewSymbolTable())
	_, err := parser.parsePredefined("1+j", reflect.TypeOf(complex128(1)))
	assert.Equals(t, "Could not convert '1+j' to type 'complex128'", err.Error(), "complex error")
}

func TestParserParsePtr(t *testing.T) {
	parser := initParser()
	aDemoStruct1 := new(demoStruct1)
	result, err := parser.parsePtr("text2", reflect.TypeOf(aDemoStruct1))
	assert.Equals(t, nil, err, fmt.Sprintf("%v error", "parsePtr Err"))
	assert.Equals(t, "text2", result.(*demoStruct1).ToString(), "ParsePtr value")
	anEmptyStruct := new(emptyStruct)
	_, err = parser.parsePtr("ok", reflect.TypeOf(anEmptyStruct))
	assert.IsTrue(t, err != nil, "Err = nil")
	assert.Equals(t, "No method Parse found for type '*slimprocessor.emptyStruct'", err.Error(), "ParsePtr Err Empty")
}

func TestParserParseEmptyAndWrongSlices(t *testing.T) {
	parser := NewParser(NewSymbolTable())
	slice := []int{}
	result, err1 := parser.parseSlice("[]", reflect.TypeOf(slice))
	assert.Equals(t, nil, err1, "no err")
	sliceOut := result.([]int)
	assert.Equals(t, 0, len(sliceOut), "len")
	_, err2 := parser.parseSlice("", reflect.TypeOf(slice))
	assert.IsTrue(t, err2 != nil, "err2 != nil")
	assert.Equals(t, "'' is not an array", err2.Error(), "wrong array")
	_, err3 := parser.parseSlice("[1", reflect.TypeOf(slice))
	assert.IsTrue(t, err2 != nil, "err3 != nil")
	assert.Equals(t, "Could not find matching ']' in '[1'", err3.Error(), "wrong array")
	_, err4 := parser.parseSlice("[[]", reflect.TypeOf(slice))
	assert.IsTrue(t, err2 != nil, "err4 != nil")
	assert.Equals(t, "Could not find matching ']' in '[[]'", err4.Error(), "wrong array")
}

func TestParserParse1DimensionalSlice(t *testing.T) {
	parser := NewParser(NewSymbolTable())
	slice := []int{}
	result, err := parser.parseSlice("[1, 2, 3, 5, 8]", reflect.TypeOf(slice))
	assert.Equals(t, nil, err, "No err")
	assert.Equals(t, "[1 2 3 5 8]", fmt.Sprintf("%v", result.([]int)), "result OK")
}

func TestParserParse1DimensionalSliceWithWrongValue(t *testing.T) {
	parser := NewParser(NewSymbolTable())
	slice := []int{}
	_, err := parser.parseSlice("[1, 2, a]", reflect.TypeOf(slice))
	assert.Equals(t, "Can't parse 'a' as element for slice '[]int'", err.Error(), "error message")
}

func TestParserParse2DimensionalSlice(t *testing.T) {
	parser := NewParser(NewSymbolTable())
	slice := [][]int{}
	result, err := parser.parseSlice("[[1, 2], [3, 5], [8, 13]]", reflect.TypeOf(slice))
	assert.Equals(t, nil, err, "No err")
	assert.Equals(t, "[[1 2] [3 5] [8 13]]", fmt.Sprintf("%v", result.([][]int)), "result OK")
}

func TestParserParseToInferredType(t *testing.T) {
	parser := NewParser(NewSymbolTable())
	assert.Equals(t, int64(5), parser.parseToInferredType("5"), "Int")
	maxUint := ^uint(0)
	assert.Equals(t, uint64(maxUint), parser.parseToInferredType(fmt.Sprintf("%v", maxUint)), "Uint")
	assert.Equals(t, float64(3.14), parser.parseToInferredType("3.14"), "Float64")
	assert.Equals(t, false, parser.parseToInferredType("False"), "Bool")
	assert.Equals(t, "q", parser.parseToInferredType("q"), "String")
}

func TestParserReplaceSymbol(t *testing.T) {
	parser := initParser()
	assert.Equals(t, nil, parser.symbols.Set("test1", "value1"), "Set Symbol test1 to string")
	assert.Equals(t, nil, parser.symbols.Set("test2", "value2"), "Set Symbol test2 to string")
	aDemoStruct1 := new(demoStruct1)
	aDemoStruct1.Parse("hi from aDemoStruct1")
	assert.Equals(t, nil, parser.symbols.Set("test3", aDemoStruct1), "Set Symbol test3 to object")
	assert.Equals(t, "Invalid symbol name: $_test3", parser.symbols.Set("$_test3", "_value3").Error(), "invalid name $_test3")

	assert.Equals(t, "value2", parser.replaceSymbolValue("$test2"), "ReplaceValue($test2) returns value=2")
	assert.Equals(t, "$_test3", parser.replaceSymbolValue("$_test3"), "ReplaceValue($_test3) returns the input")
	assert.Equals(t, "value1", parser.replaceSymbolValue("$test1"), "ReplaceValue($test1) returns value1")
	assert.Equals(t, "hi from aDemoStruct1", parser.replaceSymbolValue("$test3"), "ReplaceValue($test3) returns ToString() value")
	assert.Equals(t, "$test4", parser.replaceSymbolValue("$test4"), "Nonexisting symbol returns the input")
	assert.Equals(t, "$", parser.ReplaceSymbolsIn("$"), "just a $ returns input")
	assert.Equals(t, "1€$", parser.ReplaceSymbolsIn("1€$"), "ending with $ returns input")

	assert.Equals(t, "we see value1, value2, 'hi from aDemoStruct1', $test4 and $_test3 returned as-is",
		parser.ReplaceSymbolsIn("we see $test1, $test2, '$test3', $test4 and $_test3 returned as-is"), "3 replacements")
	listIn := slimentity.NewSlimListContaining([]slimentity.SlimEntity{"$$test1", "this is $test2", "$test3", "$test1-$"})
	listOut := parser.ReplaceSymbols(listIn).(*slimentity.SlimList)
	assert.Equals(t, `[$value1, this is value2, hi from aDemoStruct1, value1-$]`, listOut.ToString(), "list values replaced OK")
}
