protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/carbonemission.proto

go run proveCarbonEmissionCalculation.go verifyCarbonEmissionCalculation.go  main.go -run-proving-service=true -run-verification-service=true