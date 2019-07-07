package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/EdSwArchitect/learn-grpc/myservice"
	elasticsearch "github.com/elastic/go-elasticsearch"
	"google.golang.org/grpc"
)

// the type passed around for the server
type internalServerStuff struct{}

// GetIndices is the gRpc server method call
func (s *internalServerStuff) GetIndices(ctx context.Context, server *pb.EServer) (*pb.Result, error) {

	fmt.Printf("Indicies in with: %s\n", server.GetServer())

	// fmt.Printf("The context is: %+v\n", ctx)

	// var rez *pb.Result
	r := new(pb.Result)

	r.Code = 100
	r.Data = "GetIndicies '" + server.GetServer() + "'"
	return r, nil
}

// GetStatus is the gRpc server method call
func (s *internalServerStuff) GetStatus(ctx context.Context, server *pb.EServer) (*pb.Result, error) {

	fmt.Printf("Status in with: %s\n", server.GetServer())

	// fmt.Printf("The context is: %+v\n", ctx)

	// var rez *pb.Result
	r := new(pb.Result)

	r.Code = 100
	r.Data = "GetStatus '" + server.GetServer() + "'"
	return r, nil
}

// GetStatus is the gRpc server method call
func (s *internalServerStuff) QueryIndex(ctx context.Context, server *pb.Query) (*pb.QueryResult, error) {

	fmt.Printf("QueryIndex in with: %s: Size: %d. Start: %d\n", server.GetServer(), server.GetSize(), server.GetStart())

	// fmt.Printf("The context is: %+v\n", ctx)

	host := "http://" + server.Server + ":9200"

	cfg := elasticsearch.Config{
		Addresses: []string{
			host,
		},
	}

	// es, _ := elasticsearch.NewDefaultClient()
	es, _ := elasticsearch.NewClient(cfg)

	log.Println(elasticsearch.Version)
	log.Println(es.Info())

	/*
		version := elasticsearch.Version

		info, _ := es.Info()

		fmt.Printf("The version: %s\n", version)
		fmt.Printf("Info: %+v\n", info)

		fmt.Printf("-- size of query: %d\n", server.Size)
		fmt.Printf("-- start of query %d\n", server.Start)
		fmt.Printf("-- query '%s'\n", server.Query)
		fmt.Printf("-- index '%s'\n", server.Index)
		fmt.Printf("-- term '%s'\n", server.Term)
		fmt.Printf("-- server: %+v\n", server)
	*/

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				server.Term: server.Query,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("Error encoding query: %s", err.Error())
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(server.Index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
		es.Search.WithSize(int(server.Size)),
		es.Search.WithFrom(int(server.Start)),
	)
	if err != nil {
		return nil, fmt.Errorf("Error getting response: %s", err.Error())
	}

	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			// log.Fatalf("Error parsing the response body: %s", err)

			return nil, fmt.Errorf("Error parsing the response body: %s", err.Error())

		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)

			return nil, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"])
		}
	}

	buffer := new(bytes.Buffer)

	buffer.ReadFrom(res.Body)

	resultsString := buffer.String()

	/*

		var r map[string]interface{}

		res.B

		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		// Print the response status, number of results, and request duration.
		log.Printf(
			"[%s] %d hits; took: %dms",
			res.Status(),
			int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
			int(r["took"].(float64)),
		)
		// Print the ID and document source for each hit.

		var st string
		var results string

		results = fmt.Sprintf("{")

		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			// log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])

			// log.Printf("%s\n----\n", hit.(map[string]interface{})["_source"])

			goober := hit.(map[string]interface{})["_source"]

			log.Printf("What is this: %+v\n-----\n\n", goober)

			element := goober.(map[string]interface{})

			for k, v := range element {
				fmt.Printf("Key: '%s' - Value: '%+v'\n", k, v)
			}

		}

	*/

	// var rez *pb.Result
	rz := new(pb.QueryResult)

	rz.Results = resultsString
	rz.Size = 0
	rz.Code = 200

	return rz, nil
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
