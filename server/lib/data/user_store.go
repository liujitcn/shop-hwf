package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type UserStoreCondition struct {
	Id     int64
	Name   string
	UserId int64
	Status int32
}

type UserStoreRepo interface {
	Delete(ctx context.Context, userId int64, ids []int64) error
	UpdateByID(ctx context.Context, userId int64, userStore *models.UserStore) error
	Create(ctx context.Context, userStore *models.UserStore) error
	Find(ctx context.Context, condition *UserStoreCondition) (*models.UserStore, error)
	FindAll(ctx context.Context, condition *UserStoreCondition) ([]*models.UserStore, error)
	ListPage(ctx context.Context, page, size int64, condition *UserStoreCondition) ([]*models.UserStore, int64, error)
	Count(ctx context.Context, condition *UserStoreCondition) (int64, error)
}

type userStoreRepo struct {
	data *Data
}

func NewUserStoreRepo(data *Data) UserStoreRepo {
	return &userStoreRepo{data: data}
}

func (r *userStoreRepo) Delete(ctx context.Context, userId int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).UserStore
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId), q.ID.In(ids...)).Delete()
	return err
}

func (r *userStoreRepo) UpdateByID(ctx context.Context, userId int64, userStore *models.UserStore) error {
	if userStore.ID == 0 {
		return errors.New("userStore can not update without id")
	}
	m := r.data.Query(ctx).UserStore
	q := m.WithContext(ctx)
	if userId > 0 {
		q = q.Where(m.UserID.Eq(userId))
	}
	_, err := q.Updates(userStore)
	return err
}

func (r *userStoreRepo) Create(ctx context.Context, userStore *models.UserStore) error {
	q := r.data.Query(ctx).UserStore
	err := q.WithContext(ctx).Clauses().Create(userStore)
	return err
}

func (r *userStoreRepo) Find(ctx context.Context, condition *UserStoreCondition) (*models.UserStore, error) {
	m := r.data.Query(ctx).UserStore
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Eq(condition.Name))
	}
	if condition.UserId > 0 {
		q = q.Where(m.UserID.Eq(condition.UserId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

func (r *userStoreRepo) FindAll(ctx context.Context, condition *UserStoreCondition) ([]*models.UserStore, error) {
	m := r.data.Query(ctx).UserStore
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.UserId > 0 {
		q = q.Where(m.UserID.Eq(condition.UserId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.CreatedAt.Desc())
	return q.Find()
}

func (r *userStoreRepo) ListPage(ctx context.Context, page, size int64, condition *UserStoreCondition) ([]*models.UserStore, int64, error) {
	m := r.data.Query(ctx).UserStore
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.UserId > 0 {
		q = q.Where(m.UserID.Eq(condition.UserId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.CreatedAt.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *userStoreRepo) Count(ctx context.Context, condition *UserStoreCondition) (int64, error) {
	m := r.data.Query(ctx).UserStore
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.UserId > 0 {
		q = q.Where(m.UserID.Eq(condition.UserId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.Count()
}
