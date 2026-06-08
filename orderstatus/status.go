package orderstatus

import "slices"

// 主订单履约状态，只表达订单整体从下单到完成的履约进度。
// 售后、退款、换货不写入该字段，统一由 OrderAfterSalesStatus 表达。
const (
	OrderStatusCanceled         = -1 // 已取消
	OrderStatusPendingPayment   = 0  // 待支付
	OrderStatusPendingShipment  = 10 // 待发货
	OrderStatusPartShipped      = 20 // 部分发货
	OrderStatusPendingReceipt   = 30 // 待收货
	OrderStatusReceived         = 40 // 已收货
	OrderStatusPendingComment   = 50 // 待评价
	OrderStatusCompleted        = 60 // 已完成
	OrderStatusClosed           = 70 // 已关闭，通常用于全部取消或全量售后关闭履约
	OrderStatusPaymentException = 80 // 支付异常，通常用于已支付但库存等履约条件无法恢复
)

// 主订单支付状态，独立于履约状态。
// 支付成功后业务侧再把订单/明细履约推进到待发货。
const (
	OrderPayStatusFailure = -1 // 支付失败
	OrderPayStatusUnpaid  = 0  // 未支付
	OrderPayStatusSuccess = 1  // 支付成功
)

// 订单支付方式，作为 store_order.pay_type 和相关流水 pay_type 的统一取值。
const (
	OrderPayTypeWechatPay        = "wechat_pay"         // 微信支付
	OrderPayTypeWechatNative     = "wechat_pay_native"  // 微信扫码支付
	OrderPayTypeAlipay           = "alipay"             // 支付宝支付
	OrderPayTypeBalance          = "balance_pay"        // 余额支付
	OrderPayTypeOffline          = "offline_pay"        // 线下支付
	OrderPayTypeWechatFriend     = "wechat_friend_pay"  // 微信好友支付
	OrderPayTypeManualConfirmPay = "manual_confirm_pay" // 后台确认支付
	OrderPayTypeChainPay         = "chain_pay"          // 合约支付
)

// 主订单售后聚合状态，作为 store_order.after_sales_status 的统一取值。
// 该状态是明细数量和售后单流程的聚合结果，不是某一张售后单的原始流程状态。
const (
	OrderAfterSalesStatusNone          = 0  // 无售后
	OrderAfterSalesStatusProcessing    = 10 // 有售后处理中
	OrderAfterSalesStatusPartRefunded  = 20 // 部分退款完成
	OrderAfterSalesStatusAllRefunded   = 30 // 全部明细已退款
	OrderAfterSalesStatusPartExchanged = 40 // 部分换货完成
	OrderAfterSalesStatusAllDone       = 50 // 全部有效明细都有成功售后结果
)

// 订单明细履约状态，只表达单个商品明细的发货、收货、评价进度。
// 明细全量退款后可以进入 Closed，部分退款时履约状态仍保留原进度。
const (
	DetailStatusCanceled        = -1 // 已取消
	DetailStatusPendingPayment  = 0  // 待支付
	DetailStatusPendingShipment = 10 // 待发货
	DetailStatusPendingReceipt  = 20 // 待收货
	DetailStatusReceived        = 30 // 已收货
	DetailStatusPendingComment  = 40 // 待评价
	DetailStatusCompleted       = 50 // 已完成
	DetailStatusClosed          = 60 // 已关闭，不再继续履约
)

// 订单明细售后锁，不作为前端展示状态。
// 它用于限制同一明细同一时间只能有一笔售后处理中。
const (
	DetailAfterSalesStatusNone       = 0 // 无售后锁
	DetailAfterSalesStatusProcessing = 1 // 售后处理中，禁止重复申请
	DetailAfterSalesStatusRefunded   = 2 // 已全量退款
)

// 售后申请类型。web2/web3 共用同一组类型，通道差异放到扩展字段或交易表。
const (
	AfterSalesTypeRefundOnly      = "refund_only"       // 仅退款
	AfterSalesTypeRefundAndReturn = "refund_and_return" // 退货退款
	AfterSalesTypeExchange        = "exchange"          // 换货
)

// 售后流程状态，直接表达一张售后单当前处于哪个业务节点。
// 0/10/20/30/40 是公共节点；100 段是仅退款；200 段是退货退款；300 段是换货。
const (
	AfterSalesStatusDraft                        = 0  // 草稿，常用于 web3 先生成售后单再签名提交
	AfterSalesStatusPendingReview                = 10 // 待商家审核
	AfterSalesStatusRejected                     = 20 // 审核拒绝
	AfterSalesStatusCanceled                     = 30 // 用户取消
	AfterSalesStatusMerchantRejectedPendingAdmin = 40 // 商户拒绝后待平台复核

	AfterSalesStatusRefundPending = 100 // 仅退款待打款
	AfterSalesStatusRefunded      = 110 // 仅退款完成

	AfterSalesStatusReturnPending  = 200 // 退货退款待买家寄回
	AfterSalesStatusReturnShipped  = 210 // 买家已寄回，待商家收货
	AfterSalesStatusReturnReceived = 220 // 商家已收货，待退款
	AfterSalesStatusReturnRefunded = 230 // 退货退款完成

	AfterSalesStatusExchangeReturnPending = 300 // 换货待买家寄回
	AfterSalesStatusExchangeReturnShipped = 310 // 买家已寄回，待商家收货
	AfterSalesStatusExchangeSendPending   = 320 // 商家待寄出换货商品
	AfterSalesStatusExchangeShipped       = 330 // 商家已寄出换货商品
	AfterSalesStatusExchangeCompleted     = 340 // 买家确认或后台强制完成换货
)

// 售后物流类型，用于 after_sales_shipment 区分买家退回和商家换货寄出。
const (
	AfterSalesShipmentTypeBuyerReturn    = "buyer_return"
	AfterSalesShipmentTypeSellerExchange = "seller_exchange"
)

// AfterSalesProcessingStatuses 是会让主订单售后聚合显示 processing 的状态集合。
// 只要同一主订单下存在任意进行中的售后单，聚合状态就优先显示处理中。
var AfterSalesProcessingStatuses = []int{
	AfterSalesStatusPendingReview,
	AfterSalesStatusMerchantRejectedPendingAdmin,
	AfterSalesStatusRefundPending,
	AfterSalesStatusReturnPending,
	AfterSalesStatusReturnShipped,
	AfterSalesStatusReturnReceived,
	AfterSalesStatusExchangeReturnPending,
	AfterSalesStatusExchangeReturnShipped,
	AfterSalesStatusExchangeSendPending,
	AfterSalesStatusExchangeShipped,
}

// AfterSalesFinalStatuses 是售后流程终态集合，包含成功终态和关闭终态。
var AfterSalesFinalStatuses = []int{
	AfterSalesStatusRejected,
	AfterSalesStatusCanceled,
	AfterSalesStatusRefunded,
	AfterSalesStatusReturnRefunded,
	AfterSalesStatusExchangeCompleted,
}

// AfterSalesSuccessFinalStatuses 是会累计退款/换货数量的成功终态集合。
var AfterSalesSuccessFinalStatuses = []int{
	AfterSalesStatusRefunded,
	AfterSalesStatusReturnRefunded,
	AfterSalesStatusExchangeCompleted,
}

// AfterSalesClosedFinalStatuses 是不会累计成功售后数量的关闭终态集合。
var AfterSalesClosedFinalStatuses = []int{
	AfterSalesStatusRejected,
	AfterSalesStatusCanceled,
}

// IsAfterSalesProcessing 判断售后状态是否仍在处理中。
func IsAfterSalesProcessing(status int) bool {
	return slices.Contains(AfterSalesProcessingStatuses, status)
}

// IsAfterSalesSuccessFinal 判断售后状态是否是成功终态。
func IsAfterSalesSuccessFinal(status int) bool {
	return slices.Contains(AfterSalesSuccessFinalStatuses, status)
}

// IsAfterSalesClosedFinal 判断售后状态是否是拒绝、取消等关闭终态。
func IsAfterSalesClosedFinal(status int) bool {
	return slices.Contains(AfterSalesClosedFinalStatuses, status)
}

// ApprovalNextStatus 根据售后类型和是否需要退货，计算商家审核通过后的下一个流程节点。
func ApprovalNextStatus(requestType string, needReturnGoods int) int {
	switch requestType {
	case AfterSalesTypeRefundOnly:
		return AfterSalesStatusRefundPending
	case AfterSalesTypeRefundAndReturn:
		if needReturnGoods == 1 {
			return AfterSalesStatusReturnPending
		}
		return AfterSalesStatusRefundPending
	case AfterSalesTypeExchange:
		if needReturnGoods == 1 {
			return AfterSalesStatusExchangeReturnPending
		}
		return AfterSalesStatusExchangeSendPending
	default:
		return AfterSalesStatusPendingReview
	}
}

// ReturnShippedStatus 返回买家填写退回物流后的下一个状态。
// 换货和退货退款复用“买家寄回”动作，但进入不同状态段。
func ReturnShippedStatus(requestType string) int {
	if requestType == AfterSalesTypeExchange {
		return AfterSalesStatusExchangeReturnShipped
	}
	return AfterSalesStatusReturnShipped
}

// ReturnReceivedStatus 返回商家确认收到退回商品后的下一个状态。
func ReturnReceivedStatus(requestType string) int {
	if requestType == AfterSalesTypeExchange {
		return AfterSalesStatusExchangeSendPending
	}
	return AfterSalesStatusReturnReceived
}

// RefundSuccessStatus 返回退款完成后的终态。
func RefundSuccessStatus(requestType string) int {
	if requestType == AfterSalesTypeRefundAndReturn {
		return AfterSalesStatusReturnRefunded
	}
	return AfterSalesStatusRefunded
}

// CanBuyerReturn 判断当前状态是否允许买家填写退货物流。
func CanBuyerReturn(status int) bool {
	return status == AfterSalesStatusReturnPending || status == AfterSalesStatusExchangeReturnPending
}

// CanRefund 判断当前状态是否允许商家或平台执行退款完成动作。
func CanRefund(status int) bool {
	return status == AfterSalesStatusRefundPending || status == AfterSalesStatusReturnReceived
}

// CanExchangeSend 判断当前状态是否允许商家寄出换货商品。
func CanExchangeSend(status int) bool {
	return status == AfterSalesStatusExchangeSendPending
}

// CanExchangeComplete 判断当前状态是否允许确认换货完成。
func CanExchangeComplete(status int) bool {
	return status == AfterSalesStatusExchangeShipped
}
