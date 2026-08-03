package models

import (
	"encoding/json"
	"testing"
)

func TestReferralAccountSyncEventString(t *testing.T) {
	t.Parallel()

	want := ReferralAccountSyncEvent{
		TxHash:       "0xabc",
		Operation:    ReferralOperationAcceptInvitationByCode,
		Sender:       "0x1",
		ReferralCode: "ABCDEFG1",
	}

	var got ReferralAccountSyncEvent
	if err := json.Unmarshal([]byte(want.String()), &got); err != nil {
		t.Fatalf("unmarshal ReferralAccountSyncEvent: %v", err)
	}
	if got != want {
		t.Fatalf("ReferralAccountSyncEvent mismatch: got %+v, want %+v", got, want)
	}
}

func TestReferralAccountSyncEventStringOmitsEmptyCode(t *testing.T) {
	t.Parallel()

	event := ReferralAccountSyncEvent{
		TxHash:    "0xdef",
		Operation: ReferralOperationRegister,
		Sender:    "0x2",
	}
	if got := event.String(); got != `{"txHash":"0xdef","operation":"register","sender":"0x2"}` {
		t.Fatalf("unexpected JSON: %s", got)
	}
}
