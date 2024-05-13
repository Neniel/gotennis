up:
	@echo "🎾"
	docker-compose -f docker-compose.yml down
	docker-compose -f docker-compose.yml build
	docker-compose -f docker-compose.yml up -d

update-dependencies:
	@echo "\033[1;33m🎾 Updating dependencies 🎾\033[0m"
	cd app ; go get -u ./... ; go mod tidy
	cd cache ; go get -u ./... ; go mod tidy
	cd categories ; go get -u ./... ; go mod tidy
	cd database ; go get -u ./... ; go mod tidy
	cd entity ; go get -u ./... ; go mod tidy
	cd players ; go get -u ./... ; go mod tidy
	cd util ; go get -u ./... ; go mod tidy

gen-mocks:
	@echo "\033[1;33m🎾 Generating mocks 🎾\033[0m"
	mockgen -source=./app/app.go -destination=./app/mock_app.go -package=app
	mockgen -source=./database/database.go -destination=./database/mock_database.go -package=database

test:
	@echo "\033[1;33m🎾 Running tests 🎾\033[0m"
	go test ./app/...
	go test ./cache/...
	go test ./categories/...
	go test ./database/...
	go test ./entity/...
	go test ./players/...
	go test ./util/...

msg:
	@echo "\033[1;31m🎾 Hola mundo! 🎾\033[0m"
	@echo "\033[1;32m🎾 Hola mundo! 🎾\033[0m"
	@echo "\033[1;33m🎾 Hola mundo! 🎾\033[0m"
	@echo "\033[1;34m🎾 Hola mundo! 🎾\033[0m"
