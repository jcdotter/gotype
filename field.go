// Copyright 2023 james dotter. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import "unsafe"

type FIELD struct {
	typ    *rtype
	ptr    unsafe.Pointer
	f      flag
	name_  name
	name   string
	rawtag string
	index  int
}

// Interface returns the Field value as interface{}
func (f FIELD) Interface() any {
	return f.VALUE().Interface()
}

// Value returns the Field value as Value
func (f FIELD) VALUE() VALUE {
	return *(*VALUE)(unsafe.Pointer(&f))
}

// Kind returns the gotype Kind of the field
func (f FIELD) KIND() KIND {
	return f.typ.KIND()
}

// Index returns the position of FIELD in STRUCT
func (f FIELD) Index() int {
	return f.index
}

// Name returns the name of the Field
func (f FIELD) Name() string {
	return f.name
}

// String returns the string value of the Field
func (f FIELD) String() string {
	return f.VALUE().String()
}

// STRING returns the STRING value of the Field
func (f FIELD) STRING() STRING {
	return f.VALUE().STRING()
}

func (f FIELD) RawTag() string {
	return f.rawtag
}

// Tags returns a map of key vaue pairs in the field's tag
func (f FIELD) Tags() map[string]string {
	if f.rawtag != "" {
		return parseTags(f.rawtag, `"`[0])
	}
	return map[string]string{}
}

// Tag returns the value of the given tag name
func (f FIELD) Tag(name string) string {
	if f.rawtag == "" {
		f.rawtag = f.name_.tag()
	}
	return getTagValue(f.rawtag, name, `"`[0])
}

// SubTags returns a map of key vaue pairs in a given tag in the field
func (f FIELD) SubTags(tag string) map[string]string {
	return parseTags(f.Tag(tag), `'`[0])
}

// SubTag returns the value of a given subTag in the given tag
// SubTag("db","name") returns "colname" from tag `db:"name:'colname'"`
func (f FIELD) SubTag(tag string, subTag string) string {
	return getTagValue(f.Tag(tag), subTag, `'`[0])
}

func (f FIELD) Set(a any) FIELD {
	f.ptr = f.VALUE().Set(a).ptr
	return f
}

func (f FIELD) Visible() bool {
	if c := f.name[0]; c > 64 {
		return c < 91
	}
	return false
}

func parseTags(rawtag string, q byte) map[string]string {
	l := len(rawtag)
	tags := map[string]string{}
	var tag, value string
	var inTag, inValue bool
	inTag = true
	for i := 0; i < l; i++ {
		if inTag {
			tag, i, inTag, inValue = parseTagName(rawtag, i, l)
			tags[tag] = ""
		} else if inValue {
			value, i, inTag, inValue = parseTagValue(rawtag, i, l, q)
			tags[tag] = value
			tag, value = "", ""
		}
	}
	return tags
}

func getTagValue(rawtag string, tagname string, q byte) string {
	l := len(rawtag)
	var tag, value string
	var inTag, inValue bool
	inTag = true
	for i := 0; i < l; i++ {
		if inTag {
			tag, i, inTag, inValue = parseTagName(rawtag, i, l)
		} else if inValue {
			value, i, inTag, inValue = parseTagValue(rawtag, i, l, q)
			if tag == tagname {
				return value
			}
			tag, value = "", ""
		}
	}
	return ""
}

func parseTagName(t string, start int, l int) (tag string, end int, inTag bool, inValue bool) {
	if t[start] == 32 {
		panic("misformatted struct tag")
	}
	for end = start; end < l; end++ {
		if t[end] == 58 {
			inValue, tag = true, t[start:end]
			return
		}
	}
	tag = t[start:l]
	return
}

func parseTagValue(t string, start int, l int, q byte) (value string, end int, inTag bool, inValue bool) {
	if t[start] != q {
		panic("misformatted struct tag")
	}
	start++
	for end = start; end < l; end++ {
		if t[end] == q {
			inTag, value = true, t[start:end]
			if end < l-1 {
				end++
				if t[end] != 32 {
					panic("misformatted struct tag")
				}
			}
			return
		}
	}
	value = t[start:end]
	return
}
