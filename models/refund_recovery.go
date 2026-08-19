package models

const (
	RefundRecoveryActionRequest            = "request_refund"
	RefundRecoveryActionMerchantReview     = "review_by_merchant"
	RefundRecoveryActionAdminReview        = "review_by_admin"
	RefundRecoveryActionPendingAdminReview = "review_pending_by_admin"
	RefundRecoveryActionCancel             = "cancel_refund_request"
)

// RefundRecoveryEventData is the durable payload for event 34. Optional JSON
// tags preserve compatibility with rows created before recovery identity was
// added.
type RefundRecoveryEventData struct {
	TxHash              string `json:"txHash"`
	OrderDetailId       uint64 `json:"orderDetailId"`
	RefundId            uint64 `json:"refundId,omitempty"`
	Action              string `json:"action,omitempty"`
	ExpectedSender      string `json:"expectedSender,omitempty"`
	ExpirationTimestamp uint64 `json:"expirationTimestamp,omitempty"`
}
