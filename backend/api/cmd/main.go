package main

import (
	"Server/internal/handler"
	"Server/internal/helpers"
	"Server/internal/server"
	"Server/internal/store"
	"flag"

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

	if *mode != modeProduction && *mode != modeDevelopment {
		panic("invalid mode: " + *mode)
	}

	if err := godotenv.Load(".env." + *mode); err != nil {
		panic(err)
	}

	db, err := store.Connect()
	if err != nil {
		panic(err)
	}

	defer func() {
		ctx, cancel := helpers.GenerateContext()
		defer cancel()
		db.Client.Disconnect(ctx)
	}()

	server := server.NewServer(db)
	server.RegisterServices()

	handler.RegisterRoutes(server)

	if *mode == modeDevelopment {
		server.UseSwagger()
	}

	server.Start()
}
