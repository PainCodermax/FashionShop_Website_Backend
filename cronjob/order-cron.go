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

	// Lặp lại việc đọc dữ liệu cho đến khi không còn đơn hàng nào
	for {
		opt := options.FindOptions{
			Limit: utils.ParseIn64ToPointer(limit),
			Skip:  utils.ParseIn64ToPointer(offset),
		}

		// Tìm tất cả các đơn hàng với giới hạn và lùi hiện tại
		rs, err := controllers.OrderCollection.Find(ctx, bson.M{}, &opt)
		if err != nil {
			log.Println("Error finding orders:", err)
			return
		}

		// Lặp qua từng đơn hàng và cập nhật trạng thái
		for rs.Next(ctx) {
			order := models.Order{}
			if err := rs.Decode(&order); err != nil {
				log.Println("Error decoding order:", err)
				continue
			}

			go controllers.UpdateOrderStatus(order)

			log.Printf("Updated status for order %s\n", order.ID)
		}

		// Kiểm tra xem có còn đơn hàng nào nữa hay không
		if err := rs.Err(); err != nil {
			log.Println("Error iterating orders:", err)
			return
		}

		// Nếu không còn đơn hàng nào, kết thúc vòng lặp
		if !rs.Next(ctx) {
			break
		}

		// Di chuyển tới trang tiếp theo
		offset += limit
	}
}
