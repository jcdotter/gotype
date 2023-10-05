// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	"fmt"
	"os"
	"testing"
	"unsafe"
)

func TestAddrs(t *testing.T) {

	type bool_struct struct {
		V1 bool
		V2 bool
	}

	type int_struct struct {
		V1 int
		V2 int
	}

	type array_struct struct {
		V1 [2]string
		V2 [2]string
	}

	type slice_struct struct {
		V1 []string
		V2 []string
	}

	type map_struct struct {
		V1 map[string]string
		V2 map[string]string
	}

	type string_struct struct {
		V1 string
		V2 string
	}

	type struct_struct struct {
		V1 string_struct
		V2 string_struct
	}

	type bool_ptr_struct struct {
		V1 *bool
		V2 *bool
	}

	type int_ptr_struct struct {
		V1 *int
		V2 *int
	}

	type string_ptr_struct struct {
		V1 *string
		V2 *string
	}

	type array_ptr_struct struct {
		V1 *[2]string
		V2 *[2]string
	}

	type slice_ptr_struct struct {
		V1 *[]string
		V2 *[]string
	}

	type map_ptr_struct struct {
		V1 *map[string]string
		V2 *map[string]string
	}

	type struct_ptr_struct struct {
		V1 *string_struct
		V2 *string_struct
	}

	type any_struct struct {
		V1 any
		V2 any
	}

	var (
		b = false
		i = 101
		s = ""
	)
	*(*int)(ValueOf(&i).Elem().ptr) = 0
	var (
		a = [2]string{s, s}
		l = []string{s, s}
		m = map[string]string{"0": s, "1": s}
		d = string_struct{s, s}

		ab  = [2]bool{b, b}
		ai  = [2]int{i, i}
		as  = a
		aa  = [2][2]string{a, a}
		al  = [2][]string{l, l}
		am  = [2]map[string]string{m, m}
		ad  = [2]string_struct{d, d}
		apb = [2]*bool{&b, &b}
		api = [2]*int{&i, &i}
		aps = [2]*string{&s, &s}
		apa = [2]*[2]string{&a, &a}
		apl = [2]*[]string{&l, &l}
		apm = [2]*map[string]string{&m, &m}
		apd = [2]*string_struct{&d, &d}

		sb  = []bool{b, b}
		si  = []int{i, i}
		ss  = l
		sa  = [][2]string{a, a}
		sl  = [][]string{l, l}
		sm  = []map[string]string{m, m}
		sd  = []string_struct{d, d}
		spb = []*bool{&b, &b}
		spi = []*int{&i, &i}
		sps = []*string{&s, &s}
		spa = []*[2]string{&a, &a}
		spl = []*[]string{&l, &l}
		spm = []*map[string]string{&m, &m}
		spd = []*string_struct{&d, &d}

		mb  = map[string]bool{"0": b, "1": b}
		mi  = map[string]int{"0": i, "1": i}
		ms  = m
		ma  = map[string][2]string{"0": a, "1": a}
		ml  = map[string][]string{"0": l, "1": l}
		mm  = map[string]map[string]string{"0": m, "1": m}
		md  = map[string]string_struct{"0": d, "1": d}
		mpb = map[string]*bool{"0": &b, "1": &b}
		mpi = map[string]*int{"0": &i, "1": &i}
		mps = map[string]*string{"0": &s, "1": &s}
		mpa = map[string]*[2]string{"0": &a, "1": &a}
		mpl = map[string]*[]string{"0": &l, "1": &l}
		mpm = map[string]*map[string]string{"0": &m, "1": &m}
		mpd = map[string]*string_struct{"0": &d, "1": &d}

		db  = bool_struct{b, b}
		di  = int_struct{i, i}
		ds  = d
		da  = array_struct{a, a}
		dl  = slice_struct{l, l}
		dm  = map_struct{m, m}
		dd  = struct_struct{d, d}
		dpb = bool_ptr_struct{&b, &b}
		dpi = int_ptr_struct{&i, &i}
		dps = string_ptr_struct{&s, &s}
		dpa = array_ptr_struct{&a, &a}
		dpl = slice_ptr_struct{&l, &l}
		dpm = map_ptr_struct{&m, &m}
		dpd = struct_ptr_struct{&d, &d}

		pb = &b
		pi = &i
		ps = &s
		pa = &a
		pl = &l
		pm = &m
		pd = &d

		aib  = [2]any{b, b}
		aii  = [2]any{i, i}
		ais  = [2]any{s, s}
		aia  = [2]any{a, a}
		ail  = [2]any{l, l}
		aim  = [2]any{m, m}
		aid  = [2]any{d, d}
		aipb = [2]any{&b, &b}
		aipi = [2]any{&i, &i}
		aips = [2]any{&s, &s}
		aipa = [2]any{&a, &a}
		aipl = [2]any{&l, &l}
		aipm = [2]any{&m, &m}
		aipd = [2]any{&d, &d}

		sib  = []any{b, b}
		sii  = []any{i, i}
		sis  = []any{s, s}
		sia  = []any{a, a}
		sil  = []any{l, l}
		sim  = []any{m, m}
		sid  = []any{d, d}
		sipb = []any{&b, &b}
		sipi = []any{&i, &i}
		sips = []any{&s, &s}
		sipa = []any{&a, &a}
		sipl = []any{&l, &l}
		sipm = []any{&m, &m}
		sipd = []any{&d, &d}

		mib  = map[string]any{"0": b, "1": b}
		mii  = map[string]any{"0": i, "1": i}
		mis  = map[string]any{"0": s, "1": s}
		mia  = map[string]any{"0": a, "1": a}
		mil  = map[string]any{"0": l, "1": l}
		mim  = map[string]any{"0": m, "1": m}
		mid  = map[string]any{"0": d, "1": d}
		mipb = map[string]any{"0": &b, "1": &b}
		mipi = map[string]any{"0": &i, "1": &i}
		mips = map[string]any{"0": &s, "1": &s}
		mipa = map[string]any{"0": &a, "1": &a}
		mipl = map[string]any{"0": &l, "1": &l}
		mipm = map[string]any{"0": &m, "1": &m}
		mipd = map[string]any{"0": &d, "1": &d}

		dib  = any_struct{b, b}
		dii  = any_struct{i, i}
		dis  = any_struct{s, s}
		dia  = any_struct{a, a}
		dil  = any_struct{l, l}
		dim  = any_struct{m, m}
		did  = any_struct{d, d}
		dipb = any_struct{&b, &b}
		dipi = any_struct{&i, &i}
		dips = any_struct{&s, &s}
		dipa = any_struct{&a, &a}
		dipl = any_struct{&l, &l}
		dipm = any_struct{&m, &m}
		dipd = any_struct{&d, &d}
	)

	vars := []any{
		b, i, s, // 3
		ab, ai, as, aa, al, am, ad, // 10
		sb, si, ss, sa, sl, sm, sd, // 17
		mb, mi, ms, ma, ml, mm, md, // 24
		db, di, ds, da, dl, dm, dd, // 31

		aib, aii, ais, aia, ail, aim, aid, // 38
		sib, sii, sis, sia, sil, sim, sid, // 45
		mib, mii, mis, mia, mil, mim, mid, // 52
		dib, dii, dis, dia, dil, dim, did, // 59

		apb, api, aps, apa, apl, apm, apd, // 66
		spb, spi, sps, spa, spl, spm, spd, // 73
		mpb, mpi, mps, mpa, mpl, mpm, mpd, // 80
		dpb, dpi, dps, dpa, dpl, dpm, dpd, // 87

		aipb, aipi, aips, aipa, aipl, aipm, aipd, // 94
		sipb, sipi, sips, sipa, sipl, sipm, sipd, // 101
		mipb, mipi, mips, mipa, mipl, mipm, mipd, // 108
		dipb, dipi, dips, dipa, dipl, dipm, dipd, // 115

		&b, &i, &s, // 118
		&ab, &ai, &as, &aa, &al, &am, &ad, // 125
		&sb, &si, &ss, &sa, &sl, &sm, &sd, // 132
		&mb, &mi, &ms, &ma, &ml, &mm, &md, // 139
		&db, &di, &ds, &da, &dl, &dm, &dd, // 146

		&aib, &aii, &ais, &aia, &ail, &aim, &aid, // 153
		&sib, &sii, &sis, &sia, &sil, &sim, &sid, // 160
		&mib, &mii, &mis, &mia, &mil, &mim, &mid, // 167
		&dib, &dii, &dis, &dia, &dil, &dim, &did, // 174

		&apb, &api, &aps, &apa, &apl, &apm, &apd, // 181
		&spb, &spi, &sps, &spa, &spl, &spm, &spd, // 188
		&mpb, &mpi, &mps, &mpa, &mpl, &mpm, &mpd, // 195
		&dpb, &dpi, &dps, &dpa, &dpl, &dpm, &dpd, // 202

		&aipb, &aipi, &aips, &aipa, &aipl, &aipm, &aipd, // 209
		&sipb, &sipi, &sips, &sipa, &sipl, &sipm, &sipd, // 216
		&mipb, &mipi, &mips, &mipa, &mipl, &mipm, &mipd, // 223
		&dipb, &dipi, &dips, &dipa, &dipl, &dipm, &dipd, // 230

		&pb, &pi, &ps, &pa, &pl, &pm, &pd, // 237
	}

	/* v := ValueOf(vars[62])
	fmt.Println("TYPE:", v.typ, " PTR:", v.ptr, "FCT:", pfactor(v.typ), //" VAL:", ***(***[]string)(v.ptr),
		" IFACE:", fmt.Sprintf("%#v", v.Interface()), " SRL:", v.Serialize())

	n := v.NewDeep()
	fmt.Println("TYPE:", n.typ, " PTR:", n.ptr, "FCT:", pfactor(n.typ), //" VAL:", ***(***[]string)(n.ptr),
		" IFACE:", fmt.Sprintf("%#v", n.Interface()), " SRL:", n.Serialize())

	os.Exit(1) */

	for j, v := range vars {
		//fmt.Println("pfactor(" + INT(pfactor(ValueOf(v).typ)).String() + ")")
		x := ValueOf(v).NewDeep().Interface()
		val := ValueOf(x)
		vi := vget(val, 1)
		num := " (" + ValueOf(j).String() + "): "
		var o, n any
		switch vi.Kind() {
		case Bool:
			o, n = false, true
		case Int:
			o, n = 0, 1
		case String:
			o, n = "", "true"
		}
		tl := "PointerGet" + num + " Output: " + fmt.Sprintf(" Type: %T  Value: %s", val.Interface(), val)
		if err := test_validations([]validation_test{{tl, vi.Interface(), o}}); err != nil {
			t.Fatalf(fmt.Sprint(err))
		}
		vi = vset(val, 1, n)
		tl = "PointerSet" + num + " Output: " + fmt.Sprintf(" Type: %T  Value: %s", val.Interface(), val)
		if err := test_validations([]validation_test{{tl, vi.Interface(), n}}); err != nil {
			t.Fatalf(fmt.Sprint(err))
		}
		x = ValueOf(v).NewDeep().Interface()
		val = ValueOf(x)
		vreset(val, 1, true)
		vi = vget(val, 1)
		tl = "PointerSetParent" + num + " Output: " + fmt.Sprintf(" Type: %T  Value: %s", val.Interface(), val)
		if err := test_validations([]validation_test{{tl, vi.Interface(), n}}); err != nil {
			t.Fatalf(fmt.Sprint(err))
		}
		vi = vget(ValueOf(x), 1)
		tl = "PointerSetOrigin" + num + " Output: " + fmt.Sprintf(" Type: %T  Value: %s", val.Interface(), val)
		if err := test_validations([]validation_test{{tl, vi.Interface(), n}}); err != nil {
			t.Fatalf(fmt.Sprint(err))
		}
	}
}

func TestPtrElem(t *testing.T) {
	type strukt struct {
		V1 string
		V2 string
	}
	var (
		b   = false
		i   = 12345
		s   = "false"
		a   = [2]string{"false", "false"}
		m   = map[string]string{"0": "false", "1": "false"}
		l   = []string{"false", "false"}
		d   = strukt{"false", "false"}
		p   = "false"
		pp  = &p
		pb  = false
		ppb = &pb
		ps  = "false"
		pps = &ps
		pm  = map[string]string{"0": "false", "1": "false"}
		ppm = &pm
		a1  = [1]string{"false"}
		as  = "false"
		at  = "true"
		ap  = [1]*string{&as}
		/* mm  = &map[string]string{"0": "false", "1": "true"}
		ll  = &l
		lll = &ll
		sss = &p
		a1  = [1]map[string]string{{"1": "true"}} */
	)

	v := ValueOf(&ppm)
	fmt.Println(v.Elem().Elem().MapIndex("1"))
	fmt.Println(*(*string)(v.Elem().Elem().MapIndex("1").ptr))
	v.Elem().Elem().SetIndex("1", "true")
	fmt.Println(ValueOf(pm))

	os.Exit(1)

	ValueOf(&b).Elem().Set(true)
	ValueOf(&i).Elem().Set(54321)
	ValueOf(&s).Elem().Set("true")
	ValueOf(&a).Elem().Index(1).Set("true")
	ValueOf(&m).Elem().SetIndex("1", "true")
	ValueOf(&l).Elem().Index(1).Set("true")
	ValueOf(&d).Elem().Index(1).Set("true")
	ValueOf(&pp).Elem().Set("true")
	ValueOf(&ppb).Elem().Elem().Set(true)
	ValueOf(&pps).Elem().Elem().Set("true")
	ValueOf(&ppm).Elem().Elem().SetIndex("1", "true")
	ValueOf(&a1).Elem().Index(0).Set("true")
	ValueOf(&ap).Elem().Index(0).Set(&at)

	validations := []validation_test{
		{"Set Ptr Elem Bool", b, true},
		{"Set Ptr Elem Int", i, 54321},
		{"Set Ptr Elem String", s, "true"},
		{"Set Ptr Elem Array", a, [2]string{"false", "true"}},
		{"Set Ptr Elem Map", m, map[string]string{"0": "false", "1": "true"}},
		{"Set Ptr Elem Slice", l, []string{"false", "true"}},
		{"Set Ptr Elem Struct", d, strukt{"false", "true"}},
		{"Set Ptr Elem Ptr String", p, "true"},
		{"Set Ptr Elem Ptr Ptr Bool", pb, true},
		{"Set Ptr Elem Ptr Ptr String", ps, "true"},
		{"Set Ptr Elem Ptr Ptr Map", pm, map[string]string{"0": "false", "1": "true"}},
		{"Set Ptr Elem [1]Array", ap, [1]*string{&at}},
	}
	if err := test_validations(validations); err != nil {
		t.Fatalf(fmt.Sprint(err))
	}

}

func TestEmptyAddrs(t *testing.T) {

	type bool_struct struct {
		V1 bool
		V2 bool
	}

	type int_struct struct {
		V1 int
		V2 int
	}

	type array_struct struct {
		V1 [2]string
		V2 [2]string
	}

	type slice_struct struct {
		V1 []string
		V2 []string
	}

	type map_struct struct {
		V1 map[string]string
		V2 map[string]string
	}

	type struct_struct struct {
		V1 string_struct
		V2 string_struct
	}

	type bool_ptr_struct struct {
		V1 *bool
		V2 *bool
	}

	type int_ptr_struct struct {
		V1 *int
		V2 *int
	}

	type string_ptr_struct struct {
		V1 *string
		V2 *string
	}

	type array_ptr_struct struct {
		V1 *[2]string
		V2 *[2]string
	}

	type slice_ptr_struct struct {
		V1 *[]string
		V2 *[]string
	}

	type map_ptr_struct struct {
		V1 *map[string]string
		V2 *map[string]string
	}

	type struct_ptr_struct struct {
		V1 *string_struct
		V2 *string_struct
	}

	var (
		ab  [2]bool
		ai  [2]int
		as  [2]string
		aa  [2][2]string
		al  [2][]string
		am  [2]map[string]string
		ad  [2]string_struct
		apb [2]*bool
		api [2]*int
		aps [2]*string
		apa [2]*[2]string
		apl [2]*[]string
		apm [2]*map[string]string
		apd [2]*string_struct

		sb  []bool
		si  []int
		ss  []string
		sa  [][2]string
		sl  [][]string
		sm  []map[string]string
		sd  []string_struct
		spb []*bool
		spi []*int
		sps []*string
		spa []*[2]string
		spl []*[]string
		spm []*map[string]string
		spd []*string_struct

		mb  map[string]bool
		mi  map[string]int
		ms  map[string]string
		ma  map[string][2]string
		ml  map[string][]string
		mm  map[string]map[string]string
		md  map[string]string_struct
		mpb map[string]*bool
		mpi map[string]*int
		mps map[string]*string
		mpa map[string]*[2]string
		mpl map[string]*[]string
		mpm map[string]*map[string]string
		mpd map[string]*string_struct

		db  bool_struct
		di  int_struct
		ds  string_struct
		da  array_struct
		dl  slice_struct
		dm  map_struct
		dd  struct_struct
		dpb bool_ptr_struct
		dpi int_ptr_struct
		dps string_ptr_struct
		dpa array_ptr_struct
		dpl slice_ptr_struct
		dpm map_ptr_struct
		dpd struct_ptr_struct
	)

	vars := []any{
		ab, ai, as, aa, al, am, ad, apb, api, aps, apa, apl, apm, apd, // 14
		sb, si, ss, sa, sl, sm, sd, spb, spi, sps, spa, spl, spm, spd, // 28
		mb, mi, ms, ma, ml, mm, md, mpb, mpi, mps, mpa, mpl, mpm, mpd, // 42
		db, di, ds, da, dl, dm, dd, dpb, dpi, dps, dpa, dpl, dpm, dpd, // 56
		&ab, &ai, &as, &aa, &al, &am, &ad, &apb, &api, &aps, &apa, &apl, &apm, &apd, // 70
		&sb, &si, &ss, &sa, &sl, &sm, &sd, &spb, &spi, &sps, &spa, &spl, &spm, &spd, // 84
		&mb, &mi, &ms, &ma, &ml, &mm, &md, &mpb, &mpi, &mps, &mpa, &mpl, &mpm, &mpd, // 98
		&db, &di, &ds, &da, &dl, &dm, &dd, &dpb, &dpi, &dps, &dpa, &dpl, &dpm, &dpd, // 112
	}

	for j, v := range vars {
		x := ValueOf(v).New().Interface()
		val := ValueOf(x)
		pset(val, 1)
		vi := vget(val, 1)
		num := " (" + ValueOf(j).String() + "): "
		var n any
		switch vi.Kind() {
		case Bool:
			n = true
		case Int:
			n = 1
		case String:
			n = "true"
		}
		tl := "PointerSetEmpty" + num + " Output: " + fmt.Sprintf(" Type: %T  Value: %s", val.Interface(), val)
		if err := test_validations([]validation_test{{tl, vi.Interface(), n}}); err != nil {
			t.Fatalf(fmt.Sprint(err))
		}
		vi = vget(ValueOf(x), 1)
		tl = "PointerSetEmptyOrigin" + num + " Output: " + fmt.Sprintf(" Type: %T  Value: %s", val.Interface(), val)
		if err := test_validations([]validation_test{{tl, vi.Interface(), n}}); err != nil {
			t.Fatalf(fmt.Sprint(err))
		}
	}
}

func TestEmptySingleAddrs(t *testing.T) {

	type bool_struct struct {
		V1 bool
	}

	type int_struct struct {
		V1 int
	}

	type array_struct struct {
		V1 [1]string
	}

	type slice_struct struct {
		V1 []string
	}

	type map_struct struct {
		V1 map[string]string
	}

	type string_single_struct struct {
		V1 string
	}

	type struct_struct struct {
		V1 string_single_struct
	}

	type bool_ptr_struct struct {
		V1 *bool
	}

	type int_ptr_struct struct {
		V1 *int
	}

	type string_ptr_struct struct {
		V1 *string
	}

	type array_ptr_struct struct {
		V1 *[1]string
	}

	type slice_ptr_struct struct {
		V1 *[]string
	}

	type map_ptr_struct struct {
		V1 *map[string]string
	}

	type struct_ptr_struct struct {
		V1 *string_single_struct
	}

	var (
		ab  [1]bool
		ai  [1]int
		as  [1]string
		aa  [1][1]string
		al  [1][]string
		am  [1]map[string]string
		ad  [1]string_single_struct
		apb [1]*bool
		api [1]*int
		aps [1]*string
		apa [1]*[1]string
		apl [1]*[]string
		apm [1]*map[string]string
		apd [1]*string_single_struct

		sb  []bool
		si  []int
		ss  []string
		sa  [][1]string
		sl  [][]string
		sm  []map[string]string
		sd  []string_single_struct
		spb []*bool
		spi []*int
		sps []*string
		spa []*[1]string
		spl []*[]string
		spm []*map[string]string
		spd []*string_single_struct

		mb  map[string]bool
		mi  map[string]int
		ms  map[string]string
		ma  map[string][1]string
		ml  map[string][]string
		mm  map[string]map[string]string
		md  map[string]string_single_struct
		mpb map[string]*bool
		mpi map[string]*int
		mps map[string]*string
		mpa map[string]*[1]string
		mpl map[string]*[]string
		mpm map[string]*map[string]string
		mpd map[string]*string_single_struct

		db  bool_struct
		di  int_struct
		ds  string_single_struct
		da  array_struct
		dl  slice_struct
		dm  map_struct
		dd  struct_struct
		dpb bool_ptr_struct
		dpi int_ptr_struct
		dps string_ptr_struct
		dpa array_ptr_struct
		dpl slice_ptr_struct
		dpm map_ptr_struct
		dpd struct_ptr_struct
	)

	vars := []any{
		ab, ai, as, aa, al, am, ad, apb, api, aps, apa, apl, apm, apd, // 14
		sb, si, ss, sa, sl, sm, sd, spb, spi, sps, spa, spl, spm, spd, // 28
		mb, mi, ms, ma, ml, mm, md, mpb, mpi, mps, mpa, mpl, mpm, mpd, // 42
		db, di, ds, da, dl, dm, dd, dpb, dpi, dps, dpa, dpl, dpm, dpd, // 56
		&ab, &ai, &as, &aa, &al, &am, &ad, &apb, &api, &aps, &apa, &apl, &apm, &apd, // 70
		&sb, &si, &ss, &sa, &sl, &sm, &sd, &spb, &spi, &sps, &spa, &spl, &spm, &spd, // 84
		&mb, &mi, &ms, &ma, &ml, &mm, &md, &mpb, &mpi, &mps, &mpa, &mpl, &mpm, &mpd, // 98
		&db, &di, &ds, &da, &dl, &dm, &dd, &dpb, &dpi, &dps, &dpa, &dpl, &dpm, &dpd, // 112
	}

	/* v := ValueOf(vars[47]).New()
	fmt.Println("TYPE:", v.typ, " PTR:", v.ptr, "FCT:", pfactor(v.typ), " VAL:", *(*map_struct)(v.ptr),
		" IFACE:", fmt.Sprintf("%#v", v.Interface()), " SRL:", v.Serialize()) */
	//psingleset(v, 0)
	/* i := v.Index(0) //.Set(map[string]string{"0": "true"})
	i.Init()
	m := ValueOf(map[string]string{"0": "true"})
	//**(**[48]byte)(i.ptr) = **(**[48]byte)(m.ptr)
	fmt.Println(**(**[48]byte)(i.ptr))
	fmt.Println(**(**[48]byte)(m.ptr))
	*/

	/* a := (any)(dm)
	v := *(*VALUE)(unsafe.Pointer(&a))
	fmt.Println(v.typ.size) */
	/* p := mallocgc(v.typ.size, v.typ, true)
	//fmt.Println(v.ptr)
	fmt.Println(**(**[48]byte)(p)) */

	/* fmt.Println("TYPE:", v.typ, " PTR:", v.ptr, "FCT:", pfactor(v.typ), " VAL:", *(*map_struct)(v.ptr),
	" IFACE:", fmt.Sprintf("%#v", v.Interface()), " SRL:", v.Serialize()) */
	//fmt.Println(vget(v, 0))
	//os.Exit(1)

	for j, v := range vars[9:10] {
		x := ValueOf(v).New()
		fmt.Println("TYPE:", x.typ, " PTR:", x.ptr, " VAL:", *(*[1]*string)(x.ptr),
			" IFACE:", fmt.Sprintf("%#v", x.Interface()), " SRL:", x.Serialize())
		//s := "true"
		//x = x.Index(0).Set(&s)
		psingleset(x, 0)
		fmt.Println("TYPE:", x.typ, " PTR:", x.ptr, " VAL:", *(*[1]*string)(x.ptr),
			" IFACE:", fmt.Sprintf("%#v", x.Interface()), " SRL:", x.Serialize())
		fmt.Println(v.([1]*string))
		os.Exit(1)
		val := ValueOf(v).New()
		psingleset(val, 0)
		vi := vget(val, 0)
		num := " (" + ValueOf(j).String() + "): "
		var n any
		switch vi.Kind() {
		case Bool:
			n = true
		case Int:
			n = 1
		case String:
			n = "true"
		}
		tl := "PointerSetEmptySingle" + num + " Output: " + fmt.Sprintf(" Type: %T  Value: %s", val.Interface(), val)
		if err := test_validations([]validation_test{{tl, vi.Interface(), n}}); err != nil {
			t.Fatalf(fmt.Sprint(err))
		}
	}

	/* for i, v := range vars {
		val := ValueOf(v).New()
		fmt.Print(
			ValueOf(i).STRING().Width(5),
			" PointersSetEmpty: ",
			" TYPE: ", STRING(val.typ.String()).Width(32),
			" PTR:  ", STRING(fmt.Sprint(val.ptr)).Width(16),
		)
		vi := psingleset(&val, 0)
		fmt.Print(
			" TYPE: ", STRING(val.typ.String()).Width(24),
			" PTR:  ", STRING(fmt.Sprint(val.ptr)).Width(16),
			" IFC:  ", STRING(fmt.Sprintf("%#v", v)).Width(48),
			" VAL:  ", STRING(val.Serialize()).Width(40),
			" OUT:  ", vi,
			"\n",
		)
	} */
}

func TestEmptyPtrs(t *testing.T) {

	type bool_struct struct {
		V1 bool
		V2 bool
	}

	type int_struct struct {
		V1 int
		V2 int
	}

	type array_struct struct {
		V1 [2]string
		V2 [2]string
	}

	type slice_struct struct {
		V1 []string
		V2 []string
	}

	type map_struct struct {
		V1 map[string]string
		V2 map[string]string
	}

	type struct_struct struct {
		V1 string_struct
		V2 string_struct
	}

	type bool_ptr_struct struct {
		V1 *bool
		V2 *bool
	}

	type int_ptr_struct struct {
		V1 *int
		V2 *int
	}

	type string_ptr_struct struct {
		V1 *string
		V2 *string
	}

	type array_ptr_struct struct {
		V1 *[2]string
		V2 *[2]string
	}

	type slice_ptr_struct struct {
		V1 *[]string
		V2 *[]string
	}

	type map_ptr_struct struct {
		V1 *map[string]string
		V2 *map[string]string
	}

	type struct_ptr_struct struct {
		V1 *string_struct
		V2 *string_struct
	}

	var (
		b bool
		i int
		s string

		ab  [2]bool
		ai  [2]int
		as  [2]string
		aa  [2][2]string
		al  [2][]string
		am  [2]map[string]string
		ad  [2]string_struct
		apb [2]*bool
		api [2]*int
		aps [2]*string
		apa [2]*[2]string
		apl [2]*[]string
		apm [2]*map[string]string
		apd [2]*string_struct

		sb  []bool
		si  []int
		ss  []string
		sa  [][2]string
		sl  [][]string
		sm  []map[string]string
		sd  []string_struct
		spb []*bool
		spi []*int
		sps []*string
		spa []*[2]string
		spl []*[]string
		spm []*map[string]string
		spd []*string_struct

		mb  map[string]bool
		mi  map[string]int
		ms  map[string]string
		ma  map[string][2]string
		ml  map[string][]string
		mm  map[string]map[string]string
		md  map[string]string_struct
		mpb map[string]*bool
		mpi map[string]*int
		mps map[string]*string
		mpa map[string]*[2]string
		mpl map[string]*[]string
		mpm map[string]*map[string]string
		mpd map[string]*string_struct

		db  bool_struct
		di  int_struct
		ds  string_struct
		da  array_struct
		dl  slice_struct
		dm  map_struct
		dd  struct_struct
		dpb bool_ptr_struct
		dpi int_ptr_struct
		dps string_ptr_struct
		dpa array_ptr_struct
		dpl slice_ptr_struct
		dpm map_ptr_struct
		dpd struct_ptr_struct
	)

	vars := []any{
		&b, &i, &s,
		&ab, &ai, &as, &aa, &al, &am, &ad, &apb, &api, &aps, &apa, &apl, &apm, &apd, // 14
		&sb, &si, &ss, &sa, &sl, &sm, &sd, &spb, &spi, &sps, &spa, &spl, &spm, &spd, // 28
		&mb, &mi, &ms, &ma, &ml, &mm, &md, &mpb, &mpi, &mps, &mpa, &mpl, &mpm, &mpd, // 42
		&db, &di, &ds, &da, &dl, &dm, &dd, &dpb, &dpi, &dps, &dpa, &dpl, &dpm, &dpd, // 56
	}

	/* v := ValueOf(vars[15])
	fmt.Println("TYPE:", v.typ, " PTR:", v.ptr, "FCT:", pfactor(v.typ), " VAL:", **(**[]int)(v.ptr),
		" IFACE:", fmt.Sprintf("%#v", v.Interface()), " SRL:", v.Serialize()) */
	/* e := v.Elem()
	e.Init() */
	//v.Elem().SetIndex(1, 1)
	/* pset(v, 1)
	v = vget(v, 1)
	fmt.Println("TYPE:", v.typ, " PTR:", v.ptr, "FCT:", pfactor(v.typ),
		" IFACE:", fmt.Sprintf("%#v", v.Interface()), " SRL:", v.Serialize())
	v = ValueOf(vars[15])
	fmt.Println("ORIGIN:  TYPE:", v.typ, " PTR:", v.ptr, "FCT:", pfactor(v.typ),
		" IFACE:", fmt.Sprintf("%#v", v.Interface()), " SRL:", v.Serialize())
	fmt.Println(vars[15])
	*/
	//ValueOf(vars[2]).Elem().Set("1")

	v := ValueOf(vars[2])
	v.Set("true")
	//v := pset(ValueOf(vars[2]), 1)
	fmt.Println(*vars[2].(*string))
	fmt.Println(v)
	os.Exit(1)

	for j, v := range vars {
		val := ValueOf(v)
		pset(val, 1)
		vi := vget(val, 1)
		num := " (" + ValueOf(j).String() + "): "
		var n any
		switch vi.Kind() {
		case Bool:
			n = true
		case Int:
			n = 1
		case String:
			n = "true"
		}
		tl := "PointerSetEmpty" + num + " Output: " + fmt.Sprintf(" Type: %T  Value: %s", val.Interface(), val)
		if err := test_validations([]validation_test{{tl, vi.Interface(), n}}); err != nil {
			t.Fatalf(fmt.Sprint(err))
		}
		vi = vget(ValueOf(v), 1)
		tl = "PointerSetEmptyOrigin" + num + " Output: " + fmt.Sprintf(" Type: %T  Value: %s", val.Interface(), val)
		if err := test_validations([]validation_test{{tl, vi.Interface(), n}}); err != nil {
			t.Fatalf(fmt.Sprint(err))
		}
	}
}

/* func vget(v VALUE, i int) VALUE {
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
} */

func vset(v VALUE, i int, a any) VALUE {
	switch v.Kind() {
	case Bool:
		return v.Set(a)
	case Int:
		return v.Set(a)
	case String:
		return v.Set(a)
	case Pointer:
		return vset(v.Elem(), i, a)
	case Interface:
		return vset(v.SetType(), i, a)
	default:
		return vset(v.Index(i), i, a)
		/* idx := v.Index(i)
		idx.Init()
		return vset(idx, i, a) */
	}
}

func vreset(v VALUE, i int, l bool) VALUE {
	if l {
		switch v.Kind() {
		case Pointer:
			return vreset(v.Elem(), i, true)
		case Bool:
			return v.Set(true)
		case Int:
			v.Set(1)
		case String:
			v.Set("true")
		case Interface:
			return vreset(v.SetType(), i, l)
		default:
			return vreset(v.Index(i), i, false)
		}
	} else {
		k := v.Kind()
		if k == Pointer {
			k = (*ptrType)(unsafe.Pointer(v.typ)).elem.Kind()
		}
		switch k {
		case Bool:
			return v.Set(true)
		case Int:
			v.Set(1)
		case String:
			v.Set("true")
		case Interface:
			return vreset(*(*VALUE)(v.ptr), i, l)
		case Array:
			return v.Set([2]string{"true", "true"})
		case Slice:
			return v.Set([]string{"true", "true"})
		case Map:
			return v.Set(map[string]string{"0": "true", "1": "true"})
		case Struct:
			return v.Set(string_struct{"true", "true"})
		}
	}
	return v
}

func pset(v VALUE, i int) VALUE {
	switch v.Elem().Kind() {
	case Bool:
		return v.Set(true)
	case Int:
		return v.Set(1)
	case String:
		return v.Set("true")
	default:
		switch v.Elem().ElemKind() {
		case Bool:
			return v.SetIndex(i, true)
		case Int:
			return v.SetIndex(i, 1)
		case Array:
			return v.SetIndex(i, [2]string{"true", "true"})
		case Interface:
			return v.SetIndex(i, [2]string{"true", "true"})
		case Map:
			return v.SetIndex(i, map[string]string{"0": "true", "1": "true"})
		case Slice:
			return v.SetIndex(i, []string{"true", "true"})
		case String:
			return v.SetIndex(i, "true")
		case Struct:
			return v.SetIndex(i, string_struct{"true", "true"})
		case Pointer:
			switch v.Elem().elemType().elem().Kind() {
			case Bool:
				x := true
				return v.SetIndex(i, &x)
			case Int:
				x := 1
				return v.SetIndex(i, &x)
			case Array:
				x := [2]string{"true", "true"}
				return v.SetIndex(i, &x)
			case Interface:
				x := [2]string{"true", "true"}
				return v.SetIndex(i, &x)
			case Map:
				x := map[string]string{"0": "true", "1": "true"}
				return v.SetIndex(i, &x)
			case Slice:
				x := []string{"true", "true"}
				return v.SetIndex(i, &x)
			case String:
				x := "true"
				return v.SetIndex(i, &x)
			case Struct:
				x := string_struct{"true", "true"}
				return v.SetIndex(i, &x)
			}
		}
	}
	return v
}

func psingleset(v VALUE, i int) VALUE {
	switch v.Elem().ElemKind() {
	case Bool:
		return v.SetIndex(i, true)
	case Int:
		return v.SetIndex(i, 1)
	case Array:
		return v.SetIndex(i, [1]string{"true"})
	case Interface:
		return v.SetIndex(i, [1]string{"true"})
	case Map:
		return v.SetIndex(i, map[string]string{"0": "true"})
	case Slice:
		return v.SetIndex(i, []string{"true"})
	case String:
		return v.SetIndex(i, "true")
	case Struct:
		return v.SetIndex(i, string_struct_single{"true"})
	case Pointer:
		switch v.Elem().elemType().elem().Kind() {
		case Bool:
			x := true
			return v.SetIndex(i, &x)
		case Int:
			x := 1
			return v.SetIndex(i, &x)
		case Array:
			x := [1]string{"true"}
			return v.SetIndex(i, &x)
		case Interface:
			x := [1]string{"true"}
			return v.SetIndex(i, &x)
		case Map:
			x := map[string]string{"0": "true"}
			return v.SetIndex(i, &x)
		case Slice:
			x := []string{"true"}
			return v.SetIndex(i, &x)
		case String:
			x := "true"
			return v.SetIndex(i, &x)
		case Struct:
			x := string_struct_single{"true"}
			return v.SetIndex(i, &x)
		}
	}
	return v
}

func vOfArr(a any) VALUE {
	v := *(*VALUE)(unsafe.Pointer(&a))
	if (*arrayType)(unsafe.Pointer(v.typ)).elem.Kind() == Map {
		return VALUE{v.typ, unsafe.Pointer(&v.ptr), v.typ.flag()}
	}
	return v
}

func (v VALUE) iOfArr() any {
	var i any
	ef := (*VALUE)(unsafe.Pointer(&i))
	if (*arrayType)(unsafe.Pointer(v.typ)).elem.Kind() == Map {
		ef.typ, ef.ptr = v.typ, *(*unsafe.Pointer)(v.ptr)
	} else {
		ef.typ, ef.ptr = v.typ, v.ptr
	}
	return i
}

func Test1Arr(t *testing.T) {
	type strukt struct {
		V1 string
		V2 string
	}
	var (
		b   = [1]bool{false}
		i   = [1]int{12345}
		s   = [1]string{"false"}
		a   = [1][2]string{{"false"}}
		a1  = [1][1]string{{"false"}}
		m   = [1]map[string]string{{"0": "false"}}
		l   = [1][]string{{"false"}}
		d   = [1]strukt{{"false", "false"}}
		sp  = "false"
		ps  = [1]*string{&sp}
		spp = &sp
		pps = [1]**string{&spp}
		mp  = map[string]string{"0": "false"}
		pm  = [1]*map[string]string{&mp}
		mpp = &mp
		ppm = [1]**map[string]string{&mpp}
	)

	/* fmt.Println(pfactor(ValueOf(i).typ))
	os.Exit(1) */

	validations := []validation_test{
		{"[1]Array Bool", ValueOf(b).Index(0).Interface(), false},
		{"[1]Array Int", ValueOf(i).Index(0).Interface(), 12345},
		{"[1]Array String", ValueOf(s).Index(0).Interface(), "false"},
		{"[1]Array Array", ValueOf(a).Index(0).Index(0).Interface(), "false"},
		{"[1]Array [1]Array", ValueOf(a1).Index(0).Index(0).Interface(), "false"},
		{"[1]Array Map", ValueOf(m).Index(0).MapIndex("0").Interface(), "false"},
		{"[1]Array Slice", ValueOf(l).Index(0).Index(0).Interface(), "false"},
		{"[1]Array Struct", ValueOf(d).Index(0).Index(0).Interface(), "false"},
		{"[1]Array Ptr Str", ValueOf(ps).Index(0).Elem().Interface(), "false"},
		{"[1]Array Ptr Ptr Str", ValueOf(pps).Index(0).Elem().Elem().Interface(), "false"},
		{"Ptr [1]Array Ptr Ptr Str", ValueOf(&pps).Elem().Index(0).Elem().Elem().Interface(), "false"},
		{"[1]Array Ptr Map", ValueOf(pm).Index(0).Elem().MapIndex("0").Interface(), "false"},
		{"Ptr [1]Array Ptr Ptr Map", ValueOf(&ppm).Elem().Index(0).Elem().Elem().MapIndex("0").Interface(), "false"},
	}
	if err := test_validations(validations); err != nil {
		t.Fatalf(fmt.Sprint(err))
	}

	/* ValueOf(&b).Index(0).Set(true)
	ValueOf(&i).Index(0).Set(54321)
	ValueOf(&s).Index(0).Set("true")
	ValueOf(&a).Index(0).Index(0).Set("true")
	ValueOf(&a1).Index(0).Index(0).Set("true")
	ValueOf(&m).Index(0).MapIndex("0").Set("true")
	ValueOf(&l).Index(0).Index(0).Set("true")
	ValueOf(&d).Index(0).Index(0).Set("true")
	ValueOf(&pps).Index(0).Set("true")
	ValueOf(&ppm).Index(0).SetIndex("0", "true")

	validations = []validation_test{
		{"[1]Array Set Bool pfactor:" + INT(pfactor(ValueOf(b).typ)).String(), ValueOf(b).Index(0).Interface(), true},
		{"[1]Array Set Int pfactor:" + INT(pfactor(ValueOf(i).typ)).String(), ValueOf(i).Index(0).Interface(), 54321},
		{"[1]Array Set String pfactor:" + INT(pfactor(ValueOf(s).typ)).String(), ValueOf(s).Index(0).Interface(), "true"},
		{"[1]Array Array pfactor:" + INT(pfactor(ValueOf(a).typ)).String(), ValueOf(a).Index(0).Index(0).Interface(), "true"},
		{"[1]Array [1]Array pfactor:" + INT(pfactor(ValueOf(a1).typ)).String(), ValueOf(a1).Index(0).Index(0).Interface(), "true"},
		{"[1]Array Map pfactor:" + INT(pfactor(ValueOf(m).typ)).String(), ValueOf(m).Index(0).MapIndex("0").Interface(), "true"},
		{"[1]Array Slice pfactor:" + INT(pfactor(ValueOf(l).typ)).String(), ValueOf(l).Index(0).Index(0).Interface(), "true"},
		{"[1]Array Struct pfactor:" + INT(pfactor(ValueOf(d).typ)).String(), ValueOf(d).Index(0).Index(0).Interface(), "true"},
		{"Ptr [1]Array Ptr Ptr Str pfactor:" + INT(pfactor(ValueOf(pps).typ)).String(), ValueOf(pps).Index(0).Elem().Elem().Interface(), "true"},
		{"Ptr [1]Array Ptr Ptr Map pfactor:" + INT(pfactor(ValueOf(ppm).typ)).String(), ValueOf(ppm).Index(0).Elem().Elem().MapIndex("0").Interface(), "true"},
	}
	if err := test_validations(validations); err != nil {
		t.Fatalf(fmt.Sprint(err))
	} */

}

func Test_Ptr(t *testing.T) {
	type strukt struct {
		Item string
	}
	type struqt struct {
		Item int
	}
	var (
		a [2]string
		m map[string]string
		p *string
		s []float64
		d strukt
		x strukt

		aa = (any)(a)
		am = (any)(m)
		ap = (any)(p)
		as = (any)(s)
		ad = (any)(d)
		ax = (any)(x)
	)

	oa := (*VALUE)(unsafe.Pointer(&aa))
	om := (*VALUE)(unsafe.Pointer(&am))
	op := (*VALUE)(unsafe.Pointer(&ap))
	os := (*VALUE)(unsafe.Pointer(&as))
	od := (*VALUE)(unsafe.Pointer(&ad))
	ox := (*VALUE)(unsafe.Pointer(&ax))

	fmt.Println("PTR:", oa.ptr, fmt.Sprintf("%#v", oa.Interface()), ValueOf(*(*stringHeader)(oa.ptr)))
	fmt.Println("PTR:", om.ptr) //, fmt.Sprintf("%#v", *(*map[string]string)(om.ptr)))
	fmt.Println("PTR:", op.ptr) //, fmt.Sprintf("%#v", **(**string)(op.ptr)))
	fmt.Println("PTR:", os.ptr, fmt.Sprintf("%#v", *(*[]string)(os.ptr)), ValueOf(*(*sliceHeader)(os.ptr)), (*sliceHeader)(os.ptr).Data)
	fmt.Println("PTR:", od.ptr, fmt.Sprintf("%#v", *(*strukt)(od.ptr)))
	fmt.Println("PTR:", ox.ptr, fmt.Sprintf("%#v", *(*strukt)(ox.ptr)))
	fmt.Println(*(*unsafe.Pointer)(ox.ptr))
	fmt.Println()

	na := oa.New()
	nm := om.New()
	np := op.New()
	ns := os.New()
	nd := od.New()
	nx := ox.New()

	fmt.Println("PTR:", na.ptr, fmt.Sprintf("%#v", *(*[2]string)(na.ptr)))
	fmt.Println("PTR:", nm.ptr, fmt.Sprintf("%#v", *(*map[string]string)(nm.ptr)))
	fmt.Println("PTR:", np.ptr, fmt.Sprintf("%#v", **(**string)(np.ptr)))
	fmt.Println("PTR:", ns.ptr, fmt.Sprintf("%#v", *(*[]string)(ns.ptr)), ValueOf(*(*sliceHeader)(ns.ptr)), (*sliceHeader)(ns.ptr).Data)
	fmt.Println("PTR:", nd.ptr, fmt.Sprintf("%#v", *(*strukt)(nd.ptr)))
	fmt.Println("PTR:", nx.ptr, fmt.Sprintf("%#v", *(*strukt)(nx.ptr)))
	fmt.Printf("%#v\n", *(*string)(nx.ptr))
}

func TestPtr(t *testing.T) {
	type strukt struct {
		V1 string
		V2 string
	}
	var (
		b   = true
		i   = 12345
		s   = "true"
		a   = [2]string{"false", "true"}
		m   = map[string]string{"0": "false", "1": "true"}
		p   = &s
		l   = []string{"false", "true"}
		d   = strukt{"false", "true"}
		mm  = &map[string]string{"0": "false", "1": "true"}
		ll  = &l
		lll = &ll
		sss = &p
		a1  = [1]map[string]string{{"1": "true"}}
	)
	vb := ValueOf(b)
	vi := ValueOf(i)
	vs := ValueOf(s)
	va := ValueOf(a)
	vm := ValueOf(m)
	vp := ValueOf(p)
	vl := ValueOf(l)
	vd := ValueOf(d)
	vpm := ValueOf(mm)
	vppl := ValueOf(lll)
	vps := ValueOf(sss)
	va1 := ValueOf(a1)

	fmt.Println("TYPE:", vb.typ, "PTR:", vb.ptr, "VAL:", *(*bool)(vb.ptr), "IFACE:", vb.Interface())
	fmt.Println("TYPE:", vi.typ, "PTR:", vi.ptr, "VAL:", *(*int)(vi.ptr), "IFACE:", vi.Interface())
	fmt.Println("TYPE:", vs.typ, "PTR:", vs.ptr, "VAL:", *(*string)(vs.ptr), "IFACE:", vs.Interface())
	fmt.Println("TYPE:", va.typ, "PTR:", va.ptr, "VAL:", *(*[2]string)(va.ptr), "IFACE:", va.Interface())
	fmt.Println("TYPE:", vm.typ, "PTR:", vm.ptr, "VAL:", *(*map[string]string)(vm.ptr), "IFACE:", vm.Interface())
	fmt.Println("TYPE:", vp.typ, "PTR:", vp.ptr, "VAL:", **(**string)(vp.ptr), "IFACE:", vp.Elem().Interface())
	fmt.Println("TYPE:", vl.typ, "PTR:", vl.ptr, "VAL:", *(*[]string)(vl.ptr), "IFACE:", vl.Interface())
	fmt.Println("TYPE:", vd.typ, "PTR:", vd.ptr, "VAL:", *(*strukt)(vd.ptr), "IFACE:", vd.Interface())
	fmt.Println("TYPE:", vpm.typ, "PTR:", vpm.ptr, "VAL:", **(**map[string]string)(vpm.ptr), "IFACE:", vpm.Interface())
	//fmt.Println(pfactor(vp.typ), pfactor(vpm.typ), pfactor(vppl.typ))
	fmt.Println("TYPE:", vppl.typ, "PTR:", vppl.ptr, "VAL:", ***(***[]string)(vppl.ptr), "IFACE:", fmt.Sprintf("%#v", vppl.Interface()))
	//ptr := *(*unsafe.Pointer)(*(*unsafe.Pointer)(*(*unsafe.Pointer)(vpps.ptr)))
	fmt.Println("TYPE:", vps.typ, "PTR:", vps.ptr, "VAL:", ***(***string)(vps.ptr), "IFACE:", fmt.Sprintf("%#v", vps.Interface()))
	fmt.Println("TYPE:", va1.typ, "PTR:", va1.ptr, "VAL:", *(*[1]map[string]string)(va1.ptr), "IFACE:", va1.Interface())
	fmt.Println()

}

func TestStruct(t *testing.T) {
	type mstruct struct {
		V1 map[string]string
		V2 map[string]string
	}
	m := map[string]string{"0": "zero"}

	ms := (any)(mstruct{m, m})
	vs := *(*VALUE)(unsafe.Pointer(&ms))
	fmt.Println(vs.typ.size, *(*map[string]string)(vs.ptr))

	ma := (any)([2]map[string]string{m, m})
	va := *(*VALUE)(unsafe.Pointer(&ma))
	fmt.Println(va.typ.size, *(*map[string]string)(va.ptr))

	type mstruct1 struct {
		V1 map[string]string
	}
	ms1 := (any)(mstruct1{m})
	vs1 := *(*VALUE)(unsafe.Pointer(&ms1))
	fmt.Println(vs1.typ.size, *(*map[string]string)(unsafe.Pointer(&vs1.ptr)))

	ma1 := (any)([1]map[string]string{m})
	va1 := *(*VALUE)(unsafe.Pointer(&ma1))
	fmt.Println(va1.typ.size, *(*map[string]string)(unsafe.Pointer(&va1.ptr)))

	n := vs1.New()
	n.Index(0).Set(m)
	fmt.Println(n)
}
