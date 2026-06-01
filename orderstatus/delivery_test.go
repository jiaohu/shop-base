package orderstatus

import "testing"

func TestCanDeliverOrderDetailStatus(t *testing.T) {
	tests := []struct {
		name             string
		status           int
		afterSalesStatus int
		want             bool
	}{
		{name: "pending payment", status: DetailStatusPendingPayment, afterSalesStatus: DetailAfterSalesStatusNone, want: false},
		{name: "pending shipment", status: DetailStatusPendingShipment, afterSalesStatus: DetailAfterSalesStatusNone, want: true},
		{name: "after sales processing", status: DetailStatusPendingShipment, afterSalesStatus: DetailAfterSalesStatusProcessing, want: false},
		{name: "refunded", status: DetailStatusPendingShipment, afterSalesStatus: DetailAfterSalesStatusRefunded, want: false},
		{name: "pending receipt", status: DetailStatusPendingReceipt, afterSalesStatus: DetailAfterSalesStatusNone, want: false},
		{name: "closed", status: DetailStatusClosed, afterSalesStatus: DetailAfterSalesStatusNone, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanDeliverOrderDetailStatus(tt.status, tt.afterSalesStatus); got != tt.want {
				t.Fatalf("CanDeliverOrderDetailStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeliverableQuantity(t *testing.T) {
	tests := []struct {
		name               string
		quantity           uint
		delivered          uint
		afterSalesQuantity int
		refundedQuantity   int
		want               uint
	}{
		{name: "deduct delivered after sales and refunded", quantity: 10, delivered: 2, afterSalesQuantity: 1, refundedQuantity: 1, want: 6},
		{name: "after sales consumes remaining quantity", quantity: 10, delivered: 8, afterSalesQuantity: 2, refundedQuantity: 0, want: 0},
		{name: "historical over delivered data clamps to zero", quantity: 3, delivered: 5, afterSalesQuantity: 0, refundedQuantity: 0, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeliverableQuantity(tt.quantity, tt.delivered, tt.afterSalesQuantity, tt.refundedQuantity)
			if got != tt.want {
				t.Fatalf("DeliverableQuantity() = %d, want %d", got, tt.want)
			}
		})
	}
}
