// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package vals

import (
	"fmt"
	"math"
	"testing"
)

const (
	print_test_on = true
	print_fail_on = true
	print_tsrc_on = false
	print_full_v  = false
)

type validation_test struct {
	Name  string
	Value any
	Valid any
}

func TestNone(t *testing.T) {

}

func TestPrintChars(t *testing.T) {
	for i := 32; i < 127; i++ {
		fmt.Print(i, ": ", testSpace(i), string(byte(i)), "\t")
		if math.Mod(float64(i-31), 10) == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func TestPrintMax(t *testing.T) {
	l := [][]string{
		{"int8:\t", INT(math.MaxInt8).String()},
		{"int16:\t", INT(math.MaxInt16).String()},
		{"int32:\t", INT(math.MaxInt32).String()},
		{"int64:\t", INT(math.MaxInt64).String()},
		{"float32:", FLOAT(math.MaxFloat32).String()},
		{"float64:", FLOAT(math.MaxFloat64).String()},
	}
	for _, m := range l {
		fmt.Println(m[0], "\t", m[1])
	}
	fmt.Println()
}

func TestPrintHexChart(t *testing.T) {
	hextable := "0123456789abcdef"
	for i := 0; i < 256; i++ {
		if i > 0 && math.Mod(float64(i), 16) == 0 {
			fmt.Println()
		}
		fmt.Print("b: ", testSpace(i), i, " hex: ", string([]byte{hextable[i>>4], hextable[i&0x0f]}), "  |  ")
	}
}

func testSpace(i int) string {
	s := ""
	if i < 10 {
		s = "  "
	} else if i < 100 {
		s = " "
	}
	return s
}

func print_test(src bool, name string, inPtr any, outPtr any, expPtr any) {
	if !print_test_on {
		return
	}
	fmt.Println()
	fmt.Println("TEST:", name)
	if src {
		fmt.Printf("SRC-- %v\n", Source(2))
	}
	if inPtr != nil {
		print_test_elem("IN-- ", inPtr)
	}
	if outPtr != nil {
		print_test_elem("OUT--", outPtr)
	}
	if expPtr != nil {
		print_test_elem("EXP--", expPtr)
	}
}

func print_test_elem(el string, v any) {
	rv := ValueOf(v)
	typ := STRING(rv.typ.Name()).Width(10)
	bas := rv.Kind().STRING().Width(12)
	val := rv.STRING()
	fvl := STRING(fmt.Sprintf("%#v", rv.Interface()))
	if !print_full_v {
		if len(val) > 25 {
			val = STRING(val.Width(25) + " ... ")
		} else {
			val = STRING(val.Width(30))
		}
		fln := int(math.Min(45, float64(len(fvl))))
		if len(fvl) > fln {
			fvl = STRING(fvl.Width(fln) + " ... ")
		} else {
			fvl = STRING(fvl.Width(fln))
		}
	}
	fmt.Printf("%s  TYPE: %v  BASE: %v  VAL: %v  (%v)\n", el, typ, bas, val, fvl)
}

func test_validations(tests []validation_test) error {
	for _, test := range tests {
		if print_test_on {
			print_test(print_tsrc_on, test.Name, nil, test.Value, test.Valid)
		}
		if !DeepEqual(test.Value, test.Valid) {
			if print_fail_on && !print_test_on {
				print_test(print_tsrc_on, test.Name, nil, test.Value, test.Valid)
			}
			return fmt.Errorf("result match error validation %v IS NOT %v", test.Value, test.Valid)
		}
	}
	return nil
}
