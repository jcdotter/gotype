// Copyright 2023 james dotter. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func BenchmarkValueOf(b *testing.B) {
	for n, v := range getTestVars() {
		b.Run(STRING(n).Width(35), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ValueOf(v)
			}
		})
	}
}

func BenchmarkSerialize(b *testing.B) {
	for n, v := range getTestVars() {
		b.Run(STRING(n).Width(35), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ValueOf(v).json()
			}
		})
	}
}

func BenchmarkMarshaller(b *testing.B) {
	var s string
	var j []byte
	for _, v := range getTestVarsGmap() {
		b.Run(STRING(v.Key).Width(35)+"-Encode", func(b *testing.B) {
			val := v.Value.Interface()
			for i := 0; i < b.N; i++ {
				//json.Marshal(val)
				j, _ = json.Marshal(val)
				var b bytes.Buffer
				json.Indent(&b, j, "", "  ")
			}
		})
		b.Run(STRING(v.Key).Width(35)+"-Serial", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				//v.Value.json()
				s = v.Value.json()
			}
		})
		b.Run(STRING(v.Key).Width(35)+"-Marshl", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v.Value.Marshal(JsonMarshaller)
			}
		})
		p := v.Value.typ.PtrType().New().Interface()
		b.Run(STRING(v.Key).Width(35)+"-Decode", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				json.Unmarshal(j, p)
			}
		})
		b.Run(STRING(v.Key).Width(35)+"-DeSerial", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				STRING(s).UnmarshalJson()
			}
		})
		b.Run(STRING(v.Key).Width(35)+"-UnMarshl", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				JsonMarshaller.Unmarshal()
			}
		})
	}
}

func BenchmarkAppend(b *testing.B) {
	a := []byte("") //make([]byte, 0, 10000)
	s := []byte(" ")
	b.Run("append to cap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a = append(a, s...)
		}
	})
	a = s
	b.Run("append", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a = append(append(append(a, s...), s...), s...)
		}
	})
	a = s
	b.Run("append bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			//a = AppendBytes(a, s)
			a = append(a, s...)
		}
	})
}

func BenchmarkEncode(b *testing.B) {
	for n, v := range getTestVars() {
		b.Run(STRING("Encode("+n+")").Width(42), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Encode(v)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	for n, v := range getTestVars() {
		e := Encode(v)
		b.Run(STRING("Decode("+n+")").Width(42), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				e.Decode(&v)
			}
		})
	}
}

func BenchmarkCast(b *testing.B) {
	type conv struct {
		name string
		val  VALUE
		knd  KIND
	}
	var (
		_bool  = true
		_int   = 123
		_uint  = uint(123)
		_float = 123.0
		_time  = TIME(time.Unix(0, 123).UTC())
		_uuid  = UUID([16]byte{0x26, 0x7b, 0x32, 0x29, 0x25, 0x66, 0x44, 0x26, 0xa8, 0x26, 0x8d, 0x80, 0x12, 0x6e, 0x71, 0x9a})

		str_bool  = "true"
		str_int   = "123"
		str_uint  = "123"
		str_float = "123.0"
		str_time  = `1970-01-01 00:00:00.000000123`
		str_uuid  = `267b3229-2566-4426-a826-8d80126e719a`

		bytes_bool   = []byte{0x1}
		bytes_int    = []byte{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		bytes_uint   = []byte{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		bytes_float  = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0xc0, 0x5e, 0x40}
		bytes_time   = []byte{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		bytes_uuid   = []byte{0x26, 0x7b, 0x32, 0x29, 0x25, 0x66, 0x44, 0x26, 0xa8, 0x26, 0x8d, 0x80, 0x12, 0x6e, 0x71, 0x9a}
		bytes_string = []byte("true")
	)
	vars := []conv{
		{"bytes >> bool", ValueOf(bytes_bool), TypeOf(_bool).Kind()},
		{"bytes >> int", ValueOf(bytes_int), TypeOf(_int).Kind()},
		{"bytes >> uint", ValueOf(bytes_uint), TypeOf(_uint).Kind()},
		{"bytes >> float", ValueOf(bytes_float), TypeOf(_float).Kind()},
		{"bytes >> time", ValueOf(bytes_time), TypeOf(_time).KIND()},
		{"bytes >> uuid", ValueOf(bytes_uuid), TypeOf(_uuid).KIND()},

		{"bool >> bytes", ValueOf(_bool), TypeOf(bytes_bool).KIND()},
		{"bool >> int", ValueOf(_bool), TypeOf(_int).Kind()},
		{"bool >> uint", ValueOf(_bool), TypeOf(_uint).Kind()},
		{"bool >> float", ValueOf(_bool), TypeOf(_float).Kind()},
		{"bool >> string", ValueOf(_bool), TypeOf(str_bool).Kind()},

		{"int >> bytes", ValueOf(_int), TypeOf(bytes_int).KIND()},
		{"int >> bool", ValueOf(_int), TypeOf(_bool).Kind()},
		{"int >> uint", ValueOf(_int), TypeOf(_uint).Kind()},
		{"int >> float", ValueOf(_int), TypeOf(_float).Kind()},
		{"int >> string", ValueOf(_int), TypeOf(str_int).Kind()},
		{"int >> time", ValueOf(_int), TypeOf(_time).KIND()},

		{"uint >> bytes", ValueOf(_uint), TypeOf(bytes_uint).KIND()},
		{"uint >> bool", ValueOf(_uint), TypeOf(_bool).Kind()},
		{"uint >> int", ValueOf(_uint), TypeOf(_int).Kind()},
		{"uint >> float", ValueOf(_uint), TypeOf(_float).Kind()},
		{"uint >> string", ValueOf(_uint), TypeOf(str_uint).Kind()},
		{"uint >> time", ValueOf(_uint), TypeOf(_time).KIND()},

		{"float >> bytes", ValueOf(_float), TypeOf(bytes_float).KIND()},
		{"float >> bool", ValueOf(_float), TypeOf(_bool).Kind()},
		{"float >> int", ValueOf(_float), TypeOf(_int).Kind()},
		{"float >> uint", ValueOf(_float), TypeOf(_uint).Kind()},
		{"float >> string", ValueOf(_float), TypeOf(str_float).Kind()},
		{"float >> time", ValueOf(_float), TypeOf(_time).KIND()},

		{"string >> bytes", ValueOf(str_bool), TypeOf(bytes_string).KIND()},
		{"string >> bool", ValueOf(str_bool), TypeOf(_bool).Kind()},
		{"string >> int", ValueOf(str_int), TypeOf(_int).Kind()},
		{"string >> uint", ValueOf(str_uint), TypeOf(_uint).Kind()},
		{"string >> float", ValueOf(str_float), TypeOf(_float).Kind()},
		{"string >> time", ValueOf(str_time), TypeOf(_time).KIND()},
		{"string >> uuid", ValueOf(str_uuid), TypeOf(_uuid).KIND()},

		{"time >> bytes", ValueOf(_time), TypeOf(bytes_time).KIND()},
		{"time >> int", ValueOf(_time), TypeOf(_int).Kind()},
		{"time >> uint", ValueOf(_time), TypeOf(_uint).Kind()},
		{"time >> float", ValueOf(_time), TypeOf(_float).Kind()},
		{"time >> string", ValueOf(_time), TypeOf(str_time).Kind()},

		{"uuid >> bytes", ValueOf(_uuid), TypeOf(bytes_uuid).KIND()},
		{"uuid >> string", ValueOf(_uuid), TypeOf(str_uuid).Kind()},
	}

	for _, v := range vars {
		b.Run(STRING(v.name).Width(24), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v.val.Cast(v.knd)
			}
		})
	}
}

// Benchmark Cast
