package orderstatus

// CanDeliverOrderDetailStatus 判断订单明细是否允许进入发货流程。
// 只有履约状态为待发货且售后锁为空时可以发货；售后中或已全量退款的明细不能发货。
func CanDeliverOrderDetailStatus(status, afterSalesStatus int) bool {
	return status == DetailStatusPendingShipment &&
		afterSalesStatus == DetailAfterSalesStatusNone
}

// DeliverableQuantity 返回订单明细当前可发货数量。
// 计算口径：购买数量 - 已发货数量 - 售后处理中数量 - 已退款数量。
// 返回值最低为 0，避免历史异常数据导致无符号整数下溢。
func DeliverableQuantity(quantity, delivered uint, afterSalesQuantity, refundedQuantity int) uint {
	available := int(quantity) - int(delivered) - afterSalesQuantity - refundedQuantity
	if available <= 0 {
		return 0
	}
	return uint(available)
}
