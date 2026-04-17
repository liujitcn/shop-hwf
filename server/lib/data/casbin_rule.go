package data

import (
	"context"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type CasbinRuleRepo interface {
	Delete(ctx context.Context, roleKeys []string) error
	Create(ctx context.Context, casbinRule []*models.CasbinRule) error
	FindAll(ctx context.Context) ([]*models.CasbinRule, error)
}

type casbinRuleRepo struct {
	data *Data
}

func NewCasbinRuleRepo(data *Data) CasbinRuleRepo {
	return &casbinRuleRepo{data: data}
}

func (r *casbinRuleRepo) Delete(ctx context.Context, roleKeys []string) error {
	if len(roleKeys) == 0 {
		return nil
	}
	q := r.data.Query(ctx).CasbinRule
	_, err := q.WithContext(ctx).Where(q.V0.In(roleKeys...)).Delete()
	return err
}

func (r *casbinRuleRepo) Create(ctx context.Context, casbinRule []*models.CasbinRule) error {
	if len(casbinRule) == 0 {
		return nil
	}
	q := r.data.Query(ctx).CasbinRule
	err := q.WithContext(ctx).Clauses().CreateInBatches(casbinRule, 100)
	return err
}

func (r *casbinRuleRepo) FindAll(ctx context.Context) ([]*models.CasbinRule, error) {
	m := r.data.Query(ctx).CasbinRule
	q := m.WithContext(ctx)
	return q.Find()
}
