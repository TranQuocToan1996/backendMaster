containername=postgres12
rootuser=root
pw=mysecretpassword
dbname=simple_bank
dbpath=db/migration

dockerpull:
	docker pull postgres:12-alpine

runpostgres:
	docker run --name $(containername) -e POSTGRES_USER=$(rootuser) -e POSTGRES_PASSWORD=$(pw) -p 5432:5432 -d postgres:12-alpine

stoppostgre:
	docker stop $(containername)

startpostgre:
	docker start $(containername)

psqlexec:
	docker exec -it $(containername) psql -U $(rootuser)

shellexec:
	docker exec -it $(containername) /bin/sh

createdb:
	docker exec -it $(containername) createdb --username=$(rootuser) --owner=$(rootuser) $(dbname) 

dropdb:
	docker exec -it $(containername) dropdb $(dbname) 

logs:
	docker logs $(containername)

createMigrate:
	migrate create -ext sql -dir $(dbpath) -seq init_schema

migrateup:
	@echo "update to newest"
	migrate -path $(dbpath) -database "postgresql://$(rootuser):$(pw)@localhost:5432/$(dbname)?sslmode=disable" -verbose up

migrateup1:
	@echo "update one more seq"
	migrate -path $(dbpath) -database "postgresql://$(rootuser):$(pw)@localhost:5432/$(dbname)?sslmode=disable" -verbose up 1

migratedown:
	@echo "rollback all"
	migrate -path $(dbpath) -database "postgresql://$(rootuser):$(pw)@localhost:5432/$(dbname)?sslmode=disable" -verbose down

migratedown1:
	@echo "rollback last"
	migrate -path $(dbpath) -database "postgresql://$(rootuser):$(pw)@localhost:5432/$(dbname)?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

tests:
	go test -v -cover -race -timeout 1s ./...

server:
	go run main.go

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/TranQuocToan1996/backendMaster/db/sqlc Store 

.PHONY: dockerpull runpostgres stoppostgre startpostgre psqlexec shellexec createdb dropdb logs createMigrate migrateup sqlc tests server mockgen migrateup1 migratedown1