package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type BaseJobCondition struct {
	Id           int64
	Status       int32
	Name         string
	InvokeTarget string
}

type BaseJobRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseJob *models.BaseJob) error
	Create(ctx context.Context, baseJob *models.BaseJob) error
	Find(ctx context.Context, condition *BaseJobCondition) (*models.BaseJob, error)
	FindAll(ctx context.Context, condition *BaseJobCondition) ([]*models.BaseJob, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseJobCondition) ([]*models.BaseJob, int64, error)
	Count(ctx context.Context, condition *BaseJobCondition) (int64, error)
	CleanEntryId(ctx context.Context, entryId int32) error
}

type baseJobRepo struct {
	data *Data
}

func NewBaseJobRepo(data *Data) BaseJobRepo {
	return &baseJobRepo{data: data}
}

func (r *baseJobRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).BaseJob
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseJobRepo) UpdateByID(ctx context.Context, baseJob *models.BaseJob) error {
	if baseJob.ID == 0 {
		return errors.New("baseJob can not update without id")
	}
	q := r.data.Query(ctx).BaseJob
	_, err := q.WithContext(ctx).Updates(baseJob)
	return err
}

func (r *baseJobRepo) Create(ctx context.Context, baseJob *models.BaseJob) error {
	q := r.data.Query(ctx).BaseJob
	err := q.WithContext(ctx).Clauses().Create(baseJob)
	return err
}

func (r *baseJobRepo) Find(ctx context.Context, condition *BaseJobCondition) (*models.BaseJob, error) {
	m := r.data.Query(ctx).BaseJob
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Eq(condition.Name))
	}
	if condition.InvokeTarget != "" {
		q = q.Where(m.InvokeTarget.Eq(condition.InvokeTarget))
	}
	return q.First()
}

func (r *baseJobRepo) FindAll(ctx context.Context, condition *BaseJobCondition) ([]*models.BaseJob, error) {
	m := r.data.Query(ctx).BaseJob
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.InvokeTarget != "" {
		q = q.Where(m.InvokeTarget.Like(buildLikeValue(condition.InvokeTarget)))
	}
	return q.Find()
}

func (r *baseJobRepo) ListPage(ctx context.Context, page, size int64, condition *BaseJobCondition) ([]*models.BaseJob, int64, error) {
	m := r.data.Query(ctx).BaseJob
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.InvokeTarget != "" {
		q = q.Where(m.InvokeTarget.Like(buildLikeValue(condition.InvokeTarget)))
	}
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseJobRepo) Count(ctx context.Context, condition *BaseJobCondition) (int64, error) {
	m := r.data.Query(ctx).BaseJob
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.Name != "" {
		q = q.Where(m.Name.Like(buildLikeValue(condition.Name)))
	}
	if condition.InvokeTarget != "" {
		q = q.Where(m.InvokeTarget.Like(buildLikeValue(condition.InvokeTarget)))
	}
	count, err := q.Count()
	return count, err
}

func (r *baseJobRepo) CleanEntryId(ctx context.Context, entryId int32) error {
	var err error
	m := r.data.Query(ctx).BaseJob
	q := m.WithContext(ctx)
	if entryId > 0 {
		_, err = q.Where(m.EntryID.Eq(entryId)).UpdateColumn(m.EntryID, 0)
	} else {
		_, err = q.Where(m.EntryID.Gt(entryId)).Updates(&models.BaseJob{
			EntryID: 0,
		})
	}
	return err
}
