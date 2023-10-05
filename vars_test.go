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

func createTestVars() map[string]any {
	var (
		b = true
		i = 1
		s = "true"
	)

	return map[string]any{
		"bool":   b,
		"int":    i,
		"string": s,
		"array":  [2]string{s, s},
		"slice":  []string{s, s},
		"map":    map[string]string{"0": s, "1": s},
		"struct": string_struct{s, s},

		"ptr_bool":   &b,
		"ptr_int":    &i,
		"ptr_string": &s,

		"ptr_array_single":  &[1]string{s},
		"ptr_slice_single":  &[]string{s},
		"ptr_map_single":    &map[string]string{"0": s},
		"ptr_struct_single": &string_struct_single{s},

		"ptr_array":  &[2]string{s, s},
		"ptr_slice":  &[]string{s, s},
		"ptr_map":    &map[string]string{"0": s, "1": s},
		"ptr_struct": &string_struct{s, s},
		"ptr_any":    (any)(&s),

		"array_bool_single":   [1]bool{b},
		"array_int_single":    [1]int{i},
		"array_string_single": [1]string{s},
		"array_array_single":  [1][1]string{{s}},
		"array_slice_single":  [1][]string{{s}},
		"array_map_single":    [1]map[string]string{{"0": s}},
		"array_struct_single": [1]string_struct_single{{s}},
		"array_any_single":    [1]any{s},

		"array_bool":   [2]bool{b, b},
		"array_int":    [2]int{i, i},
		"array_string": [2]string{s, s},
		"array_array":  [2][2]string{{s, s}, {s, s}},
		"array_slice":  [2][]string{{s, s}, {s, s}},
		"array_map":    [2]map[string]string{{"0": s, "1": s}, {"0": s, "1": s}},
		"array_struct": [2]string_struct{{s, s}, {s, s}},
		"array_any":    [2]any{s, s},

		"array_ptr_bool_single":   [1]*bool{&b},
		"array_ptr_int_single":    [1]*int{&i},
		"array_ptr_string_single": [1]*string{&s},
		"array_ptr_array_single":  [1]*[1]string{&[1]string{s}},
		"array_ptr_slice_single":  [1]*[]string{&[]string{s}},
		"array_ptr_map_single":    [1]*map[string]string{&map[string]string{"0": s}},
		"array_ptr_struct_single": [1]*string_struct_single{&string_struct_single{s}},
		"array_ptr_any_single":    [1]any{&s},

		"array_ptr_bool":   [2]*bool{&b, &b},
		"array_ptr_int":    [2]*int{&i, &i},
		"array_ptr_string": [2]*string{&s, &s},
		"array_ptr_array":  [2]*[2]string{&[2]string{s, s}, &[2]string{s, s}},
		"array_ptr_slice":  [2]*[]string{&[]string{s, s}, &[]string{s, s}},
		"array_ptr_map":    [2]*map[string]string{&map[string]string{"0": s, "1": s}, &map[string]string{"0": s, "1": s}},
		"array_ptr_struct": [2]*string_struct{&string_struct{s, s}, &string_struct{s, s}},
		"array_ptr_any":    [2]any{&s, &s},

		"slice_bool_single":   []bool{b},
		"slice_int_single":    []int{i},
		"slice_string_single": []string{s},
		"slice_array_single":  [][1]string{{s}},
		"slice_slice_single":  [][]string{{s}},
		"slice_map_single":    []map[string]string{{"0": s}},
		"slice_struct_single": []string_struct_single{{s}},
		"slice_any_single":    []any{s},

		"slice_bool":   []bool{b, b},
		"slice_int":    []int{i, i},
		"slice_string": []string{s, s},
		"slice_array":  [][2]string{{s, s}, {s, s}},
		"slice_slice":  [][]string{{s, s}, {s, s}},
		"slice_map":    []map[string]string{{"0": s, "1": s}, {"0": s, "1": s}},
		"slice_struct": []string_struct{{s, s}, {s, s}},
		"slice_any":    []any{s, s},

		"slice_ptr_bool_single":   []*bool{&b},
		"slice_ptr_int_single":    []*int{&i},
		"slice_ptr_string_single": []*string{&s},
		"slice_ptr_array_single":  [][1]*string{{&s}},
		"slice_ptr_slice_single":  [][]*string{{&s}},
		"slice_ptr_map_single":    []*map[string]string{&map[string]string{"0": s}},
		"slice_ptr_struct_single": []*string_struct_single{&string_struct_single{s}},
		"slice_ptr_any_single":    []any{&s},

		"map_bool_single":   map[string]bool{"0": b},
		"map_int_single":    map[string]int{"0": i},
		"map_string_single": map[string]string{"0": s},
		"map_array_single":  map[string][1]string{"0": {s}},
		"map_slice_single":  map[string][]string{"0": {s}},
		"map_map_single":    map[string]map[string]string{"0": {"0": s}},
		"map_struct_single": map[string]string_struct_single{"0": {s}},
		"map_any_single":    map[string]any{"0": s},

		"map_bool":   map[string]bool{"0": b, "1": b},
		"map_int":    map[string]int{"0": i, "1": i},
		"map_string": map[string]string{"0": s, "1": s},
		"map_array":  map[string][2]string{"0": {s, s}, "1": {s, s}},
		"map_slice":  map[string][]string{"0": {s, s}, "1": {s, s}},
		"map_map":    map[string]map[string]string{"0": {"0": s, "1": s}, "1": {"0": s, "1": s}},
		"map_struct": map[string]string_struct{"0": {s, s}, "1": {s, s}},
		"map_any":    map[string]any{"0": s, "1": s},

		"map_ptr_bool_single":   map[string]*bool{"0": &b},
		"map_ptr_int_single":    map[string]*int{"0": &i},
		"map_ptr_string_single": map[string]*string{"0": &s},
		"map_ptr_array_single":  map[string]*[1]string{"0": &[1]string{s}},
		"map_ptr_slice_single":  map[string]*[]string{"0": &[]string{s}},
		"map_ptr_map_single":    map[string]*map[string]string{"0": &map[string]string{"0": s}},
		"map_ptr_struct_single": map[string]*string_struct_single{"0": &string_struct_single{s}},
		"map_ptr_any_single":    map[string]any{"0": &s},

		"map_ptr_bool":   map[string]*bool{"0": &b, "1": &b},
		"map_ptr_int":    map[string]*int{"0": &i, "1": &i},
		"map_ptr_string": map[string]*string{"0": &s, "1": &s},
		"map_ptr_array":  map[string]*[2]string{"0": &[2]string{s, s}, "1": &[2]string{s, s}},
		"map_ptr_slice":  map[string]*[]string{"0": &[]string{s, s}, "1": &[]string{s, s}},
		"map_ptr_map":    map[string]*map[string]string{"0": &map[string]string{"0": s, "1": s}, "1": &map[string]string{"0": s, "1": s}},
		"map_ptr_struct": map[string]*string_struct{"0": &string_struct{s, s}, "1": &string_struct{s, s}},
		"map_ptr_any":    map[string]any{"0": &s, "1": &s},

		"struct_bool_single":   bool_struct_single{b},
		"struct_int_single":    int_struct_single{i},
		"struct_string_single": string_struct_single{s},
		"struct_array_single":  array_struct_single{[1]string{s}},
		"struct_slice_single":  slice_struct_single{[]string{s}},
		"struct_map_single":    map_struct_single{map[string]string{"0": s}},
		"struct_struct_single": struct_struct_single{string_struct_single{s}},
		"struct_any_single":    any_struct_single{s},

		"struct_bool":   bool_struct{b, b},
		"struct_int":    int_struct{i, i},
		"struct_string": string_struct{s, s},
		"struct_array":  array_struct{[2]string{s, s}, [2]string{s, s}},
		"struct_slice":  slice_struct{[]string{s, s}, []string{s, s}},
		"struct_map":    map_struct{map[string]string{"0": s, "1": s}, map[string]string{"0": s, "1": s}},
		"struct_struct": struct_struct{string_struct{s, s}, string_struct{s, s}},
		"struct_any":    any_struct{s, s},

		"struct_ptr_bool_single":   bool_ptr_struct_single{&b},
		"struct_ptr_int_single":    int_ptr_struct_single{&i},
		"struct_ptr_string_single": string_ptr_struct_single{&s},
		"struct_ptr_array_single":  array_ptr_struct_single{&[1]string{s}},
		"struct_ptr_slice_single":  slice_ptr_struct_single{&[]string{s}},
		"struct_ptr_map_single":    map_ptr_struct_single{&map[string]string{"0": s}},
		"struct_ptr_struct_single": struct_ptr_struct_single{&string_struct_single{s}},
		"struct_ptr_any_single":    any_struct_single{&s},

		"struct_ptr_bool":   bool_ptr_struct{&b, &b},
		"struct_ptr_int":    int_ptr_struct{&i, &i},
		"struct_ptr_string": string_ptr_struct{&s, &s},
		"struct_ptr_array":  array_ptr_struct{&[2]string{s, s}, &[2]string{s, s}},
		"struct_ptr_slice":  slice_ptr_struct{&[]string{s, s}, &[]string{s, s}},
		"struct_ptr_map":    map_ptr_struct{&map[string]string{"0": s, "1": s}, &map[string]string{"0": s, "1": s}},
		"struct_ptr_struct": struct_ptr_struct{&string_struct{s, s}, &string_struct{s, s}},
		"struct_ptr_any":    any_struct{&s, &s},
	}

}
