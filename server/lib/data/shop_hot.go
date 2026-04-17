package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type ShopHotCondition struct {
	Id     int64
	Ids    []int64
	Status int32
	Title  string // 热门推荐标题
	Desc   string // 热门推荐描述
}

type ShopHotRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, shopHot *models.ShopHot) error
	Create(ctx context.Context, shopHot *models.ShopHot) error
	Find(ctx context.Context, condition *ShopHotCondition) (*models.ShopHot, error)
	FindAll(ctx context.Context, condition *ShopHotCondition) ([]*models.ShopHot, error)
	ListPage(ctx context.Context, page, size int64, condition *ShopHotCondition) ([]*models.ShopHot, int64, error)
	Count(ctx context.Context, condition *ShopHotCondition) (int64, error)
}

type shopHotRepo struct {
	data *Data
}

func NewShopHotRepo(data *Data) ShopHotRepo {
	return &shopHotRepo{data: data}
}

func (r *shopHotRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).ShopHot
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *shopHotRepo) UpdateByID(ctx context.Context, shopHot *models.ShopHot) error {
	if shopHot.ID == 0 {
		return errors.New("shopHot can not update without id")
	}
	q := r.data.Query(ctx).ShopHot
	_, err := q.WithContext(ctx).Updates(shopHot)
	return err
}

func (r *shopHotRepo) Create(ctx context.Context, shopHot *models.ShopHot) error {
	q := r.data.Query(ctx).ShopHot
	err := q.WithContext(ctx).Clauses().Create(shopHot)
	return err
}

func (r *shopHotRepo) Find(ctx context.Context, condition *ShopHotCondition) (*models.ShopHot, error) {
	m := r.data.Query(ctx).ShopHot
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Title != "" {
		q = q.Where(m.Title.Eq(condition.Title))
	}
	if condition.Desc != "" {
		q = q.Where(m.Desc.Eq(condition.Desc))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

func (r *shopHotRepo) FindAll(ctx context.Context, condition *ShopHotCondition) ([]*models.ShopHot, error) {
	m := r.data.Query(ctx).ShopHot
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Title != "" {
		q = q.Where(m.Title.Like(buildLikeValue(condition.Title)))
	}
	if condition.Desc != "" {
		q = q.Where(m.Desc.Like(buildLikeValue(condition.Desc)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.Sort.Asc())
	return q.Find()
}

func (r *shopHotRepo) ListPage(ctx context.Context, page, size int64, condition *ShopHotCondition) ([]*models.ShopHot, int64, error) {
	m := r.data.Query(ctx).ShopHot
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Title != "" {
		q = q.Where(m.Title.Like(buildLikeValue(condition.Title)))
	}
	if condition.Desc != "" {
		q = q.Where(m.Desc.Like(buildLikeValue(condition.Desc)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.Sort.Asc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *shopHotRepo) Count(ctx context.Context, condition *ShopHotCondition) (int64, error) {
	m := r.data.Query(ctx).ShopHot
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Title != "" {
		q = q.Where(m.Title.Like(buildLikeValue(condition.Title)))
	}
	if condition.Desc != "" {
		q = q.Where(m.Desc.Like(buildLikeValue(condition.Desc)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	count, err := q.Count()
	return count, err
}
