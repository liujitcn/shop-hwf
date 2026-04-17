package main

import (
	"fmt"
	_const "gitee.com/liujit/shop/mock/const"
	"gitee.com/liujit/shop/mock/image"
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

	goodsSkus := make([]*models.GoodsSku, 0)
	err = db.Find(&goodsSkus).Error
	if err != nil {
		return
	}
	for _, info := range goodsSkus {

		Picture, err := image.Upload(info.Picture, "goods")
		if err != nil {
			fmt.Println(err)
			Picture = info.Picture
		}
		info.Picture = Picture

		db.Save(info)
	}
}
