package enum

type OrderStatus string

const (
	Pending    OrderStatus = "WAITING"
	Processing OrderStatus = "PROCESSING"
	Shipping   OrderStatus = "SHIPPING"
	Delivered  OrderStatus = "Delivered"
	Cancelled  OrderStatus = "CANCELLED"
)


