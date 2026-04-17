package main

import (
	"encoding/json"
	"fmt"
	_const "gitee.com/liujit/shop/mock/const"
	"gitee.com/liujit/shop/mock/image"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Category struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Result []struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Picture  string `json:"picture"`
		Children []struct {
			Id      string `json:"id"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
		} `json:"children"`
	} `json:"result"`
}

func main() {
	url := _const.BaseUrl + fmt.Sprintf("/category/top")
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var b Category
	err = json.Unmarshal(body, &b)
	if err != nil {
		return
	}

	if b.Code == "1" {
		sqlDb, cleanup, err := sqldb.NewGormSqlDb(_const.Database)
		if err != nil {
			return
		}
		defer cleanup()

		dataData := data.NewData(sqlDb)

		goodsCategoryList := make([]*models.GoodsCategory, 0)
		for i, v := range b.Result {
			time.Sleep(time.Second)
			Id, _ := strconv.Atoi(v.Id)
			picture, err := image.Upload(v.Picture, "category")
			if err != nil {
				picture = v.Picture
			}
			goodsCategory := &models.GoodsCategory{
				ID:        int64(Id),
				ParentID:  0,
				Path:      fmt.Sprintf("/0/%d", Id),
				Picture:   picture,
				Name:      v.Name,
				Sort:      int32(i + 1),
				Status:    1,
				CreatedBy: 1,
				UpdatedBy: 1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			goodsCategoryList = append(goodsCategoryList, goodsCategory)
			if len(v.Children) > 0 {
				for idx, item := range v.Children {
					//time.Sleep(time.Second)
					itemId, _ := strconv.Atoi(item.Id)
					itemPicture, err := image.Upload(v.Picture, "category")
					if err != nil {
						picture = v.Picture
					}
					goodsCategoryList = append(goodsCategoryList, &models.GoodsCategory{
						ID:        int64(itemId),
						ParentID:  goodsCategory.ID,
						Path:      fmt.Sprintf("%s/%d", goodsCategory.Path, itemId),
						Picture:   itemPicture,
						Name:      item.Name,
						Sort:      int32(idx + 1),
						Status:    1,
						CreatedBy: 1,
						UpdatedBy: 1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					})
				}
			}
		}

		q := dataData.Query(_const.Ctx).GoodsCategory
		err = q.WithContext(_const.Ctx).CreateInBatches(goodsCategoryList, 100)
		if err != nil {
			return
		}
	}
}
