// Copyright 2023 james dotter. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import (
	"fmt"
	"math"
	"testing"
	"time"

	test "github.com/jcdotter/gtest"
)

var config = &test.Config{
	//PrintTest:   true,
	PrintFail:   true,
	PrintTrace:  true,
	PrintDetail: true,
	FailFatal:   true,
	Msg:         "%s",
}

func TestTest(t *testing.T) {
	s := STRING("testing")
	fmt.Println(TypeOf(s).Name())
	fmt.Println(TypeOf(s).NameShort())
}

func TestAll(t *testing.T) {
	TestValueOf(t)
	TestValueNew(t)
	TestValueLen(t)
	TestValueIndex(t)
	TestValueSerialize(t)
	TestValueSerializePrint(t)
	TestValueSetTyped(t)
	TestValueSetPtrPtr(t)
	TestValueSetPtrVal(t)
	TestValueSetUntyped(t)
	TestValueConversion(t)
	TestValueNewDeep(t)
	TestEncodeBasic(t)
	TestEncodeComplex(t)
}

func TestValueOf(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Interface()"
	for n, v := range getTestVars() {
		gt.Equal(ValueOf(v).Interface(), v, n)
	}
}

func TestIndirect(t *testing.T) {
	table := [][]string{{"Name", "Kind", "Type", "IfaceIndir", "flagIndir"}}
	for n, v := range getTestVars() {
		val := ValueOf(v)
		table = append(table, []string{n, val.Kind().String(), val.typ.String(), BOOL(val.typ.IfaceIndir()).String(), BOOL(val.flag&flagIndir != 0).String()})
	}
	SortByCol(table, 0)
	test.PrintTable(table, true)
}

func TestValueNew(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing VALUE.New(%s)"
	var (
		b  = true
		i  = 1
		s  = "test"
		a  = [2]string{"s", "s"}
		a1 = [1]string{"s"}
		l  = []string{"s", "s"}
		m  = map[string]string{"0": "s", "1": "s"}
		d  = string_struct{"s", "s"}
		d1 = string_struct_single{"s"}
	)

	gt.Equal(false, ValueOf(b).New().Elem().Interface(), "bool")
	gt.Equal(0, ValueOf(i).New().Elem().Interface(), "int")
	gt.Equal("", ValueOf(s).New().Elem().Interface(), "string")
	gt.Equal([2]string{"", ""}, ValueOf(a).New().Elem().Interface(), "array")
	gt.Equal([]string(nil), ValueOf(l).New().Elem().Interface(), "slice")
	gt.Equal(map[string]string(nil), ValueOf(m).New().Elem().Interface(), "map")
	gt.Equal(string_struct{}, ValueOf(d).New().Elem().Interface(), "struct")
	gt.Equal([1]string{""}, ValueOf(a1).New().Elem().Interface(), "array(1)")
	gt.Equal(string_struct_single{}, ValueOf(d1).New().Elem().Interface(), "struct(1)")

}

func TestValueLen(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing %s.Len()"
	var (
		s  = "test"
		a  = [2]string{"s", "s"}
		a1 = [1]string{"s"}
		l  = []string{"s", "s"}
		m  = map[string]string{"0": "s", "1": "s"}
		d  = string_struct{"s", "s"}
		d1 = string_struct_single{"s"}
	)
	gt.Equal(4, ValueOf(s).Len(), "string")
	gt.Equal(2, ValueOf(a).Len(), "array")
	gt.Equal(2, ValueOf(l).Len(), "slice")
	gt.Equal(2, ValueOf(m).Len(), "map")
	gt.Equal(2, ValueOf(d).Len(), "struct")
	gt.Equal(1, ValueOf(a1).Len(), "array(1)")
	gt.Equal(1, ValueOf(d1).Len(), "struct(1)")
}

func TestValueIndex(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Index(%d)"
	var (
		s  = "test"
		a  = [2]string{"s", "s"}
		l  = []string{"s", "s"}
		m  = map[string]string{"0": "s", "1": "s"}
		d  = string_struct{"s", "s"}
		a1 = [1]string{"s"}
		d1 = string_struct_single{"s"}
	)
	gt.Equal("t", ValueOf(s).Index(0).Interface(), "string", 0)
	gt.Equal("s", ValueOf(a).Index(0).Interface(), "array", 0)
	gt.Equal("s", ValueOf(l).Index(0).Interface(), "slice", 0)
	gt.Equal("s", ValueOf(m).Index(0).Interface(), "map", 0)
	gt.Equal("s", ValueOf(d).Index(0).Interface(), "struct", 0)
	gt.Equal("s", ValueOf(a1).Index(0).Interface(), "array(1)", 0)
	gt.Equal("s", ValueOf(d1).Index(0).Interface(), "struct(1)", 0)

	gt.Equal("e", ValueOf(&s).Elem().Index(1).Interface(), "*string", 1)
	gt.Equal("s", ValueOf(&a).Elem().Index(1).Interface(), "*array", 1)
	gt.Equal("s", ValueOf(&l).Elem().Index(1).Interface(), "*slice", 1)
	gt.Equal("s", ValueOf(&m).Elem().Index(1).Interface(), "*map", 1)
	gt.Equal("s", ValueOf(&d).Elem().Index(1).Interface(), "*struct", 1)
	gt.Equal("s", ValueOf(&a1).Elem().Index(0).Interface(), "*array(1)", 0)
	gt.Equal("s", ValueOf(&d1).Elem().Index(0).Interface(), "*struct(1)", 0)

	gt.Msg = "Testing ValueOf(%s).StructField(%s)"
	gt.Equal("s", ValueOf(&d).Elem().StructField("V2").Interface(), "*struct", "V2")
	gt.Equal("s", ValueOf(&d1).Elem().StructField("V1").Interface(), "*struct(1)", "V1")
}

func TestValueSerialize(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Serialize()"
	var (
		s  = "s"
		a  = [2]string{"s", "s"}
		l  = []string{"s", "s"}
		m  = map[string]string{"0": "s", "1": "s"}
		d  = string_struct{"s", "s"}
		a1 = [1]string{"s"}
		d1 = string_struct_single{"s"}

		pa  = [2]*string{&s, &s}
		pl  = []*string{&s, &s}
		pm  = map[string]*string{"0": &s, "1": &s}
		pd  = string_ptr_struct{&s, &s}
		pa1 = [1]*string{&s}
		pd1 = string_ptr_struct_single{&s}
	)
	gt.Equal(`"s"`, ValueOf(s).Serialize(), "string")
	gt.Equal(`["s","s"]`, ValueOf(a).Serialize(), "array")
	gt.Equal(`["s","s"]`, ValueOf(l).Serialize(), "slice")
	gt.Equal(`{"0":"s","1":"s"}`, ValueOf(m).Serialize(), "map")
	gt.Equal(`{"V1":"s","V2":"s"}`, ValueOf(d).Serialize(), "struct")
	gt.Equal(`["s"]`, ValueOf(a1).Serialize(), "array(1)")
	gt.Equal(`{"V1":"s"}`, ValueOf(d1).Serialize(), "struct(1)")

	gt.Equal(`"s"`, ValueOf(&s).Serialize(), "*string")
	gt.Equal(`["s","s"]`, ValueOf(&pa).Serialize(), "*array")
	gt.Equal(`["s","s"]`, ValueOf(&pl).Serialize(), "*slice")
	gt.Equal(`{"0":"s","1":"s"}`, ValueOf(&pm).Serialize(), "*map")
	gt.Equal(`{"V1":"s","V2":"s"}`, ValueOf(&pd).Serialize(), "*struct")
	gt.Equal(`["s"]`, ValueOf(&pa1).Serialize(), "*array(1)")
	gt.Equal(`{"V1":"s"}`, ValueOf(&pd1).Serialize(), "*struct(1)")
}

func TestValueSerializePrint(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Serialize() key: %s\n  value:\t%s"
	err := "%!v(PANIC=String method: runtime error: invalid memory address or nil pointer dereference)"
	for n, v := range getTestVars() {
		val := ValueOf(v)
		s := STRING(val.Serialize())
		gt.False(s == "" || s.Contains(err), val.typ, n, s)
	}
}

func TestValueSetTyped(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Set(%s)"
	vars := getTestVars()
	var (
		b  = false
		i  = 2
		s  = "updated"
		a  = [2]string{"updated", "updated"}
		l  = []string{"updated", "updated"}
		m  = map[string]string{"0": "updated", "1": "updated"}
		d  = string_struct{"updated", "updated"}
		a1 = [1]string{"updated"}
		d1 = string_struct_single{"updated"}
	)

	gt.Equal(b, ValueOf(vars["bool"]).Set(b).Interface(), "bool", "bool")
	gt.Equal(i, ValueOf(vars["int"]).Set(i).Interface(), "int", "int")
	gt.Equal(s, ValueOf(vars["string"]).Set(s).Interface(), "string", "string")
	gt.Equal(a, ValueOf(vars["[2]string"]).Set(a).Interface(), "array", "array")
	gt.Equal(l, ValueOf(vars["[]string{2}"]).Set(l).Interface(), "slice", "slice")
	gt.Equal(m, ValueOf(vars["map[string]string{2}"]).Set(m).Interface(), "map", "map")
	gt.Equal(d, ValueOf(vars["struct(string){2}"]).Set(d).Interface(), "struct", "struct")
	gt.Equal(a1, ValueOf(vars["[1]string"]).Set(a1).Interface(), "array(1)", "array(1)")
	gt.Equal(d1, ValueOf(vars["struct(string){1}"]).Set(d1).Interface(), "struct(1)", "struct(1)")

}

func TestValueSetPtrPtr(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Set(%s)"
	vars := createTestVars(1 == 2, 0, "false")
	var (
		b  = true
		i  = 2
		s  = "updated"
		a  = [2]string{"updated", "updated"}
		l  = []string{"updated", "updated"}
		m  = map[string]string{"0": "updated", "1": "updated"}
		d  = string_struct{"updated", "updated"}
		a1 = [1]string{"updated"}
		d1 = string_struct_single{"updated"}
	)

	gt.Equal(b, ValueOf(vars["*bool"]).Set(&b).Elem().Interface(), "*bool", "*bool")
	gt.Equal(i, ValueOf(vars["*int"]).Set(&i).Elem().Interface(), "*int", "*int")
	gt.Equal(s, ValueOf(vars["*string"]).Set(&s).Elem().Interface(), "*string", "*string")
	gt.Equal(a, ValueOf(vars["*[2]string"]).Set(&a).Elem().Interface(), "*array", "*array")
	gt.Equal(l, ValueOf(vars["*[]string{2}"]).Set(&l).Elem().Interface(), "*slice", "*slice")
	gt.Equal(m, ValueOf(vars["*map[string]string{2}"]).Set(&m).Elem().Interface(), "*map", "*map")
	gt.Equal(d, ValueOf(vars["*struct(string){2}"]).Set(&d).Elem().Interface(), "*struct", "*struct")
	gt.Equal(a1, ValueOf(vars["*[1]string"]).Set(&a1).Elem().Interface(), "*array(1)", "*array(1)")
	gt.Equal(d1, ValueOf(vars["*struct(string){1}"]).Set(&d1).Elem().Interface(), "*struct(1)", "*struct(1)")
}

func TestValueSetPtrVal(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Set(%s)"
	vars := getTestVars()
	var (
		b  = false
		i  = 2
		s  = "updated"
		a  = [2]string{"updated", "updated"}
		l  = []string{"updated", "updated"}
		m  = map[string]string{"0": "updated", "1": "updated"}
		d  = string_struct{"updated", "updated"}
		a1 = [1]string{"updated"}
		d1 = string_struct_single{"updated"}
	)

	gt.Equal(b, ValueOf(vars["*bool"]).Set(b).Elem().Interface(), "*bool", "bool")
	gt.Equal(i, ValueOf(vars["*int"]).Set(i).Elem().Interface(), "*int", "int")
	gt.Equal(s, ValueOf(vars["*string"]).Set(s).Elem().Interface(), "*string", "string")
	gt.Equal(a, ValueOf(vars["*[2]string"]).Set(a).Elem().Interface(), "*array", "array")
	gt.Equal(l, ValueOf(vars["*[]string{2}"]).Set(l).Elem().Interface(), "*slice", "slice")
	gt.Equal(m, ValueOf(vars["*map[string]string{2}"]).Set(m).Elem().Interface(), "*map", "map")
	gt.Equal(d, ValueOf(vars["*struct(string){2}"]).Set(d).Elem().Interface(), "*struct", "struct")
	gt.Equal(a1, ValueOf(vars["*[1]string"]).Set(a1).Elem().Interface(), "*array(1)", "array(1)")
	gt.Equal(d1, ValueOf(vars["*struct(string){1}"]).Set(d1).Elem().Interface(), "*struct(1)", "struct(1)")

}

func TestValueSetUntyped(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Set(%s)"
	var (
		o_bool    = false
		o_int     = 0
		o_uint    = uint(0)
		o_float   = 0.0
		o_array   = [2]string{"0", "0"}
		o_array1  = [1]string{"0"}
		o_slice   = []string{"0", "0"}
		o_map     = map[string]string{"0": "0", "1": "0"}
		o_struct  = string_struct{"0", "0"}
		o_struct1 = string_struct_single{"0"}
		o_time    = TIME(time.Unix(0, 0).UTC())
		o_uuid    = UUID([16]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})

		o_str_bool   = "false"
		o_str_int    = "0"
		o_str_uint   = "0"
		o_str_float  = "0"
		o_str_array  = `["0","0"]`
		o_str_slice  = `["0","0"]`
		o_str_map    = `{"0":"0","1":"0"}`
		o_str_struct = `{"V1":"0","V2":"0"}`
		o_str_time   = `1970-01-01 00:00:00.000000000`
		o_str_uuid   = `00000000-0000-0000-0000-000000000000`

		o_bytes_bool   = []byte{0x0}
		o_bytes_int    = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		o_bytes_uint   = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		o_bytes_float  = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		o_bytes_array  = []byte(o_str_array)
		o_bytes_slice  = []byte(o_str_slice)
		o_bytes_map    = []byte(o_str_map)
		o_bytes_struct = []byte(o_str_struct)
		o_bytes_time   = []byte{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		o_bytes_uuid   = []byte{0x26, 0x7b, 0x32, 0x29, 0x25, 0x66, 0x44, 0x26, 0xa8, 0x26, 0x8d, 0x80, 0x12, 0x6e, 0x71, 0x9a}
		o_bytes_string = []byte("true")

		u_bool    = true
		u_int     = 123
		u_uint    = uint(123)
		u_float   = 123.0
		u_float1  = 123.1
		u_array   = [2]string{"0", "123"}
		u_array1  = [1]string{"123"}
		u_slice   = []string{"0", "123"}
		u_map     = map[string]string{"0": "0", "1": "123"}
		u_struct  = string_struct{"0", "123"}
		u_struct1 = string_struct_single{"123"}
		u_time    = TIME(time.Unix(0, 123).UTC())
		u_uuid    = UUID([16]byte{0x26, 0x7b, 0x32, 0x29, 0x25, 0x66, 0x44, 0x26, 0xa8, 0x26, 0x8d, 0x80, 0x12, 0x6e, 0x71, 0x9a})

		u_str_bool   = "true"
		u_str_int    = "123"
		u_str_uint   = "123"
		u_str_float  = "123.0"
		u_str_float1 = "123.1"
		u_str_array  = `["0","123"]`
		u_str_slice  = `["0","123"]`
		u_str_map    = `{"0":"0","1":"123"}`
		u_str_struct = `{"V1":"0","V2":"123"}`
		u_str_time   = `1970-01-01 00:02:03.000`
		u_str_uuid   = `267b3229-2566-4426-a826-8d80126e719a`

		u_bytes_bool   = []byte{0x1}
		u_bytes_int    = []byte{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		u_bytes_uint   = []byte{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		u_bytes_float  = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0xc0, 0x5e, 0x40}
		u_bytes_array  = []byte(u_str_array)
		u_bytes_slice  = []byte(u_str_slice)
		u_bytes_map    = []byte(u_str_map)
		u_bytes_struct = []byte(u_str_struct)
		u_bytes_time   = []byte{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		u_bytes_uuid   = []byte{0x26, 0x7b, 0x32, 0x29, 0x25, 0x66, 0x44, 0x26, 0xa8, 0x26, 0x8d, 0x80, 0x12, 0x6e, 0x71, 0x9a}
		u_bytes_string = []byte("true")
		u_bytes_num    = []byte("123")

		u_int_bool   = 1
		u_uint_bool  = uint(1)
		u_float_bool = 1.0
		u_time_str   = TIME(time.Unix(123, 0).UTC())
	)

	// Test return values of Set()
	gt.Equal(u_bytes_bool, ValueOf(o_bytes_bool).Set(u_bytes_bool).Interface(), "bytes", "bytes")
	gt.Equal(u_bytes_int, ValueOf(o_bytes_int).Set(u_int).Interface(), "bytes", "int")
	gt.Equal(u_bytes_uint, ValueOf(o_bytes_uint).Set(u_uint).Interface(), "bytes", "uint")
	gt.Equal(u_bytes_float, ValueOf(o_bytes_float).Set(u_float).Interface(), "bytes", "float")
	gt.Equal(u_bytes_array, ValueOf(o_bytes_array).Set(u_array).Interface(), "bytes", "array")
	gt.Equal(u_bytes_slice, ValueOf(o_bytes_slice).Set(u_slice).Interface(), "bytes", "slice")
	gt.Equal(u_bytes_string, ValueOf(o_bytes_string).Set(u_str_bool).Interface(), "bytes", "string")
	gt.Equal(u_bytes_map, ValueOf(o_bytes_map).Set(u_map).Interface(), "bytes", "map")
	gt.Equal(u_bytes_struct, ValueOf(o_bytes_struct).Set(u_struct).Interface(), "bytes", "struct")
	gt.Equal(u_bytes_time, ValueOf(o_bytes_time).Set(u_time).Interface(), "bytes", "time")
	gt.Equal(u_bytes_uuid, ValueOf(o_bytes_uuid).Set(u_uuid).Interface(), "bytes", "uuid")

	gt.Equal(u_bool, ValueOf(o_bool).Set(u_bytes_bool).Interface(), "bool", "bytes")
	gt.Equal(u_bool, ValueOf(o_bool).Set(u_bool).Interface(), "bool", "bool")
	gt.Equal(u_bool, ValueOf(o_bool).Set(u_int).Interface(), "bool", "int")
	gt.Equal(u_bool, ValueOf(o_bool).Set(u_uint).Interface(), "bool", "uint")
	gt.Equal(u_bool, ValueOf(o_bool).Set(u_float).Interface(), "bool", "float")
	gt.Equal(u_bool, ValueOf(o_bool).Set(u_array).Interface(), "bool", "array")
	gt.Equal(u_bool, ValueOf(o_bool).Set(u_slice).Interface(), "bool", "slice")
	gt.Equal(u_bool, ValueOf(o_bool).Set(u_str_bool).Interface(), "bool", "string")
	gt.Equal(u_bool, ValueOf(o_bool).Set(u_map).Interface(), "bool", "map")
	gt.Equal(u_bool, ValueOf(o_bool).Set(u_struct).Interface(), "bool", "struct")
	gt.Equal(u_bool, ValueOf(o_bool).Set(u_time).Interface(), "bool", "time")
	gt.Equal(u_bool, ValueOf(o_bool).Set(u_uuid).Interface(), "bool", "uuid")

	gt.Equal(u_int, ValueOf(o_int).Set(u_bytes_int).Interface(), "int", "bytes")
	gt.Equal(u_int_bool, ValueOf(o_int).Set(u_bool).Interface(), "int", "bool")
	gt.Equal(u_int, ValueOf(o_int).Set(u_int).Interface(), "int", "int")
	gt.Equal(u_int, ValueOf(o_int).Set(u_uint).Interface(), "int", "uint")
	gt.Equal(u_int, ValueOf(o_int).Set(u_float).Interface(), "int", "float")
	gt.Equal(u_int, ValueOf(o_int).Set(u_str_int).Interface(), "int", "string")
	gt.Equal(u_int, ValueOf(o_int).Set(u_time).Interface(), "int", "time")

	gt.Equal(u_uint, ValueOf(o_uint).Set(u_bytes_uint).Interface(), "uint", "bytes")
	gt.Equal(u_uint_bool, ValueOf(o_uint).Set(u_bool).Interface(), "uint", "bool")
	gt.Equal(u_uint, ValueOf(o_uint).Set(u_int).Interface(), "uint", "int")
	gt.Equal(u_uint, ValueOf(o_uint).Set(u_uint).Interface(), "uint", "uint")
	gt.Equal(u_uint, ValueOf(o_uint).Set(u_float).Interface(), "uint", "float")
	gt.Equal(u_uint, ValueOf(o_uint).Set(u_str_uint).Interface(), "uint", "string")
	gt.Equal(u_uint, ValueOf(o_uint).Set(u_time).Interface(), "uint", "time")

	gt.Equal(u_float, ValueOf(o_float).Set(u_bytes_float).Interface(), "float", "bytes")
	gt.Equal(u_float_bool, ValueOf(o_float).Set(u_bool).Interface(), "float", "bool")
	gt.Equal(u_float, ValueOf(o_float).Set(u_int).Interface(), "float", "int")
	gt.Equal(u_float, ValueOf(o_float).Set(u_uint).Interface(), "float", "uint")
	gt.Equal(u_float, ValueOf(o_float).Set(u_float).Interface(), "float", "float")
	gt.Equal(u_float, ValueOf(o_float).Set(u_str_float).Interface(), "float", "string")
	gt.Equal(u_float, ValueOf(o_float).Set(u_time).Interface(), "float", "time")

	gt.Equal(u_str_bool, ValueOf(o_str_bool).Set(u_bytes_string).Interface(), "string", "bytes")
	gt.Equal(u_str_bool, ValueOf(o_str_bool).Set(u_bool).Interface(), "string", "bool")
	gt.Equal(u_str_int, ValueOf(o_str_int).Set(u_int).Interface(), "string", "int")
	gt.Equal(u_str_uint, ValueOf(o_str_uint).Set(u_uint).Interface(), "string", "uint")
	gt.Equal(u_str_float1, ValueOf(o_str_float).Set(u_float1).Interface(), "string", "float")
	gt.Equal(u_str_array, ValueOf(o_str_array).Set(u_array).Interface(), "string", "array")
	gt.Equal(u_str_slice, ValueOf(o_str_slice).Set(u_slice).Interface(), "string", "slice")
	gt.Equal(u_str_bool, ValueOf(o_str_bool).Set(u_str_bool).Interface(), "string", "string")
	gt.Equal(u_str_map, ValueOf(o_str_map).Set(u_map).Interface(), "string", "map")
	gt.Equal(u_str_struct, ValueOf(o_str_struct).Set(u_struct).Interface(), "string", "struct")
	gt.Equal(u_str_time, ValueOf(o_str_time).Set(u_time_str).Interface(), "string", "time")
	gt.Equal(u_str_uuid, ValueOf(o_str_uuid).Set(u_uuid).Interface(), "string", "uuid")

	gt.Equal(u_time, ValueOf(o_time).Set(u_bytes_time).Interface(), "time", "bytes")
	gt.Equal(u_time, ValueOf(o_time).Set(u_int).Interface(), "time", "int")
	gt.Equal(u_time, ValueOf(o_time).Set(u_uint).Interface(), "time", "uint")
	gt.Equal(u_time, ValueOf(o_time).Set(u_float).Interface(), "time", "float")
	gt.Equal(u_time_str, ValueOf(o_time).Set(u_str_time).Interface(), "time", "string")
	gt.Equal(u_time, ValueOf(o_time).Set(u_time).Interface(), "time", "time")

	gt.Equal(u_uuid, ValueOf(o_uuid).Set(u_bytes_uuid).Interface(), "uuid", "bytes")
	gt.Equal(u_uuid, ValueOf(o_uuid).Set(u_str_uuid).Interface(), "uuid", "string")
	gt.Equal(u_uuid, ValueOf(o_uuid).Set(u_uuid).Interface(), "uuid", "uuid")

	gt.Equal(u_array, ValueOf(o_array).Set(u_array).Interface(), "array", "array")
	gt.Equal(u_slice, ValueOf(o_slice).Set(u_slice).Interface(), "slice", "slice")
	gt.Equal(u_map, ValueOf(o_map).Set(u_map).Interface(), "map", "map")
	gt.Equal(u_struct, ValueOf(o_struct).Set(u_struct).Interface(), "struct", "struct")

	// Test pointer values after Set()
	gt.Equal(u_bytes_bool, testSet(o_bytes_bool, u_bytes_bool), "*bytes", "bytes")
	gt.Equal(u_bytes_int, testSet(o_bytes_int, u_int), "*bytes", "int")
	gt.Equal(u_bytes_uint, testSet(o_bytes_uint, u_uint), "*bytes", "uint")
	gt.Equal(u_bytes_float, testSet(o_bytes_float, u_float), "*bytes", "float")
	gt.Equal(u_bytes_array, testSet(o_bytes_array, u_array), "*bytes", "array")
	gt.Equal(u_bytes_slice, testSet(o_bytes_slice, u_slice), "*bytes", "slice")
	gt.Equal(u_bytes_string, testSet(o_bytes_string, u_str_bool), "*bytes", "string")
	gt.Equal(u_bytes_map, testSet(o_bytes_map, u_map), "*bytes", "map")
	gt.Equal(u_bytes_struct, testSet(o_bytes_struct, u_struct), "*bytes", "struct")
	gt.Equal(u_bytes_time, testSet(o_bytes_time, u_time), "*bytes", "time")
	gt.Equal(u_bytes_uuid, testSet(o_bytes_uuid, u_uuid), "*bytes", "uuid")

	gt.Equal(u_bool, testSet(o_bool, u_bytes_bool), "*bool", "bytes")
	gt.Equal(u_bool, testSet(o_bool, u_bool), "*bool", "bool")
	gt.Equal(u_bool, testSet(o_bool, u_int), "*bool", "int")
	gt.Equal(u_bool, testSet(o_bool, u_uint), "*bool", "uint")
	gt.Equal(u_bool, testSet(o_bool, u_float), "*bool", "float")
	gt.Equal(u_bool, testSet(o_bool, u_array), "*bool", "array")
	gt.Equal(u_bool, testSet(o_bool, u_slice), "*bool", "slice")
	gt.Equal(u_bool, testSet(o_bool, u_str_bool), "*bool", "string")
	gt.Equal(u_bool, testSet(o_bool, u_map), "*bool", "map")
	gt.Equal(u_bool, testSet(o_bool, u_struct), "*bool", "struct")
	gt.Equal(u_bool, testSet(o_bool, u_time), "*bool", "time")
	gt.Equal(u_bool, testSet(o_bool, u_uuid), "*bool", "uuid")

	gt.Equal(u_int, testSet(o_int, u_bytes_int), "*int", "bytes")
	gt.Equal(u_int_bool, testSet(o_int, u_bool), "*int", "bool")
	gt.Equal(u_int, testSet(o_int, u_int), "*int", "int")
	gt.Equal(u_int, testSet(o_int, u_uint), "*int", "uint")
	gt.Equal(u_int, testSet(o_int, u_float), "*int", "float")
	gt.Equal(u_int, testSet(o_int, u_str_int), "*int", "string")
	gt.Equal(u_int, testSet(o_int, u_time), "*int", "time")

	gt.Equal(u_uint, testSet(o_uint, u_bytes_uint), "*uint", "bytes")
	gt.Equal(u_uint_bool, testSet(o_uint, u_bool), "*uint", "bool")
	gt.Equal(u_uint, testSet(o_uint, u_int), "*uint", "int")
	gt.Equal(u_uint, testSet(o_uint, u_uint), "*uint", "uint")
	gt.Equal(u_uint, testSet(o_uint, u_float), "*uint", "float")
	gt.Equal(u_uint, testSet(o_uint, u_str_uint), "*uint", "string")
	gt.Equal(u_uint, testSet(o_uint, u_time), "*uint", "time")

	gt.Equal(u_float, testSet(o_float, u_bytes_float), "*float", "bytes")
	gt.Equal(u_float_bool, testSet(o_float, u_bool), "*float", "bool")
	gt.Equal(u_float, testSet(o_float, u_int), "*float", "int")
	gt.Equal(u_float, testSet(o_float, u_uint), "*float", "uint")
	gt.Equal(u_float, testSet(o_float, u_float), "*float", "float")
	gt.Equal(u_float, testSet(o_float, u_str_float), "*float", "string")
	gt.Equal(u_float, testSet(o_float, u_time), "*float", "time")

	gt.Equal(u_str_bool, testSet(o_str_bool, u_bytes_string), "*string", "bytes")
	gt.Equal(u_str_bool, testSet(o_str_bool, u_bool), "*string", "bool")
	gt.Equal(u_str_int, testSet(o_str_int, u_int), "*string", "int")
	gt.Equal(u_str_uint, testSet(o_str_uint, u_uint), "*string", "uint")
	gt.Equal(u_str_float1, testSet(o_str_float, u_float1), "*string", "float")
	gt.Equal(u_str_array, testSet(o_str_array, u_array), "*string", "array")
	gt.Equal(u_str_slice, testSet(o_str_slice, u_slice), "*string", "slice")
	gt.Equal(u_str_bool, testSet(o_str_bool, u_str_bool), "*string", "string")
	gt.Equal(u_str_map, testSet(o_str_map, u_map), "*string", "map")
	gt.Equal(u_str_struct, testSet(o_str_struct, u_struct), "*string", "struct")
	gt.Equal(u_str_time, testSet(o_str_time, u_time_str), "*string", "time")
	gt.Equal(u_str_uuid, testSet(o_str_uuid, u_uuid), "*string", "uuid")

	gt.Equal(u_time, testSet(o_time, u_bytes_time), "*time", "bytes")
	gt.Equal(u_time, testSet(o_time, u_int), "*time", "int")
	gt.Equal(u_time, testSet(o_time, u_uint), "*time", "uint")
	gt.Equal(u_time, testSet(o_time, u_float), "*time", "float")
	gt.Equal(u_time_str, testSet(o_time, u_str_time), "*time", "string")
	gt.Equal(u_time, testSet(o_time, u_time), "*time", "time")

	gt.Equal(u_uuid, testSet(o_uuid, u_bytes_uuid), "*uuid", "bytes")
	gt.Equal(u_uuid, testSet(o_uuid, u_str_uuid), "*uuid", "string")
	gt.Equal(u_uuid, testSet(o_uuid, u_uuid), "*uuid", "uuid")

	// Test pointer values after Index().Set()
	gt.Equal(u_array, testSetIndex(o_array, u_bytes_num, 1), "*[2]string[1]", "bytes")
	gt.Equal(u_array, testSetIndex(o_array, u_int, 1), "*[2]string[1]", "int")
	gt.Equal(u_array, testSetIndex(o_array, u_uint, 1), "*[2]string[1]", "uint")
	gt.Equal(u_array, testSetIndex(o_array, u_float, 1), "*[2]string[1]", "float")
	gt.Equal(u_array, testSetIndex(o_array, u_str_int, 1), "*[2]string[1]", "string")

	gt.Equal(u_slice, testSetIndex(o_slice, u_bytes_num, 1), "*[]string[1]", "bytes")
	gt.Equal(u_slice, testSetIndex(o_slice, u_int, 1), "*[]string[1]", "int")
	gt.Equal(u_slice, testSetIndex(o_slice, u_uint, 1), "*[]string[1]", "uint")
	gt.Equal(u_slice, testSetIndex(o_slice, u_float, 1), "*[]string[1]", "float")
	gt.Equal(u_slice, testSetIndex(o_slice, u_str_int, 1), "*[]string[1]", "string")

	gt.Equal(u_map, testSetIndex(o_map, u_bytes_num, 1), "*map[string]string[1]", "bytes")
	gt.Equal(u_map, testSetIndex(o_map, u_int, 1), "*map[string]string[1]", "int")
	gt.Equal(u_map, testSetIndex(o_map, u_uint, 1), "*map[string]string[1]", "uint")
	gt.Equal(u_map, testSetIndex(o_map, u_float, 1), "*map[string]string[1]", "float")
	gt.Equal(u_map, testSetIndex(o_map, u_str_int, 1), "*map[string]string[1]", "string")

	gt.Equal(u_struct, testSetIndex(o_struct, u_bytes_num, 1), "*string_struct[1]", "bytes")
	gt.Equal(u_struct, testSetIndex(o_struct, u_int, 1), "*string_struct[1]", "int")
	gt.Equal(u_struct, testSetIndex(o_struct, u_uint, 1), "*string_struct[1]", "uint")
	gt.Equal(u_struct, testSetIndex(o_struct, u_float, 1), "*string_struct[1]", "float")
	gt.Equal(u_struct, testSetIndex(o_struct, u_str_int, 1), "*string_struct[1]", "string")

	gt.Equal(u_array1, testSetIndex(o_array1, u_bytes_num, 0), "*[1]string[0]", "bytes")
	gt.Equal(u_array1, testSetIndex(o_array1, u_int, 0), "*[1]string[0]", "int")
	gt.Equal(u_array1, testSetIndex(o_array1, u_uint, 0), "*[1]string[0]", "uint")
	gt.Equal(u_array1, testSetIndex(o_array1, u_float, 0), "*[1]string[0]", "float")
	gt.Equal(u_array1, testSetIndex(o_array1, u_str_int, 0), "*[1]string[0]", "string")

	gt.Equal(u_struct1, testSetIndex(o_struct1, u_bytes_num, 0), "*string_struct_single[0]", "bytes")
	gt.Equal(u_struct1, testSetIndex(o_struct1, u_int, 0), "*string_struct_single[0]", "int")
	gt.Equal(u_struct1, testSetIndex(o_struct1, u_uint, 0), "*string_struct_single[0]", "uint")
	gt.Equal(u_struct1, testSetIndex(o_struct1, u_float, 0), "*string_struct_single[0]", "float")
	gt.Equal(u_struct1, testSetIndex(o_struct1, u_str_int, 0), "*string_struct_single[0]", "string")

}

func TestValueSetIndex(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Set()"
	for n := range getTestVars() {
		a := createTestVars(0 == 1, (0 * 1), "false", n)[n]
		nv := ValueOf(&a).Elem().SetType()
		testSetDeep(nv, false)
		gt.Equal("true", testGetDeep(nv), n)
	}
}

func testSet(original, new any) any {
	n := original
	ValueOf(&n).Elem().Elem().Set(new)
	return n
}

func testSetIndex(original, new any, index int) any {
	n := original
	ValueOf(&n).Elem().Elem().Index(index).Set(new)
	return n
}

func testSetDeep(v VALUE, embedded bool) {
	switch v.Kind() {
	case Pointer, Interface:
		testSetDeep(v.Elem(), true)
	case Array, Slice, Map, Struct:
		i := v.Len() - 1
		if !embedded {
			testSetDeep(v.Index(i), true)
		} else if testGetDeep(v.Index(i)) == "true" {
			panic("value already set")
		} else {
			switch v.Kind() {
			case Array:
				if _, is := v.Interface().([2]string); is {
					v.Set([2]string{"false", "true"})
				} else {
					testSetDeep(v.Index(i), true)
				}
			case Slice:
				if _, is := v.Interface().([]string); is && v.Len() == 2 {
					v.Set([]string{"false", "true"})
				} else {
					testSetDeep(v.Index(i), true)
				}
			case Map:
				if _, is := v.Interface().(map[string]string); is && v.Len() == 2 {
					v.Set(map[string]string{"0": "false", "1": "true"})
				} else {
					testSetDeep(v.Index(i), true)
				}
			case Struct:
				if _, is := v.Interface().(string_struct); is {
					v.Set(string_struct{"false", "true"})
				} else {
					testSetDeep(v.Index(i), true)
				}
			}
		}
	case Bool:
		v.Set([]byte{0x1})
	case Int:
		v.Set([]byte{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	case String:
		v.Set([]byte("true"))
	default:
		panic("value type not supported")
	}
}

func testGetDeep(v VALUE) string {
	v = testGetDeepValue(v)
	if v.Kind() == Int {
		v = v.BOOL().VALUE()
	}
	return v.String()
}

func testGetDeepValue(v VALUE) VALUE {
	switch v.Kind() {
	case Pointer, Interface:
		return testGetDeepValue(v.Elem())
	case Array, Slice, Map, Struct:
		return testGetDeepValue(v.Index(v.Len() - 1))
	case Bool, String, Int:
		return v
	}
	panic("value type not supported")
}

func TestValueConversion(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).%s()"
	var (
		_bool   = true
		_int    = 123
		_uint   = uint(123)
		_float  = 123.0
		_float1 = 123.1
		_array  = [2]string{"0", "123"}
		_slice  = []string{"0", "123"}
		_map    = map[string]string{"0": "0", "1": "123"}
		_struct = string_struct{"0", "123"}
		_time   = TIME(time.Unix(0, 123).UTC())
		_uuid   = UUID([16]byte{0x26, 0x7b, 0x32, 0x29, 0x25, 0x66, 0x44, 0x26, 0xa8, 0x26, 0x8d, 0x80, 0x12, 0x6e, 0x71, 0x9a})

		str_bool   = "true"
		str_int    = "123"
		str_uint   = "123"
		str_float  = "123.0"
		str_float1 = "123.1"
		str_array  = `["0","123"]`
		str_slice  = `["0","123"]`
		str_map    = `{"0":"0","1":"123"}`
		str_struct = `{"V1":"0","V2":"123"}`
		str_time   = `1970-01-01 00:00:00.000000123`
		str_time_c = `Jan 1, 70T00:00:00.000000123`
		str_uuid   = `267b3229-2566-4426-a826-8d80126e719a`

		bytes_bool   = []byte{0x1}
		bytes_int    = []byte{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		bytes_uint   = []byte{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		bytes_float  = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0xc0, 0x5e, 0x40}
		bytes_array  = []byte(str_array)
		bytes_slice  = []byte(str_slice)
		bytes_map    = []byte(str_map)
		bytes_struct = []byte(str_struct)
		bytes_time   = []byte{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		bytes_uuid   = []byte{0x26, 0x7b, 0x32, 0x29, 0x25, 0x66, 0x44, 0x26, 0xa8, 0x26, 0x8d, 0x80, 0x12, 0x6e, 0x71, 0x9a}
		bytes_string = []byte("true")

		int_bool   = 1
		uint_bool  = uint(1)
		float_bool = 1.0
		array_any  = [2]any{"0", "123"}
		slice_any  = []any{"0", "123"}
		map_any    = map[string]any{"0": "0", "1": "123"}
		map_struct = map[string]any{"V1": "0", "V2": "123"}

		test_struct = &string_struct{}
	)

	gt.Equal(bytes_bool, ValueOf(bytes_bool).Bytes(), "bytes", "Bytes")
	gt.Equal(bytes_bool, ValueOf(_bool).Bytes(), "bool", "Bytes")
	gt.Equal(bytes_int, ValueOf(_int).Bytes(), "int", "Bytes")
	gt.Equal(bytes_uint, ValueOf(_uint).Bytes(), "uint", "Bytes")
	gt.Equal(bytes_float, ValueOf(_float).Bytes(), "float", "Bytes")
	gt.Equal(bytes_array, ValueOf(_array).Bytes(), "array", "Bytes")
	gt.Equal(bytes_slice, ValueOf(_slice).Bytes(), "slice", "Bytes")
	gt.Equal(bytes_string, ValueOf(str_bool).Bytes(), "string", "Bool")
	gt.Equal(bytes_map, ValueOf(_map).Bytes(), "map", "Bytes")
	gt.Equal(bytes_struct, ValueOf(_struct).Bytes(), "struct", "Bytes")
	gt.Equal(bytes_time, ValueOf(_time).Bytes(), "time", "Bytes")
	gt.Equal(bytes_uuid, ValueOf(_uuid).Bytes(), "uuid", "Bytes")

	gt.Equal(_bool, ValueOf(bytes_bool).Bool(), "bytes", "Bool")
	gt.Equal(_bool, ValueOf(_bool).Bool(), "bool", "Bool")
	gt.Equal(_bool, ValueOf(_int).Bool(), "int", "Bool")
	gt.Equal(_bool, ValueOf(_uint).Bool(), "uint", "Bool")
	gt.Equal(_bool, ValueOf(_float).Bool(), "float", "Bool")
	gt.Equal(_bool, ValueOf(_array).Bool(), "array", "Bool")
	gt.Equal(_bool, ValueOf(_slice).Bool(), "slice", "Bool")
	gt.Equal(_bool, ValueOf(_map).Bool(), "map", "Bool")
	gt.Equal(_bool, ValueOf(str_bool).Bool(), "string", "Bool")
	gt.Equal(_bool, ValueOf(_struct).Bool(), "struct", "Bool")
	gt.Equal(_bool, ValueOf(_time).Bool(), "time", "Bool")
	gt.Equal(_bool, ValueOf(_uuid).Bool(), "uuid", "Bool")

	gt.Equal(_int, ValueOf(bytes_int).Int(), "bytes", "Int")
	gt.Equal(int_bool, ValueOf(_bool).Int(), "bool", "Int")
	gt.Equal(_int, ValueOf(_int).Int(), "int", "Int")
	gt.Equal(_int, ValueOf(_uint).Int(), "uint", "Int")
	gt.Equal(_int, ValueOf(_float).Int(), "float", "Int")
	gt.Equal(_int, ValueOf(str_int).Int(), "string", "Int")
	gt.Equal(_int, ValueOf(_time).Int(), "time", "Int")

	gt.Equal(_uint, ValueOf(bytes_uint).Uint(), "bytes", "Uint")
	gt.Equal(uint_bool, ValueOf(_bool).Uint(), "bool", "Uint")
	gt.Equal(_uint, ValueOf(_int).Uint(), "int", "Uint")
	gt.Equal(_uint, ValueOf(_uint).Uint(), "uint", "Uint")
	gt.Equal(_uint, ValueOf(_float).Uint(), "float", "Uint")
	gt.Equal(_uint, ValueOf(str_uint).Uint(), "string", "Uint")
	gt.Equal(_uint, ValueOf(_time).Uint(), "time", "Uint")

	gt.Equal(_float, ValueOf(bytes_float).Float64(), "bytes", "Float")
	gt.Equal(float_bool, ValueOf(_bool).Float64(), "bool", "Float")
	gt.Equal(_float, ValueOf(_int).Float64(), "int", "Float")
	gt.Equal(_float, ValueOf(_uint).Float64(), "uint", "Float")
	gt.Equal(_float, ValueOf(_float).Float64(), "float", "Float")
	gt.Equal(_float, ValueOf(str_float).Float64(), "string", "Float")
	gt.Equal(_float, ValueOf(_time).Float64(), "time", "Float")

	gt.Equal(str_bool, ValueOf(bytes_string).String(), "bytes", "String")
	gt.Equal(str_bool, ValueOf(_bool).String(), "bool", "String")
	gt.Equal(str_int, ValueOf(_int).String(), "int", "String")
	gt.Equal(str_uint, ValueOf(_uint).String(), "uint", "String")
	gt.Equal(str_float1, ValueOf(_float1).String(), "float", "String")
	gt.Equal(str_array, ValueOf(_array).String(), "array", "String")
	gt.Equal(str_slice, ValueOf(_slice).String(), "slice", "String")
	gt.Equal(str_map, ValueOf(_map).String(), "map", "String")
	gt.Equal(str_struct, ValueOf(_struct).String(), "struct", "String")
	gt.Equal(str_time, ValueOf(_time).TIME().Format(ISO8601N), "time", "String")
	gt.Equal(str_uuid, ValueOf(_uuid).String(), "uuid", "String")

	gt.Equal(array_any, ValueOf(bytes_array).JSON().ARRAY().Interface(), "bytes", "Array")
	gt.Equal(_array, ValueOf(_array).ARRAY().Interface(), "array", "Array")
	gt.Equal(_array, ValueOf(_slice).ARRAY().Interface(), "slice", "Array")
	gt.Equal(_array, ValueOf(_map).ARRAY().Interface(), "map", "Array")
	gt.Equal(array_any, ValueOf(str_array).JSON().ARRAY().Interface(), "string", "Array")
	gt.Equal(array_any, ValueOf(_struct).ARRAY().Interface(), "struct", "Array")

	gt.Equal(slice_any, ValueOf(bytes_slice).JSON().Slice(), "bytes", "Slice")
	gt.Equal(_slice, ValueOf(_array).SLICE().Interface(), "array", "Slice")
	gt.Equal(_slice, ValueOf(_slice).SLICE().Interface(), "slice", "Slice")
	gt.Equal(_slice, ValueOf(_map).SLICE().Interface(), "map", "Slice")
	gt.Equal(slice_any, ValueOf(str_slice).JSON().SLICE().Interface(), "string", "Slice")
	gt.Equal(slice_any, ValueOf(_struct).SLICE().Interface(), "struct", "Slice")

	gt.Equal(map_any, ValueOf(bytes_map).JSON().MAP().Interface(), "bytes", "Map")
	gt.Equal(map_any, ValueOf(_array).MAP().Interface(), "array", "Map")
	gt.Equal(map_any, ValueOf(_slice).MAP().Interface(), "slice", "Map")
	gt.Equal(_map, ValueOf(_map).MAP().Interface(), "map", "Map")
	gt.Equal(map_any, ValueOf(str_map).JSON().MAP().Interface(), "string", "Map")
	gt.Equal(map_struct, ValueOf(_struct).MAP().Interface(), "struct", "Map")

	gt.Equal(_time, ValueOf(bytes_time).TIME(), "bytes", "TIME")
	gt.Equal(_time, ValueOf(_int).TIME(), "int", "TIME")
	gt.Equal(_time, ValueOf(_uint).TIME(), "uint", "TIME")
	gt.Equal(_time, ValueOf(_float).TIME(), "float", "TIME")
	gt.Equal(_time, ValueOf(_time).TIME(), "time", "TIME")
	gt.Equal(_time, ValueOf(str_time).TIME(), "string", "TIME")
	gt.Equal(_time, STRING(str_time_c).ParseTime(), "string", "TIME")

	gt.Equal(_uuid, ValueOf(bytes_uuid).UUID(), "bytes", "UUID")
	gt.Equal(_uuid, ValueOf(str_uuid).UUID(), "string", "UUID")
	gt.Equal(_uuid, ValueOf(_uuid).UUID(), "uuid", "UUID")

	gt.Msg = "Testing ValueOf(%s).%s"
	test_struct = &string_struct{}
	JSON(str_struct).Scan(test_struct)
	gt.Equal(_struct, *test_struct, "bytes", "JSON().Scan(struct)")
	test_struct = &string_struct{}
	ValueOf(_array).ARRAY().Scan(test_struct)
	gt.Equal(_struct, *test_struct, "array", "ARRAY().Scan(struct)")
	test_struct = &string_struct{}
	ValueOf(_slice).SLICE().Scan(test_struct)
	gt.Equal(_struct, *test_struct, "slice", "SLICE().Scan(struct)")
	test_struct = &string_struct{}
	ValueOf(map_struct).MAP().Scan(test_struct)
	gt.Equal(_struct, *test_struct, "map", "MAP().Scan(struct)")
	test_struct = &string_struct{}
	ValueOf(str_struct).JSON().Scan(test_struct)
	gt.Equal(_struct, *test_struct, "string", "JSON().Scan(struct)")
	test_struct = &string_struct{}
	ValueOf(_struct).STRUCT().Scan(test_struct)
	gt.Equal(_struct, *test_struct, "struct", "STRUCT().Scan(struct)")

}

func TestValueNewDeep(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).NewDeep().Set()"
	// l is a list of test var keys (or the first chars of the key)
	l := []string{}
	v := getTestVars()
	for n := range v {
		// get newly created test vars to avoid preset pointers
		v := createTestVars(0 == 1, 0*1, "false", n)[n]
		// filter for test vars in l
		var proc bool
		if len(l) == 0 {
			proc = true
		} else {
			for _, a := range l {
				ln := int(math.Min(float64(len(a)), float64(len(n))))
				if n[:ln] == a {
					proc = true
					break
				}
			}
		}
		// run test if in filter
		if proc {
			r := ValueOf(v)
			nv := r.NewDeep()
			testSetDeep(nv, false)
			gt.Equal("true", testGetDeep(nv), n)
		}
	}
}

func TestEncodeBasic(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing Encode(%s).Bytes()"
	var (
		b  = true
		i  = 1
		f  = 1.0
		s  = "1"
		a  = [2]string{"0", "1"}
		l  = []string{"0", "1"}
		m  = map[string]string{"0": "0", "1": "1"}
		d  = string_struct{"0", "1"}
		a1 = [1]string{"0"}
		l1 = []string{"0"}
		d1 = string_struct_single{"0"}

		bEnc  = Encode(b)
		iEnc  = Encode(i)
		fEnc  = Encode(f)
		sEnc  = Encode(s)
		aEnc  = Encode(a)
		lEnc  = Encode(l)
		mEnc  = Encode(m)
		dEnc  = Encode(d)
		a1Enc = Encode(a1)
		l1Enc = Encode(l1)
		d1Enc = Encode(d1)

		bDec  = bEnc.Decodex()
		iDec  = iEnc.Decodex()
		fDec  = fEnc.Decodex()
		sDec  = sEnc.Decodex()
		aDec  = aEnc.Decodex()
		lDec  = lEnc.Decodex()
		mDec  = mEnc.Decodex()
		dDec  = dEnc.Decodex()
		a1Dec = a1Enc.Decodex()
		l1Dec = l1Enc.Decodex()
		d1Dec = d1Enc.Decodex()

		bBytes  = []byte{1, 1}
		iBytes  = []byte{2, 1, 0, 0, 0, 0, 0, 0, 0}
		fBytes  = []byte{14, 0, 0, 0, 0, 0, 0, 240, 63}
		sBytes  = []byte{24, 8, 1, 49}
		aBytes  = []byte{17, 24, 8, 2, 24, 8, 1, 48, 24, 8, 1, 49}
		lBytes  = []byte{23, 24, 8, 2, 24, 8, 1, 48, 24, 8, 1, 49}
		mBytes  = []byte{21, 24, 24, 8, 2, 24, 8, 1, 48, 24, 8, 1, 48, 24, 8, 1, 49, 24, 8, 1, 49}
		dBytes  = []byte{25, 8, 2, 24, 8, 1, 48, 24, 8, 1, 49}
		a1Bytes = []byte{17, 24, 8, 1, 24, 8, 1, 48}
		l1Bytes = []byte{23, 24, 8, 1, 24, 8, 1, 48}
		d1Bytes = []byte{25, 8, 1, 24, 8, 1, 48}

		bDecVal  = decodex{1, 2, ValueOf(b)}
		iDecVal  = decodex{2, 9, ValueOf(i)}
		fDecVal  = decodex{14, 9, ValueOf(f)}
		sDecVal  = decodex{24, 4, ValueOf(s)}
		aDecVal  = decodex{17, 12, ValueOf(a)}
		lDecVal  = decodex{23, 12, ValueOf(l)}
		mDecVal  = decodex{21, 21, ValueOf(m)}
		dDecVal  = decodex{25, 11, ValueOf(l)}
		a1DecVal = decodex{17, 8, ValueOf(a1)}
		l1DecVal = decodex{23, 8, ValueOf(l1)}
		d1DecVal = decodex{25, 7, ValueOf(l1)}
	)

	gt.Equal(bBytes, bEnc.Bytes(), "bool")
	gt.Equal(iBytes, iEnc.Bytes(), "int")
	gt.Equal(fBytes, fEnc.Bytes(), "float")
	gt.Equal(sBytes, sEnc.Bytes(), "string")
	gt.Equal(aBytes, aEnc.Bytes(), "array")
	gt.Equal(lBytes, lEnc.Bytes(), "slice")
	gt.Equal(mBytes, mEnc.Bytes(), "map")
	gt.Equal(dBytes, dEnc.Bytes(), "struct")
	gt.Equal(a1Bytes, a1Enc.Bytes(), "array1")
	gt.Equal(l1Bytes, l1Enc.Bytes(), "slice1")
	gt.Equal(d1Bytes, d1Enc.Bytes(), "struct1")

	gt.Msg = "Testing Encode(%s).Decodex()"
	gt.Equal(bDecVal.Serialize(), bDec.Serialize(), "bool")
	gt.Equal(iDecVal.Serialize(), iDec.Serialize(), "int")
	gt.Equal(fDecVal.Serialize(), fDec.Serialize(), "float")
	gt.Equal(sDecVal.Serialize(), sDec.Serialize(), "string")
	gt.Equal(aDecVal.Serialize(), aDec.Serialize(), "array")
	gt.Equal(lDecVal.Serialize(), lDec.Serialize(), "slice")
	gt.Equal(mDecVal.Serialize(), mDec.Serialize(), "map")
	gt.Equal(dDecVal.Serialize(), dDec.Serialize(), "struct")
	gt.Equal(a1DecVal.Serialize(), a1Dec.Serialize(), "array1")
	gt.Equal(l1DecVal.Serialize(), l1Dec.Serialize(), "slice1")
	gt.Equal(d1DecVal.Serialize(), d1Dec.Serialize(), "struct1")

	gt.Msg = "Testing Decode(Encode(%[1]s), *%[1]s)"
	var (
		bOut  = ValueOf(b).NewDeep()
		iOut  = ValueOf(i).NewDeep()
		fOut  = ValueOf(f).NewDeep()
		sOut  = ValueOf(s).NewDeep()
		aOut  = ValueOf(a).NewDeep()
		lOut  = ValueOf(l).NewDeep()
		mOut  = ValueOf(m).NewDeep()
		dOut  = ValueOf(d).NewDeep()
		a1Out = ValueOf(a1).NewDeep()
		d1Out = ValueOf(d1).NewDeep()
	)

	bEnc.Decode(bOut)
	iEnc.Decode(iOut)
	fEnc.Decode(fOut)
	sEnc.Decode(sOut)
	aEnc.Decode(aOut)
	lEnc.Decode(lOut)
	mEnc.Decode(mOut)
	dEnc.Decode(dOut)
	a1Enc.Decode(a1Out)
	l1Enc.Decode(d1Out)

	gt.Equal(b, bOut.Elem().Interface(), bOut.Elem().typ)
	gt.Equal(i, iOut.Elem().Interface(), iOut.Elem().typ)
	gt.Equal(f, fOut.Elem().Interface(), fOut.Elem().typ)
	gt.Equal(s, sOut.Elem().Interface(), sOut.Elem().typ)
	gt.Equal(a, aOut.Elem().Interface(), aOut.Elem().typ)
	gt.Equal(l, lOut.Elem().Interface(), lOut.Elem().typ)
	gt.Equal(m, mOut.Elem().Interface(), mOut.Elem().typ)
	gt.Equal(d, dOut.Elem().Interface(), dOut.Elem().typ)
	gt.Equal(a1, a1Out.Elem().Interface(), a1Out.Elem().typ)
	gt.Equal(d1, d1Out.Elem().Interface(), d1Out.Elem().typ)

}

func TestEncodeComplex(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing Decode(Encode(%[1]s), *%[1]s)"
	// l is a list of test var keys (or the first chars of the key)
	l := []string{}
	vars := getTestVars()
	for n, v := range vars {
		// filter for test vars in l
		var proc bool
		if len(l) == 0 {
			proc = true
		} else {
			for _, a := range l {
				ln := int(math.Min(float64(len(a)), float64(len(n))))
				if n[:ln] == a {
					proc = true
					break
				}
			}
		}
		// run test if in filter
		if proc {
			sv := ValueOf(v)
			testSetDeep(sv, false)
			e := Encode(sv)
			d := sv.NewDeep()
			e.Decode(d)
			gt.Equal("true", testGetDeep(d), n)
		}
	}
}

func TestValueScan(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf().Scan()"
	s := [][]string{
		{"one", "two"},
		{"three", "four"},
	}
	d := &[]*string_struct{}
	SliceOf(s).ScanList(d)
	r := `[{"V1":"one","V2":"two"},{"V1":"three","V2":"four"}]`
	gt.Equal(r, ValueOf(d).Serialize())
	m := []*map[string]int{
		{"one": 1, "two": 2},
		{"one": 3, "two": 4},
	}
	type stest struct {
		One string `json:"one"`
		Two string `json:"two"`
	}
	d1 := &[]*stest{}
	SliceOf(m).ScanList(d1, "json")
	r = `[{"One":"1","Two":"2"},{"One":"3","Two":"4"}]`
	gt.Equal(r, ValueOf(d1).Serialize())

}

func TestFuncType(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "%s"

	type testStruct struct {
		one string
		two string
	}
	f := func(s string, d *testStruct) *testStruct {
		d.one = s
		d.two = s
		return d
	}

	gt.Equal("func(string, *gotype.testStruct) *gotype.testStruct", TypeOf(f).String(), "TypeOf(func)")
	gt.Equal(2, TypeOf(f).NumIn(), "TypeOf(func).NumIn()")
	gt.Equal(1, TypeOf(f).NumOut(), "TypeOf(func).NumOut()")
	gt.Equal("string", TypeOf(f).In(0).String(), "TypeOf(func).In(0)")
	gt.Equal("*gotype.testStruct", TypeOf(f).In(1).String(), "TypeOf(func).In(1)")
	gt.Equal("*gotype.testStruct", TypeOf(f).Out(0).String(), "TypeOf(func).Out(0)")
}
