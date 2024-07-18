package controllers

import (
	"context"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/enum"
	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/PainCodermax/FashionShop_Website_Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type filter struct {
	From time.Time `json:"from,ommitempty"`
	To   time.Time `json:"to,ommitempty"`
}

func GetReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				defer cancel()

				var fil filter
				if err := c.BindJSON(&fil); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				filterCondition := bson.M{}
				filterOrderSuccess := bson.M{"status": enum.Received}
				// Nếu có ngày From và To, cập nhật điều kiện lọc
				if !fil.From.IsZero() && !fil.To.IsZero() {
					filterCondition = bson.M{
						"created_at": bson.M{
							"$gte": fil.From,
							"$lte": fil.To,
						},
					}
					filterOrderSuccess["created_at"] = bson.M{
						"$gte": fil.From,
						"$lte": fil.To,
					}

				}

				var totalAmount int = 0
				totalUser, _ := UserCollection.CountDocuments(ctx, filterCondition)
				totalOrder, _ := OrderCollection.CountDocuments(ctx, filterCondition)
				totalProduct, _ := ProductCollection.CountDocuments(ctx, filterCondition)
				totalOrderSuccess, _ := OrderCollection.CountDocuments(ctx, filterOrderSuccess)
				totalRating, _ := RatingCollection.CountDocuments(ctx, filterCondition)

				rs, err := OrderCollection.Find(ctx, filterOrderSuccess, nil)
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
					totalAmount = totalAmount + order.Price
				}

				limit := int64(20)
				offset := int64(0)
				months := make(map[string]int64)
				for i := 1; i <= 12; i++ {
					month := ""
					switch i {
					case 1:
						month = "January"
					case 2:
						month = "February"
					case 3:
						month = "March"
					case 4:
						month = "April"
					case 5:
						month = "May"
					case 6:
						month = "June"
					case 7:
						month = "July"
					case 8:
						month = "August"
					case 9:
						month = "September"
					case 10:
						month = "October"
					case 11:
						month = "November"
					case 12:
						month = "December"
					}
					months[month] = 0
				}
				for {
					opt := options.FindOptions{
						Limit: utils.ParseIn64ToPointer(limit),
						Skip:  utils.ParseIn64ToPointer(offset),
					}

					rs, err := OrderCollection.Find(ctx, bson.D{{"status", enum.Received}}, &opt)
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
						month := order.Created_At.Month().String()
						if amount, ok := months[month]; ok {
							months[month] += int64(order.Price) + amount
						} else {
							months[month] = int64(order.Price)
						}
						log.Printf("e %s\n", order.ID)
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

				var rp models.Report

				rp.TotalOrder = totalOrder
				rp.TotalOrderSuccess = totalOrderSuccess
				rp.TotalProduct = totalProduct
				rp.TotalUser = totalUser
				rp.TotalRating = totalRating
				rp.TotalAmount = totalAmount

				var amountList []models.Amount
				monthsSlice := []string{
					"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December",
				}

				monthIndex := make(map[string]int)
				for i, month := range monthsSlice {
					monthIndex[month] = i + 1 // Bắt đầu từ 1 thay vì 0
				}

				sort.Slice(monthsSlice, func(i, j int) bool {
					return monthIndex[monthsSlice[i]] < monthIndex[monthsSlice[j]]
				})
				for _, month := range monthsSlice {
					value := months[month]
					amount := models.Amount{
						TotalAmount: value,
						Month:       month,
					}
					amountList = append(amountList, amount)
				}
				rp.Amounts = amountList
				c.JSON(http.StatusOK,
					gin.H{
						"message": "get report success",
						"data":    rp,
					},
				)
			}
			return
		} else {
			c.JSON(http.StatusForbidden,
				gin.H{
					"message": "you don't have permission",
				},
			)
		}
	}
}
