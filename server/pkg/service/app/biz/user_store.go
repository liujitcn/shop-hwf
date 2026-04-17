package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
)

type UserStoreCase struct {
	tx data.Transaction
	data.UserStoreRepo
	baseAreaCase *BaseAreaCase
}

// NewUserStoreCase new a UserStore use case.
func NewUserStoreCase(tx data.Transaction,
	userStoreRepo data.UserStoreRepo,
	baseAreaCase *BaseAreaCase,
) *UserStoreCase {
	return &UserStoreCase{
		tx:            tx,
		UserStoreRepo: userStoreRepo,
		baseAreaCase:  baseAreaCase,
	}
}

func (c *UserStoreCase) ConvertToProto(ctx context.Context, item *models.UserStore) *app.UserStore {
	return &app.UserStore{
		Id:              item.ID,
		Name:            item.Name,
		Address:         str.ConvertJsonStringToStringArray(item.Address),
		Detail:          item.Detail,
		Picture:         str.ConvertJsonStringToStringArray(item.Picture),
		BusinessLicense: str.ConvertJsonStringToStringArray(item.BusinessLicense),
		AddressName:     c.baseAreaCase.GetAddressListByCode(ctx, item.Address),
		Status:          common.UserStoreStatus(item.Status),
		Remark:          item.Remark,
	}
}

func (c *UserStoreCase) ConvertToModel(userId int64, item *app.UserStoreForm) *models.UserStore {
	res := &models.UserStore{
		ID:              item.GetId(),
		UserID:          userId,
		Name:            item.GetName(),
		Address:         str.ConvertStringArrayToString(item.GetAddress()),
		Detail:          item.GetDetail(),
		Picture:         str.ConvertStringArrayToString(item.GetPicture()),
		BusinessLicense: str.ConvertStringArrayToString(item.GetBusinessLicense()),
		Status:          int32(common.UserStoreStatus_PENDING_REVIEW),
	}
	return res
}
