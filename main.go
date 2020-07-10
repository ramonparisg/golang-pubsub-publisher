package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	slog "github.com/go-eden/slf4go"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

func main() {
	err := publishMessage(os.Getenv("PROJECT_ID"), os.Getenv("TOPIC_ID"), 2)
	if err != nil {
		slog.Errorf("Error %+v", err)
	}
}

func publishMessage(projectID, topicID string, n int) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	var wg sync.WaitGroup
	var totalErrors uint64
	t := client.Topic(topicID)

	for i := 0; i < n; i++ {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("Message " + strconv.Itoa(i)),
			Attributes: map[string]string{
				"origin":   "golang",
				"username": "gcp",
			},
		})

		wg.Add(1)
		go func(i string, res *pubsub.PublishResult) {
			defer wg.Done()
			// The Get method blocks until a server-generated ID or
			// an error is returned for the published message.
			id, err := res.Get(ctx)
			if err != nil {
				// Error handling code can be added here.
				slog.Errorf("Failed to publish: %v", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
			slog.Infof("Published message %s; msg ID: %v", i, id)
		}("string", result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("%d of %d messages did not publish successfully", totalErrors, n)
	}
	return nil
}
