package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type ShopBannerCondition struct {
	Id     int64
	Ids    []int64
	Site   int32
	Type   int32
	Status int32
}

type ShopBannerRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, shopBanner *models.ShopBanner) error
	Create(ctx context.Context, shopBanner *models.ShopBanner) error
	Find(ctx context.Context, condition *ShopBannerCondition) (*models.ShopBanner, error)
	FindAll(ctx context.Context, condition *ShopBannerCondition) ([]*models.ShopBanner, error)
	ListPage(ctx context.Context, page, size int64, condition *ShopBannerCondition) ([]*models.ShopBanner, int64, error)
	Count(ctx context.Context, condition *ShopBannerCondition) (int64, error)
}

type shopBannerRepo struct {
	data *Data
}

func NewShopBannerRepo(data *Data) ShopBannerRepo {
	return &shopBannerRepo{data: data}
}

func (r *shopBannerRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).ShopBanner
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *shopBannerRepo) UpdateByID(ctx context.Context, shopBanner *models.ShopBanner) error {
	if shopBanner.ID == 0 {
		return errors.New("shopBanner can not update without id")
	}
	q := r.data.Query(ctx).ShopBanner
	_, err := q.WithContext(ctx).Updates(shopBanner)
	return err
}

func (r *shopBannerRepo) Create(ctx context.Context, shopBanner *models.ShopBanner) error {
	q := r.data.Query(ctx).ShopBanner
	err := q.WithContext(ctx).Clauses().Create(shopBanner)
	return err
}

func (r *shopBannerRepo) Find(ctx context.Context, condition *ShopBannerCondition) (*models.ShopBanner, error) {
	m := r.data.Query(ctx).ShopBanner
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Site > 0 {
		q = q.Where(m.Site.Eq(condition.Site))
	}
	if condition.Type > 0 {
		q = q.Where(m.Type.Eq(condition.Type))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

func (r *shopBannerRepo) FindAll(ctx context.Context, condition *ShopBannerCondition) ([]*models.ShopBanner, error) {
	m := r.data.Query(ctx).ShopBanner
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Site > 0 {
		q = q.Where(m.Site.Eq(condition.Site))
	}
	if condition.Type > 0 {
		q = q.Where(m.Type.Eq(condition.Type))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.Sort.Asc())
	return q.Find()
}

func (r *shopBannerRepo) ListPage(ctx context.Context, page, size int64, condition *ShopBannerCondition) ([]*models.ShopBanner, int64, error) {
	m := r.data.Query(ctx).ShopBanner
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Site > 0 {
		q = q.Where(m.Site.Eq(condition.Site))
	}
	if condition.Type > 0 {
		q = q.Where(m.Type.Eq(condition.Type))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.Sort.Asc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *shopBannerRepo) Count(ctx context.Context, condition *ShopBannerCondition) (int64, error) {
	m := r.data.Query(ctx).ShopBanner
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Site > 0 {
		q = q.Where(m.Site.Eq(condition.Site))
	}
	if condition.Type > 0 {
		q = q.Where(m.Type.Eq(condition.Type))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	count, err := q.Count()
	return count, err
}
