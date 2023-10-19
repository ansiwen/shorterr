// Package shorterr is an implementation of a short-circuit error handling
// inspired by the ? operator in Rust.
package shorterr

import (
	"errors"
	"fmt"
	"strings"
)

type shortCircuitError error

// PassTo stores the intercepted error in the variable err is pointing to. It
// must be installed with defer in the current function before the other
// short-circuit functions are used:
//
//	import se "github.com/ansiwen/shorterr"
//
//	func Foo() (err error) {
//		defer se.PassTo(&err)
//	...
func PassTo(err *error) {
	if v := recover(); v != nil {
		if e, ok := v.(shortCircuitError); ok {
			*err = e
		} else {
			panic(v)
		}
	}
}

// Check short-circuits the execution of the current function if the error is
// not nil. If the optional msg is provided, the err is wrapped with msg. PassTo
// must be installed with defer before.
func Check(err error, msg ...string) {
	if err != nil {
		msg := strings.Join(msg, " ")
		if len(msg) > 0 {
			err = fmt.Errorf("%s: %w", msg, err)
		}
		panic(shortCircuitError(err))
	}
}

// Assert short-circuits the execution of the current function if ok is false
// and returns msg as an error. PassTo must be installed with defer before.
func Assert(ok bool, msg string) {
	if !ok {
		panic(shortCircuitError(errors.New(msg)))
	}
}

// Try is a wrapper for functions that return a value and an error. It
// short-circuits the execution of the current function if the error is not nil.
// Otherwise it only returns the result value. PassTo must be installed with
// defer before.
func Try[A any](a A, err error) A {
	Check(err)
	return a
}

// Try2 is Try for functions with 2-ary results.
func Try2[A, B any](a A, b B, err error) (A, B) {
	Check(err)
	return a, b
}

// Try3 is Try for functions with 3-ary results.
func Try3[A, B, C any](a A, b B, c C, err error) (A, B, C) {
	Check(err)
	return a, b, c
}

// Try4 is Try for functions with 4-ary results.
func Try4[A, B, C, D any](a A, b B, c C, d D, err error) (A, B, C, D) {
	Check(err)
	return a, b, c, d
}

// Try5 is Try for functions with 5-ary results.
func Try5[A, B, C, D, E any](a A, b B, c C, d D, e E, err error) (A, B, C, D, E) {
	Check(err)
	return a, b, c, d, e
}

type Result[A any] struct {
	a   A
	err error
}

type Result2[A, B any] struct {
	a   A
	b   B
	err error
}

type Result3[A, B, C any] struct {
	a   A
	b   B
	c   C
	err error
}

type Result4[A, B, C, D any] struct {
	a   A
	b   B
	c   C
	d   D
	err error
}

type Result5[A, B, C, D, E any] struct {
	a   A
	b   B
	c   C
	d   D
	e   E
	err error
}

// Do is an alternative to Try that allows to wrap the short-circuit error with
// a description by appending the Or() method.
func Do[A any](a A, err error) *Result[A] {
	return &Result[A]{a, err}
}

// Do2 is Do for 2-ary results.
func Do2[A, B any](a A, b B, err error) *Result2[A, B] {
	return &Result2[A, B]{a, b, err}
}

// Do3 is Do for 3-ary results.
func Do3[A, B, C any](a A, b B, c C, err error) *Result3[A, B, C] {
	return &Result3[A, B, C]{a, b, c, err}
}

// Do4 is Do for 4-ary results.
func Do4[A, B, C, D any](a A, b B, c C, d D, err error) *Result4[A, B, C, D] {
	return &Result4[A, B, C, D]{a, b, c, d, err}
}

// Do5 is Do for 5-ary results.
func Do5[A, B, C, D, E any](a A, b B, c C, d D, e E, err error) *Result5[A, B, C, D, E] {
	return &Result5[A, B, C, D, E]{a, b, c, d, e, err}
}

// Or returns only the result value of the function called by Do if its returned
// error is nil. Otherwise it wraps the error with msg and short-circuits the
// execution of the current function. PassTo must be installed with
// defer before.
func (r *Result[A]) Or(msg string) A {
	Check(r.err, msg)
	return r.a
}

// Or for 2-ary results.
func (r *Result2[A, B]) Or(msg string) (A, B) {
	Check(r.err, msg)
	return r.a, r.b
}

// Or for 3-ary results.
func (r *Result3[A, B, C]) Or(msg string) (A, B, C) {
	Check(r.err, msg)
	return r.a, r.b, r.c
}

// Or for 4-ary results.
func (r *Result4[A, B, C, D]) Or(msg string) (A, B, C, D) {
	Check(r.err, msg)
	return r.a, r.b, r.c, r.d
}

// Or for 5-ary results.
func (r *Result5[A, B, C, D, E]) Or(msg string) (A, B, C, D, E) {
	Check(r.err, msg)
	return r.a, r.b, r.c, r.d, r.e
}
