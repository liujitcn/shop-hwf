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

	list := make([]*models.BaseAreaNew, 0)
	err = db.Where("level = 1").Find(&list).Error
	if err != nil {
		return
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Code < list[j].Code
	})

	newList := make([]*models.BaseAreaNew, 0)
	for _, v1 := range list {
		oldList := make([]*models.BaseArea, 0)
		err = db.Where("parent_id = (select id from base_area where area_code = ?)", v1.Code).Find(&oldList).Error
		if err != nil {
			return
		}
		sort.Slice(oldList, func(i, j int) bool {
			return oldList[i].AreaCode < oldList[j].AreaCode
		})
		for _, area := range oldList {
			newList = append(newList, &models.BaseAreaNew{
				ParentID: v1.ID,
				Level:    2,
				Name:     area.Name,
				Code:     area.AreaCode,
			})
		}
	}
	err = db.CreateInBatches(newList, 1000).Error
	if err != nil {
		return
	}
}
