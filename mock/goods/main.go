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

type Goods struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		Counts   int `json:"counts"`
		PageSize int `json:"pageSize"`
		Pages    int `json:"pages"`
		Page     int `json:"page"`
		Items    []struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			Desc     string `json:"desc"`
			Price    string `json:"price"`
			Picture  string `json:"picture"`
			OrderNum int    `json:"orderNum"`
		} `json:"items"`
	} `json:"result"`
}

func main() {

	url := _const.BaseUrl + fmt.Sprintf("/home/goods/guessLike?page=1&pageSize=500")
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
	var b Goods
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

		goodsList := make([]*models.Goods, 0)
		for _, v := range b.Result.Items {
			//time.Sleep(time.Second)
			Id, _ := strconv.Atoi(v.Id)
			picture, err := image.Upload(v.Picture, "goods")
			if err != nil {
				picture = v.Picture
			}
			goods := &models.Goods{
				ID:            int64(Id),
				CategoryID:    0,
				Name:          v.Name,
				Desc:          v.Desc,
				Picture:       picture,
				Banner:        "[]",
				Detail:        "[]",
				Price:         0,
				DiscountPrice: 0,
				InitSaleNum:   0,
				RealSaleNum:   0,
				Status:        1,
				CreatedBy:     1,
				UpdatedBy:     1,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			goodsList = append(goodsList, goods)
		}

		q := dataData.Query(_const.Ctx).Goods
		err = q.WithContext(_const.Ctx).CreateInBatches(goodsList, 100)
		if err != nil {
			return
		}
	}
}
