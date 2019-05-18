package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/eezhal92/belanjaaa/proto"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("p", 8080, "port to listen to")
	flag.Parse()

	fmt.Printf("Listening to port %d", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf(err.Error())
	}

	storage := NewProductStorage()
	storage.Add(&pb.Product{Id: "product-1", Title: "Golang"})
	storage.Add(&pb.Product{Id: "product-2", Title: "GRPC"})
	storage.Add(&pb.Product{Id: "product-3", Title: "Kubernetes"})

	s := grpc.NewServer()
	server := Server{Storage: storage}
	pb.RegisterProductServiceServer(s, &server)
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf(err.Error())
	}
}

type Server struct {
	Storage *ProductStorage
}

type ProductStorage struct {
	items map[string]*pb.Product
}

func NewProductStorage() *ProductStorage {
	return &ProductStorage{
		items: make(map[string]*pb.Product, 0),
	}
}

func (productStorage *ProductStorage) Add(product *pb.Product) {
	if _, ok := productStorage.items[product.Id]; !ok {
		productStorage.items[product.Id] = product
	}
}

func (s *Server) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	products := make([]*pb.Product, 0)

	for _, product := range s.Storage.items {
		products = append(products, product)
	}

	return &pb.SearchResponse{Products: products}, nil
}

func (s *Server) FindById(ctx context.Context, request *pb.FindByIdRequest) (*pb.FindByIdResponse, error) {
	if product, ok := s.Storage.items[request.Id]; ok {
		return &pb.FindByIdResponse{Product: product}, nil
	}

	return nil, errors.New("Not found")
}

func (s *Server) Add(ctx context.Context, request *pb.AddRequest) (*pb.AddResponse, error) {
	if _, ok := s.Storage.items[request.Id]; ok {
		return nil, errors.New(fmt.Sprintf("Id of %s has been taken", request.Id))
	}

	product := &pb.Product{
		Id:    request.Id,
		Title: request.Title,
	}

	s.Storage.Add(product)

	return &pb.AddResponse{Product: product}, nil
}
