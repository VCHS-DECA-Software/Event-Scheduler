PROTO = $(addprefix proto/,accounts.proto associations.proto errors.proto)

protobuf:
	protoc --go_out=. --go-grpc_out=. $(PROTO)
