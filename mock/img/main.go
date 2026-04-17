package main

import (
	"encoding/json"
	"fmt"
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

	goodsList := make([]*models.Goods, 0)
	err = db.Find(&goodsList).Error
	if err != nil {
		return
	}
	for _, info := range goodsList {
		bannerList := make([]string, 0)
		json.Unmarshal([]byte(info.Banner), &bannerList)
		newBannerList := make([]string, 0)
		for _, item := range bannerList {
			//picture, err := image.Upload(item, "goods")
			//if err != nil {
			//	fmt.Println(err)
			//	picture = item
			//}
			newBannerList = append(newBannerList, fmt.Sprintf("/%s", item))
		}
		newBanner, _ := json.Marshal(newBannerList)

		info.Banner = string(newBanner)

		detailList := make([]string, 0)
		json.Unmarshal([]byte(info.Detail), &detailList)
		newDetailList := make([]string, 0)
		for _, item := range detailList {
			//picture, err := image.Upload(item, "goods")
			//if err != nil {
			//	fmt.Println(err)
			//	picture = item
			//}
			newDetailList = append(newDetailList, fmt.Sprintf("/%s", item))
		}
		newDetail, _ := json.Marshal(newDetailList)

		info.Detail = string(newDetail)

		db.Save(info)
	}
}
