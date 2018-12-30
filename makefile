run: stop up

up:
	docker-compose -f ./docker-compose.yml up --build

stop:
	docker-compose stop

down:
	docker-compose down

test:
	docker-compose -f ./docker-compose.test.yml up --build