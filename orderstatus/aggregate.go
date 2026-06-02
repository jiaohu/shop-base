package orderstatus

// DetailSnapshot 是订单状态聚合需要读取的明细最小字段集。
// 业务服务从 store_order_detail 取数后转换成该结构，避免 shop-base 依赖具体 DAO。
type DetailSnapshot struct {
	Status             int // 明细履约状态
	Quantity           int // 购买数量
	Delivered          int // 已发货数量
	RefundedQuantity   int // 已退款数量
	ExchangedQuantity  int // 已换货完成数量
	AfterSalesQuantity int // 当前售后处理中数量
}

// AfterSalesSnapshot 是订单售后聚合需要读取的售后单最小字段集。
// 当前聚合主要使用 Status，数量字段保留给部分售后的后续扩展。
type AfterSalesSnapshot struct {
	Status           int    // 售后流程状态
	RequestType      string // 售后类型
	ApprovedQuantity int    // 商家同意数量
	ApplyQuantity    int    // 申请数量
}

// CalcDetailAfterSalesStatus 根据数量字段计算明细售后锁状态。
// 优先级：处理中数量 > 全量已退款 > 无售后。
func CalcDetailAfterSalesStatus(quantity, afterSalesQuantity, refundedQuantity int) int {
	if afterSalesQuantity > 0 {
		return DetailAfterSalesStatusProcessing
	}
	if quantity > 0 && refundedQuantity >= quantity {
		return DetailAfterSalesStatusRefunded
	}
	return DetailAfterSalesStatusNone
}

// CalcOrderStatus 根据明细履约状态聚合主订单履约状态。
// 该函数不处理售后语义；售后聚合由 CalcOrderAfterSalesStatus 负责。
func CalcOrderStatus(details []DetailSnapshot) int {
	if len(details) == 0 {
		return OrderStatusPendingPayment
	}

	active := make([]DetailSnapshot, 0, len(details))
	for _, detail := range details {
		if detail.Status == DetailStatusCanceled {
			continue
		}
		active = append(active, detail)
	}
	if len(active) == 0 {
		return OrderStatusClosed
	}

	pendingShipment := 0
	pendingReceipt := 0
	received := 0
	pendingComment := 0
	completed := 0
	closed := 0
	pendingPayment := 0
	partiallyDeliveredPendingShipment := false

	for _, detail := range active {
		switch detail.Status {
		case DetailStatusPendingPayment:
			pendingPayment++
		case DetailStatusPendingShipment:
			pendingShipment++
			if detail.Delivered > 0 {
				partiallyDeliveredPendingShipment = true
			}
		case DetailStatusPendingReceipt:
			pendingReceipt++
		case DetailStatusReceived:
			received++
		case DetailStatusPendingComment:
			pendingComment++
		case DetailStatusCompleted:
			completed++
		case DetailStatusClosed:
			closed++
		}
	}

	total := len(active)
	if pendingPayment == total {
		return OrderStatusPendingPayment
	}
	// 待发货和待收货同时存在，或单个待发货明细已经有部分发货数量时，说明订单已经部分发货。
	if pendingShipment > 0 && (pendingReceipt > 0 || partiallyDeliveredPendingShipment) {
		return OrderStatusPartShipped
	}
	if pendingShipment > 0 {
		return OrderStatusPendingShipment
	}
	if pendingReceipt > 0 {
		return OrderStatusPendingReceipt
	}
	if pendingComment > 0 {
		return OrderStatusPendingComment
	}
	if received > 0 {
		return OrderStatusReceived
	}
	if completed+closed == total {
		if completed > 0 {
			return OrderStatusCompleted
		}
		return OrderStatusClosed
	}
	return OrderStatusCompleted
}

// CalcOrderAfterSalesStatus 根据售后单流程和明细数量聚合主订单售后状态。
// 聚合优先级固定为：
// processing > all_refunded > part_refunded > all_done > part_exchanged > none。
func CalcOrderAfterSalesStatus(details []DetailSnapshot, afterSales []AfterSalesSnapshot) int {
	// 任意进行中的售后单都优先显示为售后处理中。
	for _, item := range afterSales {
		if IsAfterSalesProcessing(item.Status) {
			return OrderAfterSalesStatusProcessing
		}
	}

	totalQuantity := 0
	refundedTotal := 0
	exchangedTotal := 0
	successAfterSalesTotal := 0

	for _, detail := range details {
		if detail.Status == DetailStatusCanceled || detail.Quantity <= 0 {
			continue
		}
		totalQuantity += detail.Quantity
		// 使用 min 避免同一件商品连续换货、退款导致累计数量超过购买数量。
		refunded := minInt(detail.RefundedQuantity, detail.Quantity)
		exchanged := minInt(detail.ExchangedQuantity, detail.Quantity)
		success := minInt(detail.RefundedQuantity+detail.ExchangedQuantity, detail.Quantity)
		refundedTotal += refunded
		exchangedTotal += exchanged
		successAfterSalesTotal += success
	}

	if totalQuantity == 0 {
		return OrderAfterSalesStatusNone
	}
	if refundedTotal >= totalQuantity {
		return OrderAfterSalesStatusAllRefunded
	}
	if refundedTotal > 0 {
		return OrderAfterSalesStatusPartRefunded
	}
	if successAfterSalesTotal >= totalQuantity {
		return OrderAfterSalesStatusAllDone
	}
	if exchangedTotal > 0 {
		return OrderAfterSalesStatusPartExchanged
	}
	return OrderAfterSalesStatusNone
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
