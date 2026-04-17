package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
)

type GoodsSpecCase struct {
	data.GoodsSpecRepo
}

// NewGoodsSpecCase new a GoodsSpec use case.
func NewGoodsSpecCase(goodsSpecRepo data.GoodsSpecRepo) *GoodsSpecCase {
	return &GoodsSpecCase{
		GoodsSpecRepo: goodsSpecRepo,
	}
}

func (c *GoodsSpecCase) ListByGoodsId(ctx context.Context, goodsId int64) ([]*app.GoodsResponse_Spec, error) {
	all, err := c.FindAll(ctx, &data.GoodsSpecCondition{
		GoodsId: goodsId,
		Status:  int32(common.Status_ENABLE),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*app.GoodsResponse_Spec, 0)
	for _, item := range all {
		list = append(list, c.ConvertToProto(item))
	}
	return list, nil
}

func (c *GoodsSpecCase) ConvertToProto(item *models.GoodsSpec) *app.GoodsResponse_Spec {
	itemList := make([]*app.GoodsResponse_Spec_Item, 0)
	items := str.ConvertJsonStringToStringArray(item.Item)
	for _, item := range items {
		itemList = append(itemList, &app.GoodsResponse_Spec_Item{
			Name: item,
		})
	}

	res := &app.GoodsResponse_Spec{
		Name: item.Name,
		Item: itemList,
	}
	return res
}
