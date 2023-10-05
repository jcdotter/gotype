// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: James Dotter

package gotype

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var (
	module = `types`
	delim  = `.`
)

// GenError returns a formatted error message with source of error
// 'callOffset' indicates proceeding calls where the error orinigated
// 'callOffset' should be 0 if error occured in call to current function
func GenError(callOffset int, format string, a ...any) error {
	s := Source(callOffset + 1)
	return errors.New(string(s + ": " + Format(format, a...)))
}

// ParamError returns error message for expected and received type differences
// "[error location]: param type error expected type: 'e' received type: 'r'"
func ParamError(e string, r any) error {
	return GenError(1, "param type error\nexpected type: %s\nreceived type: %T", e, r)
}

func Source(i int) string {
	pc, fl, ln, ok := runtime.Caller(int(i + 1))
	if ok {
		fs := strings.Split(fl, `/`)
		gf := strings.Split(fs[len(fs)-1], delim)[0]
		fn := strings.Split(runtime.FuncForPC(pc).Name(), delim)
		pt := strings.Replace(fn[0], `/`, delim, -1)
		s := []string{pt, gf}
		s = append(s, fn[1:]...)
		return strings.Join(s, delim) + ` line ` + fmt.Sprint(ln)
	}
	return module + delim + `unknown.source`
}
