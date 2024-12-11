# Logger Package

A versatile logging package for Go, designed to provide a simple yet powerful logging mechanism. This package offers global logger instances, type-specific log handlers, and a variety of logging utilities.

## Installation

To install the package, run:

```sh
go get github.com/ceebydith/logger
```

## Usage
### Basic Logger

Create a new logger instance and use it for logging messages:

```go
package main

import (
    "github.com/ceebydith/logger"
)

func main() {
    log := logger.New()
    log.Print("This is a log message.")
    log.Printf("This is a formatted log message with value: %d", 42)
}
```

### LogHandler Interface

Use the `LogHandler` interface to integrate logging into your types:

```go
package main

import (
    "github.com/ceebydith/logger"
)

type MyType struct {
    logger.LogHandler
}

func main() {
    myLogger := logger.New()
    myInstance := MyType{LogHandler: myLogger}
    myInstance.Log("Logging from MyType instance.")
}
```

### Deferred Logging

Use deferred logging to capture the completion status of operations:

```go
package main

import (
    "github.com/ceebydith/logger"
)

func main() {
    log := logger.New()
    err := funcWithError()
}

func funcWithError() (err error) {
    defer log.Defer(&err, "Starting function")()
    // Simulate some work
    return nil
}
```

## Advanced Usage

### File Writer

Log messages to a file with dynamic filename formatting:

```go
package main

import (
    "context"
    "github.com/ceebydith/logger"
)

func main() {
    ctx := context.Background()
    fileLogger := logger.FileWriter(ctx, "{curdir}/{yyyy}/{mm}/logfile_{yyyy}{mm}{dd}.log", 10)
    fileLogger.Write([]byte("Log message to file"))
}
```

### Broadcast Writer

Broadcast log messages to multiple listeners:

```go
package main

import (
    "context"
    "github.com/ceebydith/logger"
)

func main() {
    ctx := context.Background()
    broadcastLogger := logger.BroadcastWriter(ctx, 10)

    ch := broadcastLogger.Listen(ctx)
    go func() {
        for msg := range ch {
            fmt.Println("Received broadcast:", string(msg))
        }
    }()

    broadcastLogger.Write([]byte("Broadcast log message"))
}
```

### Tail Writer

Keep track of the last `n` lines of log messages:

```go
package main

import (
    "context"
    "fmt"
    "github.com/ceebydith/logger"
)

func main() {
    ctx := context.Background()
    tailLogger := logger.TailWriter(ctx, 100, 10)
    tailLogger.Write([]byte("Log message 1\n"))
    tailLogger.Write([]byte("Log message 2\n"))

    fmt.Println("Tail logs:\n", tailLogger.Tail())
}
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
This package is licensed under the MIT License. See the [`LICENSE`](https://github.com/ceebydith/logger/blob/main/LICENSE) file for details.