// Copyright 2023 james dotter. All rights reserved.typVal
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

// ------------------------------------------------------------ /
// TODO:
// [x] Test Indent Json
// [x] Inc Length of Buffer
// [x] Test Inline Yaml
// [x] Test Yaml
//		rework bookends and delimiters
//		set: start, delim and end
// [ ] Finish Unmarshal
//     [x] whitespace
//     [x] quote
//     [x] comments
//     [ ] indent / object end / negative indent
//     [ ] string
//     [ ] map
//	   [ ] slice
//	   [ ] general deserializer - decides type and routes to appropriate deserializer
//     -- output types: string, map[string]any, []any (null, number and bool are all strings)
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
	curDepth    int    // the current depth of the data structure
	curIndent   int    // the current indentation level
	inline      bool   // true when marshalling in single line syntax
	hasBrackets bool   // true when marshalling with brackets around data objects
	mapItem     bool   // true when marshalling a map item
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
	SliceItem         []byte // the characters before each slice element
	MapStart          []byte // the characters that start a hash map
	MapEnd            []byte // the characters that end a hash map
	// marshalling flags
	Format           bool // when true, marshal with formatting, indentation, and line breaks
	FormatWithSpaces bool // when true, marshal with space between keys and values
	NoFormatFirst    bool // when true, do not break or indent first line of slice
	CascadeOnlyDeep  bool // when true, marshal single-depth slices and maps with inline syntax
	QuotedKey        bool // when true, marshal map keys with quotes
	QuotedString     bool // when true, marshal strings with quotes
	QuotedSpecial    bool // when true, marshal strings with quotes if they contain special characters
	QuotedNum        bool // when true, marshal numbers with quotes
	QuotedBool       bool // when true, marshal bools with quotes
	QuotedNull       bool // when true, marshal null with quotes
	RecursiveName    bool // when true, include name, string or type of recursive struct value in marshalling, otherwise, exclude all
}

type components struct {
	start, delim, end []byte
}

// ------------------------------------------------------------ /
// PRESET MARSHALLERS
// JSON, YAML...
// ------------------------------------------------------------ /

var (
	MarshallerJson = Marshaller{
		Type:              "json",
		FormatWithSpaces:  true,
		CascadeOnlyDeep:   true,
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
		Format:           true,
		FormatWithSpaces: true,
		NoFormatFirst:    true,
		QuotedSpecial:    true,
		Space:            []byte(" \t\v\f\r"),
		Indent:           []byte("  "),
		Quote:            []byte(`"'`),
		Escape:           []byte(`\`),
		KeyEnd:           []byte(":"),
		LineCommentStart: []byte("#"),
		LineCommentEnd:   []byte("\n"),
		SliceItem:        []byte("- "),
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
	if m.Quote == nil && (m.QuotedString || m.QuotedKey || m.QuotedSpecial || m.QuotedNum || m.QuotedBool || m.QuotedNull) {
		m.Quote = []byte(`"`)
	}
	if m.Escape == nil && m.Quote != nil {
		m.Escape = []byte(`\`)
	}
	m.hasBrackets = !(m.MapStart == nil || m.MapEnd == nil || m.SliceStart == nil || m.SliceEnd == nil)
	if !m.hasBrackets && !m.Format {
		panic("cannot marshal without brackets or formatting: unable to determine data structure")
	}
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
		return m.ToBuffer(m.Null)
	}
	if a.Len() == 0 {
		if !m.hasBrackets {
			return m.ToBuffer(m.Null)
		}
		return m.ToBuffer(append(m.SliceStart, m.SliceEnd...))
	}
	m.IncDepth()
	start, delim, end := m.SliceComponents(a.typ, ancestry...)
	ancestry = append([]ancestor{{a.typ, uintptr(a.ptr)}}, ancestry...)
	m.ToBuffer(start)
	a.ForEach(func(i int, k string, v VALUE) (brake bool) {
		b, recursive := m.RecursiveValue(v, ancestry)
		if recursive && b == nil {
			return
		}
		if i > 0 {
			m.ToBuffer(delim)
		}
		if b != nil {
			m.MarshalBytes(b)
			return
		}
		m.marshal(v, ancestry...)
		return
	})
	m.DecDepth()
	return m.ToBuffer(end)
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
		return m.ToBuffer(m.Null)
	}
	if hm.Len() == 0 {
		if !m.hasBrackets {
			return m.ToBuffer(m.Null)
		}
		return m.ToBuffer(append(m.MapStart, m.MapEnd...))
	}
	m.IncDepth()
	start, delim, end := m.MapComponents(hm.typ, ancestry)
	ancestry = append([]ancestor{{hm.typ, uintptr(hm.ptr)}}, ancestry...)
	m.ToBuffer(start)
	hm.ForEach(func(i int, k string, v VALUE) (brake bool) {
		b, recursive := m.RecursiveValue(v, ancestry)
		if recursive && b == nil {
			return
		}
		if i > 0 {
			m.ToBuffer(delim)
		}
		m.MarshalKey([]byte(k))
		if b != nil {
			m.MarshalBytes(b)
			return
		}
		m.marshal(v, ancestry...)
		return
	})
	m.DecDepth()
	return m.ToBuffer(end)
}

func (m *Marshaller) MarshalSlice(s SLICE, ancestry ...ancestor) ([]byte, error) {
	if s.ptr == nil {
		return m.ToBuffer(m.Null)
	}
	if s.Len() == 0 {
		if !m.hasBrackets {
			return m.ToBuffer(m.Null)
		}
		return m.ToBuffer(append(m.SliceStart, m.SliceEnd...))
	}
	m.IncDepth()
	start, delim, end := m.SliceComponents(s.typ, ancestry...)
	ancestry = append([]ancestor{{s.typ, uintptr(s.ptr)}}, ancestry...)
	m.ToBuffer(start)
	s.ForEach(func(i int, k string, v VALUE) (brake bool) {
		b, recursive := m.RecursiveValue(v, ancestry)
		if recursive && b == nil {
			return
		}
		if i > 0 {
			m.ToBuffer(delim)
		}
		if b != nil {
			m.MarshalBytes(b)
			return
		}
		m.marshal(v, ancestry...)
		return
	})
	m.DecDepth()
	return m.ToBuffer(end)
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
		return m.ToBuffer(m.Null)
	}
	if s.Len() == 0 {
		if !m.hasBrackets {
			return m.ToBuffer(m.Null)
		}
		return m.ToBuffer(append(m.MapStart, m.MapEnd...))
	}
	m.IncDepth()
	start, delim, end := m.MapComponents(s.typ, ancestry)
	ancestry = append([]ancestor{{s.typ, uintptr(s.ptr)}}, ancestry...)
	m.ToBuffer(start)
	fields, i := s.TagIndex(m.Type), 0
	for k, f := range fields {
		b, recursive := m.RecursiveValue(f.VALUE(), ancestry)
		if recursive && b == nil {
			continue
		}
		if i > 0 {
			m.ToBuffer(delim)
		}
		m.MarshalKey([]byte(k))
		if b != nil {
			m.MarshalBytes(b)
			continue
		}
		m.marshal(f.VALUE())
		i++
	}
	m.DecDepth()
	return m.ToBuffer(end)
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

func (m *Marshaller) SliceComponents(t *TYPE, ancestry ...ancestor) (start []byte, delim []byte, end []byte) {
	inline := m.CascadeOnlyDeep && !t.HasDataElem()
	mapItem := false
	if ancestry != nil {
		k := ancestry[0].typ.Kind()
		mapItem = k == Map || k == Struct
	}
	if m.Format && !inline {
		lastline := append(m.LineBreak, bytes.Repeat(m.Indent, m.curIndent)...)
		newline := lastline
		if m.hasBrackets {
			newline = append(newline, m.Indent...)
			end = JoinBytes(lastline, m.SliceEnd)
		}
		if m.NoFormatFirst && !mapItem {
			start = JoinBytes(m.SliceStart, m.SliceItem)
		} else {
			start = JoinBytes(m.SliceStart, newline, m.SliceItem)
		}
		delim = JoinBytes(m.ValEnd, newline, m.SliceItem)
		return
	}
	start, delim, end = CopyBytes(m.SliceStart), JoinBytes(m.ValEnd, m.SliceItem), CopyBytes(m.SliceEnd)
	return
}

func (m *Marshaller) MarshalSliceStart() {
	m.ToBuffer(m.SliceStart)
	m.IncDepth()
	m.MarshalBreak(0)
}

func (m *Marshaller) MarshalSliceEnd() {
	m.DecDepth()
	m.MarshalBreak(-1)
	m.ToBuffer(m.SliceEnd)
}

func (m *Marshaller) MapComponents(t *TYPE, ancestry []ancestor) (start []byte, delim []byte, end []byte) {
	inline := m.CascadeOnlyDeep && !t.HasDataElem()
	mapItem := false
	if ancestry != nil {
		k := ancestry[0].typ.Kind()
		mapItem = k == Map || k == Struct
	}
	if m.Format && !inline {
		lastline := append(m.LineBreak, bytes.Repeat(m.Indent, m.curIndent)...)
		newline := lastline
		if m.hasBrackets {
			newline = append(newline, m.Indent...)
			end = JoinBytes(lastline, m.MapEnd)
		}
		if m.NoFormatFirst && !mapItem {
			start = CopyBytes(m.MapStart)
		} else {
			start = JoinBytes(m.MapStart, newline)
		}
		delim = JoinBytes(m.ValEnd, newline)
		return
	}
	start, delim, end = CopyBytes(m.SliceStart), CopyBytes(m.ValEnd), CopyBytes(m.SliceEnd)
	return
}

func (m *Marshaller) MarshaMapStart() {
	m.ToBuffer(m.MapStart)
	m.IncDepth()
	m.MarshalBreak(0)
}

func (m *Marshaller) MarshaMapEnd() {
	m.DecDepth()
	m.MarshalBreak(-1)
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

func (m *Marshaller) MarshalNext(i int) {
	if i != 0 {
		m.ToBuffer(m.ValEnd)
		m.MarshalBreak(i)
	}
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

func (m *Marshaller) IncDepth() {
	m.curDepth++
	m.SetIndent()
}

func (m *Marshaller) DecDepth() {
	m.curDepth--
	m.SetIndent()
}

func (m *Marshaller) SetIndent() {
	if !m.hasBrackets {
		m.curIndent = m.curDepth - 1
	} else {
		m.curIndent = m.curDepth
	}
}

func (m *Marshaller) MarshalBreak(i int) {
	if m.Format && !m.inline && (!m.NoFormatFirst || i > 0 || (m.mapItem && i != -1)) {
		m.ToBuffer(m.LineBreak)
		m.ToBuffer(bytes.Repeat(m.Indent, m.curIndent))
	}
}

func (m *Marshaller) ToBuffer(b []byte) ([]byte, error) {
	if b != nil {
		m.buffer = append(m.buffer, b...)
		m.len = len(m.buffer)
	}
	return m.buffer, nil
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

func IsRecursive(v VALUE, ancestry []ancestor) bool {
	for _, a := range ancestry {
		if a.pointer == uintptr(v.ptr) && a.typ == v.typ {
			return true
		}
	}
	return false
}

func (m *Marshaller) RecursiveValue(v VALUE, ancestry []ancestor) (bytes []byte, is bool) {
	for _, a := range ancestry {
		if a.pointer == uintptr(v.ptr) && a.typ == v.typ {
			is = true
			if v.Kind() == Struct && m.RecursiveName {
				if _, ok := v.typ.Reflect().MethodByName("Name"); ok {
					bytes = v.Reflect().MethodByName("Name").Call([]reflect.Value{})[0].Bytes()
				} else if _, ok := v.typ.Reflect().MethodByName("String"); ok {
					bytes = v.Reflect().MethodByName("String").Call([]reflect.Value{})[0].Bytes()
				} else {
					bytes = []byte(v.typ.NameShort())
				}
			}
			return
		}
	}
	return nil, false
}

// ------------------------------------------------------------ /
// Unmarshal Utilities
// methods for unmarshalling from type
// ------------------------------------------------------------ /

func (m *Marshaller) UnmarshalItem() []byte {
	switch {
	case m.IsQuote():
	case m.IsBlockCommentStart():
	case m.IsLineCommentStart():
	case m.IsMapStart():
	case m.IsSliceStart():
	case m.IsSpace():
		m.UnmarshalWhitespace()
	default:
		return m.UnmarshalString()
	}
	return nil
}

func (m *Marshaller) UnmarshalKey() []byte {
	s := m.cursor
	for m.cursor < m.len {
		if m.IsKeyEnd() {
			m.Inc()
			break
		}
	}
	return m.buffer[s:m.cursor]
}

func (m *Marshaller) UnmarshalString() []byte {
	s := m.cursor
	for m.cursor < m.len {
		if m.IsValEnd() {
			m.Inc()
			break
		}
	}
	return m.buffer[s:m.cursor]
}

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

func (m *Marshaller) UnmarshalWhitespace() []byte {
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

func (m *Marshaller) UnmarshalCommentBlock() []byte {
	if m.IsBlockCommentStart() {
		s := m.cursor
		for m.cursor < m.len {
			if m.IsBlockCommentEnd() {
				m.Inc(len(m.BlockCommentEnd))
				break
			}
			m.Inc()
		}
		return m.buffer[s:m.cursor]
	}
	return nil
}

func (m *Marshaller) UnmarshalInlineComment() []byte {
	if m.IsLineCommentStart() {
		s := m.cursor
		for m.cursor < m.len {
			if m.IsLineCommentEnd() {
				m.Inc(len(m.LineCommentEnd))
				break
			}
			m.Inc()
		}
		return m.buffer[s:m.cursor]
	}
	return nil
}
func (m *Marshaller) UnmarshalTo(stop []byte) []byte {
	s := m.cursor
	e := stop[0]
	for m.cursor < m.len {
		if m.ByteIs(e) {
			if MatchBytes(m.buffer[m.cursor:m.cursor+len(stop)], stop) {
				m.Inc(len(stop))
				break
			}
		}
	}
	return m.buffer[s:m.cursor]
}

func (m *Marshaller) UnmarshalBreak() []byte {
	if m.IsLineBreak() {
		s := m.cursor
		m.Inc(len(m.LineBreak))
		return m.buffer[s:m.cursor]
	}
	return nil
}

func (m *Marshaller) UnmarshalIndents() (b []byte, i int) {
	for m.cursor < m.len {
		if m.IsIndent() {
			m.Inc(len(m.Indent))
			i++
			continue
		}
		break
	}
	return m.buffer[m.cursor : m.cursor+i], i
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

func (m *Marshaller) IsLineBreak() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.LineBreak)], m.LineBreak)
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
	if m.SliceStart != nil {
		return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.SliceStart)], m.SliceStart)
	}
	if !m.hasBrackets {
		if m.curIndent == 0 {
			return m.IsSliceItem()
		}
		return m.IsFmtObjectStart()
	}
	return false
}

func (m *Marshaller) IsSliceEnd() bool {
	if m.SliceEnd != nil {
		return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.SliceEnd)], m.SliceEnd)
	}
	if !m.hasBrackets {
		if m.curIndent == 0 {
			return m.IsSliceItem()
		}
		return m.IsFmtObjectStart()
	}
	return false
}

func (m *Marshaller) IsSliceItem() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.SliceItem)], m.SliceItem)
}

func (m *Marshaller) IsMapStart() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.MapStart)], m.MapStart)
}

func (m *Marshaller) IsMapEnd() bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(m.MapEnd)], m.MapEnd)
}

func (m *Marshaller) IsFmtObjectStart() bool {
	if !m.hasBrackets {
		start := append(m.LineBreak, bytes.Repeat(m.Indent, m.curIndent+1)...)
		return MatchBytes(m.buffer[m.cursor:m.cursor+len(start)], start)
	}
	return false
}

func (m *Marshaller) IsFmtObjectEnd() bool {
	if !m.hasBrackets {
		end := append(m.LineBreak, bytes.Repeat(m.Indent, m.curIndent-1)...)
		return MatchBytes(m.buffer[m.cursor:m.cursor+len(end)], end)
	}
	return false
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
