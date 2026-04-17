package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type OrderAddressCondition struct {
	OrderId int64
}

type OrderAddressRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, orderAddress *models.OrderAddress) error
	Create(ctx context.Context, orderAddress *models.OrderAddress) error
	Find(ctx context.Context, condition *OrderAddressCondition) (*models.OrderAddress, error)
	FindAll(ctx context.Context, condition *OrderAddressCondition) ([]*models.OrderAddress, error)
	ListPage(ctx context.Context, page, size int64, condition *OrderAddressCondition) ([]*models.OrderAddress, int64, error)
	Count(ctx context.Context, condition *OrderAddressCondition) (int64, error)
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
}

type orderAddressRepo struct {
	data *Data
}

func NewOrderAddressRepo(data *Data) OrderAddressRepo {
	return &orderAddressRepo{data: data}
}

func (r *orderAddressRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).OrderAddress
	_, err := q.WithContext(ctx).Where(q.OrderID.In(ids...)).Delete()
	return err
}

func (r *orderAddressRepo) UpdateByID(ctx context.Context, orderAddress *models.OrderAddress) error {
	if orderAddress.ID == 0 {
		return errors.New("orderAddress can not update without id")
	}
	q := r.data.Query(ctx).OrderAddress
	_, err := q.WithContext(ctx).Updates(orderAddress)
	return err
}

func (r *orderAddressRepo) Create(ctx context.Context, orderAddress *models.OrderAddress) error {
	q := r.data.Query(ctx).OrderAddress
	err := q.WithContext(ctx).Clauses().Create(orderAddress)
	return err
}

func (r *orderAddressRepo) Find(ctx context.Context, condition *OrderAddressCondition) (*models.OrderAddress, error) {
	m := r.data.Query(ctx).OrderAddress
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	return q.First()
}

func (r *orderAddressRepo) FindAll(ctx context.Context, condition *OrderAddressCondition) ([]*models.OrderAddress, error) {
	m := r.data.Query(ctx).OrderAddress
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	return q.Find()
}

func (r *orderAddressRepo) ListPage(ctx context.Context, page, size int64, condition *OrderAddressCondition) ([]*models.OrderAddress, int64, error) {
	m := r.data.Query(ctx).OrderAddress
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *orderAddressRepo) Count(ctx context.Context, condition *OrderAddressCondition) (int64, error) {
	m := r.data.Query(ctx).OrderAddress
	q := m.WithContext(ctx)
	if condition.OrderId > 0 {
		q = q.Where(m.OrderID.Eq(condition.OrderId))
	}
	count, err := q.Count()
	return count, err
}

func (r *orderAddressRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderAddress
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}
