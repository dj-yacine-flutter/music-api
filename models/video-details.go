package models


type VideoDetails struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Sound      bool   `json:"sound"`
	Bitrate    string `json:"bitrate"`
	Resolution string `json:"resolution"`
	Quality    string `json:"quality"`
}