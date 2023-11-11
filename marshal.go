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
// [x] Test Inline Yaml
// [ ] Test Yaml
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
	// marshaller state
	cursor      int    // the current position in the buffer
	curIndent   int    // the current indentation level
	inline      bool   // true when marshalling in single line syntax
	hasBrackets bool   // true when marshalling with brackets around data objects
	value       any    // the value being marshalled
	buffer      []byte // the buffer being marshalled to
	len         int    // the length of the buffer
	// marshalling syntax
	Type              string // the type of marshaller. json, yaml, etc.
	Space             []byte // the space characters
	LineBreak         []byte // the line break characters
	Indent            []byte // the indentation characters
	Quote             []byte // the quote characters, first is default, additional are alternate
	Escape            []byte // the escape characters for string special char escape
	Null              []byte // the null value
	ValEnd            []byte // the characters that separate values
	KeyEnd            []byte // the characters that separate keys and values
	BlockCommentStart []byte // the characters that start a block comment
	BlockCommentEnd   []byte // the characters that end a block comment
	LineCommentStart  []byte // the characters that start a single line comment
	LineCommentEnd    []byte // the characters that end a single line comment
	SliceStart        []byte // the characters that start a slice or array
	SliceEnd          []byte // the characters that end a slice or array
	SliceElem         []byte // the characters before each slice element
	MapStart          []byte // the characters that start a hash map
	MapEnd            []byte // the characters that end a hash map
	// marshalling flags
	Format           bool // when true, marshal with formatting, indentation, and line breaks
	FormatWithSpaces bool // when true, marshal with space between keys and values
	NoIndentFirst    bool // when true, do not indent first line of slice or map
	CascadeOnlyDeep  bool // when true, marshal single-depth slices and maps with inline syntax
	QuotedKey        bool // when true, marshal map keys with quotes
	QuotedString     bool // when true, marshal strings with quotes
	QuotedSpecial    bool // when true, marshal strings with quotes if they contain special characters
	QuotedNum        bool // when true, marshal numbers with quotes
	QuotedBool       bool // when true, marshal bools with quotes
	QuotdeNull       bool // when true, marshal null with quotes
	//Inline *Marshaller
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
		FormatWithSpaces:  true,
		CascadeOnlyDeep:   true,
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
		//NoIndentFirst: true,
		Space:  []byte(" \t\v\f\r"),
		Indent: []byte("  "),
		//lnBreak:          []byte(""),
		Quote:  []byte(`"'`),
		Escape: []byte(`\`),
		//ValEnd:           []byte("\n"),
		KeyEnd:           []byte(":"),
		LineCommentStart: []byte("#"),
		LineCommentEnd:   []byte("\n"),
		SliceElem:        []byte("- "),
		//Inline:           &MarshallerInlineYaml,
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
		if m.LineBreak == nil {
			m.LineBreak = []byte("\n")
		}
		if m.FormatWithSpaces {
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
	m.hasBrackets = !(m.MapStart == nil || m.MapEnd == nil || m.SliceStart == nil || m.SliceEnd == nil)
}

func (m *Marshaller) SetType(t string) {
	m.Type = t
}

func (m *Marshaller) Reset() {
	m.buffer = []byte{}
	m.cursor = 0
	m.curIndent = 0
	m.len = 0
	m.value = nil
}

// ------------------------------------------------------------ /
// Getter Utilities
// for getting non-exported state values from the marshaller
// ------------------------------------------------------------ /

func (m *Marshaller) Buffer() []byte {
	return m.buffer
}

func (m *Marshaller) Cursor() int {
	return m.cursor
}

func (m *Marshaller) CurIndent() int {
	return m.curIndent
}

func (m *Marshaller) Len() int {
	return m.len
}

func (m *Marshaller) Value() any {
	return m.value
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
		return m.buffer, nil
	}
	if a.Len() == 0 {
		m.ToBuffer(append(m.SliceStart, m.SliceEnd...))
		return m.buffer, nil
	}
	i := m.SetInline(a.typ)
	m.MarshalSliceStart()
	a.ForEach(func(i int, k string, v VALUE) (brake bool) {
		m.MarshalNext(i)
		m.ToBuffer(m.SliceElem)
		m.marshal(v)
		return
	})
	m.MarshalSliceEnd()
	m.UnsetInline(i)
	return m.buffer, nil
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
		return m.buffer, nil
	}
	if hm.Len() == 0 {
		m.ToBuffer(append(m.MapStart, m.MapEnd...))
		return m.buffer, nil
	}
	i := m.SetInline(hm.typ)
	m.MarshaMapStart()
	hm.ForEach(func(i int, k string, v VALUE) (brake bool) {
		m.MarshalNext(i)
		m.MarshalKey([]byte(k))
		m.marshal(v)
		return
	})
	m.MarshaMapEnd()
	m.UnsetInline(i)
	return m.buffer, nil
}

func (m *Marshaller) MarshalSlice(s SLICE, ancestry ...ancestor) ([]byte, error) {
	if s.ptr == nil {
		m.ToBuffer(m.Null)
		return m.buffer, nil
	}
	if s.Len() == 0 {
		m.ToBuffer(append(m.SliceStart, m.SliceEnd...))
		return m.buffer, nil
	}
	i := m.SetInline(s.typ)
	m.MarshalSliceStart()
	s.ForEach(func(i int, k string, v VALUE) (brake bool) {
		m.MarshalNext(i)
		m.ToBuffer(m.SliceElem)
		m.marshal(v)
		return
	})
	m.MarshalSliceEnd()
	m.UnsetInline(i)
	return m.buffer, nil
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
	return m.buffer, nil
}

func (m *Marshaller) MarshalStruct(s STRUCT, ancestry ...ancestor) ([]byte, error) {
	if s.ptr == nil {
		m.ToBuffer(m.Null)
		return m.buffer, nil
	}
	if s.Len() == 0 {
		m.ToBuffer(append(m.MapStart, m.MapEnd...))
		return m.buffer, nil
	}
	il := m.SetInline(s.typ)
	m.MarshaMapStart()
	fields, i := s.TagIndex(m.Type), 0
	for k, f := range fields {
		m.MarshalNext(i)
		m.MarshalKey([]byte(k))
		m.marshal(f.VALUE())
		i++
	}
	m.MarshaMapEnd()
	m.UnsetInline(il)
	return m.buffer, nil
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
	m.curIndent++
	m.ToBuffer(m.SliceStart)
}

func (m *Marshaller) MarshalSliceEnd() {
	m.curIndent--
	m.MarshalNext(-1)
	m.ToBuffer(m.SliceEnd)
}

func (m *Marshaller) MarshaMapStart() {
	m.curIndent++
	m.ToBuffer(m.MapStart)
}

func (m *Marshaller) MarshaMapEnd() {
	m.curIndent--
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
	m.buffer = append(m.buffer, k...)
}

func (m *Marshaller) MarshalNext(i ...int) {
	idx := 1
	if i != nil {
		idx = i[0]
	}
	if idx > 0 {
		m.ToBuffer(m.ValEnd)
	}
	if m.Format && !m.inline {
		m.ToBuffer(m.LineBreak)
		if !m.NoIndentFirst || idx > 0 {
			m.MarshalIndent()
		}
	}
}

func (m *Marshaller) MarshalIndent() {
	i := m.curIndent
	if !m.hasBrackets && i > 0 {
		i--
	}
	m.ToBuffer(bytes.Repeat(m.Indent, i))
}

func (m *Marshaller) SetInline(t *TYPE) bool {
	i := m.inline
	if m.CascadeOnlyDeep {
		m.inline = !t.HasDataElem()
	}
	return i
}

func (m *Marshaller) UnsetInline(i bool) {
	m.inline = i
}

func (m *Marshaller) ToBuffer(b []byte) {
	if b != nil {
		m.buffer = append(m.buffer, b...)
		m.len = len(m.buffer)
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
		s := m.cursor
		q := m.buffer[s]
		for m.cursor < m.len {
			if m.IsEscape() {
				m.Inc(2)
				continue
			}
			if m.ByteIs(q) {
				m.Inc()
				break
			}

		}
		return m.buffer[s:m.cursor]
	}
	return nil
}

func (m *Marshaller) ParseWhitespace() []byte {
	if m.IsSpace() {
		s := m.cursor
		for m.cursor < m.len {
			if !m.IsSpace() {
				break
			}
			m.Inc()
		}
		return m.buffer[s:m.cursor]
	}
	return nil
}

func (m *Marshaller) ByteIs(b byte) bool {
	return m.Byte() == b
}

func (m *Marshaller) Byte() byte {
	if m.len == 0 {
		return 0
	}
	return m.buffer[m.cursor]
}

func (m *Marshaller) Inc(i ...int) {
	if i == nil {
		m.cursor++
		return
	}
	m.cursor += i[0]
}

func (m *Marshaller) IsSpace() bool {
	return InBytes(m.buffer[m.cursor], m.Space)
}

func (m *Marshaller) IsIndent() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.Indent)], m.Indent)
}

func (m *Marshaller) IsQuote() bool {
	return InBytes(m.buffer[m.cursor], m.Quote)
}

func (m *Marshaller) IsEscape() bool {
	return InBytes(m.buffer[m.cursor], m.Escape)
}

func (m *Marshaller) IsValEnd() bool {
	return InBytes(m.buffer[m.cursor], m.ValEnd)
}

func (m *Marshaller) IsKeyEnd() bool {
	return InBytes(m.buffer[m.cursor], m.KeyEnd)
}

func (m *Marshaller) IsBlockCommentStart() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.BlockCommentStart)], m.BlockCommentStart)
}

func (m *Marshaller) IsBlockCommentEnd() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.BlockCommentEnd)], m.BlockCommentEnd)
}
func (m *Marshaller) IsLineCommentStart() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.LineCommentStart)], m.LineCommentStart)
}

func (m *Marshaller) IsLineCommentEnd() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.LineCommentEnd)], m.LineCommentEnd)
}

func (m *Marshaller) IsSliceStart() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.SliceStart)], m.SliceStart)
}

func (m *Marshaller) IsSliceEnd() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.SliceEnd)], m.SliceEnd)
}

func (m *Marshaller) IsMapStart() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.MapStart)], m.MapStart)
}

func (m *Marshaller) IsMapEnd() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.MapEnd)], m.MapEnd)
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
