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
// implementation of custom type of int
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

type INT int

// INT returns gotype VALUE as gotype INT
func (v VALUE) INT() INT {
	return INT(v.Int())
}

// Int returns gotype VALUE as int
func (v VALUE) Int() int {
	switch v.Kind() {
	case Int:
		return int(*(*int)(v.ptr))
	case Pointer:
		return v.Elem().Int()
	default:
		switch v.KIND() {
		case Bool:
			return BOOL(*(*bool)(v.ptr)).Int()
		case Int8:
			return int(*(*int8)(v.ptr))
		case Int16:
			return int(*(*int16)(v.ptr))
		case Int32:
			return int(*(*int32)(v.ptr))
		case Int64:
			return int(*(*int64)(v.ptr))
		case Uint:
			return UINT(*(*uint)(v.ptr)).Int()
		case Uint8:
			return UINT(*(*uint8)(v.ptr)).Int()
		case Uint16:
			return UINT(*(*uint16)(v.ptr)).Int()
		case Uint32:
			return UINT(*(*uint32)(v.ptr)).Int()
		case Uint64:
			return UINT(*(*uint64)(v.ptr)).Int()
		case Float32:
			return FLOAT(*(*float32)(v.ptr)).Int()
		case Float64:
			return FLOAT(*(*float64)(v.ptr)).Int()
		case String:
			return (*STRING)(v.ptr).Int()
		case Bytes:
			return STRING(*(*[]byte)(v.ptr)).Int()
		case Time:
			return (*(*TIME)(v.ptr)).Int()
		}
	}
	panic("cannot convert value to int")
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Natural returns gotype INT as a golang int
func (i INT) Native() int {
	return int(i)
}

// Interface returns gotype INT as a golang interface{}
func (i INT) Interface() any {
	return i.Native()
}

// Value returns gotype Int as gotype Value
func (i INT) VALUE() VALUE {
	a := (any)(i)
	return *(*VALUE)(unsafe.Pointer(&a))
}

// Encode returns a gotype encoding of INT
func (i INT) Encode() ENCODING {
	return append([]byte{byte(Int)}, i.Bytes()...)
}

// Bytes returns gotype INT as []byte
func (i INT) Bytes() []byte {
	b := make([]byte, 8)
	p := uintptr(unsafe.Pointer(&i))
	b[7] = *(*uint8)(unsafe.Pointer(p))
	b[6] = *(*uint8)(unsafe.Pointer(p + 1))
	b[5] = *(*uint8)(unsafe.Pointer(p + 2))
	b[4] = *(*uint8)(unsafe.Pointer(p + 3))
	b[3] = *(*uint8)(unsafe.Pointer(p + 4))
	b[2] = *(*uint8)(unsafe.Pointer(p + 5))
	b[1] = *(*uint8)(unsafe.Pointer(p + 6))
	b[0] = *(*uint8)(unsafe.Pointer(p + 7))
	return b
}

// String returns gotype INT as string
func (i INT) String() string {
	return strconv.Itoa(int(i))
}

// Str returns gotype INT as a gotype STRING
func (i INT) STRING() STRING {
	return STRING(strconv.Itoa(int(i)))
}

// Bool returns gotype INT as bool
// false if 0, true if 1, otherwise panics
func (i INT) Bool() bool {
	if i == 0 {
		return false
	}
	return true
}

// BOOL returns gotype INT as a gotype Bool
// false if 0, true if 1, otherwise panics
func (i INT) BOOL() BOOL {
	return BOOL(i.Bool())
}

// Int returns gotype INT as int
func (i INT) Int() int {
	return int(i)
}

// Int8 returns golang int8 of gotype INT
func (i INT) Int8() int8 {
	if i > math.MaxInt8 || i < -math.MaxInt8 {
		panic("overflow error: Int greater than max int8")
	}
	return int8(i)
}

// Int16 returns golang int16 of gotype INT
func (i INT) Int16() int16 {
	if i > math.MaxInt16 || i < -math.MaxInt16 {
		panic("overflow error: Int greater than max int16")
	}
	return int16(i)
}

// Int32 returns golang int32 of gotype INT
func (i INT) Int32() int32 {
	if i > math.MaxInt32 || i < -math.MaxInt32 {
		panic("overflow error: Int greater than max int32")
	}
	return int32(i)
}

// Int16 returns golang int64 of gotype INT
func (i INT) Int64() int64 {
	return int64(i)
}

// Uint returns gotype INT as uint
func (i INT) Uint() uint {
	if i < 0 {
		panic("overflow error: Int less than 0")
	}
	return uint(i)
}

// Uint8 returns gotype INT as uint8
func (i INT) Uint8() uint8 {
	if i < 0 || i > math.MaxUint8 {
		panic("overflow error: Int greater than max uint8 or Int less than 0")
	}
	return uint8(i)
}

// Uint16 returns gotype INT as uint16
func (i INT) Uint16() uint16 {
	if i < 0 || i > math.MaxUint16 {
		panic("overflow error: Int greater than max uint16 or Int less than 0")
	}
	return uint16(i)
}

// Uint32 returns gotype INT as uint32
func (i INT) Uint32() uint32 {
	if i < 0 || i > math.MaxUint32 {
		panic("overflow error: Int greater than max uint32 or Int less than 0")
	}
	return uint32(i)
}

// Uint64 returns gotype INT as uint64
func (i INT) Uint64() uint64 {
	if i < 0 {
		panic("overflow error: Int less than 0")
	}
	return uint64(i)
}

// Float returns gotype Int as float64
func (i INT) Float64() float64 {
	return float64(i)
}

// Float returns gotype Int as a gotype FLOAT
func (i INT) FLOAT() FLOAT {
	return FLOAT(i)
}

// Time returns gotype Int as gotype TIME using
// numberic value as unix seconds since Jan 1, 1970
func (i INT) TIME() TIME {
	return TIME(time.Unix(0, int64(i)).UTC())
}

// ------------------------------------------------------------ /
// EXPANDED FUNCTIONS
// implementations of new functions for
// int
// referenced packages: NA
// ------------------------------------------------------------ /

func (i INT) IsNegative() bool {
	return i < 0
}

func (i INT) Abs() int {
	if i < 0 {
		return int(-i)
	}
	return int(i)
}

// Add increases INT by val, panics if exceeds overflow limit
func (i INT) Add(val int) int {
	MustIntAdd(int(i), val)
	return int(i) + val
}

// Add decreases INT by val, panics if exceeds overflow limit
func (i INT) Sub(val int) int {
	MustIntAdd(int(i), -val)
	return int(i) - val
}

// Multiply mutiplies Int by val, , panics if exceeds overflow limit
func (i INT) Multiply(val int) int {
	MustIntMultiply(int(i), -val)
	return int(i) - val
}

// Divide divides Int by val and returns rounded result
func (i INT) Divide(val int) int {
	return int(math.Round(float64(i) / float64(val)))
}

// Mod returns the remainder of INT divided by val
func (i INT) Mod(val int) int {
	return int(math.Mod(float64(i), float64(val)))
}

// Pow returns INT raised to val, panics if exceeds overflow limit
func (i INT) Pow(val int) int {
	return int(math.Mod(float64(i), float64(val)))
}

// MustIntAdd panic if x + y exceeds math.MaxInt (positive or negative)
func MustIntAdd(x int, y int) {
	if x < 0 && y < 0 {
		if math.MaxInt+x+y < 0 {
			panic("overflow error: less than max int")
		}
	} else if x > 0 && y > 0 {
		if math.MaxInt-x-y < 0 {
			panic("overflow error: greater than max int")
		}
	}
}

// MustIntMultiply panic if x * y exceeds math.MaxInt (positive or negative)
func MustIntMultiply(x int, y int) {
	if math.MaxInt/AbsInt(y) < AbsInt(x) {
		panic("overflow error: greater than max int")
	}
}

// MustIntPow panic if x ^ y exceeds math.MaxInt (positive or negative)
func MustIntPow(x int, y int) {
	if math.Pow(math.MaxInt, 1/Abs(float64(y))) < Abs(float64(x)) {
		panic("overflow error: greater than max int")
	}
}
