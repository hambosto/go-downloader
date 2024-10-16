package usecase

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/hambosto/go-downloader/internal/domain"
)

type DownloadUseCase interface {
	Start(ctx context.Context) error
	GetProgress() (int64, int64, bool)
	GetJobOutputFile() string
	GetJobURL() string
}

type downloadUseCase struct {
	job         *domain.Job
	fileManager FileManager
	httpClient  HTTPClient
	mu          sync.Mutex
}

type FileManager interface {
	Create(string) error
	Write(string, int64, []byte) error
	GetSize(string) (int64, error)
}

type HTTPClient interface {
	GetFileSize(string) (int64, error)
	DownloadChunk(string, int64, int64) (io.ReadCloser, error)
}

func NewDownloadUseCase(job *domain.Job, fileManager FileManager, httpClient HTTPClient) DownloadUseCase {
	return &downloadUseCase{
		job:         job,
		fileManager: fileManager,
		httpClient:  httpClient,
	}
}

func (uc *downloadUseCase) Start(ctx context.Context) error {
	var err error
	uc.job.FileSize, err = uc.httpClient.GetFileSize(uc.job.URL)
	if err != nil {
		return fmt.Errorf("error getting file size: %v", err)
	}

	err = uc.fileManager.Create(uc.job.OutputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}

	chunks := uc.createChunks()
	chunksChan := make(chan domain.Chunk, len(chunks))
	for _, c := range chunks {
		chunksChan <- c
	}
	close(chunksChan)

	var wg sync.WaitGroup
	for i := 0; i < uc.job.NumWorkers; i++ {
		wg.Add(1)
		go uc.worker(ctx, chunksChan, &wg)
	}

	wg.Wait()

	uc.mu.Lock()
	uc.job.IsCompleted = true
	uc.mu.Unlock()

	return nil
}

func (uc *downloadUseCase) worker(ctx context.Context, chunks <-chan domain.Chunk, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case chunk, ok := <-chunks:
			if !ok {
				return
			}
			if err := uc.downloadChunk(ctx, chunk); err != nil {
				fmt.Printf("Error downloading chunk: %v\n", err)
			}
		}
	}
}

func (uc *downloadUseCase) downloadChunk(ctx context.Context, chunk domain.Chunk) error {
	resp, err := uc.httpClient.DownloadChunk(chunk.URL, chunk.Offset, chunk.Size)
	if err != nil {
		return err
	}
	defer resp.Close()

	buffer := make([]byte, 32*1024)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := resp.Read(buffer)
			if n > 0 {
				if err = uc.fileManager.Write(uc.job.OutputFile, chunk.Offset, buffer[:n]); err != nil {
					return err
				}
				chunk.Offset += int64(n)
				chunk.Size -= int64(n)

				uc.mu.Lock()
				uc.job.Downloaded += int64(n)
				uc.mu.Unlock()
			}
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}

func (uc *downloadUseCase) createChunks() []domain.Chunk {
	var chunks []domain.Chunk
	for offset := int64(0); offset < uc.job.FileSize; offset += uc.job.ChunkSize {
		size := uc.job.ChunkSize
		if offset+size > uc.job.FileSize {
			size = uc.job.FileSize - offset
		}
		chunks = append(chunks, domain.Chunk{URL: uc.job.URL, Offset: offset, Size: size})
	}
	return chunks
}

func (uc *downloadUseCase) GetProgress() (int64, int64, bool) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	return uc.job.Downloaded, uc.job.FileSize, uc.job.IsCompleted
}

func (uc *downloadUseCase) GetJobOutputFile() string {
	return uc.job.OutputFile
}

func (uc *downloadUseCase) GetJobURL() string {
	return uc.job.URL
}
