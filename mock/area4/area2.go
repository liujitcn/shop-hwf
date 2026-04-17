package main

import (
	_const "gitee.com/liujit/shop/mock/const"
	"gitee.com/liujit/shop/server/lib/data/models"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
)

func main() {
	sqlDb, cleanup, err := sqldb.NewGormSqlDb(_const.Database)
	if err != nil {
		return
	}
	defer cleanup()

	db := sqlDb.GetDb()

	list := make([]*models.BaseAreaNew, 0)
	err = db.Find(&list).Error
	if err != nil {
		return
	}
	Map := make(map[string]*models.BaseAreaNew)
	for _, v1 := range list {
		Map[v1.Code] = v1
	}
	oldList := make([]*models.BaseArea, 0)
	err = db.Find(&oldList).Error
	if err != nil {
		return
	}

	for _, areaNew := range oldList {
		if _, ok := Map[areaNew.AreaCode]; !ok {
			println(areaNew.AreaCode)
		}
	}
}
