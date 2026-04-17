package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
	"gitee.com/liujit/shop/server/pkg/service/app/util"
)

type OrderGoodsCase struct {
	data.OrderGoodsRepo
	goodsInfoCase *GoodsCase
	goodsSkuCase  *GoodsSkuCase
}

// NewOrderGoodsCase new a OrderGoods use case.
func NewOrderGoodsCase(orderGoodsRepo data.OrderGoodsRepo,
	goodsInfoCase *GoodsCase,
	goodsSkuCase *GoodsSkuCase,
) *OrderGoodsCase {
	return &OrderGoodsCase{
		OrderGoodsRepo: orderGoodsRepo,
		goodsInfoCase:  goodsInfoCase,
		goodsSkuCase:   goodsSkuCase,
	}
}

func (c *OrderGoodsCase) MapByOrderIds(ctx context.Context, orderIds []int64) (map[int64][]*app.OrderGoods, error) {
	all, err := c.FindAll(ctx, &data.OrderGoodsCondition{
		OrderIds: orderIds,
	})
	if err != nil {
		return nil, err
	}
	res := make(map[int64][]*app.OrderGoods)
	for _, item := range all {
		v, ok := res[item.OrderID]
		if !ok {
			v = make([]*app.OrderGoods, 0)
		}
		v = append(v, c.ConvertToProto(item))

		res[item.OrderID] = v
	}
	return res, nil
}

func (c *OrderGoodsCase) ListByOrderId(ctx context.Context, orderId int64) ([]*app.OrderGoods, error) {
	all, err := c.FindAll(ctx, &data.OrderGoodsCondition{
		OrderId: orderId,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*app.OrderGoods, 0)
	for _, item := range all {
		list = append(list, c.ConvertToProto(item))
	}
	return list, nil
}

func (c *OrderGoodsCase) BatchCreate(ctx context.Context, orderId int64, goods []*models.OrderGoods) error {
	if len(goods) == 0 {
		return nil
	}
	for _, item := range goods {
		item.OrderID = orderId
	}
	if len(goods) >= 0 {
		return c.OrderGoodsRepo.BatchCreate(ctx, goods)
	}
	return nil
}

func (c *OrderGoodsCase) ConvertToModelList(ctx context.Context, goods []*app.CreateOrderGoods) ([]*models.OrderGoods, error) {
	// 是否会员
	member := util.IsMember(ctx)

	orderGoodsList := make([]*models.OrderGoods, 0)
	for _, item := range goods {
		orderGoods, err := c.ConvertToModel(ctx, member, item)
		if err != nil {
			return nil, err
		}
		orderGoodsList = append(orderGoodsList, orderGoods)
	}
	return orderGoodsList, nil
}

func (c *OrderGoodsCase) ConvertToProto(item *models.OrderGoods) *app.OrderGoods {
	res := &app.OrderGoods{
		GoodsId:       item.GoodsID,
		SkuCode:       item.SkuCode,
		Picture:       item.Picture,
		Name:          item.Name,
		Num:           item.Num,
		SpecItem:      str.ConvertJsonStringToStringArray(item.SpecItem),
		Price:         item.Price,
		PayPrice:      item.PayPrice,
		TotalPrice:    item.TotalPrice,
		TotalPayPrice: item.TotalPayPrice,
	}
	return res
}

func (c *OrderGoodsCase) ConvertToProtoByCreateOrderGoods(ctx context.Context, member bool, item *app.CreateOrderGoods) (*app.OrderGoods, error) {
	// 查询商品信息和规格信息
	goodsInfo, err := c.goodsInfoCase.Find(ctx, &data.GoodsCondition{
		Id:     item.GetGoodsId(),
		Status: int32(common.Status_ENABLE),
	})
	if err != nil {
		return nil, err
	}
	var goodsSku *models.GoodsSku
	goodsSku, err = c.goodsSkuCase.Find(ctx, &data.GoodsSkuCondition{
		SkuCode: item.GetSkuCode(),
		GoodsId: item.GetGoodsId(),
	})
	if err != nil {
		return nil, err
	}
	picture := goodsInfo.Picture
	if len(goodsSku.Picture) > 0 {
		picture = goodsSku.Picture
	}

	// 支付价格
	payPrice := goodsSku.Price
	if member {
		payPrice = goodsSku.DiscountPrice
	}
	res := &app.OrderGoods{
		GoodsId:       goodsInfo.ID,
		SkuCode:       goodsSku.SkuCode,
		Picture:       picture,
		Name:          goodsInfo.Name,
		Num:           item.GetNum(),
		SpecItem:      str.ConvertJsonStringToStringArray(goodsSku.SpecItem),
		Price:         goodsSku.Price,
		PayPrice:      payPrice,
		TotalPrice:    goodsSku.Price * item.GetNum(),
		TotalPayPrice: payPrice * item.GetNum(),
	}
	return res, nil
}

func (c *OrderGoodsCase) ConvertToModel(ctx context.Context, member bool, goods *app.CreateOrderGoods) (*models.OrderGoods, error) {
	// 查询商品信息和规格信息
	goodsInfo, err := c.goodsInfoCase.Find(ctx, &data.GoodsCondition{
		Id:     goods.GoodsId,
		Status: int32(common.Status_ENABLE),
	})
	if err != nil {
		return nil, err
	}
	var goodsSku *models.GoodsSku
	goodsSku, err = c.goodsSkuCase.Find(ctx, &data.GoodsSkuCondition{
		SkuCode: goods.SkuCode,
		GoodsId: goods.GoodsId,
	})
	if err != nil {
		return nil, err
	}
	picture := goodsInfo.Picture
	if len(goodsSku.Picture) > 0 {
		picture = goodsSku.Picture
	}

	// 支付价格
	payPrice := goodsSku.Price
	if member {
		payPrice = goodsSku.DiscountPrice
	}

	res := &models.OrderGoods{
		GoodsID:       goodsInfo.ID,
		SkuCode:       goodsSku.SkuCode,
		Picture:       picture,
		Name:          goodsInfo.Name,
		Num:           goods.Num,
		SpecItem:      goodsSku.SpecItem,
		Price:         goodsSku.Price,
		PayPrice:      payPrice,
		TotalPrice:    goodsSku.Price * goods.Num,
		TotalPayPrice: payPrice * goods.Num,
	}
	return res, nil
}
