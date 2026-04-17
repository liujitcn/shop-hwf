package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/cache"
	"strconv"
	"sync"
)

var tree *common.AppTreeOptionResponse
var codeMap map[string]string
var lock sync.RWMutex

type BaseAreaCase struct {
	data.BaseAreaRepo
	cache cache.Cache
}

// NewBaseAreaCase new a BaseArea use case.
func NewBaseAreaCase(
	baseAreaRepo data.BaseAreaRepo,
) *BaseAreaCase {
	return &BaseAreaCase{
		BaseAreaRepo: baseAreaRepo,
	}
}

func (c *BaseAreaCase) Tree(ctx context.Context) (*common.AppTreeOptionResponse, error) {
	lock.RLock()
	defer lock.RUnlock()
	if tree == nil {
		list, err := c.FindAll(ctx, &data.BaseAreaCondition{})
		if err != nil {
			return nil, err
		}
		tree = &common.AppTreeOptionResponse{
			List: c.buildTree(list, 0),
		}
	}
	return tree, nil
}

func (c *BaseAreaCase) GetAddressByCode(ctx context.Context, code string) string {
	res := c.GetAddressListByCode(ctx, code)
	return str.ConvertStringArrayToString(res)
}

func (c *BaseAreaCase) GetAddressListByCode(ctx context.Context, code string) []string {
	lock.RLock()
	defer lock.RUnlock()
	res := make([]string, 0)
	if codeMap == nil {
		list, err := c.FindAll(ctx, &data.BaseAreaCondition{})
		if err != nil {
			return res
		}
		codeMap = make(map[string]string)
		for _, item := range list {
			codeMap[strconv.FormatInt(item.ID, 10)] = item.Name
		}
	}
	codeList := str.ConvertJsonStringToStringArray(code)
	for _, item := range codeList {
		if v, ok := codeMap[item]; ok {
			res = append(res, v)
		} else {
			res = append(res, item)
		}
	}
	return res
}

// buildTree 构建行政区域树状
func (c *BaseAreaCase) buildTree(list []*models.BaseArea, parentId int64) []*common.AppTreeOptionResponse_Option {
	var res []*common.AppTreeOptionResponse_Option
	for _, item := range list {
		if item.ParentID == parentId {
			option := &common.AppTreeOptionResponse_Option{
				Value:    strconv.FormatInt(item.ID, 10),
				Text:     item.Name,
				Selected: false,
				Disable:  false,
			}
			option.Children = c.buildTree(list, item.ID)
			res = append(res, option)
		}
	}
	return res
}
