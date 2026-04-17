package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
)

type GoodsSkuCase struct {
	data.GoodsSkuRepo
}

// NewGoodsSkuCase new a GoodsSku use case.
func NewGoodsSkuCase(goodsSkuRepo data.GoodsSkuRepo) *GoodsSkuCase {
	return &GoodsSkuCase{
		GoodsSkuRepo: goodsSkuRepo,
	}
}

func (c *GoodsSkuCase) MapBySkuCodes(ctx context.Context, skuCodes []string) (map[string]*models.GoodsSku, error) {
	all, err := c.FindAll(ctx, &data.GoodsSkuCondition{
		SkuCodes: skuCodes,
	})
	if err != nil {
		return nil, err
	}
	res := make(map[string]*models.GoodsSku)
	for _, item := range all {
		res[item.SkuCode] = item
	}
	return res, nil
}

func (c *GoodsSkuCase) ListByGoodsId(ctx context.Context, goodsId int64, member bool) ([]*app.GoodsResponse_Sku, error) {
	all, err := c.FindAll(ctx, &data.GoodsSkuCondition{
		GoodsId: goodsId,
		Status:  int32(common.Status_ENABLE),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*app.GoodsResponse_Sku, 0)
	for _, item := range all {
		list = append(list, c.ConvertToProto(item, member))
	}
	return list, nil
}

func (c *GoodsSkuCase) ConvertToProto(item *models.GoodsSku, member bool) *app.GoodsResponse_Sku {
	price := item.Price
	if member {
		price = item.DiscountPrice
	}
	res := &app.GoodsResponse_Sku{
		Picture:   item.Picture,
		SpecItem:  str.ConvertJsonStringToStringArray(item.SpecItem),
		SkuCode:   item.SkuCode,
		Price:     price,
		SaleNum:   item.InitSaleNum + item.RealSaleNum,
		Inventory: item.Inventory,
	}
	return res
}
