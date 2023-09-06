package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/dj-yacine-flutter/music-api/models"
	"github.com/gofiber/fiber/v2"
)

// FilterFormats removes formats with "webp" and m3u8 from the list.
func filterFormats(formats []string) []string {
	var filteredFormats []string

	for _, format := range formats {
		// Check the format
		if (strings.Contains(format, "m4a") || strings.Contains(format, "mp4")) && !strings.Contains(format, "3gp") {
			if strings.Contains(format, "http") {
				filteredFormats = append(filteredFormats, format)
			}
		}
	}

	return filteredFormats
}

func parseVideoFormat(formatStr string) (models.VideoDetails, error) {
	// Define regular expressions to extract relevant information
	idRegex := regexp.MustCompile(`(\d+)`)
	bitrateRegex := regexp.MustCompile(`(\d+k)`)
	resolutionRegex := regexp.MustCompile(`(\d+x\d+)`)
	qualityRegex := regexp.MustCompile(`(\d+)x(\d+)`)

	// Find matches in the format string
	idMatches := idRegex.FindStringSubmatch(formatStr)
	bitrateMatches := bitrateRegex.FindStringSubmatch(formatStr)
	resolutionMatches := resolutionRegex.FindStringSubmatch(formatStr)
	qualityMatches := qualityRegex.FindStringSubmatch(formatStr)

	// Check for required matches
	if len(idMatches) == 0 || len(bitrateMatches) == 0 {
		return models.VideoDetails{}, fmt.Errorf("invalid format string: %s", formatStr)
	}

	// Extract information
	id := idMatches[1]
	videoType := "video"
	if strings.Contains(formatStr, "audio only") {
		videoType = "audio"
	}

	sound := !strings.Contains(formatStr, "M.")
	bitrate := bitrateMatches[1]
	var resolution, quality string
	if len(resolutionMatches) > 1 {
		resolution = resolutionMatches[1]
	}
	if len(qualityMatches) >= 3 {
		quality = qualityMatches[2] + "p" // Add 'p' to make it like "720p"
		// Now, 'quality' contains "720p"
		// You can use 'quality' in your data structure or output
	}

	return models.VideoDetails{
		ID:         id,
		Type:       videoType,
		Sound:      sound,
		Bitrate:    bitrate,
		Resolution: resolution,
		Quality:    quality,
	}, nil
}

func YouTubeVideoFormats(c *fiber.Ctx) error {
	// Get the video URL from the query parameter
	videoURL := c.Query("url")

	// Check if the 'url' parameter is missing
	if videoURL == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing 'url' parameter",
		})
	}

	decodedVideoURL, err := url.QueryUnescape(videoURL)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to parse video url",
			"details": err.Error(),
		})
	}

	// Run yt-dlp to list available formats for the video
	cmd := exec.Command("yt-dlp", "--list-formats", decodedVideoURL)
	cmd.Stderr = os.Stderr

	// Capture the console output
	output, err := cmd.Output()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve video formats",
			"details": err.Error(),
		})
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")

	// Filter the available formats
	formats := filterFormats(lines)

	// Parse the formats into the desired JSON format
	var parsedFormats []models.VideoDetails
	for _, formatStr := range formats {
		parsedFormat, err := parseVideoFormat(formatStr)
		if err == nil {
			parsedFormats = append(parsedFormats, parsedFormat)
		}
	}

	// Present the parsed formats to the user
	return c.Status(http.StatusOK).JSON(parsedFormats)
}
