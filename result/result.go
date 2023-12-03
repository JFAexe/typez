// Package result provides a result type.
package result

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrNoValueError = errors.New("no value or error parsed")

// Result is a rust like result type.
type Result[T any] struct {
	val *T
	err error
}

// Ok returns a Result with a value.
func Ok[T any](value T) Result[T] {
	return Result[T]{
		val: &value,
		err: nil,
	}
}

// Err returns a Result with an error.
//
// Panics if error is nil.
func Err[T any](err error) Result[T] {
	if err == nil {
		panic("error can't be nil")
	}

	return Result[T]{
		val: nil,
		err: err,
	}
}

// AsResult returns a Result of passed type.
//
// If err isn't nil returns Result with an error.
func AsResult[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}

	return Ok(value)
}

// IsOk checks if Result isn't an error.
func (r *Result[T]) IsOk() bool {
	return r.val != nil && r.err == nil
}

// IsOkAnd checks if Result isn't an error with an extra check.
func (r *Result[T]) IsOkAnd(and func(v T) bool) bool {
	if r.IsOk() {
		return and(*r.val)
	}

	return false
}

// IsErr if Result is an error.
func (r *Result[T]) IsErr() bool {
	return r.val == nil && r.err != nil
}

// IsErrAnd if Result is an error with an extra check.
func (r *Result[T]) IsErrAnd(and func(e error) bool) bool {
	if r.IsErr() {
		return and(r.err)
	}

	return false
}

// Error returns error from a Result.
//
// Panics if Result has a value.
func (r *Result[T]) Error() error {
	if r.IsOk() {
		panic("can't unwrap error in the result with a value")
	}

	return r.err
}

// Value returns value from a Result.
//
// Panics if Result has an error.
func (r *Result[T]) Value() T {
	if r.IsErr() {
		panic("can't unwrap value in the result with an error")
	}

	return *r.val
}

// ValueOr returns value from a Result or passed value.
func (r *Result[T]) ValueOr(or T) T {
	if r.IsErr() {
		return or
	}

	return *r.val
}

// ValueOrElse returns value from a Result or from passed func.
func (r *Result[T]) ValueOrElse(or func() T) T {
	if r.IsErr() {
		return or()
	}

	return *r.val
}

// ValueOrDefault returns value from a Result or type's default value.
func (r *Result[T]) ValueOrDefault() T {
	if r.IsErr() {
		return *new(T)
	}

	return *r.val
}

// String prints value of an Option.
func (r Result[T]) String() string {
	if r.IsOk() {
		return fmt.Sprintf("Value: %v", *r.val)
	}

	return fmt.Sprintf("Error: %v", r.err)
}

func (r Result[T]) wrap() resultWrap[T] {
	var err string

	if r.err != nil {
		err = r.err.Error()
	}

	return resultWrap[T]{
		Value: r.val,
		Error: err,
	}
}

func (r *Result[T]) unwrap(res *resultWrap[T], err error) error {
	if err != nil {
		return err
	}

	if res.Error != "" {
		r.err = errors.New(res.Error)

		return nil
	}

	if res.Value != nil {
		r.val = res.Value

		return nil
	}

	return ErrNoValueError
}

// resultWrap is a wrapper for Result encoding.
type resultWrap[T any] struct {
	Value *T     `json:"value,omitempty"`
	Error string `json:"error,omitempty"`
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.wrap())
}

func (r *Result[T]) UnmarshalJSON(data []byte) error {
	res := new(resultWrap[T])

	return r.unwrap(res, json.Unmarshal(data, res))
}
