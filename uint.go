// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	"math"
	"strconv"
	"time"
	"unsafe"
)

// ------------------------------------------------------------ /
// GOTYPE CUSTOM TYPE IMPLEMENTATION
// implementation of custom type of uint
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

type UINT uint

// UINT returns gotype VALUE as gotype UINT
func (v VALUE) UINT() UINT {
	return UINT(v.Uint())
}

// Uint returns gotype VALUE as uint
func (v VALUE) Uint() uint {
	switch v.Kind() {
	case Uint:
		return uint(*(*uint)(v.ptr))
	case Pointer:
		return v.Elem().Uint()
	default:
		switch v.KIND() {
		case Bool:
			return uint(BOOL(*(*bool)(v.ptr)).Int())
		case Int:
			return INT(*(*int)(v.ptr)).Uint()
		case Int8:
			return INT(*(*int8)(v.ptr)).Uint()
		case Int16:
			return INT(*(*int16)(v.ptr)).Uint()
		case Int32:
			return INT(*(*int32)(v.ptr)).Uint()
		case Int64:
			return INT(*(*int64)(v.ptr)).Uint()
		case Uint8:
			return uint(*(*uint8)(v.ptr))
		case Uint16:
			return uint(*(*uint16)(v.ptr))
		case Uint32:
			return uint(*(*uint32)(v.ptr))
		case Uint64:
			return uint(*(*uint64)(v.ptr))
		case Float32:
			return FLOAT(*(*float32)(v.ptr)).Uint()
		case Float64:
			return FLOAT(*(*float64)(v.ptr)).Uint()
		case String:
			return (*STRING)(v.ptr).Uint()
		case Bytes:
			return (*BYTES)(v.ptr).Uint()
		case Time:
			return (*(*TIME)(v.ptr)).Uint()
		}
	}
	panic("cannot convert value to uint")
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Natural returns gotype INT as a golang int
func (u UINT) Native() int {
	return int(u)
}

// Interface returns gotype INT as a golang interface{}
func (u UINT) Interface() any {
	return u.Native()
}

// Value returns gotype Int as gotype Value
func (u UINT) VALUE() VALUE {
	a := (any)(u)
	return *(*VALUE)(unsafe.Pointer(&a))
}

// Encode returns a gotype encoding of INT
func (u UINT) Encode() ENCODING {
	return append([]byte{byte(Int)}, u.Bytes()...)
}

// Bytes returns gotype INT as []byte
func (u UINT) Bytes() []byte {
	b := make([]byte, 8)
	p := unsafe.Pointer(&u)
	b[0] = *(*uint8)(offset(p, 0))
	b[1] = *(*uint8)(offset(p, 1))
	b[2] = *(*uint8)(offset(p, 2))
	b[3] = *(*uint8)(offset(p, 3))
	b[4] = *(*uint8)(offset(p, 4))
	b[5] = *(*uint8)(offset(p, 5))
	b[6] = *(*uint8)(offset(p, 6))
	b[7] = *(*uint8)(offset(p, 7))
	return b
}

// String returns gotype INT as string
func (u UINT) String() string {
	return strconv.Itoa(int(u))
}

// Str returns gotype INT as a gotype STRING
func (u UINT) STRING() STRING {
	return STRING(strconv.Itoa(int(u)))
}

// Bool returns gotype INT as bool
// false if 0, otherwise true
func (u UINT) Bool() bool {
	return u != 0
}

// BOOL returns gotype INT as a gotype Bool
// false if 0, otherwise true
func (u UINT) BOOL() BOOL {
	return BOOL(u.Bool())
}

// Int returns gotype UINT as int
func (u UINT) Int() int {
	if u > math.MaxInt {
		panic("overflow error: Uint greater than max int")
	}
	return int(u)
}

// Int8 returns gotype UINT as int8
func (u UINT) Int8() int8 {
	if u > math.MaxInt8 {
		panic("overflow error: Uint greater than max int8")
	}
	return int8(u)
}

// Int16 returns gotype UINT as int16
func (u UINT) Int16() int16 {
	if u > math.MaxInt16 {
		panic("overflow error: Uint greater than max int16")
	}
	return int16(u)
}

// Int32 returns gotype UINT as int32
func (u UINT) Int32() int32 {
	if u > math.MaxInt32 {
		panic("overflow error: Uint greater than max int32")
	}
	return int32(u)
}

// Int64 returns gotype UINT as int64
func (u UINT) Int64() int64 {
	if u > math.MaxInt {
		panic("overflow error: Uint greater than max int64")
	}
	return int64(u)
}

// Uint returns gotype UINT as uint
func (u UINT) Uint() uint {
	return uint(u)
}

// Uint8 returns gotype UINT as uint8
func (u UINT) Uint8() uint8 {
	if u > math.MaxUint8 {
		panic("overflow error: Uint greater than max uint8")
	}
	return uint8(u)
}

// Uint16 returns gotype UINT as uint16
func (u UINT) Uint16() uint16 {
	if u > math.MaxUint16 {
		panic("overflow error: Int greater than max uint16")
	}
	return uint16(u)
}

// Uint32 returns gotype UINT as uint32
func (u UINT) Uint32() uint32 {
	if u > math.MaxUint32 {
		panic("overflow error: Int greater than max uint32")
	}
	return uint32(u)
}

// Uint64 returns gotype UINT as uint64
func (u UINT) Uint64() uint64 {
	return uint64(u)
}

// Float returns gotype UINT as float64
func (u UINT) Float64() float64 {
	return float64(u)
}

// Float returns gotype UINT as a gotype FLOAT
func (u UINT) FLOAT() FLOAT {
	return FLOAT(u)
}

// Time returns gotype UINT as gotype TIME using
// numberic value as unix seconds since Jan 1, 1970
func (u UINT) TIME() TIME {
	return TIME(time.Unix(0, int64(u)).UTC())
}
