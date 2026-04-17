package models

import (
	"gitee.com/liujit/shop/server/cmd/migrate/utils"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
)

func TableCreate() []interface{} {
	var m = []interface{}{
		new(models.BaseAPI),
		new(models.BaseArea),
		new(models.BaseConfig),
		new(models.BaseDept),
		new(models.BaseDict),
		new(models.BaseDictItem),
		new(models.BaseJob),
		new(models.BaseJobLog),
		new(models.BaseLog),
		new(models.BaseMenu),
		new(models.BaseRole),
		new(models.BaseUser),
		new(models.CasbinRule),
		new(models.Goods),
		new(models.GoodsCategory),
		new(models.GoodsProp),
		new(models.GoodsSku),
		new(models.GoodsSpec),
		new(models.Order),
		new(models.OrderAddress),
		new(models.OrderCancel),
		new(models.OrderGoods),
		new(models.OrderLogistics),
		new(models.OrderPayment),
		new(models.OrderRefund),
		new(models.PayBill),
		new(models.ShopBanner),
		new(models.ShopHot),
		new(models.ShopHotGoods),
		new(models.ShopHotItem),
		new(models.ShopService),
		new(models.UserAddress),
		new(models.UserCart),
		new(models.UserCollect),
		new(models.UserStore),
	}
	return m
}

func InitDb(db *gorm.DB) error {
	var sqlFileNames []string
	var dirFileNames []string
	root := "./sql"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		path = strings.ReplaceAll(path, "\\", "/")
		if strings.HasSuffix(path, ".sql") {
			if strings.Contains(strings.Trim(path, root), "/") {
				dirFileNames = append(dirFileNames, path)
			} else {
				sqlFileNames = append(sqlFileNames, path)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	sqlFileNames = append(sqlFileNames, dirFileNames...)
	for _, sqlFileName := range sqlFileNames {
		err = utils.ExecSql(db, sqlFileName)
		if err != nil {
			return err
		}
	}
	return err
}
