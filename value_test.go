package gotype

import (
	"fmt"
	"testing"

	tst "github.com/jcdotter/gotester"
)

var config = &tst.Config{
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

func TestValueNew(t *testing.T) {
	gt := tst.New(t, config)
	gt.Msg = "Testing VALUE.New(%s)"
	b, i, s, a, l, m, d := true, 1, "test", [2]string{"s", "s"}, []string{"s", "s"}, map[string]string{"0": "s", "1": "s"}, string_struct{"s", "s"}
	tests := []tst.Assertion{
		{"bool", false, ValueOf(b).New().Interface(), nil},
		{"int", 0, ValueOf(i).New().Interface(), nil},
		{"string", "", ValueOf(s).New().Interface(), nil},
		{"array", [2]string{}, ValueOf(a).New().Interface(), nil},
		{"slice", []string(nil), ValueOf(l).New(true).Interface(), nil},
		{"map", map[string]string(nil), ValueOf(m).New(true).Interface(), nil},
		{"struct", string_struct{}, ValueOf(d).New().Interface(), nil},

		{"ptr_bool", true, ValueOf(&b).New().Set(&b).Elem().Interface(), nil},
		/* {"ptr_int", 0, ValueOf(&i).New().Interface(), nil},
		{"ptr_string", "", ValueOf(&s).New().Interface(), nil}, */
		/* {"ptr_array", [2]string{}, ValueOf(&a).New().Interface(), nil},
		{"ptr_slice", []string(nil), ValueOf(&l).New(true).Interface(), nil},
		{"ptr_map", map[string]string(nil), ValueOf(&m).New(true).Interface(), nil},
		{"ptr_struct", string_struct{}, ValueOf(&d).New().Interface(), nil}, */
	}
	gt.EqualMany(tests...)
}

func TestValue(t *testing.T) {
	//gt := tst.New(t, config)

	for n, v := range createTestVars() {
		fmt.Println(n, ":", ValueOf(v).New())
	}

}

func vget(v VALUE, i int) VALUE {
	switch v.Kind() {
	case Bool, Int, String:
		return v
	case Pointer:
		return vget(v.Elem(), i)
	case Interface:
		return vget(v.SetType(), i)
	default:
		return vget(v.Index(i), i)
	}
}
