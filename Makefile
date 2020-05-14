protoc -I simplechat/ simplechat/simplechat.proto --go_out=plugins=grpc:simplechat --go_opt=paths=source_relative
