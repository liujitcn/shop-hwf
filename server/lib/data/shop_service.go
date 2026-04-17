package data

import (
	"context"
	"errors"

	"gitee.com/liujit/shop/server/lib/data/models"
	"gorm.io/gorm"
)

// ShopServiceCondition 商城服务查询条件
type ShopServiceCondition struct {
	Id     int64
	Ids    []int64
	Label  string
	Status int32
}

// ShopServiceRepo 商城服务数据接口
type ShopServiceRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, shopService *models.ShopService) error
	Create(ctx context.Context, shopService *models.ShopService) error
	Find(ctx context.Context, condition *ShopServiceCondition) (*models.ShopService, error)
	FindAll(ctx context.Context, condition *ShopServiceCondition) ([]*models.ShopService, error)
	ListPage(ctx context.Context, page, size int64, condition *ShopServiceCondition) ([]*models.ShopService, int64, error)
	Count(ctx context.Context, condition *ShopServiceCondition) (int64, error)
}

// shopServiceRepo 商城服务数据实现
type shopServiceRepo struct {
	data *Data
	db   *gorm.DB
}

// NewShopServiceRepo 创建商城服务数据实例
func NewShopServiceRepo(data *Data) ShopServiceRepo {
	return &shopServiceRepo{
		data: data,
		db:   data.db,
	}
}

// Delete 删除商城服务
func (r *shopServiceRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).ShopService
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

// UpdateByID 更新商城服务
func (r *shopServiceRepo) UpdateByID(ctx context.Context, shopService *models.ShopService) error {
	if shopService.ID == 0 {
		return errors.New("shopService can not update without id")
	}
	q := r.data.Query(ctx).ShopService
	_, err := q.WithContext(ctx).Updates(shopService)
	return err
}

// Create 创建商城服务
func (r *shopServiceRepo) Create(ctx context.Context, shopService *models.ShopService) error {
	q := r.data.Query(ctx).ShopService
	err := q.WithContext(ctx).Clauses().Create(shopService)
	return err
}

// Find 查询商城服务
func (r *shopServiceRepo) Find(ctx context.Context, condition *ShopServiceCondition) (*models.ShopService, error) {
	m := r.data.Query(ctx).ShopService
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Like(buildLikeValue(condition.Label)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

// FindAll 查询商城服务列表
func (r *shopServiceRepo) FindAll(ctx context.Context, condition *ShopServiceCondition) ([]*models.ShopService, error) {
	m := r.data.Query(ctx).ShopService
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Like(buildLikeValue(condition.Label)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.Sort.Asc())
	return q.Find()
}

// ListPage 分页查询商城服务列表
func (r *shopServiceRepo) ListPage(ctx context.Context, page, size int64, condition *ShopServiceCondition) ([]*models.ShopService, int64, error) {
	m := r.data.Query(ctx).ShopService
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Like(buildLikeValue(condition.Label)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.Sort.Asc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

// Count 查询商城服务数量
func (r *shopServiceRepo) Count(ctx context.Context, condition *ShopServiceCondition) (int64, error) {
	m := r.data.Query(ctx).ShopService
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.Label != "" {
		q = q.Where(m.Label.Like(buildLikeValue(condition.Label)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	count, err := q.Count()
	return count, err
}
