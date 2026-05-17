package main

import (
	"Server/internal/handler"
	"Server/internal/store"
	"context"
	"flag"
	"time"

	_ "Server/internal/docs"

	"github.com/joho/godotenv"
)

const (
	modeProduction  = "production"
	modeDevelopment = "development"
)

// @title Social Media API
// @version 0.1
// @description This is a social media API that allows users to create posts, like posts, and follow other users.
// @basePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type Bearer followed by a space and then the token.
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
	server.RegisterRoutes()

	if *mode == modeDevelopment {
		server.UseSwagger()
	}

	server.Start()
}
