run:
	docker-compose up --build -d
stop:
	docker-compose down
clean:
	docker-compose down -v --rmi all
test:
	go test -v `go list ./... | grep -v /tests`
integration-test:
	go test -v ./tests/...