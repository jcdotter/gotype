// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	"math"
	"unsafe"
)

// ------------------------------------------------------------ /
// GOTYPE CUSTOM TYPE IMPLEMENTATION
// implementation of custom type of encoding/json
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

type JSON []byte

// JsonOf returns 'a' as a gotype Json
// panics if cannot convert 'a' to Json
func JsonOf(a any) JSON {
	if j, isjson := a.(JSON); isjson {
		return j
	}
	switch j := a.(type) {
	case JSON:
		return j
	case []byte:
		return JSON(j)
	default:
		panic("cannot convert value to JSON")
	}
}

// IsJsonable evaluate whether a can be converted to Json
func IsJsonable(a any) bool {
	var j JSON
	switch (*VALUE)(unsafe.Pointer(&a)).Kind() {
	case Struct, Array, Slice, Map:
		return true
	default:
		switch aa := a.(type) {
		case JSON:
			j = aa
		case []byte:
			j = JSON(aa)
		case STRING:
			j = JSON(aa)
		case string:
			j = JSON(aa)
		case []rune:
			j = make(JSON, len(aa))
			for i, r := range aa {
				j[i] = uint8(r)
			}
		default:
			return false
		}
	}
	return hasJsonBookends(j)
}

func hasJsonBookends(b []byte) bool {
	return (b[0] == 91 && b[len(b)-1] == 93) ||
		(b[0] == 123 && b[len(b)-1] == 125)
}

// JSON returns VALUE as gotype JSON
func (v VALUE) JSON() JSON {
	switch v.KIND() {
	case Array:
		return (ARRAY)(v).JSON()
	case Map:
		return (MAP)(v).JSON()
	case Pointer:
		return v.ElemDeep().JSON()
	case String:
		return (*STRING)(v.ptr).JSON()
	case Slice:
		return (SLICE)(v).JSON()
	case Struct:
		return (STRUCT)(v).JSON()
	case Bytes:
		return (*BYTES)(v.ptr).JSON()
	default:
		panic("cannot convert value to array")
	}
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Natural returns underlyding value of gotype JSON as a golang interface{}
func (j JSON) Native() []byte {
	return j
}

// Interface returns gotype JSON as a golang interface{}
func (j JSON) Interface() any {
	return j
}

// Value returns gotype JSON as gotype VALUE
func (j JSON) VALUE() VALUE {
	return ValueOf([]byte(j))
}

// Bytes returns gotype JSON as []byte
func (j JSON) Bytes() []byte {
	return []byte(j)
}

// Bool returns gotype JSON as bool
// false if empty, true if a len > 0
func (j JSON) Bool() bool {
	return len(j) > 2
}

// ARRAY returns gotype JSON as gotype ARRAY
func (j JSON) ARRAY() ARRAY {
	return ArrayOf(j.Slice())
}

// String returns gotype JSON as string
func (j JSON) String() string {
	return string(j)
}

func (j JSON) STRING() STRING {
	return STRING(j)
}

// Slice returns gotype JSON as []string
func (j JSON) Slice() []any {
	_, l := STRING(j).Unserialize()
	return l
}

// SLICE returns gotype JSON as gotype SLICE
func (j JSON) SLICE() SLICE {
	return SliceOf(j.Slice())
}

// Map returns gotype JSON as map[string]string
func (j JSON) Map() map[string]any {
	o, _ := STRING(j).Unserialize()
	return o
}

// MAP returns gotype JSON as MAP
func (j JSON) MAP() MAP {
	return MapOf(j.Map())
}

// Scan reads gotype JSON into the struct provided
// if tag is empty, Field names will be used to read JSON keys into dest Struct
func (j JSON) Scan(dest any, tags ...string) {
	j.MAP().Scan(dest, tags...)
}

// Empty evaluates whether the JSON holds any values
func (j JSON) Empty() bool {
	return j.Bool()
}

// Unserialize converts a gotype serialized JSON object to a map or slice, respectively
func (j JSON) Unserialize() (object map[string]any, list []any) {
	return STRING(j).Unserialize()
}

// Format retenders JSON as string with line breaks and indents,
// indents representing the number of spaces per indent and
// breaking each element in a list onto its own line if breakList is true
func (j JSON) Format(indent int, breakList ...bool) string {
	var q, ins string
	var brkList bool
	str := []string{}
	if len(breakList) > 0 && breakList[0] {
		brkList = true
	}
	in := STRING(" ").Repeat(indent)
	l := len(j)
	for i := 0; i < l; i++ {
		b := j[i]
		switch b {
		case 34:
			q, i = j.__quote(i, l, 34)
		case 44:
			q = ",\n" + ins
		case 58:
			q = ":\t"
		case 91:
			if !brkList {
				q, i = j.__list(i, l, 34)
			} else {
				q, ins = j.__indentInc("[", ins, in)
			}
		case 93:
			q, ins = j.__indentDec("]", ins, indent)
		case 123:
			q, ins = j.__indentInc("{", ins, in)
		case 125:
			q, ins = j.__indentDec("}", ins, indent)
		default:
			q = string(b)
		}
		str = append(str, q)
	}
	return JoinStrings(str, "")
}

func (j JSON) __quote(start int, l int, q uint8) (string, int) {
	for i := start + 1; i < l; i++ {
		if j[i] == q && j[i-1] != 92 {
			return string(j[start : i+1]), i
		}
	}
	return string(j[start:l]), l
}

func (j JSON) __list(start int, l int, q uint8) (string, int) {
	for i := start + 1; i < l; i++ {
		b := j[i]
		if b == q {
			i++
			for i < l {
				if j[i] == q && j[i-1] != 92 {
					break
				}
				i++
			}
		} else if b == 93 {
			return string(j[start : i+1]), i
		} else if b == 91 {
			_, i = j.__list(i, l, 34)
		}
	}
	return string(j[start:l]), l
}

func (j JSON) __indentInc(char, curIndent, indent string) (str, ins string) {
	ins = curIndent + indent
	str = char + "\n" + ins
	return
}

func (j JSON) __indentDec(char, curIndent string, indent int) (str, ins string) {
	ins = curIndent[:int(math.Max(0, float64(len(curIndent)-indent)))]
	str = "\n" + ins + char
	return
}
