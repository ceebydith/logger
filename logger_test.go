package logger_test

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ceebydith/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const BUFFER_SIZE = 10
const NUM_MESSAGES = 1000
const NUM_LISTENERS = 100
const MAX_LINES = 100

var random *rand.Rand

func uniqueMessage(message string) string {
	return fmt.Sprintf("#%d %s", random.Intn(10000), message)
}

func TestMain(m *testing.M) {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
	os.Exit(m.Run())
}

func TestBuffer(t *testing.T) {
	var messages []string
	var result strings.Builder

	ctx, cancel := context.WithCancel(context.TODO())
	buffer := logger.NewBuffer(ctx, BUFFER_SIZE, func(buffer []byte) {
		result.Write(buffer)
		time.Sleep(100 * time.Millisecond) // Simulate slow processing
	}, nil, nil)

	for i := 0; i < NUM_MESSAGES; i++ {
		message := uniqueMessage(fmt.Sprintf("Test message #%d\n", i))
		messages = append(messages, message)
		_, err := buffer.Write([]byte(message))
		require.NoError(t, err, "Write failed")
	}

	cancel()
	<-buffer.Done()

	for _, message := range messages {
		assert.Contains(t, result.String(), message)
	}
}

func TestBroadcastWriter(t *testing.T) {
	var wg sync.WaitGroup
	var listenerBuffers []*strings.Builder

	ctx, cancel := context.WithCancel(context.TODO())
	broadcast := logger.BroadcastWriter(ctx, BUFFER_SIZE)

	for i := 0; i < NUM_LISTENERS; i++ {
		buffer := new(strings.Builder)
		listenerBuffers = append(listenerBuffers, buffer)
		wg.Add(1)
		run := make(chan struct{}, 1)
		go func(buffer *strings.Builder, run chan struct{}) {
			defer wg.Done()
			stream := broadcast.Listen(context.TODO(), 2)
			close(run)
			for data := range stream {
				buffer.Write(data)
			}
		}(buffer, run)
		<-run // Make sure the listener goroutine is running before continuing
	}

	messages := []string{}
	for i := 0; i < NUM_MESSAGES; i++ {
		message := uniqueMessage(fmt.Sprintf("Test message %d\n", i))
		broadcast.Write([]byte(message))
		messages = append(messages, message)
	}

	cancel()
	wg.Wait()

	for _, listener := range listenerBuffers {
		for _, message := range messages {
			assert.Contains(t, listener.String(), message)
		}
	}
}

func TestTailWriter(t *testing.T) {
	var messages strings.Builder

	ctx, cancel := context.WithCancel(context.TODO())
	tail := logger.TailWriter(ctx, MAX_LINES, BUFFER_SIZE)

	for i := 0; i < NUM_MESSAGES; i++ {
		message := uniqueMessage(fmt.Sprintf("Test message #%d\n", i))
		messages.WriteString(message)
		_, err := tail.Write([]byte(message))
		require.NoError(t, err, "Write failed")
	}

	cancel()
	<-tail.Done()

	assert.Contains(t, messages.String(), tail.Tail())
}

func TestFileWriter(t *testing.T) {
	var messages strings.Builder
	ctx, cancel := context.WithCancel(context.TODO())
	filename := "target.txt"
	file := logger.FileWriter(ctx, filename, BUFFER_SIZE)

	for i := 0; i < NUM_MESSAGES; i++ {
		message := uniqueMessage(fmt.Sprintf("Test message #%d\n", i))
		messages.WriteString(message)
		_, err := file.Write([]byte(message))
		require.NoError(t, err, "Write failed")
	}

	cancel()
	<-file.Done()

	assert.FileExists(t, filename)

	f, err := os.Open(filename)
	require.NoError(t, err, "Failed to open file")
	result, err := io.ReadAll(f)
	require.NoError(t, err, "Failed to read file")
	f.Close()
	os.Remove(filename)
	require.Equal(t, messages.String(), string(result))
}

type User struct {
	logger.LogHandler
}

func (u *User) Login() (err error) {
	defer u.LogfDefer(&err, "%s Logging in", "John")()
	return nil
}
func (u *User) Logout() (err error) {
	defer u.LogDefer(&err, "Logging out")()
	return nil
}

func (u *User) DoSomething() {
	u.Log("Doing something")
	u.Logf("Doing something #%d", 1)
}

func TestLoggerType(t *testing.T) {
	user := User{
		LogHandler: logger.MustHandler(nil, logger.PrefixHandler("Test", logger.PrefixHandler("User", logger.New()))),
	}
	user.Login()
	user.Logout()
	user.DoSomething()
}

func TestLoggerVar(t *testing.T) {
	var messages strings.Builder
	ctx, cancel := context.WithCancel(context.TODO())
	buffer := logger.NewBuffer(ctx, BUFFER_SIZE, func(buffer []byte) {
		messages.Write(buffer)
	}, nil, nil)
	log := logger.New(buffer)
	log.Print("Test message")
	log.Printf("Test message #%d", 1)
	func() {
		defer log.Defer(nil, "Processing")()
	}()
	func() (err error) {
		defer log.Deferf(&err, "Thinking %s", "something")()
		return fmt.Errorf("Something went wrong")
	}()
	func() (err error) {
		defer log.Defer(&err, "Smiling")()
		return nil
	}()

	cancel()
	<-buffer.Done()
	t.Log(messages.String())
}
