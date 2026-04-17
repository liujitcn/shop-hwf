package utils

import (
	"fmt"
	"gitee.com/liujit/shop/server/api/admin"
	"strconv"
	"time"
)

var week = []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}
var month = []string{"一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"}

// CalcGrowthRate 通用增长率计算函数（处理除零情况）
func CalcGrowthRate(prev, current int64) int64 {
	if prev == 0 {
		if current == 0 {
			return 0.0
		}
		return 10000 // 当基数为0时视为100%增长
	}
	res := (current - prev) * 10000 / prev
	return res
}

func GetCreatedAt(timeType admin.DashboardTimeType) (time.Time, time.Time) {
	now := time.Now()
	year, month, day := now.Date()
	var startCreatedAt, endCreatedAt time.Time
	switch timeType {
	case admin.DashboardTimeType_MONTH:
		startCreatedAt = time.Date(year, month, day, 0, 0, 0, 0, now.Location())
		endCreatedAt = startCreatedAt.AddDate(0, 0, 1)
	case admin.DashboardTimeType_WEEK:
		offset := (int(now.Weekday()) - int(time.Monday) + 7) % 7
		t := now.AddDate(0, 0, -offset)
		startCreatedAt = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		endCreatedAt = startCreatedAt.AddDate(0, 0, 7)
	default:
		startCreatedAt = time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
		endCreatedAt = startCreatedAt.AddDate(0, 1, 0)
	}
	return startCreatedAt, endCreatedAt
}

func FormatDate(timeType admin.DashboardTimeType, key int) string {
	switch timeType {
	case admin.DashboardTimeType_MONTH:
		return month[key]
	case admin.DashboardTimeType_WEEK:
		return week[key]
	default:
		return fmt.Sprintf("%d日", key+1)
	}
}

// ConvertYuanToFen 元转分
func ConvertYuanToFen(yuan string) int64 {
	float, err := strconv.ParseFloat(yuan, 64)
	if err != nil {
		return 0
	}
	return int64(float * 100)
}
