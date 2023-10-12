// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	"unsafe"
)

// ------------------------------------------------------------ /
// GOTYPE CUSTOM TYPE IMPLEMENTATION
// implementation of custom type of bool
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

type BOOL bool

// BOOL returns gotype VALUE as gotype BOOL
func (v VALUE) BOOL() BOOL {
	return BOOL(v.Bool())
}

// Bool returns gotype VALUE as bool
func (v VALUE) Bool() bool {
	switch v.Kind() {
	case Bool:
		return *(*bool)(v.ptr)
	case Pointer:
		return v.ElemDeep().Bool()
	default:
		switch v.KIND() {
		case Int:
			return INT(*(*int)(v.ptr)).Bool()
		case Int8:
			return INT(*(*int8)(v.ptr)).Bool()
		case Int16:
			return INT(*(*int16)(v.ptr)).Bool()
		case Int32:
			return INT(*(*int32)(v.ptr)).Bool()
		case Int64:
			return INT(*(*int64)(v.ptr)).Bool()
		case Uint:
			return UINT(*(*uint)(v.ptr)).Bool()
		case Uint8:
			return UINT(*(*uint8)(v.ptr)).Bool()
		case Uint16:
			return UINT(*(*uint16)(v.ptr)).Bool()
		case Uint32:
			return UINT(*(*uint32)(v.ptr)).Bool()
		case Uint64:
			return UINT(*(*uint64)(v.ptr)).Bool()
		case Float32:
			return FLOAT(*(*float32)(v.ptr)).Bool()
		case Float64:
			return FLOAT(*(*float64)(v.ptr)).Bool()
		case Array:
			return (ARRAY)(v).Bool()
		case Map:
			return (MAP)(v).Bool()
		case Slice:
			return (SLICE)(v).Bool()
		case String:
			return (*STRING)(v.ptr).Bool()
		case Time:
			return (*TIME)(v.ptr).Bool()
		case Uuid:
			return (*UUID)(v.ptr).Bool()
		case Struct:
			return (STRUCT)(v).Bool()
		case Bytes:
			return (*BYTES)(v.ptr).Bool()
		}
	}
	panic("cannot convert value to bool")
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Natural returns gotype BOOL as a golang bool
func (b BOOL) Native() bool {
	return bool(b)
}

// Interface returns gotype BOOL as a golang interface{}
func (b BOOL) Interface() any {
	return b.Native()
}

// Value returns gotype BOOL as gotype Value
func (b BOOL) VALUE() VALUE {
	a := (any)(b)
	return *(*VALUE)(unsafe.Pointer(&a))
}

// Encode returns a gotype encoding of BOOL
func (b BOOL) Encode() ENCODING {
	return []byte{byte(Bool), b.Bytes()[0]}
}

// Bytes returns gotype BOOL as []byte
func (b BOOL) Bytes() []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

// BYTES returns gotype BOOL as gotype BYTES
func (b BOOL) BYTES() BYTES {
	if b {
		return BYTES{1}
	}
	return BYTES{0}
}

// String returns gotype BOOL as string
func (b BOOL) String() string {
	if b {
		return "true"
	}
	return "false"
}

// STRING returns gotype BOOL as a gotype STRING
func (b BOOL) STRING() STRING {
	return STRING(b.String())
}

// Bool returns gotype BOOL as bool
func (b BOOL) Bool() bool {
	return bool(b)
}

// Int returns gotype BOOL as int
func (b BOOL) Int() int {
	if b {
		return 1
	}
	return 0
}

// Float returns gotype BOOL as float64
func (b BOOL) Float64() float64 {
	if b {
		return 1
	}
	return 0
}

/*
// Int returns gotype BOOL as a gotype INT
func (b BOOL) INT() INT {
	if b {
		return 1
	}
	return 0
}

// Float returns gotype BOOL as a gotype FLOAT
func (b BOOL) FLOAT() FLOAT {
	if b {
		return 1
	}
	return 0
}
*/
