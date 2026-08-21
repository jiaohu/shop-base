package models

import "encoding/json"

type OrderPayEventData struct {
	TxHash              string   `json:"txHash"`
	OrderIds            []uint64 `json:"orderIds,omitempty"`
	ExpectedSender      string   `json:"expectedSender,omitempty"`
	ExpirationTimestamp uint64   `json:"expirationTimestamp,omitempty"`
}

func (o OrderPayEventData) String() string {
	res, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(res)
}

type OrderDetailRefundEventData struct {
	OrderDetailId uint64 `json:"orderDetailId"`
	TxHash        string `json:"txHash"`
}

func (o OrderDetailRefundEventData) String() string {
	res, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(res)
}

type RefundStockSyncEventData struct {
	AfterSalesId uint64 `json:"afterSalesId"`
	SkuId        uint64 `json:"skuId"`
	Quantity     int    `json:"quantity"`
}

type OrderReceiptEventData struct {
	TxHash              string `json:"txHash"`
	OrderId             uint64 `json:"orderId"`
	ExpectedSender      string `json:"expectedSender,omitempty"`
	ExpirationTimestamp uint64 `json:"expirationTimestamp,omitempty"`
}

func (o OrderReceiptEventData) String() string {
	res, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(res)
}

// ForceConfirmReceiptsEventData records the exact platform force-confirm
// transaction intent. ExpectedFunction is either force_confirm_receipt or
// batch_force_confirm_receipts; OrderIds must match the chain events exactly.
type ForceConfirmReceiptsEventData struct {
	TxHash              string   `json:"txHash"`
	OrderIds            []uint64 `json:"orderIds"`
	ExpectedSender      string   `json:"expectedSender"`
	ExpectedFunction    string   `json:"expectedFunction"`
	ExpirationTimestamp uint64   `json:"expirationTimestamp"`
}

func (o ForceConfirmReceiptsEventData) String() string {
	res, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(res)
}

type SettlementCompletedEventData struct {
	TxHash        string   `json:"txHash"`
	SettlementIds []uint64 `json:"settlementIds"`
	CallerAddr    string   `json:"callerAddr"`
	SettledAt     int64    `json:"settledAt"`
}

// SettlementRecoveryEventData records the exact execute_settlements intent.
// The consumer verifies the transaction payload and SettlementExecutedEvent
// records before projecting the transaction into local settlement state.
type SettlementRecoveryEventData struct {
	TxHash              string   `json:"txHash"`
	SettlementIds       []uint64 `json:"settlementIds"`
	ExpectedSender      string   `json:"expectedSender"`
	ExpirationTimestamp uint64   `json:"expirationTimestamp"`
}

func (s SettlementRecoveryEventData) String() string {
	res, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(res)
}

func (s SettlementCompletedEventData) String() string {
	res, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(res)
}
