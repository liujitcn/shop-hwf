package biz

import (
	"context"
	"fmt"
	"gitee.com/liujit/shop/server/api/app"
	"gitee.com/liujit/shop/server/api/common"
	_const "gitee.com/liujit/shop/server/lib/const"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
	"gitee.com/liujit/shop/server/lib/utils/timeutil"
	"gitee.com/liujit/shop/server/lib/utils/trans"
	"gitee.com/liujit/shop/server/pkg/config"
	"gitee.com/liujit/shop/server/pkg/service/app/util"
	payBiz "gitee.com/liujit/shop/server/pkg/service/pay/biz"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"go.newcapec.cn/nctcommon/nmslib"
	authMiddleware "go.newcapec.cn/ncttools/nmskit-auth/middleware"
	"go.newcapec.cn/ncttools/nmskit/log"
	"strconv"
	"time"
)

type OrderCase struct {
	tx data.Transaction
	data.OrderRepo
	orderCancelCase    *OrderCancelCase
	orderGoodsCase     *OrderGoodsCase
	orderAddressCase   *OrderAddressCase
	orderLogisticsCase *OrderLogisticsCase
	orderPaymentCase   *OrderPaymentCase
	orderRefundCase    *OrderRefundCase
	goodsCase          *GoodsCase
	goodsSkuCase       *GoodsSkuCase
	userAddressCase    *UserAddressCase
	userCartCase       *UserCartCase
	baseDictItemCase   *BaseDictItemCase
	orderSchedulerCase *payBiz.OrderSchedulerCase
	wxPayCase          *payBiz.WxPayCase
}

// NewOrderCase new a Order use case.
func NewOrderCase(
	tx data.Transaction,
	orderRepo data.OrderRepo,
	orderCancelCase *OrderCancelCase,
	orderGoodsCase *OrderGoodsCase,
	orderAddressCase *OrderAddressCase,
	orderLogisticsCase *OrderLogisticsCase,
	orderPaymentCase *OrderPaymentCase,
	orderRefundCase *OrderRefundCase,
	goodsCase *GoodsCase,
	goodsSkuCase *GoodsSkuCase,
	userAddressCase *UserAddressCase,
	userCartCase *UserCartCase,
	baseDictItemCase *BaseDictItemCase,
	orderSchedulerCase *payBiz.OrderSchedulerCase,
	wxPayCase *payBiz.WxPayCase,
) (*OrderCase, error) {
	c := &OrderCase{
		tx:                 tx,
		OrderRepo:          orderRepo,
		orderCancelCase:    orderCancelCase,
		orderGoodsCase:     orderGoodsCase,
		orderAddressCase:   orderAddressCase,
		orderLogisticsCase: orderLogisticsCase,
		orderPaymentCase:   orderPaymentCase,
		orderRefundCase:    orderRefundCase,
		goodsCase:          goodsCase,
		goodsSkuCase:       goodsSkuCase,
		userAddressCase:    userAddressCase,
		userCartCase:       userCartCase,
		baseDictItemCase:   baseDictItemCase,
		orderSchedulerCase: orderSchedulerCase,
		wxPayCase:          wxPayCase,
	}

	// 查询全部未支付订单
	list, err := c.FindAll(context.Background(), &data.OrderCondition{
		Status: int32(common.OrderStatus_CREATED),
	})
	if err != nil {
		return nil, err
	}
	payTimeout := config.ParsePayTimeout()
	for _, item := range list {
		// 支付超时时间
		createdAt := item.CreatedAt.Add(payTimeout)
		nowTime := time.Now()
		countdown := createdAt.Sub(nowTime).Seconds()
		if countdown < 0 {
			// 自动取消订单
			err = c.Cancel(context.Background(), item.UserID, &app.CancelOrderRequest{
				OrderId: item.ID,
			})
			if err != nil {
				return nil, err
			}
		} else {
			// 添加自动取消定时任务
			c.orderSchedulerCase.AddSchedule(item.ID, time.Duration(countdown)*time.Second, func() {
				err = c.Cancel(context.Background(), item.UserID, &app.CancelOrderRequest{
					OrderId: item.ID,
				})
				if err != nil {
					log.Errorf("Cancel order %d failed: %v", item.ID, err)
				}
			})
		}
	}

	return c, nil
}

func (c *OrderCase) GetFromOrderNo(ctx context.Context, orderNo string) (int64, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return 0, common.ErrorAccessForbidden("用户认证失败")
	}

	var item *models.Order
	item, err = c.Find(ctx, &data.OrderCondition{
		OrderNo: orderNo,
		UserId:  authInfo.UserId,
	})
	if err != nil {
		return 0, err
	}
	return item.ID, nil
}

func (c *OrderCase) GetFromID(ctx context.Context, id int64) (*app.OrderResponse, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}

	var item *models.Order
	item, err = c.Find(ctx, &data.OrderCondition{
		Id:     id,
		UserId: authInfo.UserId,
	})
	if err != nil {
		return nil, err
	}
	order := c.ConvertToProto(item)
	// 支付超时时间
	payTimeout := config.ParsePayTimeout()
	createdAt := item.CreatedAt.Add(payTimeout)
	nowTime := time.Now()
	countdown := createdAt.Sub(nowTime).Seconds()

	// 商品
	order.Goods, err = c.orderGoodsCase.ListByOrderId(ctx, order.Id)
	if err != nil {
		return nil, err
	}
	// 地址
	var address *app.OrderResponse_Address
	address, err = c.orderAddressCase.GetFromByOrderId(ctx, order.Id)
	if err != nil {
		return nil, err
	}

	res := app.OrderResponse{
		Order:     order,
		Address:   address,
		Countdown: float32(countdown),
	}

	// 判断订单状态
	switch common.OrderStatus(item.Status) {
	case common.OrderStatus_PAID:
		// 在线支付
		if common.OrderPayType(item.PayType) == common.OrderPayType_ONLINE_PAY {
			var orderPayment *models.OrderPayment
			orderPayment, err = c.orderPaymentCase.Find(ctx, &data.OrderPaymentCondition{
				OrderId: order.Id,
			})
			if err != nil {
				return nil, err
			}
			res.Order.PaymentTime = timeutil.TimeToTimeString(orderPayment.SuccessTime)
		}
	case common.OrderStatus_SHIPPED, common.OrderStatus_RECEIVED:
		// 发货
		var logistics *app.OrderResponse_Logistics
		logistics, err = c.orderLogisticsCase.GetFromByOrderId(ctx, order.Id)
		if err != nil {
			return nil, err
		}
		res.Logistics = logistics
	case common.OrderStatus_CANCELED:
		// 取消
		var orderCancel *models.OrderCancel
		orderCancel, err = c.orderCancelCase.Find(ctx, &data.OrderCancelCondition{
			OrderId: order.Id,
		})
		if err != nil {
			return nil, err
		}
		res.Order.CancelTime = timeutil.TimeToTimeString(orderCancel.CreatedAt)
	case common.OrderStatus_REFUNDING:
		// 退款
		var orderRefund *models.OrderRefund
		orderRefund, err = c.orderRefundCase.Find(ctx, &data.OrderRefundCondition{
			OrderId: order.Id,
		})
		if err != nil {
			return nil, err
		}
		res.Order.RefundTime = timeutil.TimeToTimeString(orderRefund.SuccessTime)
	}
	return &res, nil
}

func (c *OrderCase) Page(ctx context.Context, req *app.PageOrderRequest) (*app.PageOrderResponse, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	condition := &data.OrderCondition{
		UserId: authInfo.UserId,
		Status: int32(req.GetStatus()),
	}
	var page []*models.Order
	var count int64
	page, count, err = c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}

	orderIds := make([]int64, 0)
	for _, item := range page {
		orderIds = append(orderIds, item.ID)
	}

	orderGoodsMap := make(map[int64][]*app.OrderGoods)
	orderGoodsMap, err = c.orderGoodsCase.MapByOrderIds(ctx, orderIds)
	if err != nil {
		return nil, err
	}

	list := make([]*app.Order, 0)
	for _, item := range page {
		order := c.ConvertToProto(item)
		if v, ok := orderGoodsMap[order.Id]; ok {
			order.Goods = v
		}
		list = append(list, order)
	}

	return &app.PageOrderResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *OrderCase) Count(ctx context.Context) (*app.CountOrderResponse, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	var res map[int32]int32
	res, err = c.MapCount(ctx, authInfo.UserId)
	if err != nil {
		return nil, err
	}
	count := make([]*app.CountOrderResponse_Count, 0)
	for k, v := range res {
		count = append(count, &app.CountOrderResponse_Count{
			Status: common.OrderStatus(k),
			Num:    v,
		})
	}
	return &app.CountOrderResponse{
		Count: count,
	}, nil
}

func (c *OrderCase) Create(ctx context.Context, request *app.CreateOrderRequest) (*app.CreateOrderResponse, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	status := int32(common.OrderStatus_CREATED)
	// 货到付款
	if common.OrderPayType(request.PayType) == common.OrderPayType_CASH_ON_DELIVERY {
		status = int32(common.OrderStatus_PAID)
	}

	// 基础信息
	tm := nmslib.Runtime.GetTmGenerate()
	order := &models.Order{
		OrderNo:      strconv.FormatInt(tm.NextVal(), 10),
		UserID:       authInfo.UserId,
		PayType:      int32(request.PayType),
		PayChannel:   int32(request.PayChannel),
		DeliveryTime: int32(request.DeliveryTime),
		Status:       status,
		Remark:       request.Remark,
	}
	err = c.tx.Transaction(ctx, func(ctx context.Context) error {
		var orderGoodsList []*models.OrderGoods
		orderGoodsList, err = c.orderGoodsCase.ConvertToModelList(ctx, request.GetGoods())
		if err != nil {
			return err
		}
		for _, orderGoods := range orderGoodsList {
			order.PayMoney += orderGoods.TotalPayPrice
			order.TotalMoney += orderGoods.TotalPrice
			order.GoodsNum += orderGoods.Num
			// 增加销量，减库存
			err = c.goodsCase.AddSaleNum(ctx, orderGoods.GoodsID, orderGoods.Num)
			if err != nil {
				return err
			}
			err = c.goodsSkuCase.AddSaleNum(ctx, orderGoods.SkuCode, orderGoods.Num)
			if err != nil {
				return err
			}
			// 删除购物车
			err = c.userCartCase.DeleteByGoodsId(ctx, authInfo.UserId, orderGoods.GoodsID, orderGoods.SkuCode)
			if err != nil {
				return err
			}
		}
		// 运费，默认都是0
		order.PostFee = 0
		err = c.OrderRepo.Create(ctx, order)
		if err != nil {
			return err
		}

		// 商品
		err = c.orderGoodsCase.BatchCreate(ctx, order.ID, orderGoodsList)
		if err != nil {
			return err
		}
		// 地址
		err = c.orderAddressCase.CreateByAddressId(ctx, authInfo.UserId, order.ID, request.GetAddressId())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 添加自动取消定时任务
	if order.Status == int32(common.OrderStatus_CREATED) {
		// 支付超时时间
		payTimeout := config.ParsePayTimeout()
		createdAt := order.CreatedAt.Add(payTimeout)
		nowTime := time.Now()
		countdown := createdAt.Sub(nowTime).Seconds()
		c.orderSchedulerCase.AddSchedule(order.ID, time.Duration(countdown)*time.Second, func() {
			err = c.Cancel(context.Background(), order.UserID, &app.CancelOrderRequest{
				OrderId: order.ID,
			})
			if err != nil {
				log.Errorf("Cancel order %d failed: %v", order.ID, err)
			}
		})
	}
	return &app.CreateOrderResponse{
		OrderId: order.ID,
	}, nil
}

func (c *OrderCase) Cancel(ctx context.Context, userId int64, req *app.CancelOrderRequest) error {
	order, err := c.Find(ctx, &data.OrderCondition{
		Id:     req.GetOrderId(),
		UserId: userId,
	})
	if err != nil {
		return err
	}
	if order.Status != int32(common.OrderStatus_CREATED) {
		return fmt.Errorf("订单状态错误：【%s】", common.OrderStatus_name[order.Status])
	}
	ids := []int64{req.GetOrderId()}
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		var orderGoodsList []*models.OrderGoods
		orderGoodsList, err = c.orderGoodsCase.FindAll(ctx, &data.OrderGoodsCondition{
			OrderIds: ids,
		})
		for _, orderGoods := range orderGoodsList {
			// 增加库存，减销量
			err = c.goodsCase.SubSaleNum(ctx, orderGoods.GoodsID, orderGoods.Num)
			if err != nil {
				return err
			}
			err = c.goodsSkuCase.SubSaleNum(ctx, orderGoods.SkuCode, orderGoods.Num)
			if err != nil {
				return err
			}
		}
		// 增加取消信息
		err = c.orderCancelCase.Create(ctx, &models.OrderCancel{
			OrderID: req.GetOrderId(),
			Reason:  int32(req.GetReason()),
		})
		if err != nil {
			return err
		}
		return c.UpdateByIds(ctx, userId, ids, &models.Order{
			Status: int32(common.OrderStatus_CANCELED),
		})
	})
}

func (c *OrderCase) Refund(ctx context.Context, req *app.RefundOrderRequest) error {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return common.ErrorAccessForbidden("用户认证失败")
	}
	var order *models.Order
	order, err = c.Find(ctx, &data.OrderCondition{
		Id:     req.GetOrderId(),
		UserId: authInfo.UserId,
	})
	if err != nil {
		return err
	}

	if order.Status != int32(common.OrderStatus_PAID) {
		return fmt.Errorf("订单状态错误：【%s】", common.OrderStatus_name[order.Status])
	}

	tm := nmslib.Runtime.GetTmGenerate()
	orderRefund := &models.OrderRefund{
		OrderID:  req.GetOrderId(),
		RefundNo: strconv.FormatInt(tm.NextVal(), 10),
		Reason:   int32(req.GetReason()),
	}
	// 在线支付
	if common.OrderPayType(order.PayType) == common.OrderPayType_ONLINE_PAY {
		// 查询退款原因
		reason := strconv.Itoa(int(orderRefund.Reason))
		var label string
		label, err = c.baseDictItemCase.FindLabelByCodeAndValue(ctx, _const.BaseDictCode_OrderRefundReason, reason)
		if err == nil {
			reason = label
		}
		// 微信支付
		if common.OrderPayChannel(order.PayChannel) == common.OrderPayChannel_WX_PAY {
			var refund *refunddomestic.Refund
			refund, err = c.wxPayCase.Refund(refunddomestic.CreateRequest{
				OutTradeNo:  trans.String(order.OrderNo),
				OutRefundNo: trans.String(orderRefund.RefundNo),
				Reason:      trans.String(reason),
				Amount: &refunddomestic.AmountReq{
					Total:    trans.Int64(order.PayMoney),
					Refund:   trans.Int64(order.PayMoney),
					Currency: trans.String("CNY"),
				},
			})
			if err != nil {
				return err
			}
			orderRefund.OrderNo = trans.StringValue(refund.OutTradeNo)
			orderRefund.ThirdOrderNo = trans.StringValue(refund.TransactionId)
			orderRefund.ThirdRefundNo = trans.StringValue(refund.RefundId)
			orderRefund.Channel = string(*refund.Channel.Ptr())
			orderRefund.UserReceivedAccount = trans.StringValue(refund.UserReceivedAccount)
			orderRefund.CreateTime = trans.TimeValue(refund.CreateTime)
			orderRefund.SuccessTime = trans.TimeValue(refund.SuccessTime)
			orderRefund.RefundState = string(*refund.Status.Ptr())
			orderRefund.FundsAccount = string(*refund.FundsAccount)
			orderRefund.Amount = str.ConvertAnyToJsonString(refund.Amount)
		}
	} else {
		t := time.Now()
		orderRefund.CreateTime = t
		orderRefund.SuccessTime = t
		orderRefund.Amount = "{}"
	}
	ids := []int64{req.GetOrderId()}
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		// 增加退款信息
		err = c.orderRefundCase.Create(ctx, orderRefund)
		if err != nil {
			return err
		}
		return c.UpdateByIds(ctx, authInfo.UserId, ids, &models.Order{
			Status: int32(common.OrderStatus_REFUNDING),
		})
	})
}

func (c *OrderCase) Receive(ctx context.Context, req *app.ReceiveOrderRequest) error {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return common.ErrorAccessForbidden("用户认证失败")
	}

	var order *models.Order
	order, err = c.Find(ctx, &data.OrderCondition{
		Id:     req.GetOrderId(),
		UserId: authInfo.UserId,
	})
	if err != nil {
		return err
	}
	if order.Status != int32(common.OrderStatus_SHIPPED) {
		return fmt.Errorf("订单状态错误：【%s】", common.OrderStatus_name[order.Status])
	}

	ids := []int64{req.GetOrderId()}
	return c.UpdateByIds(ctx, authInfo.UserId, ids, &models.Order{
		Status: int32(common.OrderStatus_RECEIVED),
	})
}

func (c *OrderCase) Delete(ctx context.Context, id int64) error {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return common.ErrorAccessForbidden("用户认证失败")
	}

	var order *models.Order
	order, err = c.Find(ctx, &data.OrderCondition{
		Id:     id,
		UserId: authInfo.UserId,
	})
	if err != nil {
		return err
	}
	if !(order.Status == int32(common.OrderStatus_RECEIVED) || order.Status == int32(common.OrderStatus_REFUNDING) || order.Status == int32(common.OrderStatus_CANCELED)) {
		return fmt.Errorf("订单状态错误：【%s】", common.OrderStatus_name[order.Status])
	}

	ids := []int64{id}
	return c.UpdateByIds(ctx, authInfo.UserId, ids, &models.Order{
		Status: int32(common.OrderStatus_DELETED),
	})
}

func (c *OrderCase) Confirm(ctx context.Context, member bool, createOrderGoods []*app.CreateOrderGoods) (*app.ConfirmOrderResponse, error) {
	newOrderGoods := make([]*app.OrderGoods, 0)
	for _, item := range createOrderGoods {
		newGoods, err := c.orderGoodsCase.ConvertToProtoByCreateOrderGoods(ctx, member, item)
		if err != nil {
			return nil, err
		}
		newOrderGoods = append(newOrderGoods, newGoods)
	}

	var summary app.OrderSummary
	for _, orderGoods := range newOrderGoods {
		summary.PayMoney += orderGoods.TotalPayPrice
		summary.TotalMoney += orderGoods.TotalPrice
		summary.GoodsNum += orderGoods.Num
	}
	// 运费 默认0
	summary.PostFee = 0
	return &app.ConfirmOrderResponse{
		Goods:   newOrderGoods,
		Summary: &summary,
	}, nil
}

func (c *OrderCase) Pre(ctx context.Context) (*app.ConfirmOrderResponse, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	member := util.IsMemberByAuthInfo(authInfo)
	var all []*models.UserCart
	all, err = c.userCartCase.FindAll(ctx, authInfo.UserId, &data.UserCartCondition{
		IsChecked: true,
	})
	if err != nil {
		return nil, err
	}
	createOrderGoods := make([]*app.CreateOrderGoods, 0)
	for _, item := range all {
		createOrderGoods = append(createOrderGoods, &app.CreateOrderGoods{
			GoodsId: item.GoodsID,
			SkuCode: item.SkuCode,
			Num:     item.Num,
		})
	}
	return c.Confirm(ctx, member, createOrderGoods)
}

func (c *OrderCase) Repurchase(ctx context.Context, req *app.OrderRepurchaseRequest) (*app.ConfirmOrderResponse, error) {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	member := util.IsMemberByAuthInfo(authInfo)
	var order *models.Order
	order, err = c.Find(ctx, &data.OrderCondition{
		Id:     req.GetOrderId(),
		UserId: authInfo.UserId,
	})
	if err != nil {
		return nil, err
	}
	// 商品
	var oldOrderGoods []*models.OrderGoods
	oldOrderGoods, err = c.orderGoodsCase.FindAll(ctx, &data.OrderGoodsCondition{
		OrderId: order.ID,
	})
	if err != nil {
		return nil, err
	}
	createOrderGoods := make([]*app.CreateOrderGoods, 0)
	for _, item := range oldOrderGoods {
		createOrderGoods = append(createOrderGoods, &app.CreateOrderGoods{
			GoodsId: item.GoodsID,
			SkuCode: item.SkuCode,
			Num:     item.Num,
		})
	}
	return c.Confirm(ctx, member, createOrderGoods)
}

func (c *OrderCase) ConvertToProto(item *models.Order) *app.Order {
	res := &app.Order{
		Id:           item.ID,
		OrderNo:      item.OrderNo,
		PayMoney:     item.PayMoney,
		TotalMoney:   item.TotalMoney,
		PostFee:      item.PostFee,
		GoodsNum:     item.GoodsNum,
		PayType:      common.OrderPayType(item.PayType),
		PayChannel:   common.OrderPayChannel(item.PayChannel),
		DeliveryTime: common.OrderDeliveryTime(item.DeliveryTime),
		Status:       common.OrderStatus(item.Status),
		Remark:       item.Remark,
		CreatedAt:    timeutil.TimeToTimeString(item.CreatedAt),
		UpdatedAt:    timeutil.TimeToTimeString(item.UpdatedAt),
	}
	return res
}
