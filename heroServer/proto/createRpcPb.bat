protoc --proto_path=%GOPATH%/go/project_go/heroServer/proto --go_out=plugins=grpc:%GOPATH%/go/project_go %GOPATH%/go/project_go/heroServer/proto/*.proto
pause