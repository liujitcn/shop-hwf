package biz

import (
	"context"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
	"gitee.com/liujit/shop/server/lib/utils/trans"
	authMiddleware "go.newcapec.cn/ncttools/nmskit-auth/middleware"
	"go.newcapec.cn/ncttools/nmskit/log"
)

type UserAddressCase struct {
	tx data.Transaction
	data.UserAddressRepo
	baseAreaCase *BaseAreaCase
}

// NewUserAddressCase new a UserAddress use case.
func NewUserAddressCase(tx data.Transaction,
	userAddressRepo data.UserAddressRepo,
	baseAreaCase *BaseAreaCase,
) *UserAddressCase {
	return &UserAddressCase{
		tx:              tx,
		UserAddressRepo: userAddressRepo,
		baseAreaCase:    baseAreaCase,
	}
}

func (c *UserAddressCase) GetFromID(ctx context.Context, id int64) (*app.UserAddressForm, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	var userAddress *models.UserAddress
	userAddress, err = c.Find(ctx, authInfo.UserId, &data.UserAddressCondition{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return c.convertToProto(ctx, userAddress), nil
}

func (c *UserAddressCase) Create(ctx context.Context, userAddress *app.UserAddressForm) error {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return common.ErrorAccessForbidden("用户认证失败")
	}
	address := c.convertToModel(authInfo.UserId, userAddress)
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		if trans.BoolValue(address.IsDefault) {
			// 修改其他为非默认
			err = c.UserAddressRepo.UpdateByUserId(ctx, authInfo.UserId, &models.UserAddress{
				IsDefault: trans.Bool(false),
			})
			if err != nil {
				return err
			}
		}
		err = c.UserAddressRepo.Create(ctx, address)
		if err != nil {
			return err
		}
		return nil
	})
}

func (c *UserAddressCase) Update(ctx context.Context, userAddress *app.UserAddressForm) error {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return common.ErrorAccessForbidden("用户认证失败")
	}
	address := c.convertToModel(authInfo.UserId, userAddress)

	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		if trans.BoolValue(address.IsDefault) {
			// 修改其他为非默认
			err = c.UserAddressRepo.UpdateByUserId(ctx, authInfo.UserId, &models.UserAddress{
				IsDefault: trans.Bool(false),
			})
			if err != nil {
				return err
			}
		}
		err = c.UpdateByID(ctx, authInfo.UserId, address)
		if err != nil {
			return err
		}
		return nil
	})
}
func (c *UserAddressCase) List(ctx context.Context) (*app.ListUserAddressResponse, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	var all []*models.UserAddress
	all, err = c.FindAll(ctx, authInfo.UserId, &data.UserAddressCondition{})
	if err != nil {
		return nil, err
	}

	list := make([]*app.UserAddress, 0)
	for _, address := range all {
		list = append(list, &app.UserAddress{
			Id:        address.ID,
			Receiver:  address.Receiver,
			Contact:   address.Contact,
			Address:   c.baseAreaCase.GetAddressListByCode(ctx, address.Address),
			Detail:    address.Detail,
			IsDefault: trans.BoolValue(address.IsDefault),
		})
	}
	return &app.ListUserAddressResponse{
		List: list,
	}, nil
}

func (c *UserAddressCase) convertToProto(ctx context.Context, item *models.UserAddress) *app.UserAddressForm {
	return &app.UserAddressForm{
		Id:          item.ID,
		Receiver:    item.Receiver,
		Contact:     item.Contact,
		Address:     str.ConvertJsonStringToStringArray(item.Address),
		Detail:      item.Detail,
		AddressName: c.baseAreaCase.GetAddressListByCode(ctx, item.Address),
		IsDefault:   trans.BoolValue(item.IsDefault),
	}
}

func (c *UserAddressCase) convertToModel(userId int64, item *app.UserAddressForm) *models.UserAddress {
	res := &models.UserAddress{
		ID:        item.GetId(),
		UserID:    userId,
		Receiver:  item.GetReceiver(),
		Contact:   item.GetContact(),
		Address:   str.ConvertStringArrayToString(item.GetAddress()),
		Detail:    item.GetDetail(),
		IsDefault: trans.Bool(item.GetIsDefault()),
	}
	return res
}
