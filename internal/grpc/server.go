package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	product "github.com/berezovskyivalerii/server-rpc-csv/proto"
)

type Server struct {
	grpcServ      *grpc.Server
	productServer product.ProductServiceServer
}

func New(productServer product.ProductServiceServer) *Server {
	return &Server{
		grpcServ:      grpc.NewServer(),
		productServer: productServer,
	}
}

func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	product.RegisterProductServiceServer(s.grpcServ, s.productServer)

	if err := s.grpcServ.Serve(lis); err != nil {
		return err
	}

	return nil
}
