export PGHOST=$NEXUS_TEST_DB_HOST
export PGUSER=$NEXUS_TEST_DB_SUPERUSER
export PGPASSWORD=$NEXUS_TEST_DB_PASSWORD
export PGDATABASE=$NEXUS_TEST_DB_NAME

psql -t -c "select datname from pg_database where datname like 'nexus_test_db_%'" \
| awk '{$1=$1};1' \
| grep -v '^$' \
| xargs -I {} -n1 psql -c "drop database {}"

psql -t -c "select rolname from pg_roles where rolname like 'nexus_test_user_%'" \
| awk '{$1=$1};1' \
| grep -v '^$' \
| xargs -I {} -n1 psql -c "drop role {}"
