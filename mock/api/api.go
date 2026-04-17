package main

import (
	_const "gitee.com/liujit/shop/mock/const"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
)

func main() {
	sqlDb, cleanup, err := sqldb.NewGormSqlDb(_const.Database)
	if err != nil {
		return
	}
	defer cleanup()

	db := sqlDb.GetDb()

	// 查询全部api
	baseApi := make([]*models.BaseAPI, 0)
	db.Find(&baseApi)
	baseApiMap := make(map[int64]string)
	for _, item := range baseApi {
		baseApiMap[item.ID] = item.Operation
	}
	// 查询全部菜单
	baseMenu := make([]*models.BaseMenu, 0)
	db.Find(&baseMenu)
	for _, item := range baseMenu {
		// 更新apis字段
		apis := make([]string, 0)
		apisList := str.ConvertJsonStringToInt64Array(item.Apis)
		for _, api := range apisList {
			if v, ok := baseApiMap[api]; ok {
				apis = append(apis, v)
			}
		}
		item.Apis = str.ConvertStringArrayToString(apis)
		//println(item.Apis)
		db.Save(item)
	}
}
