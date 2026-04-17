package biz

import (
	"context"
	"encoding/json"
	"errors"
	"gitee.com/liujit/shop/server/api/admin"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/utils/timeutil"
	"gorm.io/gorm"
)

type OrderLogisticsCase struct {
	data.OrderLogisticsRepo
}

// NewOrderLogisticsCase new a OrderLogistics use case.
func NewOrderLogisticsCase(orderLogisticsRepo data.OrderLogisticsRepo,
) *OrderLogisticsCase {
	return &OrderLogisticsCase{
		OrderLogisticsRepo: orderLogisticsRepo,
	}
}

func (c *OrderLogisticsCase) GetFromByOrderId(ctx context.Context, orderId int64) (*admin.OrderLogistics, error) {
	orderLogistics, err := c.Find(ctx, &data.OrderLogisticsCondition{
		OrderId: orderId,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &admin.OrderLogistics{}, nil
		}
		return nil, err
	}
	detail := make([]*admin.OrderLogistics_Detail, 0)
	_ = json.Unmarshal([]byte(orderLogistics.Detail), &detail)
	return &admin.OrderLogistics{
		Name:      orderLogistics.Name,
		No:        orderLogistics.No,
		Contact:   orderLogistics.Contact,
		Detail:    detail,
		CreatedAt: timeutil.TimeToTimeString(orderLogistics.CreatedAt),
	}, nil
}
