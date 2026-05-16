package main

import (
	"Server/internal/handler"
	"Server/internal/store"
	"context"
	"flag"
	"time"

	"github.com/joho/godotenv"
)

const (
	modeProduction  = "production"
	modeDevelopment = "development"
)

func main() {
	mode := flag.String("mode", modeProduction, "mode of operation (production or development)")
	flag.Parse()

	if err := godotenv.Load(".env." + *mode); err != nil {
		panic(err)
	}

	db, err := store.Connect()
	if err != nil {
		panic(err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db.Client.Disconnect(ctx)
	}()

	server := handler.NewServer(db)

	server.Start()
}
