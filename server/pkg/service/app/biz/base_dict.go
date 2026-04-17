package biz

import (
	"context"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseDictCase struct {
	data.BaseDictRepo
}

// NewBaseDictCase new a BaseDict use case.
func NewBaseDictCase(baseDictRepo data.BaseDictRepo) *BaseDictCase {
	return &BaseDictCase{
		BaseDictRepo: baseDictRepo,
	}
}

func (c *BaseDictCase) List(ctx context.Context, condition *data.BaseDictCondition) ([]*models.BaseDict, error) {
	return c.FindAll(ctx, condition)
}
