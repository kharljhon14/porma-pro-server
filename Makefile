postgres: 
	docker run --name postgres12 -p 5432:5432  -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root porma_pro

dropdb:
	docker exec -it postgres12 dropdb porma_pro

migrateup:
	migrate -path internal/db/migration -database postgresql://root:postgres@localhost:5432/porma_pro?sslmode=disable -verbose up

migratedown:
	migrate -path internal/db/migration -database postgresql://root:postgres@localhost:5432/porma_pro?sslmode=disable -verbose down

sqlc:
	sqlc generate

mockgen:
	mockgen -destination internal/db/mock/store.go github.com/kharljhon14/porma-pro-server/internal/db/sqlc Store

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mockgen