// Copyright 2023 escend llc. All rights reserved.typVal
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	"reflect"
	"time"
	"unsafe"
)

// ------------------------------------------------------------ /
// RKIND IMPLEMENTATION
// custom implementation of golang source code: reflect.Type
// with expanded functionality
// ------------------------------------------------------------ /

const (
	kindMask             = (1 << 5) - 1
	KindDirectIface      = 1 << 5
	flagStickyRO    flag = 1 << 5
	flagEmbedRO     flag = 1 << 6
	flagIndir       flag = 1 << 7
	flagAddr        flag = 1 << 8
)

var (
	__ttime       = (any)(time.Time{})
	__timefields  = (*structType)(unsafe.Pointer((*VALUE)(unsafe.Pointer(&__ttime)).typ)).fields
	__timefield0b = __timefields[0].name.bytes
	__timefield0t = __timefields[0].typ
	__timefield1b = __timefields[1].name.bytes
	__timefield1t = __timefields[1].typ
	__timefield2b = __timefields[2].name.bytes
	__timefield2t = __timefields[2].typ
)

type flag uintptr
type tflag uint8
type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype

type rval struct {
	typ *rtype
	ptr unsafe.Pointer
	flag
}

type rtype struct {
	size       uintptr
	ptrdata    uintptr // number of bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      tflag   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct field with this type
	kind       uint8   // enumeration for C
	equal      func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata     *byte   // garbage collection data
	str        nameOff // string form
	ptrToThis  typeOff // type for pointer to this type, may be zero
}

func getrtype(a any) *rtype {
	return *(**rtype)(unsafe.Pointer(&a))
}

func reflectType(t reflect.Type) *rtype {
	a := (any)(t)
	return (*rtype)((*VALUE)(unsafe.Pointer(&a)).ptr)
}

// New returns a new indirect Value of the Type
func (r *rtype) New(length ...int) VALUE {
	if r != nil {
		t := r.elem()
		var p unsafe.Pointer
		switch t.Kind() {
		case Array:
			a := (*arrayType)(unsafe.Pointer(r))
			n := unsafe_NewArray(a.elem, int(a.len))
			p = unsafe.Pointer(&n)
		case Map:
			n := makemap(t, 0, nil)
			p = unsafe.Pointer(&n)
		case Pointer:
			n := unsafe_New(t)
			p = unsafe.Pointer(&n)
		case Slice:
			ln := 0
			if len(length) > 0 {
				ln = length[0]
			}
			s := (*sliceType)(unsafe.Pointer(t))
			p = unsafe.Pointer(&sliceHeader{unsafe_NewArray(s.elem, ln), ln, ln})
		default:
			p = unsafe_New(t)
		}
		return VALUE{t.ptrType(), p, flag(Pointer)}
		//return FromReflect(reflect.New(toType(r.elem())))
	}
	panic("call to New on nil type")
}

func (r *rtype) ptrType() *rtype {
	return reflectType(reflect.PtrTo(toType(r)))
}

// newPtr returns a pointer to a new value instance of the Type
func (r *rtype) newPtr() unsafe.Pointer {
	return r.New().ptr
}

func (r *rtype) IfaceIndir() bool {
	return r.kind&KindDirectIface == 0
}

func (r *rtype) flag() flag {
	f := flag(r.kind & kindMask)
	if r.IfaceIndir() {
		f |= flagIndir
	}
	return f
}

func (r *rtype) Kind() KIND {
	return KIND(r.kind & kindMask)
}

func (r *rtype) KIND() KIND {
	k := KIND(r.kind & kindMask)
	switch k {
	case Slice: // check if byte array
		ek := (*sliceType)(unsafe.Pointer(r)).elem.kind & kindMask
		if ek == 8 || ek == 10 { // []byte or []rune
			return Bytes
		}
	case Array: // check if uuid
		a := (*arrayType)(unsafe.Pointer(r))
		if a.len == 16 && a.elem.kind&kindMask == 8 { // [16]byte
			return Uuid
		}
	case Struct: // check if time
		fs := (*structType)(unsafe.Pointer(r)).fields
		if len(fs) == 3 {
			if fs[0].name.bytes != __timefield0b || fs[0].typ != __timefield0t ||
				fs[1].name.bytes != __timefield1b || fs[1].typ != __timefield1t ||
				fs[2].name.bytes != __timefield2b || fs[2].typ != __timefield2t {
				return Struct
			} else {
				return Time
			}
		}
	}
	return k
}

func (r *rtype) elem() *rtype {
	if r.Kind() == Pointer {
		return (*ptrType)(unsafe.Pointer(r)).elem
	}
	return r
}

func (r *rtype) String() string {
	return r.Name()
}

func (r *rtype) STRING() STRING {
	return STRING(r.String())
}

func (r *rtype) Name() string {
	n := name{(*byte)(resolveNameOff(unsafe.Pointer(r), int32(r.str)))}.name()
	if r.Kind() != Pointer {
		n = n[1:]
	}
	return n
}

// matchStructType compairs the structure of 2 structs
func matchStructType(x, y *rtype) bool {
	if x.kind&kindMask == 25 && y.kind&kindMask == 25 {
		xfs := (*structType)(unsafe.Pointer(x)).fields
		yfs := (*structType)(unsafe.Pointer(y)).fields
		if len(xfs) == len(yfs) {
			for i, xf := range xfs {
				yf := yfs[i]
				if xf.name.bytes != yf.name.bytes || xf.typ != yf.typ {
					return false
				}
			}
			return true
		}
	}
	return false
}

type arrayType struct {
	rtype
	elem  *rtype // array element type
	slice *rtype // slice type
	len   uintptr
}

type mapType struct {
	rtype
	key    *rtype // map key type
	elem   *rtype // map element (value) type
	bucket *rtype // internal bucket structure
	// function for hashing keys (ptr to key, seed) -> hash
	hasher     func(unsafe.Pointer, uintptr) uintptr
	keysize    uint8  // size of key slot
	valuesize  uint8  // size of value slot
	bucketsize uint16 // size of bucket
	flags      uint32
}

type bmap struct {
	_ [bucketCnt]uint8
}

const (
	bucketCntBits = 3
	bucketCnt     = 1 << bucketCntBits
	dataOffset    = unsafe.Offsetof(struct {
		b bmap
		v int64
	}{}.v)
)

type hiter struct {
	_ unsafe.Pointer    // key
	_ unsafe.Pointer    // elem
	_ unsafe.Pointer    // t
	_ unsafe.Pointer    // h
	_ unsafe.Pointer    // buckets
	_ unsafe.Pointer    // bptr
	_ *[]unsafe.Pointer // overflow
	_ *[]unsafe.Pointer // oldoverflow
	_ uintptr           // startBucket
	_ uint8             // offset
	_ bool              // wrapped
	_ uint8             // B
	_ uint8             // i
	_ uintptr           // bucket
	_ uintptr           // checkBucket
}

type ptrType struct {
	rtype
	elem *rtype
}

type sliceType struct {
	rtype
	elem *rtype // slice element type
}

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

type stringHeader struct {
	Data uintptr
	Len  int
}

type structType struct {
	rtype
	pkgPath name
	fields  []fieldType // sorted by offset
}

type fieldType struct {
	name   name    // name is always non-empty
	typ    *rtype  // type of field
	offset uintptr // byte offset of field
}

type interfaceType struct {
	*rtype
	PkgPath name      // import path
	Methods []Imethod // sorted by hash
}

func (i *interfaceType) NumMethod() int {
	return len(i.Methods)
}

type Imethod struct {
	_ nameOff
	_ typeOff
}

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

func (k KIND) CanNil() bool {
	return k == Array || k == Interface || k == Map || k == Pointer || k == Slice || k == Struct
}

func (k KIND) Size() uintptr {
	return (*VALUE)(unsafe.Pointer(&kindVals[k])).typ.size
}

func (k KIND) NewValue() VALUE {
	return (*VALUE)(unsafe.Pointer(&kindVals[k])).typ.New()
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

// ------------------------------------------------------------ /
// NAME IMPLEMENTATION
// custom implementation of golang source code: name
// with expanded functionality
// ------------------------------------------------------------ /

type name struct {
	bytes *byte
}

func (n name) data(off int) *byte {
	return (*byte)(add(unsafe.Pointer(n.bytes), uintptr(off)))
}

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

func (n name) readVarint(off int) (int, int) {
	v := 0
	for i := 0; ; i++ {
		x := *n.data(off + i)
		v += int(x&0x7f) << (7 * i)
		if x&0x80 == 0 {
			return i + 1, v
		}
	}
}

func (n name) name() string {
	if n.bytes == nil {
		return ""
	}
	i, l := n.readVarint(1)
	return unsafe.String(n.data(1+i), l)
}

func (n name) hasTag() bool {
	return (*n.bytes)&(1<<1) != 0
}

func (n name) tag() string {
	if !n.hasTag() {
		return ""
	}
	i, l := n.readVarint(1)
	i2, l2 := n.readVarint(1 + i + l)
	return unsafe.String(n.data(1+i+l+i2), l2)
}
