package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type OrderLogisticsCondition struct {
	OrderId int64
}

type OrderLogisticsRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, orderLogistics *models.OrderLogistics) error
	Create(ctx context.Context, orderLogistics *models.OrderLogistics) error
	Find(ctx context.Context, condition *OrderLogisticsCondition) (*models.OrderLogistics, error)
	FindAll(ctx context.Context, condition *OrderLogisticsCondition) ([]*models.OrderLogistics, error)
	ListPage(ctx context.Context, page, size int64, condition *OrderLogisticsCondition) ([]*models.OrderLogistics, int64, error)
	Count(ctx context.Context, condition *OrderLogisticsCondition) (int64, error)
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
}

type orderLogisticsRepo struct {
	data *Data
}

func NewOrderLogisticsRepo(data *Data) OrderLogisticsRepo {
	return &orderLogisticsRepo{data: data}
}

func (r *orderLogisticsRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).OrderLogistics
	_, err := q.WithContext(ctx).Where(q.OrderID.In(ids...)).Delete()
	return err
}

func (r *orderLogisticsRepo) UpdateByID(ctx context.Context, orderLogistics *models.OrderLogistics) error {
	if orderLogistics.ID == 0 {
		return errors.New("orderLogistics can not update without id")
	}
	q := r.data.Query(ctx).OrderLogistics
	_, err := q.WithContext(ctx).Updates(orderLogistics)
	return err
}

func (r *orderLogisticsRepo) Create(ctx context.Context, orderLogistics *models.OrderLogistics) error {
	q := r.data.Query(ctx).OrderLogistics
	err := q.WithContext(ctx).Clauses().Create(orderLogistics)
	return err
}

func (r *orderLogisticsRepo) Find(ctx context.Context, condition *OrderLogisticsCondition) (*models.OrderLogistics, error) {
	m := r.data.Query(ctx).OrderLogistics
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	return q.First()
}

func (r *orderLogisticsRepo) FindAll(ctx context.Context, condition *OrderLogisticsCondition) ([]*models.OrderLogistics, error) {
	m := r.data.Query(ctx).OrderLogistics
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	return q.Find()
}

func (r *orderLogisticsRepo) ListPage(ctx context.Context, page, size int64, condition *OrderLogisticsCondition) ([]*models.OrderLogistics, int64, error) {
	m := r.data.Query(ctx).OrderLogistics
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *orderLogisticsRepo) Count(ctx context.Context, condition *OrderLogisticsCondition) (int64, error) {
	m := r.data.Query(ctx).OrderLogistics
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	count, err := q.Count()
	return count, err
}

func (r *orderLogisticsRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderLogistics
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}
