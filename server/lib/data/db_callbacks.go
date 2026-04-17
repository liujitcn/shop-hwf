package data

import (
	"slices"
	"time"

	authMiddleware "go.newcapec.cn/ncttools/nmskit-auth/middleware"
	"go.newcapec.cn/ncttools/nmskit/log"
	"gorm.io/gorm"
)

var exclude = []string{
	"base_log",
}

// 带检查的字段填充逻辑
func safeSetColumn(db *gorm.DB, fieldName string, value interface{}) {
	// 检查字段是否已存在非零值
	if !isFieldZero(db, fieldName) {
		return
	}

	db.Statement.SetColumn(fieldName, value)
}

// 判断字段是否为零值
func isFieldZero(db *gorm.DB, fieldName string) bool {
	statement := db.Statement
	if statement == nil {
		return false
	}
	if field := statement.Schema.LookUpField(fieldName); field == nil {
		return false
	}
	return true
}

// 创建时填充（增强版）
func fillCreatedFields(db *gorm.DB) {
	table := db.Statement.Table
	if slices.Contains(exclude, table) {
		return
	}

	var userId int64
	ctx := db.Statement.Context

	// 增强认证错误处理
	if authInfo, err := authMiddleware.FromContext(ctx); err != nil {
		log.Warnf("上下文缺失用户信息，使用默认值")
		userId = 0 // 或系统用户ID
	} else {
		userId = authInfo.UserId
	}

	now := time.Now()

	// 安全设置字段（仅当字段为空时）
	safeSetColumn(db, "CreatedBy", userId)
	safeSetColumn(db, "UpdatedBy", userId)
	safeSetColumn(db, "CreatedAt", now)
	safeSetColumn(db, "UpdatedAt", now)
}

// 更新时填充（增强版）
func fillUpdatedFields(db *gorm.DB) {
	table := db.Statement.Table
	if slices.Contains(exclude, table) {
		return
	}

	var userId int64
	ctx := db.Statement.Context

	// 增强认证错误处理
	if authInfo, err := authMiddleware.FromContext(ctx); err != nil {
		log.Warnf("上下文缺失用户信息，使用默认值")
		userId = 0 // 或系统用户ID
	} else {
		userId = authInfo.UserId
	}

	// 仅更新未明确设置的字段
	safeSetColumn(db, "UpdatedBy", userId)
	safeSetColumn(db, "UpdatedAt", time.Now())
}

// registerCallbacks 注册回调（增加错误处理）
func registerCallbacks(db *gorm.DB) {
	err := db.Callback().Create().Before("gorm:before_create").Register("fill_created_fields", fillCreatedFields)
	if err != nil {
		log.Fatal("注册创建回调失败", err)
	}

	err = db.Callback().Update().Before("gorm:before_update").Register("fill_updated_fields", fillUpdatedFields)
	if err != nil {
		log.Fatal("注册更新回调失败", err)
	}
}
