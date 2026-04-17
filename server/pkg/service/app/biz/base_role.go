package biz

import (
	"context"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseRoleCase struct {
	data.BaseRoleRepo
}

// NewBaseRoleCase new a BaseRole use case.
func NewBaseRoleCase(baseRoleRepo data.BaseRoleRepo) *BaseRoleCase {
	return &BaseRoleCase{
		BaseRoleRepo: baseRoleRepo,
	}
}

func (c *BaseRoleCase) GetFromID(ctx context.Context, id int64) (*models.BaseRole, error) {
	return c.Find(ctx, &data.BaseRoleCondition{Id: id})
}
