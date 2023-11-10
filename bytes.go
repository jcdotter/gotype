// Copyright 2023 james dotter. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import (
	b "bytes"
	"unsafe"
)

// ------------------------------------------------------------ /
// GOTYPE CUSTOM TYPE IMPLEMENTATION
// implementation of custom type of []byte
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

type BYTES []byte

// Bytes returns gotype VALUE as []byte
func (v VALUE) Bytes() []byte {
	switch v.KIND() {
	case Bytes:
		return *(*[]byte)(v.ptr)
	case Bool:
		return (*BOOL)(v.ptr).Bytes()
	case Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Float32, Float64:
		return v.BytesFixedLen()
	case String:
		return []byte(*(*string)(v.ptr))
	case Array:
		return (ARRAY)(v).Bytes()
	case Map:
		return (MAP)(v).Bytes()
	case Pointer:
		return v.Elem().Bytes()
	case Slice:
		return (SLICE)(v).Bytes()
	case Struct:
		return (STRUCT)(v).Bytes()
	case Time:
		return (*TIME)(v.ptr).Bytes()
	case Uuid:
		return (*[16]byte)(v.ptr)[:]
	}
	panic("cannot convert value to []byte")
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Natural returns gotype BYTES as a golang []byte
func (b BYTES) Native() []byte {
	return []byte(b)
}

// Interface returns gotype BYTES as a golang interface{}
func (b BYTES) Interface() any {
	return b.Native()
}

// Value returns gotype BYTES as gotype VALUE
func (b BYTES) VALUE() VALUE {
	return ValueOf([]byte(b))
}

// Encode returns a gotype encoding of BOOL
func (b BYTES) Encode() []byte {
	return append([]byte{byte(Bytes)}, append(lenBytes(len(b)), b...)...)
}

// Bytes returns gotype BYTES as []byte
func (b BYTES) Bytes() []byte {
	return []byte(b)
}

// Bool returns gotype BYTES as bool
func (b BYTES) Bool() bool {
	return b[0] != 0x0 && b[0] != 0x46 && b[0] != 0x66
}

// BOOL returns gotype BYTES as gotype BOOL
func (b BYTES) BOOL() BOOL {
	return BOOL(b.Bool())
}

// Int returns gotype BYTES as int
func (b BYTES) Int() int {
	i := 0
	l := len(b)
	switch l {
	case 1, 2, 4, 8:
		p := unsafe.Pointer(&i)
		for n := 0; n < l; n++ {
			*(*uint8)(offseti(p, n)) = b[n]
		}
	default:
		panic("cannot convert to int")
	}
	return i
}

// INT returns gotype BYTES as a gotype INT
func (b BYTES) INT() INT {
	return INT(b.Int())
}

// Uint returns gotype BYTES as uint
func (b BYTES) Uint() uint {
	i := uint(0)
	l := len(b)
	switch l {
	case 1, 2, 4, 8:
		p := unsafe.Pointer(&i)
		for n := 0; n < l; n++ {
			*(*uint8)(offseti(p, n)) = b[n]
		}
	default:
		panic("cannot convert to uint")
	}
	return i
}

// Float64 returns gotype BYTES as Float64
func (b BYTES) Float64() float64 {
	i := float64(0)
	l := len(b)
	switch l {
	case 1, 2, 4, 8:
		p := unsafe.Pointer(&i)
		for n := 0; n < l; n++ {
			*(*uint8)(offseti(p, n)) = b[n]
		}
	default:
		panic("cannot convert to float64")
	}
	return i
}

// Float returns gotype BYTES as a gotype FLOAT
func (b BYTES) FLOAT() FLOAT {
	return FLOAT(b.Float64())
}

// ARRAY returns gotype BYTES as a gotype ARRAY
func (b BYTES) ARRAY() ARRAY {
	return b.SLICE().ARRAY()
}

// String returns gotype BYTES as string
func (b BYTES) String() string {
	return string(b)
}

// STRING returns gotype BYTES as a gotype STRING
func (b BYTES) STRING() STRING {
	return STRING(b.String())
}

// SLICE returns gotype BYTES as a gotype SLICE
func (b BYTES) SLICE() SLICE {
	v := b.VALUE()
	if v.TypeMatch(JSON{}) {
		return (*JSON)(v.ptr).SLICE()
	} else {
		return (SLICE)(v)
	}
}

// TIME returns gotype BYTES as a gotype TIME
func (b BYTES) TIME() TIME {
	if len(b) == 8 {
		return b.INT().TIME()
	}
	return b.STRING().TIME()
}

// UUID returns gotype BYTES as a gotype UUID
func (b BYTES) UUID() UUID {
	if b.CanUuid() {
		return UUID(b)
	}
	panic("cannot convert bytes to uuid")
}

// JSON returns gotype BYTES as gotype JSON
func (b BYTES) JSON() JSON {
	if hasJsonBookends(b) {
		return JSON(b)
	}
	panic("cannot convert string to JSON")
}

// ------------------------------------------------------------ /
// EXPANDED FUNCTIONS
// implementations of new functions for
// bytes
// referenced packages: bytes
// ------------------------------------------------------------ /

// CanBool evaluates whether BYTES are bool
// with a len of 1 and either 0 or 1
func (b BYTES) CanBool() bool {
	return len(b) == 1 && b[0] < 2
}

// MustBool panics for bytes cannot be converted to bool
func (b BYTES) MustBool() {
	if !b.CanBool() {
		panic("cannot convert bytes to bool")
	}
}

// CanString evaluates whether BYTES can be
// written as a UTF-8 string
func (b BYTES) CanString() bool {
	for _, bt := range b {
		if bt > 127 {
			return false
		}
	}
	return true
}

// MustString panics for bytes cannot be converted to utf8 string
func (b BYTES) MustString() {
	if !b.CanString() {
		panic("cannot convert bytes to utf8 string")
	}
}

// CanInt evaluates whether BYTES can be
// written as Int, Int8, Int16, Int32 or Int64
func (b BYTES) CanInt() bool {
	switch len(b) {
	case 1, 2, 4, 8:
		return true
	default:
		return false
	}
}

// MustInt panics for bytes cannot be converted to int
func (b BYTES) MustInt() {
	if !b.CanInt() {
		panic("cannot convert bytes to int")
	}
}

// CanUint evaluates whether BYTES can be
// written as Uint, Uint8, Uint16, Uint32 or Uint64
func (b BYTES) CanUint() bool {
	if b.CanInt() && b[0] < b[len(b)-1] {
		return true
	}
	return false
}

func (b BYTES) CanFloat() bool {
	switch len(b) {
	case 4, 8:
		return true
	default:
		return false
	}
}

// MustFloat panics for bytes cannot be converted to float64
func (b BYTES) MustFloat() {
	if !b.CanFloat() {
		panic("cannot convert bytes to float64")
	}
}

// CanUuid evaluates whether BYTES can be
// written as a uuid version 1-5, variant 1
func (b BYTES) CanUuid() bool {
	if len(b) == 16 {
		if (15 < b[6] && b[6] < 96) && (127 < b[8] && b[8] < 192) {
			return true
		}
	}
	return false
}

// MustUuid panics for bytes cannot be converted to uuid
func (b BYTES) MustUuid() {
	if !b.CanUuid() {
		panic("cannot convert bytes to uuid")
	}
}

// Escaped returns string with quote chars escaped
func (b BYTES) Escaped(quote byte, esc byte) BYTES {
	p := byte(0)
	for i, c := range b {
		if c == quote && p != esc {
			b = append(append(b[:i], esc), b[i:]...)
			i++
		}
		p = c
	}
	return b
}

// JoinByteSlices concatenates the bytes slices to create a new byte slice;
// The separator sep is placed between elements in the resulting slice
func JoinByteSlices(sep []byte, bytes [][]byte) []byte {
	return b.Join(bytes, sep)
}

// JoinBytesSep concatenates the bytes to create a new byte slice;
// The separator sep is placed between elements in the resulting slice
func JoinBytesSep(sep byte, bytes ...byte) []byte {
	l := len(bytes) * 2
	s := make([]byte, l)
	for i, b := range bytes {
		s[i] = b
		s[i+1] = sep
	}
	return s[:l-1]
}

// JoinBytesSep concatenates the bytes to create a new byte slice
func JoinBytes(bytes ...byte) []byte {
	return bytes
}

// RepeatBytes returns a new byte slice consisting of count copies of bytes
func RepeatBytes(bytes []byte, n int) []byte {
	return b.Repeat(bytes, n)
}

// BytesFixed returns the bytes from the value
// of a kind with a fixed number of bytes (eg. int)
func (v VALUE) BytesFixedLen() []byte {
	l := v.typ.size
	b := make([]byte, l)
	for i := uintptr(0); i < l; i++ {
		b[i] = *(*byte)(offset(v.ptr, i))
	}
	return b
}
