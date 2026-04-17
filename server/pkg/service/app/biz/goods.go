package biz

import (
	"context"
	"fmt"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
	"gitee.com/liujit/shop/server/pkg/service/app/util"
)

type GoodsCase struct {
	data.GoodsRepo
	goodsCategoryRepo data.GoodsCategoryRepo
	goodsPropCase     *GoodsPropCase
	goodsSpecCase     *GoodsSpecCase
	goodsSkuCase      *GoodsSkuCase
}

// NewGoodsCase new a Goods use case.
func NewGoodsCase(
	goodsInfoRepo data.GoodsRepo,
	goodsCategoryRepo data.GoodsCategoryRepo,
	goodsPropCase *GoodsPropCase,
	goodsSpecCase *GoodsSpecCase,
	goodsSkuCase *GoodsSkuCase,
) *GoodsCase {
	return &GoodsCase{
		GoodsRepo:         goodsInfoRepo,
		goodsCategoryRepo: goodsCategoryRepo,
		goodsPropCase:     goodsPropCase,
		goodsSpecCase:     goodsSpecCase,
		goodsSkuCase:      goodsSkuCase,
	}
}

func (c *GoodsCase) GetFromID(ctx context.Context, id int64) (*app.GoodsResponse, error) {
	// 是否会员
	member := util.IsMember(ctx)
	info, err := c.Find(ctx, &data.GoodsCondition{
		Id:     id,
		Status: int32(common.GoodsStatus_PUT_ON),
	})
	if err != nil {
		return nil, err
	}
	price := info.Price
	if member {
		price = info.DiscountPrice
	}

	goodsInfo := &app.GoodsResponse{
		Id:         info.ID,
		CategoryId: info.CategoryID,
		Name:       info.Name,
		Desc:       info.Desc,
		Price:      price,
		SaleNum:    info.InitSaleNum + info.RealSaleNum,
		Picture:    info.Picture,
		Banner:     str.ConvertJsonStringToStringArray(info.Banner),
		Detail:     str.ConvertJsonStringToStringArray(info.Detail),
	}
	// 属性
	goodsInfo.PropList, err = c.goodsPropCase.ListByGoodsId(ctx, goodsInfo.Id)
	if err != nil {
		return nil, err
	}
	// 规格
	goodsInfo.SpecList, err = c.goodsSpecCase.ListByGoodsId(ctx, goodsInfo.Id)
	if err != nil {
		return nil, err
	}
	// SKU
	goodsInfo.SkuList, err = c.goodsSkuCase.ListByGoodsId(ctx, goodsInfo.Id, member)
	if err != nil {
		return nil, err
	}
	return goodsInfo, nil
}

func (c *GoodsCase) Page(ctx context.Context, req *app.PageGoodsRequest) (*app.PageGoodsResponse, error) {
	// 是否会员
	member := util.IsMember(ctx)
	condition := &data.GoodsCondition{
		CategoryId: req.GetCategoryId(),
		Name:       req.GetName(),
		Status:     int32(common.GoodsStatus_PUT_ON),
	}
	if condition.CategoryId > 0 {
		// 查询分类路径
		category, err := c.goodsCategoryRepo.Find(ctx, &data.GoodsCategoryCondition{Id: req.GetCategoryId()})
		if err != nil {
			return nil, err
		}
		if category.ParentID == 0 {
			condition.CategoryId = 0
			condition.CategoryPath = fmt.Sprintf("%s/", category.Path)
		}
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}
	list := make([]*app.Goods, 0)
	for _, item := range page {
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

	return &app.PageGoodsResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *GoodsCase) MapByGoodsIds(ctx context.Context, goodsIds []int64) (map[int64]*models.Goods, error) {
	all, err := c.FindAll(ctx, &data.GoodsCondition{
		Ids: goodsIds,
	})
	if err != nil {
		return nil, err
	}
	res := make(map[int64]*models.Goods)
	for _, item := range all {
		res[item.ID] = item
	}
	return res, nil
}
