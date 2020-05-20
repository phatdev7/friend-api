.PHONY: db

COMPOSE = docker-compose -f docker-compose.yml
sqlboiler: db
	sqlboiler --output internal/models/orm --pkgname orm psql

install-migration-tool:
	which migrate || go get -u -d github.com/golang-migrate/migrate/cli github.com/lib/pq go build -tags 'postgres' -o /usr/local/bin/migrate github.com/golang-migrate/migrate/cli

db:
	$(COMPOSE) up -d db && sleep 1

migrate: MOUNT_VOLUME = -v $(shell pwd)/db/migrations:/migrations
migrate:
	$(COMPOSE) run --rm $(MOUNT_VOLUME) db-migrate \
	sh -c './migrate -path /migrations -database $$DATABASE_URL up'