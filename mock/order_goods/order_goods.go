package main

import (
	"fmt"
	_const "gitee.com/liujit/shop/mock/const"
	"gitee.com/liujit/shop/server/lib/data/models"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
	"math/rand"
	"time"
)

func main() {
	sqlDb, cleanup, err := sqldb.NewGormSqlDb(_const.Database)
	if err != nil {
		return
	}
	defer cleanup()

	db := sqlDb.GetDb()

	goodsIds := make([]int64, 0)
	list := make([]*models.Goods, 0)
	err = db.Find(&list).Error
	for _, item := range list {
		goodsIds = append(goodsIds, item.ID)
	}
	for i := 1; i < 3001; i++ {
		ids := getRandomElements(goodsIds)
		for _, id := range ids {
			// 生成一个在 1 到 10 之间的随机数字
			randomNumber := getRandom(10)
			orderGoods := models.OrderGoods{
				OrderID:  int64(i),
				GoodsID:  id,
				Num:      int64(randomNumber),
				SpecItem: "[]",
			}
			err = db.Create(&orderGoods).Error
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func getRandomElements(arr []int64) []int64 {
	numElements := getRandom(3)

	// 结果切片
	var result []int64

	// 从数组中随机选择 numElements 个元素
	for i := 0; i < numElements; i++ {
		// 随机选择一个索引
		index := getRandom(len(arr))
		// 将元素添加到结果切片
		result = append(result, arr[index])
	}

	return result
}

func getRandom(n int) int {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := r.Intn(n)
	return randomNumber
}
