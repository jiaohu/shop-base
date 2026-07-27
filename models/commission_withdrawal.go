package models

import "encoding/json"

type CommissionWithdrawalPayoutRequestedEventData struct {
	WithdrawalNo string `json:"withdrawalNo"`
	Uid          uint64 `json:"uid"`
	PayoutMode   string `json:"payoutMode"`
	Channel      string `json:"channel"`
	Amount       string `json:"amount"`
	FeeAmount    string `json:"feeAmount"`
	PayAmount    string `json:"payAmount"`
	Currency     string `json:"currency"`
	PayeeMasked  string `json:"payeeMasked"`
}

func (e CommissionWithdrawalPayoutRequestedEventData) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return "{}"
	}
	return string(data)
}

type CommissionWithdrawalPayoutResultEventData struct {
	WithdrawalNo  string `json:"withdrawalNo"`
	Success       bool   `json:"success"`
	Provider      string `json:"provider"`
	ProviderReqNo string `json:"providerRequestNo"`
	ProviderTxnNo string `json:"providerTransactionNo"`
	FailureCode   string `json:"failureCode,omitempty"`
	Message       string `json:"message,omitempty"`
	CompletedAt   int64  `json:"completedAt"`
}

func (e CommissionWithdrawalPayoutResultEventData) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return "{}"
	}
	return string(data)
}
