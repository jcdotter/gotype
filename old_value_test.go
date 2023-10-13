// Copyright 2023 escend llc. All rights reserved.typVal
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

type struct_type struct {
	Zero string `json:"0" sub:"json:'0'"`
	One  string `json:"1" sub:"json:'1'"`
}

type struct_type_ struct {
	Zero string `json:"0" sub:"json:'0'"`
	One  string `json:"1" sub:"json:'1'"`
}

type struct_complex struct {
	Bool    bool
	Float   float64
	Time    TIME
	UUID    UUID
	String  string
	Array   []any
	Map     map[string]any
	Struct  struct_type
	Struct_ struct_type_
	Json    JSON
}

type strukt struct {
	Name string
	Subs []*substrukt
	SubI map[string]*substrukt
}

type substrukt struct {
	Name   string
	Struct *strukt
}

var (
	bool_bool         = true
	int_int           = 1
	float_float       = 1.0
	string_string     = "1"
	string_bool       = "true"
	string_time       = "1970-01-01 00:00:00.000000001"
	string_time1      = "Jan 1, 70T00:00:00.000000001"
	string_uuid       = "267b3229-2566-4426-a826-8d80126e719a"
	string_array      = `["0","1"]`
	string_json       = `{"0":"0","1":"1"}`
	byte_byte         = []byte("1")
	time_time         = TIME(time.Unix(0, 1).UTC())
	uuid_uuid         = UUID([16]byte{0x26, 0x7b, 0x32, 0x29, 0x25, 0x66, 0x44, 0x26, 0xa8, 0x26, 0x8d, 0x80, 0x12, 0x6e, 0x71, 0x9a})
	array_array       = []any{"0", "1"}
	map_map           = map[string]any{"0": "0", "1": "1"}
	struct_struct_old = struct_type{"0", "1"}
	struct_blank      = struct_type_{}
	struct_cmp        = struct_complex{bool_bool, float_float, time_time, uuid_uuid, string_string, array_array, map_map, struct_struct_old, struct_blank, json_array}
	json_json         = JSON(string_json)
	json_array        = JSON(string_array)
	string_serl       = `db:"dbtest" json:"jsontest" custom:"testing; notull; unique"`
	kind_kind         = KIND(0)

	time_val  = ValueOf(time_time)
	ttime_val = ValueOf(time.Time{})
)

func Test_All(t *testing.T) {
	//TestAddrs(t)
	//TestEmptyAddrs(t)
	TestValidations(t)
	TestSerialValidations(t)
	TestMapValidations(t)
	TestEncoding(t)
	//TestComplexEncoding(t)
}

func TestValidations(t *testing.T) {
	var validations = []validation_test{
		{"Conversion: Bool to Bool", ValueOf(bool_bool).Bool(), bool_bool},
		{"Conversion: Int to Bool", ValueOf(int_int).Bool(), bool_bool},
		{"Conversion: Float to Bool", ValueOf(float_float).Bool(), bool_bool},
		{"Conversion: String to Bool", ValueOf(string_string).Bool(), bool_bool},
		{"Conversion: Bytes to Bool", ValueOf(byte_byte).Bool(), bool_bool},
		{"Conversion: Time to Bool", ValueOf(time_time).Bool(), bool_bool},
		{"Conversion: UUID to Bool", ValueOf(uuid_uuid).Bool(), bool_bool},

		{"Conversion: Bool to Int", ValueOf(bool_bool).Int(), int_int},
		{"Conversion: Int to Int", ValueOf(int_int).Int(), int_int},
		{"Conversion: Float to Int", ValueOf(float_float).Int(), int_int},
		{"Conversion: String to Int", ValueOf(string_string).Int(), int_int},
		{"Conversion: Bytes to Int", ValueOf(byte_byte).Int(), int_int},
		{"Conversion: Time to Int", ValueOf(time_time).Int(), int_int},

		{"Conversion: Bool to Float", ValueOf(bool_bool).Float64(), float_float},
		{"Conversion: Int to Float", ValueOf(int_int).Float64(), float_float},
		{"Conversion: Float to Float", ValueOf(float_float).Float64(), float_float},
		{"Conversion: String to Float", ValueOf(string_string).Float64(), float_float},
		{"Conversion: Bytes to Float", ValueOf(byte_byte).Float64(), float_float},
		{"Conversion: Time to Float", ValueOf(time_time).Float64(), float_float},

		{"Conversion: Bool to String", ValueOf(bool_bool).String(), string_bool},
		{"Conversion: Int to String", ValueOf(int_int).String(), string_string},
		{"Conversion: Float to String", ValueOf(float_float).String(), string_string},
		{"Conversion: String to String", ValueOf(string_string).String(), string_string},
		{"Conversion: Bytes to String", ValueOf(byte_byte).String(), string_string},
		{"Conversion: Time to String", ValueOf(time_time).TIME().Format(ISO8601N), string_time},
		{"Conversion: UUID to String", ValueOf(uuid_uuid).String(), string_uuid},
		{"Conversion: Slice to String", ValueOf(array_array).String(), string_array},
		{"Conversion: Map to String", ValueOf(map_map).String(), string(json_json)},
		{"Conversion: Struct to String", ValueOf(struct_struct_old).STRUCT().JsonByTag("json").String(), string(json_json)},

		{"Conversion: String to Bytes", ValueOf(string_string).Bytes(), byte_byte},
		{"Conversion: Bytes to Bytes", ValueOf(byte_byte).Bytes(), byte_byte},

		{"Conversion: Int to Time", ValueOf(int_int).TIME(), time_time},
		{"Conversion: Uint to Time", ValueOf(uint(int_int)).TIME(), time_time},
		{"Conversion: Float to Time", ValueOf(float_float).TIME(), time_time},
		{"Conversion: Float32 to Time", ValueOf(float32(float_float)).TIME(), time_time},
		{"Conversion: String to Time", ValueOf(string_time).TIME(), time_time},
		{"Conversion: Time to Time", ValueOf(time_time).TIME(), time_time},
		{"Conversion: Complex String to Time", STRING(string_time1).ParseTime(), time_time},

		{"Conversion: String to UUID", ValueOf(string_uuid).UUID(), uuid_uuid},
		{"Conversion: UUID to UUID", ValueOf(uuid_uuid).UUID(), uuid_uuid},

		{"Conversion: Slice to Slice", ValueOf(array_array).SLICE().Interface(), array_array},
		{"Conversion: Map to Slice", ValueOf(map_map).SLICE().SortStrings(), SliceOf(array_array).Strings()},
		{"Conversion: Struct to Slice", ValueOf(struct_struct_old).SLICE().Interface(), array_array},
		{"Conversion: Struct to Strings Slice", ValueOf(struct_struct_old).STRUCT().Strings(), []string{"0", "1"}},
		{"Conversion: Json to Slice", ValueOf(json_array).JSON().Slice(), array_array},

		{"Conversion: Slice to Map", ValueOf(array_array).MAP().Map(), map_map},
		{"Conversion: Map to Map", ValueOf(map_map).MAP().Map(), map_map},
		{"Conversion: Struct to Map", ValueOf(struct_struct_old).STRUCT().MapByTag("json"), map_map},
		{"Conversion: Json to Map", ValueOf(json_json).MAP().Index("1").String(), map_map["1"].(string)},

		//{"Conversion: Slice to Struct", ValueOf(array_array).SLICE().StructScan(StructOf(struct_type{})).Interface(), struct_struct_old},
		//{"Conversion: Map to Struct", ValueOf(map_map).MAP().StructScan(StructOf(struct_type{}), "json", "").Interface(), struct_struct_old},
		{"Conversion: Struct to Struct", ValueOf(struct_struct_old).STRUCT().Interface(), struct_struct_old},
		//{"Conversion: Json to Struct", ValueOf(json_json).JSON().StructScan(StructOf(struct_type{}), "json", "").Interface(), struct_struct_old},
		{"Conversion: Empty Struct to Struct", ValueOf(&struct_blank).Elem().Set(struct_struct_old).STRUCT().MapByTag("json"), map_map},
		{"Conversion: Struct to Json", ValueOf(struct_struct_old).STRUCT().SerializeByTag("json"), string_json},
		{"Conversion: Match Type - Time", matchStructType(time_val.typ, ttime_val.typ), true},
		{"Conversion: Struct to Json to Slice to Value", JSON(JSON(StructOf(struct_cmp).Serialize()).MAP().Index("Json").String()).Slice()[1], "1"},

		{"Conversion: Slice to Json", ValueOf(array_array).JSON(), json_array},
		{"Conversion: Map to Json", ValueOf(map_map).JSON(), json_json},
		{"Conversion: Struct to Json", ValueOf(struct_struct_old).STRUCT().JsonByTag("json"), json_json},
		{"Conversion: Json to Json", ValueOf(json_json).JSON(), json_json},
	}
	if err := test_validations(validations); err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	validations = []validation_test{
		{"Set Value: Set Bool Pointer", ValueOf(&bool_bool).Elem().Set(false).Interface(), false},
		{"Set Value: Set Struct Pointer", StructOf(&struct_struct_old).Index(0).Set(2).Interface(), "2"},
		{"Set Value: Append to Slice", SliceOf(array_array).Append("2").Interface(), []any{"0", "1", "2"}},
	}
	if err := test_validations(validations); err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
}

func TestSerialValidations(t *testing.T) {
	var (
		tBytes       = []byte{0x31}
		tBool        = true
		tInt         = 1
		tString      = "1"
		tSlice       = []string{tString}
		tMap         = map[string]string{"string": tString}
		tStruct      = struct_type{"0", "1"}
		tSuperStruct = strukt{
			Name: "Test",
		}
		tSubStruct = substrukt{
			Name:   "SubStrukt",
			Struct: &tSuperStruct,
		}

		tSliceAny        = []any{tBytes, tBool, tInt, tString, tSlice, tMap, tStruct}
		tSliceAnyPtrs    = []any{&tBytes, &tBool, &tInt, &tString, &tSlice, &tMap, &tStruct, nil}
		tSliceStructPtrs = []*substrukt{&tSubStruct, &tSubStruct}

		tMapAny        = map[string]any{"bytes": tBytes, "bool": tBool, "int": tInt, "string": tString, "slice": tSlice, "map": tMap, "struct": tStruct}
		tMapAnyPtrs    = map[string]any{"bytes": &tBytes, "bool": &tBool, "int": &tInt, "string": &tString, "slice": &tSlice, "map": &tMap, "struct": &tStruct}
		tMapStructPtrs = map[string]*substrukt{"sub1": &tSubStruct, "sub2": &tSubStruct}
	)
	tSliceAnyPtrs[7] = &tSliceAnyPtrs
	tMapAnyPtrs["recursive"] = &tMapAnyPtrs
	tSubStruct.Struct.Subs = tSliceStructPtrs
	tSubStruct.Struct.SubI = tMapStructPtrs

	var validations = []validation_test{
		{"Serialize: Bytes", ValueOf(tBytes).Serialize(), `"1"`},
		{"Serialize: Bool", ValueOf(tBool).Serialize(), `true`},
		{"Serialize: Int", ValueOf(tInt).Serialize(), `1`},
		{"Serialize: String", ValueOf(tString).Serialize(), `"1"`},
		{"Serialize: Slice", ValueOf(tSlice).Serialize(), `["1"]`},
		{"Serialize: Map", ValueOf(tMap).Serialize(), `{"string":"1"}`},
		{"Serialize: Struct", ValueOf(tStruct).Serialize(), `{"Zero":"0","One":"1"}`},

		{"Serialize: Slice of Any", ValueOf(tSliceAny).Serialize(), `["1",true,1,"1",["1"],{"string":"1"},{"Zero":"0","One":"1"}]`},
		{"Serialize: Slice of Pointers", ValueOf(&tSliceAnyPtrs).Serialize(), `["1",true,1,"1",["1"],{"string":"1"},{"Zero":"0","One":"1"},"*recursive"]`},

		{"Serialize: Json compare on Map of Any", JSON(ValueOf(tMapAny).Serialize()).Map(), JSON(`{"bytes":"1","bool":true,"int":1,"string":"1","slice":["1"],"map":{"string":"1"},"struct":{"Zero":"0","One":"1"}}`).Map()},
		{"Serialize: Json compare on Map of Pointers", JSON(ValueOf(&tMapAnyPtrs).Serialize()).Map(), JSON(`{"bytes":"1","bool":true,"int":1,"string":"1","slice":["1"],"map":{"string":"1"},"struct":{"Zero":"0","One":"1"},"recursive":"*recursive"}`).Map()},

		{"Serialize: Json compare on Struct", JSON(ValueOf(&tSuperStruct).Serialize()).Map(), JSON(`{"Name":"Test","Subs":[{"Name":"SubStrukt","Struct":"*recursive"},{"Name":"SubStrukt","Struct":"*recursive"}],"SubI":{"sub1":{"Name":"SubStrukt","Struct":"*recursive"},"sub2":{"Name":"SubStrukt","Struct":"*recursive"}}}`).Map()},
	}
	if err := test_validations(validations); err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
}

func TestMapValidations(t *testing.T) {
	type mstruct struct {
		Hmap       map[string]string
		MapPointer *map[string]string
		MapList    []map[string]string
		MapMap     map[string]map[string]string
		MapAny     any
	}
	var (
		emap          map[string]string
		emap_pointer  *map[string]string
		emap_list     []map[string]string
		emap_map      map[string]map[string]string
		emap_struct   mstruct
		emap_list_any []any

		pmap      map[string]*string
		pmap_map  map[string]*map[string]string
		pstr_list []*string
		pmap_list []*map[string]string

		hmap         = map[string]string{"Zero": "0", "One": "1"}
		map_pointer  = &hmap
		map_list     = []map[string]string{hmap}
		map_map      = map[string]map[string]string{"Zero": hmap, "One": hmap}
		map_struct   = mstruct{hmap, map_pointer, map_list, map_map, map_map}
		map_list_any = []any{hmap}

		amap         = (any)(hmap)
		amap_pointer = (any)(map_pointer)
		amap_list    = (any)(map_list)
		amap_map     = (any)(map_map)
		amap_struct  = (any)(map_struct)
	)

	var validations = []validation_test{
		// Get Existing Index
		{"Map Index: Strings Map Value", ValueOf(amap).MapIndex("One").String(), "1"},
		{"Map Index: Map Pointer Value", ValueOf(amap_pointer).Elem().MapIndex("One").String(), "1"},
		{"Map Index: Strings Map", MapOf(amap).Index("One").String(), "1"},
		{"Map Index: Map Pointer", MapOf(amap_pointer).Index("One").String(), "1"},
		{"Map Index: Map List", ValueOf(amap_list).Index(0).MapIndex("One").String(), "1"},
		{"Map Index: Map Map", ValueOf(amap_map).MapIndex("One").MapIndex("One").String(), "1"},
		{"Map Index: Structs Map Value", ValueOf(amap_struct).StructField("MapPointer").VALUE().Elem().MapIndex("One").String(), "1"},
		{"Map Index: Any Map List", ValueOf(map_list_any).Index(0).MapIndex("One").String(), "1"},

		// Set map index, plus set from nil element
		{"Map Set: Empty Strings Map", MapOf(emap).Set("One", "1").Index("One").String(), "1"},
		{"Map Set: Empty Strings Map Pointer", MapOf(&emap).Set("One", "1").Index("One").String(), "1"},
		{"Map Set: Empty Map Pointer", MapOf(emap_pointer).Set("One", "1").Index("One").String(), "1"},
		{"Map Set: Empty Strings Map Value", ValueOf(emap_pointer).MAP().Set("One", "1").Index("One").String(), "1"},
		{"Map Set: Empty List of Maps", SliceOf(&emap_list).Extend(1).Index(0).SetIndex("One", "1").MapIndex("One").String(), "1"},
		{"Map Set: Empty Map of Maps", MapOf(emap_map).Set("One", hmap).Index("One").MapIndex("One").String(), "1"},
		{"Map Set: Empty Struct of Maps", StructOf(emap_struct).Field("Hmap").Set(hmap).VALUE().MapIndex("One").String(), "1"},
		{"Map Set: Empty List of Any Maps", SliceOf(&emap_list_any).Append(emap).Index(0).SetIndex("One", "1").MapIndex("One").String(), "1"},

		{"Map Set: Empty Map of Pointers", MapOf(pmap).Set("One", &string_string).Index("One").String(), "1"},
		{"Map Set: Empty Pointer of Map of Pointers", MapOf(&pmap).Set("One", &string_string).Index("One").String(), "1"},
		{"Map Set: Empty Map of Map Pointers", MapOf(pmap_map).Set("One", &hmap).Index("One").Elem().MapIndex("One").String(), "1"},
		{"Map Set: Empty Pointer of Map of Map Pointers", MapOf(&pmap_map).Set("One", &hmap).Index("One").Elem().MapIndex("One").String(), "1"},

		{"List Set: Empty Pointer of Slice of Pointers", SliceOf(&pstr_list).Set(0, &string_string).Index(0).String(), "1"},
		{"List Set: Empty Pointer of Slice of Map Pointers", SliceOf(&pmap_list).Set(0, &hmap).Index(0).Elem().MapIndex("One").String(), "1"},
	}
	if err := test_validations(validations); err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
}

func TestEncoding(t *testing.T) {
	// TEST DECODEX
	var (
		tBool   = true
		eBool   = Encode(tBool)
		dBool   = eBool.Decodex()
		tInt    = 1
		eInt    = Encode(tInt)
		dInt    = eInt.Decodex()
		tFloat  = float64(1)
		eFloat  = Encode(tFloat)
		dFloat  = eFloat.Decodex()
		tArray  = [1]string{"1"}
		eArray  = Encode(tArray)
		dArray  = eArray.Decodex()
		tMap    = map[string]string{"zero": "1"}
		eMap    = Encode(tMap)
		dMap    = eMap.Decodex()
		tSlice  = []string{"1"}
		eSlice  = Encode(tSlice)
		dSlice  = eSlice.Decodex()
		tString = "1"
		eString = Encode(tString)
		dString = eString.Decodex()
		tStruct = struct_type{"0", "1"}
		eStruct = Encode(tStruct)
		dStruct = eStruct.Decodex()
	)
	var validations = []validation_test{
		{"Encoding: Bool", eBool.Bytes(), []byte{1, 1}},
		{"Encoding: Int", eInt.Bytes(), []byte{2, 0, 0, 0, 0, 0, 0, 0, 1}},
		{"Encoding: Float", eFloat.Bytes(), []byte{14, 63, 240, 0, 0, 0, 0, 0, 0}},
		{"Encoding: Array", eArray.Bytes(), []byte{17, 24, 8, 1, 24, 8, 1, 49}},
		{"Encoding: Map", eMap.Bytes(), []byte{21, 24, 24, 8, 1, 24, 8, 4, 122, 101, 114, 111, 24, 8, 1, 49}},
		{"Encoding: Slice", eSlice.Bytes(), []byte{23, 24, 8, 1, 24, 8, 1, 49}},
		{"Encoding: String", eString.Bytes(), []byte{24, 8, 1, 49}},
		{"Encoding: Struct", eStruct.Bytes(), []byte{25, 8, 2, 24, 8, 1, 48, 24, 8, 1, 49}},
	}
	if err := test_validations(validations); err != nil {
		t.Fatalf(fmt.Sprint(err))
	}

	validations = []validation_test{
		{"Decodex: Bool Type", dBool.k, Bool},
		{"Decodex: Bool Value", dBool.v.Bool(), tBool},
		{"Decodex: Bool Bytes", dBool.b, 2},
		{"Decodex: Int Type", dInt.k, Int},
		{"Decodex: Int Value", dInt.v.Int(), tInt},
		{"Decodex: Int Bytes", dInt.b, 9},
		{"Decodex: Float Type", dFloat.k, Float64},
		{"Decodex: Float Value", dFloat.v.Float64(), tFloat},
		{"Decodex: Float Bytes", dFloat.b, 9},
		{"Decodex: Array Type", dArray.k, Array},
		{"Decodex: Array Value", dArray.v.Interface(), tArray},
		{"Decodex: Array Bytes", dArray.b, 8},
		{"Decodex: Map Type", dMap.k, Map},
		{"Decodex: Map Value", dMap.v.Interface(), tMap},
		{"Decodex: Map Bytes", dMap.b, 16},
		{"Decodex: Slice Type", dSlice.k, Slice},
		{"Decodex: Slice Value", dSlice.v.Interface(), tSlice},
		{"Decodex: Slice Bytes", dSlice.b, 8},
		{"Decodex: String Type", dString.k, String},
		{"Decodex: String Value", dString.v.String(), tString},
		{"Decodex: String Bytes", dString.b, 4},
		{"Decodex: Struct Type", dStruct.k, Struct},
		{"Decodex: Struct Value", dStruct.v.Interface(), []any{"0", "1"}},
		{"Decodex: Struct Bytes", dStruct.b, 11},
	}
	if err := test_validations(validations); err != nil {
		t.Fatalf(fmt.Sprint(err))
	}

	// TEST DECODE TO VAR POINTER
	var (
		rBool   = false
		rInt    = 0
		rFloat  = 0.0
		rArray  = [1]string{}
		rMap    = map[string]string{}
		rSlice  = []string{}
		rString = ""
		rStruct = struct_type{}
	)
	Decode(eBool, &rBool)
	Decode(eInt, &rInt)
	Decode(eFloat, &rFloat)
	Decode(eArray, &rArray)
	Decode(eMap, &rMap)
	Decode(eSlice, &rSlice)
	Decode(eString, &rString)
	Decode(eStruct, &rStruct)

	// TEST DECODE TO POINTER
	validations = []validation_test{
		{"Decoding: Bool", rBool, tBool},
		{"Decoding: Int", rInt, tInt},
		{"Decoding: Float", rFloat, tFloat},
		{"Decoding: Array", rArray, tArray},
		{"Decoding: Map", rMap, tMap},
		{"Decoding: Slice", rSlice, tSlice},
		{"Decoding: String", rString, tString},
		{"Decoding: Struct", rStruct, tStruct},
	}
	if err := test_validations(validations); err != nil {
		t.Fatalf(fmt.Sprint(err))
	}

}

func TestComplexEncoding(t *testing.T) {

	// TEST COMPLEX DECODING
	type strukt struct {
		Array     [1]string
		Any       any
		Map       map[string]string
		Slice     []string
		String    string
		Struct    struct_type
		PtrArray  *[1]string
		PtrMap    *map[string]string
		PtrSlice  *[]string
		PtrString *string
		PtrStruct *struct_type
	}

	var (
		tArray  = [1]string{"1"}
		tMap    = map[string]string{"zero": "1"}
		tSlice  = []string{"1"}
		tString = "1"
		tStruct = struct_type{"0", "1"}

		array_array      = [1][1]string{}
		array_any        = [1]any{}
		array_map        = [1]map[string]string{}
		array_slice      = [1][]string{}
		array_string     = [1]string{}
		array_struct     = [1]struct_type{}
		array_ptr_array  = [1]*[1]string{}
		array_ptr_map    = [1]*map[string]string{}
		array_ptr_slice  = [1]*[]string{}
		array_ptr_string = [1]*string{}
		array_ptr_struct = [1]*struct_type{}

		/* slice_array      [][1]string
		slice_any        []any
		slice_map        []map[string]string
		slice_slice      [][]string
		slice_string     []string
		slice_struct     []struct_type
		slice_ptr_array  []*[1]string
		slice_ptr_map    []*map[string]string
		slice_ptr_slice  []*[]string
		slice_ptr_string []*string
		slice_ptr_struct []*struct_type

		map_array      map[string][1]string
		map_any        map[string]any
		map_map        map[string]map[string]string
		map_slice      map[string][]string
		map_string     map[string]string
		map_struct     map[string]struct_type
		map_ptr_array  map[string]*[1]string
		map_ptr_map    map[string]*map[string]string
		map_ptr_slice  map[string]*[]string
		map_ptr_string map[string]*string
		map_ptr_struct map[string]*struct_type

		struct_test strukt */
	)

	var (
		v_array_array      = ValueOf(array_array).New().SetIndex(0, tArray)
		v_array_any        = ValueOf(array_any).New().SetIndex(0, tMap)
		v_array_map        = ValueOf(array_map).New().SetIndex(0, tMap)
		v_array_slice      = ValueOf(array_slice).New().SetIndex(0, tSlice)
		v_array_string     = ValueOf(array_string).New().SetIndex(0, tString)
		v_array_struct     = ValueOf(array_struct).New().SetIndex(0, tStruct)
		v_array_ptr_array  = ValueOf(array_ptr_array).New().SetIndex(0, &tArray)
		v_array_ptr_map    = ValueOf(array_ptr_map).New().SetIndex(0, &tMap)
		v_array_ptr_slice  = ValueOf(array_ptr_slice).New().SetIndex(0, &tSlice)
		v_array_ptr_string = ValueOf(array_ptr_string).New().SetIndex(0, &tString)
		v_array_ptr_struct = ValueOf(array_ptr_struct).New().SetIndex(0, &tStruct)
		/* )
		fmt.Println(v_array_ptr_array.ptr, *(*unsafe.Pointer)(v_array_ptr_array.ptr), v_array_ptr_array, v_array_ptr_array.Interface())
		fmt.Println(v_array_ptr_map.ptr, *(*unsafe.Pointer)(v_array_ptr_map.ptr))
		fmt.Println(v_array_ptr_slice.ptr, *(*unsafe.Pointer)(v_array_ptr_slice.ptr))
		fmt.Println(v_array_ptr_string.ptr, *(*unsafe.Pointer)(v_array_ptr_string.ptr))
		fmt.Println(v_array_ptr_struct.ptr, *(*unsafe.Pointer)(v_array_ptr_struct.ptr))
		os.Exit(1)
		var ( */
		e_array_array      = v_array_array.Encode()
		e_array_any        = v_array_any.Encode()
		e_array_map        = v_array_map.Encode()
		e_array_slice      = v_array_slice.Encode()
		e_array_string     = v_array_string.Encode()
		e_array_struct     = v_array_struct.Encode()
		e_array_ptr_array  = v_array_ptr_array.Encode()
		e_array_ptr_map    = v_array_ptr_map.Encode()
		e_array_ptr_slice  = v_array_ptr_slice.Encode()
		e_array_ptr_string = v_array_ptr_string.Encode()
		e_array_ptr_struct = v_array_ptr_struct.Encode()
		/* )
		fmt.Println(v_array_map.ptr, *(*unsafe.Pointer)(v_array_map.ptr), v_array_map.Interface())
		//fmt.Println(v_array_ptr_map.ptr, *(*unsafe.Pointer)(v_array_ptr_map.ptr))
		os.Exit(1)
		var ( */
		d_array_array      = e_array_array.Decodex()
		d_array_any        = e_array_any.Decodex()
		d_array_map        = e_array_map.Decodex()
		d_array_slice      = e_array_slice.Decodex()
		d_array_string     = e_array_string.Decodex()
		d_array_struct     = e_array_struct.Decodex()
		d_array_ptr_array  = e_array_ptr_array.Decodex()
		d_array_ptr_map    = e_array_ptr_map.Decodex()
		d_array_ptr_slice  = e_array_ptr_slice.Decodex()
		d_array_ptr_string = e_array_ptr_string.Decodex()
		d_array_ptr_struct = e_array_ptr_struct.Decodex()
	)

	/* fmt.Println(
		fmt.Sprintf("%#v", v_array_array.Interface()), "\t", d_array_array.gVal(), "\t", e_array_array, "\n",
		fmt.Sprintf("%#v", v_array_any.Interface()), "\t", d_array_any.gVal(), "\t", e_array_any, "\n",
		fmt.Sprintf("%#v", v_array_map.Interface()), "\t", d_array_map.gVal(), "\t", e_array_map, "\n",
		fmt.Sprintf("%#v", v_array_slice.Interface()), "\t", d_array_slice.gVal(), "\t", e_array_slice, "\n",
		fmt.Sprintf("%#v", v_array_string.Interface()), "\t", d_array_string.gVal(), "\t", e_array_string, "\n",
		fmt.Sprintf("%#v", v_array_struct.Interface()), "\t", d_array_struct.gVal(), "\t", e_array_struct, "\n",
		fmt.Sprintf("%#v", v_array_ptr_array.Interface()), "\t", d_array_ptr_array.gVal(), "\t", e_array_ptr_array, "\n",
		fmt.Sprintf("%#v", v_array_ptr_map.Interface()), "\t", d_array_ptr_map.gVal(), "\t", e_array_ptr_map, "\n",
		fmt.Sprintf("%#v", v_array_ptr_slice.Interface()), "\t", d_array_ptr_slice.gVal(), "\t", e_array_ptr_slice, "\n",
		fmt.Sprintf("%#v", v_array_ptr_string.Interface()), "\t", d_array_ptr_string.gVal(), "\t", e_array_ptr_string, "\n",
		fmt.Sprintf("%#v", v_array_ptr_struct.Interface()), "\t", d_array_ptr_struct.gVal(), "\t", e_array_ptr_struct,
	) */

	Decode(e_array_array, &array_array)
	Decode(e_array_any, &array_any)
	Decode(e_array_map, &array_map)
	Decode(e_array_slice, &array_slice)
	Decode(e_array_string, &array_string)
	Decode(e_array_struct, &array_struct)
	Decode(e_array_ptr_array, &array_ptr_array)
	Decode(e_array_ptr_map, &array_ptr_map)
	Decode(e_array_ptr_slice, &array_ptr_slice)
	Decode(e_array_ptr_string, &array_ptr_string)
	//Decode(e_array_ptr_struct, &array_ptr_struct)

	fmt.Println(v_array_map, v_array_map.Interface(), d_array_map.v.Interface(), array_map)
	os.Exit(1)

	var validations = []validation_test{
		{"Encoded Array: Array:", e_array_array.Bytes(), []byte{17, 17, 8, 1, 17, 24, 8, 1, 24, 8, 1, 49}},
		{"Encoded Array: Any:", e_array_any.Bytes(), []byte{17, 20, 8, 1, 21, 24, 24, 8, 1, 24, 8, 4, 122, 101, 114, 111, 24, 8, 1, 49}},
		{"Encoded Array: Map:", e_array_map.Bytes(), []byte{17, 21, 8, 1, 21, 24, 24, 8, 1, 24, 8, 4, 122, 101, 114, 111, 24, 8, 1, 49}},
		{"Encoded Array: Slice:", e_array_slice.Bytes(), []byte{17, 23, 8, 1, 23, 24, 8, 1, 24, 8, 1, 49}},
		{"Encoded Array: String:", e_array_string.Bytes(), []byte{17, 24, 8, 1, 24, 8, 1, 49}},
		{"Encoded Array: Struct:", e_array_struct.Bytes(), []byte{17, 25, 8, 1, 25, 8, 2, 24, 8, 1, 48, 24, 8, 1, 49}},
		{"Encoded Array: Array Ptr:", e_array_ptr_array.Bytes(), []byte{17, 22, 8, 1, 17, 24, 8, 1, 24, 8, 1, 49}},
		{"Encoded Array: Map Ptr:", e_array_ptr_map.Bytes(), []byte{17, 22, 8, 1, 21, 24, 24, 8, 1, 24, 8, 4, 122, 101, 114, 111, 24, 8, 1, 49}},
		{"Encoded Array: Slice Ptr:", e_array_ptr_slice.Bytes(), []byte{17, 22, 8, 1, 23, 24, 8, 1, 24, 8, 1, 49}},
		{"Encoded Array: String Ptr:", e_array_ptr_string.Bytes(), []byte{17, 22, 8, 1, 24, 8, 1, 49}},
		{"Encoded Array: Struct Ptr:", e_array_ptr_struct.Bytes(), []byte{17, 22, 8, 1, 25, 8, 2, 24, 8, 1, 48, 24, 8, 1, 49}},

		{"Decodex Array: Array:", d_array_array.v.Index(0).Index(0).Interface(), "1"},
		{"Decodex Array: Any:", d_array_any.v.Index(0).MapIndex("zero").Interface(), "1"},
		{"Decodex Array: Map:", d_array_map.v.Index(0).MapIndex("zero").Interface(), "1"},
		{"Decodex Array: Slice:", d_array_slice.v.Index(0).Index(0).Interface(), "1"},
		{"Decodex Array: String:", d_array_string.v.Index(0).Interface(), "1"},
		{"Decodex Array: Struct:", d_array_struct.v.Index(0).Index(1).Interface(), "1"},
		{"Decodex Array: Array Ptr:", d_array_ptr_array.v.Index(0).Index(0).Interface(), "1"},
		{"Decodex Array: Map Ptr:", d_array_ptr_map.v.Index(0).MapIndex("zero").Interface(), "1"},
		{"Decodex Array: Slice Ptr:", d_array_ptr_slice.v.Index(0).Index(0).Interface(), "1"},
		{"Decodex Array: String Ptr:", d_array_ptr_string.v.Index(0).Interface(), "1"},
		{"Decodex Array: Struct Ptr:", d_array_ptr_struct.v.Index(0).Index(1).Interface(), "1"},

		{"Decode to Dest Array: Array:", array_array, v_array_array.Interface()},
		{"Decode to Dest Array: Any:", array_any, v_array_any.Interface()},
		{"Decode to Dest Array: Map:", array_map, v_array_map.Interface()},
		{"Decode to Dest Array: Slice:", array_slice, v_array_slice.Interface()},
		{"Decode to Dest Array: String:", array_string, v_array_string.Interface()},
		{"Decode to Dest Array: Struct:", array_struct, v_array_struct.Interface()},
	}

	/* fmt.Println(array_map) //v_array_map.Interface())
	os.Exit(1) */

	if err := test_validations(validations); err != nil {
		t.Fatalf(fmt.Sprint(err))
	}

	/* 	s_string := ValueOf(slice_string).Append(true).Encode()
	   	s_slice := ValueOf(slice_slice).Elem().Append(tSlice).Encode()
	   	s_map := ValueOf(slice_map).Append(tMap).Encode()
	   	sp_string := ValueOf(slice_ptr_string).Append(true).Encode()
	   	sp_slice := ValueOf(slice_ptr_slice).Append(tSlice).Encode()
	   	sp_map := ValueOf(slice_ptr_map).Append(tMap).Encode()
	   	sp_struct := ValueOf(slice_ptr_struct).Append(tStruct).Encode()
	   	sp_any := ValueOf(slice_any).Append(&tMap).Encode()

	   	fmt.Println("SLICES:")
	   	fmt.Println(s_string)
	   	fmt.Println(s_slice)
	   	fmt.Println(s_map)
	   	fmt.Println(sp_string)
	   	fmt.Println(sp_slice)
	   	fmt.Println(sp_map)
	   	fmt.Println(sp_struct)
	   	fmt.Println(sp_any)

	   	fmt.Println("MAPS:")
	   	mp_string := ValueOf(map_ptr_string).SetIndex("one", true).Encode()
	   	mp_slice := ValueOf(map_ptr_slice).SetIndex("one", tSlice).Encode()
	   	mp_map := ValueOf(map_ptr_map).SetIndex("one", tMap).Encode()
	   	mp_struct := ValueOf(map_ptr_struct).SetIndex("one", tStruct).Encode()
	   	mp_any := ValueOf(map_any).SetIndex("one", &tMap).Encode()

	   	fmt.Println(mp_string)
	   	fmt.Println(mp_slice)
	   	fmt.Println(mp_map)
	   	fmt.Println(mp_struct)
	   	fmt.Println(mp_any)

	   	fmt.Println("STRUCT:")
	   	t_all := ValueOf(struct_test).
	   		SetIndex(0, true).
	   		SetIndex(1, tSlice).
	   		SetIndex(2, tMap).
	   		SetIndex(3, tStruct).
	   		SetIndex(4, &tMap).Encode()

	   	fmt.Println(t_all) */

	/*
	   encoding from and decoding to pointers
	*/
}

func TestOldTest(t *testing.T) {
	// TODOS:
	//		[1]array set map
	//		map decode impacting other map values
	// 		pointer elem should set original value
	//		encoding array of array ptr
	//		TestEmptySingle - array_ptr_struct, slice_ptr_struct
	var (
		mdArray = [1]map[string]string{}
		m       = map[string]string{"0": "1"}
		meArray = [2]map[string]string{{"zero": "1"}}
		pdArray = [1]*[1]string{}
		tArray  = [1]string{"1"}
	)
	e := Encode(meArray).Decodex().v
	fmt.Printf("%#v\n", e.Interface())
	fmt.Println(e.ptr, *(*any)(e.ptr))

	p := (ValueOf(meArray).ptr)
	fmt.Println(p, *(*map[string]string)(p), ValueOf(meArray).Interface())
	fmt.Println()
	d := ValueOf(&meArray).Elem().SetIndex(0, map[string]string{"one": "1"})
	fmt.Println(d.ptr, *(*map[string]string)(d.ptr), d.Interface())
	fmt.Println(meArray)
	fmt.Println()

	a := ValueOf(pdArray).SetIndex(0, &tArray)
	fmt.Println(a.ptr, a.Interface(), a)
	fmt.Println(a.Encode())
	aa := (any)(pdArray)
	fmt.Println((*VALUE)(unsafe.Pointer(&aa)).ptr, pdArray)
	fmt.Println()

	ma := ValueOf(mdArray).SetIndex(0, m)
	fmt.Println(ma.ptr, ma.Interface(), ma)
	fmt.Println(ma.Encode())
	maa := (any)(&mdArray)
	fmt.Println((*VALUE)(unsafe.Pointer(&maa)).ptr, unsafe.Pointer(&mdArray), mdArray)

	/* d.Encode().Decode(&mdArray)
	fmt.Println(mdArray) */

	/* fmt.Printf("%v\n", ValueOf(&mdArray).Elem().ARRAY().index(0).MAP().Set("zero", "1"))
	fmt.Println(mdArray) */
	/* fmt.Printf("%v\n", ValueOf(&meArray).Elem().ARRAY().index(0).MAP().Set("one", "1"))
	fmt.Println(meArray)*/

}

var benchValues = []any{
	true,
	map[string]string{"0": "0"},
	&map[string]string{"0": "0"},
	&[2]string{"0", "0"},
}

func BenchmarkValueOf(b *testing.B) {
	for _, v := range benchValues {
		b.Run("Reflect - ValueOf "+getrtype(v).STRING().Width(24), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				reflect.ValueOf(v)
			}
		})
		b.Run("Old     - ValueOf "+getrtype(v).STRING().Width(24), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ValueOf(v)
			}
		})
	}
}

func BenchmarkMapKeysRef(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflect.ValueOf(map[string]string{"0": "0", "1": "1"}).MapKeys()
	}
}

func BenchmarkMapKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MapOf(map[string]string{"0": "0", "1": "1"}).Keys()
	}
}

/*
func BenchmarkValueRefBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflect.ValueOf(true)
	}
}

func BenchmarkValOfBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValOf(true)
	}
}

func BenchmarkVOfBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		VOF(true)
	}
}

func BenchmarkValueOfBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValueOf(true)
	}
}

func BenchmarkValueRefMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflect.ValueOf(map[string]string{"0": "0"})
	}
}

func BenchmarkValOfPtrMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValOf(map[string]string{"0": "0"})
	}
}

func BenchmarkVOfPtrMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		VOF(map[string]string{"0": "0"})
	}
}

func BenchmarkValueOfPtrMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValueOf(map[string]string{"0": "0"})
	}
}

func BenchmarkValueRefMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflect.ValueOf(map[string]string{"0": "0"})
	}
}

func BenchmarkValOfPtrMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValOf(map[string]string{"0": "0"})
	}
}

func BenchmarkVOfPtrMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		VOF(map[string]string{"0": "0"})
	}
}

func BenchmarkValueOfPtrMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValueOf(map[string]string{"0": "0"})
	}
}
*/
/*
func BenchmarkEncNative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decNative()
	}
}

func BenchmarkEncNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decNew()
	}
}

func TestEncNative(t *testing.T) {
	fmt.Printf("%#v\n", decNative())
}

func TestEncNew(t *testing.T) {
	fmt.Printf("%#v\n", decNew())
}

func decNative() any {
	b := []byte(`{"bytes":"1","bool":true,"int":1,"string":"1","slice":["1"],"map":{"string":"1"},"struct":{"Zero":"0","One":"1"}}`)
	m := map[string]any{}
	json.Unmarshal(b, &m)
	return m
}

func decNew() any {
	b := []byte(`{"bytes":"1","bool":true,"int":1,"string":"1","slice":["1"],"map":{"string":"1"},"struct":{"Zero":"0","One":"1"}}`)
	o, _ := STRING(b).Unserialize()
	return o
}
*/
