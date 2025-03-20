package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/berezovskyivalerii/server-rpc-csv/internal/config"
	"github.com/berezovskyivalerii/server-rpc-csv/internal/grpc"
	"github.com/berezovskyivalerii/server-rpc-csv/internal/repository"
	"github.com/berezovskyivalerii/server-rpc-csv/internal/service"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("Programm started", time.Now())

	//Config init
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	//Context init
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Mongo connection
	opts := options.Client()
	opts.SetAuth(options.Credential{
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
	})
	opts.ApplyURI(cfg.DB.URI)

	dbClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := dbClient.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mongo started", time.Now())

	//Database init
	db := dbClient.Database(cfg.DB.Database)

	//repo, service, server init
	productRepo := repository.NewProduct(db)
	productService := service.NewProduct(productRepo)
	productSrv := grpc.NewProductServer(productService)
	srv := grpc.New(productSrv)
	fmt.Println("Server started", time.Now())

	if err := srv.ListenAndServe(cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}