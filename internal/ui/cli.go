package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hambosto/go-downloader/internal/usecase"
	"github.com/hambosto/go-downloader/internal/util"
	"github.com/mbndr/figlet4go"
)

type CLI struct {
	downloadUseCase usecase.DownloadUseCase
}

func NewCLI(downloadUseCase usecase.DownloadUseCase) *CLI {
	return &CLI{
		downloadUseCase: downloadUseCase,
	}
}

func (cli *CLI) Run(ctx context.Context) error {
	fileName := cli.downloadUseCase.GetJobOutputFile()
	url := cli.downloadUseCase.GetJobURL()

	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorGreen,
		figlet4go.ColorYellow,
		figlet4go.ColorCyan,
	}

	title, _ := ascii.RenderOpts("Downloader", options)
	fmt.Print(title)

	fmt.Printf("File Name:    %s\n", fileName)
	fmt.Printf("Download URL: %s\n", url)
	fmt.Println()

	errChan := make(chan error, 1)
	go func() {
		errChan <- cli.downloadUseCase.Start(ctx)
	}()

	progressTicker := time.NewTicker(500 * time.Millisecond)
	defer progressTicker.Stop()

	var lastBytes int64
	startTime := time.Now()
	lastUpdate := startTime

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			if err != nil {
				return fmt.Errorf("download error: %v", err)
			}
			fmt.Println("Download completed successfully!")
			return nil
		case <-progressTicker.C:
			currentBytes, totalBytes, isCompleted := cli.downloadUseCase.GetProgress()
			currentTime := time.Now()
			elapsed := currentTime.Sub(lastUpdate)
			totalElapsed := currentTime.Sub(startTime)

			speed := float64(currentBytes-lastBytes) / elapsed.Seconds()
			avgSpeed := float64(currentBytes) / totalElapsed.Seconds()
			percent := (float64(currentBytes) / float64(totalBytes)) * 100

			remaining := totalBytes - currentBytes
			eta := time.Duration(float64(remaining) / avgSpeed * float64(time.Second))
			fmt.Printf("\033[2K\r")
			fmt.Printf("Progress: %.2f%% | %s / %s\n", percent, humanize.Bytes(uint64(currentBytes)), humanize.Bytes(uint64(totalBytes)))
			fmt.Printf("Speed: %s/s (Avg: %s/s)\n", humanize.Bytes(uint64(speed)), humanize.Bytes(uint64(avgSpeed)))
			fmt.Printf("Elapsed: %s | Remaining: %s\n", util.FormatDuration(totalElapsed), util.FormatDuration(eta))
			fmt.Printf("\033[3A")

			lastBytes = currentBytes
			lastUpdate = currentTime

			if isCompleted {
				return nil
			}
		}
	}
}
