// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package gotype

import (
	"fmt"
	"reflect"
	"unsafe"
)

// In evaluates if 'x' is equal to any 'y'
func In(x any, ys ...any) bool {
	for _, y := range ys {
		if x == y {
			return true
		}
	}
	return false
}

// PrintStats prints stats about a to the console, including
// the line where PrintStats was called, Type, gotype DataType, and value
func PrintStats(a any) {
	fmt.Printf("Item Stats:\nLoc-1:    %v\nLoc:      %v\nType:     %T\nKind:     %v\nVal:      %v\nGoVal:    %#v\n\n", Source(2), Source(1), a, ValueOf(a).Kind(), a, a)
}

// Format inserts a into string format, similar to fmt.Sprintf
// TODO: optimize performance with alternative to fmt.Sprintf
func Format(format string, a ...any) string {
	return fmt.Sprintf(string(format), a...)
}

// DeepEqual reports whether x and y are “deeply equal”,
// as defined by golang reflect package
func DeepEqual(x any, y any) bool {
	return reflect.DeepEqual(x, y)
}

// Abs returns the absolute value of x
func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// AbsInt returns the absolute value of x
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ------------------------------------------------------------ /
// GOLANG NOESCAPES
// imported functions from standard golang library
// ------------------------------------------------------------ /

//go:noescape
//go:linkname makemap runtime.makemap
func makemap(t *rtype, cap int, hint unsafe.Pointer) unsafe.Pointer

//go:noescape
//go:linkname refmakemap reflect.makemap
func refmakemap(t *rtype, cap int) unsafe.Pointer

//go:noescape
//go:linkname maplen reflect.maplen
func maplen(unsafe.Pointer) int

//go:noescape
//go:linkname mapiterinit reflect.mapiterinit
func mapiterinit(t *rtype, m unsafe.Pointer, it *hiter)

//go:noescape
//go:linkname mapiterkey reflect.mapiterkey
func mapiterkey(it *hiter) (key unsafe.Pointer)

//go:noescape
//go:linkname mapiterelem reflect.mapiterelem
func mapiterelem(it *hiter) (elem unsafe.Pointer)

//go:noescape
//go:linkname mapiternext reflect.mapiternext
func mapiternext(it *hiter)

//go:noescape
//go:linkname mapaccess_faststr reflect.mapaccess_faststr
func mapaccess_faststr(t *rtype, m unsafe.Pointer, key string) (val unsafe.Pointer)

//go:noescape
//go:linkname mapdelete_faststr reflect.mapdelete_faststr
func mapdelete_faststr(t *rtype, m unsafe.Pointer, key string)

//go:noescape
//go:linkname mapassign_faststr reflect.mapassign_faststr
func mapassign_faststr(t *rtype, m unsafe.Pointer, key string, val unsafe.Pointer)

//go:noescape
//go:linkname toType reflect.toType
func toType(t *rtype) reflect.Type

//go:noescape
//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(*rtype) unsafe.Pointer

//go:noescape
//go:linkname unsafe_NewArray reflect.unsafe_NewArray
func unsafe_NewArray(*rtype, int) unsafe.Pointer

//go:noescape
//go:linkname growslice reflect.growslice
func growslice(t *rtype, old sliceHeader, num int) sliceHeader

//go:noescape
//go:linkname typedmemmove reflect.typedmemmove
func typedmemmove(t *rtype, dst, src unsafe.Pointer)

//go:noescape
//go:linkname resolveNameOff reflect.resolveNameOff
func resolveNameOff(ptrInModule unsafe.Pointer, off int32) unsafe.Pointer

//go:noescape
//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ *rtype, needzero bool) unsafe.Pointer

//go:noescape
//go:linkname fastrand runtime.fastrand
func fastrand() uint32

//go:noescape
//go:linkname fastrand64 runtime.fastrand64
func fastrand64() uint64
