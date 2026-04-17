package image

import (
	"errors"
	"fmt"
	"go.newcapec.cn/nctcommon/nmslib/util"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/conf"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/oss/local"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

func Upload(rawURL string, biz string) (string, error) {
	return "", errors.New("11111")
	if rawURL == "" {
		return "", nil
	}

	if strings.HasPrefix(rawURL, "shop") {
		return rawURL, nil
	}

	client := &http.Client{
		Timeout: time.Second * 3,
	}
	req, err := http.NewRequest("GET", rawURL, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	var okImgData []byte
	okImgData, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	oss := local.NewOSS(&conf.OSS_Local{
		RootDirectory: "./upload",
	})

	// 解析URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", nil
	}

	// 获取路径部分
	path := parsedURL.Path

	// 使用filepath.Base获取路径的最后一级
	base := filepath.Base(path)

	ext := filepath.Ext(base)

	tmGenerate, _ := util.NewTmGenerate()
	datePath := time.Now().Format("2006/01/02")

	fileName := fmt.Sprintf("%d%s", tmGenerate.NextVal(), ext)

	return oss.UploadByByte(fileName, fmt.Sprintf("shop/%s/%s", biz, datePath), okImgData)
}
