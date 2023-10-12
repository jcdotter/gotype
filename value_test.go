package gotype

import (
	"fmt"
	"reflect"
	"testing"

	test "github.com/jcdotter/gtest"
)

var config = &test.Config{
	PrintTest:   true,
	PrintTrace:  true,
	PrintDetail: true,
	FailFatal:   true,
	Msg:         "%s",
}

func TestTest(t *testing.T) {
	s := ValueOf([]string{}).New().Interface().([]string)
	fmt.Printf("%#v\n", s)
	s = append(s, "")
	fmt.Printf("%#v\n", s)
}

func TestValueOf(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Interface()"
	for n, v := range createTestVars() {
		gt.Equal(ValueOf(v).Interface(), v, n)
	}
}

func TestZero(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing Zero(%s)"
	for n, v := range createTestVars() {
		r := reflect.ValueOf(v)
		r.SetZero()
		fmt.Printf("%v: %v\n", n, r.Elem())
	}
}

func TestValueNew(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing VALUE.New(%s)"
	var (
		b  = true
		i  = 1
		s  = "test"
		a  = [2]string{"s", "s"}
		a1 = [1]string{"s"}
		l  = []string{"s", "s"}
		m  = map[string]string{"0": "s", "1": "s"}
		d  = string_struct{"s", "s"}
		d1 = string_struct_single{"s"}
	)

	gt.Equal(ValueOf(b).New().Interface(), false, "bool")
	gt.Equal(ValueOf(i).New().Interface(), 0, "int")
	gt.Equal(ValueOf(s).New().Interface(), "", "string")
	gt.Equal(ValueOf(a).New().Interface(), [2]string{"", ""}, "array")
	gt.Equal(ValueOf(l).New().Interface(), []string(nil), "slice")
	gt.Equal(ValueOf(m).New().Interface(), map[string]string(nil), "map")
	gt.Equal(ValueOf(d).New().Interface(), string_struct{}, "struct")
	gt.Equal(ValueOf(a1).New().Interface(), [1]string{""}, "array(1)")
	gt.Equal(ValueOf(d1).New().Interface(), string_struct_single{}, "struct(1)")

	gt.Equal(ValueOf(&b).New().Elem().Set(&b).Interface(), nil, "bool")
	/* gt.Equal(ValueOf(i).New().Interface(), 0, "int")
	gt.Equal(ValueOf(s).New().Interface(), "", "string")
	gt.Equal(ValueOf(a).New().Interface(), [2]string{"", ""}, "array")
	gt.Equal(ValueOf(l).New().Interface(), []string(nil), "slice")
	gt.Equal(ValueOf(m).New().Interface(), map[string]string(nil), "map")
	gt.Equal(ValueOf(d).New().Interface(), string_struct{}, "struct")
	gt.Equal(ValueOf(a1).New().Interface(), [1]string{""}, "array len 1")
	gt.Equal(ValueOf(d1).New().Interface(), string_struct_single{}, "struct len 1") */

}

func TestValueLen(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing %s.Len()"
	var (
		s  = "test"
		a  = [2]string{"s", "s"}
		a1 = [1]string{"s"}
		l  = []string{"s", "s"}
		m  = map[string]string{"0": "s", "1": "s"}
		d  = string_struct{"s", "s"}
		d1 = string_struct_single{"s"}
	)
	gt.Equal(4, ValueOf(s).Len(), "string")
	gt.Equal(2, ValueOf(a).Len(), "array")
	gt.Equal(2, ValueOf(l).Len(), "slice")
	gt.Equal(2, ValueOf(m).Len(), "map")
	gt.Equal(2, ValueOf(d).Len(), "struct")
	gt.Equal(1, ValueOf(a1).Len(), "array(1)")
	gt.Equal(1, ValueOf(d1).Len(), "struct(1)")
}

func TestValueIndex(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Index(%d)"
	var (
		s  = "test"
		a  = [2]string{"s", "s"}
		a1 = [1]string{"s"}
		l  = []string{"s", "s"}
		m  = map[string]string{"0": "s", "1": "s"}
		d  = string_struct{"s", "s"}
		d1 = string_struct_single{"s"}
	)
	gt.Equal("t", ValueOf(s).Index(0).Interface(), "string", 0)
	gt.Equal("s", ValueOf(a).Index(0).Interface(), "array", 0)
	gt.Equal("s", ValueOf(l).Index(0).Interface(), "s", "slice", 0)
	gt.Equal("s", ValueOf(m).Index(0).Interface(), "s", "map", 0)
	gt.Equal("s", ValueOf(d).Index(0).Interface(), "s", "struct", 0)
	gt.Equal("s", ValueOf(a1).Index(0).Interface(), "s", "array(1)", 0)
	gt.Equal("s", ValueOf(d1).Index(0).Interface(), "s", "struct(1)", 0)

	gt.Equal("e", ValueOf(&s).Elem().Index(1).Interface(), "*string", 1)
	gt.Equal("s", ValueOf(&a).Elem().Index(1).Interface(), "*array", 1)
	gt.Equal("s", ValueOf(&l).Elem().Index(1).Interface(), "*slice", 1)
	gt.Equal("s", ValueOf(&m).Elem().Index(1).Interface(), "*map", 1)
	gt.Equal("s", ValueOf(&d).Elem().Index(1).Interface(), "*struct", 1)
	gt.Equal("s", ValueOf(&a1).Elem().Index(0).Interface(), "*array(1)", 0)
	gt.Equal("s", ValueOf(&d1).Elem().Index(0).Interface(), "*struct(1)", 0)

	gt.Msg = "Testing ValueOf(%s).StructField(%d)"
	gt.Equal("s", ValueOf(&d).Elem().StructField("V2").Interface(), "*struct", "V2")
	gt.Equal("s", ValueOf(&d1).Elem().StructField("V1").Interface(), "*struct(1)", "V1")
}

func TestValueSerialize(t *testing.T) {
	gt := test.New(t, config)
	gt.Msg = "Testing ValueOf(%s).Serialize()"
	var (
		s  = "s"
		a  = [2]string{"s", "s"}
		a1 = [1]string{"s"}
		l  = []string{"s", "s"}
		m  = map[string]string{"0": "s", "1": "s"}
		d  = string_struct{"s", "s"}
		d1 = string_struct_single{"s"}

		pa = [2]*string{&s, &s}
		/* pl  = []*string{&s, &s}
		pm  = map[string]*string{"0": &s, "1": &s}
		pd  = string_ptr_struct{&s, &s}
		pa1 = [1]*string{&s}
		pd1 = string_ptr_struct_single{&s} */
	)
	gt.Equal(`"s"`, ValueOf(s).Serialize(), "string")
	gt.Equal(`["s","s"]`, ValueOf(a).Serialize(), "array")
	gt.Equal(`["s","s"]`, ValueOf(l).Serialize(), "slice")
	gt.Equal(`{"0":"s","1":"s"}`, ValueOf(m).Serialize(), "map")
	gt.Equal(`{"V1":"s","V2":"s"}`, ValueOf(d).Serialize(), "struct")
	gt.Equal(`["s"]`, ValueOf(a1).Serialize(), "array(1)")
	gt.Equal(`{"V1":"s"}`, ValueOf(d1).Serialize(), "struct(1)")

	gt.Equal(`"s"`, ValueOf(&s).Serialize(), "*string")
	gt.Equal(`["s","s"]`, ValueOf(&pa).Serialize(), "*array")
	/* gt.Equal(`["s","s"]`, ValueOf(&pl).Serialize(), "*slice")
	gt.Equal(`{"0":"s","1":"s"}`, ValueOf(&pm).Serialize(), "*map")
	gt.Equal(`{"V1":"s","V2":"s"}`, ValueOf(&pd).Serialize(), "*struct")
	gt.Equal(`["s"]`, ValueOf(&pa1).Serialize(), "*array(1)")
	gt.Equal(`{"V1":"s"}`, ValueOf(&pd1).Serialize(), "*struct(1)") */
}
