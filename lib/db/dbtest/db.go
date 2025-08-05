package dbtest

import (
	"crypto/sha1"
	"database/sql"
	"embed"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

var adminDB *sql.DB

func init() {
	var err error
	adminDB, err = db.New(
		os.Getenv("NEXUS_TEST_DB_HOST"),
		os.Getenv("NEXUS_TEST_DB_SUPERUSER"),
		os.Getenv("NEXUS_TEST_DB_PASSWORD"),
		os.Getenv("NEXUS_TEST_DB_NAME"),
		"public",
	)
	if err != nil {
		fmt.Println("Error connecting to database: " + err.Error())
		os.Exit(1)
	}
}

// New creates an isolated database for t, using application config values only
// as initial connection parameters. DB user specified in config must be a
// superuser (or at least have the capability to create roles and schemas).
// DB schema specified in config will be left untouched.
func New(t *testing.T, migrationsFS embed.FS) *sql.DB {
	t.Helper()

	testID := randomString()
	var (
		host     = os.Getenv("NEXUS_TEST_DB_HOST")
		dbname   = os.Getenv("NEXUS_TEST_DB_NAME")
		schema   = "nexus_test_schema_" + testID
		user     = "nexus_test_user_" + testID
		password = "any"
	)

	createTestRoleAndSchema(t, user, password, schema)
	return createTestTables(t, host, user, password, dbname, schema, migrationsFS)
}

func randomString() string {
	sum := sha1.Sum([]byte(uuid.New().String()))
	return hex.EncodeToString(sum[:4])
}

func createTestRoleAndSchema(t *testing.T, user, password, schema string) {
	_, err := adminDB.Exec(fmt.Sprintf(`
		create role %[1]s login password '%[2]s';
		create schema %[3]s authorization %[1]s;
		`, user, password, schema,
	))
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err = adminDB.Exec(fmt.Sprintf(`
			drop schema %s cascade;
			drop role %s;
			`, schema, user,
		))
		assert.NoError(t, err)
	})
}

func createTestTables(t *testing.T, host, user, password, dbname, schema string, migrationsFS embed.FS) *sql.DB {
	db, err := db.New(host, user, password, dbname, schema)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	sqls, err := getMigrationUpSQLs(migrationsFS)
	require.NoError(t, err)

	_, err = db.Exec(strings.Join(sqls, "\n"))
	require.NoError(t, err, "Error running SQL:\n%s")

	return db
}

func getMigrationUpSQLs(fs embed.FS) ([]string, error) {
	dir, err := fs.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("read migrations dir: %w", err)
	}

	result := []string{}
	for _, de := range dir {
		name := de.Name()
		if strings.HasSuffix(name, ".up.sql") {
			b, err := fs.ReadFile(name)
			if err != nil {
				return nil, fmt.Errorf("read file %s: %w", name, err)
			}
			result = append(result, string(b))
		}
	}

	return result, nil
}

type Rows []map[string]any

// QueryAll fetches all rows in table tableName from db.
func QueryAll(db *sql.DB, tableName, orderBy string) (Rows, error) {
	rows, err := db.Query("select * from " + tableName + " order by " + orderBy)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("get rows cols: %w", err)
	}

	result := []map[string]any{}
	for rows.Next() {
		row := make([]any, len(cols))
		rowPtr := make([]any, len(cols))
		for i := range row {
			rowPtr[i] = &row[i]
		}

		if err := rows.Scan(rowPtr...); err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		rowMap := map[string]any{}
		for i, v := range row {
			rowMap[cols[i]] = v
		}
		result = append(result, rowMap)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return result, nil
}
