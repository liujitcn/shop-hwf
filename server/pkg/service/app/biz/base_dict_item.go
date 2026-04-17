package biz

import (
	"context"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseDictItemCase struct {
	data.BaseDictItemRepo
	baseDictRepo data.BaseDictRepo
}

// NewBaseDictItemCase new a BaseDictItem use case.
func NewBaseDictItemCase(baseDictRepo data.BaseDictRepo, baseDictItemRepo data.BaseDictItemRepo) *BaseDictItemCase {
	return &BaseDictItemCase{
		baseDictRepo:     baseDictRepo,
		BaseDictItemRepo: baseDictItemRepo,
	}
}

func (c *BaseDictItemCase) List(ctx context.Context, condition *data.BaseDictItemCondition) ([]*models.BaseDictItem, error) {
	return c.FindAll(ctx, condition)
}
