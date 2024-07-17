package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var FlashSaleCache *cache.Cache
var WishListCache *cache.Cache

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

func InitWishListCache(productId, userId string) {
	if WishListCache == nil {
		c := cache.New(cache.NoExpiration, cache.NoExpiration)
		c.Set(productId+userId, "ok", cache.NoExpiration)
		WishListCache = c
	} else {
		WishListCache.Set(productId+userId, "ok", cache.NoExpiration)
	}
}

func GetWishListCache(productId, userId string) bool {
	if WishListCache == nil {
		return false
	}
	if _, ok := WishListCache.Get(productId + userId); ok {
		return true
	} else {
		return false
	}
}

func DeleteWishListCache(productId, userId string) {
	if WishListCache == nil {
		return
	}
	WishListCache.Delete(productId + userId)
}
