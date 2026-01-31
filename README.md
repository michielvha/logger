# Logger

A simple, structured logging wrapper around Go's `log/slog` package.

## Features

- JSON structured logging output
- Configurable log levels via code or `LOG_LEVEL` environment variable
- Printf-style formatting functions
- Compatible with standard `log` package interface
- Zero external dependencies (standard library only)

## Installation

```bash
go get github.com/michielvha/logger
```

## Usage

```go
package main

import "github.com/michielvha/logger"

func main() {
    // Optional: Initialize with a specific level
    // If not called, defaults to INFO (or LOG_LEVEL env var)
    logger.Init("debug")

    // Simple logging
    logger.Info("Application started")
    logger.Debug("Debug message", "key", "value")
    logger.Warn("Warning message")
    logger.Error("Error occurred", "error", err)

    // Printf-style logging
    logger.Infof("Processing %d items", count)
    logger.Errorf("Failed to connect: %v", err)

    // Fatal exits with status 1
    logger.Fatal("Critical error")
}
```

## Log Levels

Set via `logger.Init()` or `LOG_LEVEL` environment variable:

- `debug` - Detailed debugging information
- `info` - General operational messages (default)
- `warn` / `warning` - Warning messages
- `error` - Error messages

## License

Apache 2.0
