package biz

import (
	"context"
	"fmt"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/pkg/service/admin/biz"
	"strconv"
)

type ShopBannerCase struct {
	data.ShopBannerRepo
	goodsCategoryCase *biz.GoodsCategoryCase
}

// NewShopBannerCase new a ShopBanner use case.
func NewShopBannerCase(shopBannerRepo data.ShopBannerRepo, goodsCategoryCase *biz.GoodsCategoryCase) *ShopBannerCase {
	return &ShopBannerCase{
		ShopBannerRepo:    shopBannerRepo,
		goodsCategoryCase: goodsCategoryCase,
	}
}

func (c *ShopBannerCase) ConvertToProto(ctx context.Context, item *models.ShopBanner) *app.ShopBanner {
	var href string
	switch common.ShopBannerType(item.Type) {
	case common.ShopBannerType_GOODS_DETAIL:
		href = fmt.Sprintf("id=%s", item.Href)
	case common.ShopBannerType_CATEGORY_DETAIL:
		// 查询分类
		id, err := strconv.ParseInt(item.Href, 10, 64)
		if err == nil {
			var find *models.GoodsCategory
			find, err = c.goodsCategoryCase.Find(ctx, &data.GoodsCategoryCondition{
				Id: id,
			})
			if err == nil && find != nil {
				href = fmt.Sprintf("categoryId=%d&categoryName=%s", find.ID, find.Name)
			}
		}
	default:
		href = item.Href
	}
	res := &app.ShopBanner{
		Site:    common.ShopBannerSite(item.Site),
		Picture: item.Picture,
		Type:    common.ShopBannerType(item.Type),
		Href:    href,
	}
	return res
}
