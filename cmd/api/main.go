package main

import (
	"log"

	"github.com/0xBoji/web3-edu-core/config"
	"github.com/0xBoji/web3-edu-core/internal/api"
	"github.com/0xBoji/web3-edu-core/internal/database/postgres"
	"github.com/0xBoji/web3-edu-core/internal/database/redis"
)

func init() {
	// Load configuration
	config.Setup()

	// Setup database
	postgres.Setup()

	// Setup Redis
	redis.Setup()
}

func main() {
	log.Printf("Starting %s in %s mode", config.AppSetting.Name, config.ServerSetting.RunMode)

	// Create and run server
	server := api.NewServer()
	server.Run()
}
