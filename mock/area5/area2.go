package main

import (
	_const "gitee.com/liujit/shop/mock/const"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data/models"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
	"strconv"
)

func main() {
	sqlDb, cleanup, err := sqldb.NewGormSqlDb(_const.Database)
	if err != nil {
		return
	}
	defer cleanup()

	db := sqlDb.GetDb()

	oldList := make([]*models.BaseArea, 0)
	err = db.Find(&oldList).Error
	if err != nil {
		return
	}

	list := buildTree(oldList, 0)

	newList := make([]*models.BaseAreaNew, 0)

	for _, v1 := range list {
		v1Id, _ := strconv.Atoi(v1.Value)
		if len(v1.Children) > 0 {
			for _, v2 := range v1.Children {
				v2Id, _ := strconv.Atoi(v2.Value)
				if len(v2.Children) > 0 {
					for _, v3 := range v2.Children {
						v3Id, _ := strconv.Atoi(v3.Value)
						newList = append(newList, &models.BaseAreaNew{
							ID:       int64(v3Id),
							ParentID: int64(v2Id),
							Name:     v3.Text,
						})
					}
				}
				newList = append(newList, &models.BaseAreaNew{
					ID:       int64(v2Id),
					ParentID: int64(v1Id),
					Name:     v2.Text,
				})
			}
		}
		newList = append(newList, &models.BaseAreaNew{
			ID:       int64(v1Id),
			ParentID: 0,
			Name:     v1.Text,
		})
	}

	err = db.CreateInBatches(newList, 100).Error
	if err != nil {
		return
	}
}

// buildTree 构建行政区域树状
func buildTree(list []*models.BaseArea, parentId int64) []*common.AppTreeOptionResponse_Option {
	var res []*common.AppTreeOptionResponse_Option
	for _, item := range list {
		if item.ParentID == parentId {
			option := &common.AppTreeOptionResponse_Option{
				Value:    item.AreaCode,
				Text:     item.Name,
				Selected: false,
				Disable:  false,
			}
			option.Children = buildTree(list, item.ID)
			res = append(res, option)
		}
	}
	return res
}
