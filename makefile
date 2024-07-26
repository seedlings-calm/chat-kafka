protoc: 
	@./scripts/proto.sh

tree:
	@./scripts/tree.sh
kafka:
	@./scripts/kafka.sh

.PHONY: protoc tree kafka



dev: 
	go run ./ 