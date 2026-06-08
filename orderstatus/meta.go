package orderstatus

// StatusMeta 是接口层展示状态的统一结构。
// Code 对应数据库数字状态，Status 是前端稳定使用的字符串状态，Text 是默认中文文案。
type StatusMeta struct {
	Code   int
	Status string
	Text   string
}

// 主订单履约状态的接口展示映射。
var orderStatusMeta = map[int]StatusMeta{
	OrderStatusCanceled:         {Code: OrderStatusCanceled, Status: "canceled", Text: "已取消"},
	OrderStatusPendingPayment:   {Code: OrderStatusPendingPayment, Status: "pending_payment", Text: "待支付"},
	OrderStatusPendingShipment:  {Code: OrderStatusPendingShipment, Status: "pending_shipment", Text: "待发货"},
	OrderStatusPartShipped:      {Code: OrderStatusPartShipped, Status: "part_shipped", Text: "部分发货"},
	OrderStatusPendingReceipt:   {Code: OrderStatusPendingReceipt, Status: "pending_receipt", Text: "待收货"},
	OrderStatusReceived:         {Code: OrderStatusReceived, Status: "received", Text: "已收货"},
	OrderStatusPendingComment:   {Code: OrderStatusPendingComment, Status: "pending_comment", Text: "待评价"},
	OrderStatusCompleted:        {Code: OrderStatusCompleted, Status: "completed", Text: "已完成"},
	OrderStatusClosed:           {Code: OrderStatusClosed, Status: "closed", Text: "已关闭"},
	OrderStatusPaymentException: {Code: OrderStatusPaymentException, Status: "payment_exception", Text: "支付异常"},
}

// 主订单售后聚合状态的接口展示映射。
var orderAfterSalesStatusMeta = map[int]StatusMeta{
	OrderAfterSalesStatusNone:          {Code: OrderAfterSalesStatusNone, Status: "none", Text: "无售后"},
	OrderAfterSalesStatusProcessing:    {Code: OrderAfterSalesStatusProcessing, Status: "processing", Text: "售后处理中"},
	OrderAfterSalesStatusPartRefunded:  {Code: OrderAfterSalesStatusPartRefunded, Status: "part_refunded", Text: "部分退款完成"},
	OrderAfterSalesStatusAllRefunded:   {Code: OrderAfterSalesStatusAllRefunded, Status: "all_refunded", Text: "全部明细已退款"},
	OrderAfterSalesStatusPartExchanged: {Code: OrderAfterSalesStatusPartExchanged, Status: "part_exchanged", Text: "部分换货完成"},
	OrderAfterSalesStatusAllDone:       {Code: OrderAfterSalesStatusAllDone, Status: "all_done", Text: "全部明细售后完成"},
}

// 订单明细履约状态的接口展示映射。
var detailStatusMeta = map[int]StatusMeta{
	DetailStatusCanceled:        {Code: DetailStatusCanceled, Status: "canceled", Text: "已取消"},
	DetailStatusPendingPayment:  {Code: DetailStatusPendingPayment, Status: "pending_payment", Text: "待支付"},
	DetailStatusPendingShipment: {Code: DetailStatusPendingShipment, Status: "pending_shipment", Text: "待发货"},
	DetailStatusPendingReceipt:  {Code: DetailStatusPendingReceipt, Status: "pending_receipt", Text: "待收货"},
	DetailStatusReceived:        {Code: DetailStatusReceived, Status: "received", Text: "已收货"},
	DetailStatusPendingComment:  {Code: DetailStatusPendingComment, Status: "pending_comment", Text: "待评价"},
	DetailStatusCompleted:       {Code: DetailStatusCompleted, Status: "completed", Text: "已完成"},
	DetailStatusClosed:          {Code: DetailStatusClosed, Status: "closed", Text: "已关闭"},
}

// 订单明细售后锁的接口展示映射。
// 注意：该状态主要用于锁定重复售后，通常不直接作为前端主展示状态。
var detailAfterSalesStatusMeta = map[int]StatusMeta{
	DetailAfterSalesStatusNone:       {Code: DetailAfterSalesStatusNone, Status: "none", Text: "无售后"},
	DetailAfterSalesStatusProcessing: {Code: DetailAfterSalesStatusProcessing, Status: "processing", Text: "售后中"},
	DetailAfterSalesStatusRefunded:   {Code: DetailAfterSalesStatusRefunded, Status: "refunded", Text: "已退款"},
}

// 售后流程状态的接口展示映射。
var afterSalesStatusMeta = map[int]StatusMeta{
	AfterSalesStatusDraft:                        {Code: AfterSalesStatusDraft, Status: "draft", Text: "草稿，待提交"},
	AfterSalesStatusPendingReview:                {Code: AfterSalesStatusPendingReview, Status: "pending_review", Text: "待商家审核"},
	AfterSalesStatusRejected:                     {Code: AfterSalesStatusRejected, Status: "rejected", Text: "已拒绝"},
	AfterSalesStatusCanceled:                     {Code: AfterSalesStatusCanceled, Status: "canceled", Text: "已取消"},
	AfterSalesStatusMerchantRejectedPendingAdmin: {Code: AfterSalesStatusMerchantRejectedPendingAdmin, Status: "merchant_rejected_pending_admin", Text: "商户拒绝，待平台复核"},
	AfterSalesStatusRefundPending:                {Code: AfterSalesStatusRefundPending, Status: "refund_pending", Text: "待退款"},
	AfterSalesStatusRefunded:                     {Code: AfterSalesStatusRefunded, Status: "refunded", Text: "退款完成"},
	AfterSalesStatusReturnPending:                {Code: AfterSalesStatusReturnPending, Status: "return_pending", Text: "待买家寄回"},
	AfterSalesStatusReturnShipped:                {Code: AfterSalesStatusReturnShipped, Status: "return_shipped", Text: "买家已寄回，待商家收货"},
	AfterSalesStatusReturnReceived:               {Code: AfterSalesStatusReturnReceived, Status: "return_received", Text: "商家已收货，待退款"},
	AfterSalesStatusReturnRefunded:               {Code: AfterSalesStatusReturnRefunded, Status: "return_refunded", Text: "退货退款完成"},
	AfterSalesStatusExchangeReturnPending:        {Code: AfterSalesStatusExchangeReturnPending, Status: "exchange_return_pending", Text: "待买家寄回"},
	AfterSalesStatusExchangeReturnShipped:        {Code: AfterSalesStatusExchangeReturnShipped, Status: "exchange_return_shipped", Text: "买家已寄回，待商家收货"},
	AfterSalesStatusExchangeSendPending:          {Code: AfterSalesStatusExchangeSendPending, Status: "exchange_send_pending", Text: "待商家寄出换货"},
	AfterSalesStatusExchangeShipped:              {Code: AfterSalesStatusExchangeShipped, Status: "exchange_shipped", Text: "换货已寄出"},
	AfterSalesStatusExchangeCompleted:            {Code: AfterSalesStatusExchangeCompleted, Status: "exchange_completed", Text: "换货完成"},
}

func statusCodeByName(metaMap map[int]StatusMeta, status string) (int, bool) {
	for _, meta := range metaMap {
		if meta.Status == status {
			return meta.Code, true
		}
	}
	return 0, false
}

// OrderStatusMeta 返回主订单履约状态的展示信息。
func OrderStatusMeta(code int) (StatusMeta, bool) {
	meta, ok := orderStatusMeta[code]
	return meta, ok
}

// OrderStatusCode 返回主订单履约字符串状态对应的数字状态。
func OrderStatusCode(status string) (int, bool) {
	return statusCodeByName(orderStatusMeta, status)
}

// OrderAfterSalesStatusMeta 返回主订单售后聚合状态的展示信息。
func OrderAfterSalesStatusMeta(code int) (StatusMeta, bool) {
	meta, ok := orderAfterSalesStatusMeta[code]
	return meta, ok
}

// OrderAfterSalesStatusCode 返回主订单售后聚合字符串状态对应的数字状态。
func OrderAfterSalesStatusCode(status string) (int, bool) {
	return statusCodeByName(orderAfterSalesStatusMeta, status)
}

// DetailStatusMeta 返回订单明细履约状态的展示信息。
func DetailStatusMeta(code int) (StatusMeta, bool) {
	meta, ok := detailStatusMeta[code]
	return meta, ok
}

// DetailStatusCode 返回订单明细履约字符串状态对应的数字状态。
func DetailStatusCode(status string) (int, bool) {
	return statusCodeByName(detailStatusMeta, status)
}

// DetailAfterSalesStatusMeta 返回订单明细售后锁状态的展示信息。
func DetailAfterSalesStatusMeta(code int) (StatusMeta, bool) {
	meta, ok := detailAfterSalesStatusMeta[code]
	return meta, ok
}

// DetailAfterSalesStatusCode 返回订单明细售后锁字符串状态对应的数字状态。
func DetailAfterSalesStatusCode(status string) (int, bool) {
	return statusCodeByName(detailAfterSalesStatusMeta, status)
}

// AfterSalesStatusMeta 返回售后流程状态的展示信息。
func AfterSalesStatusMeta(code int) (StatusMeta, bool) {
	meta, ok := afterSalesStatusMeta[code]
	return meta, ok
}

// AfterSalesStatusCode 返回售后流程字符串状态对应的数字状态。
func AfterSalesStatusCode(status string) (int, bool) {
	return statusCodeByName(afterSalesStatusMeta, status)
}
