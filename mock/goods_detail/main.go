package main

import (
	"encoding/json"
	"fmt"
	_const "gitee.com/liujit/shop/mock/const"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type GoodsDetail struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		MainPictures []string `json:"mainPictures"`
		Specs        []struct {
			Name   string `json:"name"`
			Id     string `json:"id"`
			Values []struct {
				Name string `json:"name"`
			} `json:"values"`
		} `json:"specs"`
		Skus []struct {
			Id        string `json:"id"`
			SkuCode   string `json:"skuCode"`
			Price     string `json:"price"`
			OldPrice  string `json:"oldPrice"`
			Inventory int    `json:"inventory"`
			Picture   string `json:"picture"`
			Specs     []struct {
				Name      string `json:"name"`
				ValueName string `json:"valueName"`
			} `json:"specs"`
		} `json:"skus"`
		Categories []struct {
			Id     string `json:"id"`
			Name   string `json:"name"`
			Layer  int    `json:"layer"`
			Parent *struct {
				Id     string      `json:"id"`
				Name   string      `json:"name"`
				Layer  int         `json:"layer"`
				Parent interface{} `json:"parent"`
			} `json:"parent"`
		} `json:"categories"`
		Details struct {
			Pictures   []string `json:"pictures"`
			Properties []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"properties"`
		} `json:"details"`
	} `json:"result"`
}

func main() {

	sqlDb, cleanup, err := sqldb.NewGormSqlDb(_const.Database)
	if err != nil {
		return
	}
	defer cleanup()

	dataData := data.NewData(sqlDb)
	//transaction := data.NewTransaction(dataData)
	m := dataData.Query(_const.Ctx).Goods
	q := m.WithContext(_const.Ctx)
	q = q.Where(m.InitSaleNum.Eq(0))
	goodsList, _, err := q.FindByPage(0, 350)
	if err != nil {
		return
	}

	//goodsPropList := make([]*models.GoodsProp, 0)
	for _, goods := range goodsList {

		url := _const.BaseUrl + fmt.Sprintf("/goods?id=%d", goods.ID)
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
		var goodsDetail GoodsDetail
		err = json.Unmarshal(body, &goodsDetail)
		if err != nil {
			return
		}
		// 更新详情
		if goodsDetail.Code == "1" {
			fmt.Println("查询成功：" + goods.Name)
			goodsDetailRes := goodsDetail.Result
			goodsDetailResDetail := goodsDetailRes.Details
			goodsPropList := make([]*models.GoodsProp, 0)
			for idx, item := range goodsDetailResDetail.Properties {
				goodsPropList = append(goodsPropList, &models.GoodsProp{
					GoodsID: goods.ID,
					Label:   item.Name,
					Value:   item.Value,
					Sort:    int32(idx + 1),
				})
			}
			goodsPropData := dataData.Query(_const.Ctx).GoodsProp
			goodsPropData.WithContext(_const.Ctx).CreateInBatches(goodsPropList, 500)

			goodsSpecMap := make(map[string][]string)

			goods.RealSaleNum = 0
			goods.InitSaleNum = 0

			goodsSpecList := make([]*models.GoodsSpec, 0)

			goodsSkuList := make([]*models.GoodsSku, 0)
			for idx, item := range goodsDetailRes.Skus {
				//time.Sleep(time.Millisecond * 100)
				Price, _ := strconv.ParseFloat(item.Price, 10)
				OldPrice, _ := strconv.ParseFloat(item.OldPrice, 10)

				//picture, err := image.Upload(item.Picture, "goods")
				//if err != nil {
				//	fmt.Println(err)
				picture := item.Picture
				//}
				specItem := make([]string, 0)

				for _, item1 := range item.Specs {
					specItem = append(specItem, item1.ValueName)
					v, ok := goodsSpecMap[item1.Name]
					if !ok {
						v = make([]string, 0)
					}
					v = append(v, item1.ValueName)
					goodsSpecMap[item1.Name] = v
				}
				marshal, err := json.Marshal(specItem)
				if err != nil {
					fmt.Println(err)
				}

				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				randomNumber := r.Intn(10000)

				GoodsSku := &models.GoodsSku{
					GoodsID:       goods.ID,
					Picture:       picture,
					SkuCode:       item.SkuCode,
					SpecItem:      string(marshal),
					Price:         int64(OldPrice * 100),
					DiscountPrice: int64(Price * 100),
					InitSaleNum:   0,
					RealSaleNum:   int64(randomNumber),
					Inventory:     int64(item.Inventory),
				}

				goodsSkuList = append(goodsSkuList, GoodsSku)

				goods.RealSaleNum += GoodsSku.RealSaleNum
				goods.InitSaleNum += GoodsSku.InitSaleNum
				if idx == 0 {
					goods.Price = GoodsSku.Price
					goods.DiscountPrice = GoodsSku.DiscountPrice
				}
			}
			var idx = 0
			for k, v := range goodsSpecMap {
				idx = idx + 1
				marshal, err := json.Marshal(v)
				if err != nil {
					fmt.Println(err)
				}
				goodsSpecList = append(goodsSpecList, &models.GoodsSpec{
					GoodsID: goods.ID,
					Name:    k,
					Item:    string(marshal),
					Sort:    int32(idx),
				})
			}
			goodsSpecData := dataData.Query(_const.Ctx).GoodsSpec
			goodsSpecData.WithContext(_const.Ctx).CreateInBatches(goodsSpecList, 500)

			goodsSkuData := dataData.Query(_const.Ctx).GoodsSku
			goodsSkuData.WithContext(_const.Ctx).CreateInBatches(goodsSkuList, 500)

			detailList := make([]string, 0)
			for _, item := range goodsDetailResDetail.Pictures {
				//picture, err := image.Upload(item, "goods")
				//if err != nil {
				picture := item
				//}
				detailList = append(detailList, picture)
			}
			detail, _ := json.Marshal(detailList)
			goods.Detail = string(detail)
			bannerList := make([]string, 0)
			for _, item := range goodsDetailRes.MainPictures {
				//picture, err := image.Upload(item, "goods")
				//if err != nil {
				picture := item
				//}
				bannerList = append(bannerList, picture)
			}
			banner, _ := json.Marshal(bannerList)
			goods.Banner = string(banner)

			categories := goodsDetailRes.Categories[0]
			categoriesId, err := strconv.Atoi(categories.Id)
			if err != nil {
				fmt.Println(err)
			}
			goods.CategoryID = int64(categoriesId)
			goods.InitSaleNum = 1
			goodsData := dataData.Query(_const.Ctx).Goods
			goodsData.WithContext(_const.Ctx).Save(goods)
		}
	}
}
