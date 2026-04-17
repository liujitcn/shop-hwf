package biz

import (
	"context"
	"encoding/json"
	"gitee.com/liujit/shop/server/api/admin"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/utils/timeutil"
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

func (c *OrderRefundCase) GetFromByOrderId(ctx context.Context, orderId int64) ([]*admin.OrderRefund, error) {
	orderRefund, err := c.FindAll(ctx, &data.OrderRefundCondition{
		OrderId: orderId,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*admin.OrderRefund, 0)
	for _, item := range orderRefund {
		var amount admin.OrderRefund_Amount
		_ = json.Unmarshal([]byte(item.Amount), &amount)
		list = append(list, &admin.OrderRefund{
			RefundNo:            item.RefundNo,
			Reason:              common.OrderRefundReason(item.Reason),
			ThirdRefundNo:       item.ThirdRefundNo,
			Channel:             item.Channel,
			UserReceivedAccount: item.UserReceivedAccount,
			CreateTime:          timeutil.TimeToTimeString(item.CreateTime),
			SuccessTime:         timeutil.TimeToTimeString(item.SuccessTime),
			RefundState:         item.RefundState,
			FundsAccount:        item.FundsAccount,
			Amount:              &amount,
			Status:              common.OrderBillStatus(item.Status),
		})
	}
	return list, nil
}
