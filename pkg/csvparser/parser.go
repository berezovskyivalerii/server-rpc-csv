package csvparser

import (
	"berezovskyivalerii/server-rpc-csv/pkg/domain"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

func ParseCSV(reader io.Reader) ([]domain.Product, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	var products []domain.Product
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) != 2 {
			return nil, errors.New("invalid CSV format")
		}

		price, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		if err != nil {
			return nil, errors.New("invalid price format")
		}

		product := domain.Product{
			Name:             strings.TrimSpace(record[0]),
			Price:            price,
			PriceChangeCount: 0,
			LastUpdated:      time.Now(),
		}
		products = append(products, product)
	}

	return products, nil
}
