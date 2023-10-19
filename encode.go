// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	b "bytes"
	"encoding/gob"
	"math"
	"unsafe"
)

// ------------------------------------------------------------ /
// GOTYPE BINARY ENCODING
// implementations of functions natively available for
// binary encoding in golang
// referenced packages: encoding/gob
// ------------------------------------------------------------ /

// ENCODING FORMAT
// []byte{Kind[, ElemKind[, KeyKind]][, Len], bytes...[, EndText[, bytes..., EndText, EndTrn]]}
//
// Elem Category	Kinds									Kind	ElemKind	KeyKind		Len		EndText		EndTrn
// Fixed Len		Bool, Int, Uint, Float, Time, Uuid		√
// Variable Len		String, Bytes							√										√
// Container		Slice, Array							√		√						√					√
//					Map										√		√			√			√					√
//					Struct									√								√					√
//
// Kind: 			byte 0 to 30; always first element in []byte
// EndText: 		0x3 (Control Character: End Of Text); denotes end of Variable Len elem
// EndTrn:  		0x4 (Control Character: End Of Transmission); denotes end of Container elem
//
// FORMAT EXAMPLES:
// Int32:			[]byte{Kind, byte, byte, byte, byte}
// string: 			[]byte{Kind, bytes..., EndText}
// slice: 			[]byte{Kind, ElemKind, Len, ElemKind, bytes...[, EndText], EndTrn}
//
// EXAMPLES:
// Type				Value									Encoding
// bool 			true									[]byte{0x1, 0x1}
// int32 			1										[]byte{0x2, 0x0, 0x0, 0x0, 0x1}
// string 			"true"									[]byte{0x24, 0x116, 0x114, 0x117, 0x101}
// slice			[]string{"true", "false"}				[]byte{0x23, 0x24, 0x116, 0x114, 0x117, 0x101, 0x3, 0x24, 0x102, 0x97, 0x108, 0x115, 0x101, 0x3, 0x4}
// uuid				267b3229-2566-4426-a826-8d80126e719a	[]byte{0x29, 0x26, 0x7b, 0x32, 0x29, 0x25, 0x66, 0x44, 0x26, 0xa8, 0x26, 0x8d, 0x80, 0x12, 0x6e, 0x71, 0x9a}
// map 				map[string]string{"one":"1"}			[]byte{0x21, 0x24, 0x111, 0x110, 0x101, 0x3, 0x49, 0x3, 0x4}
// struct			struct{true, "1", []uint8{1,2,3,4}}		[]byte{0x25, 0x1, 0x1, 0x3, 0x24, 0x49, 0x3, 0x23, 0x8, 0x1, 0x8, 0x2, 0x8, 0x3, 0x8, 0x4, 0x4}

// ENCODING contains a gotype byte encoding
// for storing to disk or for other future
// decoding to golang types
type ENCODING []byte

type decodex struct {
	k KIND  // data type kind of encoding
	b int   // number of bytes decoded
	v VALUE // gotype value of encoding
}

// String return serialized string of decoded value
func (d decodex) String() string {
	return d.v.String()
}

// Serialize returns serialized string decodex
func (d decodex) Serialize() string {
	return MapOf(map[string]any{
		"kind":  d.k.String(),
		"bytes": INT(d.b).String(),
		"value": d.v,
	}).String()
}

func (v VALUE) Encode() ENCODING {
	switch v.KIND() {
	case Bool:
		return (*BOOL)(v.ptr).Encode()
	case Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Float32, Float64:
		return v.EncodeNum()
	case Array:
		return (ARRAY)(v).Encode()
	case Interface:
		v = v.SetType()
		if v.Kind() != Interface {
			return v.Encode()
		}
	case Map:
		return (MAP)(v).Encode()
	case Pointer:
		return v.ElemDeep().Encode()
	case Slice:
		return (SLICE)(v).Encode()
	case String:
		return (*STRING)(v.ptr).Encode()
	case Struct:
		return (STRUCT)(v).Encode()
	case Time:
		return (*TIME)(v.ptr).Encode()
	case Uuid:
		return (*UUID)(v.ptr).Encode()
	case Bytes:
		return (*BYTES)(v.ptr).Encode()
	}
	panic("cannot convert to encode value")
}

// Kind returns the kind of the encoded value
// panics if cannot determine kind
func (e ENCODING) KIND() KIND {
	return KIND(e[0])
}

func (e ENCODING) String() string {
	return SliceOf(e).Serialize()
}

// LenAtLeast panics if length of encoding is less than l
func (e ENCODING) LenAtLeast(l int) {
	if len(e) < l {
		panic("number of bytes do not match required length")
	}
}

// Bytes returns encoding in []byte
func (e ENCODING) Bytes() []byte {
	return e
}

// GobEncode encodes any value to []byte using
// golang encoding/gob
func GobEncode(a any) []byte {
	buf := new(b.Buffer)
	gob.NewEncoder(buf).Encode(a)
	return buf.Bytes()
}

// GobDecode decodes []byte to a provided dest using
// golang encoding/gob; dest must be a pointer and match
// []byte origine type
func GobDecode(buf []byte, dest any) {
	gob.NewDecoder(b.NewBuffer(buf)).Decode(dest)
}

// Encode encodes any value to bytes
func Encode(a any) ENCODING {
	return ValueOfV(a).Encode()
}

// Decode decodes e to poitner dest and returns
// the number of bytes decoded
func Decode(e ENCODING, dest any) int {
	return e.Decode(dest)
}

// EncodeNum returns the byte encoding from a value of a number (eg.int)
func (v VALUE) EncodeNum() ENCODING {
	l := v.typ.size + 1
	b := make([]byte, l)
	b[0] = v.typ.Kind().Byte()
	for i := uintptr(1); i < l; i++ {
		b[i] = *(*byte)(offset(v.ptr, i-1))
	}
	return b
}

// EncodeNum returns the byte encoding from a value of a number (eg.int)
func EncodeNum(n any) ENCODING {
	return (*VALUE)(unsafe.Pointer(&n)).EncodeNum()
}

// EncodeNumCompressed returns the byte encoding from a value
// of a number compressed to the smallest size repreentation of that number
func EncodeNumCompressed(i any) ENCODING {
	v := *(*VALUE)(unsafe.Pointer(&i))
	switch v.Kind() {
	case Int, Int8, Int16, Int32, Int64:
		n := *(*int)(v.ptr)
		switch {
		case n <= math.MaxInt8:
			i = (int8)(n)
		case n <= math.MaxInt16:
			i = (int16)(n)
		case n <= math.MaxInt32:
			i = (int32)(n)
		}
	case Uint, Uint8, Uint16, Uint32, Uint64, Uintptr:
		n := *(*uint)(v.ptr)
		switch {
		case n <= math.MaxUint8:
			i = (uint8)(n)
		case n <= math.MaxUint16:
			i = (uint16)(n)
		case n <= math.MaxUint32:
			i = (uint32)(n)
		}
	case Float32, Float64:
		n := *(*float64)(v.ptr)
		if n <= math.MaxFloat32 {
			i = (float32)(n)
		}
	default:
		panic("cannot encode non numeric value")
	}
	return (*VALUE)(unsafe.Pointer(&i)).EncodeNum()
}

// Decode decodes encoding and inserts values into pointer dest
// panics if dest format does not match encoding
func (e ENCODING) Decode(dest any) int {
	enc := e.Decodex()
	dVal := destValue(dest)
	decodeValueSet(enc.v, dVal)
	return enc.b
}

// decodeValueSet cascades through decoded values eVal and
// inserts values into destination dVal
func decodeValueSet(eVal VALUE, dVal VALUE) {
	eKind, dKind := eVal.Kind(), dVal.Kind()
	// check if both are basic kinds
	if (eKind.IsBasic() && dKind.IsBasic()) || dKind == Interface {
		dVal.Set(eVal)
		return
	}
	if dKind == Pointer {
		decodeValueSet(eVal, dVal.Elem())
		return
	}
	switch eKind {
	case Slice, Array:
		ea := eVal.ARRAY()
		el := ea.Len()
		switch dKind {
		case Array:
			da := dVal.ARRAY()
			if el > da.Len() {
				panic("cannot decode to array of differing length")
			}
			for i := 0; i < el; i++ {
				decodeValueSet(ea.index(i), da.index(i))
			}
			return
		case Slice:
			da := dVal.SLICE()
			n := el - da.Len()
			if n > 0 {
				da.Extend(n)
			}
			for i := 0; i < el; i++ {
				decodeValueSet(ea.index(i), da.index(i))
			}
			return
		case Struct:
			ds := (STRUCT)(dVal)
			if el > ds.Len() {
				panic("cannot decode to struct of differing length")
			}
			for i := 0; i < el; i++ {
				decodeValueSet(ea.index(i), ds.index(i))
			}
			return
		}
	case Map:
		if dKind == Map {
			em, dm := eVal.MAP(), dVal.MAP()
			t := (*mapType)(unsafe.Pointer(dm.typ)).elem
			if dm.typ.Kind().IsBasic() {
				em.ForEach(func(i int, k string, v VALUE) (brake bool) {
					dm.Set(k, v)
					return
				})
				return
			}
			em.ForEach(func(i int, k string, v VALUE) (brake bool) {
				dv := dm.Index(k)
				if dv.ptr == nil {
					dv = t.NewValue().Elem()
					dm.Set(k, dv)
				}
				decodeValueSet(v, dv)
				return
			})
			return
		}
	}
	panic("dest format does not match encoding")
}

// destValue validates a pointer dest and
// return a Value of the underlying destination
func destValue(dest any) VALUE {
	v := ValueOfV(dest)
	if v.Kind() != Pointer {
		panic("dest must be a pointer")
	}
	return v.Elem().SetType()
}

// Decodex returns the Value and Kind and number of bytes of an encoding
// returns structs as []any of the struct values and panics if
// encoding cannot be decoded (or is corrupt)
func (e ENCODING) Decodex() decodex {
	switch e.KIND() {
	case Bool:
		return e.decodexBool()
	case Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Float32, Float64:
		return e.decodexNum()
	case Bytes, String:
		return e.decodexBytes()
	case Time:
		return e.decodexTime()
	case Uuid:
		return e.decodexUuid()
	case Array:
		return e.decodexArray()
	case Map:
		return e.decodexMap()
	case Slice:
		return e.decodexSlice()
	case Struct:
		return e.decodexStruct()
	}
	panic("cannot decode encoding of " + e.KIND().String())
}

// decodexBool decodes a bool element and returns
// the Value and the number of bytes processed
func (e ENCODING) decodexBool() (d decodex) {
	e.LenAtLeast(2)
	d.k = e.KIND()
	d.b = 2
	d.v = BYTES(e[1:2]).BOOL().VALUE()
	return
}

// decodexNum decodes a number element to the indicated number type and returns
// the Value of the number and the number of bytes processed
func (e ENCODING) decodexNum() (d decodex) {
	d.k = KIND(e[0])
	d.v = d.k.NewValue().Elem()
	d.b = int(d.v.typ.size) + 1
	e[1:].LenAtLeast(d.b - 1)
	for i := 1; i < d.b; i++ {
		*(*byte)(offseti(d.v.ptr, i-1)) = e[i]
	}
	return
}

// decodexBytes decodes a []byte element to BYTES or STRING and returns
// the Value and the number of bytes processed
func (e ENCODING) decodexBytes() (d decodex) {
	d.k = e.KIND()
	l, s := e.decodeLen(1)
	d.b = s + l
	b := e[s:d.b]
	if d.k == String {
		d.v = ValueOf(string(b))
		return
	}
	d.v = ValueOf(b)
	return
}

// decodexUuid decodes a uuid element to UUID and returns
// the Value of the UUID and the number of bytes processed
func (e ENCODING) decodexUuid() (d decodex) {
	d.k = e.KIND()
	d.b = 17
	e.LenAtLeast(d.b)
	b := BYTES(e[1:17])
	b.MustUuid()
	d.v = ValueOf(UUID(b))
	return
}

// decodexTime decodes a time element to TIME and returns
// the Value of the TIME and the number of bytes processed
func (e ENCODING) decodexTime() (d decodex) {
	e[0] = Int.Byte()
	d = e.decodexNum()
	d.k = e.KIND()
	d.v = d.v.TIME().VALUE()
	return
}

// decodexArray decodes an array to the array type and returns
// the Value of the slice and the number of bytes processed
func (e ENCODING) decodexArray() (d decodex) {
	var s ARRAY
	d.k = e.KIND()
	l, i, eKind := 0, 0, KIND(e[1])
	if !eKind.IsBasic() {
		eKind = Interface
	}
	l, d.b = e.decodeLen(2)
	s = eKind.NewArray(l)
	for i < l {
		n := e[d.b:].Decodex()
		if n.k != eKind && eKind != Interface {
			panic("corrupt encoding")
		}
		if i < l {
			s.index(i).Set(n.v)
		}
		d.b += n.b
		i++
	}
	d.v = s.VALUE()
	return
}

// decodexMap decodes a map to the map type and returns the
// Value of the map and the number of bytes processed
func (e ENCODING) decodexMap() (d decodex) {
	var m MAP
	d.k = e.KIND()
	l, i, kKind, eKind := 0, 0, KIND(e[1]), KIND(e[2])
	if kKind != String {
		panic("map must have a key type of string")
	}
	if !eKind.IsBasic() {
		eKind = Interface
	}
	m = eKind.NewMap()
	l, d.b = e.decodeLen(3)
	for i < l {
		// decode key
		k := e[d.b:].Decodex()
		if k.k != String {
			panic("corrupt encoding")
		}
		d.b += k.b
		// decode value
		v := e[d.b:].Decodex()
		if v.k != eKind && eKind != Interface {
			panic("corrupt encoding")
		}
		d.b += v.b
		m.Set(k.v.String(), v.v)
		i++
	}
	d.v = m.VALUE()
	return
}

// decodexSlice decodes a slice to the slice type and returns
// the Value of the slice and the number of bytes processed
func (e ENCODING) decodexSlice() (d decodex) {
	var s SLICE
	d.k = e.KIND()
	l, i, eKind := 0, 0, KIND(e[1])
	if !eKind.IsBasic() {
		eKind = Interface
	}
	l, d.b = e.decodeLen(2)
	s = eKind.NewSlice(l)
	for i < l {
		n := e[d.b:].Decodex()
		if n.k != eKind && eKind != Interface {
			panic("corrupt encoding")
		}
		s.index(i).Set(n.v)
		d.b += n.b
		i++
	}
	d.v = s.VALUE()
	return
}

// decodexStruct decodes a struct to []any and returns the
// Value of the Slice and the number of bytes processed
func (e ENCODING) decodexStruct() (d decodex) {
	d.k = e.KIND()
	l, i := 0, 0
	l, d.b = e.decodeLen(1)
	a := make([]any, l)
	s := SliceOf(&a)
	for i < l {
		f := e[d.b:].Decodex()
		if f.k == Interface {
			panic("corrupt encoding")
		}
		s.index(i).Set(f.v)
		d.b += f.b
		i++
	}
	d.v = s.VALUE()
	return
}

func (e ENCODING) decodeLen(offset uintptr) (len int, bytes int) {
	k := KIND(e[offset])
	offset++
	p := unsafe.Pointer(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&e))) + offset)
	switch k {
	case Uint8:
		return int(*(*uint8)(p)), int(offset + 1)
	case Uint16:
		return int(*(*uint16)(p)), int(offset + 2)
	case Uint32:
		return int(*(*uint32)(p)), int(offset + 4)
	default:
		return int(*(*uint)(p)), int(offset + 8)
	}
}

func lenBytes(l int) []byte {
	p := unsafe.Pointer(&l)
	switch {
	case l < math.MaxUint8:
		return []byte{
			Uint8.Byte(),
			*(*byte)(offset(p, 0)),
		}
	case l < math.MaxUint16:
		return []byte{
			Uint16.Byte(),
			*(*byte)(offset(p, 0)),
			*(*byte)(offset(p, 1)),
		}
	case l < math.MaxUint32:
		return []byte{
			Uint32.Byte(),
			*(*byte)(offset(p, 0)),
			*(*byte)(offset(p, 1)),
			*(*byte)(offset(p, 2)),
			*(*byte)(offset(p, 3)),
		}
	default:
		return []byte{
			Uint.Byte(),
			*(*byte)(offset(p, 0)),
			*(*byte)(offset(p, 1)),
			*(*byte)(offset(p, 2)),
			*(*byte)(offset(p, 3)),
			*(*byte)(offset(p, 4)),
			*(*byte)(offset(p, 5)),
			*(*byte)(offset(p, 6)),
			*(*byte)(offset(p, 7)),
		}
	}
}
