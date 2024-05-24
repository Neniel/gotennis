SERVICE ?=
IMAGE ?=

deploy-db:
	docker-compose -f docker-compose-db.yml down
	docker-compose -f docker-compose-db.yml build
	docker-compose -f docker-compose-db.yml up -d

deploy:
	@echo "🎾"
	docker-compose -f docker-compose-db.yml down
	docker-compose -f docker-compose-db.yml build
	docker-compose -f docker-compose-db.yml up -d
	docker-compose -f docker-compose-ms.yml down ${SERVICE}
	docker-compose -f docker-compose-ms.yml build ${SERVICE}
	docker-compose --verbose -f docker-compose-ms.yml up -d ${SERVICE}

build:
	docker-compose build ${SERVICE}


push:
	docker push neniel/tennis-${IMAGE}:latest

k8s-deploy:
	helm install tennis-app .deployment/k8s/tennis-app -f .deployment/k8s/tennis-app/values-dev.yaml

k8s-redeploy:
	helm upgrade tennis-app .deployment/k8s/tennis-app -f .deployment/k8s/tennis-app/values-dev.yaml

k8s-uninstall:
	helm uninstall tennis-app

k8s-delete-cluster-kind:
	kind delete cluster

k8s-configure-minikube:
	minikube start
	minikube addons enable metrics-server
	minikube addons enable ingress
	minikube dashboard

k8s-configure-kind:
	kind create cluster --config .deployment/k8s/kind/configuration.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
	kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

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
