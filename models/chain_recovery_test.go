package models

import (
	"encoding/json"
	"testing"
)

func TestDomainRecoveryEventsIncludeTransactionIdentity(t *testing.T) {
	identity := ChainRecoveryIdentity{TxHash: "0xabc", ExpectedSender: "0x1", ExpirationTimestamp: 99}
	cases := []struct {
		name    string
		payload string
	}{
		{"token", (TokenRegistryRecoveryEventData{ChainRecoveryIdentity: identity, Operation: TokenOperationRegister, TokenAddress: "0x2", Symbol: "USDC", Decimal: 6}).String()},
		{"permission", (PermissionRecoveryEventData{ChainRecoveryIdentity: identity, Operation: PermissionOperationAddAdmin, TargetAddresses: []string{"0x2"}, Permissions: []int{1, 2}}).String()},
		{"config", (Web3ConfigRecoveryEventData{ChainRecoveryIdentity: identity, Operation: ConfigOperationUpdateSettlementPeriod, Values: map[string]string{"settlement_period": "60"}}).String()},
		{"product", (ProductRecoveryEventData{ChainRecoveryIdentity: identity, Operation: ProductOperationCreate, ProductId: 10}).String()},
		{"listing", (ListingRecoveryEventData{ChainRecoveryIdentity: identity, Operation: ListingOperationSubmit, ProductId: 10, TicketId: 11, RequiredLevel: 2}).String()},
		{"flash sale", (FlashSaleRecoveryEventData{ChainRecoveryIdentity: identity, Operation: FlashSaleOperationCreate, SaleId: 12}).String()},
		{"shipment", (OrderShipmentRecoveryEventData{ChainRecoveryIdentity: identity, OrderId: 13, AdminId: 14}).String()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var decoded ChainRecoveryIdentity
			if err := json.Unmarshal([]byte(tc.payload), &decoded); err != nil {
				t.Fatalf("decode payload: %v", err)
			}
			if decoded != identity {
				t.Fatalf("identity mismatch: got %+v want %+v", decoded, identity)
			}
		})
	}
}

func TestReferralAccountSyncEventIncludesExpiration(t *testing.T) {
	payload := (ReferralAccountSyncEvent{
		TxHash: "0xabc", Operation: ReferralOperationRegister,
		Sender: "0x1", ExpirationTimestamp: 123,
	}).String()
	var decoded ReferralAccountSyncEvent
	if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if decoded.ExpirationTimestamp != 123 {
		t.Fatalf("expiration = %d, want 123", decoded.ExpirationTimestamp)
	}
}
