package models

import (
	"encoding/json"
	"testing"
)

func TestOrderPayEventDataStringIncludesRecoveryIdentity(t *testing.T) {
	payload := OrderPayEventData{
		TxHash:              "0xabc",
		OrderIds:            []uint64{11, 12},
		ExpectedSender:      "0x1",
		ExpirationTimestamp: 123,
	}.String()
	var decoded OrderPayEventData
	if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if decoded.TxHash != "0xabc" || len(decoded.OrderIds) != 2 || decoded.ExpectedSender != "0x1" || decoded.ExpirationTimestamp != 123 {
		t.Fatalf("unexpected payload: %+v", decoded)
	}
}
