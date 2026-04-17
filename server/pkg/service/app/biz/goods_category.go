package biz

import (
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type GoodsCategoryCase struct {
	data.GoodsCategoryRepo
	goodsInfoRepo data.GoodsRepo
}

// NewGoodsCategoryCase new a GoodsCategory use case.
func NewGoodsCategoryCase(goodsCategoryRepo data.GoodsCategoryRepo, goodsInfoRepo data.GoodsRepo) *GoodsCategoryCase {
	return &GoodsCategoryCase{
		GoodsCategoryRepo: goodsCategoryRepo,
		goodsInfoRepo:     goodsInfoRepo,
	}
}

func (c *GoodsCategoryCase) ConvertToProto(item *models.GoodsCategory) *app.GoodsCategory {
	res := &app.GoodsCategory{
		Id:       item.ID,
		ParentId: item.ParentID,
		Name:     item.Name,
		Picture:  item.Picture,
		Goods:    nil,
	}
	return res
}
