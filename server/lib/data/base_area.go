package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseAreaCondition struct {
	Id      int64
	Ids     []int64
	parenId int64
	Name    string
	Code    string
}

type BaseAreaRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseArea *models.BaseArea) error
	Create(ctx context.Context, baseArea *models.BaseArea) error
	Find(ctx context.Context, condition *BaseAreaCondition) (*models.BaseArea, error)
	FindAll(ctx context.Context, condition *BaseAreaCondition) ([]*models.BaseArea, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseAreaCondition) ([]*models.BaseArea, int64, error)
	Count(ctx context.Context, condition *BaseAreaCondition) (int64, error)
}

type baseAreaRepo struct {
	data *Data
}

func NewBaseAreaRepo(data *Data) BaseAreaRepo {
	return &baseAreaRepo{data: data}
}

func (r *baseAreaRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).BaseArea
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseAreaRepo) UpdateByID(ctx context.Context, baseArea *models.BaseArea) error {
	if baseArea.ID == 0 {
		return errors.New("baseArea can not update without id")
	}
	q := r.data.Query(ctx).BaseArea
	_, err := q.WithContext(ctx).Updates(baseArea)
	return err
}

func (r *baseAreaRepo) Create(ctx context.Context, baseArea *models.BaseArea) error {
	q := r.data.Query(ctx).BaseArea
	err := q.WithContext(ctx).Clauses().Create(baseArea)
	return err
}

func (r *baseAreaRepo) Find(ctx context.Context, condition *BaseAreaCondition) (*models.BaseArea, error) {
	m := r.data.Query(ctx).BaseArea
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.parenId > 0 {
		q = q.Where(m.ParentID.Eq(condition.parenId))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Eq(condition.Name))
	}
	return q.First()
}

func (r *baseAreaRepo) FindAll(ctx context.Context, condition *BaseAreaCondition) ([]*models.BaseArea, error) {
	m := r.data.Query(ctx).BaseArea
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.parenId > 0 {
		q = q.Where(m.ParentID.Eq(condition.parenId))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	return q.Find()
}

func (r *baseAreaRepo) ListPage(ctx context.Context, page, size int64, condition *BaseAreaCondition) ([]*models.BaseArea, int64, error) {
	m := r.data.Query(ctx).BaseArea
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.parenId > 0 {
		q = q.Where(m.ParentID.Eq(condition.parenId))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseAreaRepo) Count(ctx context.Context, condition *BaseAreaCondition) (int64, error) {
	m := r.data.Query(ctx).BaseArea
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.parenId > 0 {
		q = q.Where(m.ParentID.Eq(condition.parenId))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	count, err := q.Count()
	return count, err
}
