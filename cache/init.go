package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var FlashSaleCache *cache.Cache

func InitFlashSaleCahe(t time.Duration, productId string, price int) {
	if FlashSaleCache == nil {
		c := cache.New(cache.NoExpiration, 720*time.Minute)
		c.Set(productId, price, t)
		FlashSaleCache = c
	} else {
		FlashSaleCache.Set(productId, price, t)
	}

}

func GetSalePriceByProductId(key string) int {
	if FlashSaleCache == nil {
		return 0
	}
	if price, ok := FlashSaleCache.Get(key); ok {
		return price.(int)
	} else {
		return 0
	}
}
