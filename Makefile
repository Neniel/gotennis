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
