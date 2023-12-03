// Package option provides an option type.
package option

import (
	"encoding/json"
	"fmt"
)

// Option is a rust like option type.
type Option[T any] struct {
	val *T
}

// Some returns an Option with a value.
func Some[T any](value T) Option[T] {
	return Option[T]{
		val: &value,
	}
}

// None returns an Option without a value.
func None[T any]() Option[T] {
	return Option[T]{
		val: nil,
	}
}

// IsSome checks if Option isn't empty.
func (o *Option[T]) IsSome() bool {
	return o.val != nil
}

// IsSomeAnd checks if Option isn't empty with an extra check.
func (o *Option[T]) IsSomeAnd(and func(v T) bool) bool {
	if o.IsSome() {
		return and(*o.val)
	}

	return false
}

// IsNone checks if Option is empty.
func (o *Option[T]) IsNone() bool {
	return o.val == nil
}

// IsNoneAnd checks if Option is empty with an extra check.
func (o *Option[T]) IsNoneAnd(and func() bool) bool {
	if o.IsNone() {
		return and()
	}

	return false
}

// Value returns value from an Option.
//
// Panics if Option has no value.
func (o *Option[T]) Value() T {
	if o.IsNone() {
		panic("can't unwrap none value")
	}

	return *o.val
}

// ValueOr returns value from an Option or passed value.
func (o *Option[T]) ValueOr(or T) T {
	if o.IsNone() {
		return or
	}

	return *o.val
}

// ValueOrElse returns value from an Option or from passed func.
func (o *Option[T]) ValueOrElse(or func() T) T {
	if o.IsNone() {
		return or()
	}

	return *o.val
}

// ValueOrDefault returns value from an Option or type's default value.
func (o *Option[T]) ValueOrDefault() T {
	if o.IsNone() {
		return *new(T)
	}

	return *o.val
}

// String prints value of an Option.
func (o Option[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("Some: %v", *o.val)
	}

	return "None"
}

func (o Option[T]) wrap() optionWrap[T] {
	return optionWrap[T]{
		Value: o.val,
	}
}

func (o *Option[T]) unwrap(opt *optionWrap[T], err error) error {
	if err != nil {
		return err
	}

	o.val = opt.Value

	return nil
}

// optionWrap is a wrapper for Option encoding.
type optionWrap[T any] struct {
	Value *T `json:"value,omitempty"`
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.wrap())
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	res := new(optionWrap[T])

	return o.unwrap(res, json.Unmarshal(data, res))
}
