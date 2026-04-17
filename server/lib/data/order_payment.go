package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type OrderPaymentCondition struct {
	OrderId     int64
	Status      int32
	SuccessTime string
}

type OrderPaymentRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, orderPayment *models.OrderPayment) error
	Create(ctx context.Context, orderPayment *models.OrderPayment) error
	Find(ctx context.Context, condition *OrderPaymentCondition) (*models.OrderPayment, error)
	FindAll(ctx context.Context, condition *OrderPaymentCondition) ([]*models.OrderPayment, error)
	ListPage(ctx context.Context, page, size int64, condition *OrderPaymentCondition) ([]*models.OrderPayment, int64, error)
	Count(ctx context.Context, condition *OrderPaymentCondition) (int64, error)
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
}

type orderPaymentRepo struct {
	data *Data
}

func NewOrderPaymentRepo(data *Data) OrderPaymentRepo {
	return &orderPaymentRepo{data: data}
}

func (r *orderPaymentRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).OrderPayment
	_, err := q.WithContext(ctx).Where(q.OrderID.In(ids...)).Delete()
	return err
}

func (r *orderPaymentRepo) UpdateByID(ctx context.Context, orderPayment *models.OrderPayment) error {
	if orderPayment.ID == 0 {
		return errors.New("orderPayment can not update without id")
	}
	q := r.data.Query(ctx).OrderPayment
	_, err := q.WithContext(ctx).Updates(orderPayment)
	return err
}

func (r *orderPaymentRepo) Create(ctx context.Context, orderPayment *models.OrderPayment) error {
	q := r.data.Query(ctx).OrderPayment
	err := q.WithContext(ctx).Clauses().Create(orderPayment)
	return err
}

func (r *orderPaymentRepo) Find(ctx context.Context, condition *OrderPaymentCondition) (*models.OrderPayment, error) {
	m := r.data.Query(ctx).OrderPayment
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

func (r *orderPaymentRepo) FindAll(ctx context.Context, condition *OrderPaymentCondition) ([]*models.OrderPayment, error) {
	m := r.data.Query(ctx).OrderPayment
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

func (r *orderPaymentRepo) ListPage(ctx context.Context, page, size int64, condition *OrderPaymentCondition) ([]*models.OrderPayment, int64, error) {
	m := r.data.Query(ctx).OrderPayment
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

func (r *orderPaymentRepo) Count(ctx context.Context, condition *OrderPaymentCondition) (int64, error) {
	m := r.data.Query(ctx).OrderPayment
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

func (r *orderPaymentRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderPayment
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}
