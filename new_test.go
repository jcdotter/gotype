package gotype

import (
	"fmt"
	"math"
	"testing"
	"unsafe"
)

func TestVnew(t *testing.T) {
	/* a := [2]any{}
	s := "test"
	p := &s
	vl := ValueOf(p)
	fmt.Println(vl.typ, vl.Interface())
	va := ValueOf(a)
	*(*any)(va.Index(0).ptr) = vl.Interface()
	*(*any)(va.Index(1).ptr) = vl.Interface()
	fmt.Println(va, va.Index(1).typ)
	os.Exit(1) */

	l := []string{
		/* "array_array", "array_map", "array_slice", "array_struct", "array_any",
		"array_ptr_bool", "array_ptr_int", "array_ptr_string",
		"array_ptr_array", "array_ptr_map", "array_ptr_slice", "array_ptr_struct",
		"array_array_single", "array_map_single", "array_slice_single", "array_struct_single", "array_any_single",
		"array_ptr_bool_single", "array_ptr_int_single", "array_ptr_string_single",
		"array_ptr_array_single", "array_ptr_map_single", "array_ptr_slice_single", "array_ptr_struct_single", */
		"[2]*[1]string",
	}
	v := getTestVars()
	for n, v := range v {
		var proc bool
		if len(l) == 0 {
			proc = true
		} else {
			for _, a := range l {
				ln := int(math.Min(float64(len(a)), float64(len(n))))
				if n[:ln] == a {
					proc = true
					break
				}
			}
		}
		if proc {
			r := ValueOf(v)
			fmt.Println()
			fmt.Println(n, r.Serialize())
			e := r.vNew()
			fmt.Println("new value:", e.Interface(), e.Serialize())
		}
	}

	/* os.Exit(1)
	i := 0
	for n, v := range getTestVars() {
		i++
		fmt.Println("num:", i)
		r := ValueOf(v)
		fmt.Println()
		fmt.Println(n, r.typ, r.Serialize())
		fmt.Println("new value:", r.rnewdeep().Serialize())
	} */
}

func (v VALUE) vNew() VALUE {
	fmt.Println("create new:", v.typ)
	n := v.rnew()
	indir := func(t *rtype) bool {
		k := t.Kind()
		return k == Array || k == Struct
	}
	set := func(elem VALUE, ptr unsafe.Pointer, typ *rtype) {
		elem.typ = typ
		nv := elem.vNew().Elem()
		kind := typ.Kind()
		if kind == Interface {
			*(*any)(ptr) = nv.Interface()
		} else if !indir(elem.typ) {
			/* if elem.typ.Kind() == Pointer && elem.typ.elem().Kind() == Array {
				fmt.Println("check")
				*(*unsafe.Pointer)(ptr) = *(*unsafe.Pointer)(nv.ptr)
				//typedmemmove(elem.typ, ptr, unsafe.Pointer(&nv.ptr))
				return
			} */
			typedmemmove(elem.typ, ptr, nv.ptr)
		} else {
			*(*unsafe.Pointer)(ptr) = nv.ptr
		}
	}
	switch v.Kind() {
	case Array:
		t := (*arrayType)(unsafe.Pointer(v.typ)).elem
		if !t.Kind().IsBasic() {
			a := (ARRAY)(n.Elem())
			(ARRAY)(v).ForEach(func(i int, k string, e VALUE) (brake bool) {
				set(e, offset(a.ptr, uintptr(i)*t.size), t)
				fmt.Println("set elem:", i)
				return
			})
		}
	case Map:
		*(*unsafe.Pointer)(n.ptr) = makemap(v.typ, (MAP)(v).Len(), nil)
		m := (MAP)(n.Elem())
		(MAP)(v).ForEach(func(i int, k string, e VALUE) (brake bool) {
			nv := e.vNew().Elem()
			m.Set(k, nv)
			return
		})
	case Pointer:
		fmt.Println("drill ptr:", v.typ.elem(), v.Elem().typ)
		*(*unsafe.Pointer)(n.ptr) = v.Elem().vNew().ptr
	case Slice:
		l := (*sliceHeader)(v.ptr).Len
		t := (*sliceType)(unsafe.Pointer(v.typ)).elem
		*(*unsafe.Pointer)(&n.ptr) = unsafe.Pointer(&sliceHeader{unsafe_NewArray(t, l), l, l})
		if !t.Kind().IsBasic() {
			s := (SLICE)(n.Elem())
			(SLICE)(v).ForEach(func(i int, k string, e VALUE) (brake bool) {
				set(e, offset((*sliceHeader)(s.ptr).Data, uintptr(i)*t.size), t)
				return
			})
		}
	case Struct:
		s := (STRUCT)(n.Elem())
		f := (*structType)(unsafe.Pointer(v.typ)).fields
		(STRUCT)(v).ForEach(func(i int, k string, e VALUE) (brake bool) {
			set(e, offset(s.ptr, f[i].offset), f[i].typ)
			return
		})
	}
	fmt.Println("creating from:", v.typ, v.Interface(), v.Serialize())
	fmt.Println("returning:", n.typ, n.Interface(), n.Serialize())
	return n
}
