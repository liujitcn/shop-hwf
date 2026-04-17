package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseMenuCondition struct {
	Id       int64
	ParentId *int64
	Ids      []int64
	Status   int32
	Types    []int32
}

type BaseMenuRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseMenu *models.BaseMenu) error
	Create(ctx context.Context, baseMenu *models.BaseMenu) error
	Find(ctx context.Context, condition *BaseMenuCondition) (*models.BaseMenu, error)
	FindAll(ctx context.Context, condition *BaseMenuCondition) ([]*models.BaseMenu, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseMenuCondition) ([]*models.BaseMenu, int64, error)
	Count(ctx context.Context, condition *BaseMenuCondition) (int64, error)
}

type baseMenuRepo struct {
	data *Data
}

func NewBaseMenuRepo(data *Data) BaseMenuRepo {
	return &baseMenuRepo{data: data}
}

func (r *baseMenuRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).BaseMenu
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseMenuRepo) UpdateByID(ctx context.Context, baseMenu *models.BaseMenu) error {
	if baseMenu.ID == 0 {
		return errors.New("baseMenu can not update without id")
	}
	q := r.data.Query(ctx).BaseMenu
	_, err := q.WithContext(ctx).Updates(baseMenu)
	return err
}

func (r *baseMenuRepo) Create(ctx context.Context, baseMenu *models.BaseMenu) error {
	q := r.data.Query(ctx).BaseMenu
	err := q.WithContext(ctx).Clauses().Create(baseMenu)
	return err
}

func (r *baseMenuRepo) Find(ctx context.Context, condition *BaseMenuCondition) (*models.BaseMenu, error) {
	m := r.data.Query(ctx).BaseMenu
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.ParentId != nil {
		q = q.Where(m.ParentID.Eq(*condition.ParentId))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if len(condition.Types) > 0 {
		q = q.Where(m.Type.In(condition.Types...))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

func (r *baseMenuRepo) FindAll(ctx context.Context, condition *BaseMenuCondition) ([]*models.BaseMenu, error) {
	m := r.data.Query(ctx).BaseMenu
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.ParentId != nil {
		q = q.Where(m.ParentID.Eq(*condition.ParentId))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if len(condition.Types) > 0 {
		q = q.Where(m.Type.In(condition.Types...))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.ParentID.Asc(), m.Sort.Asc())
	return q.Find()
}

func (r *baseMenuRepo) ListPage(ctx context.Context, page, size int64, condition *BaseMenuCondition) ([]*models.BaseMenu, int64, error) {
	m := r.data.Query(ctx).BaseMenu
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.ParentId != nil {
		q = q.Where(m.ParentID.Eq(*condition.ParentId))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if len(condition.Types) > 0 {
		q = q.Where(m.Type.In(condition.Types...))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.ParentID.Asc(), m.Sort.Asc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseMenuRepo) Count(ctx context.Context, condition *BaseMenuCondition) (int64, error) {
	m := r.data.Query(ctx).BaseMenu
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.ParentId != nil {
		q = q.Where(m.ParentID.Eq(*condition.ParentId))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if len(condition.Types) > 0 {
		q = q.Where(m.Type.In(condition.Types...))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	count, err := q.Count()
	return count, err
}
