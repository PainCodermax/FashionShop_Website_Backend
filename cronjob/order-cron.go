package cronjob

import (
	"context"
	"log"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/controllers"
	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/PainCodermax/FashionShop_Website_Backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateOrderStatusJob() {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	limit := int64(20)
	offset := int64(0)

	for {
		opt := options.FindOptions{
			Limit: utils.ParseIn64ToPointer(limit),
			Skip:  utils.ParseIn64ToPointer(offset),
		}

		rs, err := controllers.OrderCollection.Find(ctx, bson.M{}, &opt)
		if err != nil {
			log.Println("Error finding orders:", err)
			return
		}

		for rs.Next(ctx) {
			order := models.Order{}
			if err := rs.Decode(&order); err != nil {
				log.Println("Error decoding order:", err)
				continue
			}

			go controllers.UpdateOrderStatus(order)

			log.Printf("Updated status for order %s\n", order.ID)
		}

		if err := rs.Err(); err != nil {
			log.Println("Error iterating orders:", err)
			return
		}

		if !rs.Next(ctx) {
			break
		}

		offset += limit
	}
}
