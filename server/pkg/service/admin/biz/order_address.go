package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/admin"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/utils/str"
)

type OrderAddressCase struct {
	data.OrderAddressRepo
}

// NewOrderAddressCase new a OrderAddress use case.
func NewOrderAddressCase(orderAddressRepo data.OrderAddressRepo,
) *OrderAddressCase {
	return &OrderAddressCase{
		OrderAddressRepo: orderAddressRepo,
	}
}

func (c *OrderAddressCase) GetFromByOrderId(ctx context.Context, orderId int64) (*admin.OrderAddress, error) {
	orderAddress, err := c.Find(ctx, &data.OrderAddressCondition{
		OrderId: orderId,
	})
	if err != nil {
		return nil, err
	}
	return &admin.OrderAddress{
		Receiver: orderAddress.Receiver,
		Contact:  orderAddress.Contact,
		Address:  str.ConvertJsonStringToStringArray(orderAddress.Address),
		Detail:   orderAddress.Detail,
	}, nil
}
