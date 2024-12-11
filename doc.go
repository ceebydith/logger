/*
Package logger provides versatile logging solutions for Go applications.
It includes global logger instances, type-specific log handlers, and various
logging utilities such as file writer, broadcast writer, and tail writer.

# Installation

To install the package, run:

	go get github.com/ceebydith/logger

# Usage

# Basic Logger

Create a new logger instance and use it for logging messages:

	package main

	import (
	    "github.com/ceebydith/logger"
	)

	func main() {
	    log := logger.New()
	    log.Print("This is a log message.")
	    log.Printf("This is a formatted log message with value: %d", 42)
	}

# LogHandler Interface

Use the LogHandler interface to integrate logging into your types:

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

# Deferred Logging

Use deferred logging to capture the completion status of operations:

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

# Advanced Usage

# File Writer

Log messages to a file with dynamic filename formatting:

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

# Broadcast Writer

Broadcast log messages to multiple listeners:

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

# Tail Writer

Keep track of the last n lines of log messages:

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

# Prefix Handler

Add prefixes to log messages:

	package main

	import (
	    "github.com/ceebydith/logger"
	)

	func main() {
	    log := logger.New()
	    prefixLogger := logger.PrefixHandler("PREFIX", log)
	    prefixLogger.Log("This is a prefixed log message.")
	}
*/
package logger
