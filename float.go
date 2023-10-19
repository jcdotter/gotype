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
// implementation of custom type of float
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

type FLOAT float64

// FLOAT returns gotype VALUE as gotype FLOAT
func (v VALUE) FLOAT() FLOAT {
	return FLOAT(v.Float64())
}

// Float64 returns gotype VALUE as float64
func (v VALUE) Float64() float64 {
	switch v.Kind() {
	case Float64:
		return *(*float64)(v.ptr)
	case Pointer:
		return v.Elem().Float64()
	default:
		switch v.KIND() {
		case Bool:
			return BOOL(*(*bool)(v.ptr)).Float64()
		case Int:
			return float64(*(*int)(v.ptr))
		case Int8:
			return float64(*(*int8)(v.ptr))
		case Int16:
			return float64(*(*int16)(v.ptr))
		case Int32:
			return float64(*(*int32)(v.ptr))
		case Int64:
			return float64(*(*int64)(v.ptr))
		case Uint:
			return float64(*(*uint)(v.ptr))
		case Uint8:
			return float64(*(*uint8)(v.ptr))
		case Uint16:
			return float64(*(*uint16)(v.ptr))
		case Uint32:
			return float64(*(*uint32)(v.ptr))
		case Uint64:
			return float64(*(*uint64)(v.ptr))
		case Float32:
			return float64(*(*float32)(v.ptr))
		case String:
			return (*STRING)(v.ptr).Float64()
		case Bytes:
			return (*BYTES)(v.ptr).Float64()
		case Time:
			return (*(*TIME)(v.ptr)).Float64()
		}
	}
	panic("cannot convert value to int")
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Natural returns gotype Float as a golang float64
func (f FLOAT) Native() float64 {
	return float64(f)
}

// Interface returns gotype Float as a golang interface{}
func (f FLOAT) Interface() any {
	return f.Native()
}

// VALUE returns gotype Float as gotype VALUE
func (f FLOAT) VALUE() VALUE {
	return ValueOf(float64(f))
}

// Encode returns a gotype encoding of FLOAT
func (f FLOAT) Encode() ENCODING {
	return append([]byte{byte(Float64)}, f.Bytes()...)
}

// Bytes returns gotype FLOAT as []byte
func (f FLOAT) Bytes() []byte {
	b := make([]byte, 8)
	p := unsafe.Pointer(&f)
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

// String returns gotype Float as gotype String
func (f FLOAT) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

// Str returns gotype Float as a gotype STRING
func (f FLOAT) STRING() STRING {
	return STRING(f.String())
}

// Bool returns gotype Float as a gotype Bool
// false if 0, otherwise true
func (f FLOAT) Bool() bool {
	return f != 0
}

// BOOL returns gotype FLOAT as a gotype Bool
// false if 0, otherwise true
func (f FLOAT) BOOL() BOOL {
	return BOOL(f.Bool())
}

// Int returns gotype Float as int
func (f FLOAT) Int() int {
	if f > math.MaxInt && f < -math.MaxInt {
		panic("overflow error: Float greater than max uint or Float less than 0")
	}
	return int(math.Round(float64(f)))
}

// INT returns gotype Float as a gotype INT
func (f FLOAT) INT() INT {
	return INT(f.Int())
}

// Uint returns gotype FLOAT as uint
func (f FLOAT) Uint() uint {
	if f < 0 || f > math.MaxUint {
		panic("overflow error: Float greater than max uint or Float less than 0")
	}
	return uint(math.Round(float64(f)))
}

// UINT returns gotype Float as a gotype UINT
func (f FLOAT) UINT() UINT {
	return UINT(f.Uint())
}

// Float32 returns golang float32 of gotype FLOAT
func (f FLOAT) Float32() float32 {
	if f > math.MaxFloat32 || f < -math.MaxFloat32 {
		panic("overflow error: Float greater than max float32")
	}
	return float32(f)
}

// Float returns gotype Float as a gotype Float
func (f FLOAT) Float64() float64 {
	return float64(f)
}

// TIME returns gotype Float as gotype TIME using
// numberic value as unix seconds since Jan 1, 1970
func (f FLOAT) TIME() TIME {
	return TIME(time.Unix(0, int64(f)).UTC())
}

// ------------------------------------------------------------ /
// EXPANDED FUNCTIONS
// implementations of new functions for
// float64
// referenced packages: NA
// ------------------------------------------------------------ /

func (f FLOAT) IsNegative() bool {
	return f < 0
}

func (f FLOAT) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

// Add increases INT by val, panics if exceeds overflow limit
func (f FLOAT) Add(val float64) float64 {
	MustFloatAdd(float64(f), val)
	return float64(f) + val
}

// Add decreases INT by val, panics if exceeds overflow limit
func (f FLOAT) Sub(val float64) float64 {
	MustFloatAdd(float64(f), -val)
	return float64(f) - val
}

// Multiply mutiplies FLOAT by val, , panics if exceeds overflow limit
func (f FLOAT) Multiply(val float64) float64 {
	MustFloatMultiply(float64(f), -val)
	return float64(f) - val
}

// Divide divides FLOAT by val and returns rounded result
func (f FLOAT) Divide(val float64) float64 {
	return math.Round(float64(f) / val)
}

// Mod returns the remainder of FLOAT divided by val
func (f FLOAT) Mod(val float64) float64 {
	return math.Mod(float64(f), val)
}

// Pow returns FLOAT raised to val, panics if exceeds overflow limit
func (f FLOAT) Pow(val float64) float64 {
	return math.Mod(float64(f), val)
}

// MustFloatAdd panic if x + y exceeds math.MaxFloat64 (positive or negative)
func MustFloatAdd(x float64, y float64) {
	if x < 0 && y < 0 {
		if math.MaxFloat64+x+y < 0 {
			panic("overflow error: less than max int")
		}
	} else if x > 0 && y > 0 {
		if math.MaxFloat64-x-y < 0 {
			panic("overflow error: greater than max int")
		}
	}
}

// MustIntMultiply panic if x * y exceeds math.MaxFloat64 (positive or negative)
func MustFloatMultiply(x float64, y float64) {
	if math.MaxFloat64/Abs(y) < Abs(x) {
		panic("overflow error: greater than max int")
	}
}

// MustFloatPow panic if x ^ y exceeds math.MaxFloat64 (positive or negative)
func MustFloatPow(x float64, y float64) {
	if math.Pow(math.MaxFloat64, 1/Abs(y)) < Abs(x) {
		panic("overflow error: greater than max int")
	}
}
