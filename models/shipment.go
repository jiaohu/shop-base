package models

import (
	"fmt"
	"strconv"
	"strings"
)

// ShipmentDeliveredEventData is a logistics fact, independent of either mall.
// ShipmentId is an internal logistics ID, encoded as a decimal string; TrackingNo
// is the carrier's waybill number. DeliveredAt is UTC Unix seconds, never zero.
type ShipmentDeliveredEventData struct {
	EventId             string `json:"eventId"`
	Version             int    `json:"version"`
	ShipmentId          uint64 `json:"shipmentId,string"`
	Provider            string `json:"provider"`
	ProviderCarrierCode string `json:"providerCarrierCode"`
	TrackingNo          string `json:"trackingNo"`
	DeliveredAt         int64  `json:"deliveredAt"`
}

func ShipmentDeliveredEventKey(shipmentID uint64) string {
	return "shipment:delivered:v1:" + strconv.FormatUint(shipmentID, 10)
}

func (e ShipmentDeliveredEventData) Validate() error {
	if e.Version != 1 || e.ShipmentId == 0 || e.EventId != ShipmentDeliveredEventKey(e.ShipmentId) {
		return fmt.Errorf("invalid shipment delivered event identity/version")
	}
	if e.DeliveredAt <= 0 || e.DeliveredAt > 253402300799 {
		return fmt.Errorf("invalid deliveredAt (UTC Unix seconds required)")
	}
	for _, field := range []string{e.Provider, e.ProviderCarrierCode, e.TrackingNo} {
		if field == "" || field != strings.TrimSpace(field) || len(field) > 128 {
			return fmt.Errorf("invalid shipment waybill identity")
		}
	}
	return nil
}

type OrderShipmentRegisterEventData struct {
	OrderId uint64 `json:"orderId,string"`
}

type OrderDeliveredRecoveryEventData struct {
	ChainRecoveryIdentity
	OrderId uint64 `json:"orderId,string"`
}
