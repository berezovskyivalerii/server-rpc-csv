package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/berezovskyivalerii/server-rpc-csv/pkg/csvparser"
	"github.com/berezovskyivalerii/server-rpc-csv/pkg/domain"
	product "github.com/berezovskyivalerii/server-rpc-csv/proto"
)

type ProductRepository interface {
	Fetch(ctx context.Context, req []domain.Product) error
	List(ctx context.Context, req domain.ListRequest) (*domain.ListResponse, error)
}

type Product struct {
	repo ProductRepository
}

func NewProduct(repo ProductRepository) *Product {
	return &Product{
		repo: repo,
	}
}

func (s *Product) Fetch(ctx context.Context, req *product.FetchRequest) (*product.FetchResponse, error) {
	// Сделать запрос на сервис
	resp, err := http.Get(req.Url)
	if err != nil {
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
		for i := 0; i < len(products); i++ {
			fmt.Println(products[i])
		}
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
	// конвертация gRPC-запроса в domain-запрос
	domainReq := domain.ListRequest{
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
		SortField:  req.SortField,
		SortOrder:  req.SortOrder,
	}

	// вызов метода репозитория
	domainResp, err := s.repo.List(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Конвертация ответа репозитория в gRPC-ответ
	protoProducts := make([]*product.Product, len(domainResp.Products))
	for i, p := range domainResp.Products {
		protoProducts[i] = &product.Product{
			ProductName:      p.Name,
			Price:            p.Price,
			PriceChangeCount: p.PriceChangeCount,
			LastUpdated:      p.LastUpdated.Format(time.RFC3339), // конвертация даты в строку
		}
	}

	return &product.ListResponse{
		Products:      protoProducts,
		TotalProducts: domainResp.TotalProducts,
	}, nil
}
