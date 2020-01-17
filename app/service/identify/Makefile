
GMICRO_HOME:=/home/liangjf/ljf_home/code/go_home/study/gmicro/app/service
MODULE_NAME:=identify

.PHONY: proto
proto:
	protoc --proto_path=${GMICRO_HOME}/${MODULE_NAME}/proto/v1:. \
	--micro_out=${GMICRO_HOME}/${MODULE_NAME}/proto/v1 \
	--go_out=${GMICRO_HOME}/${MODULE_NAME}/proto/v1 \
	${GMICRO_HOME}/${MODULE_NAME}/proto/v1/api.proto

.PHONY: build
build: proto

	go build -o cmd/${MODULE_NAME}-srv cmd/main.go cmd/plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t ${MODULE_NAME}-srv:latest
