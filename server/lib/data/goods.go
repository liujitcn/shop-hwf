package data

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/liujit/shop/server/lib/data/dto"
	"gitee.com/liujit/shop/server/lib/data/models"
	"time"
)

type GoodsCondition struct {
	Id             int64
	Ids            []int64
	Name           string
	CategoryId     int64
	CategoryPath   string
	Status         int32
	StartCreatedAt *time.Time // 创建开始时间
	EndCreatedAt   *time.Time // 创建结束时间
}

type GoodsRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, goods *models.Goods) error
	Create(ctx context.Context, goods *models.Goods) error
	Find(ctx context.Context, condition *GoodsCondition) (*models.Goods, error)
	FindAll(ctx context.Context, condition *GoodsCondition) ([]*models.Goods, error)
	ListPage(ctx context.Context, page, size int64, condition *GoodsCondition) ([]*models.Goods, int64, error)
	Count(ctx context.Context, condition *GoodsCondition) (int64, error)
	AddSaleNum(ctx context.Context, id, saleNum int64) error
	SubSaleNum(ctx context.Context, id, saleNum int64) error
	GoodsCategorySummary(ctx context.Context) ([]*dto.GoodsCategorySummary, error)
}

type goodsRepo struct {
	data *Data
}

func NewGoodsRepo(data *Data) GoodsRepo {
	return &goodsRepo{data: data}
}

func (r *goodsRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).Goods
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *goodsRepo) UpdateByID(ctx context.Context, goods *models.Goods) error {
	if goods.ID == 0 {
		return errors.New("goods can not update without id")
	}
	q := r.data.Query(ctx).Goods
	_, err := q.WithContext(ctx).Updates(goods)
	return err
}

func (r *goodsRepo) Create(ctx context.Context, goods *models.Goods) error {
	q := r.data.Query(ctx).Goods
	err := q.WithContext(ctx).Clauses().Create(goods)
	return err
}

func (r *goodsRepo) Find(ctx context.Context, condition *GoodsCondition) (*models.Goods, error) {
	m := r.data.Query(ctx).Goods
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
	if condition.CategoryId > 0 {
		q = q.Where(m.CategoryID.Eq(condition.CategoryId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	return q.First()
}

func (r *goodsRepo) FindAll(ctx context.Context, condition *GoodsCondition) ([]*models.Goods, error) {
	category := r.data.Query(ctx).GoodsCategory
	m := r.data.Query(ctx).Goods
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.CategoryPath) > 0 {
		q = q.Join(category, category.ID.EqCol(m.CategoryID), category.Path.Like(fmt.Sprintf("%s%%", condition.CategoryPath)))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.CategoryId > 0 {
		q = q.Where(m.CategoryID.Eq(condition.CategoryId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	q = q.Order(m.UpdatedBy.Desc())
	return q.Find()
}

func (r *goodsRepo) ListPage(ctx context.Context, page, size int64, condition *GoodsCondition) ([]*models.Goods, int64, error) {
	category := r.data.Query(ctx).GoodsCategory
	m := r.data.Query(ctx).Goods
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.CategoryPath) > 0 {
		q = q.Join(category, category.ID.EqCol(m.CategoryID), category.Path.Like(fmt.Sprintf("%s%%", condition.CategoryPath)))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.CategoryId > 0 {
		q = q.Where(m.CategoryID.Eq(condition.CategoryId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	q = q.Order(m.UpdatedBy.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *goodsRepo) Count(ctx context.Context, condition *GoodsCondition) (int64, error) {
	goodsCategory := r.data.Query(ctx).GoodsCategory
	m := r.data.Query(ctx).Goods
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.CategoryPath) > 0 {
		q = q.Join(goodsCategory, goodsCategory.ID.EqCol(m.CategoryID), goodsCategory.Path.Like(fmt.Sprintf("%s%%", condition.CategoryPath)))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.CategoryId > 0 {
		q = q.Where(m.CategoryID.Eq(condition.CategoryId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	count, err := q.Count()
	return count, err
}

func (r *goodsRepo) AddSaleNum(ctx context.Context, id, saleNum int64) error {
	q := r.data.Query(ctx).Goods
	_, err := q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.RealSaleNum, q.RealSaleNum.Add(saleNum))
	return err
}

func (r *goodsRepo) SubSaleNum(ctx context.Context, id, saleNum int64) error {
	q := r.data.Query(ctx).Goods
	_, err := q.WithContext(ctx).Where(q.ID.Eq(id), q.RealSaleNum.Gte(saleNum)).Update(q.RealSaleNum, q.RealSaleNum.Sub(saleNum))
	return err
}

func (r *goodsRepo) GoodsCategorySummary(ctx context.Context) ([]*dto.GoodsCategorySummary, error) {
	category := r.data.Query(ctx).GoodsCategory
	m := r.data.Query(ctx).Goods
	q := m.WithContext(ctx)
	q = q.Join(category, category.ID.EqCol(m.CategoryID))

	results := make([]*dto.GoodsCategorySummary, 0)

	q.Select(category.ParentID.As("category_id"), m.ID.Count().As("goods_count")).Group(category.ParentID)
	err := q.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
