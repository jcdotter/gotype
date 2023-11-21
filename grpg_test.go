// Copyright 2023 james dotter. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import (
	"fmt"
	"testing"
)

type Model struct {
	Name          string
	Type          *Type
	Packages      Strings // list of packages to import
	Primary       *Fld
	Lock          bool `grpg:"lock"`
	Retain        bool `grpg:"retain"`
	History       bool `grpg:"track_history"`
	Fields        Fields
	CascadeDelete Fields  // populates on parent with parent field name and identifies children foreign keys to cascade delete
	CascadeUpdate Fields  // populates on parent with parent field name as key and identifies children fields to cascade to on update
	RollupFields  Rollups // populates on child with child field name as key and identifies parent fields to rollup to
	Messages      Messages
	Methods       Methods
}

type Models []*Model

type Fld struct {
	Model       *Model
	Name        string
	Type        *Type
	Primary     bool    `grpg:"primary"`
	Unique      bool    `grpg:"unique"`
	Indexed     bool    `grpg:"index"`
	Encrypted   bool    `grpg:"encrypt"`
	Required    Strings `grpg:"require"` // message names
	Default     *string `grpg:"default"`
	Lock        bool    `grpg:"lock"`
	ParentRef   *Fld    `goserv:"ref"`
	ChildRef    *Fld    `grpg:"ref"`
	Cascade     bool    `grpg:"cascade"`
	CascadeFrom *Fld    // populated on child field with reference to parent field where value is cascaded from
	RollupFrom  *Rollup // populated on parend field with reference to child field where value is rolled up from
}

type Fields []*Fld

type Message struct {
	Model    *Model
	Name     string
	Type     *Type
	Base     *TYPE
	Fields   Fields
	Required Fields
	SubMsgs  *map[string]*Message
}

type Messages []*Message

type Method struct {
	Model  *Model
	Name   string
	Input  *Message
	Output *Message
}

type Methods []*Method

type Rollup struct {
	Name        string
	Func        string
	ParentField *Fld
	ChildField  *Fld
}

type Rollups []*Rollup

type Strings []string

type Type struct {
	Name     string
	TYPE     *TYPE
	IsGrpg   bool
	Proto    string
	Postgres string
}

type User struct {
	ID      int
	Name    string
	Email   string
	Created int
}

var (
	models = Models{
		&Model{
			Name: "User",
			Type: &Type{
				Name: "User",
				TYPE: TypeOf(User{}),
			},
			Fields: Fields{
				&Fld{
					Name: "ID",
					Type: &Type{
						Name: "ID",
						TYPE: TypeOf(0),
					},
					Primary: true,
					Unique:  true,
					Indexed: true,
				},
				&Fld{
					Name: "Name",
					Type: &Type{
						Name: "Name",
						TYPE: TypeOf(""),
					},
					Unique:  true,
					Indexed: true,
				},
				&Fld{
					Name: "Email",
					Type: &Type{
						Name: "Email",
						TYPE: TypeOf(""),
					},
					Unique:  true,
					Indexed: true,
				},
				&Fld{
					Name: "Password",
					Type: &Type{
						Name: "Password",
						TYPE: TypeOf(""),
					},
					Encrypted: true,
				},
				&Fld{
					Name: "Created",
					Type: &Type{
						Name: "Created",
						TYPE: TypeOf(0),
					},
					Indexed: true,
				},
			},
			Messages: Messages{
				&Message{
					Name: "UserResponse",
					Type: &Type{
						Name: "User",
						TYPE: TypeOf(User{}),
					},
					Fields: Fields{
						&Fld{
							Name: "ID",
							Type: &Type{
								Name: "ID",
								TYPE: TypeOf(0),
							},
						},
						&Fld{
							Name: "Name",
							Type: &Type{
								Name: "Name",
								TYPE: TypeOf(""),
							},
						},
						&Fld{
							Name: "Email",
							Type: &Type{
								Name: "Email",
								TYPE: TypeOf(""),
							},
						},
						&Fld{
							Name: "Created",
							Type: &Type{
								Name: "Created",
								TYPE: TypeOf(0),
							},
						},
					},
				},
				&Message{
					Name: "UserRequest",
					Type: &Type{
						Name: "User",
						TYPE: TypeOf(User{}),
					},
					Fields: Fields{
						&Fld{
							Name: "Name",
							Type: &Type{
								Name: "Name",
								TYPE: TypeOf(""),
							},
						},
					},
					Required: Fields{
						&Fld{
							Name: "Name",
							Type: &Type{
								Name: "Name",
								TYPE: TypeOf(""),
							},
						},
					},
				},
			},
			Methods: Methods{
				&Method{
					Name: "Create",
					Input: &Message{
						Name: "UserRequest",
						Type: &Type{
							Name: "User",
							TYPE: TypeOf(User{}),
						},
						Fields: Fields{
							&Fld{
								Name: "Name",
								Type: &Type{
									Name: "Name",
									TYPE: TypeOf(""),
								},
							},
						},
						Required: Fields{
							&Fld{
								Name: "Name",
								Type: &Type{
									Name: "Name",
									TYPE: TypeOf(""),
								},
							},
						},
					},
					Output: &Message{
						Name: "UserResponse",
						Type: &Type{
							Name: "User",
							TYPE: TypeOf(User{}),
						},
						Fields: Fields{
							&Fld{
								Name: "ID",
								Type: &Type{
									Name: "ID",
									TYPE: TypeOf(0),
								},
							},
							&Fld{
								Name: "Name",

								Type: &Type{
									Name: "Name",
									TYPE: TypeOf(""),
								},
							},
							&Fld{
								Name: "Email",
								Type: &Type{
									Name: "Email",
									TYPE: TypeOf(""),
								},
							},
							&Fld{
								Name: "Created",
								Type: &Type{
									Name: "Created",
									TYPE: TypeOf(0),
								},
							},
						},
					},
				},
			},
		},
	}
)

func TestGrpg(t *testing.T) {
	/* m := YamlMarshaller.New()
	m.Marshal(models)
	fmt.Println(m.String()) */
	//j, _ := json.Marshal(models)
	j := JsonMarshaller.Marshal(models).Bytes()
	fmt.Println(string(j))
}
