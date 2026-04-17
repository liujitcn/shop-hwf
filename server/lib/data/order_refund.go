package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type OrderRefundCondition struct {
	OrderId     int64
	Status      int32
	SuccessTime string
}

type OrderRefundRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, orderRefund *models.OrderRefund) error
	Create(ctx context.Context, orderRefund *models.OrderRefund) error
	Find(ctx context.Context, condition *OrderRefundCondition) (*models.OrderRefund, error)
	FindAll(ctx context.Context, condition *OrderRefundCondition) ([]*models.OrderRefund, error)
	ListPage(ctx context.Context, page, size int64, condition *OrderRefundCondition) ([]*models.OrderRefund, int64, error)
	Count(ctx context.Context, condition *OrderRefundCondition) (int64, error)
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
}

type orderRefundRepo struct {
	data *Data
}

func NewOrderRefundRepo(data *Data) OrderRefundRepo {
	return &orderRefundRepo{data: data}
}

func (r *orderRefundRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).OrderRefund
	_, err := q.WithContext(ctx).Where(q.OrderID.In(ids...)).Delete()
	return err
}

func (r *orderRefundRepo) UpdateByID(ctx context.Context, orderRefund *models.OrderRefund) error {
	if orderRefund.ID == 0 {
		return errors.New("orderRefund can not update without id")
	}
	q := r.data.Query(ctx).OrderRefund
	_, err := q.WithContext(ctx).Updates(orderRefund)
	return err
}

func (r *orderRefundRepo) Create(ctx context.Context, orderRefund *models.OrderRefund) error {
	q := r.data.Query(ctx).OrderRefund
	err := q.WithContext(ctx).Clauses().Create(orderRefund)
	return err
}

func (r *orderRefundRepo) Find(ctx context.Context, condition *OrderRefundCondition) (*models.OrderRefund, error) {
	m := r.data.Query(ctx).OrderRefund
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.SuccessTime != "" {
		q = q.Where(m.SuccessTime.DateFormat("%Y-%m-%d").Eq(condition.SuccessTime))
	}
	return q.First()
}

func (r *orderRefundRepo) FindAll(ctx context.Context, condition *OrderRefundCondition) ([]*models.OrderRefund, error) {
	m := r.data.Query(ctx).OrderRefund
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.SuccessTime != "" {
		q = q.Where(m.SuccessTime.DateFormat("%Y-%m-%d").Eq(condition.SuccessTime))
	}
	return q.Find()
}

func (r *orderRefundRepo) ListPage(ctx context.Context, page, size int64, condition *OrderRefundCondition) ([]*models.OrderRefund, int64, error) {
	m := r.data.Query(ctx).OrderRefund
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.SuccessTime != "" {
		q = q.Where(m.SuccessTime.DateFormat("%Y-%m-%d").Eq(condition.SuccessTime))
	}
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *orderRefundRepo) Count(ctx context.Context, condition *OrderRefundCondition) (int64, error) {
	m := r.data.Query(ctx).OrderRefund
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.SuccessTime != "" {
		q = q.Where(m.SuccessTime.DateFormat("%Y-%m-%d").Eq(condition.SuccessTime))
	}
	count, err := q.Count()
	return count, err
}

func (r *orderRefundRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderRefund
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}
