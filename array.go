// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package vals

import (
	"reflect"
	"unsafe"
)

// ------------------------------------------------------------ /
// GOTYPE CUSTOM TYPE IMPLEMENTATION
// implementation of custom type of array
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

type ARRAY VALUE

// ArrayOf returns a as gotype ARRAY
// panics if a is not convertable to array
func ArrayOf(a any) ARRAY {
	if a, is := a.(ARRAY); is {
		return a
	}
	return ValueOf(a).ARRAY()
}

// ARRAY returns VALUE as gotype ARRAY
func (v VALUE) ARRAY() ARRAY {
	switch v.Kind() {
	case Array:
		return (ARRAY)(v)
	case Pointer:
		return v.ElemDeep().ARRAY()
	default:
		switch v.KIND() {
		case Map:
			return (MAP)(v).ARRAY()
		case Slice:
			return (SLICE)(v).ARRAY()
		case Struct:
			return (STRUCT)(v).ARRAY()
		case Bytes:
			j := (any)(JSON{})
			if v.typ == (*VALUE)(unsafe.Pointer(&j)).typ {
				return (*JSON)(v.ptr).ARRAY()
			}
			return (SLICE)(v).ARRAY()
		default:
			panic("cannot convert value to array")
		}
	}
}

// ------------------------------------------------------------ /
// GOLANG STANDARD IMPLEMENTATIONS
// implementations of functions natively available for
// interface and reflect.Value in golang
// referenced packages: reflect
// ------------------------------------------------------------ /

// Len returns the number of items in ARRAY
func (a ARRAY) Len() int {
	return int((*arrayType)(unsafe.Pointer(a.typ)).len)
}

// Index returns the value found at index i of ARRAY
func (a ARRAY) Index(i int) VALUE {
	if i >= a.Len() {
		panic("index is out of array range")
	}
	return a.index(i)
}

func (a ARRAY) index(i int) VALUE {
	t := (*arrayType)(unsafe.Pointer(a.typ))
	if k := t.elem.Kind(); t.len == 1 && (k == Map /* || k == Pointer */) {
		//return VALUE{t.elem, unsafe.Pointer(&a.ptr)}
	}
	return VALUE{t.elem, unsafe.Pointer(uintptr(a.ptr) + uintptr(i)*t.elem.size), t.elem.flag()}.SetType()
}

// ForEach executes function f on each item in ARRAY,
// note: k equals "" at each item
func (a ARRAY) ForEach(f func(i int, k string, v VALUE) (brake bool)) {
	for i := 0; i < a.Len(); i++ {
		if f(i, "", a.index(i)) {
			break
		}
	}
}

// Set updates the value at index i to value v,
// returns the ARRAY with the updated value
func (a ARRAY) Set(i int, v any) ARRAY {
	if i >= a.Len() {
		a = a.SLICE().Extend(i - a.Len() + 1).ARRAY()
	}
	a.Index(i).Set(v)
	return a
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Native returns gotype ARRAY as a golang any of array
func (a ARRAY) Native() any {
	return a.Interface()
}

// Interface returns gotype ARRAY as a golang interface{}
func (a ARRAY) Interface() any {
	var i any
	iface := (*VALUE)(unsafe.Pointer(&i))
	iface.typ, iface.ptr = a.typ, a.ptr
	if a.Len() == 1 && (*arrayType)(unsafe.Pointer(a.typ)).elem.Kind() == Map {
		//iface.ptr = unsafe.Pointer(&a.ptr) //*(*unsafe.Pointer)(a.ptr)
	}
	return i
}

// VALUE returns gotype ARRAY as gotype VALUE
func (a ARRAY) VALUE() VALUE {
	return (VALUE)(a)
}

// Encode returns a gotype encoding of ARRAY
func (a ARRAY) Encode() ENCODING {
	l := a.Len()
	e := append([]byte{a.typ.Kind().Byte(), (*arrayType)(unsafe.Pointer(a.typ)).elem.Kind().Byte()}, lenBytes(l)...)
	for i := 0; i < l; i++ {
		e = append(e, a.Index(i).Encode()...)
	}
	return e
}

// Bytes returns gotype ARRAY as serialized []byte
func (a ARRAY) Bytes() []byte {
	return []byte(a.String())
}

// Bool returns gotype ARRAY as bool
// false if empty, true if a len > 0
func (a ARRAY) Bool() bool {
	return a.Len() > 0
}

// Slice returns gotype ARRAY as []any
func (a ARRAY) Slice() []any {
	l := a.Len()
	r := make([]any, l)
	for i := 0; i < l; i++ {
		r[i] = a.Index(i).Interface()
	}
	return r
}

// SLICE returns gotype ARRAY as gotype SLICE
func (a ARRAY) SLICE() SLICE {
	at := (*arrayType)(unsafe.Pointer(a.typ))
	a.typ = reflectType(reflect.SliceOf(toType(at.elem)))
	l := int(at.len)
	s := reflect.SliceHeader{
		Data: uintptr(a.ptr),
		Len:  l,
		Cap:  l,
	}
	a.ptr = unsafe.Pointer(&s)
	return (SLICE)(a)
}

// Map returns gotype ARRAY as gotype Map
func (a ARRAY) Map() map[string]any {
	m := map[string]any{}
	for i := 0; i < a.Len(); i++ {
		m[INT(i).String()] = a.Index(i).Interface()
	}
	return m
}

// MAP returns gotype ARRAY as gotype Map
func (a ARRAY) MAP() MAP {
	return (MAP)(ValueOf(a.Map()))
}

// String returns gotype ARRAY as a serialized string
func (a ARRAY) String() string {
	return a.Serialize()
}

// Serialize returns gotype ARRAY as a serialized string
func (a ARRAY) Serialize(ancestry ...ancestor) (s string) {
	if a.ptr == nil {
		return "null"
	}
	if a.Len() == 0 {
		return "[]"
	}
	a.ForEach(func(i int, k string, v VALUE) (brake bool) {
		s += "," + v.serialSafe(ancestry...)
		return
	})
	return "[" + s[1:] + "]"
}

// StructScan reads the values of ARRAY into the provided Struct
// in order of the Struct fields
func (a ARRAY) StructScan(s STRUCT) STRUCT {
	for i := 0; i < a.Len() && i < s.Len(); i++ {
		s.index(i).Set(a.index(i))
	}
	return s
}

// JSON returns gotype ARRAY as gotype JSON
func (a ARRAY) JSON() JSON {
	return JSON(a.Serialize())
}
