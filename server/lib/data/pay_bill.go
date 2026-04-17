package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type PayBillCondition struct {
	Id       int64
	BillDate string
	BillType string
}

type PayBillRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, payBill *models.PayBill) error
	Create(ctx context.Context, payBill *models.PayBill) error
	Find(ctx context.Context, condition *PayBillCondition) (*models.PayBill, error)
	FindAll(ctx context.Context, condition *PayBillCondition) ([]*models.PayBill, error)
	ListPage(ctx context.Context, page, size int64, condition *PayBillCondition) ([]*models.PayBill, int64, error)
	Count(ctx context.Context, condition *PayBillCondition) (int64, error)
}

type payBillRepo struct {
	data *Data
}

func NewPayBillRepo(data *Data) PayBillRepo {
	return &payBillRepo{data: data}
}

func (r *payBillRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).PayBill
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *payBillRepo) UpdateByID(ctx context.Context, payBill *models.PayBill) error {
	if payBill.ID == 0 {
		return errors.New("payBill can not update without id")
	}
	q := r.data.Query(ctx).PayBill
	_, err := q.WithContext(ctx).Updates(payBill)
	return err
}

func (r *payBillRepo) Create(ctx context.Context, payBill *models.PayBill) error {
	q := r.data.Query(ctx).PayBill
	err := q.WithContext(ctx).Clauses().Create(payBill)
	return err
}

func (r *payBillRepo) Find(ctx context.Context, condition *PayBillCondition) (*models.PayBill, error) {
	m := r.data.Query(ctx).PayBill
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.BillDate != "" {
		q = q.Where(m.BillDate.Eq(condition.BillDate))
	}
	if condition.BillType != "" {
		q = q.Where(m.BillType.Eq(condition.BillType))
	}
	return q.First()
}

func (r *payBillRepo) FindAll(ctx context.Context, condition *PayBillCondition) ([]*models.PayBill, error) {
	m := r.data.Query(ctx).PayBill
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.BillDate != "" {
		q = q.Where(m.BillDate.Like(buildLikeValue(condition.BillDate)))
	}
	if condition.BillType != "" {
		q = q.Where(m.BillType.Eq(condition.BillType))
	}
	q = q.Order(m.BillDate.Desc())
	return q.Find()
}

func (r *payBillRepo) ListPage(ctx context.Context, page, size int64, condition *PayBillCondition) ([]*models.PayBill, int64, error) {
	m := r.data.Query(ctx).PayBill
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.BillDate != "" {
		q = q.Where(m.BillDate.Like(buildLikeValue(condition.BillDate)))
	}
	if condition.BillType != "" {
		q = q.Where(m.BillType.Eq(condition.BillType))
	}
	q = q.Order(m.BillDate.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *payBillRepo) Count(ctx context.Context, condition *PayBillCondition) (int64, error) {
	m := r.data.Query(ctx).PayBill
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.BillDate != "" {
		q = q.Where(m.BillDate.Eq(condition.BillDate))
	}
	if condition.BillType != "" {
		q = q.Where(m.BillType.Eq(condition.BillType))
	}
	count, err := q.Count()
	return count, err
}
