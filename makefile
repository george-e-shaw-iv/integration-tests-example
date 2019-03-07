run: stop up

mod:
	# This make rule requires Go 1.11+
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

up:
	docker-compose -f docker-compose.yml up -d --build

stop:
	docker-compose -f docker-compose.yml stop

down:
	docker-compose -f docker-compose.yml down

test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down --volumes

test-db-up:
	docker-compose -f docker-compose.test.yml up --build db

test-db-down:
	docker-compose -f docker-compose.test.yml down --volumes testdb

# Kubernetes Rules

# Build and tag containers
tag:
	docker build -t georgeeshawiv/listd:1.2 -f cmd/listd/deploy/Dockerfile .
	docker push georgeeshawiv/listd:1.2

# Add example.com as a host for the ingress resource
add-host:
	echo "$$(minikube ip) list.example.com" | sudo tee -a /etc/hosts

# Make sure minikube is started before running this
kube-up:
	kubectl create -f kubernetes/namespace.yaml
	kubectl create -f cmd/listd/deploy/postgres/deployment.yaml
	kubectl create -f cmd/listd/deploy/postgres/service.yaml
	kubectl create -f cmd/listd/deploy/deployment.yaml
	kubectl create -f cmd/listd/deploy/service.yaml
	kubectl create -f kubernetes/ingress.yaml

kube-down:
	kubectl delete -f kubernetes/ingress.yaml
	kubectl delete -f cmd/listd/deploy/service.yaml
	kubectl delete -f cmd/listd/deploy/deployment.yaml
	kubectl delete -f cmd/listd/deploy/postgres/service.yaml
	kubectl delete -f cmd/listd/deploy/postgres/deployment.yaml
	kubectl delete -f kubernetes/namespace.yaml
