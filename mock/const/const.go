package _const

import (
	"context"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/conf"
)

const BaseUrl = "https://pcapi-xiaotuxian-front.itheima.net"

var Ctx = context.Background()

var Database = &conf.Data_Database{
	Driver: "mysql",
	Source: "root:112233@tcp(127.0.0.1:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms",
}
