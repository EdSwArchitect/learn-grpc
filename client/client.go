package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/EdSwArchitect/learn-grpc/myservice"
	"google.golang.org/grpc"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func main() {
	fmt.Println("Client says, hi, Ed")

	flag.Parse()

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	fmt.Printf("Client connecting to address %s\n", *serverAddr)

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Cient making the connection to server\n")

	client := pb.NewMyServiceClient(conn)

	fmt.Printf("The connecton was established: %+v. \n", client)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := client.GetIndices(ctx, &pb.EServer{
		Server: "sabrina.imac",
	})

	if err != nil {
		fmt.Println("Well, that client call failed")
		log.Panic(err)
	}

	fmt.Printf("The result is: %s\n", result.Data)

}
