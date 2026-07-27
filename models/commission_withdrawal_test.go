package models

import (
	"encoding/json"
	"testing"
)

func TestCommissionWithdrawalPayoutRequestedEventDataString(t *testing.T) {
	event := CommissionWithdrawalPayoutRequestedEventData{
		WithdrawalNo: "W202607240001",
		Uid:          12,
		PayoutMode:   "manual",
		Channel:      "alipay",
		Amount:       "100.00",
		FeeAmount:    "1.00",
		PayAmount:    "99.00",
		Currency:     "CNY",
		PayeeMasked:  "a***z@example.com",
	}
	var decoded CommissionWithdrawalPayoutRequestedEventData
	if err := json.Unmarshal([]byte(event.String()), &decoded); err != nil {
		t.Fatalf("unmarshal event: %v", err)
	}
	if decoded != event {
		t.Fatalf("decoded event mismatch: %#v", decoded)
	}
}

func TestCommissionWithdrawalPayoutResultEventDataString(t *testing.T) {
	event := CommissionWithdrawalPayoutResultEventData{
		WithdrawalNo:  "W202607240001",
		Success:       true,
		Provider:      "manual",
		ProviderReqNo: "W202607240001",
		ProviderTxnNo: "FINANCE-001",
		CompletedAt:   1784822400,
	}
	var decoded CommissionWithdrawalPayoutResultEventData
	if err := json.Unmarshal([]byte(event.String()), &decoded); err != nil {
		t.Fatalf("unmarshal event: %v", err)
	}
	if decoded != event {
		t.Fatalf("decoded event mismatch: %#v", decoded)
	}
}
