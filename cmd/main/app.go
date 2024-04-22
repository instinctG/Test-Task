package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	config2 "github.com/instinctG/Test-task/internal/config"
	"github.com/instinctG/Test-task/internal/db"
	"github.com/instinctG/Test-task/internal/handler"
	"github.com/instinctG/Test-task/internal/token"
	"log"
)

func Run() error {
	gin.SetMode(gin.ReleaseMode)
	config, err := config2.LoadConfig()
	if err != nil {
		log.Println(err)
		return err
	}
	//"C:/Users/azizs/Documents/Test-task/config"

	database, err := db.NewClient(context.Background(), config.MongoURL)
	if err != nil {
		fmt.Errorf("couldn't connect to database :%w", err)
	}

	tokenService := token.NewTokenService(database)

	ginHandler := handler.NewHandler(config.Port, tokenService)
	if err := ginHandler.Serve(config.Port); err != nil {
		return err
	}
	return nil

}

func main() {
	fmt.Println("PROJECT IS GOING...")
	if err := Run(); err != nil {
		fmt.Println(err)
	}

}
