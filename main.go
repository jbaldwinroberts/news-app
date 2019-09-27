package main

import (
	"io/ioutil"
	"os"
	"time"

	_ "github.com/astaxie/beego/config/yaml"
	"github.com/josephroberts/esqimo/reader"
	"github.com/josephroberts/esqimo/server"
	"github.com/josephroberts/esqimo/store"
	Swagger "github.com/josephroberts/esqimo/swagger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
)

var defaultFeeds = []string{
	"http://feeds.bbci.co.uk/news/uk/rss.xml",
}

const defaultPort = ":1323"

type config struct {
	Port  string   `yaml:"port"`
	Feeds []string `yaml:"feeds"`
}

func main() {
	// Load config
	config := loadConfig("config.yaml")
	log.Infof("Feeds: %v\n", config.Feeds)
	log.Infof("Port: %v\n", config.Port)

	// Reader instance
	reader := reader.RSS{
		URLS: config.Feeds,
	}

	// Store instance
	store := store.New(reader, 10*time.Minute)

	// Server instance
	server := &server.Server{
		Store: store,
	}

	// Echo instance
	echo := echo.New()

	// Echo middleware
	echo.Use(middleware.Logger())
	echo.Use(middleware.Recover())

	// Register handlers
	Swagger.RegisterHandlers(echo, server)

	// Start server
	echo.Logger.Fatal(echo.Start(config.Port))
}

func loadConfig(filename string) *config {
	// Set default configs
	config := &config{
		Port:  defaultPort,
		Feeds: defaultFeeds,
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Errorf("unable to find config file: %v\n", err)
		return config
	}

	// Read the config file
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("unable to read config file: %v\n", err)
		return config
	}

	// Unmarshall the config file
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Errorf("unable to unmarshal config file: %v\n", err)
		return config
	}

	return config
}
