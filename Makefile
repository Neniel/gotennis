SERVICE ?=
IMAGE ?=

deploy:
	@echo "🎾"
	docker-compose -f docker-compose.yml down ${SERVICE}
	docker-compose -f docker-compose.yml build ${SERVICE}
	docker-compose -f docker-compose.yml up -d ${SERVICE}

build:
	docker-compose build ${SERVICE}


push:
	docker push neniel/tennis-${IMAGE}:latest

k8s-dashboard:
	minikube dashboard

k8s-apply:
	kubectl apply -f kubernetes/secrets.yml
	kubectl apply -f kubernetes/mongo-deployment.yml
	kubectl apply -f kubernetes/ingress.yml
	kubectl apply -f kubernetes/categories-deployment.yml
	kubectl apply -f kubernetes/players-deployment.yml

k8s-reset:
	kubectl delete deployment --all
	kubectl delete service --all
	kubectl delete statefulset --all
	kubectl delete daemonset --all
	kubectl delete configmap --all
	kubectl delete secret --all
	kubectl delete ingress --all
	kubectl delete pod --all

update-dependencies:
	@echo "\033[1;33m🎾 Updating dependencies 🎾\033[0m"
	cd app ; go get -u ./... ; go mod tidy
	cd cache ; go get -u ./... ; go mod tidy
	cd categories ; go get -u ./... ; go mod tidy
	cd database ; go get -u ./... ; go mod tidy
	cd entity ; go get -u ./... ; go mod tidy
	cd players ; go get -u ./... ; go mod tidy
	cd util ; go get -u ./... ; go mod tidy
	git add .
	git commit -m "update dependencies"
	git push

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

run-alloy:
	docker rm grafana_alloy -f
	docker run -d --name=grafana_alloy -v ./grafana/alloy/config.alloy:/etc/alloy/config.alloy -p 12345:12345 grafana/alloy:latest run --server.http.listen-addr=0.0.0.0:12345 --storage.path=/var/lib/alloy/data /etc/alloy/config.alloy

prometheus:
	docker stop prometheus
	docker rm prometheus
	docker run -d --name=prometheus -p 9090:9090 -v ./prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

grafana:
	docker rm grafana
	docker run -d --name=grafana -p 3000:3000 grafana/grafana-enterprise

msg:
	@echo "\033[1;31m🎾 Hola mundo! 🎾\033[0m"
	@echo "\033[1;32m🎾 Hola mundo! 🎾\033[0m"
	@echo "\033[1;33m🎾 Hola mundo! 🎾\033[0m"
	@echo "\033[1;34m🎾 Hola mundo! 🎾\033[0m"
