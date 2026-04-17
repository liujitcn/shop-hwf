package biz

import (
	"context"

	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type ShopServiceCase struct {
	data.ShopServiceRepo
}

// NewShopServiceCase new a ShopService use case.
func NewShopServiceCase(shopServiceRepo data.ShopServiceRepo) *ShopServiceCase {
	return &ShopServiceCase{
		ShopServiceRepo: shopServiceRepo,
	}
}

func (c *ShopServiceCase) List(ctx context.Context, condition *data.ShopServiceCondition) ([]*models.ShopService, error) {
	return c.FindAll(ctx, condition)
}

func (c *ShopServiceCase) ConvertToProto(ctx context.Context, item *models.ShopService) *app.ShopService {
	return &app.ShopService{
		Label: item.Label,
		Value: item.Value,
	}
}
