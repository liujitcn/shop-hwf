package main

import (
	"encoding/json"
	"fmt"
	_const "gitee.com/liujit/shop/mock/const"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Hot struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		Id            string `json:"id"`
		BannerPicture string `json:"bannerPicture"`
		Title         string `json:"title"`
		SubTypes      []struct {
			Id         string `json:"id"`
			Title      string `json:"title"`
			GoodsItems struct {
				Counts   int `json:"counts"`
				PageSize int `json:"pageSize"`
				Pages    int `json:"pages"`
				Page     int `json:"page"`
				Items    []struct {
					Id       string  `json:"id"`
					Name     string  `json:"name"`
					Desc     *string `json:"desc"`
					Price    string  `json:"price"`
					Picture  string  `json:"picture"`
					OrderNum int     `json:"orderNum"`
				} `json:"items"`
			} `json:"goodsItems,omitempty"`
		} `json:"subTypes"`
	} `json:"result"`
}

func main() {

	urlMap := make(map[int64]string)
	urlMap[897682543] = "/hot/preference"
	urlMap[896807072] = "/hot/inVogue"
	urlMap[625755297] = "/hot/oneStop"
	urlMap[814787192] = "/hot/new"

	sqlDb, cleanup, err := sqldb.NewGormSqlDb(_const.Database)
	if err != nil {
		return
	}
	defer cleanup()

	dataData := data.NewData(sqlDb)

	for k, v := range urlMap {
		m := dataData.Query(_const.Ctx).ShopHotItem
		q := m.WithContext(_const.Ctx)
		q = q.Where(m.HotID.Eq(k))
		list, _ := q.Find()
		for _, item := range list {
			url := _const.BaseUrl + fmt.Sprintf("%s?subType=%d&page=1&pageSize=1000", v, item.ID)
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
			var b Hot
			err = json.Unmarshal(body, &b)
			if err != nil {
				return
			}

			if b.Code == "1" {
				SubTypes := b.Result.SubTypes
				for _, subType := range SubTypes {
					id, _ := strconv.Atoi(subType.Id)
					goods := subType.GoodsItems.Items
					if len(goods) > 0 {
						for i, good := range goods {
							goodsId, _ := strconv.Atoi(good.Id)
							s := &models.ShopHotGoodsInfo{
								HotItemID: int64(id),
								GoodsID:   int64(goodsId),
								Sort:      int32(i + 1),
							}
							ss := fmt.Sprintf("INSERT INTO `nmskit`.`shop_hot_goods_info` (`hot_item_id`, `goods_id`, `sort`) VALUES (%d, %d, %d);", s.HotItemID, s.GoodsID, s.Sort)
							fmt.Println(ss)
						}
					}
				}
			}
		}
	}
}
