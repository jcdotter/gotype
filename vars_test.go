package gotype

type bool_struct_single struct {
	V1 bool
}

type int_struct_single struct {
	V1 int
}

type array_struct_single struct {
	V1 [1]string
}

type slice_struct_single struct {
	V1 []string
}

type map_struct_single struct {
	V1 map[string]string
}

type string_struct_single struct {
	V1 string
}

type struct_struct_single struct {
	V1 string_struct_single
}

type bool_ptr_struct_single struct {
	V1 *bool
}

type any_struct_single struct {
	V1 any
}

type int_ptr_struct_single struct {
	V1 *int
}

type string_ptr_struct_single struct {
	V1 *string
}

type array_ptr_struct_single struct {
	V1 *[1]string
}

type slice_ptr_struct_single struct {
	V1 *[]string
}

type map_ptr_struct_single struct {
	V1 *map[string]string
}

type struct_ptr_struct_single struct {
	V1 *string_struct_single
}

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

func getTestVars() map[string]any {
	return createTestVars(false, 0, "false")
}

func createTestVars(b bool, i int, s string, item ...string) map[string]any {
	v := map[string]any{
		"bool":   b,
		"int":    i,
		"string": s,

		"*bool":   &b,
		"*int":    &i,
		"*string": &s,

		"*[1]string":            &[1]string{s},
		"*[]string{1}":          &[]string{s},
		"*map[string]string{1}": &map[string]string{"0": s},
		"*struct(string){1}":    &string_struct_single{s},

		"*[2]string":         &[2]string{s, s},
		"*[]string{2}":       &[]string{s, s},
		"*map[string]string": &map[string]string{"0": s, "1": s},
		"*struct(string){2}": &string_struct{s, s},

		"[1]bool":                 [1]bool{b},
		"[1]int":                  [1]int{i},
		"[1]string":               [1]string{s},
		"[1][1]string":            [1][1]string{{s}},
		"[1][]string":             [1][]string{{s}},
		"[1]map[string]string{1}": [1]map[string]string{{"0": s}},
		"[1]struct(string){1}":    [1]string_struct_single{{s}},
		"[1]any(string)":          [1]any{s},

		"[2]bool":                 [2]bool{b, b},
		"[2]int":                  [2]int{i, i},
		"[2]string":               [2]string{s, s},
		"[2][2]string":            [2][2]string{{s, s}, {s, s}},
		"[2][1]string":            [2][1]string{{s}, {s}},
		"[2][]string":             [2][]string{{s, s}, {s, s}},
		"[2]map[string]string{2}": [2]map[string]string{{"0": s, "1": s}, {"0": s, "1": s}},
		"[2]struct(string){2}":    [2]string_struct{{s, s}, {s, s}},
		"[2]struct(string){1}":    [2]string_struct_single{{s}, {s}},
		"[2]any(string)":          [2]any{s, s},

		"[1]*bool":                 [1]*bool{&b},
		"[1]*int":                  [1]*int{&i},
		"[1]*string":               [1]*string{&s},
		"[1]*[1]string":            [1]*[1]string{{s}},
		"[1]*[]string":             [1]*[]string{{s}},
		"[1]*map[string]string{1}": [1]*map[string]string{{"0": s}},
		"[1]*struct(string){1}":    [1]*string_struct_single{{s}},
		"[1]any(*string)":          [1]any{&s},

		"[2]*bool":                 [2]*bool{&b, &b},
		"[2]*int":                  [2]*int{&i, &i},
		"[2]*string":               [2]*string{&s, &s},
		"[2]*[2]string":            [2]*[2]string{{s, s}, {s, s}},
		"[2]*[1]string":            [2]*[1]string{{s}},
		"[2]*[]string{2}":          [2]*[]string{{s, s}, {s, s}},
		"[2]*map[string]string{2}": [2]*map[string]string{{"0": s, "1": s}, {"0": s, "1": s}},
		"[2]*struct(string){2}":    [2]*string_struct{{s, s}, {s, s}},
		"[2]*struct(string){1}":    [2]*string_struct_single{{s}, {s}},
		"[2]any(*string)":          [2]any{&s, &s},

		"[]bool{1}":                []bool{b},
		"[]int{1}":                 []int{i},
		"[]string{1}":              []string{s},
		"[][1]string{1}":           [][1]string{{s}},
		"[][]string{1,1}":          [][]string{{s}},
		"[]map[string]string{1,1}": []map[string]string{{"0": s}},
		"[]struct(string){1,1}":    []string_struct_single{{s}},
		"[]any(string){1}":         []any{s},

		"[]bool{2}":                []bool{b, b},
		"[]int{2}":                 []int{i, i},
		"[]string{2}":              []string{s, s},
		"[][2]string{2}":           [][2]string{{s, s}, {s, s}},
		"[][]string{2,2}":          [][]string{{s, s}, {s, s}},
		"[]map[string]string{2,2}": []map[string]string{{"0": s, "1": s}, {"0": s, "1": s}},
		"[]struct(string){2,2}":    []string_struct{{s, s}, {s, s}},
		"[]any(string){2}":         []any{s, s},

		"[]*bool{1}":                []*bool{&b},
		"[]*int{1}":                 []*int{&i},
		"[]*string{1}":              []*string{&s},
		"[][1]*string{1}":           [][1]*string{{&s}},
		"[][]*string{1,1}":          [][]*string{{&s}},
		"[]*map[string]string{1,1}": []*map[string]string{{"0": s}},
		"[]*struct(string){1,1}":    []*string_struct_single{{s}},
		"[]any(*string){1}":         []any{&s},

		"[]*bool{2}":                []*bool{&b, &b},
		"[]*int{2}":                 []*int{&i, &i},
		"[]*string{2}":              []*string{&s, &s},
		"[][2]*string{2}":           [][2]*string{{&s, &s}, {&s, &s}},
		"[][]*string{2,2}":          [][]*string{{&s, &s}, {&s, &s}},
		"[]*map[string]string{2,2}": []*map[string]string{{"0": s, "1": s}, {"0": s, "1": s}},
		"[]*string_struct{2,2}":     []*string_struct{{s, s}, {s, s}},
		"[]any(*string){2}":         []any{&s, &s},

		"map[string]bool{1}":                map[string]bool{"0": b},
		"map[string]int{1}":                 map[string]int{"0": i},
		"map[string]string{1}":              map[string]string{"0": s},
		"map[string][1]string{1}":           map[string][1]string{"0": {s}},
		"map[string][]string{1,1}":          map[string][]string{"0": {s}},
		"map[string]map[string]string{1,1}": map[string]map[string]string{"0": {"0": s}},
		"map[string]struct(string){1,1}":    map[string]string_struct_single{"0": {s}},
		"map[string]any(string){1}":         map[string]any{"0": s},

		"map[string]bool{2}":                map[string]bool{"0": b, "1": b},
		"map[string]int{2}":                 map[string]int{"0": i, "1": i},
		"map[string]string{2}":              map[string]string{"0": s, "1": s},
		"map[string][2]string{2}":           map[string][2]string{"0": {s, s}, "1": {s, s}},
		"map[string][]string{2,2}":          map[string][]string{"0": {s, s}, "1": {s, s}},
		"map[string]map[string]string{2,2}": map[string]map[string]string{"0": {"0": s, "1": s}, "1": {"0": s, "1": s}},
		"map[string]struct(string){2,2}":    map[string]string_struct{"0": {s, s}, "1": {s, s}},
		"map[string]any(string){2}":         map[string]any{"0": s, "1": s},

		"map[string]*bool{1}":              map[string]*bool{"0": &b},
		"map[string]*int{1}":               map[string]*int{"0": &i},
		"map[string]*string{1}":            map[string]*string{"0": &s},
		"map[string]*[1]string{1}":         map[string]*[1]string{"0": {s}},
		"map[string]*[]string{1}":          map[string]*[]string{"0": {s}},
		"map[string]*map[string]string{1}": map[string]*map[string]string{"0": {"0": s}},
		"map[string]*struct(string){1}":    map[string]*string_struct_single{"0": {s}},
		"map[string]any(*string){1}":       map[string]any{"0": &s},

		"map[string]*bool{2}":               map[string]*bool{"0": &b, "1": &b},
		"map[string]*int{2}":                map[string]*int{"0": &i, "1": &i},
		"map[string]*string{2}":             map[string]*string{"0": &s, "1": &s},
		"map[string]*[2]string{2}":          map[string]*[2]string{"0": {s, s}, "1": {s, s}},
		"map[string]*[]string{2,2}":         map[string]*[]string{"0": {s, s}, "1": {s, s}},
		"ap[string]*map[string]string{2,2}": map[string]*map[string]string{"0": {"0": s, "1": s}, "1": {"0": s, "1": s}},
		"map[string]*struct(string){2}":     map[string]*string_struct{"0": {s, s}, "1": {s, s}},
		"map[string]any(*string){2}":        map[string]any{"0": &s, "1": &s},

		"struct(bool){1}":                 bool_struct_single{b},
		"struct(int){1}":                  int_struct_single{i},
		"struct(string){1}":               string_struct_single{s},
		"struct([1]string){1}":            array_struct_single{[1]string{s}},
		"struct([]string{1}){1}":          slice_struct_single{[]string{s}},
		"struct(map[string]string{1}){1}": map_struct_single{map[string]string{"0": s}},
		"struct(struct(string){1}){1}":    struct_struct_single{string_struct_single{s}},
		"struct(any(string)){1}":          any_struct_single{s},

		"struct(bool){2}":                 bool_struct{b, b},
		"struct(int){2}":                  int_struct{i, i},
		"struct(string){2}":               string_struct{s, s},
		"struct([2]string){2}":            array_struct{[2]string{s, s}, [2]string{s, s}},
		"struct([]string{2}){2}":          slice_struct{[]string{s, s}, []string{s, s}},
		"struct(map[string]string{2}){2}": map_struct{map[string]string{"0": s, "1": s}, map[string]string{"0": s, "1": s}},
		"struct(struct(string){2}){2}":    struct_struct{string_struct{s, s}, string_struct{s, s}},
		"struct(any(string)){2}":          any_struct{s, s},

		"struct(*bool){1}":                 bool_ptr_struct_single{&b},
		"struct(*int){1}":                  int_ptr_struct_single{&i},
		"struct(*string){1}":               string_ptr_struct_single{&s},
		"struct(*[1]string){1}":            array_ptr_struct_single{&[1]string{s}},
		"struct(*[]string{1}){1}":          slice_ptr_struct_single{&[]string{s}},
		"struct(*map[string]string{1}){1}": map_ptr_struct_single{&map[string]string{"0": s}},
		"struct(*struct(string){1}){1}":    struct_ptr_struct_single{&string_struct_single{s}},
		"struct(any(*string)){1}":          any_struct_single{&s},

		"struct(*bool){2}":                 bool_ptr_struct{&b, &b},
		"struct(*int){2}":                  int_ptr_struct{&i, &i},
		"struct(*string){2}":               string_ptr_struct{&s, &s},
		"struct(*[2]string){2}":            array_ptr_struct{&[2]string{s, s}, &[2]string{s, s}},
		"struct(*[]string{2}){2}":          slice_ptr_struct{&[]string{s, s}, &[]string{s, s}},
		"struct(*map[string]string{2}){2}": map_ptr_struct{&map[string]string{"0": s, "1": s}, &map[string]string{"0": s, "1": s}},
		"struct(*struct(string){2}){2}":    struct_ptr_struct{&string_struct{s, s}, &string_struct{s, s}},
		"struct(*any(string)){2}":          any_struct{&s, &s},
	}

	if len(item) > 0 {
		r := map[string]any{}
		for _, v := range item {
			r[v] = v
		}
		return r
	}
	return v
}
