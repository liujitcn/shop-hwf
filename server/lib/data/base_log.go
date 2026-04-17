package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
	"time"
)

type BaseLogCondition struct {
	Id               int64
	Operation        string
	StatusCode       int32
	RequestStartTime *time.Time
	RequestEndTime   *time.Time
}

type BaseLogRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseLog *models.BaseLog) error
	Create(ctx context.Context, baseLog *models.BaseLog) error
	Find(ctx context.Context, condition *BaseLogCondition) (*models.BaseLog, error)
	FindAll(ctx context.Context, condition *BaseLogCondition) ([]*models.BaseLog, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseLogCondition) ([]*models.BaseLog, int64, error)
	Count(ctx context.Context, condition *BaseLogCondition) (int64, error)
}

type baseLogRepo struct {
	data *Data
}

func NewBaseLogRepo(data *Data) BaseLogRepo {
	return &baseLogRepo{data: data}
}

func (r *baseLogRepo) Delete(ctx context.Context, ids []int64) error {
	q := r.data.Query(ctx).BaseLog
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseLogRepo) UpdateByID(ctx context.Context, baseLog *models.BaseLog) error {
	if baseLog.ID == 0 {
		return errors.New("baseLog can not update without id")
	}
	q := r.data.Query(ctx).BaseLog
	_, err := q.WithContext(ctx).Updates(baseLog)
	return err
}

func (r *baseLogRepo) Create(ctx context.Context, baseLog *models.BaseLog) error {
	q := r.data.Query(ctx).BaseLog
	err := q.WithContext(ctx).Clauses().Create(baseLog)
	return err
}

func (r *baseLogRepo) Find(ctx context.Context, condition *BaseLogCondition) (*models.BaseLog, error) {
	m := r.data.Query(ctx).BaseLog
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Operation != "" {
		q = q.Where(m.Operation.Eq(condition.Operation))
	}
	if condition.StatusCode > 0 {
		q = q.Where(m.StatusCode.Eq(condition.StatusCode))
	}
	if condition.RequestStartTime != nil {
		q = q.Where(m.RequestTime.Gte(*condition.RequestStartTime))
	}
	if condition.RequestEndTime != nil {
		q = q.Where(m.RequestTime.Lte(*condition.RequestEndTime))
	}
	return q.First()
}

func (r *baseLogRepo) FindAll(ctx context.Context, condition *BaseLogCondition) ([]*models.BaseLog, error) {
	m := r.data.Query(ctx).BaseLog
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Operation != "" {
		q = q.Where(m.Operation.Like(buildLikeValue(condition.Operation)))
	}
	if condition.StatusCode > 0 {
		q = q.Where(m.StatusCode.Eq(condition.StatusCode))
	}
	if condition.RequestStartTime != nil {
		q = q.Where(m.RequestTime.Gte(*condition.RequestStartTime))
	}
	if condition.RequestEndTime != nil {
		q = q.Where(m.RequestTime.Lte(*condition.RequestEndTime))
	}
	q = q.Order(m.RequestTime.Desc())
	return q.Find()
}

func (r *baseLogRepo) ListPage(ctx context.Context, page, size int64, condition *BaseLogCondition) ([]*models.BaseLog, int64, error) {
	m := r.data.Query(ctx).BaseLog
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Operation != "" {
		q = q.Where(m.Operation.Like(buildLikeValue(condition.Operation)))
	}
	if condition.StatusCode > 0 {
		q = q.Where(m.StatusCode.Eq(condition.StatusCode))
	}
	if condition.RequestStartTime != nil {
		q = q.Where(m.RequestTime.Gte(*condition.RequestStartTime))
	}
	if condition.RequestEndTime != nil {
		q = q.Where(m.RequestTime.Lte(*condition.RequestEndTime))
	}
	q = q.Order(m.RequestTime.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseLogRepo) Count(ctx context.Context, condition *BaseLogCondition) (int64, error) {
	m := r.data.Query(ctx).BaseLog
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Operation != "" {
		q = q.Where(m.Operation.Like(buildLikeValue(condition.Operation)))
	}
	if condition.StatusCode > 0 {
		q = q.Where(m.StatusCode.Eq(condition.StatusCode))
	}
	if condition.RequestStartTime != nil {
		q = q.Where(m.RequestTime.Gte(*condition.RequestStartTime))
	}
	if condition.RequestEndTime != nil {
		q = q.Where(m.RequestTime.Lte(*condition.RequestEndTime))
	}
	count, err := q.Count()
	return count, err
}
