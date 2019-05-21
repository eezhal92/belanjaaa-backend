package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pb "github.com/eezhal92/belanjaaa/proto"
	"google.golang.org/grpc"
)

func main() {
	gateway := NewAPIGateway()
	http.HandleFunc("/", gateway.List)

	address := ":9000"
	fmt.Printf("Server started on %s\n", address)

	err := http.ListenAndServe(address, nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}

// Create New instance of GRPC Client for product service
func NewProductService() pb.ProductServiceClient {
	serverAddr := "product:8080"
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect to %s", serverAddr)
	}

	return pb.NewProductServiceClient(conn)
}

type APIGateway struct{}

func NewAPIGateway() *APIGateway {
	return &APIGateway{}
}

func (gateway *APIGateway) List(w http.ResponseWriter, r *http.Request) {
	client := NewProductService()
	request := pb.SearchRequest{Query: "Hey"}
	res, err := client.Search(context.Background(), &request)

	if err != nil {
		log.Fatalf("cound not search product %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	result, err := json.Marshal(res.Products)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(result)
}
