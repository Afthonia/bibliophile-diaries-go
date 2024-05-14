run:
	go run httpd/main.go httpd/routes.go

postgres:
	sudo -u postgres psql
	
db:
	PGPASSWORD=123456 psql bibliophile -U ayse -h localhost

init_schema:
	migrate create -ext sql -dir db/migration -seq init_schema

migrate_up:
	migrate -path db/migration -database "postgresql://ayse:123456@localhost:5432/bibliophile?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://ayse:123456@localhost:5432/bibliophile?sslmode=disable" -verbose down

migrate_force:
	migrate -path db/migration -database "postgresql://ayse:123456@localhost:5432/bibliophile?sslmode=disable" force 1
	migrate -path db/migration -database "postgresql://ayse:123456@localhost:5432/bibliophile?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres db init_schema migrate_up migrate_down migrate_force sqlc