package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hambosto/go-downloader/internal/domain"
	"github.com/hambosto/go-downloader/internal/infrastructure"
	"github.com/hambosto/go-downloader/internal/ui"
	"github.com/hambosto/go-downloader/internal/usecase"
	"github.com/hambosto/go-downloader/internal/util"
)

func main() {
	urlFlag := flag.String("url", "", "URL of the file to download")
	output := flag.String("output", "", "Output file name (optional)")
	workers := flag.Int("workers", 5, "Number of concurrent workers")
	chunkSize := flag.Int64("chunk-size", 1024*1024, "Size of each chunk in bytes")
	flag.Parse()

	if *urlFlag == "" {
		flag.Usage()
		return
	}

	// Generate output file name if not provided
	outputFile := *output
	if len(outputFile) == 0 {
		outputFile = util.GenerateOutputFile(*urlFlag)
	}

	job := &domain.Job{
		URL:        *urlFlag,
		OutputFile: outputFile,
		NumWorkers: *workers,
		ChunkSize:  *chunkSize,
	}

	fileManager := infrastructure.NewFileManager()
	httpClient := infrastructure.NewHTTPClient(30*time.Second, 3, 3*time.Second)
	downloadUseCase := usecase.NewDownloadUseCase(job, fileManager, httpClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Shutting down...")
		cancel()
	}()

	cli := ui.NewCLI(downloadUseCase)
	if err := cli.Run(ctx); err != nil {
		log.Fatalf("CLI error: %v", err)
	}
}
