package dbtest

import (
	"embed"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/sampleservice1/dbmigrations"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		dbData   string
		wantData Rows
	}{
		{
			name:     "creates isolated DB unaffected by parallel test cases - #1",
			dbData:   `create table foo (id int primary key); insert into foo values(1);`,
			wantData: Rows{{"id": int64(1)}},
		},
		{
			name:     "creates isolated DB unaffected by parallel test cases - #2",
			dbData:   `create table foo (id text primary key); insert into foo values('bar');`,
			wantData: Rows{{"id": "bar"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := New(t, dbmigrations.FS)
			_, err := db.Exec(tt.dbData)
			require.NoError(t, err)

			data, err := QueryAll(db, "foo", "id")
			require.NoError(t, err)
			assert.Equal(t, tt.wantData, data)
		})
	}
}

func TestMigrationDown(t *testing.T) {
	t.Parallel()

	sqls := getMigrationDownSQLs(t, dbmigrations.FS)
	require.NotEmpty(t, sqls)

	db := New(t, dbmigrations.FS)
	for _, sql := range sqls {
		_, err := db.Exec(sql)
		require.NoError(t, err)
	}
}

func getMigrationDownSQLs(t *testing.T, fs embed.FS) []string {
	dir, err := fs.ReadDir(".")
	require.NoError(t, err)

	result := []string{}
	for _, de := range dir {
		name := de.Name()
		if strings.HasSuffix(name, ".down.sql") {
			b, err := fs.ReadFile(name)
			require.NoError(t, err)
			result = append(result, string(b))
		}
	}
	slices.Reverse(result)

	return result
}
