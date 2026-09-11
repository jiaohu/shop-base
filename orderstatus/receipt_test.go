package orderstatus

import "testing"

func TestAllShipmentsDelivered(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func([]ReceiptDetail, []ReceiptDelivery)
		ready  bool
	}{
		{"split and combined parcels", func(_ []ReceiptDetail, _ []ReceiptDelivery) {}, true},
		{"one parcel in transit", func(_ []ReceiptDetail, p []ReceiptDelivery) { p[1].DeliveredAt = 0 }, false},
		{"missing binding", func(_ []ReceiptDetail, p []ReceiptDelivery) { p[0].ShipmentID = 0 }, false},
		{"quantity not covered", func(_ []ReceiptDetail, p []ReceiptDelivery) { p[0].Quantity = 0 }, false},
		{"excess quantity", func(_ []ReceiptDetail, p []ReceiptDelivery) { p[0].Quantity = 2 }, false},
		{"duplicate delivery", func(_ []ReceiptDetail, p []ReceiptDelivery) { p[1].ID = p[0].ID }, false},
		{"unshipped detail", func(d []ReceiptDetail, _ []ReceiptDelivery) { d[1].Shipped = 0 }, false},
		{"after sales", func(d []ReceiptDetail, _ []ReceiptDelivery) { d[0].AfterSalesQuantity = 1 }, false},
		{"refund", func(d []ReceiptDetail, _ []ReceiptDelivery) { d[0].RefundedQuantity = 1 }, false},
		{"exchange", func(d []ReceiptDetail, _ []ReceiptDelivery) { d[0].ExchangedQuantity = 1 }, false},
		{"already received", func(d []ReceiptDetail, _ []ReceiptDelivery) { d[0].Status = DetailStatusReceived }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := []ReceiptDetail{{ID: 1, Quantity: 2, Shipped: 2, Status: DetailStatusPendingReceipt}, {ID: 2, Quantity: 1, Shipped: 1, Status: DetailStatusPendingReceipt}}
			p := []ReceiptDelivery{{ID: 1, DetailID: 1, ShipmentID: 10, Quantity: 1, DeliveredAt: 100}, {ID: 2, DetailID: 1, ShipmentID: 20, Quantity: 1, DeliveredAt: 200}, {ID: 3, DetailID: 2, ShipmentID: 20, Quantity: 1, DeliveredAt: 200}}
			tc.mutate(d, p)
			last, ready := AllShipmentsDelivered(d, p)
			if ready != tc.ready || ready && last != 200 {
				t.Fatalf("got (%d,%v)", last, ready)
			}
		})
	}
	if _, ready := AllShipmentsDelivered(nil, nil); ready {
		t.Fatal("empty order is not delivered")
	}
}

func TestReceiptDeadline(t *testing.T) {
	if got, err := ReceiptDeadline(100, 7); err != nil || got != 604900 {
		t.Fatalf("%d %v", got, err)
	}
	if got, err := ReceiptDeadline(100, 0); err != nil || got != 100 {
		t.Fatalf("%d %v", got, err)
	}
	if _, err := ReceiptDeadline(100, ^uint64(0)); err == nil {
		t.Fatal("overflow accepted")
	}
}
