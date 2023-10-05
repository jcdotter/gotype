// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package vals

import (
	"unsafe"

	"github.com/google/uuid"
)

// ------------------------------------------------------------ /
// GOTYPE CUSTOM TYPE IMPLEMENTATION
// implementation of custom type of uuid.UUID
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

type UUID uuid.UUID

// UUID returns gotype Value as UUID
func (v VALUE) UUID() UUID {
	switch v.KIND() {
	case String:
		return (*STRING)(v.ptr).UUID()
	case Bytes:
		return STRING(*(*[]byte)(v.ptr)).UUID()
	case Uuid:
		return *(*UUID)(v.ptr)
	case Pointer:
		return v.Elem().UUID()
	}
	panic("cannot convert value to UUID")
}

// NewUUID returns a randomly generated instance of UUID
func NewUUID() UUID {
	u, _ := uuid.NewRandom()
	return UUID(u)
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Natural returns gotype UUID as golang uuid.UUID
func (u UUID) Native() uuid.UUID {
	return uuid.UUID(u)
}

// Interface returns gotype UUID as a golang interface{}
func (u UUID) Interface() any {
	return u.Native()
}

// Value returns gotype UUID as gotype Value
func (u UUID) VALUE() VALUE {
	a := (any)(u)
	return *(*VALUE)(unsafe.Pointer(&a))
}

// Encode returns a gotype encoding of UUID
func (u UUID) Encode() ENCODING {
	return append([]byte{byte(Uuid)}, u.Bytes()...)
}

// Bytes returns gotype UUID as gotype Bytes
func (u UUID) Bytes() []byte {
	return u[:]
}

// Bool returns gotype UUID as bool
// false if empty, true if a UUID
func (u UUID) Bool() bool {
	return u != UUID{}
}

// String returns gotype UUID as string
func (u UUID) String() string {
	uu := UUID{}
	if u == uu {
		return ""
	}
	return u.Native().String()
}

// STRING returns gotype UUID as a gotype STRING
func (u UUID) STRING() STRING {
	return STRING(u.Native().String())
}

// Serialize returns gotype UUID as serialized string
func (u UUID) Serialize() string {
	return `"` + u.String() + `"`
}

// Uuid returns gotype UUID as uuid.UUID
func (u UUID) Uuid() uuid.UUID {
	return uuid.UUID(u)
}

// IsNil evaluates whether UUID is empty
func (u UUID) IsNil() bool {
	for _, b := range u {
		if b != 0 {
			return false
		}
	}
	return true
}
