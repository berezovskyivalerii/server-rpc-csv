package repository

import (
	"github.com/berezovskyivalerii/server-rpc-csv/pkg/domain"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	db *mongo.Database
}

func NewProduct(db *mongo.Database) *Product {
	return &Product{
		db: db,
	}
}

func (r *Product) Fetch(ctx context.Context, req []domain.Product) error {
	collection := r.db.Collection("products")

	for _, p := range req {
		var existingProduct domain.Product
		err := collection.FindOne(ctx, bson.M{"name": p.Name}).Decode(&existingProduct)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				// Продукта нет, создаём новый
				_, err = collection.InsertOne(ctx, bson.M{
					"name":               p.Name,
					"price":              p.Price,
					"last_updated":       time.Now(),
					"last_request":       time.Now(),
					"price_change_count": 0,
				})
				if err != nil {
					return err
				}
				continue
			}
			return err
		}

		update := bson.M{}
		if existingProduct.Price != p.Price {
			// Цена изменилась
			update = bson.M{
				"$set": bson.M{
					"price":        p.Price,
					"last_updated": time.Now(),
				},
				"$inc": bson.M{"price_change_count": 1},
			}
		} else {
			// Цена не изменилась, обновляем только last_request
			update = bson.M{
				"$set": bson.M{"last_request": time.Now()},
			}
		}

		_, err = collection.UpdateOne(ctx, bson.M{"name": p.Name}, update, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Product) List(ctx context.Context, req domain.ListRequest) (*domain.ListResponse, error){
	return nil, nil
}