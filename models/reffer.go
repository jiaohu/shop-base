package models

import "encoding/json"

const (
	ReferralOperationRegister               = "register"
	ReferralOperationAcceptInvitationByCode = "accept_invitation_by_code"
)

type RegisterByRefferCode struct {
	TxHash              string `json:"txHash"`
	RefferedByCode      string `json:"refferedByCode"`
	ExpectedSender      string `json:"expectedSender,omitempty"`
	ExpirationTimestamp uint64 `json:"expirationTimestamp,omitempty"`
}

func (r RegisterByRefferCode) String() string {
	res, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(res)
}

// ReferralAccountSyncEvent requests asynchronous projection of a confirmed
// referral entry transaction. ReferralCode is populated only for the
// accept_invitation_by_code entry and is copied from the signed transaction.
type ReferralAccountSyncEvent struct {
	TxHash              string `json:"txHash"`
	Operation           string `json:"operation"`
	Sender              string `json:"sender"`
	ReferralCode        string `json:"referralCode,omitempty"`
	ExpirationTimestamp uint64 `json:"expirationTimestamp"`
}

func (r ReferralAccountSyncEvent) String() string {
	res, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(res)
}
