package biz

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/api/admin"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/utils/timeutil"
	"gorm.io/gorm"
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

func (c *OrderCancelCase) GetFromByOrderId(ctx context.Context, orderId int64) (*admin.OrderCancel, error) {
	orderCancel, err := c.Find(ctx, &data.OrderCancelCondition{
		OrderId: orderId,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &admin.OrderCancel{}, nil
		}
		return nil, err
	}
	return &admin.OrderCancel{
		Reason:    common.OrderCancelReason(orderCancel.Reason),
		CreatedAt: timeutil.TimeToTimeString(orderCancel.CreatedAt),
	}, nil
}
