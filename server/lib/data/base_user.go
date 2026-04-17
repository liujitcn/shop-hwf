package data

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/liujit/shop/server/lib/data/models"
	"time"
)

type BaseUserCondition struct {
	Id             int64
	Ids            []int64
	DeptId         int64
	Status         int32
	DeptPath       string     // 用户部门路径
	UserName       string     // 用户账号
	NickName       string     // 用户昵称
	Phone          string     // 手机号码
	Openid         string     // 手机号码
	Keyword        string     // 关键字
	StartCreatedAt *time.Time // 创建开始时间
	EndCreatedAt   *time.Time // 创建结束时间
}

type BaseUserRepo interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, baseUser *models.BaseUser) error
	Create(ctx context.Context, baseUser *models.BaseUser) error
	Find(ctx context.Context, condition *BaseUserCondition) (*models.BaseUser, error)
	FindAll(ctx context.Context, condition *BaseUserCondition) ([]*models.BaseUser, error)
	ListPage(ctx context.Context, page, size int64, condition *BaseUserCondition) ([]*models.BaseUser, int64, error)
	Count(ctx context.Context, condition *BaseUserCondition) (int64, error)
}

type baseUserRepo struct {
	data *Data
}

func NewBaseUserRepo(data *Data) BaseUserRepo {
	return &baseUserRepo{data: data}
}

func (r *baseUserRepo) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).BaseUser
	_, err := q.WithContext(ctx).Where(q.ID.In(ids...)).Delete()
	return err
}

func (r *baseUserRepo) UpdateByID(ctx context.Context, baseUser *models.BaseUser) error {
	if baseUser.ID == 0 {
		return errors.New("baseUser can not update without id")
	}
	q := r.data.Query(ctx).BaseUser
	_, err := q.WithContext(ctx).Updates(baseUser)
	return err
}

func (r *baseUserRepo) Create(ctx context.Context, baseUser *models.BaseUser) error {
	q := r.data.Query(ctx).BaseUser
	err := q.WithContext(ctx).Clauses().Create(baseUser)
	return err
}

func (r *baseUserRepo) Find(ctx context.Context, condition *BaseUserCondition) (*models.BaseUser, error) {
	m := r.data.Query(ctx).BaseUser
	q := m.WithContext(ctx)
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.DeptId > 0 {
		q = q.Where(m.DeptID.Eq(condition.DeptId))
	}
	if len(condition.UserName) > 0 {
		q = q.Where(m.UserName.Eq(condition.UserName))
	}
	if condition.NickName != "" {
		q = q.Where(m.NickName.Eq(condition.NickName))
	}
	if condition.Phone != "" {
		q = q.Where(m.Phone.Eq(condition.Phone))
	}
	if condition.Openid != "" {
		q = q.Where(m.Openid.Eq(condition.Openid))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	return q.First()
}

func (r *baseUserRepo) FindAll(ctx context.Context, condition *BaseUserCondition) ([]*models.BaseUser, error) {
	dept := r.data.Query(ctx).BaseDept
	m := r.data.Query(ctx).BaseUser
	q := m.WithContext(ctx)
	if len(condition.DeptPath) > 0 {
		q = q.Join(dept, dept.ID.EqCol(m.DeptID), dept.Path.Like(fmt.Sprintf("%s%%", condition.DeptPath)))
	}
	if condition.DeptId > 0 {
		q = q.Where(m.DeptID.Eq(condition.DeptId))
	}
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if len(condition.UserName) > 0 {
		q = q.Where(m.UserName.Like(buildLikeValue(condition.UserName)))
	}
	if condition.NickName != "" {
		q = q.Where(m.NickName.Like(buildLikeValue(condition.NickName)))
	}
	if condition.Phone != "" {
		q = q.Where(m.Phone.Like(buildLikeValue(condition.Phone)))
	}
	if condition.Keyword != "" {
		q = q.Where(m.UserName.Like(buildLikeValue(condition.Keyword))).Or(m.NickName.Like(buildLikeValue(condition.Keyword))).Or(m.Phone.Like(buildLikeValue(condition.Phone)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	q = q.Order(m.UpdatedAt.Desc())
	return q.Find()
}

func (r *baseUserRepo) ListPage(ctx context.Context, page, size int64, condition *BaseUserCondition) ([]*models.BaseUser, int64, error) {
	dept := r.data.Query(ctx).BaseDept
	m := r.data.Query(ctx).BaseUser
	q := m.WithContext(ctx)
	if len(condition.DeptPath) > 0 {
		q = q.Join(dept, dept.ID.EqCol(m.DeptID), dept.Path.Like(fmt.Sprintf("%s%%", condition.DeptPath)))
	}
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.DeptId > 0 {
		q = q.Where(m.DeptID.Eq(condition.DeptId))
	}
	if len(condition.UserName) > 0 {
		q = q.Where(m.UserName.Like(buildLikeValue(condition.UserName)))
	}
	if condition.NickName != "" {
		q = q.Where(m.NickName.Like(buildLikeValue(condition.NickName)))
	}
	if condition.Phone != "" {
		q = q.Where(m.Phone.Like(buildLikeValue(condition.Phone)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.Keyword != "" {
		q = q.Where(m.UserName.Like(buildLikeValue(condition.Keyword))).Or(m.NickName.Like(buildLikeValue(condition.Keyword))).Or(m.Phone.Like(buildLikeValue(condition.Phone)))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	q = q.Order(m.UpdatedAt.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *baseUserRepo) Count(ctx context.Context, condition *BaseUserCondition) (int64, error) {
	dept := r.data.Query(ctx).BaseDept
	m := r.data.Query(ctx).BaseUser
	q := m.WithContext(ctx)
	if len(condition.DeptPath) > 0 {
		q = q.Join(dept, dept.ID.EqCol(m.DeptID), dept.Path.Like(fmt.Sprintf("%s%%", condition.DeptPath)))
	}
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if len(condition.Ids) > 0 {
		q = q.Where(m.ID.In(condition.Ids...))
	}
	if condition.DeptId > 0 {
		q = q.Where(m.DeptID.Eq(condition.DeptId))
	}
	if len(condition.UserName) > 0 {
		q = q.Where(m.UserName.Like(buildLikeValue(condition.UserName)))
	}
	if condition.NickName != "" {
		q = q.Where(m.NickName.Like(buildLikeValue(condition.NickName)))
	}
	if condition.Phone != "" {
		q = q.Where(m.Phone.Like(buildLikeValue(condition.Phone)))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	count, err := q.Count()
	return count, err
}
