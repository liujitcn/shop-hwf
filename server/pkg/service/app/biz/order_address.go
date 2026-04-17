package biz

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
)

type OrderAddressCase struct {
	data.OrderAddressRepo
	userAddressRepo data.UserAddressRepo
	baseAreaCase    *BaseAreaCase
}

// NewOrderAddressCase new a OrderAddress use case.
func NewOrderAddressCase(orderAddressRepo data.OrderAddressRepo,
	userAddressRepo data.UserAddressRepo,
	baseAreaCase *BaseAreaCase,
) *OrderAddressCase {
	return &OrderAddressCase{
		OrderAddressRepo: orderAddressRepo,
		userAddressRepo:  userAddressRepo,
		baseAreaCase:     baseAreaCase,
	}
}

func (c *OrderAddressCase) GetFromByOrderId(ctx context.Context, orderId int64) (*app.OrderResponse_Address, error) {
	orderAddress, err := c.Find(ctx, &data.OrderAddressCondition{
		OrderId: orderId,
	})
	if err != nil {
		return nil, err
	}
	return &app.OrderResponse_Address{
		Receiver: orderAddress.Receiver,
		Contact:  orderAddress.Contact,
		Address:  str.ConvertJsonStringToStringArray(orderAddress.Address),
		Detail:   orderAddress.Detail,
	}, nil
}

func (c *OrderAddressCase) CreateByAddressId(ctx context.Context, userId, orderId, addressId int64) error {
	userAddress, err := c.userAddressRepo.Find(ctx, userId, &data.UserAddressCondition{
		Id: addressId,
	})
	if err != nil {
		return errors.New("地址错误")
	}
	return c.Create(ctx, &models.OrderAddress{
		OrderID:  orderId,
		Receiver: userAddress.Receiver,
		Contact:  userAddress.Contact,
		Address:  c.baseAreaCase.GetAddressByCode(ctx, userAddress.Address),
		Detail:   userAddress.Detail,
	})
}
