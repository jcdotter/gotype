// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	"reflect"
	"unsafe"
)

// ------------------------------------------------------------ /
// GOTYPE CUSTOM TYPE IMPLEMENTATION
// implementation of custom type of map
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

type MAP VALUE

// MapOf returns a as gotype MAP
// panics if a is not convertable to map
func MapOf(a any) MAP {
	if a, is := a.(MAP); is {
		return a
	}
	return ValueOf(a).MAP()
}

// MAP returns VALUE as gotype MAP
func (v VALUE) MAP() MAP {
	switch v.Kind() {
	case Map:
		return (MAP)(v)
	case Pointer:
		return v.ElemDeep().MAP()
	default:
		switch v.KIND() {
		case Array:
			return (ARRAY)(v).MAP()
		case Slice:
			return (SLICE)(v).MAP()
		case Struct:
			return (STRUCT)(v).MAP()
		case Bytes:
			j := (any)(JSON{})
			if v.typ == (*VALUE)(unsafe.Pointer(&j)).typ {
				return (*JSON)(v.ptr).MAP()
			} else {
				return (SLICE)(v).MAP()
			}
		default:
			panic("cannot convert value to map")
		}
	}
}

// ------------------------------------------------------------ /
// GOLANG STANDARD IMPLEMENTATIONS
// implementations of functions natively available for
// interface and reflect.Value in golang
// referenced packages: reflect
// ------------------------------------------------------------ /

// Len returns the number of items in MAP
func (m MAP) Len() int {
	if p := *(*unsafe.Pointer)(m.ptr); p != nil {
		return *(*int)(p)
	}
	return 0
}

// Keys returns gotype MAP keys as []string
func (m MAP) Keys() []string {
	p := *(*unsafe.Pointer)(m.ptr)
	l := *(*int)(p)
	d := mallocgc(uintptr(l*16), getrtype(byte(0)), false)
	keys := *(*[]string)(unsafe.Pointer(&sliceHeader{Data: d, Len: l, Cap: l}))
	ob := uintptr(0)
	b := *(*uintptr)(unsafe.Pointer(uintptr(p) + 16))
	keysize := uintptr((*mapType)(unsafe.Pointer(m.typ)).keysize)
	bucketsize := uintptr((*mapType)(unsafe.Pointer(m.typ)).bucketsize)
	i := 0
	for i < l {
		k := 0
		for k < bucketCnt { // capture each key in bucket
			up := b + dataOffset + uintptr(k)*keysize
			if *(*unsafe.Pointer)(unsafe.Pointer(up)) == nil { // key deleted - skip
				k++
				continue
			}
			if *(*byte)(unsafe.Pointer(up + 8)) == 0 { // key empty - next bucket
				break
			}
			keys[i] = *(*string)(unsafe.Pointer(up))
			k++
			i++
		}
		// determine next bucket
		if sub := *(*uintptr)(unsafe.Pointer(b + bucketsize - 8)); sub != 0 { // has overflow
			if ob == 0 {
				ob = b
			}
			b = sub
		} else if ob != 0 { // in overflow bucket
			b = ob + bucketsize
			ob = 0
		} else {
			b += bucketsize
		}
	}
	return keys
}

// Index returns the value found at key k of Map
// returns nil pointer if key does not exist
func (m MAP) Index(k string) VALUE {
	v := VALUE{(*mapType)(unsafe.Pointer(m.typ)).elem, mapaccess_faststr(m.typ, *(*unsafe.Pointer)(m.ptr), k), m.flag}.SetType()
	if v.Kind() == Pointer && *(*unsafe.Pointer)(v.ptr) != nil {
		//v.ptr = *(*unsafe.Pointer)(v.ptr)
	}
	return v
}

// ForEach executes function f on each item in ARRAY,
// note: i is not a fixed value and may change across items
func (m MAP) ForEach(f func(i int, k string, v VALUE) (brake bool)) {
	p := *(*unsafe.Pointer)(m.ptr)
	if p == nil {
		return
	}
	t := (*mapType)(unsafe.Pointer(m.typ))
	l := *(*int)(p)
	b := uintptr(*(*unsafe.Pointer)(unsafe.Pointer(uintptr(p) + 16)))
	ob := uintptr(0)
	bucketsize := uintptr(t.bucketsize)
	keysize := uintptr(t.keysize)
	valuesize := uintptr(t.valuesize)
	voff := dataOffset + bucketCnt*keysize
	i := 0
	for i < l {
		k := 0
		for k < bucketCnt { // capture each key in bucket
			up := b + dataOffset + uintptr(k)*keysize
			if *(*unsafe.Pointer)(unsafe.Pointer(up)) == nil { // key deleted - skip
				k++
				continue
			}
			if *(*byte)(unsafe.Pointer(up + 8)) == 0 { // key empty - next bucket
				break
			}
			v := VALUE{t.elem, unsafe.Pointer(b + voff + uintptr(k)*valuesize), m.flag}.SetType()
			if v.Kind() == Pointer {
				//v.ptr = *(*unsafe.Pointer)(v.ptr)
			}
			if f(k, *(*string)(unsafe.Pointer(up)), v) {
				break
			}
			k++
		}
		i += k
		// determine next bucket
		if sub := *(*uintptr)(unsafe.Pointer(b + bucketsize - 8)); sub != 0 { // has overflow
			if ob == 0 {
				ob = b
			}
			b = sub
		} else if ob != 0 { // in overflow bucket
			b = ob + bucketsize
			ob = 0
		} else {
			b += bucketsize
		}
	}
}

// Set updates the value of key k to value v,
// returns the MAP with the updated value
func (m MAP) Set(k string, v any) MAP {
	if v == nil {
		mapdelete_faststr(m.typ, m.ptr, k)
		return m
	}
	typ := (*mapType)(unsafe.Pointer(m.typ)).elem
	knd := typ.elem().Kind()
	val := ValueOf(v)
	switch {
	case typ == val.typ:
		mapassign_faststr(m.typ, *(*unsafe.Pointer)(m.ptr), k, val.ptr)
	case typ.elem() == val.typ:
		mapassign_faststr(m.typ, *(*unsafe.Pointer)(m.ptr), k, unsafe.Pointer(&val.ptr))
	case typ.Kind() == Interface:
		if _, isVal := v.(VALUE); isVal {
			v = val.Interface()
		}
		mapassign_faststr(m.typ, *(*unsafe.Pointer)(m.ptr), k, unsafe.Pointer(&v))
	case knd == String:
		p := ValueOf(val.String()).ptr
		if typ.Kind() == Pointer {
			mapassign_faststr(m.typ, *(*unsafe.Pointer)(m.ptr), k, unsafe.Pointer(&p))
		} else {
			mapassign_faststr(m.typ, *(*unsafe.Pointer)(m.ptr), k, p)
		}
	case knd.IsBasic() && val.typ.Kind().IsBasic():
		p := ValueOf(val.retype(knd)).ptr
		if typ.Kind() == Pointer {
			mapassign_faststr(m.typ, *(*unsafe.Pointer)(m.ptr), k, unsafe.Pointer(&p))
		} else {
			mapassign_faststr(m.typ, *(*unsafe.Pointer)(m.ptr), k, p)
		}
	default:
		panic("cannot set mismatched data types")
	}
	return m
}

// Delete removes key from MAP
func (m MAP) Delete(key string) MAP {
	mapdelete_faststr(m.typ, *(*unsafe.Pointer)(m.ptr), key)
	return m
}

// ------------------------------------------------------------ /
// GOTYPE EXPANDED FUNCTIONS
// implementations of new functions for
// maps in gotype
// referenced packages: reflect
// ------------------------------------------------------------ /

// KeyPtr returns an unsafe pointer to the value at index 'key'
func (m MAP) KeyPtr(key string) unsafe.Pointer {
	return mapaccess_faststr(m.typ, *(*unsafe.Pointer)(m.ptr), key)
}

// Kind returns the kind of the map elements
func (m MAP) KIND() KIND {
	return (*mapType)(unsafe.Pointer(m.typ)).elem.Kind()
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Native returns gotype MAP as a golang any
func (m MAP) Native() any {
	return m.Interface()
}

// Interface returns gotype MAP as a golang interface{}
func (m MAP) Interface() any {
	var i any
	iface := (*VALUE)(unsafe.Pointer(&i))
	iface.typ, iface.ptr = m.typ, *(*unsafe.Pointer)(m.ptr)
	return i
}

// VALUE returns gotype MAP as gotype VALUE
func (m MAP) VALUE() VALUE {
	return (VALUE)(m)
}

// Bytes encodes gotype MAP as []byte
func (m MAP) Encode() ENCODING {
	t := (*mapType)(unsafe.Pointer(m.typ))
	e := append([]byte{
		byte(Map),
		t.key.Kind().Byte(),
		t.elem.Kind().Byte()},
		lenBytes(m.Len())...)
	m.ForEach(func(i int, k string, v VALUE) (brake bool) {
		e = append(e, STRING(k).Encode()...)
		e = append(e, m.Index(k).Encode()...)
		return
	})
	return e
}

// Bytes returns gotype MAP as serialized []byte
func (m MAP) Bytes() []byte {
	return []byte(m.String())
}

// Bool returns gotype MAP as bool
// false if empty, true if a len > 0
func (m MAP) Bool() bool {
	return m.Len() > 0
}

// ARRAY returns the values of a gotype MAP as gotype ARRAY
func (m MAP) ARRAY() ARRAY {
	return m.SLICE().ARRAY()
}

// Map returns gotype MAP as map[string]any
func (m MAP) Map() map[string]any {
	r := map[string]any{}
	m.ForEach(func(i int, k string, v VALUE) (brake bool) {
		r[k] = v.Interface()
		return
	})
	return r
}

// MapValues returns gotype MAP as a map[string]VALUE
func (m MAP) MapValues() map[string]VALUE {
	r := map[string]VALUE{}
	m.ForEach(func(i int, k string, v VALUE) (brake bool) {
		r[k] = v
		return
	})
	return r
}

// String returns gotype MAP as a serialized string
func (m MAP) String() string {
	return m.Serialize()
}

// Serialize returns gotype MAP as a serialized string
func (m MAP) Serialize(ancestry ...ancestor) (s string) {
	if m.ptr == nil || *(*unsafe.Pointer)(m.ptr) == nil {
		return "null"
	}
	if m.Len() == 0 {
		return "{}"
	}
	m.ForEach(func(i int, k string, v VALUE) (brake bool) {
		s += `,"` + k + `":` + v.serialSafe(ancestry...)
		return
	})
	return "{" + s[1:] + "}"
}

// STRING returns gotype MAP as a gotype STRING
func (m MAP) STRING() STRING {
	return STRING(m.String())
}

// Slice returns gotype MAP values as []any
func (m MAP) Slice() []any {
	vals := make([]any, m.Len())
	m.ForEach(func(i int, k string, v VALUE) (brake bool) {
		vals[i] = v.Interface()
		return
	})
	return vals
}

// Values returns gotype MAP values as []VALUE
func (m MAP) Values() []VALUE {
	vals := make([]VALUE, m.Len())
	m.ForEach(func(i int, k string, v VALUE) (brake bool) {
		vals[i] = v
		return
	})
	return vals
}

// Strings returns gotype MAP values as []string
func (m MAP) Strings() []string {
	vals := make([]string, m.Len())
	m.ForEach(func(i int, k string, v VALUE) (brake bool) {
		vals[i] = v.String()
		return
	})
	return vals
}

// SLICE returns gotype MAP values as gotype SLICE
func (m MAP) SLICE() SLICE {
	k := (*mapType)(unsafe.Pointer(m.typ)).elem.Kind()
	if k == Interface {
		return SliceOf(m.Slice())
	}
	s := k.NewSlice(m.Len())
	n := 0
	m.ForEach(func(i int, k string, v VALUE) (brake bool) {
		s.index(n).Set(v)
		n++
		return
	})
	return s
}

// Struct returns gotype MAP scanned into gotype Struct
func (m MAP) STRUCT() STRUCT {
	return m.StructTagged("")
}

// StructTagged returns gotype MAP scanned into gotype Struct
// with keys as tags with the tag lable provieded
func (m MAP) StructTagged(tag string) STRUCT {
	sfs := make([]reflect.StructField, m.Len())
	vals := map[string]reflect.Value{}
	m.ForEach(func(i int, k string, v VALUE) (brake bool) {
		n := STRING(k).ToPascal()
		rv := reflect.ValueOf(v.Interface())
		vals[n] = rv
		sf := reflect.StructField{
			Name: n,
			Type: rv.Type(),
		}
		if tag != "" {
			sf.Tag = reflect.StructTag(tag + `:"` + k + `"`)
		}
		sfs[i] = sf
		return
	})
	s := reflect.New(reflect.StructOf(sfs))
	for k, v := range vals {
		s.Elem().FieldByName(k).Set(v)
	}
	return StructOf(s.Interface())
}

// StructScan reads gotype MAP into the struct provided
// if tag is empty, Field names will be used to read Map keys into Struct
func (m MAP) StructScan(s STRUCT, tag string, subtag string) STRUCT {
	if tag == "" {
		for _, f := range (*structType)(unsafe.Pointer(s.typ)).fields {
			n := f.name.name()
			v := m.Index(n)
			if v.ptr != nil {
				VALUE{f.typ, unsafe.Pointer(uintptr(s.ptr) + f.offset), f.typ.flag()}.Set(v)
			}
		}
	} else if subtag == "" {
		for _, f := range (*structType)(unsafe.Pointer(s.typ)).fields {
			t := getTagValue(f.name.tag(), tag, `"`[0])
			v := m.Index(t)
			if v.ptr != nil {
				VALUE{f.typ, unsafe.Pointer(uintptr(s.ptr) + f.offset), f.typ.flag()}.Set(v)
			}
		}
	} else {
		for _, f := range (*structType)(unsafe.Pointer(s.typ)).fields {
			t := getTagValue(f.name.tag(), tag, `"`[0])
			st := getTagValue(t, subtag, `'`[0])
			v := m.Index(st)
			if v.ptr != nil {
				VALUE{f.typ, unsafe.Pointer(uintptr(s.ptr) + f.offset), f.typ.flag()}.Set(v)
			}
		}
	}
	return s
}

// JSON returns gotype MAP as gotype JSON
func (m MAP) JSON() JSON {
	return JSON(m.Serialize())
}
