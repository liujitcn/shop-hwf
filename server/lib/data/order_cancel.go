package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type OrderCancelCondition struct {
	OrderId int64
}

type OrderCancelRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, orderCancel *models.OrderCancel) error
	Create(ctx context.Context, orderCancel *models.OrderCancel) error
	Find(ctx context.Context, condition *OrderCancelCondition) (*models.OrderCancel, error)
	FindAll(ctx context.Context, condition *OrderCancelCondition) ([]*models.OrderCancel, error)
	ListPage(ctx context.Context, page, size int64, condition *OrderCancelCondition) ([]*models.OrderCancel, int64, error)
	Count(ctx context.Context, condition *OrderCancelCondition) (int64, error)
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
}

type orderCancelRepo struct {
	data *Data
}

func NewOrderCancelRepo(data *Data) OrderCancelRepo {
	return &orderCancelRepo{data: data}
}

func (r *orderCancelRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).OrderCancel
	_, err := q.WithContext(ctx).Where(q.OrderID.In(ids...)).Delete()
	return err
}

func (r *orderCancelRepo) UpdateByID(ctx context.Context, orderCancel *models.OrderCancel) error {
	if orderCancel.ID == 0 {
		return errors.New("orderCancel can not update without id")
	}
	q := r.data.Query(ctx).OrderCancel
	_, err := q.WithContext(ctx).Updates(orderCancel)
	return err
}

func (r *orderCancelRepo) Create(ctx context.Context, orderCancel *models.OrderCancel) error {
	q := r.data.Query(ctx).OrderCancel
	err := q.WithContext(ctx).Clauses().Create(orderCancel)
	return err
}

func (r *orderCancelRepo) Find(ctx context.Context, condition *OrderCancelCondition) (*models.OrderCancel, error) {
	m := r.data.Query(ctx).OrderCancel
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	return q.First()
}

func (r *orderCancelRepo) FindAll(ctx context.Context, condition *OrderCancelCondition) ([]*models.OrderCancel, error) {
	m := r.data.Query(ctx).OrderCancel
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	return q.Find()
}

func (r *orderCancelRepo) ListPage(ctx context.Context, page, size int64, condition *OrderCancelCondition) ([]*models.OrderCancel, int64, error) {
	m := r.data.Query(ctx).OrderCancel
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *orderCancelRepo) Count(ctx context.Context, condition *OrderCancelCondition) (int64, error) {
	m := r.data.Query(ctx).OrderCancel
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	count, err := q.Count()
	return count, err
}

func (r *orderCancelRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderCancel
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}
