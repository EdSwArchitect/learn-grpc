package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/EdSwArchitect/learn-grpc/myservice"
	"google.golang.org/grpc"
)

// the type passed around for the server
type internalServerStuff struct{}

// GetIndices is the gRpc server method call
func (s *internalServerStuff) GetIndices(ctx context.Context, server *pb.EServer) (*pb.Result, error) {

	fmt.Printf("Came in with: %s\n", server.GetServer())

	// fmt.Printf("The context is: %+v\n", ctx)

	// var rez *pb.Result
	r := new(pb.Result)

	r.Code = 100
	r.Data = "Hi, Ed '" + server.GetServer() + "'"
	return r, nil
}

var (
	// tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	// certFile   = flag.String("cert_file", "", "The TLS cert file")
	// keyFile    = flag.String("key_file", "", "The TLS key file")
	// jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port = flag.Int("port", 10000, "The server port")
)

// create the type for the server interface
func newServer() *internalServerStuff {
	s := new(internalServerStuff)

	return s

}
func main() {
	fmt.Println("Hi, Ed")

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	fmt.Println("Registering service")

	pb.RegisterMyServiceServer(grpcServer, newServer())

	fmt.Println("About to server gRpc requests, yo!")

	grpcServer.Serve(lis)

}
