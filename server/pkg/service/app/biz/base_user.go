package biz

import (
	"context"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseUserCase struct {
	data.BaseUserRepo
	baseDeptRepo data.BaseDeptRepo
}

// NewBaseUserCase new a BaseUser use case.
func NewBaseUserCase(baseUserRepo data.BaseUserRepo, deptRepo data.BaseDeptRepo) *BaseUserCase {
	return &BaseUserCase{
		BaseUserRepo: baseUserRepo,
		baseDeptRepo: deptRepo,
	}
}

func (c *BaseUserCase) GetFromID(ctx context.Context, id int64) (*models.BaseUser, error) {
	return c.Find(ctx, &data.BaseUserCondition{Id: id})
}

func (c *BaseUserCase) GetFromOpenid(ctx context.Context, openid string) (*models.BaseUser, error) {
	return c.Find(ctx, &data.BaseUserCondition{Openid: openid})
}

func (c *BaseUserCase) GetFromPhone(ctx context.Context, userName string) (*models.BaseUser, error) {
	return c.Find(ctx, &data.BaseUserCondition{Phone: userName})
}
