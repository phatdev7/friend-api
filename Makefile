.PHONY: db

COMPOSE = docker-compose -f docker-compose.yml
sqlboiler:
	sqlboiler --output models/orm --pkgname orm psql

boil:
	sqlboiler --output internal/models2 --pkgname models2 postgres

install-migration-tool:
	which migrate || go get -u -d github.com/golang-migrate/migrate/cli github.com/lib/pq go build -tags 'postgres' -o /usr/local/bin/migrate github.com/golang-migrate/migrate/cli
APP_NAME := friend-api
RUN_COMPOSE = $(COMPOSE) run --rm --service-ports -w /go/src/$(APP_NAME) $(MOUNT_VOLUME) go

run: db sqlboiler start
db:
	$(COMPOSE) up -d db && sleep 5
	make migrate

migrate: MOUNT_VOLUME = -v $(shell pwd)/db/migrations:/migrations
migrate:
	$(COMPOSE) run --rm $(MOUNT_VOLUME) db-migrate \
	sh -c './migrate -path /migrations -database $$DATABASE_URL up'

start:
	go run -mod=readonly main.go

# run:
# 	@$(RUN_COMPOSE) env $(shell cat .env | egrep -v '^#|DATABASE_URL' | xargs) \
# 	go run -mod=vendor main.go