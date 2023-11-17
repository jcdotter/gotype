// Copyright 2023 james dotter. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

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

func getTestVarsGmap() Gmap {
	m := MapOf(createTestVars(false, 0, "false")).Gmap()
	m.SortByKeys()
	return m
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

		"*[2]string":            &[2]string{s, s},
		"*[]string{2}":          &[]string{s, s},
		"*map[string]string{2}": &map[string]string{"0": s, "1": s},
		"*struct(string){2}":    &string_struct{s, s},

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
		"[2]*[1]string":            [2]*[1]string{{s}, {s}},
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
		"[]*[1]string{1}":           []*[1]string{{s}},
		"[]*[1]string{2}":           []*[1]string{{s}, {s}},
		"[]*[]string{1,1}":          []*[]string{{s}},
		"[]*map[string]string{1,1}": []*map[string]string{{"0": s}},
		"[]*struct(string){1,1}":    []*string_struct_single{{s}},
		"[]any(*string){1}":         []any{&s},

		"[]*bool{2}":                []*bool{&b, &b},
		"[]*int{2}":                 []*int{&i, &i},
		"[]*string{2}":              []*string{&s, &s},
		"[]*[2]string{2}":           []*[2]string{{s, s}, {s, s}},
		"[]*[]string{2,2}":          []*[]string{{s, s}, {s, s}},
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

		"map[string]*bool{2}":                map[string]*bool{"0": &b, "1": &b},
		"map[string]*int{2}":                 map[string]*int{"0": &i, "1": &i},
		"map[string]*string{2}":              map[string]*string{"0": &s, "1": &s},
		"map[string]*[2]string{2}":           map[string]*[2]string{"0": {s, s}, "1": {s, s}},
		"map[string]*[]string{2,2}":          map[string]*[]string{"0": {s, s}, "1": {s, s}},
		"map[string]*map[string]string{2,2}": map[string]*map[string]string{"0": {"0": s, "1": s}, "1": {"0": s, "1": s}},
		"map[string]*struct(string){2}":      map[string]*string_struct{"0": {s, s}, "1": {s, s}},
		"map[string]any(*string){2}":         map[string]any{"0": &s, "1": &s},

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

var YAML = `- Name: Team
Type: 'models.Team'
Packages:
  - github.com/jcdotter/grpg/example/models
  - github.com/jcdotter/grpg/types
Primary: Id
Retain: true
Fields:
  - Name: Id
	Type: 'types.UUID'
	Primary: true
	Unique: true
	Indexed: true
  - Name: Name
	Type: 'types.String'
	Unique: true
	Indexed: true
	Required:
	  - Create
  - Name: Users
	Type: '[]*models.User'
	ChildRef: User.Team
Messages:
  - Name: TeamCreate
	Model: Team
	Type: '*models.TeamCreate'
	Base: models.TeamCreate
	Fields:
	  - Team.Name
	Required:
	  - Team.Name
  - Name: TeamResponse
	Model: Team
	Type: '*models.TeamResponse'
	Base: models.TeamResponse
	Fields:
	  - Team.Id
	  - Team.Name
  - Name: TeamAccess
	Model: Team
	Type: '*models.TeamAccess'
	Base: models.TeamAccess
	Fields:
	  - Team.Id
	Required:
	  - Team.Id
  - Name: TeamUpdate
	Model: Team
	Type: '*models.TeamUpdate'
	Base: models.TeamUpdate
	Fields:
	  - Team.Id
	  - Team.Name
	Required:
	  - Team.Id
  - Name: TeamDetail
	Model: Team
	Type: '*models.TeamDetail'
	Base: models.TeamDetail
	Fields:
	  - Team.Id
	  - Team.Name
	  - Team.Users
	SubMsgs:
	  Users: TeamDetail_UserResponse
  - Name: TeamDetail_UserResponse
	Model: User
	Type: '[]*models.UserResponse'
	Base: models.UserResponse
	Fields:
	  - User.Id
	  - User.Number
	  - User.Name
	  - User.Email
- Name: User
Type: 'models.User'
Packages:
  - github.com/jcdotter/grpg/example/models
  - github.com/jcdotter/grpg/types
Primary: Id
Retain: true
Fields:
  - Name: Id
	Type: 'types.UUID'
	Primary: true
	Unique: true
	Indexed: true
  - Name: Number
	Type: 'types.Serial'
  - Name: Name
	Type: 'types.String'
	Unique: true
	Indexed: true
  - Name: Email
	Type: 'types.String'
	Unique: true
	Indexed: true
  - Name: Password
	Type: 'types.String'
	Encrypted: true
  - Name: Rate
	Type: 'types.Decimal'
  - Name: ActionCount
	Type: 'types.Int64'
	RollupFrom: Action.Count(Actions.Id)
  - Name: Actions
	Type: '[]*models.Action'
	ChildRef: Action.Actor
	Cascade: true
  - Name: Team
	Type: '*models.Team'
	ParentRef: Team.Id
CascadeDelete:
  - User.Actions
CascadeUpdate:
  - Action.Rate
Messages:
  - Name: UserCreate
	Model: User
	Type: '*models.UserCreate'
	Base: models.UserCreate
	Fields:
	  - User.Name
	  - User.Email
	  - User.Password
	Required:
	  - User.Name
	  - User.Email
	  - User.Password
  - Name: UserResponse
	Model: User
	Type: '*models.UserResponse'
	Base: models.UserResponse
	Fields:
	  - User.Id
	  - User.Number
	  - User.Name
	  - User.Email
  - Name: UserAccess
	Model: User
	Type: '*models.UserAccess'
	Base: models.UserAccess
	Fields:
	  - User.Id
	Required:
	  - User.Id
  - Name: UserUpdate
	Model: User
	Type: '*models.UserUpdate'
	Base: models.UserUpdate
	Fields:
	  - User.Id
	  - User.Name
	  - User.Email
	  - User.Password
	Required:
	  - User.Id
  - Name: UserDetail
	Model: User
	Type: '*models.UserDetail'
	Base: models.UserDetail
	Fields:
	  - User.Id
	  - User.Number
	  - User.Name
	  - User.Email
	  - User.Rate
	  - User.ActionCount
	  - User.Actions
	  - User.Team
	SubMsgs:
	  Actions: UserDetail_ActionResponse
	  Team: UserDetail_TeamResponse
  - Name: UserDetail_ActionResponse
	Model: Action
	Type: '[]*models.ActionResponse'
	Base: models.ActionResponse
	Fields:
	  - Action.Id
	  - Action.Type
	  - Action.Actor
	  - Action.Rate
	  - Action.Time
	SubMsgs:
	  Actor: UserDetail_ActionResponse_UserResponse
  - Name: UserDetail_ActionResponse_UserResponse
	Model: User
	Type: '*models.UserResponse'
	Base: models.UserResponse
	Fields:
	  - User.Id
  - Name: UserDetail_TeamResponse
	Model: Team
	Type: '*models.TeamResponse'
	Base: models.TeamResponse
	Fields:
	  - Team.Id
	  - Team.Name
- Name: Action
Type: 'models.Action'
Packages:
  - github.com/jcdotter/grpg/example/models
  - github.com/jcdotter/grpg/types
Primary: Id
Lock: true
Retain: true
Fields:
  - Name: Id
	Type: 'types.UUID'
	Primary: true
	Unique: true
	Indexed: true
  - Name: Type
	Type: 'types.String'
  - Name: Actor
	Type: '*models.User'
	ParentRef: User.Id
  - Name: Rate
	Type: 'types.Decimal'
	CascadeFrom: User.Rate
  - Name: MainAct
	Type: '*models.Action'
	ParentRef: Action.Id
  - Name: SubActs
	Type: '[]*models.Action'
	ChildRef: Action.MainAct
  - Name: Time
	Type: 'types.Time'
RollupFields:
  - Name: Count(Actions.Id)
	Func: COUNT
	ParentField: User.ActionCount
	ChildField: Action.Id
Messages:
  - Name: ActionCreate
	Model: Action
	Type: '*models.ActionCreate'
	Base: models.ActionCreate
	Fields:
	  - Action.Type
	  - Action.Actor
	  - Action.MainAct
	Required:
	  - Action.Type
	  - Action.Actor
	  - Action.MainAct
	SubMsgs:
	  Actor: ActionCreate_UserAccess
	  MainAct: ActionCreate_ActionAccess
  - Name: ActionCreate_ActionAccess
	Model: Action
	Type: '*models.ActionAccess'
	Base: models.ActionAccess
	Fields:
	  - Action.Id
	Required:
	  - Action.Id
  - Name: ActionCreate_UserAccess
	Model: User
	Type: '*models.UserAccess'
	Base: models.UserAccess
	Fields:
	  - User.Id
	Required:
	  - User.Id
  - Name: ActionResponse
	Model: Action
	Type: '*models.ActionResponse'
	Base: models.ActionResponse
	Fields:
	  - Action.Id
	  - Action.Type
	  - Action.Actor
	  - Action.Rate
	  - Action.Time
	SubMsgs:
	  Actor: ActionResponse_UserResponse
  - Name: ActionResponse_UserResponse
	Model: User
	Type: '*models.UserResponse'
	Base: models.UserResponse
	Fields:
	  - User.Id
	  - User.Number
	  - User.Name
	  - User.Email
  - Name: ActionAccess
	Model: Action
	Type: '*models.ActionAccess'
	Base: models.ActionAccess
	Fields:
	  - Action.Id
	Required:
	  - Action.Id
  - Name: ActionUpdate
	Model: Action
	Type: '*models.ActionUpdate'
	Base: models.ActionUpdate
	Fields:
	  - Action.Id
	  - Action.Type
	  - Action.Actor
	Required:
	  - Action.Id
	SubMsgs:
	  Actor: ActionUpdate_UserAccess
  - Name: ActionUpdate_UserAccess
	Model: User
	Type: '*models.UserAccess'
	Base: models.UserAccess
	Fields:
	  - User.Id
	Required:
	  - User.Id
  - Name: ActionDetail
	Model: Action
	Type: '*models.ActionDetail'
	Base: models.ActionDetail
	Fields:
	  - Action.Id
	  - Action.Type
	  - Action.Actor
	  - Action.Rate
	  - Action.MainAct
	  - Action.SubActs
	  - Action.Time
	SubMsgs:
	  MainAct: ActionDetail_ActionResponse
	  SubActs: ActionDetail_ActionResponse
	  Actor: ActionDetail_UserResponse
  - Name: ActionDetail_ActionResponse
	Model: Action
	Type: '*models.ActionResponse'
	Base: models.ActionResponse
	Fields:
	  - Action.Id
	  - Action.Type
	  - Action.Actor
	  - Action.Rate
	  - Action.Time
	SubMsgs:
	  Actor: ActionDetail_ActionResponse_UserResponse
  - Name: ActionDetail_ActionResponse_UserResponse
	Model: User
	Type: '*models.UserResponse'
	Base: models.UserResponse
	Fields:
	  - User.Id
  - Name: ActionDetail_UserResponse
	Model: User
	Type: '*models.UserResponse'
	Base: models.UserResponse
	Fields:
	  - User.Id
	  - User.Number
	  - User.Name
	  - User.Email
`
