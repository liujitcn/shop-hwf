package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
	"time"
)

type BaseJobLogCondition struct {
	Id               int64
	JobId            int64
	Status           int32
	ExecuteStartTime *time.Time
	ExecuteEndTime   *time.Time
}

type BaseJobLogRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseJobLog *models.BaseJobLog) error
	Create(ctx context.Context, baseJobLog *models.BaseJobLog) error
	Find(ctx context.Context, condition *BaseJobLogCondition) (*models.BaseJobLog, error)
	FindAll(ctx context.Context, condition *BaseJobLogCondition) ([]*models.BaseJobLog, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseJobLogCondition) ([]*models.BaseJobLog, int64, error)
	Count(ctx context.Context, condition *BaseJobLogCondition) (int64, error)
}

type baseJobLogRepo struct {
	data *Data
}

func NewBaseJobLogRepo(data *Data) BaseJobLogRepo {
	return &baseJobLogRepo{data: data}
}

func (r *baseJobLogRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).BaseJobLog
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseJobLogRepo) UpdateByID(ctx context.Context, baseJobLog *models.BaseJobLog) error {
	if baseJobLog.ID == 0 {
		return errors.New("baseJobLog can not update without id")
	}
	q := r.data.Query(ctx).BaseJobLog
	_, err := q.WithContext(ctx).Updates(baseJobLog)
	return err
}

func (r *baseJobLogRepo) Create(ctx context.Context, baseJobLog *models.BaseJobLog) error {
	q := r.data.Query(ctx).BaseJobLog
	err := q.WithContext(ctx).Clauses().Create(baseJobLog)
	return err
}

func (r *baseJobLogRepo) Find(ctx context.Context, condition *BaseJobLogCondition) (*models.BaseJobLog, error) {
	m := r.data.Query(ctx).BaseJobLog
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.JobId > 0 {
		q = q.Where(m.JobID.Eq(condition.JobId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.ExecuteStartTime != nil {
		q = q.Where(m.ExecuteTime.Gte(*condition.ExecuteStartTime))
	}
	if condition.ExecuteEndTime != nil {
		q = q.Where(m.ExecuteTime.Lte(*condition.ExecuteEndTime))
	}
	return q.First()
}

func (r *baseJobLogRepo) FindAll(ctx context.Context, condition *BaseJobLogCondition) ([]*models.BaseJobLog, error) {
	m := r.data.Query(ctx).BaseJobLog
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.JobId > 0 {
		q = q.Where(m.JobID.Eq(condition.JobId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.ExecuteStartTime != nil {
		q = q.Where(m.ExecuteTime.Gte(*condition.ExecuteStartTime))
	}
	if condition.ExecuteEndTime != nil {
		q = q.Where(m.ExecuteTime.Lte(*condition.ExecuteEndTime))
	}
	return q.Find()
}

func (r *baseJobLogRepo) ListPage(ctx context.Context, page, size int64, condition *BaseJobLogCondition) ([]*models.BaseJobLog, int64, error) {
	m := r.data.Query(ctx).BaseJobLog
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.JobId > 0 {
		q = q.Where(m.JobID.Eq(condition.JobId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.ExecuteStartTime != nil {
		q = q.Where(m.ExecuteTime.Gte(*condition.ExecuteStartTime))
	}
	if condition.ExecuteEndTime != nil {
		q = q.Where(m.ExecuteTime.Lte(*condition.ExecuteEndTime))
	}
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseJobLogRepo) Count(ctx context.Context, condition *BaseJobLogCondition) (int64, error) {
	m := r.data.Query(ctx).BaseJobLog
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.JobId > 0 {
		q = q.Where(m.JobID.Eq(condition.JobId))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.ExecuteStartTime != nil {
		q = q.Where(m.ExecuteTime.Gte(*condition.ExecuteStartTime))
	}
	if condition.ExecuteEndTime != nil {
		q = q.Where(m.ExecuteTime.Lte(*condition.ExecuteEndTime))
	}
	count, err := q.Count()
	return count, err
}
