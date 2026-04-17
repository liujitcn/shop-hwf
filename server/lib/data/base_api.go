package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseApiCondition struct {
	Id          int64
	Ids         []int64
	ServiceName string // 服务名
	ServiceDesc string // 服务描述
	Desc        string // 描述
	Path        string // 请求地址
}

type BaseApiRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseApi *models.BaseAPI) error
	Create(ctx context.Context, baseApi *models.BaseAPI) error
	Find(ctx context.Context, condition *BaseApiCondition) (*models.BaseAPI, error)
	FindAll(ctx context.Context, condition *BaseApiCondition) ([]*models.BaseAPI, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseApiCondition) ([]*models.BaseAPI, int64, error)
	Count(ctx context.Context, condition *BaseApiCondition) (int64, error)
	BatchCreate(ctx context.Context, list []*models.BaseAPI) error
}

type baseApiRepo struct {
	data *Data
}

func NewBaseApiRepo(data *Data) BaseApiRepo {
	return &baseApiRepo{data: data}
}

func (r *baseApiRepo) Delete(ctx context.Context, ids []int64) error {
	q := r.data.Query(ctx).BaseAPI
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseApiRepo) UpdateByID(ctx context.Context, baseApi *models.BaseAPI) error {
	if baseApi.ID == 0 {
		return errors.New("baseApi can not update without id")
	}
	q := r.data.Query(ctx).BaseAPI
	_, err := q.WithContext(ctx).Updates(baseApi)
	return err
}

func (r *baseApiRepo) Create(ctx context.Context, baseApi *models.BaseAPI) error {
	q := r.data.Query(ctx).BaseAPI
	err := q.WithContext(ctx).Clauses().Create(baseApi)
	return err
}

func (r *baseApiRepo) Find(ctx context.Context, condition *BaseApiCondition) (*models.BaseAPI, error) {
	m := r.data.Query(ctx).BaseAPI
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.ServiceName != "" {
		q = q.Where(m.ServiceName.Like(buildLikeValue(condition.ServiceDesc)))
	}
	if condition.ServiceDesc != "" {
		q = q.Where(m.ServiceDesc.Like(buildLikeValue(condition.ServiceDesc)))
	}
	if condition.Desc != "" {
		q = q.Where(m.Desc.Like(buildLikeValue(condition.Desc)))
	}
	if condition.Path != "" {
		q = q.Where(m.Path.Like(buildLikeValue(condition.Path)))
	}
	return q.First()
}

func (r *baseApiRepo) FindAll(ctx context.Context, condition *BaseApiCondition) ([]*models.BaseAPI, error) {
	m := r.data.Query(ctx).BaseAPI
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.ServiceName != "" {
		q = q.Where(m.ServiceName.Like(buildLikeValue(condition.ServiceDesc)))
	}
	if condition.ServiceDesc != "" {
		q = q.Where(m.ServiceDesc.Like(buildLikeValue(condition.ServiceDesc)))
	}
	if condition.Desc != "" {
		q = q.Where(m.Desc.Like(buildLikeValue(condition.Desc)))
	}
	if condition.Path != "" {
		q = q.Where(m.Path.Like(buildLikeValue(condition.Path)))
	}
	return q.Find()
}

func (r *baseApiRepo) ListPage(ctx context.Context, page, size int64, condition *BaseApiCondition) ([]*models.BaseAPI, int64, error) {
	m := r.data.Query(ctx).BaseAPI
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.ServiceName != "" {
		q = q.Where(m.ServiceName.Like(buildLikeValue(condition.ServiceDesc)))
	}
	if condition.ServiceDesc != "" {
		q = q.Where(m.ServiceDesc.Like(buildLikeValue(condition.ServiceDesc)))
	}
	if condition.Desc != "" {
		q = q.Where(m.Desc.Like(buildLikeValue(condition.Desc)))
	}
	if condition.Path != "" {
		q = q.Where(m.Path.Like(buildLikeValue(condition.Path)))
	}
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseApiRepo) Count(ctx context.Context, condition *BaseApiCondition) (int64, error) {
	m := r.data.Query(ctx).BaseAPI
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.ServiceName != "" {
		q = q.Where(m.ServiceName.Like(buildLikeValue(condition.ServiceDesc)))
	}
	if condition.ServiceDesc != "" {
		q = q.Where(m.ServiceDesc.Like(buildLikeValue(condition.ServiceDesc)))
	}
	if condition.Desc != "" {
		q = q.Where(m.Desc.Like(buildLikeValue(condition.Desc)))
	}
	if condition.Path != "" {
		q = q.Where(m.Path.Like(buildLikeValue(condition.Path)))
	}
	count, err := q.Count()
	return count, err
}
func (r *baseApiRepo) BatchCreate(ctx context.Context, list []*models.BaseAPI) error {
	q := r.data.Query(ctx).BaseAPI
	err := q.WithContext(ctx).Clauses().CreateInBatches(list, 100)
	return err
}
