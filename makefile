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
	docker-compose -f docker-compose.test.yml down

testdb-up:
	docker-compose -f docker-compose.test.yml up --build testdb

testdb-down:
	docker-compose -f docker-compose.test.yml down testdb

# Kubernetes Rules

# Build and tag containers
tag:
	docker build -t georgeeshawiv/listd:1.1 -f cmd/listd/deploy/Dockerfile .
	docker push georgeeshawiv/listd:1.1

# Add example.com as a host for the ingress resource
add-host:
	echo "$$(minikube ip) example.com" | sudo tee -a /etc/hosts

# Make sure minikube is started before running this
kube-up:
	kubectl create -f k8s/postgres-deployment.yaml
	kubectl create -f k8s/postgres-service.yaml
	kubectl create -f k8s/listd-deployment.yaml
	kubectl create -f k8s/listd-service.yaml
	kubectl create -f k8s/ingress.yaml

kube-down:
	kubectl delete -f k8s/ingress.yaml
	kubectl delete -f k8s/listd-service.yaml
	kubectl delete -f k8s/listd-deployment.yaml
	kubectl delete -f k8s/postgres-service.yaml
	kubectl delete -f k8s/postgres-deployment.yaml