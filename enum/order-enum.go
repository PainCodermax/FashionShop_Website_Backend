package enum

type OrderStatus string

const (
	Pending    OrderStatus = "WAITING"
	Submitted  OrderStatus = "SUBMITED"
	Processing OrderStatus = "PROCESSING"
	Shipping   OrderStatus = "SHIPPING"
	Delivery   OrderStatus = "DELIVERY"
	Cancelled  OrderStatus = "CANCELLED"
	Received  OrderStatus = "RECEIVED"
)
