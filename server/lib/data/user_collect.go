package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type UserCollectCondition struct {
	Id      int64
	GoodsId int64
}

type UserCollectRepo interface {
	Delete(ctx context.Context, userId int64, ids []int64) error
	UpdateByID(ctx context.Context, userId int64, userCollect *models.UserCollect) error
	Create(ctx context.Context, userCollect *models.UserCollect) error
	Find(ctx context.Context, userId int64, condition *UserCollectCondition) (*models.UserCollect, error)
	FindAll(ctx context.Context, userId int64, condition *UserCollectCondition) ([]*models.UserCollect, error)
	ListPage(ctx context.Context, userId int64, page, size int64, condition *UserCollectCondition) ([]*models.UserCollect, int64, error)
	Count(ctx context.Context, userId int64, condition *UserCollectCondition) (int64, error)
}

type userCollectRepo struct {
	data *Data
}

func NewUserCollectRepo(data *Data) UserCollectRepo {
	return &userCollectRepo{data: data}
}

func (r *userCollectRepo) Delete(ctx context.Context, userId int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).UserCollect
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId), q.ID.In(ids...)).Delete()
	return err
}

func (r *userCollectRepo) UpdateByID(ctx context.Context, userId int64, userCollect *models.UserCollect) error {
	if userCollect.ID == 0 {
		return errors.New("userCollect can not update without id")
	}
	q := r.data.Query(ctx).UserCollect
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId)).Updates(userCollect)
	return err
}

func (r *userCollectRepo) Create(ctx context.Context, userCollect *models.UserCollect) error {
	q := r.data.Query(ctx).UserCollect
	err := q.WithContext(ctx).Clauses().Create(userCollect)
	return err
}

func (r *userCollectRepo) Find(ctx context.Context, userId int64, condition *UserCollectCondition) (*models.UserCollect, error) {
	m := r.data.Query(ctx).UserCollect
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	return q.First()
}

func (r *userCollectRepo) FindAll(ctx context.Context, userId int64, condition *UserCollectCondition) ([]*models.UserCollect, error) {
	m := r.data.Query(ctx).UserCollect
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	q = q.Order(m.CreatedAt.Desc())
	return q.Find()
}

func (r *userCollectRepo) ListPage(ctx context.Context, userId int64, page, size int64, condition *UserCollectCondition) ([]*models.UserCollect, int64, error) {
	m := r.data.Query(ctx).UserCollect
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	q = q.Order(m.CreatedAt.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *userCollectRepo) Count(ctx context.Context, userId int64, condition *UserCollectCondition) (int64, error) {
	m := r.data.Query(ctx).UserCollect
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	return q.Count()
}
