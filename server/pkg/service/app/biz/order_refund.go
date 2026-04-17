package biz

import (
	"gitee.com/liujit/shop/server/lib/data"
)

type OrderRefundCase struct {
	data.OrderRefundRepo
}

// NewOrderRefundCase new a OrderRefund use case.
func NewOrderRefundCase(orderRefundRepo data.OrderRefundRepo,
) *OrderRefundCase {
	return &OrderRefundCase{
		OrderRefundRepo: orderRefundRepo,
	}
}
