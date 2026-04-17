package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseConfigCondition struct {
	Id     int64
	Ids    []int64
	Site   int32
	Name   string
	Type   int32
	Key    string
	Status int32
}

type BaseConfigRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseConfig *models.BaseConfig) error
	Create(ctx context.Context, baseConfig *models.BaseConfig) error
	Find(ctx context.Context, condition *BaseConfigCondition) (*models.BaseConfig, error)
	FindAll(ctx context.Context, condition *BaseConfigCondition) ([]*models.BaseConfig, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseConfigCondition) ([]*models.BaseConfig, int64, error)
	Count(ctx context.Context, condition *BaseConfigCondition) (int64, error)
}

type baseConfigRepo struct {
	data *Data
}

func NewBaseConfigRepo(data *Data) BaseConfigRepo {
	return &baseConfigRepo{data: data}
}

func (r *baseConfigRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).BaseConfig
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseConfigRepo) UpdateByID(ctx context.Context, baseConfig *models.BaseConfig) error {
	if baseConfig.ID == 0 {
		return errors.New("baseConfig can not update without id")
	}
	q := r.data.Query(ctx).BaseConfig
	_, err := q.WithContext(ctx).Updates(baseConfig)
	return err
}

func (r *baseConfigRepo) Create(ctx context.Context, baseConfig *models.BaseConfig) error {
	q := r.data.Query(ctx).BaseConfig
	err := q.WithContext(ctx).Clauses().Create(baseConfig)
	return err
}

func (r *baseConfigRepo) Find(ctx context.Context, condition *BaseConfigCondition) (*models.BaseConfig, error) {
	m := r.data.Query(ctx).BaseConfig
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
	if condition.Name != "" {
		q = q.Where(m.Name.Eq(condition.Name))
	}
	if condition.Type > 0 {
		q = q.Where(m.Type.Eq(condition.Type))
	}
	if condition.Key != "" {
		q = q.Where(m.Key.Eq(condition.Key))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

func (r *baseConfigRepo) FindAll(ctx context.Context, condition *BaseConfigCondition) ([]*models.BaseConfig, error) {
	m := r.data.Query(ctx).BaseConfig
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
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.Type > 0 {
		q = q.Where(m.Type.Eq(condition.Type))
	}
	if condition.Key != "" {
		q = q.Where(m.Key.Like(buildLikeValue(condition.Key)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.UpdatedBy.Desc())
	return q.Find()
}

func (r *baseConfigRepo) ListPage(ctx context.Context, page, size int64, condition *BaseConfigCondition) ([]*models.BaseConfig, int64, error) {
	m := r.data.Query(ctx).BaseConfig
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
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.Type > 0 {
		q = q.Where(m.Type.Eq(condition.Type))
	}
	if condition.Key != "" {
		q = q.Where(m.Key.Like(buildLikeValue(condition.Key)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	q = q.Order(m.UpdatedBy.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseConfigRepo) Count(ctx context.Context, condition *BaseConfigCondition) (int64, error) {
	m := r.data.Query(ctx).BaseConfig
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
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.Type > 0 {
		q = q.Where(m.Type.Eq(condition.Type))
	}
	if condition.Key != "" {
		q = q.Where(m.Key.Like(buildLikeValue(condition.Key)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	count, err := q.Count()
	return count, err
}
