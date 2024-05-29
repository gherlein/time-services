
TSPB := ./timeservice-pb


proto-clean:
	-rm ${TSPB}/*.go
	ls ${TSPB}



proto: proto-clean
	@echo "creating go files from .proto files"
	protoc -I${TSPB} --go_out=${TSPB} --go-grpc_out=${TSPB} --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative ${TSPB}/timeservice.proto 
	ls ${TSPB}


run-server:
	go run server/main.go


run-client:
	go run client/main.go	
