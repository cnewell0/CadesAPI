package main

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"time"

	"cloud.google.com/go/pubsub"
)

//PullMsgs pulls sum messages
// func PullMsgs(w io.Writer, projectID, subID string) error {
// 	projectID = "kochava-testing"
// 	subID = "test-sub"
// 	ctx := context.Background()
// 	client, err := pubsub.NewClient(ctx, projectID)
// 	if err != nil {
// 		return fmt.Errorf("pubsub.NewClient: %v", err)
// 	}

// 	// Consume 10 messages.
// 	var mu sync.Mutex
// 	received := 0
// 	sub := client.Subscription(subID)
// 	cctx, cancel := context.WithCancel(ctx)
// 	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
// 		mu.Lock()
// 		defer mu.Unlock()
// 		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
// 		msg.Ack()
// 		received++
// 		if received == 10 {
// 			cancel()
// 		}
// 	})
// 	if err != nil {
// 		return fmt.Errorf("Receive: %v", err)
// 	}
// 	return nil
// }

func pullMsgsConcurrenyControl(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)
	// Must set ReceiveSettings.Synchronous to false (or leave as default) to enable
	// concurrency settings. Otherwise, NumGoroutines will be set to 1.
	sub.ReceiveSettings.Synchronous = false
	// NumGoroutines is the number of goroutines sub.Receive will spawn to pull messages concurrently.
	sub.ReceiveSettings.NumGoroutines = runtime.NumCPU()

	// Receive messages for 10 seconds.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Create a channel to handle messages to as they come in.
	cm := make(chan *pubsub.Message)
	defer close(cm)
	// Handle individual messages in a goroutine.
	go func() {
		for msg := range cm {
			fmt.Fprintf(w, "Got message :%q\n", string(msg.Data))
			msg.Ack()
		}
	}()

	// Receive blocks until the context is cancelled or an error occurs.
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}

	return nil
}
