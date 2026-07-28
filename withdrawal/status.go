package withdrawal

const (
	StatusPendingReview = 10
	StatusPayoutPending = 20
	StatusProcessing    = 30
	StatusQuerying      = 40
	StatusSucceeded     = 50
	StatusRejected      = 60
	StatusCanceled      = 70
	StatusFailed        = 80
)

const (
	PayoutResultNone    = 0
	PayoutResultSuccess = 1
	PayoutResultFailed  = 2
)

const (
	PayoutModeManual = "manual"
	PayoutModeAlipay = "alipay"
)

func IsFinal(status int) bool {
	switch status {
	case StatusSucceeded, StatusRejected, StatusCanceled, StatusFailed:
		return true
	default:
		return false
	}
}

func CanCancel(status int) bool {
	return status == StatusPendingReview
}

func CanApprove(status int) bool {
	return status == StatusPendingReview
}

func CanSubmitPayoutResult(status int) bool {
	return status == StatusPayoutPending
}

func CanApplyPayoutResult(status int) bool {
	return status == StatusPayoutPending ||
		status == StatusProcessing ||
		status == StatusQuerying
}
