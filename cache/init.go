package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var FlashSaleCache *cache.Cache

func InitFlashSaleCahe(t time.Duration, productId string, price int) {
	c := cache.New(cache.NoExpiration, 720*time.Minute)
	c.Set(productId, price, t)
	FlashSaleCache = c
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
