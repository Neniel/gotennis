SERVICE ?=
IMAGE ?=
TAG ?=
ENV ?=

deploy-db:
	docker-compose -f docker-compose-db.yml down
	docker-compose -f docker-compose-db.yml build
	docker-compose -f docker-compose-db.yml up -d

MICROSERVICE ?=
VERSION ?=
buildv2:
ifeq ($(MICROSERVICE),)
	@echo "Please provide a value for MICROSERVICE"
else
ifeq ($(VERSION),)
	@echo "Please provide a value for VERSION"
else
	docker build ${MICROSERVICE} -t neniel/tennis-${MICROSERVICE}:${VERSION}
endif
endif

pushv2:
ifeq ($(MICROSERVICE),)
	@echo "Please provide a value for MICROSERVICE"
else
ifeq ($(VERSION),)
	@echo "Please provide a value for VERSION"
else
	docker push neniel/tennis-${MICROSERVICE}:${VERSION}
endif
endif



build:
	docker-compose -f docker-compose.yml build ${SERVICE}

deploy:
ifeq ($(ENV),)
	@echo "🎾"
	docker-compose -f docker-compose.yml down ${SERVICE}
	docker-compose -f docker-compose.yml build ${SERVICE}
	docker-compose --verbose -f docker-compose.yml up -d ${SERVICE}
else
	docker-compose -f docker-compose-${ENV}.yml down ${SERVICE}
	docker-compose -f docker-compose-${ENV}.yml build ${SERVICE}
	docker-compose --verbose -f docker-compose-${ENV}.yml up -d ${SERVICE}
endif

push:
ifeq ($(IMAGE),)
	echo "Error: IMAGE is not set!"
	false
endif

ifneq ($(TAG),)
	docker push neniel/tennis-${IMAGE}:${TAG}
else
	docker push neniel/tennis-${IMAGE}:latest
endif

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
  --timeout=180s

	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.5/cert-manager.yaml
	cmctl check api --wait=2m

k8s-deploy:
ifneq ($(ENV),)
	helm install tennis-app .deployment/k8s/tennis-app -f .deployment/k8s/tennis-app/values-${ENV}.yaml
else
	helm install tennis-app .deployment/k8s/tennis-app -f .deployment/k8s/tennis-app/values.yaml
endif

k8s-redeploy:
ifneq ($(ENV),)
	helm upgrade tennis-app .deployment/k8s/tennis-app -f .deployment/k8s/tennis-app/values-${ENV}.yaml
else
	helm upgrade tennis-app .deployment/k8s/tennis-app -f .deployment/k8s/tennis-app/values.yaml
endif

k8s-uninstall:
	helm uninstall tennis-app

PROJECT_PATH ?=
update-dependencies:
ifeq ($(PROJECT_PATH),)
	echo Please prvide a PROJECT_PATH
else
	@echo "\033[1;33m🎾 Updating dependencies 🎾\033[0m"
	cd ${PROJECT_PATH} ; go get -u ./... ; go mod tidy
endif

gen-mocks:
	@echo "\033[1;33m🎾 Generating mocks 🎾\033[0m"
	mockgen -source=./lib/app/app.go -destination=./lib/app/mock_app.go -package=app
	mockgen -source=./lib/database/database.go -destination=./lib/database/mock_database.go -package=database

test:
	@echo "\033[1;33m🎾 Running tests 🎾\033[0m"
	go test ./cache/...
	go test ./categories/...
	go test ./players/...
	go test ./tournaments/...
	go test ./lib/app/...
	go test ./lib/database/...
	go test ./lib/entity/...
	go test ./lib/util/...

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
