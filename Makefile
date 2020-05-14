sqlboiler: db
	sqlboiler --output internal/models/orm --pkgname orm psql

install-migration-tool:
	which migrate || go get -u -d github.com/golang-migrate/migrate/cli github.com/lib/pq go build -tags 'postgres' -o /usr/local/bin/migrate github.com/golang-migrate/migrate/cli