package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type GoodsSpecCondition struct {
	Id      int64
	Ids     []int64
	GoodsId int64
	Name    string
	Status  int32
}

type GoodsSpecRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, goodsSpec *models.GoodsSpec) error
	Create(ctx context.Context, goodsSpec *models.GoodsSpec) error
	Find(ctx context.Context, condition *GoodsSpecCondition) (*models.GoodsSpec, error)
	FindAll(ctx context.Context, condition *GoodsSpecCondition) ([]*models.GoodsSpec, error)
	ListPage(ctx context.Context, page, size int64, condition *GoodsSpecCondition) ([]*models.GoodsSpec, int64, error)
	Count(ctx context.Context, condition *GoodsSpecCondition) (int64, error)
	BatchCreate(ctx context.Context, list []*models.GoodsSpec) error
	DeleteByGoodsId(ctx context.Context, goodsId int64) error
}

type goodsSpecRepo struct {
	data *Data
}

func NewGoodsSpecRepo(data *Data) GoodsSpecRepo {
	return &goodsSpecRepo{data: data}
}

func (r *goodsSpecRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).GoodsSpec
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *goodsSpecRepo) UpdateByID(ctx context.Context, goodsSpec *models.GoodsSpec) error {
	if goodsSpec.ID == 0 {
		return errors.New("goodsSpec can not update without id")
	}
	q := r.data.Query(ctx).GoodsSpec
	_, err := q.WithContext(ctx).Updates(goodsSpec)
	return err
}

func (r *goodsSpecRepo) Create(ctx context.Context, goodsSpec *models.GoodsSpec) error {
	q := r.data.Query(ctx).GoodsSpec
	err := q.WithContext(ctx).Clauses().Create(goodsSpec)
	return err
}

func (r *goodsSpecRepo) Find(ctx context.Context, condition *GoodsSpecCondition) (*models.GoodsSpec, error) {
	m := r.data.Query(ctx).GoodsSpec
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
	if condition.Name != "" {
		q = q.Where(m.Name.Eq(condition.Name))
	}
	return q.First()
}

func (r *goodsSpecRepo) FindAll(ctx context.Context, condition *GoodsSpecCondition) ([]*models.GoodsSpec, error) {
	m := r.data.Query(ctx).GoodsSpec
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
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	q = q.Order(m.Sort.Asc())
	return q.Find()
}

func (r *goodsSpecRepo) ListPage(ctx context.Context, page, size int64, condition *GoodsSpecCondition) ([]*models.GoodsSpec, int64, error) {
	m := r.data.Query(ctx).GoodsSpec
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
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	q = q.Order(m.Sort.Asc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *goodsSpecRepo) Count(ctx context.Context, condition *GoodsSpecCondition) (int64, error) {
	m := r.data.Query(ctx).GoodsSpec
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
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	count, err := q.Count()
	return count, err
}

func (r *goodsSpecRepo) BatchCreate(ctx context.Context, list []*models.GoodsSpec) error {
	q := r.data.Query(ctx).GoodsSpec
	err := q.WithContext(ctx).Clauses().CreateInBatches(list, 100)
	return err
}
func (r *goodsSpecRepo) DeleteByGoodsId(ctx context.Context, goodsId int64) error {
	q := r.data.Query(ctx).GoodsSpec
	_, err := q.WithContext(ctx).Where(q.GoodsID.Eq(goodsId)).Delete()
	return err
}
