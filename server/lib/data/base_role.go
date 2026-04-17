package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseRoleCondition struct {
	Id     int64
	Ids    []int64
	Status int32
	Name   string
	Code   string
}

type BaseRoleRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseRole *models.BaseRole) error
	Create(ctx context.Context, baseRole *models.BaseRole) error
	Find(ctx context.Context, condition *BaseRoleCondition) (*models.BaseRole, error)
	FindAll(ctx context.Context, condition *BaseRoleCondition) ([]*models.BaseRole, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseRoleCondition) ([]*models.BaseRole, int64, error)
	Count(ctx context.Context, condition *BaseRoleCondition) (int64, error)
}

type baseRoleRepo struct {
	data *Data
}

func NewBaseRoleRepo(data *Data) BaseRoleRepo {
	return &baseRoleRepo{data: data}
}

func (r *baseRoleRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).BaseRole
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseRoleRepo) UpdateByID(ctx context.Context, baseRole *models.BaseRole) error {
	if baseRole.ID == 0 {
		return errors.New("baseRole can not update without id")
	}
	q := r.data.Query(ctx).BaseRole
	_, err := q.WithContext(ctx).Updates(baseRole)
	return err
}

func (r *baseRoleRepo) Create(ctx context.Context, baseRole *models.BaseRole) error {
	q := r.data.Query(ctx).BaseRole
	err := q.WithContext(ctx).Clauses().Create(baseRole)
	return err
}

func (r *baseRoleRepo) Find(ctx context.Context, condition *BaseRoleCondition) (*models.BaseRole, error) {
	m := r.data.Query(ctx).BaseRole
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if len(condition.Name) > 0 {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if len(condition.Code) > 0 {
		q = q.Where(m.Code.Like(buildLikeValue(condition.Code)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

func (r *baseRoleRepo) FindAll(ctx context.Context, condition *BaseRoleCondition) ([]*models.BaseRole, error) {
	m := r.data.Query(ctx).BaseRole
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if len(condition.Name) > 0 {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if len(condition.Code) > 0 {
		q = q.Where(m.Code.Like(buildLikeValue(condition.Code)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.UpdatedAt.Desc())
	return q.Find()
}

func (r *baseRoleRepo) ListPage(ctx context.Context, page, size int64, condition *BaseRoleCondition) ([]*models.BaseRole, int64, error) {
	m := r.data.Query(ctx).BaseRole
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if len(condition.Name) > 0 {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if len(condition.Code) > 0 {
		q = q.Where(m.Code.Like(buildLikeValue(condition.Code)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.UpdatedAt.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseRoleRepo) Count(ctx context.Context, condition *BaseRoleCondition) (int64, error) {
	m := r.data.Query(ctx).BaseRole
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Name) > 0 {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if len(condition.Code) > 0 {
		q = q.Where(m.Code.Like(buildLikeValue(condition.Code)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	count, err := q.Count()
	return count, err
}
