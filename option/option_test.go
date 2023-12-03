package option_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/JFAexe/typez/option"
)

func Test_Option(t *testing.T) {
	t.Parallel()

	var (
		num = 42
		alt = 1337
		def = 0
	)

	var (
		Some = option.Some[int](num)
		None = option.None[int]()
	)

	var (
		someAnd = func(v int) bool { return v == num }
		noneAnd = func() bool { return true }
		orElse  = func() int { return alt }
	)

	t.Run("IsSome", func(t *testing.T) {
		t.Parallel()

		require.True(t, Some.IsSome())
		require.False(t, None.IsSome())
	})

	t.Run("IsSomeAnd", func(t *testing.T) {
		t.Parallel()

		require.True(t, Some.IsSomeAnd(someAnd))
		require.False(t, None.IsSomeAnd(someAnd))
	})

	t.Run("IsNone", func(t *testing.T) {
		t.Parallel()

		require.False(t, Some.IsNone())
		require.True(t, None.IsNone())
	})

	t.Run("IsNoneAnd", func(t *testing.T) {
		t.Parallel()

		require.False(t, Some.IsNoneAnd(noneAnd))
		require.True(t, None.IsNoneAnd(noneAnd))
	})

	t.Run("Value", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, num, Some.Value())
		require.Panics(t, func() { None.Value() })
	})

	t.Run("ValueOr", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, num, Some.ValueOr(alt))
		require.Equal(t, alt, None.ValueOr(alt))
	})

	t.Run("ValueOrElse", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, num, Some.ValueOrElse(orElse))
		require.Equal(t, alt, None.ValueOrElse(orElse))
	})

	t.Run("ValueOrDefault", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, num, Some.ValueOrDefault())
		require.Equal(t, def, None.ValueOrDefault())
	})

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()

		var (
			wrapError error
			someBytes []byte
			noneBytes []byte
		)

		someBytes, wrapError = json.Marshal(Some)
		require.NoError(t, wrapError)
		require.JSONEq(t, fmt.Sprintf("{\"value\":%v}", Some.Value()), string(someBytes))

		noneBytes, wrapError = json.Marshal(None)
		require.NoError(t, wrapError)
		require.Equal(t, "{}", string(noneBytes))

		var (
			unwrappedSome option.Option[int]
			unwrappedNone option.Option[int]
		)

		wrapError = json.Unmarshal(someBytes, &unwrappedSome)
		require.NoError(t, wrapError)
		require.Equal(t, Some.Value(), unwrappedSome.Value())

		wrapError = json.Unmarshal(noneBytes, &unwrappedNone)
		require.NoError(t, wrapError)
		require.Panics(t, func() { unwrappedNone.Value() })
	})
}
