package biz

import (
	"context"
	"encoding/json"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/lib/data"
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

func (c *OrderLogisticsCase) GetFromByOrderId(ctx context.Context, orderId int64) (*app.OrderResponse_Logistics, error) {
	orderLogistics, err := c.Find(ctx, &data.OrderLogisticsCondition{
		OrderId: orderId,
	})
	if err != nil {
		return nil, err
	}
	detail := make([]*app.OrderResponse_Logistics_Detail, 0)
	_ = json.Unmarshal([]byte(orderLogistics.Detail), &detail)
	return &app.OrderResponse_Logistics{
		Name:    orderLogistics.Name,
		No:      orderLogistics.No,
		Contact: orderLogistics.Contact,
		Detail:  detail,
	}, nil
}
