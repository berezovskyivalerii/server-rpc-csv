package grpc

import (
	"context"
	"fmt"
	product "github.com/berezovskyivalerii/server-rpc-csv/proto"
)

type ProductService interface {
	Fetch(ctx context.Context, req *product.FetchRequest) (*product.FetchResponse, error)
	List(ctx context.Context, req *product.ListRequest) (*product.ListResponse, error)
}

type ProductServer struct {
	product.UnimplementedProductServiceServer
	service ProductService
}

func NewProductServer(service ProductService) *ProductServer {
	return &ProductServer{
		service: service,
	}
}

func (s *ProductServer) Fetch(ctx context.Context, req *product.FetchRequest) (*product.FetchResponse, error) {
	fmt.Println("[Fetch]")
	return s.service.Fetch(ctx, req)
}

func (s *ProductServer) List(ctx context.Context, req *product.ListRequest) (*product.ListResponse, error) {
	fmt.Println("[List]")
	return s.service.List(ctx, req)
}
