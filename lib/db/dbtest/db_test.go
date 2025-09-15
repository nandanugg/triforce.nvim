package dbtest

import (
	"context"
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

			db := New(t, testmigrations.FS)
			_, err := db.Exec(context.Background(), fmt.Sprintf("insert into foo values (1, %d)", i))
			require.NoError(t, err)

			data, err := QueryAll(db, "foo", "foo")
			require.NoError(t, err)
			assert.Equal(t, Rows{{"foo": int32(1), "bar": int32(i)}}, data)
		})
	}
}
