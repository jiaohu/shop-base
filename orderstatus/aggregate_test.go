package orderstatus

import "testing"

func TestCalcOrderStatusMixedFulfillment(t *testing.T) {
	for _, status := range []int{DetailStatusPendingReceipt, DetailStatusReceived, DetailStatusPendingComment, DetailStatusCompleted} {
		details := []DetailSnapshot{
			{Status: DetailStatusPendingShipment, Quantity: 2},
			{Status: status, Quantity: 1, Delivered: 1},
		}
		if got := CalcOrderStatus(details); got != OrderStatusPartShipped {
			t.Errorf("unshipped + detail status %d: got %d, want part shipped", status, got)
		}
		details[0], details[1] = details[1], details[0]
		if got := CalcOrderStatus(details); got != OrderStatusPartShipped {
			t.Errorf("reversed details with status %d: got %d, want part shipped", status, got)
		}
	}
}

func TestCalcOrderStatusFulfillmentBoundaries(t *testing.T) {
	for _, tc := range []struct {
		name    string
		details []DetailSnapshot
		want    int
	}{
		{"empty", nil, OrderStatusPendingPayment},
		{"all canceled", []DetailSnapshot{{Status: DetailStatusCanceled}}, OrderStatusClosed},
		{"unshipped and canceled", []DetailSnapshot{{Status: DetailStatusPendingShipment}, {Status: DetailStatusCanceled, Delivered: 1}}, OrderStatusPendingShipment},
		{"unshipped and closed", []DetailSnapshot{{Status: DetailStatusPendingShipment}, {Status: DetailStatusClosed}}, OrderStatusPendingShipment},
		{"all shipped", []DetailSnapshot{{Status: DetailStatusPendingReceipt}, {Status: DetailStatusPendingComment}}, OrderStatusPendingReceipt},
		{"all received", []DetailSnapshot{{Status: DetailStatusPendingComment}, {Status: DetailStatusCompleted}}, OrderStatusPendingComment},
		{"all completed", []DetailSnapshot{{Status: DetailStatusCompleted}, {Status: DetailStatusCompleted}}, OrderStatusCompleted},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := CalcOrderStatus(tc.details); got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}

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
