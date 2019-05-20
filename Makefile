proto:
	protoc -I api -I=/usr/local/include service.proto --go_out=plugins=grpc:api

build-product-service:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o product product-service/main.go

build-api-gateway:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o api api-gateway/main.go

docker:
	docker build -t eezhal92/product-service -f Dockerfile-Product .
	docker build -t eezhal92/api-gateway -f Dockerfile-APIGateway .
