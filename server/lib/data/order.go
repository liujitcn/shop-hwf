package data

import (
	"context"
	"errors"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data/dto"
	"gitee.com/liujit/shop/server/lib/data/models"
	"time"
)

type OrderCondition struct {
	Id             int64
	UserId         int64
	OrderNo        string
	Status         int32
	PayType        int32      // 支付方式，1为在线支付，2为货到付款
	PayChannel     int32      // 支付渠道：支付渠道，1支付宝、2微信--支付方式为在线支付时，传值，为货到付款时，不传值
	StartCreatedAt *time.Time // 创建开始时间
	EndCreatedAt   *time.Time // 创建结束时间
}

type OrderRepo interface {
	Delete(ctx context.Context, userId int64, ids []int64) error
	UpdateByID(ctx context.Context, userId int64, orderInfo *models.Order) error
	Create(ctx context.Context, orderInfo *models.Order) error
	Find(ctx context.Context, condition *OrderCondition) (*models.Order, error)
	FindAll(ctx context.Context, condition *OrderCondition) ([]*models.Order, error)
	ListPage(ctx context.Context, page, size int64, condition *OrderCondition) ([]*models.Order, int64, error)
	Count(ctx context.Context, condition *OrderCondition) (int64, error)
	UpdateByIds(ctx context.Context, userId int64, ids []int64, orderInfo *models.Order) error
	MapCount(ctx context.Context, userId int64) (map[int32]int32, error)
	FindByOrderNo(ctx context.Context, orderNo string) (*models.Order, error)
	Sum(ctx context.Context, condition *OrderCondition) (int64, error)
	OrderSummary(ctx context.Context, timeType int32, condition *OrderCondition) ([]*dto.OrderSummary, error)
}

type orderRepo struct {
	data *Data
}

func NewOrderRepo(data *Data) OrderRepo {
	return &orderRepo{data: data}
}

func (r *orderRepo) Delete(ctx context.Context, userId int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).Order
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId), q.ID.In(ids...)).Delete()
	return err
}

func (r *orderRepo) UpdateByID(ctx context.Context, userId int64, orderInfo *models.Order) error {
	if orderInfo.ID == 0 {
		return errors.New("orderInfo can not update without id")
	}
	q := r.data.Query(ctx).Order
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId)).Updates(orderInfo)
	return err
}

func (r *orderRepo) Create(ctx context.Context, orderInfo *models.Order) error {
	q := r.data.Query(ctx).Order
	err := q.WithContext(ctx).Clauses().Create(orderInfo)
	return err
}

func (r *orderRepo) Find(ctx context.Context, condition *OrderCondition) (*models.Order, error) {
	m := r.data.Query(ctx).Order
	q := m.WithContext(ctx)
	if condition.UserId > 0 {
		q = q.Where(m.UserID.Eq(condition.UserId))
	}
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.OrderNo != "" {
		q = q.Where(m.OrderNo.Eq(condition.OrderNo))
	}
	if condition.PayType > 0 {
		q = q.Where(m.PayType.Eq(condition.PayType))
	}
	if condition.PayChannel > 0 {
		q = q.Where(m.PayChannel.Eq(condition.PayChannel))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	return q.First()
}

func (r *orderRepo) FindAll(ctx context.Context, condition *OrderCondition) ([]*models.Order, error) {
	m := r.data.Query(ctx).Order
	q := m.WithContext(ctx)
	if condition.UserId > 0 {
		q = q.Where(m.UserID.Eq(condition.UserId))
	}
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.OrderNo != "" {
		q = q.Where(m.OrderNo.Like(buildLikeValue(condition.OrderNo)))
	}
	if condition.PayType > 0 {
		q = q.Where(m.PayType.Eq(condition.PayType))
	}
	if condition.PayChannel > 0 {
		q = q.Where(m.PayChannel.Eq(condition.PayChannel))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	} else {
		q = q.Where(m.Status.Neq(int32(common.OrderStatus_DELETED)))
	}
	q = q.Order(m.CreatedAt.Desc())
	return q.Find()
}

func (r *orderRepo) ListPage(ctx context.Context, page, size int64, condition *OrderCondition) ([]*models.Order, int64, error) {
	m := r.data.Query(ctx).Order
	q := m.WithContext(ctx)
	if condition.UserId > 0 {
		q = q.Where(m.UserID.Eq(condition.UserId))
	}
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.OrderNo != "" {
		q = q.Where(m.OrderNo.Like(buildLikeValue(condition.OrderNo)))
	}
	if condition.PayType > 0 {
		q = q.Where(m.PayType.Eq(condition.PayType))
	}
	if condition.PayChannel > 0 {
		q = q.Where(m.PayChannel.Eq(condition.PayChannel))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	} else {
		q = q.Where(m.Status.Neq(int32(common.OrderStatus_DELETED)))
	}
	q = q.Order(m.CreatedAt.Desc())
	offset, limit := convertPageSize(page, size)
	return q.FindByPage(offset, limit)
}
func (r *orderRepo) Count(ctx context.Context, condition *OrderCondition) (int64, error) {
	m := r.data.Query(ctx).Order
	q := m.WithContext(ctx)
	if condition.UserId > 0 {
		q = q.Where(m.UserID.Eq(condition.UserId))
	}
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.OrderNo != "" {
		q = q.Where(m.OrderNo.Like(buildLikeValue(condition.OrderNo)))
	}
	if condition.PayType > 0 {
		q = q.Where(m.PayType.Eq(condition.PayType))
	}
	if condition.PayChannel > 0 {
		q = q.Where(m.PayChannel.Eq(condition.PayChannel))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	} else {
		q = q.Where(m.Status.Neq(int32(common.OrderStatus_DELETED)))
	}
	return q.Count()
}

func (r *orderRepo) UpdateByIds(ctx context.Context, userId int64, ids []int64, orderInfo *models.Order) error {
	if len(ids) == 0 {
		return errors.New("orderInfo can not update without id")
	}
	q := r.data.Query(ctx).Order
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId), q.ID.In(ids...)).Updates(orderInfo)
	return err
}

func (r *orderRepo) MapCount(ctx context.Context, userId int64) (map[int32]int32, error) {
	m := r.data.Query(ctx).Order
	q := m.WithContext(ctx)
	var results []struct {
		Status int32
		Count  int32
	}
	q = q.Select(m.Status.As("status"), m.Status.Count().As("count")).Where(m.UserID.Eq(userId)).Group(m.Status)
	err := q.Scan(&results)
	if err != nil {
		return nil, err
	}
	total := make(map[int32]int32)
	for _, result := range results {
		total[result.Status] = result.Count
	}
	return total, nil
}

func (r *orderRepo) FindByOrderNo(ctx context.Context, orderNo string) (*models.Order, error) {
	m := r.data.Query(ctx).Order
	q := m.WithContext(ctx)
	q = q.Where(m.OrderNo.Eq(orderNo))
	return q.First()
}
func (r *orderRepo) Sum(ctx context.Context, condition *OrderCondition) (int64, error) {
	m := r.data.Query(ctx).Order
	q := m.WithContext(ctx)
	if condition.UserId > 0 {
		q = q.Where(m.UserID.Eq(condition.UserId))
	}
	if condition.Id > 0 {
		q = q.Where(m.ID.Eq(condition.Id))
	}
	if condition.OrderNo != "" {
		q = q.Where(m.OrderNo.Like(buildLikeValue(condition.OrderNo)))
	}
	if condition.PayType > 0 {
		q = q.Where(m.PayType.Eq(condition.PayType))
	}
	if condition.PayChannel > 0 {
		q = q.Where(m.PayChannel.Eq(condition.PayChannel))
	}
	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	} else {
		q = q.Where(m.Status.Neq(int32(common.OrderStatus_DELETED)))
	}
	var results struct{ Total int64 }
	q.Select(m.PayMoney.Sum().As("total"))
	err := q.Scan(&results)
	if err != nil {
		return 0, err
	}
	return results.Total, nil
}

func (r *orderRepo) OrderSummary(ctx context.Context, timeType int32, condition *OrderCondition) ([]*dto.OrderSummary, error) {
	m := r.data.Query(ctx).Order
	q := m.WithContext(ctx)

	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	} else {
		q = q.Where(m.Status.Neq(int32(common.OrderStatus_DELETED)))
	}
	results := make([]*dto.OrderSummary, 0)
	switch timeType {
	case 2:
		q.Select(m.CreatedAt.Month().As("key"),
			m.ID.Count().As("order_count"),
			m.PayMoney.Sum().As("sale_amount")).Group(m.CreatedAt.Month()).Order(m.CreatedAt.Month())
	case 1:
		q.Select(m.CreatedAt.DayOfWeek().As("key"),
			m.ID.Count().As("order_count"),
			m.PayMoney.Sum().As("sale_amount")).Group(m.CreatedAt.DayOfWeek()).Order(m.CreatedAt.DayOfWeek())
	default:
		q.Select(m.CreatedAt.DayOfMonth().As("key"),
			m.ID.Count().As("order_count"),
			m.PayMoney.Sum().As("sale_amount")).Group(m.CreatedAt.DayOfMonth()).Order(m.CreatedAt.DayOfMonth())
	}
	err := q.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
