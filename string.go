// Copyright 2023 escend llc. All rights reserved.typVal
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/google/uuid"
)

type STRING string

// STRING returns gotype VALUE as gotype STRING
func (v VALUE) STRING() STRING {
	return STRING(v.String())
}

// String returns gotype VALUE as golang string
func (v VALUE) String() string {
	switch v.Kind() {
	case String:
		return *(*string)(v.ptr)
	case Pointer:
		return v.Elem().String()
	default:
		return v.string()
	}
}

func (v VALUE) string() string {
	if v.ptr == nil {
		return "null"
	}
	v = v.SetType()
	switch v.KIND() {
	case Bool:
		if *(*byte)(v.ptr) == 0 {
			return "false"
		}
		return "true"
	case Int:
		return strconv.FormatInt(*(*int64)(v.ptr), 10)
	case Int8:
		return strconv.FormatInt(int64(*(*int8)(v.ptr)), 10)
	case Int16:
		return strconv.FormatInt(int64(*(*int16)(v.ptr)), 10)
	case Int32:
		return strconv.FormatInt(int64(*(*int32)(v.ptr)), 10)
	case Int64:
		return strconv.FormatInt(*(*int64)(v.ptr), 10)
	case Uint:
		return strconv.FormatUint(*(*uint64)(v.ptr), 10)
	case Uint8:
		return strconv.FormatUint(uint64(*(*uint8)(v.ptr)), 10)
	case Uint16:
		return strconv.FormatUint(uint64(*(*uint16)(v.ptr)), 10)
	case Uint32:
		return strconv.FormatUint(uint64(*(*uint32)(v.ptr)), 10)
	case Uint64:
		return strconv.FormatUint(*(*uint64)(v.ptr), 10)
	case Uintptr:
		return strconv.FormatUint(*(*uint64)(v.ptr), 10)
	case Float32:
		return strconv.FormatFloat(float64(*(*float32)(v.ptr)), 'f', -1, 64)
	case Float64:
		return strconv.FormatFloat(*(*float64)(v.ptr), 'f', -1, 64)
	case Complex64:
		return strconv.FormatComplex(complex128(*(*complex64)(v.ptr)), 'f', -1, 128)
	case Complex128:
		return strconv.FormatComplex(*(*complex128)(v.ptr), 'f', -1, 128)
	case Array:
		return (ARRAY)(v).String()
	case Chan:
		return "*Channel"
	case Func:
		return `"` + fmt.Sprintf("%v", v.Interface()) + `"`
	case Interface:
		if *(*unsafe.Pointer)(v.ptr) == nil {
			return "null"
		}
		return STRING(fmt.Sprint(v.Interface())).Serialize()
	case Map:
		return (MAP)(v).String()
	case Pointer:
		return v.Elem().String()
	case Slice:
		return (SLICE)(v).String()
	case String:
		return *(*string)(v.ptr)
	case Struct:
		return (STRUCT)(v).String()
	case UnsafePointer:
		return `"` + fmt.Sprint(v.ptr) + `"`
	case Field:
		return (*FIELD)(v.ptr).VALUE().String()
	case Time:
		return (*TIME)(v.ptr).String()
	case Uuid:
		return (*UUID)(v.ptr).String()
	case Bytes:
		return (*BYTES)(v.ptr).String()
	}
	panic("cannot convert value to string")
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Native returns gotype STRING as a golang string
func (s STRING) Native() string {
	return string(s)
}

// Interface returns gotype STRING as a golang interface{}
func (s STRING) Interface() any {
	return s.Native()
}

// VALUE returns gotype SLICE as gotype VALUE
func (s STRING) VALUE() VALUE {
	a := (any)(s)
	return *(*VALUE)(unsafe.Pointer(&a))
}

// Encode returns a gotype encoding of STRING
func (s STRING) Encode() ENCODING {
	return append([]byte{byte(String)}, append(lenBytes(len(s)), s.Bytes()...)...)
	//return append([]byte{byte(String)}, append(s.Bytes(), 0x3)...)
}

// Bytes returns gotype STRING as gotype Bytes
func (s STRING) Bytes() []byte {
	return []byte(s)
}

// Bool returns gotype STRING as Bool
func (s STRING) Bool() bool {
	if s == "true" || s == "True" || s == "TRUE" || s == "1" || s == "t" || s == "T" {
		return true
	} else if s == "" || s == "false" || s == "False" || s == "FALSE" || s == "0" || s == "f" || s == "F" {
		return false
	} else {
		panic(`STRING must be one of "t", "true", "1", "f", "false", "0"`)
	}
}

// BOOL returns gotype STRING as a gotype BOOL
func (s STRING) BOOL() BOOL {
	return BOOL(s.Bool())
}

// Int returns gotype STRING as gotype Int
func (s STRING) Int() int {
	i, e := strconv.ParseInt(string(s), 10, 64)
	if e != nil {
		panic("cannot convert string to int")
	}
	return int(i)
}

// INT returns gotype STRING as a gotype INT
func (s STRING) INT() INT {
	return INT(s.Int())
}

// Uint returns gotype STRING as gotype Uint
func (s STRING) Uint() uint {
	i, e := strconv.ParseUint(string(s), 10, 64)
	if e != nil {
		panic("cannot convert string to uint")
	}
	return uint(i)
}

// UINT returns gotype STRING as a gotype UINT
func (s STRING) UINT() UINT {
	return UINT(s.Uint())
}

// Float64 returns gotype STRING as gotype Float64
func (s STRING) Float64() float64 {
	i, e := strconv.ParseFloat(string(s), 64)
	if e != nil {
		panic("cannot convert string to float64")
	}
	return i
}

// Float returns gotype STRING as a gotype Float
// (parses commas and converts parenthesis to negative num)
func (s STRING) ParseFloat() float64 {
	f, can, _ := s.CanParseFloat()
	if !can {
		panic("cannot parse string as float")
	}
	return f
}

// CanFloat evaluates whether STRING can be parsed to Float
// (parses commas and converts parenthesis to negative num)
// returns gotype STRING as float64
func (s STRING) CanParseFloat() (float float64, can bool, err error) {
	if s == "" {
		return 0, true, nil
	}
	str := strings.ReplaceAll(string(s), ",", "")
	// if first char is '(' and last ')'
	if str[0] == 40 && str[len(str)-1] == 41 {
		str = "-" + string(str[1:len(str)-1])
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, false, err
	}
	return f, true, nil
}

// FLOAT returns gotype STRING as a gotype FLOAT
func (s STRING) FLOAT() FLOAT {
	return FLOAT(s.Float64())
}

// String returns gotype STRING as golang string
func (s STRING) String() string {
	return string(s)
}

// STRING returns gotype STRING as a gotype STRING
func (s STRING) STRING() STRING {
	return s
}

// Serialize returns gotype STRING as serialized string with quotes escaped
func (s STRING) Serialize() string {
	return `"` + s.Escaped(`"`, `\`) + `"`
}

// Time returns gotype STRING as gotype Time
// panics if STRING is not formatted as date or time:
// 2006-01-02 [15:04:05.000]
func (s STRING) TIME() TIME {
	t, can, _ := s.CanTime()
	if !can {
		panic("cannot parse string as time")
	}
	return t
}

// CanTime evaluates whether STRING can convert to Time
// returns gotype STRING as gotype Time
// panics if STRING is not formatted as date or time
func (s STRING) CanTime() (t TIME, can bool, err error) {
	if s == "" {
		return TIME{}, true, nil
	}
	var tt time.Time
	if tt, err = time.Parse(ISO8601, string(s)); err == nil {
		return (TIME)(tt.UTC()), true, nil
	} else if tt, err = time.Parse(SqlDate, string(s)); err == nil {
		return (TIME)(tt.UTC()), true, nil
	} else if tt, err = time.Parse(TimeFormat, string(s)); err == nil {
		return (TIME)(tt.UTC()), true, nil
	} else if tt, err = time.Parse(DateFormat, string(s)); err == nil {
		return (TIME)(tt.UTC()), true, nil
	}
	return (TIME)(tt), false, err
}

// Time returns gotype STRING as gotype Time
// panics if STRING is not formatted as date or time:
// 2006-01-02 [15:04:05.000]
func (s STRING) ParseTime() TIME {
	t, _ := stringToTime(string(s))
	return t
}

var (
	dateparts = map[string]string{
		"sun": "ddd", "mon": "ddd", "tue": "ddd", "wed": "ddd", "thu": "ddd", "fri": "ddd", "sat": "ddd",
		"sunday": "dddd", "monday": "dddd", "tuesday": "dddd", "wednesday": "dddd", "thursday": "dddd", "friday": "dddd", "saturday": "dddd",
		"jan": "mmm", "feb": "mmm", "mar": "mmm", "apr": "mmm", "jun": "mmm", "jul": "mmm", "aug": "mmm", "sep": "mmm", "oct": "mmm", "nov": "mmm", "dec": "mmm",
		"january": "mmmm", "february": "mmmm", "march": "mmmm", "april": "mmmm", "may": "mmmm", "june": "mmmm",
		"july": "mmmm", "august": "mmmm", "september": "mmmm", "october": "mmmm", "november": "mmmm", "december": "mmmm",
		"pm": "pm", "am": "pm",
	}
	dateformats = map[string]string{
		"yy": "06", "yyyy": "2006",
		"mm": "01", "mmm": "Jan", "mmmm": "January",
		"_d": "_2", "dd": "02", "ddd": "Mon", "dddd": "Monday",
		"hh": "15",
		"nn": "04",
		"ss": "05",
		"pm": "PM",
	}
	datemonths = map[string]int{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4, "jun": 6, "jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
		"january": 1, "february": 2, "march": 3, "april": 4, "may": 5, "june": 6,
		"july": 7, "august": 8, "september": 9, "october": 10, "november": 11, "december": 12,
	}
	//datelocations = map[string]
)

func stringToTime(s string) (tm TIME, format string) {
	l, p, t := len(s), "", map[string]*int{"Y": nil, "M": nil, "D": nil, "h": nil, "m": nil, "s": nil, "n": nil}
	loc := time.UTC
	for i := 0; i < l; i++ {
		b := s[i]
		if b == 32 || (43 < b && b < 48) || b == 58 || b == 84 || b == 90 { // is sep
			format += string(b)
		} else {
			p, i = parseDatePart(s, i, l)
			pl := len(p)
			switch {
			case 47 < b && b < 58: // is num
				n := STRING(p).Int()
				switch {
				case pl == 4:
					if t["Y"] == nil {
						t["Y"] = &n
						format += dateformats["yyyy"]
					} else {
						format += p
					}
				case pl == 2:
					switch {
					case t["D"] == nil && t["M"] != nil:
						t["D"] = &n
						format += dateformats["dd"]
					case t["M"] == nil && t["Y"] != nil:
						t["M"] = &n
						format += dateformats["mm"]
					case t["Y"] == nil:
						if n < STRING(Now().Format("06")).Int()+5 {
							n += 2000
						} else {
							n += 1900
						}
						t["Y"] = &n
						format += dateformats["yy"]
					case t["h"] == nil:
						t["h"] = &n
						format += dateformats["hh"]
					case t["m"] == nil:
						t["m"] = &n
						format += dateformats["nn"]
					case t["s"] == nil:
						t["s"] = &n
						format += dateformats["ss"]
					default:
						format += p
					}
				case pl == 1:
					switch {
					case t["D"] == nil:
						t["D"] = &n
						format += dateformats["_d"]
					default:
						format += p
					}
				case pl < 10 && format[len(format)-1] == 46:
					n = n * int(math.Pow10(9-pl))
					t["n"] = &n
					format += strings.Repeat("9", pl)
				default:
					format += p
				}
			case (64 < b && b < 91) || (96 < b && b < 123): // is letter
				plower := STRING(p).ToLower()
				if f, found := dateparts[plower]; found {
					format += dateformats[f]
					if m, found := datemonths[plower]; found {
						t["M"] = &m
					}
				} else {
					format += p
				}
			default:
				format += p
			}
			if i < l {
				format += string(s[i])
			}
		}
	}
	d := 0
	for k, e := range t {
		if e == nil {
			t[k] = &d
		}
	}
	tm = NewTime(*t["Y"], *t["M"], *t["D"], *t["h"], *t["m"], *t["s"], *t["n"], loc)
	return
}

func parseDatePart(s string, i int, l int) (part string, end int) {
	for end = i; end < l; end++ {
		b := s[end]
		if b == 32 || (43 < b && b < 48) || b == 58 || b == 84 || b == 90 {
			break
		}
	}
	part = s[i:end]
	return
}

// UUID returns gotype STRING as gotype UUID
func (s STRING) UUID() UUID {
	if s == "" {
		return UUID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	}
	ui, err := uuid.Parse(string(s))
	if err != nil {
		panic("cannot convert string to UUID")
	}
	return UUID(ui)
}

// Slice returns gotype STRING as []any
func (s STRING) Slice() []any {
	a := make([]any, len(s))
	for i := 0; i < len(s); i++ {
		a[i] = string(s[i])
	}
	return a
}

// SLICE returns gotype STRING as gotype SLICE
func (s STRING) SLICE() SLICE {
	return SliceOf(s.Strings())
}

// Strings returns gotype STRING as []string
func (s STRING) Strings() []string {
	a := make([]string, len(s))
	for i := 0; i < len(s); i++ {
		a[i] = string(s[i])
	}
	return a
}

// JSON returns gotype STRING as gotype JSON
func (s STRING) JSON() JSON {
	j := JSON(s)
	if hasJsonBookends(j) {
		return j
	}
	panic("cannot convert string to JSON")
}

// ------------------------------------------------------------ /
// GOLANG STANDARD IMPLEMENTATIONS
// implementations of functions natively available for
// strings in golang
// referenced packages: strings, regexp
// ------------------------------------------------------------ /

// Len return the number of chars in string
func (s STRING) Len() int {
	return len(s)
}

// ToUpper returns s with all Unicode letters mapped to their upper case
func (s STRING) ToUpper() string {
	return strings.ToUpper(string(s))
}

// ToLower returns s with all Unicode letters mapped to their lower case
func (s STRING) ToLower() string {
	return strings.ToLower(string(s))
}

// Contains evaluates whether string contains substr
func (s STRING) Contains(substr string) bool {
	return strings.Contains(string(s), string(substr))
}

// ContainsAny evaluates whether stromg contains any char in chars
func (s STRING) ContainsAny(chars string) bool {
	return strings.ContainsAny(string(s), string(chars))
}

// ContainsAny evaluates whether string contains any char in chars
func (s STRING) Count(substr string) int {
	return strings.Count(string(s), string(substr))
}

// Index returns the char VALUE at index i of the string
func (s STRING) Index(i int) VALUE {
	v := string(s[i])
	return *(*VALUE)(unsafe.Pointer(&v))
}

// IndexOf returns the index of the first instance of substr
// or -1 if substr is not present in string
func (s STRING) IndexOf(substr string) int {
	return strings.Index(string(s), string(substr))
}

// LastIndexOf returns the index of the last instance of substr
// or -1 if substr is not present in string
func (s STRING) LastIndexOf(substr string) int {
	return strings.LastIndex(string(s), string(substr))
}

// IndexOfAny returns the index of the first instance of any char in chars
// or -1 if substr is not present in string
func (s STRING) IndexOfAny(chars string) int {
	return strings.IndexAny(string(s), string(chars))
}

// LastIndexOfAny returns the index of the last instance of any char in chars
// or -1 if substr is not present in string
func (s STRING) LastIndexOfAny(chars string) int {
	return strings.LastIndexAny(string(s), string(chars))
}

// Replace replaces all matches of old in string with new, up to
// N number of matches in string (or all if N < 0), and
// returns new string containing replacements
func (s STRING) Replace(old string, new string, N int) string {
	return strings.Replace(string(s), old, new, N)
}

// Split slices string into substrings separated by sep, up to
// N number of occurances of sep matched in string (or all if N < 0) and
// returns a slice of the substrings between (and not including) those sep matches
func (s STRING) Split(sep string, N int) []string {
	return strings.SplitN(string(s), sep, N)
}

// Two slices string into two substrings separated by sep
func (s STRING) Two(sep string) (x string, y string) {
	p := strings.Split(string(s), sep)
	x, y = p[0], ""
	if len(p) > 1 {
		y = p[1]
	}
	return
}

// Trim returns a string with all leading and
// trailing chars in cutset removed
func (s STRING) Trim(cutset string) string {
	return strings.Trim(string(s), cutset)
}

// Repeat returns Str releated N times as one string
func (s STRING) Repeat(N int) string {
	return strings.Repeat(string(s), int(N))
}

// RegexMatch evaluates whether string contains any match of the regex
func (s STRING) RegexMatch(regex string) bool {
	re, err := regexp.Compile(regex)
	if err != nil {
		panic("valid regular expression")
	}
	return re.MatchString(string(s))
}

// RegexFind slices string into substrings that match the regex, up to
// N number of matches in string (or all if N < 0), and
// returns a slice of the substrings of those regex matches
func (s STRING) RegexFind(regex string, N int) []string {
	re, err := regexp.Compile(regex)
	if err != nil {
		panic("valid regular expression")
	}
	return re.FindAllString(string(s), N)
}

// RegexReplace replaces all matches of regex in string with repl, up to
// N number of matches in string (or all if N < 0), and
// returns new Str containing replacements
func (s STRING) RegexReplace(regex string, repl string, N int) string {
	re, err := regexp.Compile(regex)
	if err != nil {
		panic("valid regular expression")
	}
	if N < 0 {
		return re.ReplaceAllString(string(s), string(repl))
	}
	var end bool
	var i, occurances int
	for occurances < N && !end {
		loc := re.FindStringIndex(string(s[i:]))
		if loc == nil || loc[1] == len(s)-1 {
			end = true
			continue
		}
		occurances++
		s = s[:loc[0]] + STRING(repl) + s[loc[1]:]
	}
	return string(s)
}

// RegexSplit slices Str into substrings separated by the regex, up to
// N number of occurances of regex matched in string (or all if N < 0) and
// returns a slice of the substrings between (and not including) those regex matches
func (s STRING) RegexSplit(regex string, N int) []string {
	re, err := regexp.Compile(regex)
	if err != nil {
		panic("valid regular expression")
	}
	return re.Split(string(s), N)
}

// Join concatenates the elements of its first argument to create a single string.
// The separator Str sep is placed between elements in the resulting string.
func JoinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

// ------------------------------------------------------------ /
// GOTYPE EXPANDED FUNCTIONS
// implementations of new functions for
// strings in gotype
// referenced packages: strings, regexp
// ------------------------------------------------------------ /

func (s STRING) SetIndex(i int, c string) string {
	n := string(s[:i]) + c + string(s[i+len(c):])
	return n
}

// ToPascal converts STRING example_string
// to pascal case format ExampleString
func (s STRING) ToPascal() string {
	s = STRING(s.toCamelCase())
	return s[:1].ToUpper() + string(s[1:])
}

// ToCamel converts STRING example_string
// to camel case format exampleStr
func (s STRING) ToCamel() string {
	s = STRING(s.toCamelCase())
	return s[:1].ToLower() + string(s[1:])
}

// ToCamel converts STRING example_string
// to camel case format exampleString
func (s STRING) toCamelCase() string {
	np := regexp.MustCompile(`[_ \n\t]`).Split(string(s), -1)
	for i := 1; i < len(np); i++ {
		np[i] = strings.ToUpper(np[i][:1]) + np[i][1:]
	}
	return strings.Join(np, ``)
}

// ToSnake converts STRING exampleStr
// to snake case format example_string
func (s STRING) ToSnake() string {
	return s.toSnakeCase(`_`, true)
}

// ToPhrase converts STRING exampleString
// to phrase case format Example string and
// if case sensative 'c', creating new word at each capital letter
func (s STRING) ToPhrase(c bool) string {
	s = STRING(s.toSnakeCase(` `, c))
	return s[:1].ToUpper() + string(s[1:])
}

// toSnakeCase converts STRING example String
// to snake case format example_string
// using separator 'sep' to join the words in the string and
// if case sensative 'c', creating new word at each capital letter
func (s STRING) toSnakeCase(sep string, c bool) string {
	var r string
	var w bool
	for i := 0; i < len(s); i++ {
		b := []byte{s[i]}
		// check for spacing
		if regexp.MustCompile(`[_ \n\t]`).Match(b) {
			if w {
				r += sep
				w = false
			}
		} else {
			// check for beginning of word if case sensative
			if c && regexp.MustCompile(`[A-Z]`).Match(b) && w {
				r += sep
			}
			r += string(s[i])
			w = true
		}
	}
	if c {
		r = STRING(r).ToLower()
	}
	return r
}

// Width returns Str trucated or expanded with whitespaces
// to meet the N number of chars provided
func (s STRING) Width(N int) string {
	l := N - len(s)
	if l > 0 {
		s += STRING(STRING(" ").Repeat(l))
	}
	if l < 0 {
		s = s[:N]
	}
	return string(s)
}

// RemoveWhitespace removes all spaces and breaks
// in Str not contained within quotes
// returns error if no all quotes are paired
func (s STRING) RemoveWhitespace() (string, error) {
	space := []any{'\t', '\n', '\v', '\f', '\r', ' '}
	quotes := []any{'"', '`', '\''}
	out := rune(0)
	quote := out
	clean := make([]rune, 0, len(s))
	escape := false
	for _, c := range s {
		if !escape {
			if In(c, quotes...) {
				if quote == out {
					quote = c
				} else if quote == c {
					quote = out
				}
			}
		}
		escape = !escape && c == '\\'
		if quote != out || !In(c, space...) {
			clean = append(clean, c)
		}
	}
	if quote != out {
		return "", errors.New("quotes paired in string")
	}
	return string(clean), nil
}

// EscapedSplit splits string 's' into []string
// of substrings between 'sep', unless 'sep' is
// contained within quotes in string 's'
func (s STRING) EscapedSplit(sep string) []string {
	whitespace := "\t\n\v\f\r "
	sp := s.breakQuotedSTRING()
	esp := []string{}
	ps := 0
	for i, s := range sp {
		if i%2 == 0 {
			parts := STRING(s).Split(sep, -1)
			if len(parts) > 1 {
				for p := 0; p < len(parts)-1; p++ {
					part := append(sp[ps:i], parts[p])
					esp = append(esp, STRING(JoinStrings(part, ``)).Trim(whitespace))
					ps = i
				}
				sp[i] = parts[len(parts)-1]
			}
		}
	}
	l := STRING(JoinStrings(sp[ps:cap(sp)], ``)).Trim(whitespace)
	if l != `` {
		esp = append(esp, l)
	}
	return esp
}

// TrimWhitespace removes all spaces, tabs and breaks
// in Str not contained within quotes
// panics if no all quotes are paired
func (s STRING) TrimWhitespace() string {
	whitespace := "\t\n\v\f\r "
	sp := s.breakQuotedSTRING()
	esp := []string{}
	for i, s := range sp {
		if i%2 == 0 {
			for i, w := range whitespace {
				if i < len(whitespace)-1 {
					s = STRING(s).Replace(string(w), " ", -1)
				} else {
					parts := STRING(s).Split(string(w), -1)
					subp := []string{}
					for _, p := range parts {
						if p != "" {
							subp = append(subp, p)
						}
					}
					s = JoinStrings(subp, ` `)
				}
			}
		}
		esp = append(esp, s)
	}
	return STRING(JoinStrings(esp, ` `)).Trim(whitespace)
}

func (s STRING) breakQuotedSTRING() []string {
	parts := s.breakQuotes()
	sp := []string{}
	for _, part := range parts {
		sp = append(sp, string(part))
	}
	return sp
}

func (s STRING) breakQuotes() [][]rune {
	quotes := []any{'"', '`', '\''}
	parts := [][]rune{{}}
	out := rune(0)
	quote := out
	escape := false
	part := 0
	for _, c := range s {
		i := false
		if !escape {
			if In(c, quotes...) {
				if quote == out {
					quote = c
				} else if quote == c {
					quote = out
					parts[part] = append(parts[part], c)
					i = true
				}
				parts = append(parts, []rune{})
				part++
			}
		}
		escape = !escape && c == '\\'
		if !i {
			parts[part] = append(parts[part], c)
		}
	}
	if quote != out {
		panic("quotes not paired in string")
	}
	return parts
}

// Escaped returns string with quote chars escaped
func (s STRING) Escaped(quote string, esc string) string {
	p, q, e := byte(0), quote[0], esc[0]
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b == q && p != e {
			s = s[:i] + STRING(e) + s[i:]
			i++
		}
		p = b
	}
	return string(s)
}

// StrFormat represents the case format of a string
// Pascal:	ExampleStr
// Camel: 	exampleStr
// Snake: 	example_string
// Phrase: 	Example string
type StrFormat uint8

const (
	None StrFormat = iota
	Pascal
	Camel
	Snake
	Phrase
	Upper
	Lower
)

// Format converts string 's' to the StrFormat
func (f StrFormat) Format(s string) string {
	switch f {
	case Pascal:
		return STRING(s).ToPascal()
	case Camel:
		return STRING(s).ToCamel()
	case Snake:
		return STRING(s).ToSnake()
	case Phrase:
		return STRING(s).ToPhrase(true)
	case Upper:
		return STRING(s).ToUpper()
	case Lower:
		return STRING(s).ToLower()
	default:
		return s
	}
}

func (s STRING) UnserializeMap(kvSep string, itemSep string) map[string]string {
	return unserializeMap(string(s), CharToByte(itemSep), CharToByte(kvSep))
}

func (s STRING) UnserializeList(itemSep string) []string {
	return unserializeList(string(s), CharToByte(itemSep))
}

func (s STRING) GetSerialValue(key string, kvSep string, itemSep string) (string, bool) {
	return getSerialValue(key, string(s), CharToByte(itemSep), CharToByte(kvSep))
}

func CharToByte(sep string) byte {
	if sep == "" {
		return byte(0)
	} else {
		return sep[0]
	}
}

func unserializeMap(s string, elSep byte, kvSep byte) map[string]string {
	l := len(s)
	m := map[string]string{}
	var key, value string
	var inKey, inValue bool
	inKey = true
	for i := 0; i < l; i++ {
		if inKey {
			key, i = parseSerialItem(s, i, l, kvSep)
			m[key], inKey, inValue = "", false, true
		} else if inValue {
			value, i = parseSerialItem(s, i, l, elSep)
			m[key], inKey, inValue = value, true, false
			key, value = "", ""
		}
	}
	return m
}

func unserializeList(s string, elSep byte) []string {
	l := len(s)
	m := []string{}
	var value string
	var item int
	for i := 0; i < l; i++ {
		value, i = parseSerialItem(s, i, l, elSep)
		m = append(m, value)
		item++
	}
	return m
}

func getSerialValue(s string, k string, elSep byte, kvSep byte) (string, bool) {
	l := len(s)
	var key, value string
	var inKey, inValue bool
	inKey = true
	for i := 0; i < l; i++ {
		if inKey {
			key, i = parseSerialItem(s, i, l, kvSep)
			inKey, inValue = false, true
		} else if inValue {
			value, i = parseSerialItem(s, i, l, elSep)
			if key == k {
				return value, true
			}
			inKey, inValue = true, false
			key, value = "", ""
		}
	}
	return "", false
}

func parseSerialItem(s string, start int, l int, sep byte) (item string, end int) {
	start = skipWhitespace(s, start, l)
	if q := s[start]; IsQuote(q) {
		start++
		item, end = parseQuote(s, start, l, q)
		end++
		if end == l || s[end] == sep {
			return
		}
		end = skipWhitespace(s, end, l)
		if end == l || s[end] == sep {
			return
		}
		panic("misformatted serial string")
	}
	for end = start; end < l; end++ {
		if s[end] == sep {
			item = s[start:end]
			return
		}
	}
	item = s[start:end]
	return
}

func parseQuote(s string, start int, l int, q byte) (quote string, end int) {
	for end = start; end < l; end++ {
		if s[end] == q {
			quote = s[start:end]
			return
		}
	}
	panic("unpaired quote in string")
}

func skipWhitespace(s string, start int, l int) (end int) {
	for end = start; end < l; end++ {
		if !IsWhitespace(s[end]) {
			return
		}
	}
	return
}

func IsWhitespace(b byte) bool {
	return b == 32 || b < 14
}

func IsQuote(b byte) bool {
	return b == 34 || b == 39 || b == 36
}

// InStrings evaluates if string is list of strings
func (s STRING) InStrings(xs []string) bool {
	for _, x := range xs {
		if x == string(s) {
			return true
		}
	}
	return false
}

// StringsInStrings evaluates whether all x's exist in y's
func StringsInStrings(xs []string, ys []string) bool {
	for _, x := range xs {
		if !STRING(x).InStrings(ys) {
			return false
		}
	}
	return true
}

// CoalesceStrings concatenates a series of strings into one string
func CoalesceStrings(strings ...string) (s string) {
	for _, s = range strings {
		if s != "" {
			return
		}
	}
	return
}

// Unserialize converts a json serialized STRING to a map or slice, respectively
func (s STRING) Unserialize() (object map[string]any, list []any) {
	l := len(s)
	start := s.skipWhitespace(0, l)
	switch s[start] {
	case `{`[0]:
		object, _ = s.parseObject(start+1, l, `}`[0], `,`[0], `:`[0])
		return
	case `[`[0]:
		list, _ = s.parseList(start+1, l, `]`[0], `,`[0])
		return
	default:
		panic("malformatted serialized string")
	}
}

func (s STRING) parseObject(start int, l int, objEnd byte, elSep byte, kvSep byte) (hmap map[string]any, end int) {
	start, hmap = s.skipWhitespace(start, l), map[string]any{}
	var key, value any
	var inKey, inValue bool
	inKey, end = true, start
	for end < l {
		if inKey {
			end = s.skipWhitespace(end, l)
			key, end = s.parseSerialItem(end, l, objEnd, kvSep)
			hmap[key.(string)], inKey = nil, false
		} else {
			panic("malformatted serialized string")
		}
		if s[end] == kvSep {
			end++
			inValue = true
		} else {
			panic("malformatted serialized string")
		}
		if inValue {
			value, end = s.parseSerialItem(end, l, objEnd, elSep)
			hmap[key.(string)], inValue = value, false
			key, value = "", ""
		} else {
			panic("malformatted serialized string")
		}
		if s[end] == elSep {
			end++
			inKey = true
			continue
		}
		if s[end] == objEnd {
			end++
			return
		} else {
			panic("malformatted serialized string")
		}
	}
	return
}

func (s STRING) parseList(start int, l int, listEnd byte, elSep byte) (slice []any, end int) {
	start, slice = s.skipWhitespace(start, l), []any{}
	var value any
	end = start
	for end < l {
		value, end = s.parseSerialItem(end, l, listEnd, elSep)
		slice = append(slice, value)
		if s[end] == elSep {
			end++
			continue
		}
		if s[end] == listEnd {
			end++
			return
		} else {
			panic("malformatted serialized string")
		}
	}
	return

}

func (s STRING) parseQuote(start int, l int, q byte) (quote string, end int) {
	var skip int
	for end = start; end < l; end++ {
		if s[end] == q {
			quote = string(s[start:end])
			end = end + skip + 1
			return
		}
		if s[end] == `\`[0] {
			if s[end+1] == q {
				l--
				skip++
				s = s[:end] + s[end+1:]
			}
		}
	}
	panic("unpaired quote in string")
}

func (s STRING) parseSerialItem(start int, l int, elend byte, sep byte) (item any, end int) {
	start = s.skipWhitespace(start, l)
	if q := s[start]; IsQuote(q) {
		start++
		item, end = s.parseQuote(start, l, q)
	} else if q == `{`[0] {
		item, end = s.parseObject(start+1, l, `}`[0], sep, `:`[0])
		return
	} else if q == `[`[0] {
		item, end = s.parseList(start+1, l, `]`[0], sep)
		return
	} else {
		for end = start; end < l; end++ {
			if s[end] == sep || s[end] == elend {
				item = s[start:end].typeSerialValue()
				return
			}
		}
		item = s[start:end].typeSerialValue()
		return
	}
	end = s.skipWhitespace(end, l)
	return
}

func (s STRING) skipWhitespace(start int, l int) (end int) {
	for end = start; end < l; end++ {
		if !IsWhitespace(s[end]) {
			return
		}
	}
	return
}

func (s STRING) skipWhitespaceTail(start int, l int) (end int) {
	for end = l - 1; end >= start; end-- {
		if !IsWhitespace(s[end]) {
			end++
			return
		}
	}
	return
}

func (s STRING) typeSerialValue() any {
	s = s[:s.skipWhitespaceTail(0, len(s))]
	if s == "true" || s == "false" {
		return s.Bool()
	}
	if s == "null" {
		return nil
	}
	if i, e := strconv.ParseInt(string(s), 10, 64); e == nil {
		return int(i)
	}
	if f, can, _ := s.CanParseFloat(); can {
		return f
	}
	return s
}
