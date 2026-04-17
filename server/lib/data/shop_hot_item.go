package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type ShopHotItemCondition struct {
	Id     int64
	HotId  int64
	Status int32
	Title  string
}

type ShopHotItemRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, shopHotItem *models.ShopHotItem) error
	Create(ctx context.Context, shopHotItem *models.ShopHotItem) error
	Find(ctx context.Context, condition *ShopHotItemCondition) (*models.ShopHotItem, error)
	FindAll(ctx context.Context, condition *ShopHotItemCondition) ([]*models.ShopHotItem, error)
	ListPage(ctx context.Context, page, size int64, condition *ShopHotItemCondition) ([]*models.ShopHotItem, int64, error)
	Count(ctx context.Context, condition *ShopHotItemCondition) (int64, error)
}

type shopHotItemRepo struct {
	data *Data
}

func NewShopHotItemRepo(data *Data) ShopHotItemRepo {
	return &shopHotItemRepo{data: data}
}

func (r *shopHotItemRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).ShopHotItem
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *shopHotItemRepo) UpdateByID(ctx context.Context, shopHotItem *models.ShopHotItem) error {
	if shopHotItem.ID == 0 {
		return errors.New("shopHotItem can not update without id")
	}
	q := r.data.Query(ctx).ShopHotItem
	_, err := q.WithContext(ctx).Updates(shopHotItem)
	return err
}

func (r *shopHotItemRepo) Create(ctx context.Context, shopHotItem *models.ShopHotItem) error {
	q := r.data.Query(ctx).ShopHotItem
	err := q.WithContext(ctx).Clauses().Create(shopHotItem)
	return err
}

func (r *shopHotItemRepo) Find(ctx context.Context, condition *ShopHotItemCondition) (*models.ShopHotItem, error) {
	m := r.data.Query(ctx).ShopHotItem
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.HotId > 0 {
		q = q.Where(m.HotID.Eq(condition.HotId))
	}
	if condition.Title != "" {
		q = q.Where(m.Title.Eq(condition.Title))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

func (r *shopHotItemRepo) FindAll(ctx context.Context, condition *ShopHotItemCondition) ([]*models.ShopHotItem, error) {
	m := r.data.Query(ctx).ShopHotItem
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.HotId > 0 {
		q = q.Where(m.HotID.Eq(condition.HotId))
	}
	if condition.Title != "" {
		q = q.Where(m.Title.Like(buildLikeValue(condition.Title)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.Sort.Asc())
	return q.Find()
}

func (r *shopHotItemRepo) ListPage(ctx context.Context, page, size int64, condition *ShopHotItemCondition) ([]*models.ShopHotItem, int64, error) {
	m := r.data.Query(ctx).ShopHotItem
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.HotId > 0 {
		q = q.Where(m.HotID.Eq(condition.HotId))
	}
	if condition.Title != "" {
		q = q.Where(m.Title.Like(buildLikeValue(condition.Title)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.Sort.Asc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *shopHotItemRepo) Count(ctx context.Context, condition *ShopHotItemCondition) (int64, error) {
	m := r.data.Query(ctx).ShopHotItem
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.HotId > 0 {
		q = q.Where(m.HotID.Eq(condition.HotId))
	}
	if condition.Title != "" {
		q = q.Where(m.Title.Like(buildLikeValue(condition.Title)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	count, err := q.Count()
	return count, err
}
