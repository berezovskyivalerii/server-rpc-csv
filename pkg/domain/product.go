package domain

import "time"

type Product struct {
	ID               int64     `bson:"id"`
	Name             string    `bson:"name"`
	Price            float64   `bson:"price"`
	PriceChangeCount int32     `bson:"price_change_count"`
	LastUpdated      time.Time `bson:"last_updated"`
}

type ListRequest struct {
	PageNumber int32
	PageSize   int32
	SortField  string
	SortOrder  string
}
