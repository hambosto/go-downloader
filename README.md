# go-downloader

go-downloader is a concurrent file downloader written in Go. It allows you to download files from the internet using multiple workers, potentially speeding up the download process for large files.

## Features

- Concurrent downloading using multiple workers
- Customizable chunk size for downloads
- Ability to specify output file name
- Graceful shutdown on interrupt signals

## Installation

### From Binary Release

1. Go to the [Releases](https://github.com/hambosto/go-downloader/releases) page of the go-downloader repository.
2. Download the latest release for your operating system and architecture:

   - Windows: Choose `go-downloader-<version>-windows-<arch>.exe`
   - macOS: Choose `go-downloader-<version>-darwin-<arch>`
   - Linux: Choose `go-downloader-<version>-linux-<arch>`

   Replace `<version>` with the latest version number and `<arch>` with your system architecture (amd64 for 64-bit, 386 for 32-bit, or arm64 for ARM-based systems).

3. Make the downloaded file executable (macOS and Linux only):

   ```bash
   chmod +x go-downloader-<version>-<os>-<arch>
   ```

4. Move the executable to a directory in your system's PATH. For example:
   - macOS/Linux:
     ```bash
     sudo mv go-downloader-<version>-<os>-<arch> /usr/local/bin/go-downloader
     ```
     ```bash
     chmod +x /usr/local/bin/go-downloader
     ```

Now you can run `go-downloader` from anywhere in your terminal.

### From Source

To install go-downloader from source, make sure you have Go installed on your system, then run:

```bash
go get github.com/hambosto/go-downloader/cmd/go-downloader@latest
```

## Usage

You can run go-downloader using the following command:

```bash
go-downloader -url <URL> [options]
```

### Options

- `-url`: URL of the file to download (required)
- `-output`: Output file name (optional)
- `-workers`: Number of concurrent workers (default: 5)
- `-chunk-size`: Size of each chunk in bytes (default: 1MB)

### Example

```bash
go-downloader -url https://example.com/largefile.zip -workers 10 -chunk-size 2097152 -output my_large_file.zip
```

This command will download the file from the specified URL using 10 workers, with a chunk size of 2MB, and save it as "my_large_file.zip".

## License

[MIT License](LICENSE)

## Acknowledgments

- This project uses the standard Go libraries for concurrent programming and HTTP requests.
- Thanks to all contributors who have helped to improve this project.
