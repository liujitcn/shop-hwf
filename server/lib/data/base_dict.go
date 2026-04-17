package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseDictCondition struct {
	Id     int64
	Status int32
	Name   string   // 字典名称
	Code   string   // 字典代码
	Codes  []string // 字典代码
}

type BaseDictRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseDict *models.BaseDict) error
	Create(ctx context.Context, baseDict *models.BaseDict) error
	Find(ctx context.Context, condition *BaseDictCondition) (*models.BaseDict, error)
	FindAll(ctx context.Context, condition *BaseDictCondition) ([]*models.BaseDict, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseDictCondition) ([]*models.BaseDict, int64, error)
	Count(ctx context.Context, condition *BaseDictCondition) (int64, error)
}

type baseDictRepo struct {
	data *Data
}

func NewBaseDictRepo(data *Data) BaseDictRepo {
	return &baseDictRepo{data: data}
}

func (r *baseDictRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).BaseDict
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseDictRepo) UpdateByID(ctx context.Context, baseDict *models.BaseDict) error {
	if baseDict.ID == 0 {
		return errors.New("baseDict can not update without id")
	}
	q := r.data.Query(ctx).BaseDict
	_, err := q.WithContext(ctx).Updates(baseDict)
	return err
}

func (r *baseDictRepo) Create(ctx context.Context, baseDict *models.BaseDict) error {
	q := r.data.Query(ctx).BaseDict
	err := q.WithContext(ctx).Clauses().Create(baseDict)
	return err
}

func (r *baseDictRepo) Find(ctx context.Context, condition *BaseDictCondition) (*models.BaseDict, error) {
	m := r.data.Query(ctx).BaseDict
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Eq(condition.Name))
	}
	if condition.Code != "" {
		q = q.Where(m.Code.Eq(condition.Code))
	}
	if len(condition.Codes) > 0 {
		q = q.Where(m.Code.In(condition.Codes...))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

func (r *baseDictRepo) FindAll(ctx context.Context, condition *BaseDictCondition) ([]*models.BaseDict, error) {
	m := r.data.Query(ctx).BaseDict
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.Code != "" {
		q = q.Where(m.Code.Like(buildLikeValue(condition.Code)))
	}
	if len(condition.Codes) > 0 {
		q = q.Where(m.Code.In(condition.Codes...))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.UpdatedBy.Desc())
	return q.Find()
}

func (r *baseDictRepo) ListPage(ctx context.Context, page, size int64, condition *BaseDictCondition) ([]*models.BaseDict, int64, error) {
	m := r.data.Query(ctx).BaseDict
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.Code != "" {
		q = q.Where(m.Code.Like(buildLikeValue(condition.Code)))
	}
	if len(condition.Codes) > 0 {
		q = q.Where(m.Code.In(condition.Codes...))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.UpdatedBy.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseDictRepo) Count(ctx context.Context, condition *BaseDictCondition) (int64, error) {
	m := r.data.Query(ctx).BaseDict
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.Code != "" {
		q = q.Where(m.Code.Like(buildLikeValue(condition.Code)))
	}
	if len(condition.Codes) > 0 {
		q = q.Where(m.Code.In(condition.Codes...))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	count, err := q.Count()
	return count, err
}
