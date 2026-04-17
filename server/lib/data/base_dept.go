package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseDeptCondition struct {
	Id       int64
	ParentId *int64
	Status   int32
}

type BaseDeptRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseDept *models.BaseDept) error
	Create(ctx context.Context, baseDept *models.BaseDept) error
	Find(ctx context.Context, condition *BaseDeptCondition) (*models.BaseDept, error)
	FindAll(ctx context.Context, condition *BaseDeptCondition) ([]*models.BaseDept, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseDeptCondition) ([]*models.BaseDept, int64, error)
	Count(ctx context.Context, condition *BaseDeptCondition) (int64, error)
}

type baseDeptRepo struct {
	data *Data
}

func NewBaseDeptRepo(data *Data) BaseDeptRepo {
	return &baseDeptRepo{data: data}
}

func (r *baseDeptRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).BaseDept
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseDeptRepo) UpdateByID(ctx context.Context, baseDept *models.BaseDept) error {
	if baseDept.ID == 0 {
		return errors.New("baseDept can not update without id")
	}
	q := r.data.Query(ctx).BaseDept
	_, err := q.WithContext(ctx).Updates(baseDept)
	return err
}

func (r *baseDeptRepo) Create(ctx context.Context, baseDept *models.BaseDept) error {
	q := r.data.Query(ctx).BaseDept
	err := q.WithContext(ctx).Clauses().Create(baseDept)
	return err
}

func (r *baseDeptRepo) Find(ctx context.Context, condition *BaseDeptCondition) (*models.BaseDept, error) {
	if condition.Id == 0 {
		return nil, errors.New("baseDept can not find without id")
	}
	m := r.data.Query(ctx).BaseDept
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.ParentId != nil {
		q = q.Where(m.ParentID.Eq(*condition.ParentId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.Sort.Asc())
	return q.First()
}

func (r *baseDeptRepo) FindAll(ctx context.Context, condition *BaseDeptCondition) ([]*models.BaseDept, error) {
	m := r.data.Query(ctx).BaseDept
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.ParentId != nil {
		q = q.Where(m.ParentID.Eq(*condition.ParentId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.ParentID.Asc(), m.Sort.Asc())
	return q.Find()
}

func (r *baseDeptRepo) ListPage(ctx context.Context, page, size int64, condition *BaseDeptCondition) ([]*models.BaseDept, int64, error) {
	m := r.data.Query(ctx).BaseDept
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.ParentId != nil {
		q = q.Where(m.ParentID.Eq(*condition.ParentId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.ParentID.Asc(), m.Sort.Asc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseDeptRepo) Count(ctx context.Context, condition *BaseDeptCondition) (int64, error) {
	m := r.data.Query(ctx).BaseDept
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.ParentId != nil {
		q = q.Where(m.ParentID.Eq(*condition.ParentId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	count, err := q.Count()
	return count, err
}
