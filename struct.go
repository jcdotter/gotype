// Copyright 2023 james dotter. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import (
	"reflect"
	"unsafe"
)

// ------------------------------------------------------------ /
// GOTYPE CUSTOM TYPE IMPLEMENTATION
// implementation of custom type of struct
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

type STRUCT VALUE

// StructOf returns a as gotype STRUCT
// panics if a is not convertable to struct
func StructOf(a any) STRUCT {
	if a, is := a.(STRUCT); is {
		return a
	}
	return ValueOf(a).STRUCT()
}

// STRUCT returns VALUE as gotype STRUCT
func (v VALUE) STRUCT() STRUCT {
	switch v.Kind() {
	case Struct:
		return (STRUCT)(v)
	case Map:
		return (MAP)(v).STRUCT()
	case Pointer:
		return v.ElemDeep().STRUCT()
	default:
		panic("cannot convert value to struct")
	}
}

// ------------------------------------------------------------ /
// GOLANG STANDARD IMPLEMENTATIONS
// implementations of functions natively available for
// interface and reflect.Value in golang
// referenced packages: reflect
// ------------------------------------------------------------ /

// Len returns the number of fields in STRUCT
func (s STRUCT) Len() int {
	return len((*structType)(unsafe.Pointer(s.typ)).fields)
}

// NumField returns the number of fields in STRUCT
func (s STRUCT) NumField() int {
	return len((*structType)(unsafe.Pointer(s.typ)).fields)
}

// Index returns the value found at index i of STRUCT
func (s STRUCT) Index(i int) VALUE {
	if i >= s.Len() {
		panic("index is out of struct range")
	}
	return s.index(i)
}

// Index returns the value found at index i of STRUCT
func (s STRUCT) index(i int) VALUE {
	fs := (*structType)(unsafe.Pointer(s.typ)).fields
	f := fs[i]
	return VALUE{
		f.typ,
		unsafe.Pointer(uintptr(s.ptr) + f.offset),
		s.flag&(flagStickyRO|flagIndir|flagAddr) | flag(f.typ.Kind()),
	}.SetType()
}

// Index returns FIELD with name n of STRUCT
func (s STRUCT) Field(n string) (r FIELD) {
	s.ForFields(false, func(i int, f FIELD) (brake bool) {
		if f.name = f.name_.name(); n == f.name {
			r = f
			r.rawtag = r.name_.tag()
			brake = true
		}
		return
	})
	if r.ptr == nil {
		panic("struct field does not exists")
	}
	return
}

// Index returns FIELD with tag t having value v in STRUCT
func (s STRUCT) Tag(t string, v string) (r FIELD) {
	s.ForFields(false, func(i int, f FIELD) (brake bool) {
		f.rawtag = f.name_.tag()
		tv := getTagValue(f.rawtag, t, `"`[0])
		if tv != "" && tv == v {
			r = f
			r.name = r.name_.name()
			brake = true
		}
		return
	})
	if r.ptr == nil {
		panic("struct tag value does not exists")
	}
	return
}

// ForEach executes function f on each item in STRUCT,
// note: f is the field name
func (s STRUCT) ForEach(f func(i int, f string, v VALUE) (brake bool)) {
	for i, fld := range (*structType)(unsafe.Pointer(s.typ)).fields {
		v := VALUE{
			fld.typ,
			unsafe.Pointer(uintptr(s.ptr) + fld.offset),
			s.flag&(flagStickyRO|flagIndir|flagAddr) | flag(fld.typ.Kind()),
		}.SetType()
		if f(i, fld.name.name(), v) {
			break
		}
	}
}

// ForFields executes function f on each field in STRUCT;
// note: f is the struct field; tag and name are populated when inclDetail is true
func (s STRUCT) ForFields(inclDetail bool, f func(i int, f FIELD) (brake bool)) {
	for i, fld := range (*structType)(unsafe.Pointer(s.typ)).fields {
		fo := FIELD{
			typ:   fld.typ,
			ptr:   unsafe.Pointer(uintptr(s.ptr) + fld.offset),
			f:     s.flag&(flagStickyRO|flagIndir|flagAddr) | flag(fld.typ.Kind()),
			name_: fld.name,
			index: i,
		}
		if inclDetail {
			fo.name, fo.rawtag = fo.name_.name(), fo.name_.tag()
		}
		if f(i, fo) {
			break
		}
	}
}

// Set updates the value at index i to value v,
// returns the StrUCT with the updated value
func (s STRUCT) Set(i int, v any) STRUCT {
	if i >= s.Len() {
		panic("index is out of struct range")
	}
	s.Index(i).Set(v)
	return s
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Natural returns underlyding value of gotype STRUCT as a golang interface{}
func (s STRUCT) Native() any {
	return s.Interface()
}

// Interface returns gotype STRUCT as a golang interface{}
func (s STRUCT) Interface() any {
	var i any
	iface := (*VALUE)(unsafe.Pointer(&i))
	iface.typ, iface.ptr = s.typ, s.ptr
	return i
}

// VALUE returns gotype SLICE as gotype VALUE
func (s STRUCT) VALUE() VALUE {
	return (VALUE)(s)
}

// TYPE returns the TYPE of gotype STRUCT
func (s STRUCT) TYPE() *TYPE {
	return s.typ
}

// Pointer returns the pointer to gotype STRUCT
func (s STRUCT) Pointer() unsafe.Pointer {
	return s.ptr
}

// Bytes returns gotype STRUCT as []byte
func (s STRUCT) Encode() ENCODING {
	l := s.Len()
	e := append([]byte{byte(Struct)}, lenBytes(l)...)
	for i := 0; i < l; i++ {
		e = append(e, s.Index(i).Encode()...)
	}
	return e
}

// Bytes returns gotype STRUCT as serialized json []byte
func (s STRUCT) Bytes() []byte {
	return (VALUE)(s).Marshal(JsonMarshaller).Bytes()
}

// Bool returns gotype STRUCT as bool
// true if struct has at least one field
func (s STRUCT) Bool() bool {
	return s.Len() > 0
}

// ARRAY returns the values of a gotype STRUCT as gotype ARRAY
func (s STRUCT) ARRAY() ARRAY {
	return s.SLICE().ARRAY()
}

// Map returns gotype STRUCT as map[string]any
func (s STRUCT) Map() map[string]any {
	m := map[string]any{}
	for _, f := range (*structType)(unsafe.Pointer(s.typ)).fields {
		var v any
		e := (*VALUE)(unsafe.Pointer(&v))
		e.typ, e.ptr = f.typ, unsafe.Pointer(uintptr(s.ptr)+f.offset)
		m[f.name.name()] = v
	}
	return m
}

// MAP returns gotype STRUCT as MAP
func (s STRUCT) MAP() MAP {
	return MapOf(s.Map())
}

// String returns gotype STRUCT as a serialized json string
func (s STRUCT) String() string {
	if _, ok := s.ReflectType().MethodByName("String"); ok {
		return s.ReflectValue().MethodByName("String").Call([]reflect.Value{})[0].String()
	} else {
		return (VALUE)(s).Marshal(JsonMarshaller).String()
	}
}

// Name returns the name of the struct type as a string
func (s STRUCT) Name() string {
	return s.typ.Name()
}

// json returns gotype STRUCT as a serialized json string
func (s STRUCT) json(ancestry ...ancestor) (S string) {
	if s.ptr == nil {
		return "null"
	}
	if s.Len() == 0 {
		return "{}"
	}
	s.ForEach(func(i int, k string, v VALUE) (brake bool) {
		sval, recursive := v.jsonSafe(ancestry...)
		if !recursive {
			S += `,"` + k + `":` + sval
		}
		return
	})
	return "{" + S[1:] + "}"
}

// jsonByTag returns gotype STRUCT as a serialized json string with the provided tag as keys
func (s STRUCT) jsonByTag(tag string, ancestry ...ancestor) (S string) {
	if s.ptr == nil {
		return "null"
	}
	if s.Len() == 0 {
		return "{}"
	}
	s.ForFields(false, func(i int, f FIELD) (brake bool) {
		sval, recursive := f.VALUE().jsonSafe(ancestry...)
		if !recursive {
			S += `,"` + f.Tag(tag) + `":` + sval
		}
		return
	})
	if len(S) < 2 {
		return "{}"
	}
	return "{" + S[1:] + "}"
}

// Slice returns gotype STRUCT field values as []any
func (s STRUCT) Slice() []any {
	sf := (*structType)(unsafe.Pointer(s.typ)).fields
	l := len(sf)
	a := make([]any, l)
	for i := 0; i < l; i++ {
		var v any
		e := (*VALUE)(unsafe.Pointer(&v))
		e.typ, e.ptr = sf[i].typ, unsafe.Pointer(uintptr(s.ptr)+sf[i].offset)
		a[i] = v
	}
	return a
}

// SLICE returns gotype STRUCT as SLICE
func (s STRUCT) SLICE() SLICE {
	return SliceOf(s.Slice())
}

// Strings returns gotype STRUCT field values as []string
func (s STRUCT) Strings() []string {
	sf := (*structType)(unsafe.Pointer(s.typ)).fields
	l := len(sf)
	a := make([]string, l)
	for i := 0; i < l; i++ {
		var v VALUE
		e := (*VALUE)(unsafe.Pointer(&v))
		e.typ, e.ptr = sf[i].typ, unsafe.Pointer(uintptr(s.ptr)+sf[i].offset)
		a[i] = v.String()
	}
	return a
}

// Struct returns gotype STRUCT as interface()
func (s STRUCT) Struct() any {
	return s.Interface()
}

// JSON returns gotype STRUCT as gotype JSON
func (s STRUCT) JSON() JSON {
	return (VALUE)(s).Marshal(JsonMarshaller).Bytes()
}

// JsonByTag returns gotype STRUCT as gotype JSON using provided tag as keys
func (s STRUCT) JsonByTag(tag string) JSON {
	return JSON(s.jsonByTag(tag))
}

// ------------------------------------------------------------ /
// GOTYPE EXPANDED FUNCTIONS
// implementations of new functions for
// structs in gotype
// referenced packages: reflect
// ------------------------------------------------------------ /

// FieldNames returns a []string of names of fields in struct
func (s STRUCT) FieldNames() []string {
	names := make([]string, s.Len())
	for i, f := range (*structType)(unsafe.Pointer(s.typ)).fields {
		names[i] = f.name.name()
	}
	return names
}

// FieldIndex returns an index of field names in the Struct to Fields
func (s *STRUCT) FieldIndex() (index map[string]FIELD) {
	index = map[string]FIELD{}
	s.ForFields(true, func(i int, f FIELD) (brake bool) {
		index[f.name] = f
		return
	})
	return index
}

// TagIndex returns an index of tag values to Field in the Struct
// returns field names as keys if tag value is not unique across fields
func (s *STRUCT) TagIndex(tag string) map[string]FIELD {
	hasTag, index, fIndex := true, map[string]FIELD{}, map[string]FIELD{}
	s.ForFields(true, func(i int, f FIELD) (brake bool) {
		fIndex[f.name] = f
		if hasTag {
			tval := getTagValue(f.rawtag, tag, `"`[0])
			if _, found := index[tval]; tval == "" || found {
				hasTag = false
			} else {
				index[tval] = f
			}
		}
		return
	})
	if hasTag {
		return index
	}
	return fIndex
}

// SubTagIndex returns an index of subtag values to Field in the Struct
// returns field names as keys if subtag value is not unique across fields
func (s *STRUCT) SubTagIndex(tag string, subTag string) map[string]FIELD {
	hasTag, index, fIndex := true, map[string]FIELD{}, map[string]FIELD{}
	s.ForFields(true, func(i int, f FIELD) (brake bool) {
		index[f.name] = f
		if hasTag {
			tval := getTagValue(f.rawtag, tag, `"`[0])
			stval := getTagValue(tval, subTag, `'`[0])
			if _, found := index[stval]; stval == "" || found {
				hasTag = false
			} else {
				index[stval] = f
			}
		}
		return
	})
	if hasTag {
		return index
	}
	return fIndex
}

// ReflectValue returns the reflect.Value of gotype STRUCT
func (s STRUCT) ReflectValue() reflect.Value {
	return *(*reflect.Value)(unsafe.Pointer(&VALUE{s.typ, s.ptr, flag(Struct)}))
}

// ReflectType returns the reflect.Type of gotype STRUCT
func (s STRUCT) ReflectType() reflect.Type {
	return toType(s.typ)
}

func (s STRUCT) MapByTagMap(tag string) MAP {
	return MapOf(s.MapByTag(tag))
}

func (s STRUCT) MapByTag(tag string) (m map[string]any) {
	m = map[string]any{}
	for _, f := range (*structType)(unsafe.Pointer(s.typ)).fields {
		tagvalue := getTagValue(f.name.tag(), tag, `"`[0])
		if tagvalue != "" {
			var v any
			e := (*VALUE)(unsafe.Pointer(&v))
			e.typ, e.ptr = f.typ, unsafe.Pointer(uintptr(s.ptr)+f.offset)
			m[tagvalue] = v
		}
	}
	return
}

func (s STRUCT) MapFormatted(format StrFormat) (m map[string]any) {
	m = map[string]any{}
	for _, f := range (*structType)(unsafe.Pointer(s.typ)).fields {
		var v any
		e := (*VALUE)(unsafe.Pointer(&v))
		e.typ, e.ptr = f.typ, unsafe.Pointer(uintptr(s.ptr)+f.offset)
		m[string(format.Format(f.name.name()))] = v
	}
	return
}

// Scan reads the values of STRUCT into the provided Struct pointer dest
// by mapping the field names (or field tags) to those of the dest Struct
func (s STRUCT) Scan(dest any, tags ...string) {
	d := ValueOfV(dest).Elem().STRUCT()
	var dest_idx map[string]FIELD
	switch len(tags) {
	case 0:
		dest_idx = d.FieldIndex()
		s.ForFields(false, func(i int, f FIELD) (brake bool) {
			if dest_fld, found := dest_idx[f.name_.name()]; found {
				dest_fld.Set(f.VALUE())
			}
			return
		})
	case 1:
		tag := tags[0]
		dest_idx = d.TagIndex(tag)
		s.ForFields(true, func(i int, f FIELD) (brake bool) {
			if dest_fld, found := dest_idx[f.Tag(tag)]; found {
				dest_fld.Set(f.VALUE())
			}
			return
		})
	default:
		tag := tags[0]
		subtag := tags[1]
		dest_idx = d.SubTagIndex(tag, subtag)
		s.ForFields(true, func(i int, f FIELD) (brake bool) {
			if dest_fld, found := dest_idx[f.SubTag(tag, subtag)]; found {
				dest_fld.Set(f.VALUE())
			}
			return
		})
	}
}

// Gmap returns gotype STRUCT as gotype Gmap
func (s STRUCT) Gmap() Gmap {
	gm := make(Gmap, s.Len())
	s.ForEach(func(i int, k string, v VALUE) (brake bool) {
		gm[i] = GmapEl{k, v}
		return
	})
	return gm
}

// GmapByTag returns gotype STRUCT as gotype Gmap
func (s STRUCT) GmapByTag(tag string) Gmap {
	gm, gmn, hasTag := make(Gmap, s.Len()), make(Gmap, s.Len()), true
	s.ForFields(true, func(i int, f FIELD) (brake bool) {
		if hasTag {
			t := f.Tag(tag)
			if t == "" {
				hasTag = false
			} else {
				gm[i] = GmapEl{t, f.VALUE()}
			}
		}
		gmn[i] = GmapEl{f.name, f.VALUE()}
		return
	})
	if hasTag {
		return gm
	}
	return gmn
}
