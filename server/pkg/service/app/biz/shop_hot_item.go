package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/pkg/service/app/util"
)

type ShopHotItemCase struct {
	data.ShopHotItemRepo
	shopHotRepo      data.ShopHotRepo
	shopHotGoodsRepo data.ShopHotGoodsRepo
	goodsInfoRepo    data.GoodsRepo
}

// NewShopHotItemCase new a ShopHotItem use case.
func NewShopHotItemCase(shopHotRepo data.ShopHotRepo, shopHotItemRepo data.ShopHotItemRepo, shopHotGoodsRepo data.ShopHotGoodsRepo, goodsInfoRepo data.GoodsRepo) *ShopHotItemCase {
	return &ShopHotItemCase{
		ShopHotItemRepo:  shopHotItemRepo,
		shopHotRepo:      shopHotRepo,
		shopHotGoodsRepo: shopHotGoodsRepo,
		goodsInfoRepo:    goodsInfoRepo,
	}
}

func (c *ShopHotItemCase) List(ctx context.Context, condition *data.ShopHotItemCondition) ([]*models.ShopHotItem, error) {
	return c.FindAll(ctx, condition)
}

func (c *ShopHotItemCase) PageGoods(ctx context.Context, req *app.PageShopHotGoodsRequest) (*app.PageShopHotGoodsResponse, error) {
	// 是否会员
	member := util.IsMember(ctx)
	// 查询商品id
	goodsIds, count, err := c.shopHotGoodsRepo.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), &data.ShopHotGoodsCondition{
		HotItemId: req.HotItemId,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*app.Goods, 0)
	if count > 0 {
		var all []*models.Goods
		all, err = c.goodsInfoRepo.FindAll(ctx, &data.GoodsCondition{
			Ids:    goodsIds,
			Status: int32(common.Status_ENABLE),
		})
		if err != nil {
			return nil, err
		}
		for _, item := range all {
			price := item.Price
			if member {
				price = item.DiscountPrice
			}
			list = append(list, &app.Goods{
				Id:      item.ID,
				Name:    item.Name,
				Desc:    item.Desc,
				Picture: item.Picture,
				SaleNum: item.InitSaleNum + item.RealSaleNum,
				Price:   price,
			})
		}
	}
	return &app.PageShopHotGoodsResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *ShopHotItemCase) ConvertToProto(item *models.ShopHotItem) *app.ShopHotItem {
	res := &app.ShopHotItem{
		Id:    item.ID,
		Title: item.Title,
	}
	return res
}
