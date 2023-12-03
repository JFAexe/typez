package result_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/JFAexe/typez/result"
)

func Test_Result(t *testing.T) {
	t.Parallel()

	var (
		num = 42
		alt = 1337
		def = 0
		err = errors.New("wrong number")
	)

	var (
		Ok  = result.Ok[int](num)
		Err = result.Err[int](err)
	)

	var (
		okAnd  = func(v int) bool { return v == num }
		errAnd = func(e error) bool { return errors.Is(e, err) }
		orElse = func() int { return alt }
	)

	t.Run("Err", func(t *testing.T) {
		t.Parallel()

		require.Panics(t, func() { result.Err[int](nil) })
	})

	t.Run("AsResult", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, result.Err[int](err), result.AsResult[int](num, err))
	})

	t.Run("IsOk", func(t *testing.T) {
		t.Parallel()

		require.True(t, Ok.IsOk())
		require.False(t, Err.IsOk())
	})

	t.Run("IsOkAnd", func(t *testing.T) {
		t.Parallel()

		require.True(t, Ok.IsOkAnd(okAnd))
		require.False(t, Err.IsOkAnd(okAnd))
	})

	t.Run("IsErr", func(t *testing.T) {
		t.Parallel()

		require.False(t, Ok.IsErr())
		require.True(t, Err.IsErr())
	})

	t.Run("IsErrAnd", func(t *testing.T) {
		t.Parallel()

		require.False(t, Ok.IsErrAnd(errAnd))
		require.True(t, Err.IsErrAnd(errAnd))
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		require.Panics(t, func() { Ok.Error() })
		require.ErrorIs(t, Err.Error(), err)
	})

	t.Run("Value", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, num, Ok.Value())
		require.Panics(t, func() { Err.Value() })
	})

	t.Run("ValueOr", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, num, Ok.ValueOr(alt))
		require.Equal(t, alt, Err.ValueOr(alt))
	})

	t.Run("ValueOrElse", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, num, Ok.ValueOrElse(orElse))
		require.Equal(t, alt, Err.ValueOrElse(orElse))
	})

	t.Run("ValueOrDefault", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, num, Ok.ValueOrDefault())
		require.Equal(t, def, Err.ValueOrDefault())
	})

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()

		var (
			wrapError error
			okBytes   []byte
			errBytes  []byte
		)

		okBytes, wrapError = json.Marshal(Ok)
		require.NoError(t, wrapError)
		require.JSONEq(t, fmt.Sprintf("{\"value\":%v}", Ok.Value()), string(okBytes))

		errBytes, wrapError = json.Marshal(Err)
		require.NoError(t, wrapError)
		require.Equal(t, "{\"error\":\"wrong number\"}", string(errBytes))

		var (
			unwrappedOk  result.Result[int]
			unwrappedErr result.Result[int]
		)

		wrapError = json.Unmarshal(okBytes, &unwrappedOk)
		require.NoError(t, wrapError)
		require.Equal(t, Ok.Value(), unwrappedOk.Value())

		wrapError = json.Unmarshal(errBytes, &unwrappedErr)
		require.NoError(t, wrapError)
		require.Equal(t, Err.Error(), unwrappedErr.Error())
	})
}
