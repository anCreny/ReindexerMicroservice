start-microservice:
	CGO_ENABLED=0 GOOS=linux go build -o ./build/pkg main.go
	docker-compose -f ./build/pkg/docker-compose.yml build
	docker-compose -f ./build/pkg/docker-compose.yml up

stop-microservice:
	rm -rf ./build/pkg/main
	docker-compose -f ./build/pkg/docker-compose.yml down -v

restart-microservice: stop-microservice start-microservice



start-database:
	docker-compose -f ./third_party/build/pkg/docker-compose.yml build
	docker-compose -f ./third_party/build/pkg/docker-compose.yml up

stop-database:
	docker-compose -f ./third_party/build/pkg/docker-compose.yml down -v

restart-database: stop-database start-database