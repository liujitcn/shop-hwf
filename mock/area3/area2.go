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
	err = db.Where("level = 2").Find(&list).Error
	if err != nil {
		return
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Code < list[j].Code
	})

	for _, v1 := range list {
		println(v1.Code, v1.Name)
		oldList := make([]*models.BaseArea, 0)
		err = db.Where("parent_id = (select id from base_area where area_code = ?)", v1.Code).Find(&oldList).Error
		if err != nil {
			return
		}
		sort.Slice(oldList, func(i, j int) bool {
			return oldList[i].AreaCode < oldList[j].AreaCode
		})
		newList := make([]*models.BaseAreaNew, 0)
		for _, area := range oldList {
			newList = append(newList, &models.BaseAreaNew{
				ParentID: v1.ID,
				Level:    3,
				Name:     area.Name,
				Code:     area.AreaCode,
			})
		}
		err = db.CreateInBatches(newList, 100).Error
		if err != nil {
			return
		}

		v1.Level = 99
		db.Save(v1)
	}
}
