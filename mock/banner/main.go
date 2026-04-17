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

type Banner struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Result []struct {
		Id      string `json:"id"`
		ImgUrl  string `json:"imgUrl"`
		HrefUrl string `json:"hrefUrl"`
		Type    string `json:"type"`
	} `json:"result"`
}

func main() {
	var site = 2

	url := _const.BaseUrl + fmt.Sprintf("/home/banner?distributionSite=%d", site)
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
	var b Banner
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

		shopBannerList := make([]*models.ShopBanner, 0)
		for i, v := range b.Result {
			//time.Sleep(time.Second)
			Type, _ := strconv.Atoi(v.Type)

			picture, err := image.Upload(v.ImgUrl, "banner")
			if err != nil {
				picture = v.ImgUrl
			}

			shopBannerList = append(shopBannerList, &models.ShopBanner{
				Site:      int32(site),
				Picture:   picture,
				Type:      int32(Type),
				Href:      v.HrefUrl,
				Sort:      int32(i + 1),
				Status:    1,
				CreatedBy: 1,
				UpdatedBy: 1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}

		q := dataData.Query(_const.Ctx).ShopBanner
		err = q.WithContext(_const.Ctx).CreateInBatches(shopBannerList, 100)
		if err != nil {
			return
		}
	}
}
