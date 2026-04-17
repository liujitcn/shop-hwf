package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseDictItemCondition struct {
	Id      int64
	DictId  int64
	DictIds []int64
	Status  int32
	Label   string
	Value   string
}

type BaseDictItemRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseDictItem *models.BaseDictItem) error
	Create(ctx context.Context, baseDictItem *models.BaseDictItem) error
	Find(ctx context.Context, condition *BaseDictItemCondition) (*models.BaseDictItem, error)
	FindAll(ctx context.Context, condition *BaseDictItemCondition) ([]*models.BaseDictItem, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseDictItemCondition) ([]*models.BaseDictItem, int64, error)
	Count(ctx context.Context, condition *BaseDictItemCondition) (int64, error)
	FindLabelByCodeAndValue(ctx context.Context, code, value string) (string, error)
}

type baseDictItemRepo struct {
	data *Data
}

func NewBaseDictItemRepo(data *Data) BaseDictItemRepo {
	return &baseDictItemRepo{data: data}
}

func (r *baseDictItemRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).BaseDictItem
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseDictItemRepo) UpdateByID(ctx context.Context, baseDictItem *models.BaseDictItem) error {
	if baseDictItem.ID == 0 {
		return errors.New("baseDictItem can not update without id")
	}
	q := r.data.Query(ctx).BaseDictItem
	_, err := q.WithContext(ctx).Updates(baseDictItem)
	return err
}

func (r *baseDictItemRepo) Create(ctx context.Context, baseDictItem *models.BaseDictItem) error {
	q := r.data.Query(ctx).BaseDictItem
	err := q.WithContext(ctx).Clauses().Create(baseDictItem)
	return err
}

func (r *baseDictItemRepo) Find(ctx context.Context, condition *BaseDictItemCondition) (*models.BaseDictItem, error) {
	m := r.data.Query(ctx).BaseDictItem
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.DictId > 0 {
		q = q.Where(m.DictID.Eq(condition.DictId))
	}
	if len(condition.DictIds) > 0 {
		q = q.Where(m.DictID.In(condition.DictIds...))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Eq(condition.Label))
	}
	if condition.Value != "" {
		q = q.Where(m.Value.Eq(condition.Value))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

func (r *baseDictItemRepo) FindAll(ctx context.Context, condition *BaseDictItemCondition) ([]*models.BaseDictItem, error) {
	m := r.data.Query(ctx).BaseDictItem
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.DictId > 0 {
		q = q.Where(m.DictID.Eq(condition.DictId))
	}
	if len(condition.DictIds) > 0 {
		q = q.Where(m.DictID.In(condition.DictIds...))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Like(buildLikeValue(condition.Label)))
	}
	if condition.Value != "" {
		q = q.Where(m.Value.Like(buildLikeValue(condition.Value)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.Sort.Asc())
	return q.Find()
}

func (r *baseDictItemRepo) ListPage(ctx context.Context, page, size int64, condition *BaseDictItemCondition) ([]*models.BaseDictItem, int64, error) {
	m := r.data.Query(ctx).BaseDictItem
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.DictId > 0 {
		q = q.Where(m.DictID.Eq(condition.DictId))
	}
	if len(condition.DictIds) > 0 {
		q = q.Where(m.DictID.In(condition.DictIds...))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Like(buildLikeValue(condition.Label)))
	}
	if condition.Value != "" {
		q = q.Where(m.Value.Like(buildLikeValue(condition.Value)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.Sort.Asc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseDictItemRepo) Count(ctx context.Context, condition *BaseDictItemCondition) (int64, error) {
	m := r.data.Query(ctx).BaseDictItem
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.DictId > 0 {
		q = q.Where(m.DictID.Eq(condition.DictId))
	}
	if len(condition.DictIds) > 0 {
		q = q.Where(m.DictID.In(condition.DictIds...))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Like(buildLikeValue(condition.Label)))
	}
	if condition.Value != "" {
		q = q.Where(m.Value.Like(buildLikeValue(condition.Value)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	count, err := q.Count()
	return count, err
}

func (r *baseDictItemRepo) FindLabelByCodeAndValue(ctx context.Context, code, value string) (string, error) {
	m := r.data.Query(ctx).BaseDictItem
	dict := r.data.Query(ctx).BaseDict
	q := m.WithContext(ctx).Select(m.Label).Join(dict, m.DictID.EqCol(dict.ID))
	q = q.Where(dict.Code.Eq(code))
	q = q.Where(m.Value.Eq(value))

	var label string
	err := q.Scan(&label)
	return label, err
}
