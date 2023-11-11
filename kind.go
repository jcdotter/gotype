// Copyright 2023 james dotter. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import (
	"unsafe"
)

// ------------------------------------------------------------ /
// KIND IMPLEMENTATION
// custom implementation of golang source code: reflect.Kind
// with expanded functionality
// ------------------------------------------------------------ /

type KIND uint8

const (
	Invalid KIND = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Pointer
	Slice
	String
	Struct
	UnsafePointer
	Field
	Time
	Uuid
	Bytes
)

var kindNames = []string{
	Invalid:       "invalid",
	Bool:          "bool",
	Int:           "int",
	Int8:          "int8",
	Int16:         "int16",
	Int32:         "int32",
	Int64:         "int64",
	Uint:          "uint",
	Uint8:         "uint8",
	Uint16:        "uint16",
	Uint32:        "uint32",
	Uint64:        "uint64",
	Uintptr:       "uintptr",
	Float32:       "float32",
	Float64:       "float64",
	Complex64:     "complex64",
	Complex128:    "complex128",
	Array:         "array",
	Chan:          "chan",
	Func:          "func",
	Interface:     "interface",
	Map:           "map",
	Pointer:       "ptr",
	Slice:         "slice",
	String:        "string",
	Struct:        "struct",
	UnsafePointer: "unsafe.Pointer",
	Field:         "field",
	Time:          "time",
	Uuid:          "uuid",
	Bytes:         "bytes",
}

var kindVals = []any{
	Invalid:    0,
	Bool:       false,
	Int:        int(0),
	Int8:       int8(0),
	Int16:      int16(0),
	Int32:      int32(0),
	Int64:      int64(0),
	Uint:       uint(0),
	Uint8:      uint8(0),
	Uint16:     uint16(0),
	Uint32:     uint32(0),
	Uint64:     uint64(0),
	Uintptr:    uintptr(0),
	Float32:    float32(0),
	Float64:    float64(0),
	Complex64:  complex64(0),
	Complex128: complex128(0),
	Array: []any{
		Bool:       [0]bool{},
		Int:        [0]int{},
		Int8:       [0]int8{},
		Int16:      [0]int16{},
		Int32:      [0]int32{},
		Int64:      [0]int64{},
		Uint:       [0]uint{},
		Uint8:      [0]uint8{},
		Uint16:     [0]uint16{},
		Uint32:     [0]uint32{},
		Uint64:     [0]uint64{},
		Float32:    [0]float32{},
		Float64:    [0]float64{},
		Complex64:  [0]complex64{},
		Complex128: [0]complex128{},
		Array:      [0]float64{},
	},
	Chan:          make(chan any),
	Func:          func() {},
	Interface:     (any)(0),
	Map:           map[any]any{},
	Pointer:       VALUE{}.ptr,
	Slice:         []any{},
	String:        "",
	Struct:        struct{}{},
	UnsafePointer: VALUE{}.ptr,
	Field:         FIELD{},
	Time:          INT(0).TIME(),
	Uuid:          UUID{},
	Bytes:         []byte{},
}

func KindOf(a any) KIND {
	return ValueOf(a).KIND()
}

func (k KIND) String() string {
	return kindNames[uint8(k)]
}

func (k KIND) STRING() STRING {
	return STRING(k.String())
}

func (k KIND) Byte() byte {
	return byte(k)
}

func (k KIND) IsBasic() bool {
	return k > 0 && (k < 15 || k == 24 || k > 27)
}

func (k KIND) IsNumeric() bool {
	return k == Int || k == Int8 || k == Int16 || k == Int32 || k == Int64 ||
		k == Uint || k == Uint8 || k == Uint16 || k == Uint32 || k == Uint64 ||
		k == Float32 || k == Float64
}

func (k KIND) IsData() bool {
	return k == Array || k == Chan || k == Map || k == Slice || k == Struct || k == Bytes || k == Interface
}

func (k KIND) CanNil() bool {
	return k == Array || k == Interface || k == Map || k == Pointer || k == Slice || k == Struct
}

func (k KIND) Size() uintptr {
	return (*VALUE)(unsafe.Pointer(&kindVals[k])).typ.size
}

func (k KIND) NewValue() VALUE {
	return (*VALUE)(unsafe.Pointer(&kindVals[k])).typ.New()
}

func (k KIND) NewArray(size int) ARRAY {
	return k.NewSlice(size).ARRAY()
}

func (k KIND) NewSlice(size int) SLICE {
	switch k {
	case Bool:
		a := make([]bool, size)
		return SliceOf(&a)
	case Int:
		a := make([]int, size)
		return SliceOf(&a)
	case Int8:
		a := make([]int8, size)
		return SliceOf(&a)
	case Int16:
		a := make([]int16, size)
		return SliceOf(&a)
	case Int32:
		a := make([]int32, size)
		return SliceOf(&a)
	case Int64:
		a := make([]int64, size)
		return SliceOf(&a)
	case Uint:
		a := make([]uint, size)
		return SliceOf(&a)
	case Uint8:
		a := make([]uint8, size)
		return SliceOf(&a)
	case Uint16:
		a := make([]uint16, size)
		return SliceOf(&a)
	case Uint32:
		a := make([]uint32, size)
		return SliceOf(&a)
	case Uint64:
		a := make([]uint64, size)
		return SliceOf(&a)
	case Uintptr:
		a := make([]uintptr, size)
		return SliceOf(&a)
	case Float32:
		a := make([]float32, size)
		return SliceOf(&a)
	case Float64:
		a := make([]float64, size)
		return SliceOf(&a)
	case Complex64:
		a := make([]complex64, size)
		return SliceOf(&a)
	case Complex128:
		a := make([]complex128, size)
		return SliceOf(&a)
	case String:
		a := make([]string, size)
		return SliceOf(&a)
	case Time:
		a := make([]TIME, size)
		return SliceOf(&a)
	case Uuid:
		a := make([]UUID, size)
		return SliceOf(&a)
	case Bytes:
		a := make([][]byte, size)
		return SliceOf(&a)
	}
	a := make([]any, size)
	return SliceOf(&a)
}

func (k KIND) NewMap() MAP {
	switch k {
	case Bool:
		return MapOf(make(map[string]bool))
	case Int:
		return MapOf(make(map[string]int))
	case Int8:
		return MapOf(make(map[string]int8))
	case Int16:
		return MapOf(make(map[string]int16))
	case Int32:
		return MapOf(make(map[string]int32))
	case Int64:
		return MapOf(make(map[string]int64))
	case Uint:
		return MapOf(make(map[string]uint))
	case Uint8:
		return MapOf(make(map[string]uint8))
	case Uint16:
		return MapOf(make(map[string]uint16))
	case Uint32:
		return MapOf(make(map[string]uint32))
	case Uint64:
		return MapOf(make(map[string]uint64))
	case Uintptr:
		return MapOf(make(map[string]uintptr))
	case Float32:
		return MapOf(make(map[string]float32))
	case Float64:
		return MapOf(make(map[string]float64))
	case Complex64:
		return MapOf(make(map[string]complex64))
	case Complex128:
		return MapOf(make(map[string]complex128))
	case String:
		return MapOf(make(map[string]string))
	case Time:
		return MapOf(make(map[string]TIME))
	case Uuid:
		return MapOf(make(map[string]UUID))
	case Bytes:
		return MapOf(make(map[string][]byte))
	}
	return MapOf(make(map[string]any))
}
