package models

const (
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"

	StockIn     = "IN"
	StockOut    = "OUT"
	StockAdjust = "ADJUST"

	MediaTypeImage = "IMAGE"
	MediaTypeVideo = "VIDEO"

	CacheKeyProductDetail = "product:detail:"
	CacheKeyProductList   = "products:list:"
	CacheKeyCategoryAll   = "categories:all"
	CacheKeyBrandAll      = "brands:all"
)
