package dbtest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest/testmigrations"
)

func TestNew(t *testing.T) {
	t.Parallel()

	for i := range 10 {
		t.Run(fmt.Sprintf("safe to execute in parallel - %d", i), func(t *testing.T) {
			t.Parallel()

			db := New(t, "test", testmigrations.FS)
			_, err := db.Exec(fmt.Sprintf("insert into test.foo values (1, %d)", i))
			require.NoError(t, err)

			data, err := QueryAll(db, "test.foo", "foo")
			require.NoError(t, err)
			assert.Equal(t, Rows{{"foo": int64(1), "bar": int64(i)}}, data)
		})
	}
}
