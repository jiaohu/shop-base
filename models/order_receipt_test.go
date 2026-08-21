package models

import (
	"encoding/json"
	"testing"
)

func TestOrderReceiptEventDataStringIncludesRecoveryIdentity(t *testing.T) {
	payload := OrderReceiptEventData{
		TxHash:              "0xabc",
		OrderId:             11,
		ExpectedSender:      "0x1",
		ExpirationTimestamp: 123,
	}.String()
	var decoded OrderReceiptEventData
	if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if decoded.TxHash != "0xabc" || decoded.OrderId != 11 || decoded.ExpectedSender != "0x1" || decoded.ExpirationTimestamp != 123 {
		t.Fatalf("unexpected payload: %+v", decoded)
	}
}

func TestOrderReceiptEventDataReadsLegacyAndIgnoresUnknownFields(t *testing.T) {
	var decoded OrderReceiptEventData
	if err := json.Unmarshal([]byte(`{"txHash":"0xabc","orderId":11,"futureField":true}`), &decoded); err != nil {
		t.Fatalf("decode legacy payload: %v", err)
	}
	if decoded.TxHash != "0xabc" || decoded.OrderId != 11 || decoded.ExpectedSender != "" || decoded.ExpirationTimestamp != 0 {
		t.Fatalf("unexpected legacy payload: %+v", decoded)
	}
}

func TestForceConfirmReceiptsEventDataStringIncludesExactIntent(t *testing.T) {
	payload := ForceConfirmReceiptsEventData{
		TxHash:              "0xdef",
		OrderIds:            []uint64{11, 12},
		ExpectedSender:      "0x2",
		ExpectedFunction:    "batch_force_confirm_receipts",
		ExpirationTimestamp: 456,
	}.String()
	var decoded ForceConfirmReceiptsEventData
	if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if decoded.TxHash != "0xdef" || len(decoded.OrderIds) != 2 || decoded.ExpectedSender != "0x2" || decoded.ExpectedFunction != "batch_force_confirm_receipts" || decoded.ExpirationTimestamp != 456 {
		t.Fatalf("unexpected payload: %+v", decoded)
	}
}

func TestSettlementRecoveryEventDataStringIncludesExactIntent(t *testing.T) {
	payload := SettlementRecoveryEventData{
		TxHash:              "0xabc",
		SettlementIds:       []uint64{21, 22},
		ExpectedSender:      "0x3",
		ExpirationTimestamp: 789,
	}.String()
	var decoded SettlementRecoveryEventData
	if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if decoded.TxHash != "0xabc" || len(decoded.SettlementIds) != 2 || decoded.ExpectedSender != "0x3" || decoded.ExpirationTimestamp != 789 {
		t.Fatalf("unexpected payload: %+v", decoded)
	}
}
