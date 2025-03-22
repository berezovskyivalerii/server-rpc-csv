package repository

import (
	"github.com/berezovskyivalerii/server-rpc-csv/pkg/domain"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
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
		//Проверяем ли существует продукт (проверка по name)
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
					fmt.Println("Failed Inserting new product:", p.Name)
					return err
				}
				fmt.Println("Inserting new product:", p.Name)

				continue
			}
			return err
		}

		// Если продукт есть, но цена поменялась
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
	collection := r.db.Collection("products")

	//Параметры пагинации 
	skip := (req.PageNumber - 1) * req.PageSize
	limit := int64(req.PageSize)

	// Параметры сортировки
	sortOrder := 1 // по умолчанию ASC
	if req.SortOrder == "desc"{
		sortOrder = -1
	}

	sortField := req.SortField
	if sortField == ""{
		sortField = "name"
	}

	// Запрос сортировки
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: sortField, Value: sortOrder}})
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(limit)

	// Запрос данных
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil{
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	for cursor.Next(ctx){
		var product domain.Product
		if err := cursor.Decode(&product); err != nil{
			return nil, err
		}
		products = append(products, product)
	}

	// Подсчёт общего количества документов
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	return &domain.ListResponse{
		Products:      products,
		TotalProducts: int32(totalCount),
	}, nil
}