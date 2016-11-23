# This script requires PGPASSWORD environment variable to be set to run
{
  "${PGPASSWORD?Need to set PGPASSWORD env var}"
} &> /dev/null

echo "DROP DATABASE meirldb; \q\n" | psql -U postgres

psql -U postgres -f ./sql/create_database.sql
psql -U postgres -d meirldb -f ./sql/create_schema.sql
psql -U postgres -d meirldb -f ./sql/seed.sql