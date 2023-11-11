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
// TODO:
// [x] Test Indent Json
// [x] Inc Length of Buffer
// [ ] Test Yaml and Inline Yaml
// [ ] Finish Unmarshal
// [ ] Handle Comments
// [ ] Test Special Characters
// [ ] Handle Recursive Subvalues
// [ ] Remove Serialize functions
//     [ ] Add Yaml Type
//     [ ] Update JSON Type
//     [ ] Update Testing to use Marshaller
// ------------------------------------------------------------ /

// ------------------------------------------------------------ /
// MARSHALLER IMPLEMENTATION
// a generic marshaller for serializing and deserializing
// golang values to and from []byte in a specific format
// ------------------------------------------------------------ /

type Marshaller struct {
	Type      string
	Cursor    int
	CurIndent int
	Value     any
	Len       int
	Buffer,
	Space,
	lnBreak,
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
	SliceElem,
	MapStart,
	MapEnd []byte
	Format,
	Fspace,
	hasDataBrackets,
	inline,
	QuotedKey,
	QuotedString,
	QuotedSpecial,
	QuotedNum,
	QuotedBool,
	QuotdeNull bool
	Inline *Marshaller
}

// ------------------------------------------------------------ /
// PRESET MARSHALLERS
// JSON, YAML...
// ------------------------------------------------------------ /

var (
	MarshallerJson = Marshaller{
		Type:              "json",
		QuotedKey:         true,
		QuotedString:      true,
		Fspace:            true,
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
		Type:          "yaml",
		Format:        true,
		QuotedSpecial: true,
		Space:         []byte(" \t\v\f\r"),
		Indent:        []byte("  "),
		//lnBreak:          []byte(""),
		Quote:            []byte(`"'`),
		Escape:           []byte(`\`),
		ValEnd:           []byte("\n"),
		KeyEnd:           []byte(":"),
		LineCommentStart: []byte("#"),
		LineCommentEnd:   []byte("\n"),
		SliceElem:        []byte("- "),
		Inline:           &MarshallerInlineYaml,
	}
	MarshallerInlineYaml = Marshaller{
		Type:             "yaml",
		QuotedSpecial:    true,
		Space:            []byte(" \t\v\f\r"),
		Indent:           []byte("  "),
		Quote:            []byte(`"'`),
		Escape:           []byte(`\`),
		ValEnd:           []byte(", "),
		KeyEnd:           []byte(": "),
		LineCommentStart: []byte("#"),
		LineCommentEnd:   []byte("\n"),
		SliceStart:       []byte("["),
		SliceEnd:         []byte("]"),
		MapStart:         []byte("{"),
		MapEnd:           []byte("}"),
	}
)

// ------------------------------------------------------------ /
// Init Utilities
// methods for setting up the marshaller
// ------------------------------------------------------------ /

func (m *Marshaller) Init() {
	m.Reset()
	if m.Null == nil {
		m.Null = []byte("null")
	}
	if m.Space == nil {
		m.Space = []byte(" \t\n\v\f\r")
	}
	if m.Format {
		if m.Indent == nil {
			m.Indent = []byte("  ")
		}
		if m.lnBreak == nil {
			if m.ValEnd[0] == '\n' {
				m.lnBreak = []byte("")
			} else {
				m.lnBreak = []byte("\n")
			}
		}
		if m.Fspace {
			m.KeyEnd = append(m.KeyEnd, m.Space[0])
			m.ValEnd = append(m.ValEnd, m.Space[0])
		}
	}
	if m.Quote == nil && (m.QuotedString || m.QuotedKey || m.QuotedSpecial || m.QuotedNum || m.QuotedBool || m.QuotdeNull) {
		m.Quote = []byte(`"`)
	}
	if m.Escape == nil && m.Quote != nil {
		m.Escape = []byte(`\`)
	}
	m.hasDataBrackets = !(m.MapStart == nil || m.MapEnd == nil || m.SliceStart == nil || m.SliceEnd == nil)
}

func (m *Marshaller) SetType(t string) {
	m.Type = t
}

func (m *Marshaller) Reset() {
	m.Buffer = []byte{}
	m.Cursor = 0
	m.CurIndent = 0
	m.Len = 0
	m.Value = nil
}

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
		return m.MarshalString(*(*string)(v.ptr))
	case Struct:
		return m.MarshalStruct((STRUCT)(v), ancestry...)
	case UnsafePointer:
		return m.MarshalUnsafePointer(*(*unsafe.Pointer)(v.ptr))
	case Field:
		return m.MarshalField(*(*FIELD)(unsafe.Pointer(&v)))
	case Time:
		return m.MarshalTime(*(*TIME)(v.ptr))
	case Uuid:
		return m.MarshalUuid(*(*UUID)(v.ptr))
	case Bytes:
		return m.MarshalBytes(*(*BYTES)(v.ptr))
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
	i := m.inline
	m.inline = !a.TYPE().HasDataElem()
	m.MarshalSliceStart()
	a.ForEach(func(i int, k string, v VALUE) (brake bool) {
		m.MarshalNext(i)
		m.ToBuffer(m.SliceElem)
		m.marshal(v)
		return
	})
	m.MarshalSliceEnd()
	m.inline = i
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
	i := m.inline
	m.inline = !hm.TYPE().HasDataElem()
	m.MarshaMapStart()
	hm.ForEach(func(i int, k string, v VALUE) (brake bool) {
		m.MarshalNext(i)
		m.MarshalKey([]byte(k))
		m.marshal(v)
		return
	})
	m.MarshaMapEnd()
	m.inline = i
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
	i := m.inline
	m.inline = !s.TYPE().HasDataElem()
	m.MarshalSliceStart()
	s.ForEach(func(i int, k string, v VALUE) (brake bool) {
		m.MarshalNext(i)
		m.ToBuffer(m.SliceElem)
		m.marshal(v)
		return
	})
	m.MarshalSliceEnd()
	m.inline = i
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
	il := m.inline
	m.inline = !s.TYPE().HasDataElem()
	m.MarshaMapStart()
	fields, i := s.TagIndex(m.Type), 0
	for k, f := range fields {
		m.MarshalNext(i)
		m.MarshalKey([]byte(k))
		m.marshal(f.VALUE())
		i++
	}
	m.MarshaMapEnd()
	m.inline = il
	return m.Buffer, nil
}

func (m *Marshaller) MarshalUnsafePointer(p unsafe.Pointer) (bytes []byte, err error) {
	return m.MarshalString(fmt.Sprintf("%p", p))
}

func (m *Marshaller) MarshalField(f FIELD) (bytes []byte, err error) {
	return m.marshal(f.VALUE())
}

func (m *Marshaller) MarshalTime(t TIME) (bytes []byte, err error) {
	return m.MarshalString(t.String())
}

func (m *Marshaller) MarshalUuid(u UUID) (bytes []byte, err error) {
	return m.MarshalString(u.String())
}

func (m *Marshaller) MarshalBytes(b BYTES) (bytes []byte, err error) {
	return m.MarshalString(string(b))
}

func (m *Marshaller) MarshalSliceStart() {
	m.CurIndent++
	m.ToBuffer(m.SliceStart)
}

func (m *Marshaller) MarshalSliceEnd() {
	m.CurIndent--
	m.MarshalNext(-1)
	m.ToBuffer(m.SliceEnd)
}

func (m *Marshaller) MarshaMapStart() {
	m.CurIndent++
	m.ToBuffer(m.MapStart)
}

func (m *Marshaller) MarshaMapEnd() {
	m.CurIndent--
	m.MarshalNext(-1)
	m.ToBuffer(m.MapEnd)
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

func (m *Marshaller) MarshalNext(i ...int) {
	idx := 1
	if i != nil {
		idx = i[0]
	}
	if idx > 0 {
		m.ToBuffer(m.ValEnd)
	}
	if m.Format /* && !m.inline */ {
		m.ToBuffer(m.lnBreak)
		m.MarshalIndent()
	}
}

func (m *Marshaller) MarshalIndent() {
	i := m.CurIndent
	if !m.hasDataBrackets && i > 0 {
		i--
	}
	m.ToBuffer(bytes.Repeat(m.Indent, i))
}

func (m *Marshaller) ToBuffer(b []byte) {
	if b != nil {
		m.Buffer = append(m.Buffer, b...)
		m.Len = len(m.Buffer)
	}
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
