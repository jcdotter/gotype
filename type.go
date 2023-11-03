// Copyright 2023 james dotter. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import (
	"reflect"
	"time"
	"unsafe"
)

// ------------------------------------------------------------ /
// TYPE IMPLEMENTATION
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

	tflagUncommon tflag = 1 << 0
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
type typeOff int32 // offset to a *TYPE

type TYPE struct {
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

// TypeOf returns the gotype of value a
func TypeOf(a any) *TYPE {
	return *(**TYPE)(unsafe.Pointer(&a))
}

// FromReflectType returns the gotype of reflect.Type t
func FromReflectType(t reflect.Type) *TYPE {
	a := (any)(t)
	return (*TYPE)((*VALUE)(unsafe.Pointer(&a)).ptr)
}

// New returns an empty pointer to a new value of the TYPE
func (t *TYPE) New() VALUE {
	if t != nil {
		return VALUE{t.PtrType(), unsafe_New(t), flag(Pointer)}
	}
	panic("call to New on nil type")
}

// New returns a pointer to a new (non nil) value of the TYPE
func (t *TYPE) NewValue() VALUE {
	n := t.New()
	switch t.Kind() {
	case Map:
		*(*unsafe.Pointer)(n.ptr) = makemap(t, 0, nil)
	case Pointer:
		*(*unsafe.Pointer)(n.ptr) = t.Elem().NewValue().ptr
	case Slice:
		t := (*sliceType)(unsafe.Pointer(t)).elem
		*(*unsafe.Pointer)(&n.ptr) = unsafe.Pointer(&sliceHeader{unsafe_NewArray(t, 0), 0, 0})
	}
	return n
}

// Reflect returns the reflect.Type of the TYPE
func (t *TYPE) Reflect() reflect.Type {
	return toType(t)
}

// PtrType returns a new TYPE of a pointer to the TYPE
func (t *TYPE) PtrType() *TYPE {
	return FromReflectType(reflect.PtrTo(toType(t)))
}

// IfaceIndir returns true if the TYPE is an indirect value
func (t *TYPE) IfaceIndir() bool {
	return t.kind&KindDirectIface == 0
}

// flag returns the flag of the TYPE
func (t *TYPE) flag() flag {
	f := flag(t.kind & kindMask)
	if t.IfaceIndir() {
		f |= flagIndir
	}
	return f
}

// Kind returns the KIND of the TYPE synonomous with reflect.Kind
func (t *TYPE) Kind() KIND {
	return KIND(t.kind & kindMask)
}

// KIND returns the gotype KIND of the TYPE
// which includes Bytes, Field, Time, Uuid
func (t *TYPE) KIND() KIND {
	k := KIND(t.kind & kindMask)
	switch k {
	case Slice: // check if byte array
		ek := (*sliceType)(unsafe.Pointer(t)).elem.kind & kindMask
		if ek == 8 || ek == 10 { // []byte or []rune
			return Bytes
		}
	case Array: // check if uuid
		a := (*arrayType)(unsafe.Pointer(t))
		if a.len == 16 && a.elem.kind&kindMask == 8 { // [16]byte
			return Uuid
		}
	case Struct: // check if time
		fs := (*structType)(unsafe.Pointer(t)).fields
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

// Elem returns the TYPE of the element of the TYPE
func (t *TYPE) Elem() *TYPE {
	switch t.Kind() {
	case Array:
		return (*arrayType)(unsafe.Pointer(t)).elem
	case Map:
		return (*mapType)(unsafe.Pointer(t)).elem
	case Pointer:
		return (*ptrType)(unsafe.Pointer(t)).elem
	case Slice:
		return (*sliceType)(unsafe.Pointer(t)).elem
	}
	return t
}

// String returns the string representation of the TYPE
func (t *TYPE) String() string {
	return t.Name()
}

// STRING returns the gotype STRING representation of the TYPE
func (t *TYPE) STRING() STRING {
	return STRING(t.String())
}

// Name returns the name of the TYPE
func (t *TYPE) Name() string {
	n := name{(*byte)(resolveNameOff(unsafe.Pointer(t), int32(t.str)))}.name()
	if t.Kind() != Pointer {
		n = n[1:]
	}
	return n
}

// NameShort returns the short name of the TYPE
// excluding the package path, module name and pointer indicator
func (t *TYPE) NameShort() string {
	n := t.STRING()
	return string(n[n.LastIndexOf(".")+1:])
}

// ------------------------------------------------------------ /
// STURCTURED TYPES
// implementation of golang types for data structures:
// array, map, ptr, slice, string, struct, field, interface
// ------------------------------------------------------------ /

type arrayType struct {
	TYPE
	elem  *TYPE // array element type
	slice *TYPE // slice type
	len   uintptr
}

type funcType struct {
	TYPE
	inCount  uint16
	outCount uint16
}

type mapType struct {
	TYPE
	key    *TYPE // map key type
	elem   *TYPE // map element (value) type
	bucket *TYPE // internal bucket structure
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
	TYPE
	elem *TYPE
}

type sliceType struct {
	TYPE
	elem *TYPE // slice element type
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
	TYPE
	pkgPath name
	fields  []fieldType // sorted by offset
}

type fieldType struct {
	name   name    // name is always non-empty
	typ    *TYPE   // type of field
	offset uintptr // byte offset of field
}

type interfaceType struct {
	*TYPE
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
// STRUCT TYPE IMPLEMENTATION
// custom implementation of golang struct type
// ------------------------------------------------------------ /

// IsStruct returns true if the TYPE is a struct
func (t *TYPE) IsStruct() bool {
	return t.Kind() == Struct
}

// PkgPath returns the package path of a struct TYPE
func (t *TYPE) PkgPath() string {
	return (*structType)(unsafe.Pointer(t)).pkgPath.name()
}

// NumField returns the number of fields in a struct TYPE
func (t *TYPE) NumField() int {
	return len((*structType)(unsafe.Pointer(t)).fields)
}

// Field returns the TYPE of the field at index i in a struct TYPE
func (t *TYPE) Field(i int) *TYPE {
	return (*structType)(unsafe.Pointer(t)).fields[i].typ
}

// FieldByName returns the TYPE of the field with name in a struct TYPE
func (t *TYPE) FieldByName(name string) *TYPE {
	fs := (*structType)(unsafe.Pointer(t)).fields
	for _, f := range fs {
		if f.name.name() == name {
			return f.typ
		}
	}
	return nil
}

// FieldByTag returns the TYPE of the field with tag value in a struct TYPE
func (t *TYPE) FieldByTag(tag string, value string) *TYPE {
	fs := (*structType)(unsafe.Pointer(t)).fields
	for _, f := range fs {
		v := getTagValue(f.name.tag(), tag, 34)
		if v == value {
			return f.typ
		}
	}
	return nil
}

// FieldByIndex returns the TYPE of the field at index in a struct TYPE
func (t *TYPE) FieldByIndex(index []int) *TYPE {
	if len(index) == 0 {
		return t
	}
	return (*structType)(unsafe.Pointer(t)).fields[index[0]].typ.FieldByIndex(index[1:])
}

// FieldName returns the name of the field at index i in a struct TYPE
func (t *TYPE) FieldName(i int) string {
	return (*structType)(unsafe.Pointer(t)).fields[i].name.name()
}

// FieldIndex returns the index of the field with name in a struct TYPE
func (t *TYPE) FieldIndex(name string) int {
	fs := (*structType)(unsafe.Pointer(t)).fields
	for i, f := range fs {
		if f.name.name() == name {
			return i
		}
	}
	return 0
}

// IndexTag returns the tag of the field at index i in a struct TYPE
func (t *TYPE) IndexTag(i int) string {
	return (*structType)(unsafe.Pointer(t)).fields[i].name.tag()
}

// FieldTag returns the tag of the field with name in a struct TYPE
func (t *TYPE) FieldTag(name string) string {
	fs := (*structType)(unsafe.Pointer(t)).fields
	for _, f := range fs {
		if f.name.name() == name {
			return f.name.tag()
		}
	}
	return ""
}

// IndexTagValue returns the value of the tag of the field at index i in a struct TYPE
func (t *TYPE) IndexTagValue(i int, tag string) string {
	return getTagValue(t.IndexTag(i), tag, 34)
}

// FieldTagValue returns the value of the tag of the field with name in a struct TYPE
func (t *TYPE) FieldTagValue(name string, tag string) string {
	return getTagValue(t.FieldTag(name), tag, 34)
}

// ForFields iterates over the fields of a struct TYPE and calls
// the function f with the index and TYPE of each field
func (t *TYPE) ForFields(f func(i int, n string, typ *TYPE) (brake bool)) {
	for i, fld := range (*structType)(unsafe.Pointer(t)).fields {
		if brake := f(i, fld.name.name(), fld.typ); brake {
			break
		}
	}
}

// matchStructType compairs the structure of 2 structs
func matchStructType(x, y *TYPE) bool {
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

// ------------------------------------------------------------ /
// FUNC TYPE IMPLEMENTATION
// custom implementation of golang func type
// ------------------------------------------------------------ /

// IsFunc returns true if the TYPE is a func
func (t *TYPE) IsFunc() bool {
	return t.Kind() == Func
}

// NumIn returns the number of input parameters in a func TYPE
func (t *TYPE) NumIn() int {
	return int((*funcType)(unsafe.Pointer(t)).inCount)
}

// NumOut returns the number of output parameters in a func TYPE
func (t *TYPE) NumOut() int {
	return int((*funcType)(unsafe.Pointer(t)).outCount)
}

// In returns the TYPE of the input parameter at index i in a func TYPE
func (t *TYPE) In(i int) *TYPE {
	return (*funcType)(unsafe.Pointer(t)).in()[i]
}

// Out returns the TYPE of the output parameter at index i in a func TYPE
func (t *TYPE) Out(i int) *TYPE {
	return (*funcType)(unsafe.Pointer(t)).out()[i]
}

// in returns a slice of TYPEs of the input parameters in a func TYPE
func (t *funcType) in() []*TYPE {
	uadd := unsafe.Sizeof(*t)
	if t.tflag&tflagUncommon != 0 {
		uadd += 32 //unsafe.Sizeof(uncommonType{})
	}
	if t.inCount == 0 {
		return nil
	}
	return (*[1 << 20]*TYPE)(add(unsafe.Pointer(t), uadd))[:t.inCount:t.inCount]
}

// out returns a slice of TYPEs of the output parameters in a func TYPE
func (t *funcType) out() []*TYPE {
	uadd := unsafe.Sizeof(*t)
	if t.tflag&tflagUncommon != 0 {
		uadd += 32 //size of uncommonType
	}
	outCount := t.outCount & (1<<15 - 1)
	if outCount == 0 {
		return nil
	}
	return (*[1 << 20]*TYPE)(add(unsafe.Pointer(t), uadd))[t.inCount : t.inCount+outCount : t.inCount+outCount]
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
