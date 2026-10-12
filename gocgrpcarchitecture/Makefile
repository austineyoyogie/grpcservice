client:
	@air -c .air.clientend.toml

server:
	@air -c .air.backend.toml

genProto:
	@protoc --proto_path=protobuf \
            --go_out=generate/proto --go_opt=paths=source_relative \
            --go-grpc_out=generate/proto  --go-grpc_opt=paths=source_relative \
           protobuf/userProto/userProto.proto

# protobuf/documentProto/documentProto.proto

userProto:
	@protoc \
		--proto_path=protobuf "protobuf/userProto/userProto.proto" \
		--go_out=generate/proto --go_opt=paths=source_relative \
  	    --go-grpc_out=generate/proto --go-grpc_opt=paths=source_relative

documentProto:
	@protoc \
		--proto_path=protobuf "protobuf/documentProto/documentProto.proto" \
		--go_out=generate/proto/documentProto --go_opt=paths=source_relative \
  	--go-grpc_out=generate/proto/documentProto --go-grpc_opt=paths=source_relative