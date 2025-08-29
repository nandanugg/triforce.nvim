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

var (
	adminDB         *sql.DB
	testDBHost      string
	testDBName      string
	testDBSchema    string
	testDBSuperuser string
	testDBPassword  string
)

func init() {
	testDBHost = os.Getenv("NEXUS_TEST_DB_HOST")
	testDBName = os.Getenv("NEXUS_TEST_DB_NAME")
	testDBSchema = os.Getenv("NEXUS_TEST_DB_SCHEMA")
	testDBSuperuser = os.Getenv("NEXUS_TEST_DB_SUPERUSER")
	testDBPassword = os.Getenv("NEXUS_TEST_DB_PASSWORD")

	var err error
	adminDB, err = db.New(
		testDBHost,
		testDBSuperuser,
		testDBPassword,
		testDBName,
		testDBSchema)
	if err != nil {
		fmt.Println("Error connecting to database: " + err.Error())
		os.Exit(1)
	}
}

// New creates an isolated database for t, using application config values only
// as initial connection parameters.
func New(t *testing.T, migrationsFS embed.FS) *sql.DB {
	t.Helper()

	testID := randomString()
	var (
		dbname   = "nexus_test_db_" + testID
		user     = "nexus_test_user_" + testID
		schema   = "nexus_test_schema_" + testID
		password = "any"
	)

	// Using superuser, create new DB and app user:
	createTestDBAndRole(t, user, password, dbname)

	// Do necessary preparation in new DB before applying migration files:
	prepareTestDB(t, user, password, dbname, schema)

	// Apply migration files:
	return applyMigrations(t, user, password, dbname, schema, migrationsFS)
}

func randomString() string {
	sum := sha1.Sum([]byte(uuid.New().String()))
	return hex.EncodeToString(sum[:4])
}

func createTestDBAndRole(t *testing.T, user, password, dbname string) {
	_, err := adminDB.Exec("create database " + dbname)
	require.NoError(t, err)

	_, err = adminDB.Exec(fmt.Sprintf(`
		create role %[1]s login password '%[2]s';
		grant create on database %[3]s to %[1]s;
		`, user, password, dbname,
	))
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err := adminDB.Exec("drop database " + dbname)
		assert.NoError(t, err)
		_, err = adminDB.Exec("drop role " + user)
		assert.NoError(t, err)
	})
}

func prepareTestDB(t *testing.T, user, password, dbname, schema string) {
	d, err := db.New(testDBHost, testDBSuperuser, testDBPassword, dbname, schema)
	require.NoError(t, err)
	defer d.Close()

	// NOTE: 'update pg_language' dibutuhkan di DB kepegawaian (search
	// "LANGUAGE c" di file migrations).
	_, err = d.Exec("update pg_language set lanpltrusted = true where lanname = 'c'")
	require.NoError(t, err)

	d2, err := db.New(testDBHost, user, password, dbname, schema)
	require.NoError(t, err)
	defer d2.Close()

	_, err = d2.Exec("create schema " + schema)
	require.NoError(t, err)
}

func applyMigrations(t *testing.T, user, password, dbname, schema string, migrationsFS embed.FS) *sql.DB {
	db, err := db.New(testDBHost, user, password, dbname, schema)
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
