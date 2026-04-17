package main

import (
	"encoding/json"
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

	shopHots := make([]*models.ShopHot, 0)
	err = db.Find(&shopHots).Error
	if err != nil {
		return
	}
	for _, info := range shopHots {
		picList := make([]string, 0)
		json.Unmarshal([]byte(info.Picture), &picList)
		newPicList := make([]string, 0)
		for _, item := range picList {

			newPicList = append(newPicList, fmt.Sprintf("/%s", item))
		}
		newPic, _ := json.Marshal(newPicList)

		info.Picture = string(newPic)
		//

		banner, err := image.Upload(info.Banner, "hot")
		if err != nil {
			fmt.Println(err)
			banner = info.Banner
		}
		info.Banner = banner

		db.Save(info)
	}
}
