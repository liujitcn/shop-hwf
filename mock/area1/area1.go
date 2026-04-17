package main

import (
	_const "gitee.com/liujit/shop/mock/const"
	"gitee.com/liujit/shop/server/lib/data/models"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
	"sort"
)

func main() {
	sqlDb, cleanup, err := sqldb.NewGormSqlDb(_const.Database)
	if err != nil {
		return
	}
	defer cleanup()

	db := sqlDb.GetDb()

	list := make([]*models.BaseArea, 0)
	err = db.Where("level = 1").Find(&list).Error
	if err != nil {
		return
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].AreaCode < list[j].AreaCode
	})

	newList := make([]*models.BaseAreaNew, 0)
	for _, v1 := range list {
		newList = append(newList, &models.BaseAreaNew{
			ParentID: 0,
			Level:    1,
			Name:     v1.Name,
			Code:     v1.AreaCode,
		})
	}
	err = db.CreateInBatches(newList, 100).Error
	if err != nil {
		return
	}
}
