run:
	go run main.go

start:
	up run

up:
	cd docker/ && docker-compose up -d 

down:
	cd docker && docker-compose down

test:
	go test -check.vv

check_install:
	which swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger:
	swagger generate spec -o ./swagger.yaml --scan-models && swagger serve -F=swagger swagger.yaml