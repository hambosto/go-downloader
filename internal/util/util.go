package util

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"
)

func GenerateOutputFile(urlString string) string {
	parsedURL, err := url.Parse(urlString)
	if err != nil || parsedURL.Path == "" {
		return "downloaded_file_" + time.Now().Format("20060102_150405")
	}

	fileName := path.Base(parsedURL.Path)
	if fileName == "." || fileName == "/" || fileName == "" || strings.Contains(fileName, "?") {
		return "downloaded_file_" + time.Now().Format("20060102_150405")
	}

	return fileName
}

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
