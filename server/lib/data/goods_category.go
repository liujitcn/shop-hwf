package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type GoodsCategoryCondition struct {
	Id       int64
	Ids      []int64
	Name     string
	ParentId *int64
	Status   int32
}

type GoodsCategoryRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, goodsCategory *models.GoodsCategory) error
	Create(ctx context.Context, goodsCategory *models.GoodsCategory) error
	Find(ctx context.Context, condition *GoodsCategoryCondition) (*models.GoodsCategory, error)
	FindAll(ctx context.Context, condition *GoodsCategoryCondition) ([]*models.GoodsCategory, error)
	ListPage(ctx context.Context, page, size int64, condition *GoodsCategoryCondition) ([]*models.GoodsCategory, int64, error)
	Count(ctx context.Context, condition *GoodsCategoryCondition) (int64, error)
}

type goodsCategoryRepo struct {
	data *Data
}

func NewGoodsCategoryRepo(data *Data) GoodsCategoryRepo {
	return &goodsCategoryRepo{data: data}
}

func (r *goodsCategoryRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).GoodsCategory
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *goodsCategoryRepo) UpdateByID(ctx context.Context, goodsCategory *models.GoodsCategory) error {
	if goodsCategory.ID == 0 {
		return errors.New("goodsCategory can not update without id")
	}
	q := r.data.Query(ctx).GoodsCategory
	_, err := q.WithContext(ctx).Updates(goodsCategory)
	return err
}

func (r *goodsCategoryRepo) Create(ctx context.Context, goodsCategory *models.GoodsCategory) error {
	q := r.data.Query(ctx).GoodsCategory
	err := q.WithContext(ctx).Clauses().Create(goodsCategory)
	return err
}

func (r *goodsCategoryRepo) Find(ctx context.Context, condition *GoodsCategoryCondition) (*models.GoodsCategory, error) {
	m := r.data.Query(ctx).GoodsCategory
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Eq(condition.Name))
	}
	if condition.ParentId != nil {
		q = q.Where(m.ParentID.Eq(*condition.ParentId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

func (r *goodsCategoryRepo) FindAll(ctx context.Context, condition *GoodsCategoryCondition) ([]*models.GoodsCategory, error) {
	m := r.data.Query(ctx).GoodsCategory
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
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

func (r *goodsCategoryRepo) ListPage(ctx context.Context, page, size int64, condition *GoodsCategoryCondition) ([]*models.GoodsCategory, int64, error) {
	m := r.data.Query(ctx).GoodsCategory
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
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

func (r *goodsCategoryRepo) Count(ctx context.Context, condition *GoodsCategoryCondition) (int64, error) {
	m := r.data.Query(ctx).GoodsCategory
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
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
