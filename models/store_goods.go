package models

import (
	"encoding/json"
	"time"
)

// 商品事件，传id就够了
type StoreGoodsEventData struct {
	Id int64 `json:"id"`
}

func (s StoreGoodsEventData) String() string {
	res, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(res)
}

type StoreGoodsAddReq struct {
	Name                string                       `json:"name" v:"required#商品名称不能为空" dc:"商品名称"`
	Summary             string                       `json:"summary"  dc:"商品描述"`
	Imgs                StringArray                  `json:"imgs"  dc:"封面图"`
	Thumb               string                       `json:"thumb"  dc:"缩略图"`
	Keywords            string                       `json:"keywords"  dc:"关键字"`
	CateId              Uint64Array                  `json:"cateId"  dc:"分类"`
	UnitName            string                       `json:"unitName" v:"required#单位不能为空" dc:"单位"`
	Sort                int                          `json:"sort"  dc:"排序"`
	IsShow              int                          `json:"isShow"  dc:"是否显示"`
	IsHot               int                          `json:"isHot"  dc:"是否热卖"`
	IsBest              int                          `json:"isBest"  dc:"是否精品"`
	IsNew               int                          `json:"isNew"  dc:"是否新品"`
	SpecificationType   int                          `json:"specificationType" dc:"规格"`
	GoodsType           int                          `json:"goodsType"  dc:"虚拟商品类型"`
	IsPostage           int                          `json:"isPostage"  dc:"是否包邮"`
	GiftPoints          uint                         `json:"giftPoints"  dc:"赠送积分"`
	Hits                int64                        `json:"hits"  dc:"点击量"`
	VideoUrl            string                       `json:"videoUrl"  dc:"视频地址"`
	FreightTempId       int                          `json:"freightTempId"  dc:"运费模板"`
	Spu                 string                       `json:"spu"  dc:"商品SPU"`
	LabelId             Uint64Array                  `json:"labelId"  dc:"商品标签ID"`
	IsVip               int                          `json:"isVip"  dc:"是否VIP商品"`
	Freight             float64                      `json:"freight"  dc:"运费"`
	Logistics           Uint64Array                  `json:"logistics"  dc:"物流方式"`
	MinOrderNum         uint                         `json:"minOrderNum"  dc:"起购数量"`
	MaxOrderNum         uint                         `json:"maxOrderNum"  dc:"限购数量"`
	PurchaseLimited     uint                         `json:"purchaseLimited"  dc:"是否限购"`
	PurchaseLimitedType uint                         `json:"purchaseLimitedType"  dc:"限购类型，0:单次限购，1:单人限购"`
	Presell             uint                         `json:"presell"  dc:"是否预售"`
	PresellBeginTime    *time.Time                   `json:"presellBeginTime"  dc:"预售开始时间"`
	PresellEndTime      *time.Time                   `json:"presellEndTime"  dc:"预售结束时间"`
	PresellDeliveryTime *time.Time                   `json:"presellDeliveryTime"  dc:"预售结束发货时间"`
	Attr                string                       `json:"attr" dc:"商品属性json串"`
	AttrValues          []*StoreGoodsAttrValueAddReq `json:"attrValues" dc:"多规格商品属性值"`
	AttrValue           *StoreGoodsAttrValueAddReq   `json:"attrValue" dc:"单规格商品属性值"`
	IsoCode             string                       `json:"isoCode" dc:"商品所属国家"`
	SignedTx            string                       `json:"signedTx" dc:"签名后的交易体, bcs hex string"`
}

type StoreGoodsAttrValueAddReq struct {
	GoodsId           uint64                     `json:"goodsId"  dc:"商品ID"`
	Suk               string                     `json:"suk"  dc:"属性值索引"`
	Stock             int                        `json:"stock"  dc:"库存"`
	Sales             int                        `json:"sales"  dc:"销量"`
	Price             float64                    `json:"price"  dc:"价格"`
	Image             string                     `json:"image"  dc:"图片"`
	Uuid              string                     `json:"uuid"  dc:"编号"`
	Cost              float64                    `json:"cost"  dc:"成本"`
	BarCode           string                     `json:"barCode"  dc:"条码"`
	OriginalPrice     float64                    `json:"originalPrice"  dc:"原价"`
	VipPrice          float64                    `json:"vipPrice"  dc:"会员价"`
	Weight            float64                    `json:"weight"  dc:"重量(kg)"`
	Volume            float64                    `json:"volume"  dc:"体积(m3)"`
	IsVirtual         int                        // 是否是虚拟商品
	VirtualType       int                        // 虚拟商品类型（1、固定卡密；2、一次性卡密）
	VirtualInfo       string                     // 虚拟商品内容
	CouponId          uint64                     // 关联优惠券
	CardPasswordGoods *StoreGoodsPasswordCardReq `json:"cardPasswordGoods" dc:"卡密商品"`
	Type              int                        `json:"type" dc:"活动类型"` // 0=商品，1=秒杀，2=砍价，3=拼团
	Points            float64                    `json:"points" dc:"购买奖励积分"`
}

// 卡密模型
type StoreGoodsPasswordCardReq struct {
	VirtualType   int                         `json:"virtualType" dc:"虚拟商品类型"`
	VirtualInfo   string                      `json:"virtualInfo" dc:"虚拟商品信息"`
	FixedNumber   string                      `json:"fixedNumber" dc:"固定卡密数量"`
	SingleUseList []*SingleUsePasswordCardReq `json:"singleUseList" dc:"一次性卡密列表"`
}
type SingleUsePasswordCardReq struct {
	Id               uint64 `json:"id" dc:"主键id"`                                                         //主键id
	GoodsId          uint64 `json:"goodsId" v:"goodsId@integer#商品id需为整数" dc:"商品id"`                       //商品id
	GoodsAttrValueId uint64 `json:"goodsAttrValueId" v:"goodsAttrValueId@integer#商品规格id需为整数" dc:"商品规格id"` //商品规格id
	CardNumber       string `json:"cardNumber" dc:"卡密卡号"`                                                 //卡密卡号
	CardPassword     string `json:"cardPassword" dc:"卡密密码"`                                               //卡密密码
	SnowId           string `json:"snowId" dc:"卡密雪花id"`                                                   //卡密雪花id
}

type Web3StoreGoodsAddEventData struct {
	TxHash string           `json:"txHash"`
	Param  StoreGoodsAddReq `json:"param"`
}

func (w Web3StoreGoodsAddEventData) String() string {
	res, err := json.Marshal(w)
	if err != nil {
		return ""
	}

	return string(res)
}

type StoreGoodsAttrValueEditReq struct {
	Id                uint64                     `json:"id"  dc:"商品属性值ID"`
	GoodsId           uint64                     `json:"goodsId"  dc:"商品ID"`
	Suk               string                     `json:"suk"  dc:"属性值索引"`
	Stock             int                        `json:"stock"  dc:"库存"`
	Sales             int                        `json:"sales"  dc:"销量"`
	Price             float64                    `json:"price"  dc:"价格"`
	Image             string                     `json:"image"  dc:"图片"`
	Uuid              string                     `json:"uuid"  dc:"编号"`
	Cost              float64                    `json:"cost"  dc:"成本"`
	BarCode           string                     `json:"barCode"  dc:"条码"`
	OriginalPrice     float64                    `json:"originalPrice"  dc:"原价"`
	VipPrice          float64                    `json:"vipPrice"  dc:"会员价"`
	Weight            float64                    `json:"weight"  dc:"重量(kg)"`
	Volume            float64                    `json:"volume"  dc:"体积(m3)"`
	IsVirtual         int                        `json:"isVirtual"`   // 是否是虚拟商品
	VirtualType       int                        `json:"virtualType"` // 虚拟商品类型（1、固定卡密；2、一次性卡密）
	VirtualInfo       string                     `json:"virtualInfo"` // 虚拟商品内容
	CouponId          uint64                     `json:"couponId"`    // 关联优惠券
	CardPasswordGoods *StoreGoodsPasswordCardReq `json:"cardPasswordGoods" dc:"卡密商品"`
	Type              int                        `json:"type" dc:"活动类型"` // 0=商品，1=秒杀，2=砍价，3=拼团
	//DeletedAt         string                     `json:"deletedAt" dc:"删除时间"`
	Points float64 `json:"points" dc:"购买奖励积分"`
}

type StoreGoodsEditReq struct {
	Id                  uint64                        `json:"id" v:"required#主键ID不能为空" dc:"ID"`
	Name                string                        `json:"name" v:"required#商品名称不能为空" dc:"商品名称"`
	Summary             string                        `json:"summary"  dc:"商品描述"`
	Imgs                StringArray                   `json:"imgs"  dc:"封面图"`
	Thumb               string                        `json:"thumb"  dc:"缩略图"`
	Keywords            string                        `json:"keywords"  dc:"关键字"`
	CateId              Uint64Array                   `json:"cateId"  dc:"分类"`
	UnitName            string                        `json:"unitName" v:"required#单位不能为空" dc:"单位"`
	Sort                int                           `json:"sort"  dc:"排序"`
	IsShow              int                           `json:"isShow"  dc:"是否显示"`
	IsHot               int                           `json:"isHot"  dc:"是否热卖"`
	IsBest              int                           `json:"isBest"  dc:"是否精品"`
	IsNew               int                           `json:"isNew"  dc:"是否新品"`
	SpecificationType   int                           `json:"specificationType" dc:"规格"`
	GoodsType           int                           `json:"goodsType"  dc:"虚拟商品类型"`
	IsPostage           int                           `json:"isPostage"  dc:"是否包邮"`
	GiftPoints          uint                          `json:"giftPoints"  dc:"赠送积分"`
	Hits                int64                         `json:"hits"  dc:"点击量"`
	VideoUrl            string                        `json:"videoUrl"  dc:"视频地址"`
	FreightTempId       int                           `json:"freightTempId"  dc:"运费模板"`
	Spu                 string                        `json:"spu"  dc:"商品SPU"`
	LabelId             Uint64Array                   `json:"labelId"  dc:"商品标签ID"`
	IsVip               int                           `json:"isVip"  dc:"是否VIP商品"`
	Freight             float64                       `json:"freight"  dc:"运费"`
	Logistics           Uint64Array                   `json:"logistics"  dc:"物流方式"`
	MinOrderNum         uint                          `json:"minOrderNum"  dc:"起购数量"`
	MaxOrderNum         uint                          `json:"maxOrderNum"  dc:"限购数量"`
	PurchaseLimited     uint                          `json:"purchaseLimited"  dc:"是否限购"`
	PurchaseLimitedType uint                          `json:"purchaseLimitedType"  dc:"限购类型，0:单次限购，1:单人限购"`
	Presell             uint                          `json:"presell"  dc:"是否预售"`
	PresellBeginTime    *time.Time                    `json:"presellBeginTime"  dc:"预售开始时间"`
	PresellEndTime      *time.Time                    `json:"presellEndTime"  dc:"预售结束时间"`
	PresellDeliveryTime *time.Time                    `json:"presellDeliveryTime"  dc:"预售结束发货时间"`
	Attr                string                        `json:"attr" dc:"商品属性json串"`
	AttrValues          []*StoreGoodsAttrValueEditReq `json:"attrValues" dc:"多规格商品属性值"`
	AttrValue           *StoreGoodsAttrValueEditReq   `json:"attrValue" dc:"单规格商品属性值"`
	IsoCode             string                        `json:"isoCode" dc:"商品所属国家"`
}

type Web3StoreGoodsListingEventData struct {
	TxHash  string `json:"txHash"`
	GoodsId uint64 `json:"goodsId"`
	Sender  string `json:"sender"`
}

func (s Web3StoreGoodsListingEventData) String() string {
	res, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(res)
}

type Web3StoreGoodsListingApproveEventData struct {
	TxHash  string `json:"txHash"`
	GoodsId uint64 `json:"goodsId"`
	Type    int    `json:"type"` // 1: approve1; 2: approve2; 3: approve3
	Sender  string `json:"sender"`
}

func (s Web3StoreGoodsListingApproveEventData) String() string {
	res, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(res)
}
