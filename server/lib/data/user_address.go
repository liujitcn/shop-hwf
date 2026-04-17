package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type UserAddressCondition struct {
	Id int64
}

type UserAddressRepo interface {
	Delete(ctx context.Context, userId int64, ids []int64) error
	UpdateByID(ctx context.Context, userId int64, userAddress *models.UserAddress) error
	Create(ctx context.Context, userAddress *models.UserAddress) error
	Find(ctx context.Context, userId int64, condition *UserAddressCondition) (*models.UserAddress, error)
	FindAll(ctx context.Context, userId int64, condition *UserAddressCondition) ([]*models.UserAddress, error)
	ListPage(ctx context.Context, userId int64, page, size int64, condition *UserAddressCondition) ([]*models.UserAddress, int64, error)
	Count(ctx context.Context, userId int64, condition *UserAddressCondition) (int64, error)
	UpdateByUserId(ctx context.Context, userId int64, userAddress *models.UserAddress) error
}

type userAddressRepo struct {
	data *Data
}

func NewUserAddressRepo(data *Data) UserAddressRepo {
	return &userAddressRepo{data: data}
}

func (r *userAddressRepo) Delete(ctx context.Context, userId int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).UserAddress
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId), q.ID.In(ids...)).Delete()
	return err
}

func (r *userAddressRepo) UpdateByID(ctx context.Context, userId int64, userAddress *models.UserAddress) error {
	if userAddress.ID == 0 {
		return errors.New("userAddress can not update without id")
	}
	q := r.data.Query(ctx).UserAddress
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId)).Updates(userAddress)
	return err
}

func (r *userAddressRepo) Create(ctx context.Context, userAddress *models.UserAddress) error {
	q := r.data.Query(ctx).UserAddress
	err := q.WithContext(ctx).Clauses().Create(userAddress)
	return err
}

func (r *userAddressRepo) Find(ctx context.Context, userId int64, condition *UserAddressCondition) (*models.UserAddress, error) {
	m := r.data.Query(ctx).UserAddress
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	return q.First()
}

func (r *userAddressRepo) FindAll(ctx context.Context, userId int64, condition *UserAddressCondition) ([]*models.UserAddress, error) {
	m := r.data.Query(ctx).UserAddress
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	q = q.Order(m.IsDefault.Desc(), m.UpdatedAt.Desc())
	return q.Find()
}

func (r *userAddressRepo) ListPage(ctx context.Context, userId int64, page, size int64, condition *UserAddressCondition) ([]*models.UserAddress, int64, error) {
	m := r.data.Query(ctx).UserAddress
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	q = q.Order(m.IsDefault.Desc(), m.UpdatedAt.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *userAddressRepo) Count(ctx context.Context, userId int64, condition *UserAddressCondition) (int64, error) {
	m := r.data.Query(ctx).UserAddress
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	return q.Count()
}

func (r *userAddressRepo) UpdateByUserId(ctx context.Context, userId int64, userAddress *models.UserAddress) error {
	q := r.data.Query(ctx).UserAddress
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId)).Updates(userAddress)
	return err
}
