// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	"reflect"
	"strconv"
	"unsafe"
)

// ------------------------------------------------------------ /
// VALUE IMPLEMENTATION
// inspired by golang standard reflect.Value
// with expanded methods and type conversations
// ------------------------------------------------------------ /

type VALUE struct {
	typ *rtype
	ptr unsafe.Pointer
	flag
}

func ValueOf(a any) VALUE {
	v := reflect.ValueOf(a)
	return FromReflect(v)
}

func ValueOfV(a any) VALUE {
	if n, is := a.(VALUE); !is {
		return ValueOf(a)
	} else {
		return n
	}
}

func FromReflect(v reflect.Value) VALUE {
	return *(*VALUE)(unsafe.Pointer(&v))
}

func (v VALUE) Reflect() reflect.Value {
	return *(*reflect.Value)(unsafe.Pointer(&v))
}

func (v VALUE) Init() VALUE {
	if v.ptr == nil {
		return v.typ.NewValue().Elem()
	}
	k := v.Kind()
	if (k == Map || k == Slice || k == Pointer) && *(*unsafe.Pointer)(v.ptr) == nil {
		*(*unsafe.Pointer)(v.ptr) = *(*unsafe.Pointer)(v.typ.NewValue().Elem().ptr)
	}
	return v
}

// ------------------------------------------------------------ /
// VALUE INITIALIZATION
// generates new mem address for VALUE data type when VALUE is nil
// ------------------------------------------------------------ /

// SetType sets the actual data type of interface VALUE
func (v VALUE) SetType() VALUE {
	if v.Kind() == Interface {
		var e *VALUE
		if (*interfaceType)(unsafe.Pointer(v.typ)).NumMethod() != 0 {
			e = (*VALUE)(unsafe.Pointer((*interface{ M() })(v.ptr)))
		} else {
			e = (*VALUE)(v.ptr)
		}
		if e.typ != nil && e.typ.kind != 0 {
			return VALUE{e.typ, e.ptr, e.typ.flag()}
		}
	}
	return v
}

// New returns a new empty value of VALUE type,
// if ptr is nil and init != true, returns nil ptr
func (v VALUE) New(init ...bool) VALUE {
	if v.ptr == nil && !init[0] {
		return v
	}
	return v.typ.New()
}

// NewDeep returns a new empty value of VALUE type
// with matching number of elements and no nil spaces
func (v VALUE) NewDeep() VALUE {
	v = v.SetType()
	n := v.New()
	indir := func(t *rtype) bool {
		k := t.Kind()
		return k == Array || k == Struct
	}
	set := func(e VALUE, p unsafe.Pointer, t *rtype) {
		if e.ptr == nil {
			return
		}
		nv := e.NewDeep().Elem()
		kind := t.Kind()
		if kind == Interface {
			*(*any)(p) = nv.Interface()
		} else if !indir(e.typ) {
			typedmemmove(e.typ, p, nv.ptr)
		} else {
			*(*unsafe.Pointer)(p) = nv.ptr
		}
	}
	switch v.Kind() {
	case Array:
		t := (*arrayType)(unsafe.Pointer(v.typ)).elem
		if !t.Kind().IsBasic() {
			a := (ARRAY)(n.Elem())
			(ARRAY)(v).ForEach(func(i int, k string, e VALUE) (brake bool) {
				set(e, offset(a.ptr, uintptr(i)*t.size), t)
				return
			})
		}
	case Map:
		*(*unsafe.Pointer)(n.ptr) = makemap(v.typ, (MAP)(v).Len(), nil)
		m := (MAP)(n.Elem())
		(MAP)(v).ForEach(func(i int, k string, e VALUE) (brake bool) {
			if e.ptr == nil {
				return
			}
			m.Set(k, e.NewDeep().Elem())
			return
		})
	case Pointer:
		e := v.Elem()
		if e.ptr != nil {
			*(*unsafe.Pointer)(n.ptr) = e.NewDeep().ptr
		}
	case Slice:
		l := (*sliceHeader)(v.ptr).Len
		t := (*sliceType)(unsafe.Pointer(v.typ)).elem
		*(*unsafe.Pointer)(&n.ptr) = unsafe.Pointer(&sliceHeader{unsafe_NewArray(t, l), l, l})
		if !t.Kind().IsBasic() {
			s := (SLICE)(n.Elem())
			(SLICE)(v).ForEach(func(i int, k string, e VALUE) (brake bool) {
				set(e, offset((*sliceHeader)(s.ptr).Data, uintptr(i)*t.size), t)
				return
			})
		}
	case Struct:
		s := (STRUCT)(n.Elem())
		f := (*structType)(unsafe.Pointer(v.typ)).fields
		(STRUCT)(v).ForEach(func(i int, k string, e VALUE) (brake bool) {
			set(e, offset(s.ptr, f[i].offset), f[i].typ)
			return
		})
	}
	return n
}

// ------------------------------------------------------------ /
// GOLANG STANDARD IMPLEMENTATIONS
// implementations of functions natively available for
// interface and reflect.Value in golang
// referenced packages: reflect
// ------------------------------------------------------------ /

// Kind returns the kind of data type of Value
func (v VALUE) KIND() KIND {
	return v.typ.KIND()
}

// Kind returns the kind of data type of Value
func (v VALUE) Kind() KIND {
	return v.typ.Kind()
}

// Interface returns VALUE as interface {}
func (v VALUE) Interface() any {
	return v.Reflect().Interface()
}

// UnsafePointer returns the an unsafe.Pointer
// to the underlying Value
func (v VALUE) Pointer() unsafe.Pointer {
	if v.flag&flagIndir != 0 {
		return *(*unsafe.Pointer)(v.ptr)
	}
	return v.ptr
}

// PointerTo returns a pointer to the underlying value,
// if the value is a string, returns *string as VALUE
func (v VALUE) PointerTo() VALUE {
	return VALUE{
		v.typ.ptrType(),
		unsafe.Pointer(&v.ptr),
		flag(Pointer),
	}
}

// Uintptr returns the uintptr representation of
// a pointer to the underlying Value
func (v VALUE) Uintptr() uintptr {
	return uintptr(v.Pointer())
}

// Len returns the number of items in VALUE
// panics if Value is not of a struct, array, map, slice, string
func (v VALUE) Len() int {
	switch v.Kind() {
	case Array:
		return (ARRAY)(v).Len()
	case Map:
		return (MAP)(v).Len()
	case Slice:
		return (SLICE)(v).Len()
	case String:
		return (*stringHeader)(v.ptr).Len
	case Struct:
		return (STRUCT)(v).Len()
	}
	panic("cannot call Len on type " + v.Kind().String())
}

// Elem returns the underlying value of a pointer
func (v VALUE) Elem() VALUE {
	return FromReflect(v.Reflect().Elem())
}

// ElemDeep cascades a series of pointers to return the underlying VALUE
func (v VALUE) ElemDeep() VALUE {
	for v.Kind() == Pointer {
		v = v.Elem()
	}
	return v
}

// Index returns the value found at index i of VALUE
// returns VALUE if not an Array, Map, Pointer, Slice, String or Struct
// does not panic if i is greater than len of VALUE
func (v VALUE) Index(i int) VALUE {
	k := v.Kind()
	switch k {
	case Array:
		return (ARRAY)(v).Index(i)
	case Map:
		return (MAP)(v).Index(strconv.FormatInt(int64(i), 10))
	case Pointer:
		return v.Elem().Index(i)
	case Slice:
		return (SLICE)(v).Index(i)
	case String:
		h := stringHeader{
			Data: (*stringHeader)(v.ptr).Data + uintptr(i),
			Len:  1,
		}
		v.ptr = unsafe.Pointer(&h)
	case Struct:
		return (STRUCT)(v).Index(i)
	}
	return v
}

// MapIndex returns the value found at index (key) i in map,
// panics if VALUE is not a map
func (v VALUE) MapIndex(i string) VALUE {
	if v.Kind() == Map {
		return (MAP)(v).Index(i)
	}
	panic("value is not a map")
}

// StructField returns the field with name f in struct,
// panics if VALUE is not a struct
func (v VALUE) StructField(f string) FIELD {
	if v.Kind() == Struct {
		return (STRUCT)(v).Field(f)
	}
	panic("value is not a struct")
}

// ForEach executes function f on each item in VALUE,
// where VALUE is an Array, Map, Slice, String, or Struct,
// k is "" for Array, Slice and String, i is not fixed for Maps
func (v VALUE) ForEach(f func(i int, k string, v VALUE) (brake bool)) {
	switch v.Kind() {
	case Array:
		(ARRAY)(v).ForEach(f)
	case Map:
		(MAP)(v).ForEach(f)
	case Slice:
		(SLICE)(v).ForEach(f)
	case String:
		for i := 0; i < (*stringHeader)(v.ptr).Len; i++ {
			h := stringHeader{
				Data: (*stringHeader)(v.ptr).Data + uintptr(i),
				Len:  1,
			}
			if f(i, "", VALUE{v.typ, unsafe.Pointer(&h), v.flag}) {
				break
			}
		}
	case Struct:
		(STRUCT)(v).ForEach(f)
	}
}

// SetIndex sets the value of a key (or index) for a
// slice, array, map, struct or string
func (v VALUE) SetIndex(key any, val any) VALUE {
	k := ValueOfV(key)
	switch v.Kind() {
	case Array:
		(ARRAY)(v).Index(k.Int()).Set(val)
	case Interface:
		v = v.SetType()
		if v.Kind() != Interface {
			v.SetIndex(key, val)
		} else {
			panic("cannot set index of untyped interface")
		}
	case Pointer:
		v.ElemDeep().SetIndex(key, val)
	case Map:
		(MAP)(v).Set(k.String(), val)
	case Slice:
		return (SLICE)(v).Set(k.Int(), val).VALUE()
	case String:
		v.Set(v.STRING().SetIndex(k.Int(), ValueOf(val).String()))
	case Struct:
		switch k.Kind() {
		case Int:
			(STRUCT)(v).Index(k.Int()).Set(val)
		case String:
			(STRUCT)(v).Field(k.String()).Set(val)
		default:
			panic("invalid struct key")
		}
	default:
		panic("cannot set index on value type")
	}
	return v
}

// Set updates the VALUE to a and returns VALUE
func (v VALUE) Set(a any) VALUE {
	v = v.SetType()
	n := ValueOfV(a).SetType()
	switch {
	case v.typ == n.typ:
		return v.setMatched(n)
	case v.Kind() == Pointer && v.typ.elem() == n.typ:
		v.Elem().setMatched(n)
		return v
	default:
		return v.setUnmatched(n)
	}
}

func (v VALUE) setMatched(n VALUE) VALUE {
	switch v.Kind() {
	case Bool, Int8, Uint8:
		*(*[1]byte)(v.ptr) = *(*[1]byte)(n.ptr)
	case Int16, Uint16:
		*(*[2]byte)(v.ptr) = *(*[2]byte)(n.ptr)
	case Int32, Uint32, Float32:
		*(*[4]byte)(v.ptr) = *(*[4]byte)(n.ptr)
	case Int, Int64, Uint, Uint64, Uintptr, Float64, Complex64:
		*(*[8]byte)(v.ptr) = *(*[8]byte)(n.ptr)
	case Interface:
		*(*any)(v.ptr) = *(*any)(n.ptr)
	case Map: // hmap header size
		*(*[48]byte)(v.Pointer()) = *(*[48]byte)(n.Pointer())
	case Pointer:
		*(*unsafe.Pointer)(&v.ptr) = n.ptr
	case Slice: // slice header size
		*(*[24]byte)(v.ptr) = *(*[24]byte)(n.ptr)
	case String: // string header size
		*(*[16]byte)(v.ptr) = *(*[16]byte)(n.ptr)
	default:
		typedmemmove(v.typ, v.ptr, n.ptr)
	}
	return v
}

func (v VALUE) setUnmatched(n VALUE) VALUE {
	switch v.KIND() {
	case Bool:
		*(*bool)(v.ptr) = n.Bool()
	case Int:
		*(*int)(v.ptr) = n.Int()
	case Int8:
		*(*int8)(v.ptr) = n.INT().Int8()
	case Int16:
		*(*int16)(v.ptr) = n.INT().Int16()
	case Int32:
		*(*int32)(v.ptr) = n.INT().Int32()
	case Int64:
		*(*int64)(v.ptr) = n.INT().Int64()
	case Uint:
		*(*uint)(v.ptr) = n.Uint()
	case Uint8:
		*(*uint8)(v.ptr) = n.UINT().Uint8()
	case Uint16:
		*(*uint16)(v.ptr) = n.UINT().Uint16()
	case Uint32:
		*(*uint32)(v.ptr) = n.UINT().Uint32()
	case Uint64:
		*(*uint64)(v.ptr) = n.UINT().Uint64()
	case Float32:
		*(*float32)(v.ptr) = n.FLOAT().Float32()
	case Float64:
		*(*float64)(v.ptr) = n.Float64()
	case Array, Map, Slice:
		// must have identical type match
		panic("type mismatch on set value")
	case Pointer:
		v, n = v.ElemDeep(), n.ElemDeep()
		if v.typ == n.typ {
			v.setMatched(n)
		} else {
			v.setUnmatched(n)
		}
	case Interface:
		*(*any)(v.ptr) = n.Interface()
	case String:
		*(*string)(v.ptr) = n.String()
	case Struct:
		if matchStructType(v.typ, n.typ) {
			typedmemmove(v.typ, v.ptr, n.ptr)
		}
	case Time:
		*(*TIME)(v.ptr) = n.TIME()
	case Uuid:
		*(*UUID)(v.ptr) = n.UUID()
	case Bytes:
		*(*[]byte)(v.ptr) = n.Bytes()
	default:
		panic("type mismatch on set value")
	}
	return v
}

// Append adds the provided a to the end of a Slice,
// panics if VALUE is not a Slice
func (v VALUE) Append(a ...any) VALUE {
	switch v.Kind() {
	case Array:
		return (ARRAY)(v).SLICE().Append(a...).ARRAY().VALUE()
	case Slice:
		return (SLICE)(v).Append(a...).VALUE()
	case Pointer:
		return v.ElemDeep().Append(a...)
	}
	panic("can only append to slice value")
}

// Extend adds n elements to a Slice,
// panics if VALUE is not a Slice
func (v VALUE) Extend(n int) VALUE {
	switch v.Kind() {
	case Array:
		return (ARRAY)(v).SLICE().Extend(n).ARRAY().VALUE()
	case Slice:
		return (SLICE)(v).Extend(n).VALUE()
	}
	panic("can only extend slice value")
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

func (v VALUE) Cast(k KIND) any {
	vk := v.typ.KIND()
	if vk == k {
		return v.Interface()
	}
	if !vk.IsBasic() || !k.IsBasic() {
		panic("cannot convert to type")
	}
	switch k {
	case Bool:
		return v.Bool()
	case Int:
		return v.Int()
	case Int8:
		return v.INT().Int8()
	case Int16:
		return v.INT().Int16()
	case Int32:
		return v.INT().Int32()
	case Int64:
		return v.INT().Int64()
	case Uint:
		return v.Uint()
	case Uint8:
		return v.UINT().Uint8()
	case Uint16:
		return v.UINT().Uint16()
	case Uint32:
		return v.UINT().Uint32()
	case Uint64:
		return v.UINT().Int64()
	case Float32:
		return v.FLOAT().Float32()
	case Float64:
		return v.Float64()
	case String:
		return v.String()
	case Time:
		return v.TIME()
	case Uuid:
		return v.UUID()
	case Bytes:
		return v.Bytes()
	default:
		panic("cannot convert to type")
	}
}

func (v VALUE) Convert(a any) any {
	return v.convert(getrtype(a)).Interface()
}

func (v VALUE) convert(typ *rtype) VALUE {
	v = v.SetType()
	vk, tk := v.typ.KIND(), typ.Kind()
	switch {
	case v.typ == typ:
		return v
	case vk == Pointer && v.typ.elem() == typ:
		return v.Elem()
	case tk == Pointer && typ.elem() == v.typ:
		return VALUE{
			typ,
			unsafe.Pointer(&v.ptr),
			v.flag&(flagIndir|flagAddr) | flag(Pointer),
		}
	case vk.IsBasic() && tk.IsBasic():
		return ValueOf(v.Cast(tk))
	default:
		return typ.New().Elem().Set(v)
	}
}

func (v VALUE) Serialize() string {
	if v.ptr == nil {
		return "null"
	}
	switch v.KIND() {
	default:
		return v.String()
	case String:
		return (*STRING)(v.ptr).Serialize()
	case Bytes:
		return (*BYTES)(v.ptr).STRING().Serialize()
	case Array:
		return (ARRAY)(v).Serialize(ancestor{v.typ, v.Uintptr()})
	case Map:
		return (MAP)(v).Serialize(ancestor{v.typ, v.Uintptr()})
	case Pointer:
		p := v.Elem()
		if p.ptr == nil {
			return "null"
		}
		return p.SetType().Serialize()
	case Slice:
		return (SLICE)(v).Serialize(ancestor{v.typ, v.Uintptr()})
	case Struct:
		return (STRUCT)(v).Serialize(ancestor{v.typ, v.Uintptr()})
	case Time:
		return (*TIME)(v.ptr).Serialize()
	case Uuid:
		return (*UUID)(v.ptr).Serialize()
	}
}

type ancestor struct {
	typ     *rtype
	pointer uintptr
}

func (v VALUE) serialSafe(ancestry ...ancestor) (s string, recursive bool) {
	if v.ptr == nil {
		return "null", false
	}
	k := v.KIND()
	if k.IsBasic() {
		return v.Serialize(), false
	}
	if k != Array && k != Struct && v.Pointer() == nil {
		return "null", false
	}
	uptr := v.Uintptr()
	for _, a := range ancestry {
		if uptr == a.pointer && v.typ == a.typ {
			return "", true
		}
	}
	ancestry = append(ancestry, ancestor{v.typ, uptr})
	switch k {
	case Array:
		return (ARRAY)(v).Serialize(ancestry...), false
	case Map:
		return (MAP)(v).Serialize(ancestry...), false
	case Pointer:
		return v.Elem().SetType().serialSafe(ancestry...)
	case Slice:
		return (SLICE)(v).Serialize(ancestry...), false
	case Struct:
		if v.typ == getrtype(VALUE{}) {
			return v.Interface().(VALUE).serialSafe(ancestry...)
		}
		return (STRUCT)(v).Serialize(ancestry...), false
	default:
		return v.Serialize(), false
	}
}

// ------------------------------------------------------------ /
// EXPANDED FUNCTIONS
// implementations of new functions for
// interface and reflect.Value
// referenced packages: reflect
// ------------------------------------------------------------ /

// IsPtr evaluates whether Value is a pointer to a value
func (v VALUE) IsPtr() bool {
	return v.Kind() == Pointer
}

// TypeMatch evaluates whether the type of a is the same type of Value
func (v VALUE) TypeMatch(a any) bool {
	return v.typ == ValueOf(a).typ
}

// PtrKind returns the Kind of the underlying Value
func (v VALUE) PtrKind() KIND {
	return v.ElemDeep().KIND()
}

// ElemKind returns the KIND of the underlying Value(s)
func (v VALUE) ElemKind() KIND {
	return v.elemType().KIND()
}

// elemType returns the *rtype of the underlying Value(s)
func (v VALUE) elemType() *rtype {
	switch v.Kind() {
	case Array:
		return (*arrayType)(unsafe.Pointer(v.typ)).elem
	case Interface:
		return v.SetType().typ
	case Map:
		return (*mapType)(unsafe.Pointer(v.typ)).elem
	case Pointer:
		return (*ptrType)(unsafe.Pointer(v.typ)).elem
	case Slice:
		return (*sliceType)(unsafe.Pointer(v.typ)).elem
	case Struct:
		return (*structType)(unsafe.Pointer(v.typ)).fields[0].typ
	default:
		return v.typ
	}
}
