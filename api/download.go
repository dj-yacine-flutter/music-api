package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/dj-yacine-flutter/music-api/constant"
	"github.com/gofiber/fiber/v2"
)

func YouTubeVideoDownload(c *fiber.Ctx) error {
	// Get the video URL from the query parameter
	videoURL := c.Query("url")

	// Check if the 'url' parameter is missing
	if videoURL == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing 'URL' parameter",
		})
	}

	decodedVideoURL, err := url.QueryUnescape(videoURL)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to parse video url",
			"details": err.Error(),
		})
	}
	// Get the video URL from the query parameter
	videoFormatID := c.Query("id")

	// Check if the 'url' parameter is missing
	if videoFormatID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing 'ID' parameter",
		})
	}

	titleCmd := exec.Command("yt-dlp", "--skip-download", "--get-filename", "-o", "%(title)s.%(ext)s", "-f", videoFormatID, decodedVideoURL)
	titleCmd.Stderr = os.Stderr

	titleOutput, err := titleCmd.Output()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get video title",
			"details": err.Error(),
		})
	}

	fileName := strings.TrimSpace(string(titleOutput))

	// Run yt-dlp to download the video as an MP3 file
	downloadCmd := exec.Command("yt-dlp", "-f", videoFormatID, "--output", fmt.Sprintf("%s/%s", constant.MusicPath, fileName), decodedVideoURL)
	downloadCmd.Stderr = os.Stderr

	if err := downloadCmd.Run(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to download video",
			"details": err.Error(),
		})
	}

	// Construct the full path to the downloaded file
	filePath := fmt.Sprintf("%s/%s", constant.MusicPath, fileName)

	// Set the response headers for file download
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Set("Content-Type", "application/octet-stream")

	// Send the file as a response
	if err := c.SendFile(filePath); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to send file contents",
			"details": err.Error(),
		})
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete file",
			"details": err.Error(),
		})
	}

	return nil
}

/* 		// Convert the downloaded file to mp3 using FFmpeg
convertCmd := exec.Command("ffmpeg", "-i", fileName, "-c:a", "libmp3lame", "-q:a", "0", newFileName, "-y")
convertCmd.Stderr = os.Stderr

// Execute the command to convert to mp3
if err := convertCmd.Run(); err != nil {
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"error":   "Failed to convert video to mp3",
		"details": err.Error(),
	})
} */
