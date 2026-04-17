package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/lib/data/models"
)

type UserCartCondition struct {
	Id        int64
	GoodsId   int64
	SkuCode   string
	IsChecked bool
}

type UserCartRepo interface {
	Delete(ctx context.Context, userId int64, ids []int64) error
	UpdateByID(ctx context.Context, userId int64, userCart *models.UserCart) error
	Create(ctx context.Context, userCart *models.UserCart) error
	Find(ctx context.Context, userId int64, condition *UserCartCondition) (*models.UserCart, error)
	FindAll(ctx context.Context, userId int64, condition *UserCartCondition) ([]*models.UserCart, error)
	ListPage(ctx context.Context, userId int64, page, size int64, condition *UserCartCondition) ([]*models.UserCart, int64, error)
	Count(ctx context.Context, userId int64, condition *UserCartCondition) (int64, error)
	UpdateByUserId(ctx context.Context, userId int64, userCart *models.UserCart) error
	DeleteByGoodsId(ctx context.Context, userId int64, goodsId int64, skuCode string) error
}

type userCartRepo struct {
	data *Data
}

func (r *userCartRepo) DeleteByGoodsId(ctx context.Context, userId int64, goodsId int64, skuCode string) error {
	if goodsId == 0 || len(skuCode) == 0 {
		return nil
	}
	q := r.data.Query(ctx).UserCart
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId), q.GoodsID.Eq(goodsId), q.SkuCode.Eq(skuCode)).Delete()
	return err
}

func NewUserCartRepo(data *Data) UserCartRepo {
	return &userCartRepo{data: data}
}

func (r *userCartRepo) Delete(ctx context.Context, userId int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).UserCart
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId), q.ID.In(ids...)).Delete()
	return err
}

func (r *userCartRepo) UpdateByID(ctx context.Context, userId int64, userCart *models.UserCart) error {
	if userCart.ID == 0 {
		return errors.New("userCart can not update without id")
	}
	q := r.data.Query(ctx).UserCart
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId)).Updates(userCart)
	return err
}

func (r *userCartRepo) Create(ctx context.Context, userCart *models.UserCart) error {
	q := r.data.Query(ctx).UserCart
	err := q.WithContext(ctx).Clauses().Create(userCart)
	return err
}

func (r *userCartRepo) Find(ctx context.Context, userId int64, condition *UserCartCondition) (*models.UserCart, error) {
	m := r.data.Query(ctx).UserCart
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	if condition.SkuCode != "" {
		q = q.Where(m.SkuCode.Eq(condition.SkuCode))
	}
	if condition.IsChecked {
		q = q.Where(m.IsChecked.Is(condition.IsChecked))
	}
	return q.First()
}

func (r *userCartRepo) FindAll(ctx context.Context, userId int64, condition *UserCartCondition) ([]*models.UserCart, error) {
	m := r.data.Query(ctx).UserCart
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	if condition.SkuCode != "" {
		q = q.Where(m.SkuCode.Eq(condition.SkuCode))
	}
	if condition.IsChecked {
		q = q.Where(m.IsChecked.Is(condition.IsChecked))
	}
	q = q.Order(m.UpdatedAt.Desc())
	return q.Find()
}

func (r *userCartRepo) ListPage(ctx context.Context, userId int64, page, size int64, condition *UserCartCondition) ([]*models.UserCart, int64, error) {
	m := r.data.Query(ctx).UserCart
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	if condition.SkuCode != "" {
		q = q.Where(m.SkuCode.Eq(condition.SkuCode))
	}
	if condition.IsChecked {
		q = q.Where(m.IsChecked.Is(condition.IsChecked))
	}
	q = q.Order(m.UpdatedAt.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}

func (r *userCartRepo) Count(ctx context.Context, userId int64, condition *UserCartCondition) (int64, error) {
	m := r.data.Query(ctx).UserCart
	q := m.WithContext(ctx)
	q = q.Where(m.UserID.Eq(userId))
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.GoodsId > 0 {
		q = q.Where(m.GoodsID.Eq(condition.GoodsId))
	}
	if condition.SkuCode != "" {
		q = q.Where(m.SkuCode.Eq(condition.SkuCode))
	}
	if condition.IsChecked {
		q = q.Where(m.IsChecked.Is(condition.IsChecked))
	}
	return q.Count()
}

func (r *userCartRepo) UpdateByUserId(ctx context.Context, userId int64, userCart *models.UserCart) error {
	q := r.data.Query(ctx).UserCart
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId)).Updates(userCart)
	return err
}
