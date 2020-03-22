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
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

const absoluteZeroInCelsius float64 = -273.15
const absoluteZeroInFahrenheit float64 = -459.67
const parseError = "Expected float with suffix F, C or K but got '%v'"
const scaleError = "Unrecognized temperature scale: %v"

// NewTemperatureFactory returns a Temperature Factory
func NewTemperatureFactory() *TemperatureFactory {
	return new(TemperatureFactory)
}

// TemperatureFactory is an example fixture factory.
type TemperatureFactory struct{}

// NewTemperatureConverter creates a TemperatureConverter.
func (factory *TemperatureFactory) NewTemperatureConverter() *TemperatureConverter {
	return new(TemperatureConverter)
}

// NewTemperature creates a Temperature.
func (factory *TemperatureFactory) NewTemperature(input string) *Temperature {
	temperature := new(Temperature)
	temperature.Parse(input)
	return temperature
}

// TemperatureConverter shows how to use objects as parameters.
type TemperatureConverter struct{}

// ConvertTo converts temperatures between scales.
func (temperatureConverter *TemperatureConverter) ConvertTo(input *Temperature, scale string) float64 {
	return input.ValueIn(scale)
}

// Temperature is an example parsable object.
type Temperature struct {
	value float64
}

// ToString serializes a Temperature.
func (temperature *Temperature) ToString() string {
	return fmt.Sprintf("%v K", temperature.value)
}

// Parse deserializes a string into a Temperature.
func (temperature *Temperature) Parse(input string) {
	if input == "" {
		panic(fmt.Errorf(parseError, ""))
	}
	scale := input[len(input)-1:]
	baseValue := strings.TrimSpace(input[:len(input)-1])
	temperatureValue, err := strconv.ParseFloat(baseValue, 64)
	if err != nil {
		panic(fmt.Errorf(parseError, input))
	}
	switch scale {
	case "F":
		temperature.value = math.Round(10000.0*(temperatureValue-absoluteZeroInFahrenheit)*5.0/9.0) / 10000.0
	case "C":
		temperature.value = temperatureValue - absoluteZeroInCelsius
	case "K":
		temperature.value = temperatureValue
	default:
		panic(fmt.Errorf(parseError, input))
	}
}

// ValueIn returns the temperature value in the required scale (F, C or K).
func (temperature *Temperature) ValueIn(scale string) float64 {
	if scale == "" {
		panic(fmt.Errorf(scaleError, ""))
	}
	switch unicode.ToUpper(rune(scale[0])) {
	case 'F':
		return math.Round(10000.0*(temperature.value*9.0/5.0+absoluteZeroInFahrenheit)) / 10000.0
	case 'C':
		return temperature.value + absoluteZeroInCelsius
	case 'K':
		return temperature.value
	default:
		panic(fmt.Errorf(scaleError, scale))
	}
}
