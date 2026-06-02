package orderstatus

import "testing"

func TestCalcOrderStatusPartShippedByDeliveredQuantity(t *testing.T) {
	got := CalcOrderStatus([]DetailSnapshot{
		{
			Status:    DetailStatusPendingShipment,
			Quantity:  10,
			Delivered: 3,
		},
	})
	if got != OrderStatusPartShipped {
		t.Fatalf("CalcOrderStatus() = %d, want %d", got, OrderStatusPartShipped)
	}
}

func TestCalcOrderStatusPendingShipmentWithoutDeliveredQuantity(t *testing.T) {
	got := CalcOrderStatus([]DetailSnapshot{
		{
			Status:    DetailStatusPendingShipment,
			Quantity:  10,
			Delivered: 0,
		},
	})
	if got != OrderStatusPendingShipment {
		t.Fatalf("CalcOrderStatus() = %d, want %d", got, OrderStatusPendingShipment)
	}
}
