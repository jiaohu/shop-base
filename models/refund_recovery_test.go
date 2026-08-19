package models

import (
	"encoding/json"
	"testing"
)

func TestRefundRecoveryEventDataJSONContract(t *testing.T) {
	payload, err := json.Marshal(RefundRecoveryEventData{
		TxHash: "0xabc", OrderDetailId: 11, RefundId: 12,
		Action: RefundRecoveryActionRequest, ExpectedSender: "0x1", ExpirationTimestamp: 123,
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	var decoded RefundRecoveryEventData
	if err = json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if decoded.TxHash != "0xabc" || decoded.OrderDetailId != 11 || decoded.RefundId != 12 || decoded.Action != RefundRecoveryActionRequest {
		t.Fatalf("unexpected payload: %+v", decoded)
	}
}
