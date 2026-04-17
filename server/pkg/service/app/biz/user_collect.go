package biz

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/pkg/service/app/util"
	authMiddleware "go.newcapec.cn/ncttools/nmskit-auth/middleware"
	"go.newcapec.cn/ncttools/nmskit/log"
	"gorm.io/gorm"
)

type UserCollectCase struct {
	data.UserCollectRepo
	goodsInfoCase *GoodsCase
	goodsSkuCase  *GoodsSkuCase
}

// NewUserCollectCase new a UserCollect use case.
func NewUserCollectCase(
	userCollectRepo data.UserCollectRepo,
	goodsInfoCase *GoodsCase,
	goodsSkuCase *GoodsSkuCase,
) *UserCollectCase {
	return &UserCollectCase{
		UserCollectRepo: userCollectRepo,
		goodsInfoCase:   goodsInfoCase,
		goodsSkuCase:    goodsSkuCase,
	}
}

func (c *UserCollectCase) GetFromID(ctx context.Context, id int64) (*models.UserCollect, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	return c.Find(ctx, authInfo.UserId, &data.UserCollectCondition{
		Id: id,
	})
}

func (c *UserCollectCase) Create(ctx context.Context, userCollect *app.UserCollectForm) error {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return common.ErrorAccessForbidden("用户认证失败")
	}
	member := util.IsMemberByAuthInfo(authInfo)
	// 查询是否存在
	var find *models.UserCollect
	find, err = c.Find(ctx, authInfo.UserId, &data.UserCollectCondition{
		GoodsId: userCollect.GetGoodsId(),
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var sku *models.Goods
			sku, err = c.goodsInfoCase.Find(ctx, &data.GoodsCondition{
				Id: userCollect.GetGoodsId(),
			})
			if err != nil {
				return err
			}
			price := sku.Price
			if member {
				price = sku.DiscountPrice
			}

			userCard := &models.UserCollect{
				UserID:  authInfo.UserId,
				GoodsID: userCollect.GetGoodsId(),
				Price:   price,
			}
			return c.UserCollectRepo.Create(ctx, userCard)
		} else {
			return err
		}
	}

	// 删除
	return c.Delete(ctx, authInfo.UserId, []int64{find.ID})
}

func (c *UserCollectCase) Page(ctx context.Context, req *app.PageUserCollectRequest) (*app.PageUserCollectResponse, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	member := util.IsMemberByAuthInfo(authInfo)
	var page []*models.UserCollect
	var count int64
	page, count, err = c.ListPage(ctx, authInfo.UserId, req.GetPageNum(), req.GetPageSize(), &data.UserCollectCondition{})

	goodsIds := make([]int64, 0)
	for _, info := range page {
		goodsIds = append(goodsIds, info.GoodsID)
	}

	var goodsInfoMap map[int64]*models.Goods
	goodsInfoMap, err = c.goodsInfoCase.MapByGoodsIds(ctx, goodsIds)
	if err != nil {
		return nil, err
	}

	list := make([]*app.UserCollect, 0)
	for _, item := range page {
		goods, ok := goodsInfoMap[item.GoodsID]
		if !ok {
			goods = &models.Goods{}
		}

		price := goods.Price
		if member {
			price = goods.DiscountPrice
		}

		collect := &app.UserCollect{
			Id:        item.ID,
			GoodsId:   item.GoodsID,
			Name:      goods.Name,
			Desc:      goods.Desc,
			Picture:   goods.Picture,
			SaleNum:   goods.InitSaleNum + goods.RealSaleNum,
			Price:     price,
			JoinPrice: item.Price,
		}
		list = append(list, collect)
	}
	return &app.PageUserCollectResponse{
		List:  list,
		Total: int32(count),
	}, nil
}
