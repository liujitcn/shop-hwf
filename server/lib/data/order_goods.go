package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/dto"
	"gitee.com/liujit/shop/server/lib/data/models"
	"time"
)

type OrderGoodsCondition struct {
	OrderId  int64
	OrderIds []int64
}

type OrderGoodsRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, orderGoods *models.OrderGoods) error
	Create(ctx context.Context, orderGoods *models.OrderGoods) error
	Find(ctx context.Context, condition *OrderGoodsCondition) (*models.OrderGoods, error)
	FindAll(ctx context.Context, condition *OrderGoodsCondition) ([]*models.OrderGoods, error)
	ListPage(ctx context.Context, page, size int64, condition *OrderGoodsCondition) ([]*models.OrderGoods, int64, error)
	Count(ctx context.Context, condition *OrderGoodsCondition) (int64, error)
	BatchCreate(ctx context.Context, list []*models.OrderGoods) error
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
	OrderGoodsStatusSummary(ctx context.Context, startCreatedAt, endCreatedAt *time.Time) ([]*dto.OrderGoodsStatusSummary, error)
	OrderGoodsSummary(ctx context.Context, top int64, startCreatedAt, endCreatedAt *time.Time) ([]*dto.OrderGoodsSummary, error)
}

type orderGoodsRepo struct {
	data *Data
}

func NewOrderGoodsRepo(data *Data) OrderGoodsRepo {
	return &orderGoodsRepo{data: data}
}

func (r *orderGoodsRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).OrderGoods
	_, err := q.WithContext(ctx).Where(q.OrderID.In(ids...)).Delete()
	return err
}

func (r *orderGoodsRepo) UpdateByID(ctx context.Context, orderGoods *models.OrderGoods) error {
	if orderGoods.ID == 0 {
		return errors.New("orderGoods can not update without id")
	}
	q := r.data.Query(ctx).OrderGoods
	_, err := q.WithContext(ctx).Updates(orderGoods)
	return err
}

func (r *orderGoodsRepo) Create(ctx context.Context, orderGoods *models.OrderGoods) error {
	q := r.data.Query(ctx).OrderGoods
	err := q.WithContext(ctx).Clauses().Create(orderGoods)
	return err
}

func (r *orderGoodsRepo) Find(ctx context.Context, condition *OrderGoodsCondition) (*models.OrderGoods, error) {
	m := r.data.Query(ctx).OrderGoods
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	return q.First()
}

func (r *orderGoodsRepo) FindAll(ctx context.Context, condition *OrderGoodsCondition) ([]*models.OrderGoods, error) {
	m := r.data.Query(ctx).OrderGoods
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	if len(condition.OrderIds) > 0 {
		q = q.Where(m.OrderID.In(condition.OrderIds...))
	}
	return q.Find()
}

func (r *orderGoodsRepo) ListPage(ctx context.Context, page, size int64, condition *OrderGoodsCondition) ([]*models.OrderGoods, int64, error) {
	m := r.data.Query(ctx).OrderGoods
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	if len(condition.OrderIds) > 0 {
		q = q.Where(m.OrderID.In(condition.OrderIds...))
	}
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *orderGoodsRepo) Count(ctx context.Context, condition *OrderGoodsCondition) (int64, error) {
	m := r.data.Query(ctx).OrderGoods
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	if len(condition.OrderIds) > 0 {
		q = q.Where(m.OrderID.In(condition.OrderIds...))
	}
	count, err := q.Count()
	return count, err
}

func (r *orderGoodsRepo) BatchCreate(ctx context.Context, list []*models.OrderGoods) error {
	m := r.data.Query(ctx).OrderGoods
	return m.WithContext(ctx).Clauses().CreateInBatches(list, 100)
}

func (r *orderGoodsRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderGoods
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}

func (r *orderGoodsRepo) OrderGoodsStatusSummary(ctx context.Context, startCreatedAt, endCreatedAt *time.Time) ([]*dto.OrderGoodsStatusSummary, error) {
	order := r.data.Query(ctx).Order
	goods := r.data.Query(ctx).Goods
	goodsCategory := r.data.Query(ctx).GoodsCategory
	m := r.data.Query(ctx).OrderGoods
	q := m.WithContext(ctx)
	q = q.Join(order, order.ID.EqCol(m.OrderID))
	q = q.Join(goods, goods.ID.EqCol(m.GoodsID))
	q = q.Join(goodsCategory, goodsCategory.ID.EqCol(goods.CategoryID))

	if startCreatedAt != nil {
		q = q.Where(order.CreatedAt.Gte(*startCreatedAt))
	}
	if endCreatedAt != nil {
		q = q.Where(order.CreatedAt.Lt(*endCreatedAt))
	}
	results := make([]*dto.OrderGoodsStatusSummary, 0)

	q.Select(goodsCategory.ParentID.As("category_id"), order.Status.As("status"), m.Num.Sum().As("goods_count")).Group(goodsCategory.ParentID, order.Status)
	err := q.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *orderGoodsRepo) OrderGoodsSummary(ctx context.Context, top int64, startCreatedAt, endCreatedAt *time.Time) ([]*dto.OrderGoodsSummary, error) {
	order := r.data.Query(ctx).Order
	m := r.data.Query(ctx).OrderGoods
	q := m.WithContext(ctx)
	q = q.Join(order, order.ID.EqCol(m.OrderID))

	if startCreatedAt != nil {
		q = q.Where(order.CreatedAt.Gte(*startCreatedAt))
	}
	if endCreatedAt != nil {
		q = q.Where(order.CreatedAt.Lt(*endCreatedAt))
	}
	results := make([]*dto.OrderGoodsSummary, 0)

	q.Select(m.GoodsID.As("goods_id"), m.Num.Sum().As("goods_count")).Group(m.GoodsID).Order(m.Num.Sum().Desc()).Limit(int(top))
	err := q.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
