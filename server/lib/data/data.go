package data

import (
	"context"
	"fmt"
	"gitee.com/liujit/shop/server/lib/data/query"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
	"gorm.io/gorm"
)

type contextTxKey struct{}

type Data struct {
	query *query.Query
	db    *gorm.DB
}

// NewData .
func NewData(sqlDb sqldb.SqlDb) *Data {
	db := sqlDb.GetDb()
	if db == nil {
		panic(fmt.Errorf("NewQuery need init by db"))
	}
	// 注册回调到全局 GORM 实例
	// 创建前回调：填充创建人和时间
	registerCallbacks(db)

	d := &Data{
		query: query.Use(db),
		db:    db,
	}
	return d
}

type Transaction interface {
	Transaction(context.Context, func(ctx context.Context) error) error
}

func NewTransaction(d *Data) Transaction {
	return d
}

func (d *Data) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.query.Transaction(func(tx *query.Query) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) Query(ctx context.Context) *query.Query {
	tx, ok := ctx.Value(contextTxKey{}).(*query.Query)
	if ok {
		return tx
	}
	return d.query
}

func convertPageSize(page, size int64) (offset, limit int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	offset = int((page - 1) * size)
	limit = int(size)
	return
}

func buildLikeValue(key string) string {
	return fmt.Sprintf("%%%s%%", key)
}
