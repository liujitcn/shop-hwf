package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type GoodsSkuCondition struct {
	Id       int64
	Ids      []int64
	GoodsId  int64
	SkuCode  string
	SkuCodes []string
	Status   int32
}

type GoodsSkuRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, goodsSku *models.GoodsSku) error
	Create(ctx context.Context, goodsSku *models.GoodsSku) error
	Find(ctx context.Context, condition *GoodsSkuCondition) (*models.GoodsSku, error)
	FindAll(ctx context.Context, condition *GoodsSkuCondition) ([]*models.GoodsSku, error)
	ListPage(ctx context.Context, page, size int64, condition *GoodsSkuCondition) ([]*models.GoodsSku, int64, error)
	Count(ctx context.Context, condition *GoodsSkuCondition) (int64, error)
	BatchCreate(ctx context.Context, list []*models.GoodsSku) error
	DeleteByGoodsId(ctx context.Context, goodsId int64) error
	AddSaleNum(ctx context.Context, skuCode string, saleNum int64) error
	SubSaleNum(ctx context.Context, skuCode string, saleNum int64) error
}

type goodsSkuRepo struct {
	data *Data
}

func NewGoodsSkuRepo(data *Data) GoodsSkuRepo {
	return &goodsSkuRepo{data: data}
}

func (r *goodsSkuRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).GoodsSku
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *goodsSkuRepo) UpdateByID(ctx context.Context, goodsSku *models.GoodsSku) error {
	if goodsSku.ID == 0 {
		return errors.New("goodsSku can not update without id")
	}
	q := r.data.Query(ctx).GoodsSku
	_, err := q.WithContext(ctx).Updates(goodsSku)
	return err
}

func (r *goodsSkuRepo) Create(ctx context.Context, goodsSku *models.GoodsSku) error {
	q := r.data.Query(ctx).GoodsSku
	err := q.WithContext(ctx).Clauses().Create(goodsSku)
	return err
}

func (r *goodsSkuRepo) Find(ctx context.Context, condition *GoodsSkuCondition) (*models.GoodsSku, error) {
	m := r.data.Query(ctx).GoodsSku
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
	if condition.SkuCode != "" {
		q = q.Where(m.SkuCode.Eq(condition.SkuCode))
	}
	if len(condition.SkuCodes) > 0 {
		q = q.Where(m.SkuCode.In(condition.SkuCodes...))
	}
	return q.First()
}

func (r *goodsSkuRepo) FindAll(ctx context.Context, condition *GoodsSkuCondition) ([]*models.GoodsSku, error) {
	m := r.data.Query(ctx).GoodsSku
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
	if condition.SkuCode != "" {
		q = q.Where(m.SkuCode.Like(buildLikeValue(condition.SkuCode)))
	}
	if len(condition.SkuCodes) > 0 {
		q = q.Where(m.SkuCode.In(condition.SkuCodes...))
	}
	return q.Find()
}

func (r *goodsSkuRepo) ListPage(ctx context.Context, page, size int64, condition *GoodsSkuCondition) ([]*models.GoodsSku, int64, error) {
	m := r.data.Query(ctx).GoodsSku
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
	if condition.SkuCode != "" {
		q = q.Where(m.SkuCode.Like(buildLikeValue(condition.SkuCode)))
	}
	if len(condition.SkuCodes) > 0 {
		q = q.Where(m.SkuCode.In(condition.SkuCodes...))
	}
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *goodsSkuRepo) Count(ctx context.Context, condition *GoodsSkuCondition) (int64, error) {
	m := r.data.Query(ctx).GoodsSku
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
	if condition.SkuCode != "" {
		q = q.Where(m.SkuCode.Like(buildLikeValue(condition.SkuCode)))
	}
	if len(condition.SkuCodes) > 0 {
		q = q.Where(m.SkuCode.In(condition.SkuCodes...))
	}
	count, err := q.Count()
	return count, err
}

func (r *goodsSkuRepo) BatchCreate(ctx context.Context, list []*models.GoodsSku) error {
	q := r.data.Query(ctx).GoodsSku
	err := q.WithContext(ctx).Clauses().CreateInBatches(list, 100)
	return err
}
func (r *goodsSkuRepo) DeleteByGoodsId(ctx context.Context, goodsId int64) error {
	q := r.data.Query(ctx).GoodsSku
	_, err := q.WithContext(ctx).Where(q.GoodsID.Eq(goodsId)).Delete()
	return err
}

func (r *goodsSkuRepo) AddSaleNum(ctx context.Context, skuCode string, saleNum int64) error {
	q := r.data.Query(ctx).GoodsSku
	updates := map[string]interface{}{
		"real_sale_num": q.RealSaleNum.Add(saleNum),
		"inventory":     q.Inventory.Sub(saleNum),
	}
	_, err := q.WithContext(ctx).Where(q.SkuCode.Eq(skuCode), q.Inventory.Gte(saleNum)).Updates(updates)
	return err
}

func (r *goodsSkuRepo) SubSaleNum(ctx context.Context, skuCode string, saleNum int64) error {
	q := r.data.Query(ctx).GoodsSku
	updates := map[string]interface{}{
		"real_sale_num": q.RealSaleNum.Sub(saleNum),
		"inventory":     q.Inventory.Add(saleNum),
	}
	_, err := q.WithContext(ctx).Where(q.SkuCode.Eq(skuCode), q.RealSaleNum.Gte(saleNum)).Updates(updates)
	return err
}
