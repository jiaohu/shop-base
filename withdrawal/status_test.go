package withdrawal

import "testing"

func TestStatusRules(t *testing.T) {
	if !CanCancel(StatusPendingReview) || CanCancel(StatusPayoutPending) {
		t.Fatal("cancel rule mismatch")
	}
	if !CanApprove(StatusPendingReview) || CanApprove(StatusSucceeded) {
		t.Fatal("approve rule mismatch")
	}
	if !CanApplyPayoutResult(StatusProcessing) || CanApplyPayoutResult(StatusRejected) {
		t.Fatal("payout result rule mismatch")
	}
	if !IsFinal(StatusSucceeded) || IsFinal(StatusQuerying) {
		t.Fatal("final rule mismatch")
	}
}
