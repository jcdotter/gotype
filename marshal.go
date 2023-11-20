// Copyright 2023 james dotter. All rights reserved.typVal
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

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
	SliceItem         []byte // the characters before each slice element
	MapStart          []byte // the characters that start a hash map
	MapEnd            []byte // the characters that end a hash map
	InlineSyntax      *InlineSyntax
	// marshalling flags
	Format           bool // when true, marshal with formatting, indentation, and line breaks
	FormatWithSpaces bool // when true, marshal with space between keys and values
	CascadeOnlyDeep  bool // when true, marshal single-depth slices and maps with inline syntax
	QuotedKey        bool // when true, marshal map keys with quotes
	QuotedString     bool // when true, marshal strings with quotes
	QuotedSpecial    bool // when true, marshal strings with quotes if they contain special characters
	QuotedNum        bool // when true, marshal numbers with quotes
	QuotedBool       bool // when true, marshal bools with quotes
	QuotedNull       bool // when true, marshal null with quotes
	RecursiveName    bool // when true, include name, string or type of recursive struct value in marshalling, otherwise, exclude all
	UnmarshalTyped   bool // when true, unmarshal to typed values (int, float64, bool, string) instead of just strings
	MarshalMethods   bool // when true, marshal structs with a Marshal method by calling the method
	IncludeNulls     bool // when true, marshal nil values as null, otherwise, exclude them
	// marshaller cache
	space      []byte
	quote      []byte
	escape     []byte
	valEnd     []byte
	keyEnd     []byte
	ivalEnd    []byte
	sliceParts map[string][3][]byte
	mapParts   map[string][3][]byte
	hasTag     map[*TYPE]bool
	tagKeys    map[*TYPE][]string
	hasMethod  map[*TYPE]bool
	methods    map[*TYPE]reflect.Value
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
	JsonMarshaller = &Marshaller{
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
	YamlMarshaller = &Marshaller{
		Type:             "yaml",
		Format:           true,
		FormatWithSpaces: true,
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

func init() {
	JsonMarshaller.Init()
	YamlMarshaller.Init()
}

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
	m.space = m.Space[:1]
	m.quote = m.Quote[:1]
	m.escape = m.Escape[:1]
	m.keyEnd = m.KeyEnd
	m.valEnd = m.ValEnd
	m.sliceParts = map[string][3][]byte{}
	m.mapParts = map[string][3][]byte{}
	m.initFormat()
}

func (m *Marshaller) initFormat() {
	if m.Format {
		if m.Indent == nil {
			m.Indent = []byte("  ")
		}
		if m.LineBreak == nil {
			m.LineBreak = []byte("\n")
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
		if m.FormatWithSpaces {
			m.keyEnd = AppendBytes(m.KeyEnd, m.space)
			m.valEnd = AppendBytes(m.ValEnd, m.space)
			m.ivalEnd = AppendBytes(m.InlineSyntax.ValEnd, m.space)
		}
	}
}

func (m *Marshaller) New() *Marshaller {
	n := *m
	if m.InlineSyntax != nil {
		s := *m.InlineSyntax
		n.InlineSyntax = &s
	}
	return &n
}

func (m *Marshaller) Reset() {
	m.buffer = []byte{}
	m.cursor = 0
	m.curIndent = 0
	m.len = 0
	m.value = nil
}

func (m *Marshaller) ResetCursor() {
	m.cursor = 0
	m.curDepth = 0
	m.curIndent = 0
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

func (m Marshaller) Bytes() []byte {
	return m.buffer
}

func (m Marshaller) String() string {
	return string(m.buffer)
}

func (m Marshaller) Formatted(indent ...int) string {
	difIndent := false
	oldIndent := m.Indent
	oldBreak := m.LineBreak
	m.initFormat()
	if len(indent) > 0 {
		in := bytes.Repeat([]byte(" "), indent[0])
		if !bytes.Equal(m.Indent, in) {
			m.Indent = in
			difIndent = true
		}
	}
	if (difIndent || m.buffer == nil || !m.Format || !bytes.Equal(oldBreak, m.LineBreak)) && m.value != nil {
		oldFormat := m.Format
		m.Format = true
		m.initFormat()
		m.ResetCursor()
		m.Marshal(m.value)
		m.Indent = oldIndent
		m.Format = oldFormat
	}
	m.LineBreak = oldBreak
	return m.String()
}

func (m Marshaller) Map() map[string]any {
	v := ValueOf(m.value)
	switch v.Kind() {
	case Map:
		return m.value.(map[string]any)
	case Slice:
		return (SLICE)(v).Map()
	}
	return nil
}

func (m Marshaller) Slice() []any {
	v := ValueOf(m.value)
	switch v.Kind() {
	case Map:
		return (MAP)(v).Slice()
	case Slice:
		return m.value.([]any)
	}
	return nil
}

// ------------------------------------------------------------ /
// Marshal Utilities
// methods for marshalling to type
// ------------------------------------------------------------ /

func (m *Marshaller) Marshal(a any) *Marshaller {
	m.Reset()
	m.marshal(ValueOf(a))
	return m
}

func (v VALUE) Marshal(m *Marshaller) *Marshaller {
	m.Reset()
	m.marshal(v)
	return m
}

func (m *Marshaller) marshal(v VALUE, ancestry ...ancestor) {
	if v.IsNil() {
		m.ToBuffer(m.Null)
		return
	}
	switch v.KIND() {
	case Bool:
		m.marshalBool(v.Bool())
	case Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Float32, Float64, Complex64, Complex128:
		m.marshalNum(v)
	case Array:
		m.marshalArray((ARRAY)(v), ancestry...)
	case Func:
		m.marshalFunc(v)
	case Interface:
		m.marshalInterface(v, ancestry...)
	case Map:
		m.marshalMap((MAP)(v), ancestry...)
	case Pointer:
		m.marshalPointer(v, ancestry...)
	case Slice:
		m.marshalSlice((SLICE)(v), ancestry...)
	case String:
		m.marshalString(*(*string)(v.ptr))
	case Struct:
		m.marshalStruct((STRUCT)(v), ancestry...)
	case UnsafePointer:
		m.marshalUnsafePointer(*(*unsafe.Pointer)(v.ptr))
	case Field:
		m.marshalField(*(*FIELD)(unsafe.Pointer(&v)))
	case Time:
		m.marshalTime(*(*TIME)(v.ptr))
	case Uuid:
		m.marshalUuid(*(*UUID)(v.ptr))
	case Bytes:
		m.marshalBytes(*(*BYTES)(v.ptr))
	default:
		panic("cannot marshal type '" + v.typ.String() + "'")
	}
}

func (m *Marshaller) marshalBool(b bool) {
	var bytes []byte
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
}

func (m *Marshaller) marshalNum(v VALUE) {
	var bytes []byte
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
		panic("cannot marshal type '" + v.typ.String() + "'")
	}
	if m.QuotedNum {
		q := m.Quote[:1]
		bytes = append(q, append(bytes, q...)...)
	}
	m.ToBuffer(bytes)
}

func (m *Marshaller) marshalArray(a ARRAY, ancestry ...ancestor) {
	if a.Len() == 0 {
		m.marshalEmptySlice()
		return
	}
	delim, end, ancestry := m.marshalSliceStart((VALUE)(a), ancestry)
	var j int
	a.ForEach(func(i int, k string, v VALUE) (brake bool) {
		j = m.marshalElem(j, delim, nil, v, ancestry)
		return
	})
	m.marshalEnd(end)
}

func (m *Marshaller) marshalFunc(v VALUE) {
	m.ToBuffer([]byte(v.typ.Name()))
}

func (m *Marshaller) marshalInterface(v VALUE, ancestry ...ancestor) {
	v = v.SetType()
	if v.Kind() != Interface {
		m.marshal(v, ancestry...)
		return
	}
	m.Marshal(fmt.Sprint(v.Interface()))
}

func (m *Marshaller) marshalPointer(v VALUE, ancestry ...ancestor) {
	ancestry = append([]ancestor{{v.typ, v.Uintptr()}}, ancestry...)
	m.marshal(v.Elem(), ancestry...)
}

func (m *Marshaller) marshalMap(hm MAP, ancestry ...ancestor) {
	if hm.Len() == 0 {
		m.marshalEmptyMap()
		return
	}
	delim, end, ancestry := m.marshalMapStart((VALUE)(hm), ancestry)
	var j int
	hm.ForEach(func(i int, k string, v VALUE) (brake bool) {
		j = m.marshalElem(j, delim, []byte(k), v, ancestry)
		return
	})
	m.marshalEnd(end)
}

func (m *Marshaller) marshalSlice(s SLICE, ancestry ...ancestor) {
	if s.Len() == 0 {
		m.marshalEmptySlice()
		return
	}
	delim, end, ancestry := m.marshalSliceStart((VALUE)(s), ancestry)
	var j int
	s.ForEach(func(i int, k string, v VALUE) (brake bool) {
		j = m.marshalElem(j, delim, nil, v, ancestry)
		return
	})
	m.marshalEnd(end)
}

func (m *Marshaller) marshalString(s string) {
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
}

func (m *Marshaller) marshalStruct(s STRUCT, ancestry ...ancestor) {
	if m.marshaltStructByMethod(s) {
		return
	}
	if s.Len() == 0 {
		m.marshalEmptyMap()
		return
	}
	delim, end, ancestry := m.marshalMapStart((VALUE)(s), ancestry)
	m.marshalStructTag(s)
	j, has, keys := 0, m.hasTag[s.typ], m.tagKeys[s.typ]
	s.ForEach(func(i int, k string, v VALUE) (brake bool) {
		if has {
			k = keys[i]
		}
		j = m.marshalElem(j, delim, []byte(k), v, ancestry)
		return
	})
	m.marshalEnd(end)
}

func (m *Marshaller) marshaltStructByMethod(s STRUCT) bool {
	if s.typ == TypeOf(TYPE{}) {
		m.marshalString((*TYPE)(s.ptr).Name())
		return true
	}
	if !m.MarshalMethods {
		return false
	}
	if m.hasMethod != nil {
		if has, ok := m.hasMethod[s.typ]; ok {
			if has {
				m.ToBuffer([]byte(m.methods[s.typ].Call([]reflect.Value{})[0].String()))
			}
			return has
		}
	} else {
		m.hasMethod = map[*TYPE]bool{s.typ: false}
		m.methods = map[*TYPE]reflect.Value{}
	}
	n := STRING(m.Type[:1]).ToUpper() + m.Type[1:]
	methods := []string{n, "Marshal" + n}
	for _, name := range methods {
		meth, exists := s.typ.Reflect().MethodByName(name)
		if exists {
			in, out := meth.Type.NumIn(), meth.Type.NumOut()
			if in == 1 && out > 0 {
				if k := FromReflectType(meth.Type.Out(0)).KIND(); k == String || k == Bytes {
					m.hasMethod[s.typ] = true
					m.methods[s.typ] = (VALUE)(s).Reflect().Method(meth.Index)
					b := []byte(m.methods[s.typ].Call([]reflect.Value{})[0].String())
					m.ToBuffer(b)
					return true
				}
			}
		}
	}
	return false
}

func (m *Marshaller) marshalStructTag(s STRUCT) {
	if m.hasTag == nil {
		vals, has := s.typ.TagValues(m.Type)
		m.hasTag = map[*TYPE]bool{s.typ: has}
		m.tagKeys = map[*TYPE][]string{s.typ: vals}
	} else if _, ok := m.hasTag[s.typ]; !ok {
		m.tagKeys[s.typ], m.hasTag[s.typ] = s.typ.TagValues(m.Type)
	}
}

func (m *Marshaller) marshalUnsafePointer(p unsafe.Pointer) {
	m.marshalString(fmt.Sprintf("%p", p))
}

func (m *Marshaller) marshalField(f FIELD) {
	m.marshal(f.VALUE())
}

func (m *Marshaller) marshalTime(t TIME) {
	m.marshalString(t.String())
}

func (m *Marshaller) marshalUuid(u UUID) {
	m.marshalString(u.String())
}

func (m *Marshaller) marshalBytes(b BYTES) {
	m.marshalString(string(b))
}

func (m *Marshaller) marshalSliceComponents(v VALUE, ancestry []ancestor) (start []byte, delim []byte, end []byte) {
	if !m.Format {
		return m.SliceStart, m.ValEnd, m.SliceEnd
	}
	path := m.ancestryPath(ancestry)
	if parts, ok := m.sliceParts[path]; ok {
		return parts[0], parts[1], parts[2]
	}
	switch {
	case m.CascadeOnlyDeep && !v.HasDataElem():
		start, delim, end = m.InlineSyntax.SliceStart, m.ivalEnd, m.InlineSyntax.SliceEnd
	case m.hasBrackets:
		start, delim, end = m.formattedSliceComponents()
	default:
		start, delim, end = m.bracketlessSliceComponents(ancestry)
	}
	m.sliceParts[path] = [3][]byte{start, delim, end}
	return
}

func (m *Marshaller) formattedSliceComponents() (start []byte, delim []byte, end []byte) {
	in := bytes.Repeat(m.Indent, m.curIndent)
	start = AppendBytes(m.SliceStart, m.LineBreak, in, m.Indent, m.SliceItem)
	delim = AppendBytes(m.valEnd, m.LineBreak, in, m.Indent, m.SliceItem)
	end = AppendBytes(m.LineBreak, in, m.SliceEnd)
	return
}

func (m *Marshaller) bracketlessSliceComponents(ancestry []ancestor) (start []byte, delim []byte, end []byte) {
	in := bytes.Repeat(m.Indent, m.curIndent)
	switch m.itemOf(ancestry) {
	case Map:
		start = AppendBytes(m.LineBreak, in, m.Indent, m.SliceItem)
		delim = AppendBytes(m.ValEnd, start)
	case Slice:
		start = m.SliceItem
		delim = AppendBytes(m.ValEnd, m.LineBreak, in, m.Indent, m.Indent, m.SliceItem)
	default:
		start = m.SliceItem
		delim = AppendBytes(m.ValEnd, m.LineBreak, in, m.SliceItem)
	}
	return
}

func (m *Marshaller) marshalMapComponents(v VALUE, ancestry []ancestor) (start []byte, delim []byte, end []byte) {
	if !m.Format {
		return m.MapStart, m.ValEnd, m.MapEnd
	}
	path := m.ancestryPath(ancestry)
	if parts, ok := m.mapParts[path]; ok {
		return parts[0], parts[1], parts[2]
	}
	switch {
	case m.Format && m.CascadeOnlyDeep && !v.HasDataElem():
		start, delim, end = m.InlineSyntax.MapStart, m.ivalEnd, m.InlineSyntax.MapEnd
	case m.hasBrackets:
		start, delim, end = m.formattedMapComponents()
	default:
		start, delim, end = m.bracketlessMapComponents(ancestry)
	}
	m.mapParts[path] = [3][]byte{start, delim, end}
	return
}

func (m *Marshaller) formattedMapComponents() (start []byte, delim []byte, end []byte) {
	in := bytes.Repeat(m.Indent, m.curIndent)
	start = AppendBytes(m.MapStart, m.LineBreak, in, m.Indent)
	delim = AppendBytes(m.valEnd, m.LineBreak, in, m.Indent)
	end = AppendBytes(m.LineBreak, in, m.MapEnd)
	return
}

func (m *Marshaller) bracketlessMapComponents(ancestry []ancestor) (start []byte, delim []byte, end []byte) {
	in := bytes.Repeat(m.Indent, m.curIndent)
	switch m.itemOf(ancestry) {
	case Map:
		start = AppendBytes(m.LineBreak, in, m.Indent)
		delim = AppendBytes(m.ValEnd, m.LineBreak, in, m.Indent)
	case Slice:
		delim = AppendBytes(m.ValEnd, m.LineBreak, in, m.Indent)
	default:
		delim = AppendBytes(m.ValEnd, m.LineBreak, in)
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

func (m *Marshaller) marshalEmptySlice() {
	if !m.hasBrackets {
		m.ToBuffer(m.Null)
		return
	}
	m.SetBuffer(m.buffer, m.SliceStart, m.SliceEnd)
}

func (m *Marshaller) marshalEmptyMap() {
	if !m.hasBrackets {
		m.ToBuffer(m.Null)
		return
	}
	m.SetBuffer(m.buffer, m.MapStart, m.MapEnd)
}

func (m *Marshaller) marshalSliceStart(v VALUE, a []ancestor) (delim []byte, end []byte, ancestry []ancestor) {
	start, delim, end := m.marshalSliceComponents(v, a)
	m.ToBuffer(start)
	m.IncDepth()
	return delim, end, append([]ancestor{{v.typ, v.Uintptr()}}, a...)
}

func (m *Marshaller) marshalMapStart(v VALUE, a []ancestor) (delim []byte, end []byte, ancestry []ancestor) {
	start, delim, end := m.marshalMapComponents(v, a)
	m.ToBuffer(start)
	m.IncDepth()
	return delim, end, append([]ancestor{{v.typ, v.Uintptr()}}, a...)
}

func (m *Marshaller) marshalElem(i int, delim, k []byte, v VALUE, ancestry []ancestor) int {
	if i == 0 {
		delim = nil
	}
	if k != nil {
		var q []byte
		if m.QuotedKey {
			q = m.quote
		}
		k = AppendBytes(q, k, q, m.keyEnd)
	}
	if v.IsNil() {
		if m.IncludeNulls {
			i++
			m.SetBuffer(m.buffer, delim, k, m.Null)
			return i
		}
		return i
	}
	b, recursive := m.recursiveValue(v, ancestry)
	if recursive && b == nil {
		return i
	}
	i++
	m.SetBuffer(m.buffer, delim, k)
	if b != nil {
		m.marshalBytes(b)
	} else {
		m.marshal(v, ancestry...)
	}
	return i
}

func (m *Marshaller) marshalEnd(end []byte) {
	m.ToBuffer(end)
	m.DecDepth()
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
func (m *Marshaller) ToBuffer(b []byte) {
	if b != nil {
		m.buffer = AppendBytes(m.buffer, b)
		m.len = len(m.buffer)
	}
}

// first arg should be the buffer to append,
// or not, to replace the buffer
func (m *Marshaller) SetBuffer(b ...[]byte) {
	m.buffer = AppendBytes(b...)
	m.len = len(m.buffer)
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
		if a.pointer == v.Uintptr() && a.typ == v.typ {
			is = true
			if v.Kind() == Struct && m.RecursiveName {
				if m, ok := v.typ.Reflect().MethodByName("Name"); ok {
					bytes = []byte(v.Reflect().Method(m.Index).Call([]reflect.Value{})[0].String())
				} else if m, ok := v.typ.Reflect().MethodByName("String"); ok {
					bytes = []byte(v.Reflect().Method(m.Index).Call([]reflect.Value{})[0].String())
				} else {
					bytes = []byte(v.typ.NameShort())
				}
			}
			return
		}
	}
	if v.Kind() == Pointer {
		return m.recursiveValue(v.Elem(), ancestry)
	}
	return nil, false
}

func (m *Marshaller) ancestryPath(ancestry []ancestor) (path string) {
	for _, a := range ancestry {
		switch a.typ.Kind() {
		case Map, Struct:
			path += ":Map"
		default:
			path += ":Slice"
		}
	}
	return
}

// ------------------------------------------------------------ /
// Unmarshal Utilities
// methods for unmarshalling from type
// ------------------------------------------------------------ /

func (m *Marshaller) unmarshalError(err string) {
	var start, mid, end int
	switch {
	case m.cursor == 0:
		start = 0
		mid = 0
		end = int(math.Min(float64(m.len-1), 25))
	case m.cursor >= m.len:
		start = int(math.Max(0, float64(m.len-1-25)))
		mid = m.len
		end = m.len
	default:
		start = int(math.Max(0, float64(m.cursor-25)))
		mid = m.cursor
		end = int(math.Min(float64(m.len-1), float64(m.cursor+25)))
	}
	panic("unmarshal error at position " + strconv.Itoa(m.cursor) + " of " + strconv.Itoa(m.len) + "\n" +
		"tried to unmarshal:\n" + string(m.buffer[start:mid]) + "<ERROR>" + string(m.buffer[mid:end]) + "\n" +
		"error: " + err)
}

func (m *Marshaller) Unmarshal(bytes ...[]byte) *Marshaller {
	m.ResetCursor()
	if len(bytes) > 0 {
		m.buffer = bytes[0]
		m.len = len(m.buffer)
	} else if m.value != nil && m.value != any(nil) {
		return m
	}
	m.value = nil
	var slice []any
	var hmap map[string]any
	var value any
	for m.cursor < m.len {
		slice, hmap = m.unmarshalObject()
		if slice != nil {
			m.value = slice
			m.ResetCursor()
			return m
		}
		if hmap != nil {
			m.value = hmap
			m.ResetCursor()
			return m
		}
		value = m.unmarshalItem(nil)
		if value != any(nil) {
			m.value = value
			m.ResetCursor()
			return m
		}
		m.Inc()
	}
	m.ResetCursor()
	return m
}

func (m *Marshaller) unmarshalObject(ancestry ...ancestor) (slice []any, hmap map[string]any) {
	if delim, end, isSlice := m.unmarshalSliceStart(ancestry); isSlice {
		return m.unmarshalSlice(delim, end, ancestry...), nil
	}
	if delim, end, isMap := m.unmarshalMapStart(ancestry); isMap {
		return nil, m.unmarshalMap(delim, end, ancestry...)
	}
	return nil, nil
}

func (m *Marshaller) unmarshalSlice(delim, end []byte, ancestry ...ancestor) (slice []any) {
	ancestry = append([]ancestor{{&TYPE{kind: 23}, 0}}, ancestry...)
	for m.cursor < m.len {
		slice = append(slice, m.unmarshalItem([][]byte{delim, end}, ancestry...))
		m.unmarshalNonData()
		if m.isMatch(delim) {
			m.Inc(len(delim))
			continue
		}
		if m.isMatch(end) || end == nil {
			m.Inc(len(end))
			m.DecDepth()
			break
		}
		m.unmarshalError("failed to find end of slice element")
	}
	return
}

func (m *Marshaller) unmarshalMap(delim, end []byte, ancestry ...ancestor) map[string]any {
	ancestry = append([]ancestor{{&TYPE{kind: 53}, 0}}, ancestry...)
	hmap := map[string]any{}
	for m.cursor < m.len {
		hmap[m.unmarshalKey()] = m.unmarshalItem([][]byte{delim, end}, ancestry...)
		m.unmarshalNonData()
		if m.isMatch(delim) {
			m.Inc(len(delim))
			continue
		} else if m.isMatch(end) || end == nil {
			m.Inc(len(end))
			m.DecDepth()
			break
		}
		m.unmarshalError("failed to find end of map element")
	}
	return hmap
}

func (m *Marshaller) unmarshalItem(endings [][]byte, ancestry ...ancestor) any {
	m.unmarshalNonData()
	switch {
	case m.isQuote():
		return m.unmarshalQuote()
	case m.isNull():
		return m.unmarshalNull()
	default:
		slice, hmap := m.unmarshalObject(ancestry...)
		switch {
		case slice != nil:
			return slice
		case hmap != nil:
			return hmap
		default:
			return m.unmarshalAny(endings...)
		}
	}
}

func (m *Marshaller) unmarshalKey() (key string) {
	m.unmarshalNonData()
	if m.isQuote() {
		key = m.unmarshalQuote()
	}
	s := m.cursor
	for m.cursor < m.len {
		if m.isKeyEnd() {
			m.Inc()
			break
		}
		m.Inc()
	}
	if key == "" {
		key = string(m.buffer[s : m.cursor-1])
	}
	return
}

func (m *Marshaller) unmarshalAny(end ...[]byte) any {
	s := m.cursor
	for m.cursor < m.len {
		if m.isMatch(end...) || m.isMatch(m.LineBreak) {
			break
		}
		m.Inc()
	}
	a := STRING(m.buffer[s:m.cursor]).Trim(string(m.Space))
	if a == "" {
		return nil
	}
	if m.UnmarshalTyped {
		if a == "true" {
			return true
		}
		if a == "false" {
			return false
		}
		if a == string(m.Null) {
			return nil
		}
		if i, e := strconv.ParseInt(a, 10, 64); e == nil {
			return int(i)
		}
		if f, e := strconv.ParseFloat(a, 64); e == nil {
			return f
		}
	}
	return a
}

func (m *Marshaller) unmarshalQuote() string {
	q := m.buffer[m.cursor]
	m.Inc()
	s := m.cursor
	for m.cursor < m.len {
		if m.isEscape() {
			m.Inc(2)
			continue
		}
		if m.ByteIs(q) {
			m.Inc()
			break
		}
		m.Inc()
	}
	return string(m.buffer[s : m.cursor-1])
}

func (m *Marshaller) unmarshalNull() any {
	m.Inc(len(m.Null))
	return nil
}

func (m *Marshaller) unmarshalSpace() string {
	s := m.cursor
	for m.cursor < m.len {
		if !m.isSpace() {
			break
		}
		m.Inc()
	}
	return string(m.buffer[s : m.cursor-1])
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
		return m.buffer[s : m.cursor+len(m.BlockCommentEnd)-1]
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
		return m.buffer[s : m.cursor+len(m.LineCommentEnd)-1]
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

func (m *Marshaller) unmarshalSliceStart(ancestry []ancestor) (delim, end []byte, is bool) {
	if m.hasBrackets {
		if m.isMatch(m.SliceStart) {
			m.Inc(len(m.SliceStart))
			m.IncDepth()
			return m.ValEnd, m.SliceEnd, true
		}
	} else {
		var start []byte
		if start, delim, end = m.bracketlessSliceComponents(ancestry); m.isMatch(start) {
			m.Inc(len(start))
			m.IncDepth()
			return delim, end, true
		}
	}
	if m.Format {
		if m.isMatch(m.InlineSyntax.SliceStart) {
			m.Inc(len(m.InlineSyntax.SliceStart))
			m.IncDepth()
			return m.InlineSyntax.ValEnd, m.InlineSyntax.SliceEnd, true
		}
	}
	return nil, nil, false
}

func (m *Marshaller) unmarshalMapStart(ancestry []ancestor) (delim, end []byte, is bool) {
	if m.hasBrackets {
		if m.isMatch(m.MapStart) {
			m.Inc(len(m.MapStart))
			m.IncDepth()
			return m.ValEnd, m.MapEnd, true
		}
	} else {
		var start []byte
		if start, delim, end = m.bracketlessMapComponents(ancestry); start != nil {
			if m.isMatch(start) {
				m.Inc(len(start))
				m.IncDepth()
				return delim, end, true
			}
		} else if m.isKey() {
			m.IncDepth()
			return delim, end, true
		}
	}
	if m.Format {
		if m.isMatch(m.InlineSyntax.MapStart) {
			m.Inc(len(m.InlineSyntax.MapStart))
			m.IncDepth()
			return m.InlineSyntax.ValEnd, m.InlineSyntax.MapEnd, true
		}
	}
	return nil, nil, false
}

func (m *Marshaller) Inc(i ...int) {
	if i == nil {
		m.cursor++
		return
	}
	m.cursor += i[0]
}

func (m *Marshaller) Byte() byte {
	if m.len == 0 {
		return 0
	}
	return m.buffer[m.cursor]
}

func (m *Marshaller) ByteIs(b byte) bool {
	return m.Byte() == b
}

func (m *Marshaller) isSpace() bool {
	return InBytes(m.buffer[m.cursor], m.Space)
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

func (m *Marshaller) isMatch(b ...[]byte) bool {
	for _, s := range b {
		if m.len-m.cursor < len(s) {
			continue
		}
		if MatchBytes(m.buffer[m.cursor:m.cursor+len(s)], s) {
			return true
		}
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
	return bytes.Equal(a, b)
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
