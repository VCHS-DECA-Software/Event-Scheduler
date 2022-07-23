PROTO_PATH = $(abspath .)
PROTO = $(addprefix proto/,request.proto types.proto)

rpc:
	protoc \
		--proto_path=. \
		--go_out=. \
		--go-grpc_out=. \
		$(PROTO)

test:
	cd tests/components && \
		cd common && go test
