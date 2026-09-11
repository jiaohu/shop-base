package models

import (
	"encoding/json"
	"testing"
)

func TestShipmentDeliveredContract(t *testing.T) {
	e := ShipmentDeliveredEventData{EventId: ShipmentDeliveredEventKey(9007199254740993), Version: 1, ShipmentId: 9007199254740993, Provider: "kuaidi100", ProviderCarrierCode: "shunfeng", TrackingNo: "SF123", DeliveredAt: 1700000000}
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	var got ShipmentDeliveredEventData
	if err = json.Unmarshal(b, &got); err != nil || got != e {
		t.Fatalf("roundtrip: %s %v", b, err)
	}
	if err = got.Validate(); err != nil {
		t.Fatal(err)
	}
	got.Version = 0
	if got.Validate() == nil {
		t.Fatal("missing version accepted")
	}
	got = e
	got.EventId = "other"
	if got.Validate() == nil {
		t.Fatal("unstable identity accepted")
	}
}
