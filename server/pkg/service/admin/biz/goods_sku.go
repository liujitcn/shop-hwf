package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/admin"
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
func (c *GoodsSkuCase) GetFromID(ctx context.Context, id int64) (*models.GoodsSku, error) {
	return c.Find(ctx, &data.GoodsSkuCondition{
		Id: id,
	})
}

func (c *GoodsSkuCase) Page(ctx context.Context, req *admin.PageGoodsSkuRequest) (*admin.PageGoodsSkuResponse, error) {
	condition := &data.GoodsSkuCondition{
		GoodsId: req.GetGoodsId(),
		SkuCode: req.GetSkuCode(),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}
	list := make([]*admin.GoodsSku, 0)
	for _, item := range page {
		list = append(list, c.ConvertToProto(item))
	}

	return &admin.PageGoodsSkuResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *GoodsSkuCase) BatchCreate(ctx context.Context, goodsId int64, sku []*admin.GoodsSku) error {
	// 查询旧数据
	oldSkuList, err := c.FindAll(ctx, &data.GoodsSkuCondition{
		GoodsId: goodsId,
	})
	if err != nil {
		return err
	}
	oldSkuByID := make(map[int64]*models.GoodsSku)
	oldSkuByCode := make(map[string]*models.GoodsSku)
	for _, oldSku := range oldSkuList {
		oldSkuByID[oldSku.ID] = oldSku
		oldSkuByCode[oldSku.SkuCode] = oldSku
	}

	skuList := make([]*models.GoodsSku, 0)
	for _, item := range sku {
		if oldSku, ok := oldSkuByID[item.Id]; ok {
			item.Id = oldSku.ID
			item.GoodsId = goodsId
			err = c.UpdateByID(ctx, c.ConvertToModel(item))
			if err != nil {
				return err
			}
			delete(oldSkuByID, oldSku.ID)
			delete(oldSkuByCode, oldSku.SkuCode)
		} else if oldSku, ok := oldSkuByCode[item.SkuCode]; ok {
			item.Id = oldSku.ID
			item.GoodsId = goodsId
			err = c.UpdateByID(ctx, c.ConvertToModel(item))
			if err != nil {
				return err
			}
			delete(oldSkuByID, oldSku.ID)
			delete(oldSkuByCode, oldSku.SkuCode)
		} else {
			item.Id = 0
			item.GoodsId = goodsId
			skuList = append(skuList, c.ConvertToModel(item))
		}
	}
	if len(oldSkuByID) > 0 {
		oldSkuId := make([]int64, 0)
		for id := range oldSkuByID {
			oldSkuId = append(oldSkuId, id)
		}
		err = c.Delete(ctx, oldSkuId)
		if err != nil {
			return err
		}
	}
	if len(skuList) > 0 {
		return c.GoodsSkuRepo.BatchCreate(ctx, skuList)
	}
	return nil
}

func (c *GoodsSkuCase) ListByGoodsId(ctx context.Context, goodsId int64) ([]*admin.GoodsSku, error) {
	all, err := c.FindAll(ctx, &data.GoodsSkuCondition{
		GoodsId: goodsId,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*admin.GoodsSku, 0)
	for _, item := range all {
		list = append(list, c.ConvertToProto(item))
	}
	return list, nil
}

func (c *GoodsSkuCase) ConvertToProto(item *models.GoodsSku) *admin.GoodsSku {
	res := &admin.GoodsSku{
		Id:            item.ID,
		GoodsId:       item.GoodsID,
		Picture:       item.Picture,
		SpecItem:      str.ConvertJsonStringToStringArray(item.SpecItem),
		SkuCode:       item.SkuCode,
		Price:         item.Price,
		DiscountPrice: item.DiscountPrice,
		InitSaleNum:   item.InitSaleNum,
		RealSaleNum:   item.RealSaleNum,
		Inventory:     item.Inventory,
	}
	return res
}

func (c *GoodsSkuCase) ConvertToModel(item *admin.GoodsSku) *models.GoodsSku {
	res := &models.GoodsSku{
		ID:            item.GetId(),
		GoodsID:       item.GetGoodsId(),
		Picture:       item.GetPicture(),
		SpecItem:      str.ConvertStringArrayToString(item.SpecItem),
		SkuCode:       item.GetSkuCode(),
		Price:         item.GetPrice(),
		DiscountPrice: item.GetDiscountPrice(),
		InitSaleNum:   item.GetInitSaleNum(),
		RealSaleNum:   item.GetRealSaleNum(),
		Inventory:     item.GetInventory(),
	}
	return res
}
