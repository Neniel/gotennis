up:
	docker-compose -f docker-compose.yml down
	docker-compose -f docker-compose.yml build
	docker-compose -f docker-compose.yml up -d

update-dependencies:
	cd app ; go get -u ./... ; go mod tidy
	cd cache ; go get -u ./... ; go mod tidy
	cd categories ; go get -u ./... ; go mod tidy
	cd database ; go get -u ./... ; go mod tidy
	cd entity ; go get -u ./... ; go mod tidy
	cd players ; go get -u ./... ; go mod tidy
	cd util ; go get -u ./... ; go mod tidy

gen-mocks:
	mockgen -source=./app/app.go -destination=./app/mock_app.go -package=app
	mockgen -source=./database/database.go -destination=./database/mock_database.go -package=database

test:
	cd app ; go test ./...
	cd cache ; go test ./...
	cd categories ; go test ./...
	cd database ; go test ./...
	cd entity ; go test ./...
	cd players ; go test ./...
	cd util ; go test ./...
