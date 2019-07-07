# learn-grpc

To create the stub code from the protofile.

protoc -I myservice/ myservice/my_service.proto --go_out=plugins=grpc:myservice

