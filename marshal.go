// Copyright 2023 james dotter. All rights reserved.typVal
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import (
	"errors"
)

// ------------------------------------------------------------ /
// Marshaller
// a generic marshaller
// ------------------------------------------------------------ /

type Marshaller struct {
	Type       string
	Cursor     int
	CurIndent  int
	UseIndent  bool
	Value      any
	Len        int
	QuotedKeys bool
	Buffer,
	Space,
	Indent,
	Quote,
	Escape,
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
}

// ------------------------------------------------------------ /
// Presets
// JSON, YAML...
// ------------------------------------------------------------ /

var (
	MarshallerJson = Marshaller{
		Type:              "json",
		Space:             []byte("\t\n\v\f\r "),
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
		Space:            []byte("\t\v\f\r "),
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
		Space:            []byte("\t\v\f\r "),
		Indent:           []byte("  "),
		Quote:            []byte(`"'`),
		Escape:           []byte(`\`),
		ValEnd:           []byte(",\n"),
		KeyEnd:           []byte(":\n"),
		LineCommentStart: []byte("#"),
		LineCommentEnd:   []byte("\n"),
		SliceStart:       []byte("["),
		SliceEnd:         []byte("]"),
		MapStart:         []byte("{"),
		MapEnd:           []byte("}"),
	}
)

func (m *Marshaller) Marshal(a any) ([]byte, error) {
	return m.marshal(ValueOf(a))
}

func (m *Marshaller) marshal(v VALUE) ([]byte, error) {
	switch v.Kind() {
	case Struct:
		return m.MarshalStruct((STRUCT)(v))
	case Array:
		return m.MarshalArray(v)
	case Slice:
		return m.MarshalSlice(v)
	case Map:
		return m.MarshalMap(v)
	case Pointer:
		return m.marshal(v.Elem())
	default:
		return nil, errors.New("cannot marshal type '" + v.typ.String() + "'")
	}
}

func (m *Marshaller) MarshalStruct(s STRUCT, ancestry ...ancestor) ([]byte, error) {
	if s.ptr == nil {
		return []byte("null"), nil
	}
	if s.Len() == 0 {
		return []byte("{}"), nil
	}
	b := []byte{}
	fields := s.TagIndex(m.Type)
	for _, f := range fields {

	}
	return nil, nil
}

func (m *Marshaller) ParseQuote() []byte {
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
	return inBytes(m.Buffer[m.Cursor], m.Space)
}

func (m *Marshaller) IsIndent() bool {
	return matchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.Indent)], m.Indent)
}

func (m *Marshaller) IsQuote() bool {
	return inBytes(m.Buffer[m.Cursor], m.Quote)
}

func (m *Marshaller) IsEscape() bool {
	return inBytes(m.Buffer[m.Cursor], m.Escape)
}

func (m *Marshaller) IsValEnd() bool {
	return inBytes(m.Buffer[m.Cursor], m.ValEnd)
}

func (m *Marshaller) IsKeyEnd() bool {
	return inBytes(m.Buffer[m.Cursor], m.KeyEnd)
}

func (m *Marshaller) IsBlockCommentStart() bool {
	return matchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.BlockCommentStart)], m.BlockCommentStart)
}

func (m *Marshaller) IsBlockCommentEnd() bool {
	return matchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.BlockCommentEnd)], m.BlockCommentEnd)
}
func (m *Marshaller) IsLineCommentStart() bool {
	return matchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.LineCommentStart)], m.LineCommentStart)
}

func (m *Marshaller) IsLineCommentEnd() bool {
	return matchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.LineCommentEnd)], m.LineCommentEnd)
}

func (m *Marshaller) IsSliceStart() bool {
	return matchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.SliceStart)], m.SliceStart)
}

func (m *Marshaller) IsSliceEnd() bool {
	return matchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.SliceEnd)], m.SliceEnd)
}

func (m *Marshaller) IsMapStart() bool {
	return matchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.MapStart)], m.MapStart)
}

func (m *Marshaller) IsMapEnd() bool {
	return matchBytes(m.Buffer[m.Cursor:m.Cursor+len(m.MapEnd)], m.MapEnd)
}

func matchBytes(a, b []byte) bool {
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

func inBytes(a byte, b []byte) bool {
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
