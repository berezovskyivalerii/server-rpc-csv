package grpc

import (
	"context"
)

type ProductService interface {
	Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error)
	List(ctx context.Context, req *ListRequest) (*ListResponse, error)
}

type ProductServer struct {
	UnimplementedProductServiceServer
	service ProductService
}

func NewProductServer(service ProductService) *ProductServer {
	return &ProductServer{
		service: service,
	}
}

func (s *ProductServer) Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error) {
	return s.service.Fetch(ctx, req)
}

func (s *ProductServer) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	return s.service.List(ctx, req)
}
