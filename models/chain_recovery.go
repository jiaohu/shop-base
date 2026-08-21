package models

import "encoding/json"

const (
	TokenOperationRegister   = "register_token"
	TokenOperationUnregister = "unregister_token"

	PermissionOperationAddAdmin            = "add_admin"
	PermissionOperationUpdateAdmin         = "update_admin"
	PermissionOperationRemoveAdmin         = "remove_admin"
	PermissionOperationTransferSuperAdmin  = "transfer_super_admin"
	PermissionOperationAddReviewer         = "add_listing_reviewer"
	PermissionOperationUpdateReviewerLevel = "update_listing_reviewer_level"
	PermissionOperationRemoveReviewer      = "remove_listing_reviewer"
	PermissionOperationRemoveReviewerLevel = "remove_listing_reviewer_level"
	PermissionOperationAddMerchants        = "add_merchants"
	PermissionOperationRemoveMerchants     = "remove_merchants"

	ConfigOperationUpdateSettlementPeriod = "update_lock_period"
	ConfigOperationUpdateReceipt          = "update_receipt_confirmation_config"
	ConfigOperationUpdateOrder            = "update_order_config"

	ProductOperationCreate     = "create_product"
	ProductOperationUpdate     = "update_product"
	ProductOperationActivate   = "activate_product"
	ProductOperationDeactivate = "deactivate_product"

	ListingOperationSubmit = "submit_for_listing_approval"
	ListingOperationReview = "review_listing"

	FlashSaleOperationCreate = "create_flash_sale"
	FlashSaleOperationUpdate = "update_flash_sale"
	FlashSaleOperationCancel = "cancel_flash_sale"
	FlashSaleOperationEnd    = "end_flash_sale"
)

// ChainRecoveryIdentity identifies the exact signed transaction whose local
// projection must be recovered. It is embedded by every domain recovery event.
type ChainRecoveryIdentity struct {
	TxHash              string `json:"txHash"`
	ExpectedSender      string `json:"expectedSender"`
	ExpirationTimestamp uint64 `json:"expirationTimestamp"`
}

type TokenRegistryRecoveryEventData struct {
	ChainRecoveryIdentity
	Operation    string `json:"operation"`
	TokenAddress string `json:"tokenAddress"`
	Symbol       string `json:"symbol,omitempty"`
	Name         string `json:"name,omitempty"`
	Decimal      int    `json:"decimal"`
}

func (e TokenRegistryRecoveryEventData) String() string { return marshalRecoveryEvent(e) }

type PermissionRecoveryEventData struct {
	ChainRecoveryIdentity
	Operation       string   `json:"operation"`
	TargetAddresses []string `json:"targetAddresses"`
	Permissions     []int    `json:"permissions"`
	Level           int      `json:"level,omitempty"`
}

func (e PermissionRecoveryEventData) String() string { return marshalRecoveryEvent(e) }

type Web3ConfigRecoveryEventData struct {
	ChainRecoveryIdentity
	Operation string            `json:"operation"`
	Values    map[string]string `json:"values"`
}

func (e Web3ConfigRecoveryEventData) String() string { return marshalRecoveryEvent(e) }

type ProductRecoveryEventData struct {
	ChainRecoveryIdentity
	Operation string `json:"operation"`
	ProductId uint64 `json:"productId"`
}

func (e ProductRecoveryEventData) String() string { return marshalRecoveryEvent(e) }

type ListingRecoveryEventData struct {
	ChainRecoveryIdentity
	Operation     string `json:"operation"`
	ProductId     uint64 `json:"productId"`
	TicketId      uint64 `json:"ticketId"`
	RequiredLevel int    `json:"requiredLevel"`
	Approved      bool   `json:"approved"`
}

func (e ListingRecoveryEventData) String() string { return marshalRecoveryEvent(e) }

type FlashSaleRecoveryEventData struct {
	ChainRecoveryIdentity
	Operation string `json:"operation"`
	SaleId    uint64 `json:"saleId"`
}

func (e FlashSaleRecoveryEventData) String() string { return marshalRecoveryEvent(e) }

type OrderShipmentRecoveryEventData struct {
	ChainRecoveryIdentity
	OrderId        uint64                      `json:"orderId"`
	AdminId        uint64                      `json:"adminId"`
	DeliveryType   int                         `json:"deliveryType"`
	ExpressCompany uint64                      `json:"expressCompany,omitempty"`
	ExpressNo      string                      `json:"expressNo,omitempty"`
	Deliveryman    uint64                      `json:"deliveryman,omitempty"`
	Remark         string                      `json:"remark,omitempty"`
	IsSplitDeliver bool                        `json:"isSplitDeliver"`
	DeliveryList   []OrderShipmentDeliveryItem `json:"deliveryList,omitempty"`
}

type OrderShipmentDeliveryItem struct {
	DeliveryId        uint64 `json:"deliveryId"`
	DetailId          uint64 `json:"detailId"`
	CurrentDeliverNum uint   `json:"currentDeliverNum"`
}

func (e OrderShipmentRecoveryEventData) String() string { return marshalRecoveryEvent(e) }

func marshalRecoveryEvent(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(data)
}
