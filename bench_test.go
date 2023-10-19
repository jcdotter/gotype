package gotype

import (
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
				ValueOf(v).Serialize()
			}
		})
	}
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
		{"bytes >> bool", ValueOf(bytes_bool), getrtype(_bool).Kind()},
		{"bytes >> int", ValueOf(bytes_int), getrtype(_int).Kind()},
		{"bytes >> uint", ValueOf(bytes_uint), getrtype(_uint).Kind()},
		{"bytes >> float", ValueOf(bytes_float), getrtype(_float).Kind()},
		{"bytes >> time", ValueOf(bytes_time), getrtype(_time).KIND()},
		{"bytes >> uuid", ValueOf(bytes_uuid), getrtype(_uuid).KIND()},

		{"bool >> bytes", ValueOf(_bool), getrtype(bytes_bool).KIND()},
		{"bool >> int", ValueOf(_bool), getrtype(_int).Kind()},
		{"bool >> uint", ValueOf(_bool), getrtype(_uint).Kind()},
		{"bool >> float", ValueOf(_bool), getrtype(_float).Kind()},
		{"bool >> string", ValueOf(_bool), getrtype(str_bool).Kind()},

		{"int >> bytes", ValueOf(_int), getrtype(bytes_int).KIND()},
		{"int >> bool", ValueOf(_int), getrtype(_bool).Kind()},
		{"int >> uint", ValueOf(_int), getrtype(_uint).Kind()},
		{"int >> float", ValueOf(_int), getrtype(_float).Kind()},
		{"int >> string", ValueOf(_int), getrtype(str_int).Kind()},
		{"int >> time", ValueOf(_int), getrtype(_time).KIND()},

		{"uint >> bytes", ValueOf(_uint), getrtype(bytes_uint).KIND()},
		{"uint >> bool", ValueOf(_uint), getrtype(_bool).Kind()},
		{"uint >> int", ValueOf(_uint), getrtype(_int).Kind()},
		{"uint >> float", ValueOf(_uint), getrtype(_float).Kind()},
		{"uint >> string", ValueOf(_uint), getrtype(str_uint).Kind()},
		{"uint >> time", ValueOf(_uint), getrtype(_time).KIND()},

		{"float >> bytes", ValueOf(_float), getrtype(bytes_float).KIND()},
		{"float >> bool", ValueOf(_float), getrtype(_bool).Kind()},
		{"float >> int", ValueOf(_float), getrtype(_int).Kind()},
		{"float >> uint", ValueOf(_float), getrtype(_uint).Kind()},
		{"float >> string", ValueOf(_float), getrtype(str_float).Kind()},
		{"float >> time", ValueOf(_float), getrtype(_time).KIND()},

		{"string >> bytes", ValueOf(str_bool), getrtype(bytes_string).KIND()},
		{"string >> bool", ValueOf(str_bool), getrtype(_bool).Kind()},
		{"string >> int", ValueOf(str_int), getrtype(_int).Kind()},
		{"string >> uint", ValueOf(str_uint), getrtype(_uint).Kind()},
		{"string >> float", ValueOf(str_float), getrtype(_float).Kind()},
		{"string >> time", ValueOf(str_time), getrtype(_time).KIND()},
		{"string >> uuid", ValueOf(str_uuid), getrtype(_uuid).KIND()},

		{"time >> bytes", ValueOf(_time), getrtype(bytes_time).KIND()},
		{"time >> int", ValueOf(_time), getrtype(_int).Kind()},
		{"time >> uint", ValueOf(_time), getrtype(_uint).Kind()},
		{"time >> float", ValueOf(_time), getrtype(_float).Kind()},
		{"time >> string", ValueOf(_time), getrtype(str_time).Kind()},

		{"uuid >> bytes", ValueOf(_uuid), getrtype(bytes_uuid).KIND()},
		{"uuid >> string", ValueOf(_uuid), getrtype(str_uuid).Kind()},
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
