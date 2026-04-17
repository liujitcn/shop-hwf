package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/admin"
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
func (c *GoodsPropCase) GetFromID(ctx context.Context, id int64) (*models.GoodsProp, error) {
	return c.Find(ctx, &data.GoodsPropCondition{
		Id: id,
	})
}

func (c *GoodsPropCase) Page(ctx context.Context, req *admin.PageGoodsPropRequest) (*admin.PageGoodsPropResponse, error) {
	condition := &data.GoodsPropCondition{
		GoodsId: req.GetGoodsId(),
		Label:   req.GetLabel(),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}
	list := make([]*admin.GoodsProp, 0)
	for _, item := range page {
		list = append(list, c.ConvertToProto(item))
	}

	return &admin.PageGoodsPropResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *GoodsPropCase) ListByGoodsId(ctx context.Context, goodsId int64) ([]*admin.GoodsProp, error) {
	all, err := c.FindAll(ctx, &data.GoodsPropCondition{
		GoodsId: goodsId,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*admin.GoodsProp, 0)
	for _, item := range all {
		list = append(list, c.ConvertToProto(item))
	}
	return list, nil
}

func (c *GoodsPropCase) BatchCreate(ctx context.Context, goodsId int64, prop []*admin.GoodsProp) error {
	// 查询旧数据
	oldPropList, err := c.FindAll(ctx, &data.GoodsPropCondition{
		GoodsId: goodsId,
	})
	if err != nil {
		return err
	}
	oldPropByID := make(map[int64]*models.GoodsProp)
	oldPropByLabel := make(map[string]*models.GoodsProp)
	for _, oldProp := range oldPropList {
		oldPropByID[oldProp.ID] = oldProp
		oldPropByLabel[oldProp.Label] = oldProp
	}

	propList := make([]*models.GoodsProp, 0)
	for _, item := range prop {
		if oldProp, ok := oldPropByID[item.Id]; ok {
			item.Id = oldProp.ID
			item.GoodsId = goodsId
			err = c.UpdateByID(ctx, c.ConvertToModel(item))
			if err != nil {
				return err
			}
			delete(oldPropByID, oldProp.ID)
			delete(oldPropByLabel, oldProp.Label)
		} else if oldProp, ok := oldPropByLabel[item.Label]; ok {
			item.Id = oldProp.ID
			item.GoodsId = goodsId
			err = c.UpdateByID(ctx, c.ConvertToModel(item))
			if err != nil {
				return err
			}
			delete(oldPropByID, oldProp.ID)
			delete(oldPropByLabel, oldProp.Label)
		} else {
			item.Id = 0
			item.GoodsId = goodsId
			propList = append(propList, c.ConvertToModel(item))
		}
	}
	if len(oldPropByID) > 0 {
		oldPropId := make([]int64, 0)
		for id := range oldPropByID {
			oldPropId = append(oldPropId, id)
		}
		err = c.Delete(ctx, oldPropId)
		if err != nil {
			return err
		}
	}
	if len(propList) > 0 {
		return c.GoodsPropRepo.BatchCreate(ctx, propList)
	}
	return nil
}

func (c *GoodsPropCase) ConvertToProto(item *models.GoodsProp) *admin.GoodsProp {
	res := &admin.GoodsProp{
		Id:      item.ID,
		GoodsId: item.GoodsID,
		Label:   item.Label,
		Value:   item.Value,
		Sort:    item.Sort,
	}
	return res
}

func (c *GoodsPropCase) ConvertToModel(item *admin.GoodsProp) *models.GoodsProp {
	res := &models.GoodsProp{
		ID:      item.GetId(),
		GoodsID: item.GetGoodsId(),
		Label:   item.GetLabel(),
		Value:   item.GetValue(),
		Sort:    item.GetSort(),
	}
	return res
}
