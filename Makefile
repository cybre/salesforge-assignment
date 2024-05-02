run:
	docker-compose up -d
stop:
	docker-compose down
clean:
	docker-compose down -v --rmi all
test:
	go test -v ./...