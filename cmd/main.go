package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/Vineet0197/shoppin-webcrawler/config"
	"github.com/Vineet0197/shoppin-webcrawler/internal/api"
	"github.com/Vineet0197/shoppin-webcrawler/internal/crawler"
	"github.com/Vineet0197/shoppin-webcrawler/internal/services"
	"github.com/Vineet0197/shoppin-webcrawler/pkg/utils"
	"github.com/labstack/echo"
)

var (
	cfgFile *string
)

func init() {
	cfgFile = flag.String("f", "../config/config.yaml", "Path to config file")
	flag.Parse()
}

func main() {
	defer handlePanic()

	// Load config
	err := config.LoadConfig(*cfgFile)
	if err != nil {
		panic(err)
	}

	// Initialize the Kafka storage
	ks, err := services.NewKafkaStorage(config.GetConfig().Kafka.Brokers, config.GetConfig().Kafka.Topic, config.GetConfig().Kafka.GroupID)
	if err != nil {
		panic(err)
	}

	// Start the server to take Domain and URL
	e := echo.New()

	e.POST("/crawl", func(ctx echo.Context) error {
		reqBody := []string{}
		// Bind the request body to the reqBody
		err := json.NewDecoder(ctx.Request().Body).Decode(&reqBody)
		if err != nil {
			return ctx.JSON(400, map[string]string{"error": "Invalid request body"})
		}

		return api.AcceptDomainAndURL(ctx, reqBody, ks)
	})

	// Start crawling the web pages
	c := crawler.NewCrawler(ks, utils.NewHTTPClient(60))
	go c.StartCrawling()

	log.Printf("Server started on port [:]%s", config.GetConfig().Server.Port)
	log.Fatal(e.Start(":" + config.GetConfig().Server.Port))
}

func handlePanic() {
	if r := recover(); r != nil {
		// Log the panic message
		log.Printf("Recovered from panic: %+v", r)
	}
}
