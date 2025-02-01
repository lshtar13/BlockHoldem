protoc --proto_path=./node --go_out=./node/global --go_opt=paths=source_relative\
  --go-grpc_out=./node/global --go-grpc_opt=paths=source_relative node/global/*.proto 

protoc --proto_path=./node --go_out=./node/local --go_opt=paths=source_relative\
  --go-grpc_out=./node/local --go-grpc_opt=paths=source_relative node/local/*.proto 

protoc --proto_path=./node --go_out=./node/protos --go_opt=paths=source_relative\
  --go-grpc_out=./node/protos --go-grpc_opt=paths=source_relative node/protos/*.proto 