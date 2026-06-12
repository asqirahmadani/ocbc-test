package main

import (
	"boilerplate/config"
	"boilerplate/repo/postgres"
	"boilerplate/router"
	"boilerplate/server"
	"boilerplate/utils/logger"
	"log"
)

func main() {
	if err := logger.InitZapSugaredLogger(); err != nil {
		log.Fatal(err)
	}

	cfg := config.Load()
	db, err := postgres.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	defer func() { _ = logger.Log.Sync() }()

	router := router.SetupRouter(*cfg)
	server.StartHTTPServer(router, cfg.App)
}