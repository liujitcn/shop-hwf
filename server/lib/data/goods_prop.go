package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type GoodsPropCondition struct {
	Id      int64
	Ids     []int64
	GoodsId int64
	Label   string
	Status  int32
}

type GoodsPropRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, goodsProp *models.GoodsProp) error
	Create(ctx context.Context, goodsProp *models.GoodsProp) error
	Find(ctx context.Context, condition *GoodsPropCondition) (*models.GoodsProp, error)
	FindAll(ctx context.Context, condition *GoodsPropCondition) ([]*models.GoodsProp, error)
	ListPage(ctx context.Context, page, size int64, condition *GoodsPropCondition) ([]*models.GoodsProp, int64, error)
	Count(ctx context.Context, condition *GoodsPropCondition) (int64, error)
	BatchCreate(ctx context.Context, list []*models.GoodsProp) error
	DeleteByGoodsId(ctx context.Context, goodsId int64) error
}

type goodsPropRepo struct {
	data *Data
}

func NewGoodsPropRepo(data *Data) GoodsPropRepo {
	return &goodsPropRepo{data: data}
}

func (r *goodsPropRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).GoodsProp
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *goodsPropRepo) UpdateByID(ctx context.Context, goodsProp *models.GoodsProp) error {
	if goodsProp.ID == 0 {
		return errors.New("goodsProp can not update without id")
	}
	q := r.data.Query(ctx).GoodsProp
	_, err := q.WithContext(ctx).Updates(goodsProp)
	return err
}

func (r *goodsPropRepo) Create(ctx context.Context, goodsProp *models.GoodsProp) error {
	q := r.data.Query(ctx).GoodsProp
	err := q.WithContext(ctx).Clauses().Create(goodsProp)
	return err
}

func (r *goodsPropRepo) Find(ctx context.Context, condition *GoodsPropCondition) (*models.GoodsProp, error) {
	m := r.data.Query(ctx).GoodsProp
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Eq(condition.Label))
	}
	return q.First()
}

func (r *goodsPropRepo) FindAll(ctx context.Context, condition *GoodsPropCondition) ([]*models.GoodsProp, error) {
	m := r.data.Query(ctx).GoodsProp
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Like(buildLikeValue(condition.Label)))
	}
	q = q.Order(m.Sort.Asc())
	return q.Find()
}

func (r *goodsPropRepo) ListPage(ctx context.Context, page, size int64, condition *GoodsPropCondition) ([]*models.GoodsProp, int64, error) {
	m := r.data.Query(ctx).GoodsProp
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Like(buildLikeValue(condition.Label)))
	}
	q = q.Order(m.Sort.Asc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *goodsPropRepo) Count(ctx context.Context, condition *GoodsPropCondition) (int64, error) {
	m := r.data.Query(ctx).GoodsProp
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Like(buildLikeValue(condition.Label)))
	}
	count, err := q.Count()
	return count, err
}

func (r *goodsPropRepo) BatchCreate(ctx context.Context, list []*models.GoodsProp) error {
	m := r.data.Query(ctx).GoodsProp
	return m.WithContext(ctx).Clauses().CreateInBatches(list, 100)
}

func (r *goodsPropRepo) DeleteByGoodsId(ctx context.Context, goodsId int64) error {
	q := r.data.Query(ctx).GoodsProp
	_, err := q.WithContext(ctx).Where(q.GoodsID.Eq(goodsId)).Delete()
	return err
}
