package biz

import (
	"gitee.com/liujit/shop/server/lib/data"
)

type PayBillCase struct {
	data.PayBillRepo
}

// NewPayBillCase new a ShopPayBill use case.
func NewPayBillCase(
	payBillRepo data.PayBillRepo,
) *PayBillCase {
	return &PayBillCase{
		PayBillRepo: payBillRepo,
	}
}
