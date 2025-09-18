package typeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCast(t *testing.T) {
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "abcd", Cast[string]("abcd"))
		assert.Empty(t, Cast[string](1))
	})

	t.Run("time.Time", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), Cast[time.Time](time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)))
		assert.Empty(t, Cast[time.Time](""))
	})

	t.Run("int", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, 12, Cast[int](12))
		assert.Empty(t, Cast[int](uint(12)))
	})
}
