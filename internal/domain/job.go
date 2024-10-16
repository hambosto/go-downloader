package domain

type Job struct {
	URL         string
	OutputFile  string
	NumWorkers  int
	ChunkSize   int64
	FileSize    int64
	Downloaded  int64
	IsCompleted bool
}
