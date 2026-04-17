package main

import (
	_const "gitee.com/liujit/shop/mock/const"
	"gitee.com/liujit/shop/server/lib/data/models"
	"go.newcapec.cn/nctcommon/nmslib/util"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	sqlDb, cleanup, err := sqldb.NewGormSqlDb(_const.Database)
	if err != nil {
		return
	}
	defer cleanup()

	db := sqlDb.GetDb()

	status := []int32{2, 3, 4, 97, 98, 99, 2, 3, 4, 97}

	tmGenerate, _ := util.NewTmGenerate()
	for i := 0; i < 3000; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNumber := r.Intn(10000)
		t := getRandomTimeInRange()
		order := models.Order{
			OrderNo:      strconv.FormatInt(tmGenerate.NextVal(), 10),
			UserID:       1,
			PayMoney:     int64(randomNumber),
			TotalMoney:   int64(randomNumber),
			PostFee:      0,
			GoodsNum:     0,
			PayType:      1,
			PayChannel:   1,
			deliveryTime: 1,
			Status:       status[i%10],
			Remark:       "",
			CreatedAt:    t,
			UpdatedAt:    t,
		}

		err = db.Create(&order).Error
		if err != nil {
			return
		}
	}
}

func getRandomTimeInRange() time.Time {
	// 获取开始和结束时间的 Unix 时间戳
	// 设置时间范围
	startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) // 2023-04-01 00:00:00
	endTime := time.Date(2025, 4, 14, 0, 0, 0, 0, time.UTC)  // 2023-04-10 00:00:00
	startUnix := startTime.Unix()
	endUnix := endTime.Unix()

	// 在时间范围内生成随机的 Unix 时间戳
	randomUnix := rand.Int63n(endUnix-startUnix) + startUnix

	// 将随机的 Unix 时间戳转换为 time.Time 类型
	randomTime := time.Unix(randomUnix, 0)

	return randomTime
}
