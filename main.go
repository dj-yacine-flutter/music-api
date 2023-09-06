package main

import (
	"log"

	"github.com/dj-yacine-flutter/music-api/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/details", api.YouTubeVideoDetails)
	app.Get("/format", api.YouTubeVideoFormats)
	app.Get("/download", api.YouTubeVideoDownload)

	port := ":8080"
	log.Fatal(app.Listen(port))
}
