deps:
	go get github.com/gogo/protobuf/proto
	go get github.com/gogo/protobuf/protoc-gen-gogo
	go get github.com/gogo/protobuf/gogoproto

generate:
	protoc -I. -I$GOPATH/src/github.com/gogo/protobuf/gogoproto -I$GOPATH/src/github.com/gogo/protobuf/protobuf/ --gogo_out=. *.proto