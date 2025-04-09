recall_BINARY=recallApp
# ==================================================================================== # 
# HELPERS 
# ==================================================================================== #

## help: print this help message
.PHONY: help confirm run/api db/psql db/migrate/up db/migrate/down db/migrate/upt db/migrate/downt audit vendor build/api build/docker run/docker swag
help: 
	@echo 'Usage:' 
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


confirm: 
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== # 
# DEVELOPMENT 
# ==================================================================================== #

## run/api: run the cmd/api application 
swag:
	@echo 'Generating Swagger Docs...'
	@cd cmd/main && swag init

## run/api: run the cmd/api application
run/api:
	@echo 'Starting server...'
	go run ./cmd/main

## db/psql: connect to the database using psql and docker
db/psql:
	@echo 'connecting to db'
	docker exec -it post-db bash
	psql postgres://itojudb:itojudb@localhost/itojudb?sslmode=disable
 
## db/migrate/up: apply all up database migrations
db/migrate/up:
	echo 'Running up migrations...'
	@cd internal/sql/migrations/ && goose postgres postgres://itojudb:itojudb@localhost/itojudb up && goose postgres postgres://koyeb-adm:rcHo1Ck7BYmf@ep-tiny-mode-a2d0vyca.eu-central-1.pg.koyeb.app/Itoju-ky up && goose postgres postgres://djjsagev:WG11sRXwe2q1C0I9-3XhTZywTnhbZQPJ@stampy.db.elephantsql.com/djjsagev up

db/migrate/upt:
	echo 'Running up migrations...'
	@cd internal/sql/migrations/ && goose postgres postgres://role:password@localhost:5432/recall_king?sslmode=disable up 

db/migrate/downt:
	@echo 'Running down migrations...'
	@cd internal/sql/migrations/ && goose postgres postgres://role:password@localhost:5432/recall_king?sslmode=disable down

## db/migrate/down: apply all down database migrations
.PHONY: db/migrate/down
db/migrate/down:
	@echo 'Running down migrations...'
	@cd internal/sql/migrations/ && goose postgres postgres://role:password@localhost:5432/recall_king?sslmode=disable down

# ==================================================================================== # 
# QUALITY CONTROL 
# ==================================================================================== # 
## audit: tidy dependencies and format, vet and test all code 
 
audit: vendor
	@echo 'Formatting code...' 
	go fmt ./... 
	@echo 'Vetting code...' 
	go vet ./... 
	staticcheck ./... 
	@echo 'Running tests...' 
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies  
vendor: 
	@echo 'Tidying and verifying module dependencies...' 
	go mod tidy 
	go mod verify 
	@echo 'Vendoring dependencies...' 
	go mod vendor

# ==================================================================================== # 
# BUILD 
# ==================================================================================== # 

## build/api: build the cmd/api application 
build/api: 
	@echo 'Building cmd/api...' 
	env GOOS=linux CGO_ENABLED=0 go build -o bin/${recall_BINARY} ./cmd/main
	# go build -ldflags='-s' -o=./bin/api ${recall_BINARY} ./cmd/main
	# GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

## build/docker: build the docker application 
build/docker: build/api
	@echo 'Building docker...' 
	docker build --platform linux/amd64 -t recall-king .
	# docker build --platform linux/amd64 -t goliathoh/recall-king:latest .



run/docker: build/docker
	@echo 'Building docker...' 
	# docker run -e DB_URL=postgres://djjsagev:WG11sRXwe2q1C0I9-3XhTZywTnhbZQPJ@stampy.db.elephantsql.com/djjsagev itojuapp
	docker run -d --name recall-king-api --network recall-king-network -p 8080:8080 recall-king
