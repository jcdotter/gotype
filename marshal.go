// Copyright 2023 james dotter. All rights reserved.typVal
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

// ------------------------------------------------------------ /
// Marshaller
// a generic marshaller
// ------------------------------------------------------------ /

type Marshaller struct {
	Type      string
	Cursor    int
	CurIndent int
	Value     any
	Len       int
	Buffer,
	Space,
	Indent,
	Quote,
	Escape,
	Null,
	ValEnd,
	KeyEnd,
	BlockCommentStart,
	BlockCommentEnd,
	LineCommentStart,
	LineCommentEnd,
	SliceStart,
	SliceEnd,
	MapStart,
	MapEnd []byte
	UseIndent,
	QuotedKey,
	QuotedString,
	QuotedSpecial,
	QuotedNum,
	QuotedBool,
	QuotdeNull bool
}

// ------------------------------------------------------------ /
// Presets
// JSON, YAML...
// ------------------------------------------------------------ /

var (
	MarshallerJson = Marshaller{
		Type:              "json",
		QuotedKey:         true,
		QuotedString:      true,
		Space:             []byte(" \t\n\v\f\r"),
		Indent:            []byte("  "),
		Quote:             []byte(`"`),
		Escape:            []byte(`\`),
		ValEnd:            []byte(","),
		KeyEnd:            []byte(":"),
		BlockCommentStart: []byte("/*"),
		BlockCommentEnd:   []byte("*/"),
		LineCommentStart:  []byte("//"),
		LineCommentEnd:    []byte("\n"),
		SliceStart:        []byte("["),
		SliceEnd:          []byte("]"),
		MapStart:          []byte("{"),
		MapEnd:            []byte("}"),
	}
	MarshallerYaml = Marshaller{
		Type:             "yaml",
		UseIndent:        true,
		QuotedSpecial:    true,
		Space:            []byte(" \t\v\f\r"),
		Indent:           []byte("  "),
		Quote:            []byte(`"'`),
		Escape:           []byte(`\`),
		ValEnd:           []byte("\n"),
		KeyEnd:           []byte(":\n"),
		LineCommentStart: []byte("#"),
		LineCommentEnd:   []byte("\n"),
		SliceStart:       []byte("- "),
		SliceEnd:         []byte("]"),
		MapStart:         []byte("{"),
		MapEnd:           []byte("}"),
	}
	MarshallerInlineYaml = Marshaller{
		Type:             "yaml",
		QuotedKey:        true,
		QuotedSpecial:    true,
		Space:            []byte(" \t\v\f\r"),
		Indent:           []byte("  "),
		Quote:            []byte(`"'`),
		Escape:           []byte(`\`),
		ValEnd:           []byte(","),
		KeyEnd:           []byte(":"),
		LineCommentStart: []byte("#"),
		LineCommentEnd:   []byte("\n"),
		SliceStart:       []byte("["),
		SliceEnd:         []byte("]"),
		MapStart:         []byte("{"),
		MapEnd:           []byte("}"),
	}
)

// ------------------------------------------------------------ /
// Marshal Utilities
// methods for marshalling to type
// ------------------------------------------------------------ /

func (m *Marshaller) Marshal(a any) ([]byte, error) {
	return m.marshal(ValueOf(a))
}

func (m *Marshaller) marshal(v VALUE, ancestry ...ancestor) ([]byte, error) {
	v = v.SetType()
	if v.ptr == nil {
		return []byte("null"), nil
	}
	switch v.KIND() {
	case Bool:
		return m.MarshalBool(v.Bool())
	case Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Float32, Float64, Complex64, Complex128:
		return m.MarshalNum(v)
	case Array:
		return m.MarshalArray((ARRAY)(v), ancestry...)
	case Func:
		return m.MarshalFunc(v)
	case Interface:
		return m.MarshalInterface(v)
	case Map:
		return m.MarshalMap((MAP)(v), ancestry...)
	case Pointer:
		return m.marshal(v.Elem())
	case Slice:
		return m.MarshalSlice((SLICE)(v), ancestry...)
	case String:
		return m.MarshalString(v.String())
	case Struct:
		return m.MarshalStruct((STRUCT)(v), ancestry...)
	case UnsafePointer:
	case Field:
	case Time:
	case Uuid:
	case Bytes:
		return m.MarshalString(string(v.Bytes()))
	default:
		return nil, errors.New("cannot marshal type '" + v.typ.String() + "'")
	}
}

func (m *Marshaller) MarshalBool(b bool) (bytes []byte, err error) {
	if b {
		bytes = []byte("true")
	} else {
		bytes = []byte("false")
	}
	if m.QuotedBool {
		q := m.Quote[:1]
		bytes = append(q, append(bytes, q...)...)
	}
	m.ToBuffer(bytes)
	return
}

func (m *Marshaller) MarshalNum(v VALUE) (bytes []byte, err error) {
	switch v.Kind() {
	case Int:
		bytes = []byte(strconv.FormatInt(*(*int64)(v.ptr), 10))
	case Int8:
		bytes = []byte(strconv.FormatInt(int64(*(*int8)(v.ptr)), 10))
	case Int16:
		bytes = []byte(strconv.FormatInt(int64(*(*int16)(v.ptr)), 10))
	case Int32:
		bytes = []byte(strconv.FormatInt(int64(*(*int32)(v.ptr)), 10))
	case Int64:
		bytes = []byte(strconv.FormatInt(*(*int64)(v.ptr), 10))
	case Uint:
		bytes = []byte(strconv.FormatUint(*(*uint64)(v.ptr), 10))
	case Uint8:
		bytes = []byte(strconv.FormatUint(uint64(*(*uint8)(v.ptr)), 10))
	case Uint16:
		bytes = []byte(strconv.FormatUint(uint64(*(*uint16)(v.ptr)), 10))
	case Uint32:
		bytes = []byte(strconv.FormatUint(uint64(*(*uint32)(v.ptr)), 10))
	case Uint64:
		bytes = []byte(strconv.FormatUint(*(*uint64)(v.ptr), 10))
	case Uintptr:
		bytes = []byte(strconv.FormatUint(*(*uint64)(v.ptr), 10))
	case Float32:
		bytes = []byte(strconv.FormatFloat(float64(*(*float32)(v.ptr)), 'f', -1, 64))
	case Float64:
		bytes = []byte(strconv.FormatFloat(*(*float64)(v.ptr), 'f', -1, 64))
	case Complex64:
		bytes = []byte(strconv.FormatComplex(complex128(*(*complex64)(v.ptr)), 'f', -1, 128))
	case Complex128:
		bytes = []byte(strconv.FormatComplex(*(*complex128)(v.ptr), 'f', -1, 128))
	default:
		return nil, errors.New("cannot marshal type '" + v.typ.String() + "'")
	}
	if m.QuotedNum {
		q := m.Quote[:1]
		bytes = append(q, append(bytes, q...)...)
	}
	m.ToBuffer(bytes)
	return
}

func (m *Marshaller) MarshalArray(a ARRAY, ancestry ...ancestor) ([]byte, error) {
	if a.ptr == nil {
		m.ToBuffer(m.Null)
		return m.Buffer, nil
	}
	if a.Len() == 0 {
		m.ToBuffer(append(m.SliceStart, m.SliceEnd...))
		return m.Buffer, nil
	}
	m.MarshalSliceStart()
	a.ForEach(func(i int, k string, v VALUE) (brake bool) {
		m.CurIndent++
		m.marshal(v)
		m.CurIndent--
		return
	})
	m.MarshalSliceEnd()
	return m.Buffer, nil
}

func (m *Marshaller) MarshalFunc(v VALUE) (bytes []byte, err error) {
	bytes = []byte(v.typ.Name())
	m.ToBuffer(bytes)
	return
}

func (m *Marshaller) MarshalInterface(v VALUE) (bytes []byte, err error) {
	if v.ptr != nil {
		if *(*unsafe.Pointer)(v.ptr) != nil {
			v = v.SetType()
			if v.Kind() != Interface {
				return m.marshal(v)
			}
			return m.Marshal(fmt.Sprint(v.Interface()))
		}
	}
	return m.Null, nil
}

func (m *Marshaller) MarshalMap(hm MAP, ancestry ...ancestor) ([]byte, error) {
	if hm.ptr == nil {
		m.ToBuffer(m.Null)
		return m.Buffer, nil
	}
	if hm.Len() == 0 {
		m.ToBuffer(append(m.MapStart, m.MapEnd...))
		return m.Buffer, nil
	}
	m.MarshaMapStart()
	hm.ForEach(func(i int, k string, v VALUE) (brake bool) {
		m.MarshalKey([]byte(k))
		m.CurIndent++
		m.marshal(v)
		m.CurIndent--
		return
	})
	m.MarshaMapEnd()
	return m.Buffer, nil
}

func (m *Marshaller) MarshalSlice(s SLICE, ancestry ...ancestor) ([]byte, error) {
	if s.ptr == nil {
		m.ToBuffer(m.Null)
		return m.Buffer, nil
	}
	if s.Len() == 0 {
		m.ToBuffer(append(m.SliceStart, m.SliceEnd...))
		return m.Buffer, nil
	}
	m.MarshalSliceStart()
	s.ForEach(func(i int, k string, v VALUE) (brake bool) {
		m.CurIndent++
		m.marshal(v)
		m.CurIndent--
		return
	})
	m.MarshalSliceEnd()
	return m.Buffer, nil
}

func (m *Marshaller) MarshalString(s string) ([]byte, error) {
	b := []byte(s)
	quoted := m.QuotedString
	if !quoted && m.QuotedSpecial {
		if ContainsSpecial(s) {
			quoted = true
		}
	}
	if quoted {
		q := m.Quote[:1]
		b = append(append(q, BYTES(b).Escaped(q[0], m.Escape[0])...), q...)
	}
	m.ToBuffer(b)
	return m.Buffer, nil
}

func (m *Marshaller) MarshalStruct(s STRUCT, ancestry ...ancestor) ([]byte, error) {
	if s.ptr == nil {
		m.ToBuffer(m.Null)
		return m.Buffer, nil
	}
	if s.Len() == 0 {
		m.ToBuffer(append(m.MapStart, m.MapEnd...))
		return m.Buffer, nil
	}
	m.MarshaMapStart()
	fields := s.TagIndex(m.Type)
	for k, f := range fields {
		m.MarshalKey([]byte(k))
		m.CurIndent++
		m.marshal(f.VALUE())
		m.CurIndent--
	}
	m.MarshaMapEnd()
	return m.Buffer, nil
}

func (m *Marshaller) MarshalSliceStart() {
	m.ToBuffer(m.SliceStart)
	m.CurIndent++
	m.MarshalNext()
}

func (m *Marshaller) MarshalSliceEnd() {
	m.CurIndent--
	m.MarshalNext()
	m.ToBuffer(m.SliceStart)
}

func (m *Marshaller) MarshaMapStart() {
	m.ToBuffer(m.MapStart)
	m.CurIndent++
	m.MarshalNext()
}

func (m *Marshaller) MarshaMapEnd() {
	m.CurIndent--
	m.MarshalNext()
	m.ToBuffer(m.MapStart)
}

func (m *Marshaller) MarshalKey(k []byte) {
	if m.QuotedKey {
		q := m.Quote[:1]
		k = append(append(q, k...), q...)
	}
	if m.KeyEnd != nil {
		k = append(k, m.KeyEnd...)
	}
	m.Buffer = append(m.Buffer, k...)
}

func (m *Marshaller) MarshalNext() {
	if m.UseIndent && m.Indent != nil {
		m.ToBuffer([]byte("\n"))
		m.ToBuffer(bytes.Repeat(m.Indent, m.CurIndent))
	}
}

func (m *Marshaller) ToBuffer(b []byte) {
	m.Buffer = append(m.Buffer, b...)
}

func ContainsSpecial(s string) bool {
	for _, c := range s {
		if IsSpecialChar(byte(c)) {
			return true
		}
	}
	return false
}

func IsSpecialChar(b byte) bool {
	return b < 0x30 || (b > 0x3a && b < 0x41) || (b > 0x5a && b < 0x61) || b > 0x7a
}

// ------------------------------------------------------------ /
// Unmarshal Utilities
// methods for unmarshalling from type
// ------------------------------------------------------------ /

func (m *Marshaller) UnmarshalQuote() []byte {
	if m.IsQuote() {
		s := m.Cursor
		q := m.Buffer[s]
		for m.Cursor < m.Len {
			if m.IsEscape() {
				m.Inc(2)
				continue
			}
			if m.ByteIs(q) {
				m.Inc()
				break
			}

		}
		return m.Buffer[s:m.Cursor]
	}
	return nil
}

func (m *Marshaller) ParseWhitespace() []byte {
	if m.IsSpace() {
		s := m.Cursor
		for m.Cursor < m.Len {
			if !m.IsSpace() {
				break
			}
			m.Inc()
		}
		return m.Buffer[s:m.Cursor]
	}
	return nil
}

func (m *Marshaller) ByteIs(b byte) bool {
	return m.Byte() == b
}

func (m *Marshaller) Byte() byte {
	if m.Len == 0 {
		return 0
	}
	return m.Buffer[m.Cursor]
}

func (m *Marshaller) Inc(i ...int) {
	if i == nil {
		m.Cursor++
		return
	}
	m.Cursor += i[0]
}

func (m *Marshaller) IsSpace() bool {
	return InBytes(m.Buffer[m.Cursor], m.Space)
}

func (m *Marshaller) IsIndent() bool {
	return MatchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.Indent)], m.Indent)
}

func (m *Marshaller) IsQuote() bool {
	return InBytes(m.Buffer[m.Cursor], m.Quote)
}

func (m *Marshaller) IsEscape() bool {
	return InBytes(m.Buffer[m.Cursor], m.Escape)
}

func (m *Marshaller) IsValEnd() bool {
	return InBytes(m.Buffer[m.Cursor], m.ValEnd)
}

func (m *Marshaller) IsKeyEnd() bool {
	return InBytes(m.Buffer[m.Cursor], m.KeyEnd)
}

func (m *Marshaller) IsBlockCommentStart() bool {
	return MatchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.BlockCommentStart)], m.BlockCommentStart)
}

func (m *Marshaller) IsBlockCommentEnd() bool {
	return MatchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.BlockCommentEnd)], m.BlockCommentEnd)
}
func (m *Marshaller) IsLineCommentStart() bool {
	return MatchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.LineCommentStart)], m.LineCommentStart)
}

func (m *Marshaller) IsLineCommentEnd() bool {
	return MatchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.LineCommentEnd)], m.LineCommentEnd)
}

func (m *Marshaller) IsSliceStart() bool {
	return MatchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.SliceStart)], m.SliceStart)
}

func (m *Marshaller) IsSliceEnd() bool {
	return MatchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.SliceEnd)], m.SliceEnd)
}

func (m *Marshaller) IsMapStart() bool {
	return MatchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.MapStart)], m.MapStart)
}

func (m *Marshaller) IsMapEnd() bool {
	return MatchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.MapEnd)], m.MapEnd)
}

func MatchBytes(a, b []byte) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i, s := range a {
		if s != b[i] {
			return false
		}
	}
	return true
}

func InBytes(a byte, b []byte) bool {
	if b == nil {
		return false
	}
	for _, s := range b {
		if s == a {
			return true
		}
	}
	return false
}
