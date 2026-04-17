package biz

import (
	"gitee.com/liujit/shop/server/lib/data"
)

type OrderCancelCase struct {
	data.OrderCancelRepo
}

// NewOrderCancelCase new a OrderCancel use case.
func NewOrderCancelCase(orderCancelRepo data.OrderCancelRepo,
) *OrderCancelCase {
	return &OrderCancelCase{
		OrderCancelRepo: orderCancelRepo,
	}
}
