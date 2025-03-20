package service

import (
	"berezovskyivalerii/server-rpc-csv/pkg/domain"
	product "berezovskyivalerii/server-rpc-csv/internal/grpc"
	"berezovskyivalerii/server-rpc-csv/pkg/csvparser"
	"context"
	"net/http"
)

type ProductRepository interface{
	Fetch(ctx context.Context, req []domain.Product) error
	List(ctx context.Context, req *product.ListRequest) (*product.ListResponse, error)
}

type Product struct{
	repo ProductRepository
}

func NewProduct(repo ProductRepository) *Product{
	return &Product{
		repo: repo,
	}
}

func (s *Product) Fetch(ctx context.Context, req *product.FetchRequest) (*product.FetchResponse, error) {
	// Сделать запрос на сервис
	resp, err := http.Get(req.Url)
	if err != nil{
		return &product.FetchResponse{
			Success: false,
			Message: "Failed to fetch CSV: " + err.Error(),
		}, nil
	}
	defer resp.Body.Close()

	//Проверить статус возврата
	if resp.StatusCode != http.StatusOK {
		return &product.FetchResponse{
			Success: false,
			Message: "Received non-OK response from URL: " + resp.Status,
		}, nil
	}

	// Распарсить ответ в []domain.Product
	products, err := csvparser.ParseCSV(resp.Body)
	if err != nil {
		return &product.FetchResponse{
			Success: false,
			Message: "Failed to parse CSV: " + err.Error(),
		}, nil
	}

	// Отправить на уровень репозитория
	err = s.repo.Fetch(ctx, products)
	if err != nil {
		return &product.FetchResponse{
			Success: false,
			Message: "Failed to save products to repository: " + err.Error(),
		}, nil
	}

	// Успешно
	return &product.FetchResponse{
		Success: true,
		Message: "Successfully fetched and stored products.",
	}, nil
}

func (s *Product) List(ctx context.Context, req *product.ListRequest) (*product.ListResponse, error) {
	return s.repo.List(ctx, req)
}