package orderstatus

import "fmt"

type ReceiptDetail struct {
	ID, Quantity, Shipped, AfterSalesQuantity, RefundedQuantity, ExchangedQuantity uint64
	Status, AfterSalesStatus                                                       int
}

type ReceiptDelivery struct {
	ID, DetailID, ShipmentID, Quantity uint64
	DeliveredAt                        int64 // UTC Unix seconds; zero means unknown/not delivered.
}

// AllShipmentsDelivered requires full quantities and every parcel for every
// detail. Missing, virtual, refunded or exchanged details are not silently skipped.
func AllShipmentsDelivered(details []ReceiptDetail, deliveries []ReceiptDelivery) (int64, bool) {
	if len(details) == 0 || len(deliveries) == 0 {
		return 0, false
	}
	remaining := make(map[uint64]uint64, len(details))
	for _, d := range details {
		if d.ID == 0 || d.Quantity == 0 || d.Shipped != d.Quantity || d.Status != DetailStatusPendingReceipt ||
			d.AfterSalesStatus != DetailAfterSalesStatusNone || d.AfterSalesQuantity != 0 || d.RefundedQuantity != 0 || d.ExchangedQuantity != 0 {
			return 0, false
		}
		if _, duplicate := remaining[d.ID]; duplicate {
			return 0, false
		}
		remaining[d.ID] = d.Quantity
	}
	seen := make(map[uint64]bool, len(deliveries))
	var last int64
	for _, d := range deliveries {
		quantity, exists := remaining[d.DetailID]
		if !exists || d.ID == 0 || seen[d.ID] || d.ShipmentID == 0 || d.Quantity == 0 || d.Quantity > quantity || d.DeliveredAt <= 0 {
			return 0, false
		}
		seen[d.ID] = true
		remaining[d.DetailID] -= d.Quantity
		if d.DeliveredAt > last {
			last = d.DeliveredAt
		}
	}
	for _, quantity := range remaining {
		if quantity != 0 {
			return 0, false
		}
	}
	return last, true
}

// ReceiptDeadline adds whole 24-hour days. The caller handles Web2's disabled
// policy (zero days); Web3 permits zero days. Reject values outside MySQL dates.
func ReceiptDeadline(deliveredAt int64, days uint64) (int64, error) {
	const latest = int64(253402300799)
	if deliveredAt <= 0 || deliveredAt > latest || days > uint64((latest-deliveredAt)/86400) {
		return 0, fmt.Errorf("receipt deadline outside supported date range")
	}
	return deliveredAt + int64(days)*86400, nil
}
