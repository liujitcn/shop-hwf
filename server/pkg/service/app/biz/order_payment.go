package biz

import (
	"gitee.com/liujit/shop/server/lib/data"
)

type OrderPaymentCase struct {
	data.OrderPaymentRepo
}

// NewOrderPaymentCase new a OrderPayment use case.
func NewOrderPaymentCase(orderPaymentRepo data.OrderPaymentRepo,
) *OrderPaymentCase {
	return &OrderPaymentCase{
		OrderPaymentRepo: orderPaymentRepo,
	}
}
