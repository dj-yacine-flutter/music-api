package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"

	"github.com/dj-yacine-flutter/music-api/models"
	"github.com/gofiber/fiber/v2"
)

func YouTubeVideoDetails(c *fiber.Ctx) error {
	// Get the video URL from the query parameter
	videoURL := c.Query("url")

	// Check if the 'url' parameter is missing
	if videoURL == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing 'url' parameter",
		})
	}
	// Encode the videoURL
	decodedVideoURL, err := url.QueryUnescape(videoURL)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to parse video url",
			"details": err.Error(),
		})
	}

	// Run yt-dlp with --dump-json to get video details
	cmd := exec.Command("yt-dlp", "--skip-download", "--get-id", decodedVideoURL)
	output, err := cmd.Output()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve video details",
			"details": err.Error(),
		})
	}

	IDs := getIDs(string(output))
	var videosInfo []models.VideoInfo
	for _, ID := range IDs {

		// Run yt-dlp with --dump-json to get video details
		cmd := exec.Command("yt-dlp", "--skip-download", "--dump-json", "--no-playlist", fmt.Sprintf("https://www.youtube.com/watch?v=%s", ID))
		output, err := cmd.Output()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve video details",
				"details": err.Error(),
			})
		}

		// Parse the JSON response into a VideoDetails struct
		var videoInfo models.VideoInfo
		err = json.Unmarshal(output, &videoInfo)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to parse video details",
				"details": err.Error(),
			})
		}
		videosInfo = append(videosInfo, videoInfo)
	}

	// Return the video details as JSON response
	return c.Status(http.StatusOK).JSON(videosInfo)
}

func getIDs(output string) []string {
	// Define a regular expression to match YouTube video IDs.
	idRegex := regexp.MustCompile(`[a-zA-Z0-9_-]{11}`)

	// Find all matches in the output string.
	matches := idRegex.FindAllString(output, -1)

	return matches
}
