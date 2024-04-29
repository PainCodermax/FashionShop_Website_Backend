package worker

import (
	"context"
)

func Worker(ctx context.Context, orderChannel <-chan string) {
	// time.Sleep(10 * time.Second)
	for {
		select {
		// case orderId := <-orderChannel:
		// Perform the update operation
		// go controllers.UpdateOrder(ctx, orderId)
		case <-ctx.Done():
			return
		}
	}
}
