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
	cursor    int // the current position in the buffer
	curDepth  int // the current depth of the data structure
	curIndent int // the current indentation level
	//inline      bool   // true when marshalling in single line syntax
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
	InlineSyntax      *InlineSyntax
	// marshalling flags
	Format           bool // when true, marshal with formatting, indentation, and line breaks
	FormatWithSpaces bool // when true, marshal with space between keys and values
	//NoFormatFirst    bool // when true, do not break or indent first line of slice
	CascadeOnlyDeep bool // when true, marshal single-depth slices and maps with inline syntax
	QuotedKey       bool // when true, marshal map keys with quotes
	QuotedString    bool // when true, marshal strings with quotes
	QuotedSpecial   bool // when true, marshal strings with quotes if they contain special characters
	QuotedNum       bool // when true, marshal numbers with quotes
	QuotedBool      bool // when true, marshal bools with quotes
	QuotedNull      bool // when true, marshal null with quotes
	RecursiveName   bool // when true, include name, string or type of recursive struct value in marshalling, otherwise, exclude all
}

type InlineSyntax struct {
	ValEnd     []byte // the characters that separate values
	SliceStart []byte // the characters that start a slice or array
	SliceEnd   []byte // the characters that end a slice or array
	MapStart   []byte // the characters that start a hash map
	MapEnd     []byte // the characters that end a hash map
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
		//NoFormatFirst:    true,
		QuotedSpecial:    true,
		Space:            []byte(" \t\v\f\r"),
		Indent:           []byte("  "),
		Quote:            []byte(`"'`),
		Escape:           []byte(`\`),
		KeyEnd:           []byte(":"),
		LineCommentStart: []byte("#"),
		LineCommentEnd:   []byte("\n"),
		SliceItem:        []byte("- "),
		InlineSyntax: &InlineSyntax{
			ValEnd:     []byte(", "),
			SliceStart: []byte("["),
			SliceEnd:   []byte("]"),
			MapStart:   []byte("{"),
			MapEnd:     []byte("}"),
		},
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
	if m.SliceItem != nil && m.hasBrackets {
		panic("slice item reserved for bracketless marshalling")
	}
	if !m.hasBrackets && InBytes('\n', m.Space) {
		panic("cannot marshal without brackets when line breaks in space characters")
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
		if m.InlineSyntax == nil {
			m.InlineSyntax = &InlineSyntax{
				ValEnd:     m.ValEnd,
				SliceStart: m.SliceStart,
				SliceEnd:   m.SliceEnd,
				MapStart:   m.MapStart,
				MapEnd:     m.MapEnd,
			}
		}
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
		return m.marshalBool(v.Bool())
	case Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Float32, Float64, Complex64, Complex128:
		return m.marshalNum(v)
	case Array:
		return m.marshalArray((ARRAY)(v), ancestry...)
	case Func:
		return m.marshalFunc(v)
	case Interface:
		return m.marshalInterface(v)
	case Map:
		return m.marshalMap((MAP)(v), ancestry...)
	case Pointer:
		return m.marshal(v.Elem())
	case Slice:
		return m.marshalSlice((SLICE)(v), ancestry...)
	case String:
		return m.marshalString(*(*string)(v.ptr))
	case Struct:
		return m.marshalStruct((STRUCT)(v), ancestry...)
	case UnsafePointer:
		return m.marshalUnsafePointer(*(*unsafe.Pointer)(v.ptr))
	case Field:
		return m.marshalField(*(*FIELD)(unsafe.Pointer(&v)))
	case Time:
		return m.marshalTime(*(*TIME)(v.ptr))
	case Uuid:
		return m.marshalUuid(*(*UUID)(v.ptr))
	case Bytes:
		return m.marshalBytes(*(*BYTES)(v.ptr))
	default:
		return nil, errors.New("cannot marshal type '" + v.typ.String() + "'")
	}
}

func (m *Marshaller) marshalBool(b bool) (bytes []byte, err error) {
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

func (m *Marshaller) marshalNum(v VALUE) (bytes []byte, err error) {
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

func (m *Marshaller) marshalArray(a ARRAY, ancestry ...ancestor) ([]byte, error) {
	if a.ptr == nil {
		return m.ToBuffer(m.Null)
	}
	if a.Len() == 0 {
		if !m.hasBrackets {
			return m.ToBuffer(m.Null)
		}
		return m.ToBuffer(append(m.SliceStart, m.SliceEnd...))
	}
	start, delim, end := m.marshalSliceComponents(a.typ, ancestry)
	ancestry = append([]ancestor{{a.typ, uintptr(a.ptr)}}, ancestry...)
	m.ToBuffer(start)
	m.IncDepth()
	a.ForEach(func(i int, k string, v VALUE) (brake bool) {
		b, recursive := m.recursiveValue(v, ancestry)
		if recursive && b == nil {
			return
		}
		if i > 0 {
			m.ToBuffer(delim)
		}
		if b != nil {
			m.marshalBytes(b)
			return
		}
		m.marshal(v, ancestry...)
		return
	})
	m.DecDepth()
	return m.ToBuffer(end)
}

func (m *Marshaller) marshalFunc(v VALUE) (bytes []byte, err error) {
	bytes = []byte(v.typ.Name())
	m.ToBuffer(bytes)
	return
}

func (m *Marshaller) marshalInterface(v VALUE) (bytes []byte, err error) {
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

func (m *Marshaller) marshalMap(hm MAP, ancestry ...ancestor) ([]byte, error) {
	if hm.ptr == nil {
		return m.ToBuffer(m.Null)
	}
	if hm.Len() == 0 {
		if !m.hasBrackets {
			return m.ToBuffer(m.Null)
		}
		return m.ToBuffer(append(m.MapStart, m.MapEnd...))
	}
	start, delim, end := m.marshalMapComponents(hm.typ, ancestry)
	ancestry = append([]ancestor{{hm.typ, uintptr(hm.ptr)}}, ancestry...)
	m.ToBuffer(start)
	m.IncDepth()
	hm.ForEach(func(i int, k string, v VALUE) (brake bool) {
		b, recursive := m.recursiveValue(v, ancestry)
		if recursive && b == nil {
			return
		}
		if i > 0 {
			m.ToBuffer(delim)
		}
		m.marshalKey([]byte(k))
		if b != nil {
			m.marshalBytes(b)
			return
		}
		m.marshal(v, ancestry...)
		return
	})
	m.DecDepth()
	return m.ToBuffer(end)
}

func (m *Marshaller) marshalSlice(s SLICE, ancestry ...ancestor) ([]byte, error) {
	if s.ptr == nil {
		return m.ToBuffer(m.Null)
	}
	if s.Len() == 0 {
		if !m.hasBrackets {
			return m.ToBuffer(m.Null)
		}
		return m.ToBuffer(append(m.SliceStart, m.SliceEnd...))
	}
	start, delim, end := m.marshalSliceComponents(s.typ, ancestry)
	ancestry = append([]ancestor{{s.typ, uintptr(s.ptr)}}, ancestry...)
	m.ToBuffer(start)
	m.IncDepth()
	s.ForEach(func(i int, k string, v VALUE) (brake bool) {
		b, recursive := m.recursiveValue(v, ancestry)
		if recursive && b == nil {
			return
		}
		if i > 0 {
			m.ToBuffer(delim)
		}
		if b != nil {
			m.marshalBytes(b)
			return
		}
		m.marshal(v, ancestry...)
		return
	})
	m.DecDepth()
	return m.ToBuffer(end)
}

func (m *Marshaller) marshalString(s string) ([]byte, error) {
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

func (m *Marshaller) marshalStruct(s STRUCT, ancestry ...ancestor) ([]byte, error) {
	if s.ptr == nil {
		return m.ToBuffer(m.Null)
	}
	if s.Len() == 0 {
		if !m.hasBrackets {
			return m.ToBuffer(m.Null)
		}
		return m.ToBuffer(append(m.MapStart, m.MapEnd...))
	}
	start, delim, end := m.marshalMapComponents(s.typ, ancestry)
	ancestry = append([]ancestor{{s.typ, uintptr(s.ptr)}}, ancestry...)
	m.ToBuffer(start)
	m.IncDepth()
	fields, i := s.TagIndex(m.Type), 0
	for k, f := range fields {
		b, recursive := m.recursiveValue(f.VALUE(), ancestry)
		if recursive && b == nil {
			continue
		}
		if i > 0 {
			m.ToBuffer(delim)
		}
		m.marshalKey([]byte(k))
		if b != nil {
			m.marshalBytes(b)
			continue
		}
		m.marshal(f.VALUE())
		i++
	}
	m.DecDepth()
	return m.ToBuffer(end)
}

func (m *Marshaller) marshalUnsafePointer(p unsafe.Pointer) (bytes []byte, err error) {
	return m.marshalString(fmt.Sprintf("%p", p))
}

func (m *Marshaller) marshalField(f FIELD) (bytes []byte, err error) {
	return m.marshal(f.VALUE())
}

func (m *Marshaller) marshalTime(t TIME) (bytes []byte, err error) {
	return m.marshalString(t.String())
}

func (m *Marshaller) marshalUuid(u UUID) (bytes []byte, err error) {
	return m.marshalString(u.String())
}

func (m *Marshaller) marshalBytes(b BYTES) (bytes []byte, err error) {
	return m.marshalString(string(b))
}

func (m *Marshaller) marshalSliceComponents(t *TYPE, ancestry []ancestor) (start []byte, delim []byte, end []byte) {
	switch {
	case !m.Format:
		return m.SliceStart, m.ValEnd, m.SliceEnd
	case m.CascadeOnlyDeep && !t.HasDataElem():
		return m.InlineSyntax.SliceStart, m.InlineSyntax.ValEnd, m.InlineSyntax.SliceEnd
	case m.hasBrackets:
		return m.formattedSliceComponents()
	default:
		return m.bracketlessSliceComponents(ancestry)
	}
}

func (m *Marshaller) formattedSliceComponents() (start []byte, delim []byte, end []byte) {
	end = AppendBytes(m.LineBreak, bytes.Repeat(m.Indent, m.curIndent))
	nl := append(end, m.Indent...)
	start = AppendBytes(m.SliceStart, nl, m.SliceItem)
	delim = AppendBytes(m.ValEnd, nl, m.SliceItem)
	end = AppendBytes(end, m.SliceEnd)
	return
}

func (m *Marshaller) bracketlessSliceComponents(ancestry []ancestor) (start []byte, delim []byte, end []byte) {
	nl := AppendBytes(m.LineBreak, bytes.Repeat(m.Indent, m.curIndent))
	switch m.itemOf(ancestry) {
	case Map:
		start = AppendBytes(nl, m.Indent, m.SliceItem)
		delim = AppendBytes(m.ValEnd, start)
	case Slice:
		start = m.SliceItem
		delim = AppendBytes(m.ValEnd, nl, m.Indent, m.Indent, m.SliceItem)
	default:
		start = m.SliceItem
		delim = AppendBytes(m.ValEnd, nl, m.SliceItem)
	}
	return
}

func (m *Marshaller) marshalMapComponents(t *TYPE, ancestry []ancestor) (start []byte, delim []byte, end []byte) {
	switch {
	case !m.Format:
		return m.MapStart, m.ValEnd, m.MapEnd
	case m.Format && m.CascadeOnlyDeep && !t.HasDataElem():
		return m.InlineSyntax.MapStart, m.InlineSyntax.ValEnd, m.InlineSyntax.MapEnd
	case m.hasBrackets:
		return m.formattedMapComponents()
	default:
		return m.bracketlessMapComponents(ancestry)
	}
}

func (m *Marshaller) formattedMapComponents() (start []byte, delim []byte, end []byte) {
	end = AppendBytes(m.LineBreak, bytes.Repeat(m.Indent, m.curIndent))
	nl := append(end, m.Indent...)
	start = AppendBytes(m.MapStart, nl)
	delim = AppendBytes(m.ValEnd, nl)
	end = AppendBytes(end, m.MapEnd)
	return
}

func (m *Marshaller) bracketlessMapComponents(ancestry []ancestor) (start []byte, delim []byte, end []byte) {
	nl := AppendBytes(m.LineBreak, bytes.Repeat(m.Indent, m.curIndent))
	switch m.itemOf(ancestry) {
	case Map:
		start = AppendBytes(nl, m.Indent)
		delim = AppendBytes(m.ValEnd, nl, m.Indent)
	case Slice:
		delim = AppendBytes(m.ValEnd, nl, m.Indent)
	default:
		delim = AppendBytes(m.ValEnd, nl)
	}
	return
}

func (m *Marshaller) itemOf(ancestry []ancestor) KIND {
	if ancestry != nil {
		k := ancestry[0].typ.Kind()
		if k == Map || k == Struct {
			return Map
		}
		if (k == Slice || k == Array) && m.SliceItem != nil {
			return Slice
		}
	}
	return Invalid
}

func (m *Marshaller) marshalKey(k []byte) {
	if m.QuotedKey {
		q := m.Quote[:1]
		k = append(append(q, k...), q...)
	}
	if m.KeyEnd != nil {
		k = append(k, m.KeyEnd...)
	}
	m.buffer = append(m.buffer, k...)
}

func (m *Marshaller) IncDepth() {
	m.curDepth++
	m.setIndent()
}

func (m *Marshaller) DecDepth() {
	m.curDepth--
	m.setIndent()
}

func (m *Marshaller) setIndent() {
	if !m.hasBrackets && m.curDepth > 0 {
		m.curIndent = m.curDepth - 1
	} else {
		m.curIndent = m.curDepth
	}
}
func (m *Marshaller) ToBuffer(b []byte) ([]byte, error) {
	if b != nil {
		m.buffer = AppendBytes(m.buffer, b)
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

func (m *Marshaller) recursiveValue(v VALUE, ancestry []ancestor) (bytes []byte, is bool) {
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

/*
	Unmarshal Object: Comment, Space, Slice, Map.
	Unmarshal MapItem: Comment, Space, Key, Value.
	Unmarshal SliceItem: Comment, Space Value.


*/

func (m *Marshaller) Unmarshal() error {
	for m.cursor < m.len {
		m.unmarshalObject()
	}
	return nil
}

func (m *Marshaller) unmarshalObject(ancestry ...ancestor) []byte {
	m.unmarshalNonData()
	if start, delim, end, isSlice := m.unmarshalSliceStart(ancestry); isSlice {
		return m.unmarshalSlice(start, delim, end)
	}
	if start, delim, end, isMap := m.unmarshalMapStart(ancestry); isMap {
		return m.unmarshalMap(start, delim, end)
	}
	return nil
}

func (m *Marshaller) unmarshalSlice(start, delim, end []byte) []byte {
	// unmarshal item
	// check for end or delim
	// if delim, unmarshal item, otherwise end
	return nil
}

func (m *Marshaller) unmarshalMap(start, delim, end []byte) []byte {
	return nil
}

func (m *Marshaller) unmarshalItem() []byte {
	m.unmarshalNonData()
	// order: quote, null, slice, map, string
	switch {
	case m.isQuote():
		return m.unmarshalQuote()
	case m.isNull():
		return m.unmarshalNull()
	}
	return nil
}

func (m *Marshaller) unmarshalKey() []byte {
	m.unmarshalNonData()
	s := m.cursor
	for m.cursor < m.len {
		if m.isKeyEnd() {
			m.Inc()
			break
		}
	}
	return m.buffer[s:m.cursor]
}

func (m *Marshaller) unmarshalString() []byte {
	s := m.cursor
	for m.cursor < m.len {
		if m.isValEnd() {
			m.Inc()
			break
		}
	}
	return m.buffer[s:m.cursor]
}

func (m *Marshaller) unmarshalQuote() []byte {
	if m.isQuote() {
		s := m.cursor
		q := m.buffer[s]
		for m.cursor < m.len {
			if m.isEscape() {
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

func (m *Marshaller) unmarshalNull() []byte {
	if m.isNull() {
		m.Inc(len(m.Null))
		return m.Null
	}
	return nil
}

func (m *Marshaller) unmarshalSpace() []byte {
	if m.isSpace() {
		s := m.cursor
		for m.cursor < m.len {
			if !m.isSpace() {
				break
			}
			m.Inc()
		}
		return m.buffer[s:m.cursor]
	}
	return nil
}

func (m *Marshaller) unmarshalCommentBlock() []byte {
	if m.isBlockCommentStart() {
		s := m.cursor
		for m.cursor < m.len {
			if m.isBlockCommentEnd() {
				m.Inc(len(m.BlockCommentEnd))
				break
			}
			m.Inc()
		}
		return m.buffer[s:m.cursor]
	}
	return nil
}

func (m *Marshaller) unmarshalInlineComment() []byte {
	if m.isLineCommentStart() {
		s := m.cursor
		for m.cursor < m.len {
			if m.isLineCommentEnd() {
				m.Inc(len(m.LineCommentEnd))
				break
			}
			m.Inc()
		}
		return m.buffer[s:m.cursor]
	}
	return nil
}

func (m *Marshaller) unmarshalNonData() []byte {
	data, s := false, m.cursor
	for m.cursor < m.len && !data {
		switch {
		case m.isBlockCommentStart():
			m.unmarshalCommentBlock()
		case m.isLineCommentStart():
			m.unmarshalInlineComment()
		case m.isSpace():
			m.unmarshalSpace()
		default:
			data = true
		}
	}
	return m.buffer[s:m.cursor]
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

func (m *Marshaller) unmarshalBreak() []byte {
	if m.isLineBreak() {
		s := m.cursor
		m.Inc(len(m.LineBreak))
		return m.buffer[s:m.cursor]
	}
	return nil
}

func (m *Marshaller) unmarshalIndents() (b []byte, i int) {
	for m.cursor < m.len {
		if m.isIndent() {
			m.Inc(len(m.Indent))
			i++
			continue
		}
		break
	}
	return m.buffer[m.cursor : m.cursor+i], i
}

func (m *Marshaller) unmarshalSliceStart(ancestry []ancestor) (start, delim, end []byte, is bool) {
	if m.hasBrackets {
		if m.isMatch(m.SliceStart) {
			m.Inc(len(m.SliceStart))
			m.IncDepth()
			return m.SliceStart, m.ValEnd, m.SliceEnd, true
		}
	} else {
		if start, delim, end = m.bracketlessSliceComponents(ancestry); m.isMatch(start) {
			m.Inc(len(start))
			m.IncDepth()
			return start, delim, end, true
		}
	}
	if m.isMatch(m.InlineSyntax.SliceStart) {
		m.Inc(len(m.InlineSyntax.SliceStart))
		m.IncDepth()
		return m.InlineSyntax.SliceStart, m.InlineSyntax.ValEnd, m.InlineSyntax.SliceEnd, true
	}
	return nil, nil, nil, false
}

func (m *Marshaller) unmarshalMapStart(ancestry []ancestor) (start, delim, end []byte, is bool) {
	if m.hasBrackets {
		if m.isMatch(m.MapStart) {
			m.Inc(len(m.MapStart))
			m.IncDepth()
			return m.MapStart, m.ValEnd, m.MapEnd, true
		}
	} else {
		if start, delim, end = m.bracketlessMapComponents(ancestry); start != nil {
			if m.isMatch(start) {
				m.Inc(len(start))
				m.IncDepth()
				return start, delim, end, true
			}
		} else if m.isKey() {
			return start, delim, end, true
		}

	}
	if m.isMatch(m.InlineSyntax.MapStart) {
		m.Inc(len(m.InlineSyntax.MapStart))
		m.IncDepth()
		return m.InlineSyntax.MapStart, m.InlineSyntax.ValEnd, m.InlineSyntax.MapEnd, true
	}
	return nil, nil, nil, false
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

func (m *Marshaller) isSpace() bool {
	return InBytes(m.buffer[m.cursor], m.Space)
}

func (m *Marshaller) isIndent() bool {
	return m.isMatch(m.Indent)
}

func (m *Marshaller) isLineBreak() bool {
	return m.isMatch(m.LineBreak)
}

func (m *Marshaller) isQuote() bool {
	return InBytes(m.buffer[m.cursor], m.Quote)
}

func (m *Marshaller) isEscape() bool {
	return InBytes(m.buffer[m.cursor], m.Escape)
}

func (m *Marshaller) isNull() bool {
	return m.isMatch(m.Null)
}

func (m *Marshaller) isValEnd() bool {
	return InBytes(m.buffer[m.cursor], m.ValEnd)
}

func (m *Marshaller) isKey() bool {
	i := m.cursor
	l := len(m.KeyEnd)
	for i < m.len {
		if m.buffer[i] == m.LineBreak[0] {
			return false
		}
		if MatchBytes(m.buffer[i:i+l], m.KeyEnd) {
			return true
		}
		i++
	}
	return false
}

func (m *Marshaller) isKeyEnd() bool {
	return InBytes(m.buffer[m.cursor], m.KeyEnd)
}

func (m *Marshaller) isBlockCommentStart() bool {
	return m.isMatch(m.BlockCommentStart)
}

func (m *Marshaller) isBlockCommentEnd() bool {
	return m.isMatch(m.BlockCommentEnd)
}
func (m *Marshaller) isLineCommentStart() bool {
	return m.isMatch(m.LineCommentStart)
}

func (m *Marshaller) isLineCommentEnd() bool {
	return m.isMatch(m.LineCommentEnd)
}

func (m *Marshaller) isMatch(b []byte) bool {
	return MatchBytes(m.buffer[m.cursor:m.cursor+len(b)], b)
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
