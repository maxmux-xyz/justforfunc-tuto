# Summary

GRPC server - CLI client.

1. pb: https://www.youtube.com/watch?v=_jQ3i_fyqGA&t=84s
2. grpc: https://www.youtube.com/watch?v=uolTUtioIrc&t=68s
3. 


(1) Notes with GRPC
```bash
go env -w GO111MODULE=off #optional
go get -u google.golang.org/grpc
go get -u google.golang.org/grpc
cd proto && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative todo.proto
cd cmd/server && go build main.go && mv server ../..
cd ../.. && ./server # Server is running and has access to the db in the root folder
```

Open another terminal
```bash
go run main.go add newtask
go run main.go add list
# ...
```


