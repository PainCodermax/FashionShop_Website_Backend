package worker

import (
	"context"

	"github.com/PainCodermax/FashionShop_Website_Backend/controllers"
)

func Worker(ctx context.Context, orderChannel <-chan string) {
	for {
		select {
		case orderId := <-orderChannel:
			// Perform the update operation
			controllers.UpdateOrder(ctx, orderId)
		case <-ctx.Done():
			return
		}
	}
}
