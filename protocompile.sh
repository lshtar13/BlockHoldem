protoc --proto_path=./node/global --go_out=./node/global --go_opt=paths=source_relative\
  --go-grpc_out=./node/global --go-grpc_opt=paths=source_relative node/global/*.proto 
