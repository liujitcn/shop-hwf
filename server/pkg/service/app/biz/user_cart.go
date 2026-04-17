package biz

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
	"gitee.com/liujit/shop/server/lib/utils/trans"
	"gitee.com/liujit/shop/server/pkg/service/app/util"
	authMiddleware "go.newcapec.cn/ncttools/nmskit-auth/middleware"
	"go.newcapec.cn/ncttools/nmskit/log"
	"gorm.io/gorm"
)

type UserCartCase struct {
	data.UserCartRepo
	goodsInfoCase *GoodsCase
	goodsSkuCase  *GoodsSkuCase
}

// NewUserCartCase new a UserCart use case.
func NewUserCartCase(
	userCartRepo data.UserCartRepo,
	goodsInfoCase *GoodsCase,
	goodsSkuCase *GoodsSkuCase,
) *UserCartCase {
	return &UserCartCase{
		UserCartRepo:  userCartRepo,
		goodsInfoCase: goodsInfoCase,
		goodsSkuCase:  goodsSkuCase,
	}
}

func (c *UserCartCase) GetFromID(ctx context.Context, id int64) (*models.UserCart, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	return c.Find(ctx, authInfo.UserId, &data.UserCartCondition{
		Id: id,
	})
}

func (c *UserCartCase) Create(ctx context.Context, userCart *app.CreateUserCartRequest) error {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return common.ErrorAccessForbidden("用户认证失败")
	}
	member := util.IsMemberByAuthInfo(authInfo)
	// 查询是否存在
	var find *models.UserCart
	find, err = c.Find(ctx, authInfo.UserId, &data.UserCartCondition{
		GoodsId: userCart.GetGoodsId(),
		SkuCode: userCart.GetSkuCode(),
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var sku *models.GoodsSku
			sku, err = c.goodsSkuCase.Find(ctx, &data.GoodsSkuCondition{
				SkuCode: userCart.GetSkuCode(),
			})
			if err != nil {
				return err
			}
			price := sku.Price
			if member {
				price = sku.DiscountPrice
			}

			userCard := &models.UserCart{
				UserID:    authInfo.UserId,
				GoodsID:   userCart.GetGoodsId(),
				SkuCode:   userCart.GetSkuCode(),
				Num:       userCart.GetNum(),
				Price:     price,
				IsChecked: trans.Bool(true),
			}
			return c.UserCartRepo.Create(ctx, userCard)
		} else {
			return err
		}
	}

	// 更新
	find.Num += userCart.GetNum()
	return c.UpdateByID(ctx, authInfo.UserId, find)
}

func (c *UserCartCase) List(ctx context.Context) (*app.ListUserCartResponse, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	member := util.IsMemberByAuthInfo(authInfo)
	var all []*models.UserCart
	all, err = c.FindAll(ctx, authInfo.UserId, &data.UserCartCondition{})

	goodsIds := make([]int64, 0)
	skuCodes := make([]string, 0)
	for _, info := range all {
		goodsIds = append(goodsIds, info.GoodsID)
		skuCodes = append(skuCodes, info.SkuCode)
	}

	var goodsInfoMap map[int64]*models.Goods
	goodsInfoMap, err = c.goodsInfoCase.MapByGoodsIds(ctx, goodsIds)
	if err != nil {
		return nil, err
	}
	var goodsSkuMap map[string]*models.GoodsSku
	goodsSkuMap, err = c.goodsSkuCase.MapBySkuCodes(ctx, skuCodes)
	if err != nil {
		return nil, err
	}

	list := make([]*app.UserCart, 0)
	for _, item := range all {
		sku, ok1 := goodsSkuMap[item.SkuCode]
		if !ok1 {
			sku = &models.GoodsSku{}
		}
		goods, ok2 := goodsInfoMap[item.GoodsID]
		if !ok2 {
			goods = &models.Goods{}
		}

		picture := goods.Picture
		if len(sku.Picture) > 0 {
			picture = sku.Picture
		}

		price := sku.Price
		if member {
			price = sku.DiscountPrice
		}

		cart := &app.UserCart{
			Id:        item.ID,
			GoodsId:   item.GoodsID,
			SkuCode:   item.SkuCode,
			Picture:   picture,
			Name:      goods.Name,
			Num:       item.Num,
			SpecItem:  str.ConvertJsonStringToStringArray(sku.SpecItem),
			Inventory: sku.Inventory,
			Price:     price,
			JoinPrice: item.Price,
			IsChecked: trans.BoolValue(item.IsChecked),
		}
		list = append(list, cart)
	}
	return &app.ListUserCartResponse{
		List: list,
	}, nil
}
