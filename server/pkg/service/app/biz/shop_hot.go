package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
)

type ShopHotCase struct {
	data.ShopHotRepo
}

// NewShopHotCase new a ShopHot use case.
func NewShopHotCase(shopHotRepo data.ShopHotRepo) *ShopHotCase {
	return &ShopHotCase{
		ShopHotRepo: shopHotRepo,
	}
}
func (c *ShopHotCase) GetFromID(ctx context.Context, id int64) (*models.ShopHot, error) {
	return c.Find(ctx, &data.ShopHotCondition{
		Id:     id,
		Status: int32(common.Status_ENABLE),
	})
}

func (c *ShopHotCase) List(ctx context.Context, condition *data.ShopHotCondition) ([]*models.ShopHot, error) {
	return c.FindAll(ctx, condition)
}

func (c *ShopHotCase) ConvertToProto(item *models.ShopHot) *app.ShopHot {
	return &app.ShopHot{
		Id:      item.ID,
		Title:   item.Title,
		Desc:    item.Desc,
		Picture: str.ConvertJsonStringToStringArray(item.Picture),
	}
}
