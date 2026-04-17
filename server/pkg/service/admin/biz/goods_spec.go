package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/admin"
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

func (c *GoodsSpecCase) ListByGoodsId(ctx context.Context, goodsId int64) ([]*admin.GoodsSpec, error) {
	all, err := c.FindAll(ctx, &data.GoodsSpecCondition{
		GoodsId: goodsId,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*admin.GoodsSpec, 0)
	for _, item := range all {
		list = append(list, c.ConvertToProto(item))
	}
	return list, nil
}

func (c *GoodsSpecCase) BatchCreate(ctx context.Context, goodsId int64, spec []*admin.GoodsSpec) error {
	// 查询旧数据
	oldSpecList, err := c.FindAll(ctx, &data.GoodsSpecCondition{
		GoodsId: goodsId,
	})
	if err != nil {
		return err
	}
	oldSpecByID := make(map[int64]*models.GoodsSpec)
	oldSpecByName := make(map[string]*models.GoodsSpec)
	for _, oldSpec := range oldSpecList {
		oldSpecByID[oldSpec.ID] = oldSpec
		oldSpecByName[oldSpec.Name] = oldSpec
	}
	specList := make([]*models.GoodsSpec, 0)
	for _, item := range spec {
		if oldSpec, ok := oldSpecByID[item.Id]; ok {
			item.Id = oldSpec.ID
			item.GoodsId = goodsId
			err = c.UpdateByID(ctx, c.ConvertToModel(item))
			if err != nil {
				return err
			}
			delete(oldSpecByID, oldSpec.ID)
			delete(oldSpecByName, oldSpec.Name)
		} else if oldSpec, ok := oldSpecByName[item.Name]; ok {
			item.Id = oldSpec.ID
			item.GoodsId = goodsId
			err = c.UpdateByID(ctx, c.ConvertToModel(item))
			if err != nil {
				return err
			}
			delete(oldSpecByID, oldSpec.ID)
			delete(oldSpecByName, oldSpec.Name)
		} else {
			item.Id = 0
			item.GoodsId = goodsId
			specList = append(specList, c.ConvertToModel(item))
		}
	}
	if len(oldSpecByID) > 0 {
		oldSpecId := make([]int64, 0)
		for id := range oldSpecByID {
			oldSpecId = append(oldSpecId, id)
		}
		err = c.Delete(ctx, oldSpecId)
		if err != nil {
			return err
		}
	}
	if len(specList) > 0 {
		return c.GoodsSpecRepo.BatchCreate(ctx, specList)
	}
	return nil
}
func (c *GoodsSpecCase) ConvertToProto(item *models.GoodsSpec) *admin.GoodsSpec {
	res := &admin.GoodsSpec{
		Id:      item.ID,
		GoodsId: item.GoodsID,
		Name:    item.Name,
		Item:    str.ConvertJsonStringToStringArray(item.Item),
		Sort:    item.Sort,
	}
	return res
}

func (c *GoodsSpecCase) ConvertToModel(item *admin.GoodsSpec) *models.GoodsSpec {
	res := &models.GoodsSpec{
		ID:      item.GetId(),
		GoodsID: item.GetGoodsId(),
		Name:    item.GetName(),
		Item:    str.ConvertStringArrayToString(item.GetItem()),
		Sort:    item.GetSort(),
	}
	return res
}
