// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	"reflect"
	"sort"
	"unsafe"
)

// ------------------------------------------------------------ /
// GOTYPE CUSTOM TYPE IMPLEMENTATION
// implementation of custom type of slice
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

type SLICE VALUE

// SliceOf returns a as gotype SLICE
// panics if a is not convertable to slice
func SliceOf(a any) SLICE {
	if a, is := a.(SLICE); is {
		return a
	}
	return ValueOf(a).SLICE()
}

// SLICE returns VALUE as gotype SLICE
func (v VALUE) SLICE() SLICE {
	switch v.Kind() {
	case Slice:
		return (SLICE)(v)
	case Pointer:
		return v.ElemDeep().SLICE()
	default:
		switch v.KIND() {
		case Array:
			return (ARRAY)(v).SLICE()
		case Map:
			return (MAP)(v).SLICE()
		case Struct:
			return (STRUCT)(v).SLICE()
		case Bytes:
			j := (any)(JSON{})
			if v.typ == (*VALUE)(unsafe.Pointer(&j)).typ {
				return (*JSON)(v.ptr).SLICE()
			}
			return (SLICE)(v)
		default:
			panic("cannot convert value to slice")
		}
	}
}

// ------------------------------------------------------------ /
// GOLANG STANDARD IMPLEMENTATIONS
// implementations of functions natively available for
// interface and reflect.Value in golang
// referenced packages: reflect
// ------------------------------------------------------------ /

// Len returns the number of items in SLICE
func (s SLICE) Len() int {
	return (*sliceHeader)(s.ptr).Len
}

// Index returns the value found at index i of SLICE
func (s SLICE) Index(i int) VALUE {
	if i >= s.Len() {
		panic("index is out of slice range")
	}
	return s.index(i)
}

func (s SLICE) index(i int) VALUE {
	t := (*sliceType)(unsafe.Pointer(s.typ))
	return VALUE{
		t.elem,
		unsafe.Pointer(uintptr((*sliceHeader)(s.ptr).Data) + uintptr(i)*t.elem.size),
		flagAddr | flagIndir | flag(t.elem.Kind()),
	}.SetType()
}

// ForEach executes function f on each item in SLICE,
// note: k equals "" at each item
func (s SLICE) ForEach(f func(i int, k string, v VALUE) (brake bool)) {
	for i := 0; i < s.Len(); i++ {
		if f(i, "", s.index(i)) {
			break
		}
	}
}

func (s SLICE) Extend(n int) SLICE {
	h := (*sliceHeader)(s.ptr)
	t := (*sliceType)(unsafe.Pointer(s.typ))
	if h.Cap < h.Len+n {
		*h = growslice(t.elem, *h, n)
	}
	if t.elem.Kind().CanNil() {
		np := t.elem.newPtr()
		for i := 0; i < n; i++ {
			*(*unsafe.Pointer)(unsafe.Pointer(uintptr(h.Data) + (uintptr(h.Len)+uintptr(i))*t.elem.size)) = np
		}
	}
	h.Len += n
	return s
}

// Append adds x values to the end of an existing SLICE
// panics x does not match slice type
func (s SLICE) Append(x ...any) SLICE {
	l := s.Len()
	s = s.Extend(len(x))
	for i, e := range x {
		s.index(l + i).Set(e)
	}
	return s
}

// Set updates the value at index i to value v,
// returns the SLICE with the updated value
func (s SLICE) Set(i int, v any) SLICE {
	if i >= s.Len() {
		s = s.Extend(i - s.Len() + 1)
	}
	s.index(i).Set(v)
	return s
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Native returns gotype SLICE as a golang any of slice
func (s SLICE) Native() any {
	return s.Interface()
}

// Interface returns gotype SLICE as a golang interface{}
func (s SLICE) Interface() any {
	return (VALUE)(s).Interface()
}

// VALUE returns gotype SLICE as gotype VALUE
func (s SLICE) VALUE() VALUE {
	return (VALUE)(s)
}

// Encode returns a gotype encoding of SLICE
func (s SLICE) Encode() ENCODING {
	l := s.Len()
	e := append([]byte{s.typ.Kind().Byte(), (*sliceType)(unsafe.Pointer(s.typ)).elem.Kind().Byte()}, lenBytes(l)...)
	for i := 0; i < l; i++ {
		e = append(e, s.Index(i).Encode()...)
	}
	return e
}

// Bytes returns gotype SLICE as JSON []byte
func (s SLICE) Bytes() []byte {
	return []byte(s.String())
}

// Bool returns gotype SLICE as bool
// false if empty, true if a len > 0
func (s SLICE) Bool() bool {
	return s.Len() > 0
}

// ARRAY reutrns gotype SLICE as gotype ARRAY
func (s SLICE) ARRAY() ARRAY {
	h := (*sliceHeader)(s.ptr)
	s.typ = reflectType(reflect.ArrayOf(h.Len, toType((*sliceType)(unsafe.Pointer(s.typ)).elem)))
	s.ptr = unsafe.Pointer(h.Data)
	return (ARRAY)(s)
}

// String returns gotype SLICE as a serialized string
func (s SLICE) String() string {
	return s.Serialize()
}

// Serialize returns gotype SLICE as a serialized string
func (s SLICE) Serialize(ancestry ...ancestor) (S string) {
	if s.ptr == nil || *(*unsafe.Pointer)(s.ptr) == nil {
		return "null"
	}
	if s.Len() == 0 {
		return "[]"
	}
	s.ForEach(func(i int, k string, v VALUE) (brake bool) {
		sval, recursive := v.serialSafe(ancestry...)
		if !recursive {
			S += "," + sval
		}
		return
	})
	return "[" + S[1:] + "]"
}

// Slice returns gotype SLICE as []any
func (s SLICE) Slice() []any {
	l := s.Len()
	r := make([]any, l)
	for i := 0; i < l; i++ {
		r[i] = s.index(i).Interface()
	}
	return r
}

// Values returns SLICE as []Value
// panics if cannot convert an element to Value
func (s SLICE) Values() []VALUE {
	l := s.Len()
	r := make([]VALUE, l)
	for i := 0; i < s.Len(); i++ {
		r[i] = s.index(i)
	}
	return r
}

// Ints returns SLICE as []int
// panics if cannot convert an element to int
func (s SLICE) Ints() []int {
	l := s.Len()
	r := make([]int, l)
	for i := 0; i < s.Len(); i++ {
		r[i] = s.index(i).Int()
	}
	return r
}

// Floats returns SLICE as []float
// panics if cannot convert an element to float
func (s SLICE) Floats() []float64 {
	l := s.Len()
	r := make([]float64, l)
	for i := 0; i < l; i++ {
		r[i] = s.index(i).Float64()
	}
	return r
}

// Strings returns SLICE as []string
// panics if cannot convert an element to string
func (s SLICE) Strings() []string {
	l := s.Len()
	r := make([]string, l)
	for i := 0; i < l; i++ {
		r[i] = s.index(i).String()
	}
	return r
}

// Map returns gotype SLICE as gotype Map
func (s SLICE) Map() map[string]any {
	m := map[string]any{}
	for i := 0; i < s.Len(); i++ {
		m[INT(i).String()] = s.Index(i).Interface()
	}
	return m
}

// MAP returns gotype SLICE as gotype Map
func (s SLICE) MAP() MAP {
	return (MAP)(ValueOf(s.Map()))
}

// MapKeys returns a Map of the values of SLICE
// mapped to keys in the provided SLICE
func (s SLICE) MapKeys(keys []string) (m map[string]any) {
	m = map[string]any{}
	for i := 0; i < s.Len(); i++ {
		m[keys[i]] = s.Index(i).Interface()
	}
	return
}

// MapKeys returns a Map of the values of SLICE
// mapped to keys in the provided SLICE
func (s SLICE) MapKeysMap(keys []string) MAP {
	return MapOf(s.MapKeys(keys))
}

// MapValues returns a Map of SLICE as keys
// mapped to values in the provided SLICE
func (s SLICE) MapValues(values SLICE) (m map[string]any) {
	m = map[string]any{}
	for i := 0; i < s.Len(); i++ {
		m[s.Index(i).String()] = values.Index(i).Interface()
	}
	return
}

// MapValues returns a Map of SLICE as keys
// mapped to values in the provided SLICE
func (s SLICE) MapValuesMap(values SLICE) MAP {
	return MapOf(s.MapValues(values))
}

// Scan reads the values of SLICE into the provided destination pointer,
// the number of elements in dest must be greater than or equal to
// the number of elements in SLICE, otherwise Scan will panic
func (s SLICE) Scan(dest any) {
	d := ValueOf(dest).Elem()
	for i := 0; i < s.Len() && i < d.Len(); i++ {
		d.Index(i).Set(s.index(i))
	}
}

/* func (s SLICE) ScanList(dest any) {
	d := SliceOf(dest)
	s.ForEach(func(i int, k string, v VALUE) (brake bool) {

		d = d.Append(v.Interface())
		return
	}
} */

// JSON returns gotype SLICE as gotype JSON
func (s SLICE) JSON() JSON {
	return JSON(s.Serialize())
}

// COMPLEX SORTING
// allows for sorting of slices by inherent values
// for example, number strings will be sorted as numbers
// varying types will be sorted in order of numbers, time and strings
// slices with values not convertable to these types will panic

type sort_value struct {
	Original any
	Float    float64
	Time     int
	Str      string
}

type value_sorter struct {
	values []sort_value
	by     func(v1, v2 *sort_value) bool
}

type by func(v1, v2 *sort_value) bool

func (s *value_sorter) Len() int {
	return len(s.values)
}

func (s *value_sorter) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

func (s *value_sorter) Less(i, j int) bool {
	return s.by(&s.values[i], &s.values[j])
}

func (by by) Sort(values []sort_value) {
	vs := &value_sorter{
		values: values,
		by:     by,
	}
	sort.Sort(vs)
}

// SortComplex sorts SLICE in progressive order of
// numbers, time and strings while keeping DataTypes consistent
// NOTE: this sort is resource intensive, use DataType specific
// sort for less costly sorting
func (A SLICE) SortComplex() []any {

	// build sort_value lists by data type
	l := A.Len()
	a, floats, times, strings := make([]any, l), make([]sort_value, l), make([]sort_value, l), make([]sort_value, l)
	var float_i, time_i, string_i int
	for i := 0; i < l; i++ {
		v := A.Index(i)
		k := v.Kind()
		switch {
		case k.IsNumeric():
			floats[float_i] = sort_value{Original: v, Float: v.Float64()}
			float_i++
		case k == Time:
			times[time_i] = sort_value{Original: v, Time: v.TIME().Int()}
			time_i++
		case k == String || k == Bytes:
			str := v.STRING()
			if f, can, _ := str.CanParseFloat(); can {
				floats[float_i] = sort_value{Original: v, Float: f}
				float_i++
			} else if t, can, _ := str.CanTime(); can {
				times[time_i] = sort_value{Original: v, Time: t.Int()}
				time_i++
			} else {
				strings[string_i] = sort_value{Original: v, Str: str.String()}
				string_i++
			}
		default:
			panic("slice must contain numeric, time or string values")
		}
	}

	// trim each data type list to remove excess values
	floats = floats[:float_i]
	times = times[:time_i]
	strings = strings[:string_i]

	// sort and populate final slice with each data type list
	//a = make(slice, float_i+time_i+string_i)
	if len(floats) > 0 {
		by(func(v1, v2 *sort_value) bool { return v1.Float < v2.Float }).Sort(floats)
		for i, v := range floats {
			a[i] = v.Original
		}
	}
	if len(times) > 0 {
		by(func(v1, v2 *sort_value) bool { return v1.Time < v2.Time }).Sort(times)
		for i, v := range times {
			a[float_i+i] = v.Original
		}
	}
	if len(strings) > 0 {
		by(func(v1, v2 *sort_value) bool { return v1.Str < v2.Str }).Sort(strings)
		for i, v := range strings {
			a[float_i+time_i+i] = v.Original
		}
	}

	return a
}

// SortStrs sorts SLICE by first converting all values
// to string and then returning a new []string in sort order
func (s SLICE) SortStrings() []string {
	n := s.Strings()
	sort.Strings(n)
	return n
}

// SortInts sorts SLICE by first converting all values
// to int and then returning a new []int in sort order
func (s SLICE) SortInts() []int {
	n := s.Ints()
	sort.Ints(n)
	return n
}

// SortFloats sorts SLICE by first converting all values
// to float and then returning a new []float64 in sort order
func (s SLICE) SortFloats() []float64 {
	n := s.Floats()
	sort.Float64s(n)
	return n
}
