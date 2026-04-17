package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type GoodsPropCase struct {
	data.GoodsPropRepo
}

// NewGoodsPropCase new a GoodsProp use case.
func NewGoodsPropCase(goodsPropRepo data.GoodsPropRepo) *GoodsPropCase {
	return &GoodsPropCase{
		GoodsPropRepo: goodsPropRepo,
	}
}

func (c *GoodsPropCase) ListByGoodsId(ctx context.Context, goodsId int64) ([]*app.GoodsResponse_Prop, error) {
	all, err := c.FindAll(ctx, &data.GoodsPropCondition{
		GoodsId: goodsId,
		Status:  int32(common.GoodsStatus_PUT_ON),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*app.GoodsResponse_Prop, 0)
	for _, item := range all {
		list = append(list, c.ConvertToProto(item))
	}
	return list, nil
}

func (c *GoodsPropCase) ConvertToProto(item *models.GoodsProp) *app.GoodsResponse_Prop {
	res := &app.GoodsResponse_Prop{
		Label: item.Label,
		Value: item.Value,
	}
	return res
}
